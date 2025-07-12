package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/api"
	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/configutils"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/otelutil"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	"github.com/AlexGustafsson/cupdate/internal/platform/docker"
	"github.com/AlexGustafsson/cupdate/internal/platform/kubernetes"
	"github.com/AlexGustafsson/cupdate/internal/slogutil"
	"github.com/AlexGustafsson/cupdate/internal/store"
	"github.com/AlexGustafsson/cupdate/internal/web"
	"github.com/AlexGustafsson/cupdate/internal/worker"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

func main() {
	slog.SetDefault(slog.New(slogutil.NewHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})).With(slog.String("service.version", Version)).With(slog.String("service.name", "cupdate")))

	config, err := ParseConfigFromEnv()
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

	registryAuth, err := config.RegistryAuth()
	if err != nil {
		slog.Error("Failed to parse registry auth", slog.Any("error", err))
		os.Exit(1)
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
		kubernetesConfig, err := config.KubernetesClientConfig()
		if err != nil {
			slog.ErrorContext(ctx, "Failed to configure Kubernetes client", slog.Any("error", err))
			os.Exit(1)
		}

		targetPlatform, err = kubernetes.NewPlatform(kubernetesConfig, &kubernetes.Options{IncludeOldReplicaSets: config.Kubernetes.IncludeOldReplicaSets})
		if err != nil {
			slog.ErrorContext(ctx, "Failed to create kubernetes source", slog.Any("error", err))
			os.Exit(1)
		}
	} else {
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

	databaseURI, err := config.DatabaseURI()
	if err != nil {
		slog.Error("Failed to resolve database path", slog.Any("error", err))
		os.Exit(1)
	}

	if err := store.Initialize(ctx, databaseURI); err != nil {
		slog.ErrorContext(ctx, "Failed to initialize database", slog.Any("error", err))
		os.Exit(1)
	}

	readStore, err := store.New(databaseURI, true)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load database", slog.Any("error", err))
		os.Exit(1)
	}
	writeStore, err := store.New(databaseURI, false)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load database", slog.Any("error", err))
		os.Exit(1)
	}

	var wg errgroup.Group

	processQueue := worker.NewQueue[oci.Reference](config.Processing.QueueBurst, config.Processing.QueueRate)
	prometheus.DefaultRegisterer.MustRegister(processQueue)

	wg.Go(func() error {
		ticker := time.NewTicker(config.Processing.Interval)
		defer ticker.Stop()
		defer processQueue.Close()

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
						slog.ErrorContext(ctx, "Unexpectedly failed to parse reference from store", slog.Any("error", err), slog.String("reference", image.Reference))
						cancel()
						return err
					}

					processQueue.Push(reference)
				}

				cancel()
			}
		}
	})

	httpClient := httputil.NewClient(cache, config.Cache.MaxAge)
	httpClient.UserAgent = config.HTTP.UserAgent
	prometheus.DefaultRegisterer.MustRegister(httpClient)

	worker := worker.New(httpClient, writeStore, registryAuth)
	prometheus.DefaultRegisterer.MustRegister(worker)

	wg.Go(func() error {
		for reference := range processQueue.Pull() {
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

		grapher, ok := targetPlatform.(platform.ContinuousGrapher)
		if !ok {
			slog.DebugContext(ctx, "Platform lacks native continuous graphing support. Falling back to polling", slog.Duration("interval", config.Processing.Interval))
			grapher = &platform.PollGrapher{
				Grapher:  targetPlatform,
				Interval: config.Processing.Interval,
			}
		}

		graphs, err := grapher.GraphContinuously(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to start graphing platform", slog.Any("error", err))
			return err
		}

		for graph := range graphs {
			slog.DebugContext(ctx, "Got updated platform graph")

			// Delete ignored images / trees
			graph.DeleteFunc(func(n platform.Node) bool {
				return n.Labels().Ignore()
			})

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
							Domain:         "kubernetes",
							Type:           string(n.Kind()),
							Name:           n.Name(),
							Labels:         n.Labels().RemoveUnsupported(),
							InternalLabels: n.InternalLabels(),
						}
						if node.Type() == "kubernetes/"+kubernetes.ResourceKindCoreV1Namespace {
							namespaceNode = &node
						}
					case docker.Resource:
						mappedNodes[node.ID()] = models.GraphNode{
							Domain:         "docker",
							Type:           string(n.Kind()),
							Name:           n.Name(),
							Labels:         n.Labels().RemoveUnsupported(),
							InternalLabels: n.InternalLabels(),
						}
						if node.Type() == "docker/"+docker.ResourceKindSwarmNamespace || node.Type() == "docker/"+docker.ResourceKindComposeProject {
							namespaceNode = &node
						}
					case platform.ImageNode:
						// This node is added later on
					default:
						panic(fmt.Sprintf("mapping unimplemented node type: %s", node.Type()))
					}
				}

				// Resolve labels for the image node. The nearest label takes precedence
				resolvedLabels := make(map[string]string)
				queue := []string{root.ID()}
				for len(queue) > 0 {
					id := queue[0]
					queue = queue[1:]

					for k, v := range mappedNodes[id].Labels {
						_, ok := resolvedLabels[k]
						if !ok {
							resolvedLabels[k] = v
						}
					}

					for adjacent, isChild := range edges[id] {
						if isChild {
							queue = append(queue, adjacent)
						}
					}
				}
				mappedNodes[root.ID()] = models.GraphNode{
					Domain:         "oci",
					Type:           "image",
					Name:           imageNode.Reference.String(),
					Labels:         resolvedLabels,
					InternalLabels: nil,
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
					processQueue.Push(imageNode.Reference)
				}
			}

			allReferences := make([]string, 0)
			for _, root := range roots {
				imageNode := root.(platform.ImageNode)
				allReferences = append(allReferences, imageNode.Reference.String())
			}

			slog.DebugContext(ctx, "Cleaning up removed images")
			removed, err := writeStore.DeleteNonPresent(context.TODO(), allReferences)
			if err == nil {
				slog.DebugContext(ctx, "Cleaned up removed images successfully", slog.Int64("removed", removed))
			} else {
				slog.ErrorContext(ctx, "Failed to clean up removed images", slog.Any("error", err))
			}
		}

		return nil
	})

	wg.Go(func() error {
		ticker := time.NewTicker(config.Workflow.CleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
				slog.DebugContext(ctx, "Cleaning up old workflow runs")
				removed, err := writeStore.DeleteWorkflowRuns(context.TODO(), time.Now().Add(-config.Workflow.CleanupMaxAge))
				cancel()
				if err == nil {
					slog.DebugContext(ctx, "Cleaned up old workflow runs successfully", slog.Int64("removed", removed))
				} else {
					slog.ErrorContext(ctx, "Failed to clean up old workflow runs", slog.Any("error", err))
				}
			}
		}
	})

	logoProxy := api.CompoundProxy{
		Proxies: []api.LogoProxy{
			&api.LogoFSProxy{
				FS: os.DirFS(config.Logos.Path),
			},
			&api.LogoHTTPProxy{
				Client: httpClient,
				GetURL: readStore.GetImageLogo,
			},
		},
	}

	mux := http.NewServeMux()

	apiServer := api.NewServer(readStore, worker.Hub, processQueue, logoProxy)
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
		mux.Handle("/", web.MustNewEmbeddedServer())
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
