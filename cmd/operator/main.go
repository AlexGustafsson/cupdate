package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/AlexGustafsson/cupdate/internal/api"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/platform/kubernetes"
	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/rest"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

	config := &rest.Config{
		Host: "http://localhost:8001",
	}

	platform, err := kubernetes.NewPlatform(config)
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

	apiServer := api.NewServer(data)

	mux := http.NewServeMux()

	mux.Handle("/api/v1/", apiServer)

	ctx, cancel := context.WithCancel(context.Background())
	var wg errgroup.Group

	wg.Go(func() error {
		// TODO: Listen on events and react on them once running, rather than just
		// check once or poll

		slog.Info("Fetching initial state")
		graph, err := platform.Graph(ctx)
		if err != nil {
			return err
		}

		fmt.Println(graph.String())

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
