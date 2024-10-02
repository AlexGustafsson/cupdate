package jobs

import (
	"log/slog"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/pipeline"
	"github.com/AlexGustafsson/cupdate/internal/registry"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

type GetImageManifestsJob struct {
	Output GetImageManifestsJobOutput

	client string
}

type GetImageManifestsJobOutput struct {
	Manifests string
}

func GetImageManifests(client string) *GetImageManifestsJob {
	return &GetImageManifestsJob{
		client: client,
		Output: GetImageManifestsJobOutput{
			Manifests: "get-image-manifests/manifests",
		},
	}
}

func (j GetImageManifestsJob) Execute(ctx pipeline.Context[ImageData]) error {
	ctx.Lock()
	defer ctx.Unlock()

	var client registry.Client
	ctx.MustGetOutput(j.client, &client)

	cacheKey := "pipeline/get-image-manifests-v1/" + ctx.Data.ImageReference.String()
	var manifests []oci.Manifest
	if err := ctx.Cache().GetJSON(ctx, cacheKey, &manifests, 24*time.Hour); err != nil {
		slog.Error("Failed to get cache", slog.Any("error", err))
		// Fallthrough
	}

	log := slog.With(slog.String("domain", ctx.Data.ImageReference.Domain), slog.String("imageReference", ctx.Data.ImageReference.String()))

	if manifests == nil {
		log.Debug("Fetching annotations")

		var err error
		manifests, err = client.GetManifests(ctx, ctx.Data.ImageReference)
		if err != nil {
			slog.Error("Failed to get manifests", slog.Any("error", err))
			return err
		}

		if err := ctx.Cache().SetJSON(ctx, cacheKey, &manifests); err != nil {
			slog.Error("Failed to set cache", slog.Any("error", err))
			// Fallthrough
		}
	}

	if manifests == nil {
		slog.Error("Image manifests not found")
		return nil
	}

	if len(manifests) == 0 {
		slog.Error("Got zero manifests")
		return nil
	}

	ctx.SetOutput(j.Output.Manifests, manifests)

	return nil
}
