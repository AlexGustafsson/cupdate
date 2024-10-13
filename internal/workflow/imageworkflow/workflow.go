package imageworkflow

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/github"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/registry/docker"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

// TODO: Let each step take an optional cache value instead. If it exists,
// perform request caching of all 200 responses instead - more like is stated
// in the cache code. That way we naturally always parse the values and don't
// have to care about special types in the same way. The cache is always just
// bytes - no need to cache JSON (yay!). Learning: Typed caches suck. Cache
// repsonses instead.

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

						if reference == nil {
							return nil, nil
						}

						data.LatestReference = *reference
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
					// TODO
				},
			},
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
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						data.InsertTag("github")
						return nil, nil
					}),
					GetGitHubRepsitory().
						WithID("repository").
						With("httpClient", httpClient).
						With("manifests", workflow.Ref{Key: "job.oci.step.manifests.manifests"}).
						With("reference", data.LatestReference),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						endpoint, err := workflow.GetValue[string](ctx, "step.repository.endpoint")
						if err != nil {
							return nil, err
						}

						owner, err := workflow.GetValue[string](ctx, "step.repository.owner")
						if err != nil {
							return nil, err
						}

						repository, err := workflow.GetValue[string](ctx, "step.repository.name")
						if err != nil {
							return nil, err
						}

						data.InsertLink(models.ImageLink{
							Type: "github",
							URL:  fmt.Sprintf("%s/%s/%s", endpoint, url.PathEscape(owner), url.PathEscape(repository)),
						})
						return nil, nil
					}),
					// TODO: Get latest version based on github instead if possible
					// TODO: Set short description from repository if not already exists
					GetGitHubRelease().
						WithID("release").
						With("httpClient", httpClient).
						With("endpoint", workflow.Ref{Key: "step.repository.endpoint"}).
						With("owner", workflow.Ref{Key: "step.repository.owner"}).
						With("repository", workflow.Ref{Key: "step.repository.name"}).
						With("reference", data.LatestReference),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						release, err := workflow.GetValue[*github.Release](ctx, "step.release.release")
						if err != nil {
							return nil, err
						}

						if release == nil {
							return nil, nil
						}

						data.ReleaseNotes = &models.ImageReleaseNotes{
							Title: release.Title,
							HTML:  release.Description,
						}
						return nil, nil
					}),
				},
			},
		},
	}
}
