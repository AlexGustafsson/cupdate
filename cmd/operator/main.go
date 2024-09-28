package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/AlexGustafsson/cupdate/internal/api"
	"github.com/AlexGustafsson/cupdate/internal/models"
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

	// cache, err := cache.NewDiskCache("./cache")
	// if err != nil {
	// 	slog.Error("Failed to serve", slog.Any("error", err))
	// 	os.Exit(1)
	// }

	data := &api.InMemoryAPI{
		Store: &models.Store{
			Tags:         []*models.Tag{},
			Images:       []*models.Image{},
			Descriptions: map[string]*models.ImageDescription{},
			ReleaseNotes: map[string]*models.ImageReleaseNotes{},
			Graphs:       map[string]*models.Graph{},
		},
	}

	store := &models.UnprocessedStore{
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
		Graphs:       map[string][]*models.Graph{},
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

			imageName := ref.Name()

			imageTag := "latest"
			if ref.HasTag {
				imageTag = ref.Tag
			}

			image := &models.Image{
				Name:           imageName,
				CurrentVersion: imageTag,
				LatestVersion:  imageTag,
				// TODO: Tags should include pod, job, cron job, deployment set etc.
				// Everything's a pod, so try to use the topmost descriptor
				Tags:  []string{"k8s", "pod"},
				Links: []*models.ImageLink{},
				Image: "",
			}

			children := graph.Children(root.ID())
			f := func(rootID string)

			origins := graph.Origins(ref)
			for _, origin := range origins {
				// TODO: Build actual graph. We don't handle duplicates right now...
				root := &models.GraphNode{
					Domain: "oci",
					Type:   "image",
					Name:   imageName,
				}

				origin := origin.(*kubernetes.Origin)

				container := &models.GraphNode{
					Domain: "kubernetes",
					Type:   "core/v1/container",
					Name:   origin.Container.Name,
				}
				root.Parents = []*models.GraphNode{container}

				pod := &models.GraphNode{
					Domain: "kubernetes",
					Type:   "core/v1/pod",
					Name:   origin.Container.Pod.Name,
				}
				if pod.Name == "" && origin.Container.Pod.IsTemplate {
					pod.Name = "(template)"
				}
				container.Parents = []*models.GraphNode{pod}

				tag := "pod"
				currentNode := pod
				currentParent := origin.Container.Pod.Parent
				for currentParent != nil {
					node := &models.GraphNode{
						Domain:  "kubernetes",
						Type:    string(currentParent.ResourceKind),
						Name:    currentParent.Name,
						Parents: make([]*models.GraphNode, 0),
					}
					currentNode.Parents = []*models.GraphNode{node}

					switch currentParent.ResourceKind {
					case kubernetes.ResourceKindAppsV1Deployment:
						tag = "deployment"
					case kubernetes.ResourceKindAppsV1DaemonSet:
						tag = "daemon set"
					case kubernetes.ResourceKindAppsV1ReplicaSet:
						tag = "replica set"
					case kubernetes.ResourceKindBatchV1CronJob:
						tag = "cron job"
					case kubernetes.ResourceKindBatchV1Job:
						tag = "job"
					case kubernetes.ResourceKindAppsV1StatefulSet:
						tag = "stateful set"
					}

					currentNode = node
					currentParent = currentParent.Parent
				}

				image.Tags = append(image.Tags, tag)

				// Namespace is implicit
				currentNode.Parents = []*models.GraphNode{{
					Domain:  "kubernetes",
					Type:    "core/v1/namespace",
					Name:    origin.Container.Namespace,
					Parents: make([]*models.GraphNode, 0),
				}}

				graph := &models.Graph{
					Root: root,
				}

				// TODO: Can overwrite. The graph should be shared among all ways the
				// image is used
				_, ok := graphs[ref.String()]
				if !ok {
					g := make([]*models.Graph, 0)
					graphs[ref.String()] = g
				}

				graphs[ref.String()] = append(graphs[ref.String()], graph)
			}

			images = append(images, image)
		}

		store.Images = images
		store.Graphs = graphs

		// pipeline := pipeline.New(cache, jobs.DefaultJobs())
		// TODO: How will deduplication work with this when we invoke just one image
		// at a time?
		// for _, image := range images {
		// 	processedStore, err := pipeline.Run(ctx, jobs.Image)
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		// data.Store = processedStore
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
