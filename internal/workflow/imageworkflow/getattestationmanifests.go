package imageworkflow

import (
	"errors"

	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetAttestationManifests() workflow.Step {
	return workflow.Step{
		Name: "Get attestation manifests",

		Main: func(ctx workflow.Context) (workflow.Command, error) {
			registryClient, err := workflow.GetInput[*oci.Client](ctx, "registryClient", true)
			if err != nil {
				return nil, err
			}

			image, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if errors.Is(err, workflow.ErrInvalidType) {
				imageRef, err := workflow.GetInput[*oci.Reference](ctx, "reference", true)
				if err != nil {
					return nil, err
				} else if imageRef != nil {
					image = *imageRef
				}
			} else if err != nil {
				return nil, err
			}

			// TODO: We could use the image graph to determine what architectures to
			// keep, right now we fetch all attestation manifests
			attestationManifests, err := registryClient.GetAttestationManifests(ctx, image, nil)
			if err != nil {
				return nil, err
			}

			return workflow.SetOutput("manifests", attestationManifests), nil
		},
	}
}
