package imageworkflow

import (
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/registry/docker"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func SetupRegistryClient() workflow.Step {
	return workflow.Step{
		Name: "Setup registry client",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			image, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			// TODO: Support other registries
			if image.Domain != "docker.io" {
				return nil, fmt.Errorf("unsupported registry domain: %s", image.Domain)
			}
			client := &docker.Client{}

			return workflow.Batch(
				workflow.SetOutput("client", client),
				workflow.SetOutput("domain", image.Domain),
			), nil
		},
	}
}
