package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/store"
	"github.com/AlexGustafsson/cupdate/internal/workflow/imageworkflow"
)

type Worker struct {
	httpClient *httputil.Client
	store      *store.Store
}

func New(httpClient *httputil.Client, store *store.Store) *Worker {
	return &Worker{
		httpClient: httpClient,
		store:      store,
	}
}

func (w *Worker) ProcessOldReferences(ctx context.Context, n int, notUpdatedSince time.Time) error {
	slog.Debug("Identifying old references to process")
	images, err := w.store.ListRawImages(ctx, &store.ListRawImagesOptions{
		NotUpdatedSince: notUpdatedSince,
		Limit:           n,
	})
	if err != nil {
		return err
	}

	if len(images) == 0 {
		slog.Debug("Found no old references, skipping run")
		return nil
	}

	slog.Debug("Processing old references", slog.Int("n", len(images)))
	for _, image := range images {
		if err := w.ProcessRawImage(ctx, image); err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) ProcessRawImage(ctx context.Context, image models.RawImage) error {
	reference, err := oci.ParseReference(image.Reference)
	if err != nil {
		return err
	}

	log := slog.With(slog.String("reference", reference.String()))
	log.Debug("Processing reference")

	// Try to update the image's process time
	// NOTE: There's a race here if the entry has been modified or removed since
	// it was loaded from the store. It will eventually be corrent and consistent,
	// though. And it's unlikely to happen. So let's not keep a transaction during
	// processing for now. If it becomes important, we could keep an "etag" /
	// generation id in the document and throw an error if the expectation fails.
	// NOTE: Always update immediately as a failure to process or update the image
	// could be a reoccuring issue, so try to process other images before retrying
	// the failing image.
	image.LastProcessed = time.Now()
	if err := w.store.InsertRawImage(ctx, &image); err != nil {
		return err
	}

	log.Debug("Running workflow")
	data := &imageworkflow.Data{
		ImageReference:  reference,
		Image:           "",
		LatestReference: reference,
		Tags:            make([]string, 0),
		Description:     "",
		FullDescription: nil,
		ReleaseNotes:    nil,
		Links:           make([]models.ImageLink, 0),
		Graph:           image.Graph,
	}

	for _, tag := range image.Tags {
		data.InsertTag(tag)
	}

	workflow := imageworkflow.New(w.httpClient, data)
	if err := workflow.Run(ctx); err != nil {
		log.Error("Failed to run pipeline for image", slog.Any("error", err))
		// Fallthrough - insert what we have
	}

	// Add some basic tags
	if data.ImageReference == data.LatestReference {
		data.Tags = append(data.Tags, "up-to-date")
	} else {
		data.Tags = append(data.Tags, "outdated")
	}

	if err := w.store.InsertImage(context.TODO(), &models.Image{
		Reference:       data.ImageReference.String(),
		LatestReference: data.LatestReference.String(),
		Description:     data.Description,
		// TODO: Tags should include pod, job, cron job, deployment set etc.
		// Everything's a pod, so try to use the topmost descriptor
		Tags:         data.Tags,
		Image:        data.Image,
		Links:        data.Links,
		LastModified: time.Now(),
	}); err != nil {
		log.Error("Failed to insert image", slog.Any("error", err))
		// Fallthrough - try to insert what we have
	}

	if data.FullDescription != nil {
		if err := w.store.InsertImageDescription(ctx, reference.String(), data.FullDescription); err != nil {
			log.Error("Failed to insert image description", slog.Any("error", err))
			// Fallthrough - try to insert what we have
		}
	}

	if data.ReleaseNotes != nil {
		if err := w.store.InsertImageReleaseNotes(ctx, reference.String(), data.ReleaseNotes); err != nil {
			log.Error("Failed to insert image description", slog.Any("error", err))
			// Fallthrough - try to insert what we have
		}
	}

	if err := w.store.InsertImageGraph(ctx, reference.String(), &data.Graph); err != nil {
		log.Error("Failed to insert image description", slog.Any("error", err))
		// Fallthrough - try to insert what we have
	}

	log.Debug("Updated data")
	return nil
}
