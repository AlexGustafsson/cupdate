package jobs

import (
	"log/slog"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/pipeline"
	"github.com/AlexGustafsson/cupdate/internal/registry"
	"github.com/AlexGustafsson/cupdate/internal/registry/docker"
)

type GetDockerHubLatestVersionJob struct {
	Output GetDockerHubLatestVersionJobOutput
}

type GetDockerHubLatestVersionJobOutput struct {
	Image string
}

func GetDockerHubLatestVersion() *GetDockerHubLatestVersionJob {
	return &GetDockerHubLatestVersionJob{
		Output: GetDockerHubLatestVersionJobOutput{
			Image: "get-docker-hub-latest-version/latest-version",
		},
	}
}

func (j GetDockerHubLatestVersionJob) Execute(ctx pipeline.Context[ImageData]) error {
	log := slog.With(slog.String("imageReference", ctx.Data.ImageReference.String()))

	if !ctx.Data.ImageReference.HasTag {
		log.Info("Skipping non-tagged image")
		return nil
	}

	cacheKey := "pipeline/get-docker-hub-latest-version-v1/" + ctx.Data.ImageReference.String()
	var image *registry.Image
	if err := ctx.Cache().GetJSON(ctx, cacheKey, &image, 24*time.Hour); err != nil {
		log.Error("Failed to get cache", slog.Any("error", err))
		// Fallthrough
	}

	if image == nil {
		log.Debug("Fetching latest image")

		client := &docker.Client{}

		var err error
		image, err = client.GetLatestVersion(ctx, ctx.Data.ImageReference)
		if err != nil {
			log.Error("Failed to get latest image", slog.Any("error", err))
			return err
		}

		if err := ctx.Cache().SetJSON(ctx, cacheKey, &image); err != nil {
			log.Error("Failed to set cache", slog.Any("error", err))
			// Fallthrough
		}
	}

	if image == nil {
		log.Info("No new version found")
		return nil
	}

	ctx.SetOutput(j.Output.Image, image)

	return nil
}
