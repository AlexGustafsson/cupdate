package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/api"
	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	"github.com/AlexGustafsson/cupdate/internal/platform/kubernetes"
	"github.com/AlexGustafsson/cupdate/internal/store"
	"github.com/AlexGustafsson/cupdate/internal/worker"
	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/rest"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	config := &rest.Config{
		Host: "http://localhost:8001",
	}

	kubernetesPlatform, err := kubernetes.NewPlatform(config)
	if err != nil {
		slog.Error("Failed to create kubernetes source", slog.Any("error", err))
		os.Exit(1)
	}

	cache, err := cache.NewDiskCache("./cache")
	if err != nil {
		slog.Error("Failed to serve", slog.Any("error", err))
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		slog.Error("Failed to identify working directory", slog.Any("error", err))
		os.Exit(1)
	}

	readStore, err := store.New("file://"+path.Join(cwd, "dbv1.sqlite"), true)
	if err != nil {
		slog.Error("Failed to load database", slog.Any("error", err))
		os.Exit(1)
	}
	writeStore, err := store.New("file://"+path.Join(cwd, "dbv1.sqlite"), false)
	if err != nil {
		slog.Error("Failed to load database", slog.Any("error", err))
		os.Exit(1)
	}

	// Insert default tags
	defaultTags := []models.Tag{
		{
			Name:        "k8s",
			Description: "Kubernetes",
			Color:       "#DBEAFE",
		},
		{
			Name:        "pod",
			Description: "A kubernetes pod",
			Color:       "#FFEDD5",
		},
		{
			Name:        "job",
			Description: "A kubernetes job",
			Color:       "#DBEAFE",
		},
		{
			Name:        "cron job",
			Description: "A kubernetes cron job",
			Color:       "#DBEAFE",
		},
		{
			Name:        "deployment",
			Description: "A kubernetes deployment",
			Color:       "#DBEAFE",
		},
		{
			Name:        "replica set",
			Description: "A kubernetes replica set",
			Color:       "#DBEAFE",
		},
		{
			Name:        "daemon set",
			Description: "A kubernetes daemon set",
			Color:       "#DBEAFE",
		},
		{
			Name:        "stateful set",
			Description: "A kubernetes stateful set",
			Color:       "#DBEAFE",
		},
		{
			Name:        "docker",
			Description: "A docker image",
			Color:       "#FEE2E2",
		},
		{
			Name:        "ghcr",
			Description: "A ghcr image",
			Color:       "#FEE2E2",
		},
		{
			Name:        "github",
			Description: "A github project",
			Color:       "#FEE2E2",
		},
		{
			Name:        "up-to-date",
			Description: "Up-to-date images",
			Color:       "#DCFCE7",
		},
		{
			Name:        "outdated",
			Description: "Outdated images",
			Color:       "#FEE2E2",
		},
	}
	for _, tag := range defaultTags {
		if err := writeStore.InsertTag(context.Background(), &tag); err != nil {
			slog.Error("Failed to insert default tags", slog.Any("error", err))
			os.Exit(1)
		}
	}

	apiServer := api.NewServer(readStore)

	mux := http.NewServeMux()

	mux.Handle("/api/v1/", apiServer)

	ctx, cancel := context.WithCancel(context.Background())
	var wg errgroup.Group

	wg.Go(func() error {
		httpClient := httputil.NewClient(cache, 24*time.Hour)
		worker := worker.New(httpClient, writeStore)
		// TODO: Temp interval for now
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				// TODO: Temp time for now
				err := worker.ProcessOldReferences(ctx, 1, time.Now().Add(-2*time.Minute))
				if err != nil {
					slog.Error("Failed to process old references", slog.Any("error", err))
				}
			}
		}
	})

	wg.Go(func() error {
		// TODO: Listen on events and react on them once running, rather than just
		// check once or poll

		slog.Info("Fetching initial state")
		graph, err := kubernetesPlatform.Graph(ctx)
		if err != nil {
			return err
		}

		roots := graph.Roots()

		for _, root := range roots {
			imageNode := root.(platform.ImageNode)

			subgraph := graph.Subgraph(root.ID())

			edges := subgraph.Edges()
			nodes := subgraph.Nodes()

			mappedNodes := make(map[string]models.GraphNode)
			for _, node := range nodes {
				switch n := node.(type) {
				case kubernetes.Resource:
					mappedNodes[node.ID()] = models.GraphNode{
						Domain: "kubernetes",
						Type:   string(n.Kind()),
						Name:   n.Name(),
					}
				case platform.ImageNode:
					mappedNodes[node.ID()] = models.GraphNode{
						Domain: "oci",
						Type:   "image",
						Name:   imageNode.Reference.String(),
					}
				default:
					panic(fmt.Sprintf("unimplemented node type: %s", node.Type()))
				}
			}

			// TODO: platform-discovered tags
			tags := []string{}

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
			slog.Debug("Inserting raw image", slog.String("reference", rawImage.Reference))
			if err := writeStore.InsertRawImage(context.TODO(), rawImage); err != nil {
				slog.Error("Failed to insert raw image", slog.Any("error", err))
				return err
			}
		}

		allReferences := make([]string, 0)
		for _, root := range roots {
			imageNode := root.(platform.ImageNode)
			allReferences = append(allReferences, imageNode.Reference.String())
		}

		slog.Debug("Cleaning up removed images")
		removed, err := writeStore.DeleteNonPresent(context.TODO(), allReferences)
		if err != nil {
			slog.Error("Failed to clean up removed images", slog.Any("error", err))
			return err
		}
		slog.Debug("Cleaned up removed images successfully", slog.Int64("removed", removed))

		return nil
	})

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: apiServer,
	}

	wg.Go(func() error {
		slog.Info("Starting HTTP server")
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to serve", slog.Any("error", err))
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
				slog.Info("Caught signal, exiting gracefully")
				cancel()
				httpServer.Close()
				if err := readStore.Close(); err != nil {
					slog.Error("Failed to close read store", slog.Any("error", err))
					// Fallthrough
				}
				if err := writeStore.Close(); err != nil {
					slog.Error("Failed to close write store", slog.Any("error", err))
					// Fallthrough
				}
			} else {
				slog.Info("Caught signal, exiting now")
				os.Exit(1)
			}
		}
	}()

	if err := wg.Wait(); err != nil && err != ctx.Err() {
		slog.Error("Failed to run", slog.Any("error", err))
		os.Exit(1)
	}
}
