package imageworkflow

import (
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func New(cache cache.Cache, data *Data) error {
	w := workflow.Workflow{
		Name: "Process image",
		Jobs: []workflow.Job{
			{
				ID:   "oci",
				Name: "Get OCI information",
				Steps: []workflow.Step{
					SetupRegistryClient(data).WithID("registry"),
					InsertTags(data, func(ctx workflow.Context) []string {
						domain, ok := ctx.Output("steps.registry.domain")
						if !ok {
							return nil
						}

						return []string{domain.(string)}
					}),
					workflow.UnlessCached(
						cache,
						"abcd",
						24*time.Hour,
						GetManifests().
							WithInput("client", "steps.registry.client").
							WithInput("image", "") // TODO: How do we handle inputs like these?
							WithID("manifests"),
					),
					workflow.Cache(cache, "abcd", func(ctx workflow.Context) (any, bool) {
						return ctx.Output("steps.manifests.manifests")
					}),
					InsertLinks(data, func(ctx workflow.Context) []models.ImageLink {
						return []models.ImageLink{}
					}),
				},
			},
			{
				ID:        "docker",
				Name:      "Get Docker Hub information",
				DependsOn: []string{"oci"},
				ShouldRun: func(ctx workflow.Context) (bool, error) {
					// Is docker registry
					return false, nil
				},
				Steps: []workflow.Step{
					workflow.Cache(cache, "abcd", "repository.repository"),
					workflow.UnlessCached(
						cache,
						"abcd",
						GetDockerHubRepository().WithID("repository"),
					),
					InsertDescription(data, func(ctx workflow.Context) *models.ImageDescription {
						return nil
					}),
					InsertLinks(data, func(ctx workflow.Context) []models.ImageLink {
						return []models.ImageLink{}
					}),
					workflow.Cache(cache, "abcd", "tags.tags"),
					workflow.UnlessCached(
						cache,
						"abcd",
						GetDockerHubTags().WithID("tags"),
					),
					InsertLatestVersion(data, func(ctx workflow.Context) *oci.Reference {
						return data.LatestVersion
					}),
				},
			},
			{
				ID:   "github",
				Name: "Get release from GitHub",
				// Depend on whatever provides us with the latest image version
				DependsOn: []string{"oci", "docker"},
				ShouldRun: func(ctx workflow.Context) (bool, error) {
					_, ok := ctx.Output("jobs.oci.steps.manifest")
					// TODO: only if github link?
					return ok, nil
				},
				Steps: []workflow.Step{
					workflow.Cache(cache, "abcd", "steps.release.release"),
					workflow.UnlessCached(
						cache,
						"abcd",
						GetGitHubRelease("jobs.oci.steps.manifest").WithID("release"),
					),
					InsertReleaseNotes(data, func(ctx workflow.Context) *models.ImageReleaseNotes {
						return nil
					}),
				},
			},
		},
	}
}
