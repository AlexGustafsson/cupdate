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
		Tags: []api.Tag{
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
		Images:       []api.Image{},
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

		images := make([]api.Image, 0)
		graphs := make(map[string]*api.Graph)

		for _, entry := range entries {
			origin := entry.Origin.(*k8s.Origin)

			// For now, ignore definitions referencing images, just use running pods
			if origin.Container.Pod.IsTemplate {
				continue
			}

			images = append(images, api.Image{
				Name:           entry.Image,
				CurrentVersion: entry.Version,
				LatestVersion:  entry.Version,
				Tags:           []string{"k8s", "pod"},
				Links:          []api.ImageLink{},
				Image:          "",
			})

			// TODO: Build actual graph. Can occur more than once
			graphs[entry.Image+":"+entry.Version] = &api.Graph{
				Root: api.GraphNode{
					Domain:  "oic",
					Type:    "image",
					Name:    entry.Image,
					Parents: make([]api.GraphNode, 0),
				},
			}
		}

		data.Images = images

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
