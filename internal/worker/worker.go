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

func (w *Worker) ProcessOldReferences(ctx context.Context, minAge time.Duration) error {
	page, err := w.store.ListImages(ctx, &store.ListImageOptions{
		SortProperty: store.SortPropertyLastModified,
		Limit:        1,
		Page:         0,
	})
	if err != nil {
		return err
	}

	// Should we just do one image table with just references and then keep
	// details elsewhere?
	// The issue is that we'll have event-driven pushes (inserts / upserts) to
	// something with incomplete data from k8s. Fields like description will be
	// empty. The only fields that are well defined are tags and graphs?
	// Instead, perhaps maintain a table that the platform populates with
	// whatever data it has (such as the graph, reference and tags). Then, every
	// now and then, run this worker which takes data from the unprocessed "source
	// of truth" tables and pushes it to the enriched tables that currently exist.

	// This works for now, because we only read data once and this likely runs
	// after that.

	// This is also a sign of that we need to change how the data is stored;
	// we get all information about the images here, but we just wanted to loop
	// over outdated references...

	// Perhaps let each platform own the "platform" table (can be shared, it's not
	// likely that we at ever point will support k8s+docker simultaneously).
	for _, image := range page.Images {
		reference, err := oci.ParseReference(image.Reference)
		if err != nil {
			return err
		}

		if err := w.ProcessReference(ctx, reference); err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) ProcessReference(ctx context.Context, reference oci.Reference) error {
	slog.Debug("Running workflow", slog.String("image", reference.String()))
	data := &imageworkflow.Data{
		ImageReference:  reference,
		Image:           "",
		LatestReference: reference,
		Tags:            make([]string, 0),
		Description:     nil,
		ReleaseNotes:    nil,
		Links:           make([]models.ImageLink, 0),
	}
	workflow := imageworkflow.New(w.httpClient, data)

	if err := workflow.Run(ctx); err != nil {
		slog.Error("Failed to run pipeline for image", slog.Any("error", err))
		// Fallthrough - insert what we have
	}

	if err := w.store.InsertImage(context.TODO(), &models.Image{
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
		if err := w.store.InsertImageDescription(context.TODO(), reference.String(), data.Description); err != nil {
			slog.Error("Failed to insert image description", slog.Any("error", err))
			// Fallthrough - insert what we have
		}
	}

	if data.ReleaseNotes != nil {
		if err := w.store.InsertImageReleaseNotes(context.TODO(), reference.String(), data.ReleaseNotes); err != nil {
			slog.Error("Failed to insert image description", slog.Any("error", err))
			// Fallthrough - insert what we have
		}
	}

	return nil
}
