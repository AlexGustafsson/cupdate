package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/store"
	"github.com/AlexGustafsson/cupdate/internal/worker"
)

func HandleScheduling(ctx context.Context, config *Config, processQueue *worker.Queue[oci.Reference], readStore *store.Store) {
	ticker := time.NewTicker(config.Processing.Interval)
	defer ticker.Stop()
	defer slog.Debug("Closed scheduling")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			slog.Debug("Identifying old references to process")
			ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
			images, err := readStore.ListRawImages(ctx, &store.ListRawImagesOptions{
				NotUpdatedSince: time.Now().Add(-config.Processing.MinAge),
				Limit:           config.Processing.Items,
			})
			cancel()
			if err != nil {
				slog.Error("Failed to process old references", slog.Any("error", err))
				continue
			}

			for _, image := range images {
				reference, err := oci.ParseReference(image.Reference)
				if err != nil {
					slog.Error("Unexpectedly failed to parse reference from store", slog.Any("error", err), slog.String("reference", image.Reference))
					return
				}

				processQueue.PushBack(reference)
			}
		}
	}
}

func HandleProcessing(ctx context.Context, config *Config, worker *worker.Worker, processQueue *worker.Queue[oci.Reference]) {
	defer slog.Debug("Closed processing")
	for reference := range processQueue.Pull() {
		ctx, cancel := context.WithTimeout(ctx, config.Processing.Timeout)
		err := worker.ProcessRawImage(ctx, reference)
		cancel()
		if err != nil {
			slog.Error("Failed to process queued raw image", slog.Any("error", err), slog.String("reference", reference.String()))
			continue
		}
	}
}
