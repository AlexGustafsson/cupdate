package jobs

import (
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/github"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/pipeline"
	"github.com/AlexGustafsson/cupdate/internal/registry"
	"github.com/AlexGustafsson/cupdate/internal/registry/docker"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

func DefaultJobs() pipeline.Job[ImageData] {
	setupRegistryClient := SetupRegistryClient()
	getManifests := GetImageManifests(setupRegistryClient.Output.Client)

	getDockerHubRepository := GetDockerHubRepository()
	getDockerHubRepositoryOwner := GetDockerHubRepositoryOwner(getDockerHubRepository.Output.Repository)

	getDockerHubLatestVersion := GetDockerHubLatestVersion()

	getGitHubRelease := GetGitHubRelease(getManifests.Output.Manifests)

	return pipeline.Series[ImageData]{
		setupRegistryClient,
		pipeline.Parallel[ImageData]{
			// SVC
			pipeline.Series[ImageData]{
				getManifests,
				// Add link to SVC
				pipeline.JobFunc[ImageData](func(ctx pipeline.Context[ImageData]) error {
					var manifests []oci.Manifest
					if ok := ctx.GetOutput(getManifests.Output.Manifests, &manifests); !ok {
						return nil
					}

					ctx.Lock()
					defer ctx.Unlock()
					for _, manifest := range manifests {
						source := manifest.SourceAnnotation()
						if source != "" {
							// TODO: Identify source more granularly (GitHub, GitLab, source
							// hut etc.)
							*ctx.Data.Links = append(*ctx.Data.Links, models.ImageLink{
								Type: "svc",
								URL:  source,
							})
						}
					}

					return nil
				}),
			},

			// Docker Hub
			pipeline.Series[ImageData]{
				getDockerHubRepository,
				getDockerHubRepositoryOwner,
				pipeline.Parallel[ImageData]{
					getDockerHubLatestVersion,
					// Set description
					pipeline.JobFunc[ImageData](func(ctx pipeline.Context[ImageData]) error {
						var repository *docker.Repository
						if ok := ctx.GetOutput(getDockerHubRepository.Output.Repository, &repository); !ok {
							return nil
						}

						*ctx.Data.Description = repository.FullDescription

						return nil
					}),
					// Set image
					pipeline.JobFunc[ImageData](func(ctx pipeline.Context[ImageData]) error {
						var owner *docker.Entity
						if ok := ctx.GetOutput(getDockerHubRepositoryOwner.Output.Owner, &owner); !ok {
							return nil
						}

						*ctx.Data.Image = owner.GravatarURL

						return nil
					}),
					// Add link to Docker Hub
					pipeline.JobFunc[ImageData](func(ctx pipeline.Context[ImageData]) error {
						var repository *docker.Repository
						if ok := ctx.GetOutput(getDockerHubRepository.Output.Repository, &repository); !ok {
							return nil
						}

						owner := repository.Namespace
						if owner == "library" {
							owner = "_"
						}

						ctx.Lock()
						defer ctx.Unlock()
						*ctx.Data.Links = append(*ctx.Data.Links, models.ImageLink{
							Type: "docker",
							URL:  fmt.Sprintf("https://hub.docker.com/%s/%s", owner, repository.Name),
						})

						return nil
					}),
					// Set latest version
					pipeline.JobFunc[ImageData](func(ctx pipeline.Context[ImageData]) error {
						var image *registry.Image
						if ok := ctx.GetOutput(getDockerHubLatestVersion.Output.Image, &image); !ok {
							return nil
						}

						*ctx.Data.LatestVersion = image.Name

						return nil
					}),
				},
			},
		},
		// Once we have identified the latest image by any means, as well as the SVC
		// repository of the image, let's try to find some release notes
		getGitHubRelease,
		// Set release notes
		pipeline.JobFunc[ImageData](func(ctx pipeline.Context[ImageData]) error {
			var release *github.Release
			if ok := ctx.GetOutput(getGitHubRelease.Output.Release, &release); !ok {
				return nil
			}

			*ctx.Data.ReleaseNotes = release.Description

			return nil
		}),
		// TODO: Deduplicate tags
		// TODO: Deduplicate links
	}
}
