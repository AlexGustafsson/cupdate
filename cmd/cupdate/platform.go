package main

import (
	"context"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"

	"github.com/AlexGustafsson/cupdate/internal/configutils"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	"github.com/AlexGustafsson/cupdate/internal/platform/docker"
	"github.com/AlexGustafsson/cupdate/internal/platform/kubernetes"
	"github.com/AlexGustafsson/cupdate/internal/platform/static"
)

func ConfigurePlatform(ctx context.Context, config *Config) (platform.ContinuousGrapher, error) {
	// 1. Docker
	if len(config.Docker.Hosts) > 0 {
		graphers := make([]platform.Grapher, 0)
		for _, host := range config.Docker.Hosts {
			options := &docker.Options{
				IncludeAllContainers: config.Docker.IncludeAllContainers,
			}

			if config.Docker.TLSPath != "" {
				uri, err := url.Parse(host)
				if err != nil {
					slog.ErrorContext(ctx, "Failed to parse docker URI", slog.Any("error", err), slog.String("host", host))
					os.Exit(1)
				}

				tlsConfig, err := configutils.LoadTLSConfig(
					filepath.Join(config.Docker.TLSPath, uri.Hostname()),
					config.Docker.TLSPath,
				)
				if err != nil {
					slog.ErrorContext(ctx, "Failed to read docker TLS files", slog.Any("error", err), slog.String("host", host))
				}

				options.TLSClientConfig = tlsConfig
			}

			platform, err := docker.NewPlatform(ctx, host, options)
			if err != nil {
				slog.ErrorContext(ctx, "Failed to create docker source", slog.Any("error", err), slog.String("host", host))
				os.Exit(1)
			}

			graphers = append(graphers, platform)
		}

		return platform.NewPollGrapher(
			&platform.CompoundGrapher{
				Graphers: graphers,
			},
			config.Processing.Interval,
		), nil
	}

	// 2. Static platform
	if config.Static.FilePath != "" {
		return platform.NewPollGrapher(
			&static.Platform{
				FilePath: config.Static.FilePath,
			},
			config.Processing.Interval,
		), nil
	}

	// 3. Default to Kubernetes auto-discovery
	kubernetesConfig, err := config.KubernetesClientConfig()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to configure Kubernetes client", slog.Any("error", err))
		os.Exit(1)
	}

	options := &kubernetes.Options{
		IncludeOldReplicaSets: config.Kubernetes.IncludeOldReplicaSets,
		DebounceInterval:      config.Kubernetes.DebounceInterval,
	}
	return kubernetes.NewPlatform(kubernetesConfig, options)
}
