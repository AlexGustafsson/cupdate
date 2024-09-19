package jobs

import (
	"log/slog"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/pipeline"
	"github.com/AlexGustafsson/cupdate/internal/registry/docker"
)

type GetDockerHubRepositoryJob struct {
	Output GetDockerHubRepositoryJobOutput
}

type GetDockerHubRepositoryJobOutput struct {
	Repository string
}

func GetDockerHubRepository() *GetDockerHubRepositoryJob {
	return &GetDockerHubRepositoryJob{
		Output: GetDockerHubRepositoryJobOutput{
			Repository: "get-docker-hub-repository/repository",
		},
	}
}

func (j GetDockerHubRepositoryJob) Execute(ctx pipeline.Context[ImageData]) error {
	log := slog.With(slog.String("imageReference", ctx.Data.ImageReference.String()))

	cacheKey := "pipeline/get-docker-hub-repository-v1/" + ctx.Data.ImageReference.Name()
	var repository *docker.Repository
	if err := ctx.Cache().GetJSON(ctx, cacheKey, &repository, 24*time.Hour); err != nil {
		log.Error("Failed to get cache", slog.Any("error", err))
		// Fallthrough
	}

	if repository == nil {
		log.Debug("Fetching repository")

		client := &docker.Client{}

		var err error
		repository, err = client.GetRepository(ctx, ctx.Data.ImageReference)
		if err != nil {
			log.Error("Failed to get repository", slog.Any("error", err))
			return err
		}

		if err := ctx.Cache().SetJSON(ctx, cacheKey, &repository); err != nil {
			log.Error("Failed to set cache", slog.Any("error", err))
			// Fallthrough
		}
	}

	if repository == nil {
		log.Info("Image repository not found")
		return nil
	}

	ctx.SetOutput(j.Output.Repository, repository)

	return nil
}
