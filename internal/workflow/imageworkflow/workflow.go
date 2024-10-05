package imageworkflow

import (
	"fmt"
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
					InsertTag().
						With("data", data).
						With("tag", workflow.Ref{Key: "step.registry.domain"}),
					GetManifests().
						WithID("manifests").
						With("registryClient", workflow.Ref{Key: "step.registry.client"}).
						With("reference", data.ImageReference),
					InsertLink().
						With("data", data).
						With("link", workflow.Func{Func: func(ctx workflow.Context) (any, error) {
							domain, err := workflow.GetValue[string](ctx, "step.registry.domain")
							if err != nil {
								return nil, err
							}

							return models.ImageLink{
								Type: "oci-registry",
								URL:  domain,
							}, nil
						}}),
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
					GetDockerHubOwner().
						WithID("owner").
						With("httpClient", httpClient).
						With("reference", data.ImageReference),
					InsertDescription().
						With("data", data).
						With("description", workflow.Func{Func: func(ctx workflow.Context) (any, error) {
							repository, err := workflow.GetValue[*docker.Repository](ctx, "step.repository.repository")
							if err != nil {
								return nil, err
							}

							return &models.ImageDescription{
								Markdown: repository.FullDescription,
							}, nil
						}}),
					InsertLink().
						With("data", data).
						With("link", models.ImageLink{
							Type: "docker",
							URL:  docker.RepositoryPath(data.ImageReference),
						}),
					GetDockerHubLatestVersion().
						WithID("latest").
						With("reference", data.ImageReference).
						With("httpClient", httpClient),
					InsertLatestVersion().
						With("data", data).
						With("reference", workflow.Ref{Key: "step.latest.reference"}),
				},
			},
			{
				ID:   "github",
				Name: "Get GitHub information",
				// Depend on whatever provides us with the latest image version
				DependsOn: []string{"oci", "docker"},
				// Only run for images with a reference to GitHub
				If: func(ctx workflow.Context) (bool, error) {
					manifests, err := workflow.GetValue[[]oci.Manifest](ctx, "job.oci.step.manifests.manifests")
					if err != nil {
						return false, err
					}

					if manifests == nil {
						return false, nil
					}

					for _, manifest := range manifests {
						if strings.Contains(manifest.SourceAnnotation(), "github.com") {
							fmt.Println(manifest.SourceAnnotation())
							return true, nil
						}
					}

					return false, nil
				},
				Steps: []workflow.Step{
					InsertTag().
						With("data", data).
						With("tag", "github"),
					InsertLinks().
						With("data", data).
						With("links", workflow.Func{Func: func(ctx workflow.Context) (any, error) {
							manifests, err := workflow.GetValue[[]oci.Manifest](ctx, "job.oci.step.manifests.manifests")
							if err != nil {
								return []models.ImageLink{}, nil
							}

							if manifests == nil {
								return []models.ImageLink{}, nil
							}

							links := make([]models.ImageLink, 0)
							for _, manifest := range manifests {
								if strings.Contains(manifest.SourceAnnotation(), "github.com") {
									links = append(links, models.ImageLink{
										Type: "github",
										URL:  manifest.SourceAnnotation(),
									})
								}
							}

							return links, nil
						}}),
					GetGitHubRelease().
						WithID("release").
						With("httpClient", httpClient).
						With("manifests", workflow.Ref{Key: "job.oci.step.manifests.manifests"}).
						With("reference", data.ImageReference),
					InsertReleaseNotes().
						With("data", data).
						With("releaseNotes", workflow.Func{Func: func(ctx workflow.Context) (any, error) {
							release, err := workflow.GetValue[*github.Release](ctx, "step.release.release")
							if err != nil {
								return nil, err
							}

							if release == nil {
								return (*models.ImageReleaseNotes)(nil), nil
							}

							return &models.ImageReleaseNotes{
								Title: release.Title,
								HTML:  release.Description,
							}, nil
						}}),
				},
			},
		},
	}
}
