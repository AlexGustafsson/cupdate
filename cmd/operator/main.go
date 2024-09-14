package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/AlexGustafsson/cupdate/internal/api"
	"github.com/AlexGustafsson/cupdate/internal/source/k8s"
	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/rest"
)

func main() {
	config := &rest.Config{
		Host: "http://localhost:8001",
	}

	source, err := k8s.New(config)
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
		Tags: []*api.Tag{
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
				Name:        "chron job",
				Description: "A kubernetes chron job",
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
		Images:       []*api.Image{},
		Descriptions: map[string]*api.ImageDescription{},
		ReleaseNotes: map[string]*api.ImageReleaseNotes{},
		Graphs:       map[string]*api.Graph{},
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
		entries, err := source.Entries(ctx)
		if err != nil {
			return err
		}

		images := make([]*api.Image, 0)
		graphs := make(map[string]*api.Graph)

		for _, entry := range entries {
			origin := entry.Origin.(*k8s.Origin)

			// For now, ignore definitions referencing images, just use running pods
			if origin.Container.Pod.IsTemplate {
				continue
			}

			image := &api.Image{
				Name:           entry.Image,
				CurrentVersion: entry.Version,
				LatestVersion:  entry.Version,
				// TODO: Tags should include pod, job, chron job, deployment set etc.
				// Everything's a pod, so try to use the topmost descriptor
				Tags:  []string{"k8s"},
				Links: []*api.ImageLink{},
				Image: "",
			}

			// TODO: Build actual graph. We don't handle duplicates right now...
			root := &api.GraphNode{
				Domain: "oci",
				Type:   "image",
				Name:   entry.Image,
			}

			container := &api.GraphNode{
				Domain: "kubernetes",
				Type:   "core/v1/container",
				Name:   origin.Container.Name,
			}
			root.Parents = []*api.GraphNode{container}

			pod := &api.GraphNode{
				Domain: "kubernetes",
				Type:   "core/v1/pod",
				Name:   origin.Container.Pod.Name,
			}
			container.Parents = []*api.GraphNode{pod}

			tag := "pod"
			currentNode := pod
			currentParent := origin.Container.Pod.Parent
			for currentParent != nil {
				node := &api.GraphNode{
					Domain:  "kubernetes",
					Type:    string(currentParent.ResourceKind),
					Name:    currentParent.Name,
					Parents: make([]*api.GraphNode, 0),
				}
				currentNode.Parents = []*api.GraphNode{node}

				switch currentParent.ResourceKind {
				case k8s.ResourceKindAppsV1Deployment:
					tag = "deployment"
				case k8s.ResourceKindAppsV1DaemonSet:
					tag = "daemon set"
				case k8s.ResourceKindAppsV1ReplicaSet:
					tag = "replica set"
				case k8s.ResourceKindBatchV1CronJob:
					tag = "chron job"
				case k8s.ResourceKindBatchV1Job:
					tag = "job"
				case k8s.ResourceKindAppsV1StatefulSet:
					tag = "stateful set"
				}

				currentNode = node
				currentParent = currentParent.Parent
			}

			image.Tags = append(image.Tags, tag)
			images = append(images, image)

			// Namespace is implicit
			currentNode.Parents = []*api.GraphNode{{
				Domain:  "kubernetes",
				Type:    "core/v1/namespace",
				Name:    origin.Container.Namespace,
				Parents: make([]*api.GraphNode, 0),
			}}

			graph := &api.Graph{
				Root: root,
			}

			// TODO: Can overwrite. The graph should be shared among all ways the
			// image is used
			graphs[entry.Image+":"+entry.Version] = graph
		}

		data.Images = images
		data.Graphs = graphs

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
