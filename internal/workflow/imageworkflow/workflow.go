package imageworkflow

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/github"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/registry/docker"
	"github.com/AlexGustafsson/cupdate/internal/registry/ghcr"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func New(httpClient *httputil.Client, data *Data) workflow.Workflow {
	return workflow.Workflow{
		Name: "Process image",
		Jobs: []workflow.Job{
			{
				ID:   "oci",
				Name: "Get OCI information",
				Steps: []workflow.Step{
					SetupRegistryClient().
						WithID("registry").
						With("httpClient", httpClient).
						With("reference", data.ImageReference),
					GetManifests().
						WithID("manifests").
						With("registryClient", workflow.Ref{Key: "step.registry.client"}).
						With("reference", data.ImageReference),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						domain, err := workflow.GetValue[string](ctx, "step.registry.domain")
						if err != nil {
							return nil, err
						}

						data.InsertLink(models.ImageLink{
							Type: "oci-registry",
							URL:  "https://" + domain,
						})
						return nil, nil
					}),
				},
			},
			{
				ID:        "docker",
				Name:      "Get Docker Hub information",
				DependsOn: []string{"oci"},
				// Only run for Docker images
				If: func(ctx workflow.Context) (bool, error) {
					domain, err := workflow.GetValue[string](ctx, "job.oci.step.registry.domain")
					if err != nil {
						return false, err
					}

					return domain == "docker.io", nil
				},
				Steps: []workflow.Step{
					GetDockerHubRepository().
						WithID("repository").
						With("httpClient", httpClient).
						With("reference", data.ImageReference),
					GetDockerHubRepositoryOwner().
						WithID("owner").
						With("httpClient", httpClient).
						With("repository", workflow.Ref{Key: "step.repository.repository"}),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						repository, err := workflow.GetValue[*docker.Repository](ctx, "step.repository.repository")
						if err != nil {
							return nil, err
						}

						data.Description = repository.Description
						data.FullDescription = &models.ImageDescription{
							Markdown: repository.FullDescription,
						}
						return nil, nil
					}),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						owner, err := workflow.GetValue[*docker.Entity](ctx, "step.owner.owner")
						if err != nil {
							return nil, err
						}

						data.Image = owner.GravatarURL
						return nil, nil
					}),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						data.InsertLink(models.ImageLink{
							Type: "docker",
							URL:  docker.RepositoryUIPath(data.ImageReference),
						})
						return nil, nil
					}),
					GetDockerHubLatestVersion().
						WithID("latest").
						With("reference", data.ImageReference).
						With("httpClient", httpClient),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						reference, err := workflow.GetValue[*oci.Reference](ctx, "step.latest.reference")
						if err != nil {
							return nil, err
						}

						data.LatestReference = reference
						return nil, nil
					}),
				},
			},
			{
				ID:        "ghcr",
				Name:      "Get GHCR information",
				DependsOn: []string{"oci"},
				// Only run for GHCR images
				If: func(ctx workflow.Context) (bool, error) {
					domain, err := workflow.GetValue[string](ctx, "job.oci.step.registry.domain")
					if err != nil {
						return false, err
					}

					return domain == "ghcr.io", nil
				},
				Steps: []workflow.Step{
					GetGithubPackage().
						WithID("package").
						With("httpClient", httpClient).
						With("reference", latestOrCurrentReference(data)),
					GetGitHubDescription().
						WithID("description").
						With("httpClient", httpClient).
						With("owner", workflow.Ref{Key: "step.package.owner"}).
						With("repository", workflow.Ref{Key: "step.package.repository"}),
					GetGithubPackageREADME().
						WithID("readme").
						With("httpClient", httpClient).
						With("package", workflow.Ref{Key: "step.package.package"}),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						pkg, err := workflow.GetValue[*github.Package](ctx, "step.package.package")
						if err != nil {
							return nil, err
						}

						currentVersion, err := oci.ParseVersion(data.ImageReference.Tag)
						if err == nil && currentVersion != nil {
							for _, tag := range pkg.Tags {
								if tag.Name == "" {
									continue
								}

								newVersion, err := oci.ParseVersion(tag.Name)
								if err != nil || newVersion == nil {
									continue
								}

								if newVersion.IsCompatible(currentVersion) && newVersion.Compare(currentVersion) >= 0 {
									ref := data.ImageReference
									ref.Tag = tag.Name
									data.LatestReference = &ref
									break
								}
							}
						}

						description, err := workflow.GetValue[string](ctx, "step.description.description")
						if err != nil {
							return nil, err
						}

						data.InsertLink(models.ImageLink{
							Type: "github",
							URL:  fmt.Sprintf("%s/%s/%s", "https://github.com", url.PathEscape(pkg.Owner), url.PathEscape(pkg.Repository)),
						})

						data.InsertLink(models.ImageLink{
							Type: "ghcr",
							URL:  ghcr.PackagePath(data.ImageReference),
						})
						data.Description = description

						readme, err := workflow.GetValue[string](ctx, "step.readme.readme")
						if err == nil {
							data.FullDescription = &models.ImageDescription{
								HTML: readme,
							}
						}

						return nil, nil
					}),
				},
			},
			{
				ID:        "quay",
				Name:      "Get Quay information",
				DependsOn: []string{"oci"},
				// Only run for quay images
				If: func(ctx workflow.Context) (bool, error) {
					domain, err := workflow.GetValue[string](ctx, "job.oci.step.registry.domain")
					if err != nil {
						return false, err
					}

					return domain == "quay.io", nil
				},
				Steps: []workflow.Step{
					GetQuayLatestVersion().
						WithID("latest").
						With("reference", data.ImageReference).
						With("httpClient", httpClient),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						reference, err := workflow.GetValue[*oci.Reference](ctx, "step.latest.reference")
						if err != nil {
							return nil, err
						}

						data.LatestReference = reference
						return nil, nil
					}),
				},
			},
			// TODO: Improve this to work well for adding information for both docker,
			// and ghcr. For now, basically no docker image supports the required
			// annotations, so this might not be worth it?
			{
				ID:   "github",
				Name: "Get GitHub information",
				// Depend on whatever provides us with the latest image version
				DependsOn: []string{"oci", "docker", "ghcr"},
				// Only run for images with a reference to GitHub
				If: func(ctx workflow.Context) (bool, error) {
					if data.ImageReference.Domain == "ghcr.io" {
						return true, nil
					}

					manifests, err := workflow.GetValue[[]oci.Manifest](ctx, "job.oci.step.manifests.manifests")
					if err != nil {
						return false, err
					}

					if manifests == nil {
						return false, nil
					}

					for _, manifest := range manifests {
						if strings.Contains(manifest.SourceAnnotation(), "github.com") {
							return true, nil
						}
					}

					return false, nil
				},
				Steps: []workflow.Step{
					// workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
					// 	data.InsertTag("github")
					// 	return nil, nil
					// }),
					// // TODO: Now this only works for manifest-based images
					// GetGitHubRepsitory().
					// 	With("manifests", workflow.Ref{Key: "job.oci.step.manifests.manifests"}).
					// 	With("reference", data.LatestReference),
					// workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
					// 	endpoint, err := workflow.GetValue[string](ctx, "step.repository.endpoint")
					// 	if err != nil {
					// 		return nil, err
					// 	}

					// 	owner, err := workflow.GetValue[string](ctx, "step.repository.owner")
					// 	if err != nil {
					// 		return nil, err
					// 	}

					// 	repository, err := workflow.GetValue[string](ctx, "step.repository.name")
					// 	if err != nil {
					// 		return nil, err
					// 	}

					// 	data.InsertLink(models.ImageLink{
					// 		Type: "github",
					// 		URL:  fmt.Sprintf("%s/%s/%s", endpoint, url.PathEscape(owner), url.PathEscape(repository)),
					// 	})
					// 	return nil, nil
					// }),
					// // TODO: Get latest version based on github instead if possible
					// // TODO: Get description if not found
					// GetGitHubRelease().
					// 	WithID("release").
					// 	With("httpClient", httpClient).
					// 	With("endpoint", workflow.Ref{Key: "step.repository.endpoint"}).
					// 	With("owner", workflow.Ref{Key: "step.repository.owner"}).
					// 	With("repository", workflow.Ref{Key: "step.repository.name"}).
					// 	With("reference", data.LatestReference),
					// workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
					// 	release, err := workflow.GetValue[*github.Release](ctx, "step.release.release")
					// 	if err != nil {
					// 		return nil, err
					// 	}

					// 	if release == nil {
					// 		return nil, nil
					// 	}

					// 	data.ReleaseNotes = &models.ImageReleaseNotes{
					// 		Title: release.Title,
					// 		HTML:  release.Description,
					// 	}
					// 	return nil, nil
					// }),
				},
			},
		},
	}
}

// latestOrCurrentReference returns the latest reference contained in data if it
// is set, or the base reference if a latest reference is not (yet) found.
func latestOrCurrentReference(data *Data) oci.Reference {
	if data.LatestReference != nil {
		return *data.LatestReference
	}

	return data.ImageReference
}
