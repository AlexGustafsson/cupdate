package jobs

import (
	"log/slog"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/pipeline"
	"github.com/AlexGustafsson/cupdate/internal/registry/docker"
)

type GetDockerHubRepositoryOwnerJob struct {
	Output GetDockerHubRepositoryOwnerJobOutput

	repository string
}

type GetDockerHubRepositoryOwnerJobOutput struct {
	Owner string
}

func GetDockerHubRepositoryOwner(repository string) *GetDockerHubRepositoryOwnerJob {
	return &GetDockerHubRepositoryOwnerJob{
		Output: GetDockerHubRepositoryOwnerJobOutput{
			Owner: "get-docker-hub-repository-owner/owner",
		},
		repository: repository,
	}
}

func (j GetDockerHubRepositoryOwnerJob) Execute(ctx pipeline.Context[ImageData]) error {
	ctx.Lock()
	defer ctx.Unlock()

	log := slog.With(slog.String("imageReference", ctx.Data.ImageReference.String()))
	var repository *docker.Repository
	if !ctx.GetOutput(j.repository, &repository) {
		log.Debug("Skipping job - no repository found")
		return nil
	}

	cacheKey := "pipeline/get-docker-hub-repository-owner-v1/" + repository.Namespace
	var owner *docker.Entity
	if err := ctx.Cache().GetJSON(ctx, cacheKey, &owner, 24*time.Hour); err != nil {
		log.Error("Failed to get cache", slog.Any("error", err))
		// Fallthrough
	}

	if owner == nil {
		log.Debug("Fetching owner")

		client := &docker.Client{}

		var err error
		owner, err = client.GetOrganizationOrUser(ctx, repository.Namespace)
		if err != nil {
			log.Error("Failed to get owner", slog.Any("error", err))
			return err
		}

		if err := ctx.Cache().SetJSON(ctx, cacheKey, &owner); err != nil {
			log.Error("Failed to set cache", slog.Any("error", err))
			// Fallthrough
		}
	}

	if owner == nil {
		log.Info("Image owner not found")
		return nil
	}

	ctx.SetOutput(j.Output.Owner, owner)

	return nil
}
