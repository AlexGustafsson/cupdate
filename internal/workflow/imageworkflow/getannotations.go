package imageworkflow

import (
	"runtime"

	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetAnnotations() workflow.Step {
	return workflow.Step{
		Name: "Get annotations",

		Main: func(ctx workflow.Context) (workflow.Command, error) {
			registryClient, err := workflow.GetInput[*oci.Client](ctx, "registryClient", true)
			if err != nil {
				return nil, err
			}

			image, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			manifests, err := workflow.GetInput[[]oci.Manifest](ctx, "manifests", true)
			if err != nil {
				return nil, err
			}

			annotations, err := registryClient.GetAnnotations(ctx, image, &oci.GetAnnotationsOptions{
				Manifests:    manifests,
				Architecture: runtime.GOARCH,
			})
			if err != nil {
				return nil, err
			}

			return workflow.SetOutput("annotations", annotations), nil
		},
	}
}
