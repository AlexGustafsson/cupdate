package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/store"
)

type ImageWorkerScheduler struct {
	store *store.Store
}

func NewImageWorkerScheduler(store *store.Store) *ImageWorkerScheduler {
	return &ImageWorkerScheduler{
		store: store,
	}
}

// PushTo continuously pushes old references to the queue at a set interval.
func (s *ImageWorkerScheduler) PushTo(ctx context.Context, queue *Queue[oci.Reference], interval time.Duration, minAge time.Duration, batchSize int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			s.queueBatch(ctx, queue, time.Now().Add(-minAge), batchSize)
			cancel()
		}
	}
}

func (s *ImageWorkerScheduler) queueBatch(ctx context.Context, queue *Queue[oci.Reference], notUpdatedSince time.Time, batchSize int) {
	slog.DebugContext(ctx, "Identifying old references to process")
	images, err := s.store.ListRawImages(ctx, &store.ListRawImagesOptions{
		NotUpdatedSince: notUpdatedSince,
		Limit:           batchSize,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to process old references", slog.Any("error", err))
		return
	}

	for _, image := range images {
		reference, err := oci.ParseReference(image.Reference)
		if err != nil {
			slog.ErrorContext(ctx, "Unexpectedly failed to parse reference from store", slog.Any("error", err), slog.String("reference", image.Reference))
			return
		}

		queue.Push(reference)
	}
}
