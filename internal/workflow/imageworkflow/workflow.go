package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
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
					GetDockerHubRepository().
						WithID("repository").
						With("httpClient", httpClient).
						With("reference", data.Image),
					InsertDescription().
						With("data", data).
						With("description", models.ImageDescription{ /*TODO*/ }),
					InsertLink().
						With("data", data).
						With("link", models.ImageLink{ /*TODO*/ }),
					GetDockerHubTags().
						WithID("tags").
						WithCondition("step.tagsCache.miss").
						With("reference", data.ImageReference).
						With("httpClient", httpClient),
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
					GetGitHubRelease().
						WithCondition("step.releaseCache.miss").
						WithID("release").
						With("httpClient", httpClient),
					InsertReleaseNotes().
						With("data", data).
						With("releaseNotes", models.ImageReleaseNotes{ /*TODO*/ }),
				},
			},
		},
	}
}
