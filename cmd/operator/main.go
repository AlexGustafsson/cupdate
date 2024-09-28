package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/AlexGustafsson/cupdate/internal/api"
	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/pipeline"
	"github.com/AlexGustafsson/cupdate/internal/pipeline/jobs"
	"github.com/AlexGustafsson/cupdate/internal/platform"
	"github.com/AlexGustafsson/cupdate/internal/platform/kubernetes"
	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/rest"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

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

	data := &api.InMemoryAPI{
		Store: &models.Store{
			Tags:         []*models.Tag{},
			Images:       []*models.Image{},
			Descriptions: map[string]*models.ImageDescription{},
			ReleaseNotes: map[string]*models.ImageReleaseNotes{},
			Graphs:       map[string]models.Graph{},
		},
	}

	processedStore := &models.Store{
		Tags: []*models.Tag{
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
		},
		Images:       []*models.Image{},
		Descriptions: map[string]*models.ImageDescription{},
		ReleaseNotes: map[string]*models.ImageReleaseNotes{},
		Graphs:       map[string]models.Graph{},
	}

	apiServer := api.NewServer(data)

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
			ref := imageNode.Reference

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
						Name:   ref.String(),
					}
				default:
					panic(fmt.Sprintf("unimplemented node type: %s", node.Type()))
				}
			}

			mappedGraph := models.Graph{
				Edges: edges,
				Nodes: mappedNodes,
			}

			processedStore.Graphs[ref.String()] = mappedGraph
		}

		pipeline := pipeline.New(cache, jobs.DefaultJobs())
		for _, root := range roots {
			imageNode := root.(platform.ImageNode)

			var image string
			latestVersion := imageNode.Reference
			tags := make([]string, 0)
			var description *models.ImageDescription
			var releaseNotes *models.ImageReleaseNotes
			links := make([]models.ImageLink, 0)

			_, err := pipeline.Run(ctx, jobs.ImageData{
				ImageReference: imageNode.Reference,
				Image:          &image,
				LatestVersion:  &latestVersion,
				Tags:           &tags,
				Description:    &description,
				ReleaseNotes:   &releaseNotes,
				Links:          &links,
			})
			if err != nil {
				slog.Error("Failed to run pipeline for image", slog.Any("error", err))
				continue
			}

			processedStore.Images = append(processedStore.Images, &models.Image{
				Name: imageNode.Reference.Name(),
				// TODO: Handle digests, not just tags
				CurrentVersion: imageNode.Reference.Tag,
				LatestVersion:  latestVersion.Tag,
				// TODO: Tags should include pod, job, cron job, deployment set etc.
				// Everything's a pod, so try to use the topmost descriptor
				Tags:  tags,
				Image: image,
				Links: links,
			})

			if description != nil {
				processedStore.Descriptions[imageNode.Reference.String()] = description
			}
			if releaseNotes != nil {
				processedStore.ReleaseNotes[imageNode.Reference.String()] = releaseNotes
			}
		}

		data.Store = processedStore
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
