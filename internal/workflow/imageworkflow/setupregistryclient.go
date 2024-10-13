package imageworkflow

import (
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/registry"
	"github.com/AlexGustafsson/cupdate/internal/registry/docker"
	"github.com/AlexGustafsson/cupdate/internal/registry/ghcr"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
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

			// TODO: Support other registries
			var client registry.Client
			switch image.Domain {
			case "docker.io":
				client = &docker.Client{
					Client: httpClient,
				}
			case "ghcr.io":
				client = &ghcr.Client{
					Client: httpClient,
				}
			default:
				return nil, fmt.Errorf("unsupported registry domain: %s", image.Domain)
			}

			return workflow.Batch(
				workflow.SetOutput("client", client),
				workflow.SetOutput("domain", "https://"+image.Domain),
			), nil
		},
	}
}
