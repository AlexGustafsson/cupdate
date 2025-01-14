package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/api"
	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/otelutil"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	"github.com/AlexGustafsson/cupdate/internal/platform/docker"
	"github.com/AlexGustafsson/cupdate/internal/platform/kubernetes"
	"github.com/AlexGustafsson/cupdate/internal/ratelimit"
	"github.com/AlexGustafsson/cupdate/internal/slogutil"
	"github.com/AlexGustafsson/cupdate/internal/store"
	"github.com/AlexGustafsson/cupdate/internal/web"
	"github.com/AlexGustafsson/cupdate/internal/worker"
	"github.com/caarlos0/env/v10"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/rest"
)

// Version is the current Cupdate version. Overwritten at build time.
var Version string = "development build"

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

	Kubernetes struct {
		Host                  string `env:"HOST"`
		IncludeOldReplicaSets bool   `env:"INCLUDE_OLD_REPLICAS"`
	} `envPrefix:"KUBERNETES_"`

	Docker struct {
		Hosts                []string `env:"HOST"`
		IncludeAllContainers bool     `env:"INCLUDE_ALL_CONTAINERS"`
	} `envPrefix:"DOCKER_"`

	OTEL struct {
		Target   string `env:"TARGET"`
		Insecure bool   `env:"INSECURE"`
	} `envPrefix:"OTEL_"`

	Registry struct {
		Secrets string `env:"SECRETS"`
	} `envPrefix:"REGISTRY_"`
}

func main() {
	slog.SetDefault(slog.New(slogutil.NewHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})).With(slog.String("service.version", Version)).With(slog.String("service.name", "cupdate")))

	var config Config
	err := env.ParseWithOptions(&config, env.Options{
		Prefix: "CUPDATE_",
	})
	if err != nil {
		slog.Error("Failed to parse config from environment variables", slog.Any("error", err))
		os.Exit(1)
	}

	var logLevel slog.Level
	switch config.Log.Level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		slog.Error("Failed to parse config - invalid log level")
		os.Exit(1)
	}
	slog.SetDefault(slog.New(slogutil.NewHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})).With(slog.String("service.version", Version)).With(slog.String("service.name", "cupdate")))

	slog.Debug("Parsed config", slog.Any("config", config))

	registryAuth := httputil.NewAuthMux()
	if config.Registry.Secrets != "" {
		file, err := os.Open(config.Registry.Secrets)
		if err != nil {
			slog.Error("Failed to read registry secrets", slog.Any("error", err))
			os.Exit(1)
		}

		var dockerConfig *docker.ConfigFile
		err = json.NewDecoder(file).Decode(&dockerConfig)
		file.Close()
		if err != nil {
			slog.Error("Failed to parse registry secrets", slog.Any("error", err))
			os.Exit(1)
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
				registryAuth.Handle(pattern, httputil.BearerToken(auth.Auth))
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	if config.OTEL.Target != "" {
		shutdown, err := otelutil.Init(ctx, config.OTEL.Target, config.OTEL.Insecure)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to initialize otel", slog.Any("error", err))
			os.Exit(1)
		}

		// TODO: This won't be invoked on exit, make it part of the shutdown
		// procedure
		defer shutdown(ctx)
	}

	// Set up the configured platform (Docker if specified, auto discovery of
	// Kubernetes otherwise)
	var targetPlatform platform.Grapher
	if len(config.Docker.Hosts) == 0 {
		var kubernetesConfig *rest.Config
		if config.Kubernetes.Host == "" {
			var err error
			kubernetesConfig, err = rest.InClusterConfig()
			if err != nil {
				slog.ErrorContext(ctx, "Failed to configure Kubernetes client", slog.Any("error", err))
				os.Exit(1)
			}
		} else {
			kubernetesConfig = &rest.Config{
				Host: config.Kubernetes.Host,
			}
		}

		targetPlatform, err = kubernetes.NewPlatform(kubernetesConfig, &kubernetes.Options{IncludeOldReplicaSets: config.Kubernetes.IncludeOldReplicaSets})
		if err != nil {
			slog.ErrorContext(ctx, "Failed to create kubernetes source", slog.Any("error", err))
			os.Exit(1)
		}
	} else {
		graphers := make([]platform.Grapher, 0)
		for _, host := range config.Docker.Hosts {
			platform, err := docker.NewPlatform(context.Background(), host, &docker.Options{
				IncludeAllContainers: config.Docker.IncludeAllContainers,
			})
			if err != nil {
				slog.ErrorContext(ctx, "Failed to create docker source", slog.Any("error", err))
				os.Exit(1)
			}

			graphers = append(graphers, platform)
		}
		targetPlatform = &platform.CompoundGrapher{
			Graphers: graphers,
		}
	}

	cache, err := cache.NewDiskCache(config.Cache.Path)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create disk cache", slog.Any("error", err))
		os.Exit(1)
	}
	prometheus.DefaultRegisterer.MustRegister(cache)

	absoluteDatabasePath, err := filepath.Abs(config.Database.Path)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to resolve database path", slog.Any("error", err))
		os.Exit(1)
	}

	readStore, err := store.New("file://"+absoluteDatabasePath, true)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load database", slog.Any("error", err))
		os.Exit(1)
	}
	writeStore, err := store.New("file://"+absoluteDatabasePath, false)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load database", slog.Any("error", err))
		os.Exit(1)
	}

	var wg errgroup.Group

	processQueue := make(chan oci.Reference, config.Processing.QueueSize)

	wg.Go(func() error {
		ticker := time.NewTicker(config.Processing.Interval)
		defer ticker.Stop()
		defer close(processQueue)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(ctx, 30*time.Second)

				slog.DebugContext(ctx, "Identifying old references to process")
				images, err := readStore.ListRawImages(ctx, &store.ListRawImagesOptions{
					NotUpdatedSince: time.Now().Add(-config.Processing.MinAge),
					Limit:           config.Processing.Items,
				})
				if err != nil {
					slog.ErrorContext(ctx, "Failed to process old references", slog.Any("error", err))
					cancel()
					continue
				}

				for _, image := range images {
					reference, err := oci.ParseReference(image.Reference)
					if err != nil {
						slog.ErrorContext(ctx, "Unexpectedly failed to parse reference from store", slog.Any("error", err))
						cancel()
						return err
					}

					select {
					case <-ctx.Done():
						slog.ErrorContext(ctx, "Failed to queue old references")
					case processQueue <- reference:
					}
				}

				cancel()
			}
		}
	})

	wg.Go(func() error {
		httpClient := httputil.NewClient(cache, config.Cache.MaxAge)
		httpClient.UserAgent = config.HTTP.UserAgent
		prometheus.DefaultRegisterer.MustRegister(httpClient)

		worker := worker.New(httpClient, writeStore, registryAuth)
		prometheus.DefaultRegisterer.MustRegister(worker)

		gauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "cupdate",
			Subsystem: "worker",
			Name:      "available_burst",
		})
		prometheus.DefaultRegisterer.MustRegister(gauge)

		for reference, burst := range ratelimit.Channel(ctx, config.Processing.QueueBurst, config.Processing.QueueRate, processQueue) {
			gauge.Set(float64(burst))

			ctx, cancel := context.WithTimeout(ctx, config.Processing.Timeout)
			err := worker.ProcessRawImage(ctx, reference)
			cancel()
			if err != nil {
				slog.ErrorContext(ctx, "Failed to process queued raw image", slog.Any("error", err), slog.String("reference", reference.String()))
			}
		}

		return nil
	})

	wg.Go(func() error {
		slog.InfoContext(ctx, "Starting platform grapher")

		grapher, ok := targetPlatform.(platform.ContinousGrapher)
		if !ok {
			slog.DebugContext(ctx, "Platform lacks native continous graphing support. Falling back to polling", slog.Duration("interval", config.Processing.Interval))
			grapher = &platform.PollGrapher{
				Grapher:  targetPlatform,
				Interval: config.Processing.Interval,
			}
		}

		graphs, err := grapher.GraphContinously(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to start graphing platform", slog.Any("error", err))
			return err
		}

		for graph := range graphs {
			slog.DebugContext(ctx, "Got updated platform graph")
			roots := graph.Roots()

			for _, root := range roots {
				imageNode := root.(platform.ImageNode)

				subgraph := graph.Subgraph(root.ID())

				edges := subgraph.Edges()
				nodes := subgraph.Nodes()

				var namespaceNode *platform.Node

				mappedNodes := make(map[string]models.GraphNode)
				for _, node := range nodes {
					switch n := node.(type) {
					case kubernetes.Resource:
						mappedNodes[node.ID()] = models.GraphNode{
							Domain: "kubernetes",
							Type:   string(n.Kind()),
							Name:   n.Name(),
						}
						if node.Type() == "kubernetes/"+kubernetes.ResourceKindCoreV1Namespace {
							namespaceNode = &node
						}
					case docker.Resource:
						mappedNodes[node.ID()] = models.GraphNode{
							Domain: "docker",
							Type:   string(n.Kind()),
							Name:   n.Name(),
						}
						if node.Type() == "docker/"+docker.ResourceKindSwarmNamespace || node.Type() == "docker/"+docker.ResourceKindComposeProject {
							namespaceNode = &node
						}
					case platform.ImageNode:
						mappedNodes[node.ID()] = models.GraphNode{
							Domain: "oci",
							Type:   "image",
							Name:   imageNode.Reference.String(),
						}
					default:
						panic(fmt.Sprintf("mapping unimplemented node type: %s", node.Type()))
					}
				}

				tags := []string{}

				// Set tags for resources
				if namespaceNode != nil {
					children := edges[(*namespaceNode).ID()]
					for childID, isParent := range children {
						if isParent {
							continue
						}

						var childNode *platform.Node
						for _, node := range nodes {
							var n = node
							if node.ID() == childID {
								childNode = &n
								break
							}
						}

						if childNode != nil {
							switch resource := (*childNode).(type) {
							case kubernetes.Resource:
								kind := resource.Kind()
								if kind.IsSupported() {
									tags = append(tags, kubernetes.TagName(resource.Kind()))
								}
							case docker.Resource:
								tags = append(tags, docker.TagName(resource.Kind()))
							}
						}
					}
				}

				mappedGraph := models.Graph{
					Edges: edges,
					Nodes: mappedNodes,
				}

				rawImage := &models.RawImage{
					Reference: imageNode.Reference.String(),
					Tags:      tags,
					Graph:     mappedGraph,
				}

				// TODO: Do this inside of the worker as well?
				slog.DebugContext(ctx, "Inserting raw image", slog.String("reference", rawImage.Reference))
				inserted, err := writeStore.InsertRawImage(context.TODO(), rawImage)
				if err != nil {
					slog.ErrorContext(ctx, "Failed to insert raw image", slog.Any("error", err))
					return err
				}

				// Try to schedule the image for processing
				if inserted {
					slog.DebugContext(ctx, "Raw image inserted for first time - scheduling for processing")
					select {
					case processQueue <- imageNode.Reference:
					default:
						slog.WarnContext(ctx, "Failed to schedule inserted raw image for processing - queue full")
					}
				}
			}

			allReferences := make([]string, 0)
			for _, root := range roots {
				imageNode := root.(platform.ImageNode)
				allReferences = append(allReferences, imageNode.Reference.String())
			}

			slog.DebugContext(ctx, "Cleaning up removed images")
			removed, err := writeStore.DeleteNonPresent(context.TODO(), allReferences)
			if err != nil {
				slog.ErrorContext(ctx, "Failed to clean up removed images", slog.Any("error", err))
				return err
			}
			slog.DebugContext(ctx, "Cleaned up removed images successfully", slog.Int64("removed", removed))
		}

		return nil
	})

	mux := http.NewServeMux()

	apiServer := api.NewServer(readStore, writeStore.Hub, processQueue)
	apiServer.WebAddress = config.Web.Address
	mux.Handle("/api/v1/", apiServer)
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/livez", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	mux.Handle("/readyz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Figure out what checks to have
		w.WriteHeader(http.StatusOK)
	}))

	if !config.Web.Disabled {
		mux.Handle("/", web.MustNewServer())
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.API.Address, config.API.Port),
		Handler: mux,
	}

	wg.Go(func() error {
		slog.InfoContext(ctx, "Starting HTTP server")
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "Failed to serve", slog.Any("error", err))
			return err
		}
		return nil
	})

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	caught := 0
	go func() {
		for range signals {
			caught++
			if caught == 1 {
				slog.InfoContext(ctx, "Caught signal, exiting gracefully")
				if err := httpServer.Close(); err != nil {
					slog.ErrorContext(ctx, "Failed to close server", slog.Any("error", err))
					// Fallthrough
				}
				if err := cache.Close(); err != nil {
					slog.ErrorContext(ctx, "Failed to close cache", slog.Any("error", err))
					// Fallthrough
				}
				if err := readStore.Close(); err != nil {
					slog.ErrorContext(ctx, "Failed to close read store", slog.Any("error", err))
					// Fallthrough
				}
				if err := writeStore.Close(); err != nil {
					slog.ErrorContext(ctx, "Failed to close write store", slog.Any("error", err))
					// Fallthrough
				}
				// Cancel goroutines started in main last as to block on all of the
				// above calls
				cancel()
			} else {
				slog.InfoContext(ctx, "Caught signal, exiting now")
				os.Exit(1)
			}
		}
	}()

	if err := wg.Wait(); err != nil && err != ctx.Err() {
		slog.ErrorContext(ctx, "Failed to run", slog.Any("error", err))
		os.Exit(1)
	}
}
