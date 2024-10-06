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
	"github.com/AlexGustafsson/cupdate/internal/workflow/imageworkflow"
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

	client := httputil.NewClient(cache, 24*time.Hour)

	cwd, err := os.Getwd()
	if err != nil {
		slog.Error("Failed to identify working directory", slog.Any("error", err))
		os.Exit(1)
	}
	store, err := store.New("file://" + path.Join(cwd, "dbv1.sqlite"))
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
			Description: "A docker container",
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
		if err := store.InsertTag(context.Background(), &tag); err != nil {
			slog.Error("Failed to insert default tags", slog.Any("error", err))
			os.Exit(1)
		}
	}

	apiServer := api.NewServer(store)

	mux := http.NewServeMux()

	mux.Handle("/api/v1/", apiServer)

	ctx, cancel := context.WithCancel(context.Background())
	var wg errgroup.Group

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

			slog.Debug("Running workflow", slog.String("image", imageNode.Reference.String()))
			data := &imageworkflow.Data{
				ImageReference:  imageNode.Reference,
				Image:           "",
				LatestReference: imageNode.Reference,
				Tags:            make([]string, 0),
				Description:     nil,
				ReleaseNotes:    nil,
				Links:           make([]models.ImageLink, 0),
			}
			workflow := imageworkflow.New(client, data)

			if err := workflow.Run(ctx); err != nil {
				slog.Error("Failed to run pipeline for image", slog.Any("error", err))
				// Fallthrough - insert what we have
			}

			if err := store.InsertImage(context.TODO(), &models.Image{
				Reference:       data.ImageReference.String(),
				LatestReference: data.LatestReference.String(),
				// TODO:
				Description: "",
				// TODO: Tags should include pod, job, cron job, deployment set etc.
				// Everything's a pod, so try to use the topmost descriptor
				Tags:         data.Tags,
				Image:        data.Image,
				Links:        data.Links,
				LastModified: time.Now(),
			}); err != nil {
				slog.Error("Failed to insert image graph", slog.Any("error", err))
				// Fallthrough - insert what we have
			}

			if data.Description != nil {
				if err := store.InsertImageDescription(context.TODO(), imageNode.Reference.String(), data.Description); err != nil {
					slog.Error("Failed to insert image description", slog.Any("error", err))
					// Fallthrough - insert what we have
				}
			}

			if data.ReleaseNotes != nil {
				if err := store.InsertImageReleaseNotes(context.TODO(), imageNode.Reference.String(), data.ReleaseNotes); err != nil {
					slog.Error("Failed to insert image description", slog.Any("error", err))
					// Fallthrough - insert what we have
				}
			}

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

			mappedGraph := models.Graph{
				Edges: edges,
				Nodes: mappedNodes,
			}

			if err := store.InsertImageGraph(context.TODO(), imageNode.Reference.String(), &mappedGraph); err != nil {
				slog.Error("Failed to insert image graph", slog.Any("error", err))
				return err
			}
		}

		return nil
	})

	wg.Go(func() error {
		slog.Info("Starting HTTP server")
		err := http.ListenAndServe(":8080", apiServer)
		if err != nil && err != ctx.Err() {
			slog.Error("Failed to serve", slog.Any("error", err))
		}
		return err
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
			} else {
				slog.Info("Caught signal, exiting now")
				os.Exit(1)
			}
		}
	}()

	if err := wg.Wait(); err != nil {
		slog.Error("Failed to run", slog.Any("error", err))
		os.Exit(1)
	}
}
