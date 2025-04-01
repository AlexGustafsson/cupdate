package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/platform/docker"
	"github.com/caarlos0/env/v11"
	"k8s.io/client-go/rest"
)

type Config struct {
	Log struct {
		Level string `env:"LEVEL" envDefault:"info"`
	} `envPrefix:"LOG_"`

	API struct {
		Address string `env:"ADDRESS" envDefault:"0.0.0.0"`
		Port    uint16 `env:"PORT" envDefault:"8080"`
	} `envPrefix:"API_"`

	Web struct {
		Disabled bool   `env:"DISABLED"`
		Address  string `env:"ADDRESS"`
	} `envPrefix:"WEB_"`

	HTTP struct {
		UserAgent string `env:"USER_AGENT" envDefault:"Cupdate/1.0"`
	} `envPrefix:"HTTP_"`

	Cache struct {
		Path   string        `env:"PATH" envDefault:"cachev1.boltdb"`
		MaxAge time.Duration `env:"MAX_AGE" envDefault:"24h"`
	} `envPrefix:"CACHE_"`

	Database struct {
		Path string `env:"PATH" envDefault:"dbv1.sqlite"`
	} `envPrefix:"DB_"`

	Processing struct {
		Interval   time.Duration `env:"INTERVAL" envDefault:"1h"`
		Items      int           `env:"ITEMS" envDefault:"10"`
		MinAge     time.Duration `env:"MIN_AGE" envDefault:"72h"`
		Timeout    time.Duration `env:"TIMEOUT" envDefault:"2m"`
		QueueSize  int           `env:"QUEUE_SIZE" envDefault:"50"`
		QueueBurst int           `env:"QUEUE_BURST" envDefault:"10"`
		QueueRate  time.Duration `env:"QUEUE_RATE" envDefault:"1m"`
	} `envPrefix:"PROCESSING_"`

	Workflow struct {
		CleanupMaxAge   time.Duration `env:"CLEANUP_MAX_AGE" envDefault:"48h"`
		CleanupInterval time.Duration `env:"CLEANUP_INTERVAL" envDefault:"1h"`
	} `envPrefix:"WORKFLOW_"`

	Kubernetes struct {
		Host                  string `env:"HOST"`
		IncludeOldReplicaSets bool   `env:"INCLUDE_OLD_REPLICAS"`
	} `envPrefix:"KUBERNETES_"`

	Docker struct {
		Hosts                []string `env:"HOST"`
		IncludeAllContainers bool     `env:"INCLUDE_ALL_CONTAINERS"`
		TLSPath              string   `env:"TLS_PATH"`
	} `envPrefix:"DOCKER_"`

	OTEL struct {
		Target   string `env:"TARGET"`
		Insecure bool   `env:"INSECURE"`
	} `envPrefix:"OTEL_"`

	Registry struct {
		Secrets string `env:"SECRETS"`
	} `envPrefix:"REGISTRY_"`
}

func (c *Config) LogLevel() (slog.Level, error) {
	var logLevel slog.Level

	switch c.Log.Level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		return logLevel, fmt.Errorf("invalid log level")
	}

	return logLevel, nil
}

func (c *Config) RegistryAuth() (*httputil.AuthMux, error) {
	registryAuth := httputil.NewAuthMux()

	if c.Registry.Secrets != "" {
		file, err := os.Open(c.Registry.Secrets)
		if err != nil {
			return nil, fmt.Errorf("failed to read registry secrets: %w", err)
		}

		var dockerConfig *docker.ConfigFile
		err = json.NewDecoder(file).Decode(&dockerConfig)
		file.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to parse registry secrets: %w", err)
		}

		for k, v := range dockerConfig.HttpHeaders {
			registryAuth.SetHeader(k, v)
		}

		for pattern, auth := range dockerConfig.Auths {
			if auth.Auth == "" {
				registryAuth.Handle(pattern, httputil.BasicAuthHandler{
					Username: auth.Username,
					Password: auth.Password,
				})
			} else {
				value, err := base64.StdEncoding.DecodeString(auth.Auth)
				if err != nil {
					return nil, fmt.Errorf("invalid registry secrets file: %w", err)
				}

				username, password, ok := strings.Cut(string(value), ":")
				if !ok {
					return nil, fmt.Errorf("invalid registry secrets file: invalid auth field")
				}

				registryAuth.Handle(pattern, httputil.BasicAuthHandler{
					Username: username,
					Password: password,
				})
			}
		}
	}

	return registryAuth, nil
}

func (c *Config) KubernetesClientConfig() (*rest.Config, error) {
	if c.Kubernetes.Host == "" {
		return rest.InClusterConfig()
	}

	return &rest.Config{
		Host: c.Kubernetes.Host,
	}, nil
}

func (c *Config) DatabaseURI() (string, error) {
	absolutePath, err := filepath.Abs(c.Database.Path)
	if err != nil {
		return "", err
	}

	return "file://" + absolutePath, nil
}

func ParseConfigFromEnv() (*Config, error) {
	var config Config
	if err := env.ParseWithOptions(&config, env.Options{Prefix: "CUPDATE_"}); err != nil {
		return nil, err
	}

	return &config, nil
}
