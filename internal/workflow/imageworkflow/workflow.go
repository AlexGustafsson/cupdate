package imageworkflow

import (
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/registry/docker"
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
					SetupRegistryClient().WithID("registry"),
					InsertTag().
						With("data", data).
						With("tag", workflow.Ref("steps.registry.domain")),
					workflow.Cache[[]oci.Manifest]().
						WithID("manifestsCache").
						With("cache", cache).
						With("cacheKey", "abcd").
						With("valueKey", "manifests"),
					GetManifests().
						WithID("manifests").
						WithCondition("step.manifestsCache.miss").
						With("registry", workflow.Ref("steps.registry.client")).
						With("reference", data.ImageReference),
					workflow.StoreValue().
						With("name", "manifests").
						With("value", workflow.Ref("steps.manifests.manifests")),
					// TODO:
					InsertLink().
						With("data", data).
						With("link", models.ImageLink{}),
				},
			},
			{
				ID:        "docker",
				Name:      "Get Docker Hub information",
				DependsOn: []string{"oci"},
				If: func(ctx workflow.Context) (bool, error) {
					// TODO: Is docker registry
					return true, nil
				},
				Steps: []workflow.Step{
					workflow.Cache(cache, "abcd", "repository", 24*time.Hour).
						WithID("repositoryCache"),
					GetDockerHubRepository().
						InsertDescription(data, func(ctx workflow.Context) *models.ImageDescription {
							return nil
						}),
					// TODO:
					InsertLink().
						With("data", data).
						With("link", models.ImageLink{}),
					workflow.Cache[[]docker.Tag]().
						WithID("tagsCache").
						With("cache", cache).
						With("cacheKey", "abcd").
						With("valueKey", "tags"),
					workflow.Cache(cache, "abcd", "tags", 24*time.Hour).
						WithID("tagsCache"),
					GetDockerHubTags().
						WithID("tags").
						WithCondition("steps.tagsCache.miss").
						With("reference", data.ImageReference),
					workflow.StoreValue().
						With("name", "tags").
						With("value", workflow.Ref("step.tags.tags")),
					// TODO:
					InsertLatestVersion().
						With("data", data).
						With("reference", nil),
				},
			},
			{
				ID:   "github",
				Name: "Get release from GitHub",
				// Depend on whatever provides us with the latest image version
				DependsOn: []string{"oci", "docker"},
				If: func(ctx workflow.Context) (bool, error) {
					// TODO: Has GitHub link
					return true, nil
				},
				Steps: []workflow.Step{
					workflow.Cache(cache, "abcd", "release", 24*time.Hour).
						WithID("releaseCache"),
					GetGitHubRelease().
						WithCondition("step.releaseCache.miss").
						WithID("release"),
					InsertReleaseNotes(data, func(ctx workflow.Context) *models.ImageReleaseNotes {
						return nil
					}),
				},
			},
		},
	}
}
