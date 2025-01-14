package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/dockerhub"
	"github.com/AlexGustafsson/cupdate/internal/ghcr"
	"github.com/AlexGustafsson/cupdate/internal/gitlab"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func SetupRegistryClient() workflow.Step {
	return workflow.Step{
		Name: "Setup registry client",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			httpClient, err := workflow.GetInput[*httputil.Client](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			image, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			registryAuth, err := workflow.GetInput[*httputil.AuthMux](ctx, "registryAuth", true)
			if err != nil {
				return nil, err
			}

			// TODO: Support the www-authenticate return header mandated by the OCI
			// distribution spec. Would help support currently unknown registries that are
			// well-behaved
			baseAuth := httputil.NewAuthMux()
			baseAuth.Handle("*.docker.io", &dockerhub.Client{
				Client: httpClient,
			})
			baseAuth.Handle("ghcr.io", &ghcr.Client{
				Client: httpClient,
			})
			// Linux Server mirrors images, but the default is GitHub and I've never
			// seen any other backend being used. For now, assume GitHub
			baseAuth.Handle("lscr.io", &ghcr.Client{
				Client: httpClient,
			})
			baseAuth.Handle("registry.gitlab.com", &gitlab.Client{
				Client: httpClient,
			})

			// Apply user configuration
			baseAuth.Copy(registryAuth)

			client := &oci.Client{
				Client:   httpClient,
				AuthFunc: baseAuth.HandleAuth,
			}

			return workflow.Batch(
				workflow.SetOutput("client", client),
				workflow.SetOutput("domain", image.Domain),
			), nil
		},
	}
}
