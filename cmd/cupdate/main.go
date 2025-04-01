package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/api"
	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/otelutil"
	"github.com/AlexGustafsson/cupdate/internal/slogutil"
	"github.com/AlexGustafsson/cupdate/internal/store"
	"github.com/AlexGustafsson/cupdate/internal/web"
	"github.com/AlexGustafsson/cupdate/internal/worker"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ErrNonZeroExitCode is a special error which, when caught in main, will exit
// with a non-zero exit code without logging the error, assuming it was already
// handled.
var ErrNonZeroExitCode = errors.New("unexpected error")

func main() {
	InitDefaultLogger()

	runCtx, cancelRunCtx := context.WithCancel(context.Background())

	var mutex sync.Mutex
	shutdownFuncs := []func(context.Context){}
	registerShutdownFunc := func(f func(context.Context)) {
		mutex.Lock()
		defer mutex.Unlock()

		shutdownFuncs = append(shutdownFuncs, f)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	caught := 0
	go func() {
		for range signals {
			caught++
			if caught == 1 {
				slog.Info("Caught signal, exiting gracefully")

				shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), 15*time.Second)
				mutex.Lock()
				for _, shutdownFunc := range shutdownFuncs {
					shutdownFunc(shutdownCtx)
				}
				cancelShutdownCtx()

				// Cancel the run context last as to make sure the cleanup funcs were
				// run
				cancelRunCtx()
			} else {
				slog.Warn("Caught signal, exiting now")
				os.Exit(1)
			}
		}
	}()

	if err := run(runCtx, registerShutdownFunc); err != nil {
		if err != ErrNonZeroExitCode {
			slog.Error("Exiting with non-zero exit code", slog.Any("error", err))
		}
		os.Exit(1)
	}
}

func run(ctx context.Context, registerShutdownFunc func(func(context.Context))) error {
	slog.SetDefault(slog.New(slogutil.NewHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})).With(slog.String("service.version", Version)).With(slog.String("service.name", "cupdate")))

	config, err := ParseConfigFromEnv()
	if err != nil {
		slog.Error("Failed to parse config from environment variables", slog.Any("error", err))
		return ErrNonZeroExitCode
	}

	logLevel, err := config.LogLevel()
	if err != nil {
		slog.Error("Failed to parse config - invalid log level")
		return ErrNonZeroExitCode
	}
	slog.SetLogLoggerLevel(logLevel)

	slog.Debug("Parsed config", slog.Any("config", config))

	registryAuth, err := config.RegistryAuth()
	if err != nil {
		slog.Error("Failed to parse registry auth", slog.Any("error", err))
		return ErrNonZeroExitCode
	}

	if config.OTEL.Target != "" {
		shutdown, err := otelutil.Init(ctx, config.OTEL.Target, config.OTEL.Insecure)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to initialize otel", slog.Any("error", err))
			return ErrNonZeroExitCode
		}

		defer shutdown(ctx)
	}

	targetPlatform, err := ConfigurePlatformGrapher(ctx, config)
	if err != nil {
		slog.Error("Failed to configure platform", slog.Any("error", err))
		return ErrNonZeroExitCode
	}

	cache, err := cache.NewDiskCache(config.Cache.Path)
	if err != nil {
		slog.Error("Failed to create disk cache", slog.Any("error", err))
		return ErrNonZeroExitCode
	}
	defer func() {
		if err := cache.Close(); err != nil {
			slog.Error("Failed to close cache", slog.Any("error", err))
		}
	}()
	if err := prometheus.DefaultRegisterer.Register(cache); err != nil {
		slog.Error("Failed to register prometheus metrics for cache", slog.Any("error", err))
		return ErrNonZeroExitCode
	}

	databaseURI, err := config.DatabaseURI()
	if err != nil {
		slog.Error("Failed to resolve database path", slog.Any("error", err))
		return ErrNonZeroExitCode
	}

	if err := store.Initialize(ctx, databaseURI); err != nil {
		slog.Error("Failed to initialize database", slog.Any("error", err))
		return ErrNonZeroExitCode
	}

	readStore, err := store.New(databaseURI, true)
	if err != nil {
		slog.Error("Failed to open database for reading", slog.Any("error", err))
		return ErrNonZeroExitCode
	}
	defer func() {
		if err := readStore.Close(); err != nil {
			slog.Error("Failed to close read-only database", slog.Any("error", err))
		}
	}()

	writeStore, err := store.New(databaseURI, false)
	if err != nil {
		slog.Error("Failed to open database for writing", slog.Any("error", err))
		return ErrNonZeroExitCode
	}
	defer func() {
		if err := writeStore.Close(); err != nil {
			slog.Error("Failed to close writable database", slog.Any("error", err))
		}
	}()

	workerQueue := worker.NewQueue[oci.Reference](config.Processing.QueueBurst, config.Processing.QueueRate)
	if err := prometheus.DefaultRegisterer.Register(workerQueue); err != nil {
		slog.Error("Failed to register prometheus metrics for worker queue", slog.Any("error", err))
		return ErrNonZeroExitCode
	}

	imageWorkerScheduler := worker.NewImageWorkerScheduler(readStore)

	httpClient := httputil.NewClient(cache, config.Cache.MaxAge)
	httpClient.UserAgent = config.HTTP.UserAgent
	if err := prometheus.DefaultRegisterer.Register(httpClient); err != nil {
		slog.Error("Failed to register prometheus metrics for HTTP client", slog.Any("error", err))
		return ErrNonZeroExitCode
	}

	imageWorker := worker.NewImageWorker(httpClient, writeStore, registryAuth)
	if err := prometheus.DefaultRegisterer.Register(imageWorker); err != nil {
		slog.Error("Failed to register prometheus metrics for worker", slog.Any("error", err))
		return ErrNonZeroExitCode
	}

	serveMux := http.NewServeMux()

	apiServer := api.NewServer(readStore, imageWorker.Hub, workerQueue)
	apiServer.WebAddress = config.Web.Address
	serveMux.Handle("/api/v1/", apiServer)
	serveMux.Handle("/metrics", promhttp.Handler())
	serveMux.Handle("/livez", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	serveMux.Handle("/readyz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Figure out what checks to have
		w.WriteHeader(http.StatusOK)
	}))

	if !config.Web.Disabled {
		serveMux.Handle("/", web.MustNewEmbeddedServer())
	}

	httpServer := &http.Server{
		Addr: fmt.Sprintf("%s:%d", config.API.Address, config.API.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := w
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.Header().Set("Content-Encoding", "gzip")
				gzip := &httputil.GzipWriter{ResponseWriter: w}
				defer gzip.Close()

				writer = gzip
			}

			serveMux.ServeHTTP(writer, r)
		}),
	}
	registerShutdownFunc(func(ctx context.Context) {
		if err := httpServer.Shutdown(ctx); err != nil {
			slog.Error("Failed to shutdown HTTP server", slog.Any("error", err))
			// Fallthrough
		}
	})

	var wg sync.WaitGroup

	// Run the worker scheduler to push jobs to the worker queue
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer workerQueue.Close()

		imageWorkerScheduler.PushTo(ctx, workerQueue, config.Processing.Interval, config.Processing.MinAge, config.Processing.Items)
	}()

	// Run the worker to handle pull jobs from the worker queue
	wg.Add(1)
	go func() {
		defer wg.Done()

		imageWorker.PullFrom(ctx, workerQueue, config.Processing.Timeout)
	}()

	// Keep available images up-to-date by reacting on changes made to the
	// platform
	slog.DebugContext(ctx, "Starting platform grapher")
	wg.Add(1)
	go func() {
		defer wg.Done()

	}()

	// Start cleaning up on an interval
	wg.Add(1)
	go func() {
		ticker := time.NewTicker(config.Workflow.CleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
				slog.Debug("Cleaning up old workflow runs")
				removed, err := writeStore.DeleteWorkflowRuns(ctx, time.Now().Add(-config.Workflow.CleanupMaxAge))
				cancel()
				if err == nil {
					slog.Debug("Cleaned up old workflow runs successfully", slog.Int64("removed", removed))
				} else {
					slog.Error("Failed to clean up old workflow runs", slog.Any("error", err))
				}
			}
		}
	}()

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to serve", slog.Any("error", err))
			// Fallthrough
		}
	}()

	wg.Wait()
	return ctx.Err()
}
