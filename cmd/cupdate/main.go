package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/events"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	"github.com/AlexGustafsson/cupdate/internal/store"
	"github.com/AlexGustafsson/cupdate/internal/syncutil"
	"github.com/AlexGustafsson/cupdate/internal/worker"
	"github.com/prometheus/client_golang/prometheus"
)

type ExitCode int

const (
	ExitCodeOK    ExitCode = 0
	ExitCodeNotOK ExitCode = 1
)

// setup configures base resources and returns them.
// If any resource failes to initialize, those that succeeded are still
// returned. The caller should therefore make sure to clean up the initialized
// resources upon receiving an error.
func setup(ctx context.Context, config *Config) (*cache.DiskCache, *store.Store, *store.Store, platform.ContinuousGrapher, error) {
	if err := ConfigureOtel(ctx, config); err != nil {
		slog.ErrorContext(ctx, "Failed to initialize otel", slog.Any("error", err))
		return nil, nil, nil, nil, err
	}

	cache, err := cache.NewDiskCache(config.Cache.Path)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create disk cache", slog.Any("error", err))
		return nil, nil, nil, nil, err
	}
	prometheus.DefaultRegisterer.MustRegister(cache)

	if err := store.Initialize(ctx, config.DatabaseURI()); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	readStore, err := store.New(ctx, config.DatabaseURI(), true)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to configure read database: %w", err)
	}

	writeStore, err := store.New(ctx, config.DatabaseURI(), false)
	if err != nil {
		return cache, readStore, nil, nil, fmt.Errorf("failed to configure write database: %w", err)
	}

	targetPlatform, err := ConfigurePlatform(ctx, config)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to configure platform", slog.Any("error", err))
		return cache, readStore, writeStore, nil, err
	}

	return cache, readStore, writeStore, targetPlatform, nil
}

func run(environ []string, signals <-chan os.Signal) ExitCode {
	// Configure logging once with defaults, before any configuration is available
	// and then again once configuration is available
	ConfigureLogging(nil)
	config, err := ParseConfigFromEnv(environ)
	if err != nil {
		slog.Error("Failed to parse config from environment variables", slog.Any("error", err))
		return ExitCodeNotOK
	}
	ConfigureLogging(config)
	slog.Debug("Parsed config", slog.Any("config", config))

	setupCtx, cancelSetup := context.WithTimeout(context.Background(), 60*time.Second)
	cache, readStore, writeStore, targetPlatform, startErr := setup(setupCtx, config)
	cancelSetup()
	if startErr != nil {
		slog.Error("Failed to set up resources", slog.Any("error", err))
		// Fallthrough - clean up resources
	}

	var wg sync.WaitGroup
	runCtx, cancelRun := context.WithCancel(context.Background())
	if startErr == nil {
		processQueue := worker.NewQueue[oci.Reference](config.Processing.QueueBurst, config.Processing.QueueRate)
		prometheus.DefaultRegisterer.MustRegister(processQueue)

		httpClient := httputil.NewClient(cache, config.Cache.MaxAge)
		httpClient.UserAgent = config.HTTP.UserAgent
		prometheus.DefaultRegisterer.MustRegister(httpClient)

		worker := worker.New(httpClient, writeStore, config.RegistryAuth())
		prometheus.DefaultRegisterer.MustRegister(worker)

		platformHub := events.NewHub[models.PlatformEvent]()

		httpServer := ConfigureServer(config, httpClient, readStore, worker.Hub, platformHub, processQueue, targetPlatform)

		wg.Go(func() {
			HandleScheduling(runCtx, config, processQueue, readStore)
			processQueue.Close()
		})

		wg.Go(func() {
			HandleProcessing(runCtx, config, worker, processQueue)
		})

		wg.Go(func() {
			HandleGraphs(runCtx, targetPlatform, platformHub, writeStore, processQueue)
		})

		wg.Go(func() {
			HandleWorkflowCleanup(runCtx, config, writeStore)
		})

		// Close the HTTP signal on shutdown
		go func() {
			<-signals
			slog.Info("Caught signal, exiting gracefully")

			shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 15*time.Second)
			slog.Debug("Shutting down HTTP server")
			err := httpServer.Shutdown(shutdownCtx)
			cancelShutdown()
			if err != nil {
				slog.Error("Failed to shut down HTTP server", slog.Any("error", err))
				return
			}
		}()

		slog.InfoContext(runCtx, "Starting HTTP server")
		err = httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(runCtx, "Failed to serve", slog.Any("error", err))
			// Fallthrough
		}
	}

	slog.Debug("Canceling processing")
	cancelRun()

	// Close ungracefully on another signal
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 60*time.Second)
	go func() {
		<-signals

		slog.WarnContext(runCtx, "Caught signal, exiting ungracefully immediately")
		cancelShutdown()
	}()

	slog.Debug("Waiting for processing to end")
	if err := syncutil.WaitContext(shutdownCtx, &wg); err != nil {
		slog.Error("Failed to wait for processing to end", slog.Any("error", err))
		return ExitCodeNotOK
	}

	// Gracefully shut down resources

	if cache != nil {
		wg.Go(func() {
			slog.Debug("Closing cache")
			err := cache.Close()
			if err != nil {
				slog.Error("Failed to close cache", slog.Any("error", err))
				return
			}
		})
	}

	if readStore != nil {
		wg.Go(func() {
			slog.Debug("Closing read store")
			err := readStore.Close()
			if err != nil {
				slog.Error("Failed to close read store", slog.Any("error", err))
				return
			}
		})
	}

	if writeStore != nil {
		wg.Go(func() {
			slog.Debug("Closing write store")
			err := writeStore.Close()
			if err != nil {
				slog.Error("Failed to close read store", slog.Any("error", err))
				return
			}
		})
	}

	if targetPlatform != nil {
		wg.Go(func() {
			slog.Debug("Closing platform")
			err := targetPlatform.Close()
			if err != nil {
				slog.Error("Failed to close platform", slog.Any("error", err))
				return
			}
		})
	}

	slog.Debug("Waiting for resources to close")
	if err := syncutil.WaitContext(shutdownCtx, &wg); err != nil {
		slog.Error("Failed to wait for resources to close", slog.Any("error", err))
		return ExitCodeNotOK
	}

	slog.Info("Shut down successfully")
	if startErr == nil {
		return ExitCodeOK
	} else {
		return ExitCodeNotOK
	}
}

func main() {
	// NOTE: In order to allow for testing and use of defer, think twice before
	// adding logic to the main function
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	exitCode := run(os.Environ(), signals)
	os.Exit(int(exitCode))
}
