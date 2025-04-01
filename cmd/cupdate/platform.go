package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"path/filepath"

	"github.com/AlexGustafsson/cupdate/internal/configutils"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	"github.com/AlexGustafsson/cupdate/internal/platform/docker"
	"github.com/AlexGustafsson/cupdate/internal/platform/kubernetes"
)

func ConfigurePlatformGrapher(ctx context.Context, config *Config) (platform.ContinuousGrapher, error) {
	// Set up the configured platform (Docker if specified, auto discovery of
	// Kubernetes otherwise)
	if len(config.Docker.Hosts) == 0 {
		kubernetesConfig, err := config.KubernetesClientConfig()
		if err != nil {
			return nil, err
		}

		platform, err := kubernetes.NewPlatform(kubernetesConfig, &kubernetes.Options{IncludeOldReplicaSets: config.Kubernetes.IncludeOldReplicaSets})
		if err != nil {
			return nil, err
		}

		return platform, nil
	}

	graphers := make([]platform.Grapher, 0)
	for _, host := range config.Docker.Hosts {
		options := &docker.Options{
			IncludeAllContainers: config.Docker.IncludeAllContainers,
		}

		if config.Docker.TLSPath != "" {
			uri, err := url.Parse(host)
			if err != nil {
				return nil, fmt.Errorf("failed to parse docker host '%s': %w", host, err)
			}

			tlsConfig, err := configutils.LoadTLSConfig(
				filepath.Join(config.Docker.TLSPath, uri.Hostname()),
				config.Docker.TLSPath,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to read docker host '%s' TLS files: %w", host, err)
			}

			options.TLSClientConfig = tlsConfig
		}

		platform, err := docker.NewPlatform(ctx, host, options)
		if err != nil {
			return nil, err
		}

		graphers = append(graphers, platform)
	}

	slog.Debug("Platform lacks native continuous graphing support. Falling back to polling", slog.Duration("interval", config.Processing.Interval))
	return &platform.PollGrapher{
		Grapher: &platform.CompoundGrapher{
			Graphers: graphers,
		},
		Interval: config.Processing.Interval,
	}, nil
}
