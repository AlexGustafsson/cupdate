package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/registry"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetManifests() workflow.Step {
	return workflow.Step{
		Name: "Get manifests",

		Main: func(ctx workflow.Context) (workflow.Command, error) {
			registryClient, err := workflow.GetInput[registry.Client](ctx, "registryClient", true)
			if err != nil {
				return nil, err
			}

			image, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			manifests, err := registryClient.GetManifests(ctx, image)
			if err != nil {
				return nil, err
			}

			return workflow.SetOutput("manifests", manifests), nil
		},
	}
}