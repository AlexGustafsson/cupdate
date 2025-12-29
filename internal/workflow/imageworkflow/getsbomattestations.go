package imageworkflow

import (
	"encoding/json"
	"errors"

	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetSBOMAttestations() workflow.Step {
	return workflow.Step{
		Name: "Get SBOM attestations",

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

			manifests, err := workflow.GetInput[map[string]*oci.ImageManifest](ctx, "manifests", true)
			if err != nil {
				return nil, err
			}

			// TODO: Instead of getting the attestations for all images (typically the
			// case for multi-arch images), we could use the host / node information
			// from the graph to only get data for the architectures in use
			attestations := make(map[string]oci.SBOMAttestation)
			for manifestDigest, attestationManifest := range manifests {
				_, sbomBlobDigest, ok := attestationManifest.SBOMLayerDigest()
				if !ok {
					return nil, nil
				}

				blob, err := registryClient.GetBlob(ctx, image, sbomBlobDigest, true)
				if err != nil {
					return nil, err
				}
				defer blob.Close()

				var attestation oci.SBOMAttestation
				if err := json.NewDecoder(blob).Decode(&attestation); err != nil {
					return nil, err
				}

				// NOTE: Docker Hardened Images reports CyclonedX, but its predicate is
				// just null, ignore them
				if attestation.SBOM == "null" {
					continue
				}

				attestations[manifestDigest] = attestation
			}

			return workflow.SetOutput("attestations", attestations), nil
		},
	}
}
