package imageworkflow

import (
	"errors"

	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetManifest() workflow.Step {
	return workflow.Step{
		Name: "Get manifest",

		Main: func(ctx workflow.Context) (workflow.Command, error) {
			registryClient, err := workflow.GetInput[*oci.Client](ctx, "registryClient", true)
			if err != nil {
				return nil, err
			}

			ref, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if errors.Is(err, workflow.ErrInvalidType) {
				imageRef, err := workflow.GetInput[*oci.Reference](ctx, "reference", true)
				if err != nil {
					return nil, err
				} else if imageRef != nil {
					ref = *imageRef
				}
			} else if err != nil {
				return nil, err
			}

			manifest, err := registryClient.GetManifest(ctx, ref)
			if err != nil {
				return nil, err
			}

			return workflow.SetOutput("manifest", manifest), nil
		},
	}
}
