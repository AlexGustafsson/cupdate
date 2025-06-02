package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func DetermineSource() workflow.Step {
	return workflow.Step{
		Name: "Determine source",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			reference, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			annotations, err := workflow.GetInput[oci.Annotations](ctx, "annotations", false)
			if err != nil {
				return nil, err
			}

			attestations, err := workflow.GetInput[map[string]oci.ProvenanceAttestation](ctx, "attestations", false)
			if err != nil {
				return nil, err
			}

			registry := reference.Domain
			repository := ""

			// Identify repository from OCI annotations. Prioritize these as they are
			// typically manually set by the image author
			if repository == "" && annotations != nil {
				if uri := annotations.Source(); uri != "" {
					repository = uri
				}
			}

			// Fall back to identifying the repository from provenance attestations
			if repository == "" && attestations != nil {
				for _, attestation := range attestations {
					if uri := attestation.Source; uri != "" {
						repository = uri
						break
					}
				}
			}

			return workflow.Batch(
				workflow.SetOutput("registry", registry),
				workflow.SetOutput("repository", repository),
			), nil
		},
	}
}
