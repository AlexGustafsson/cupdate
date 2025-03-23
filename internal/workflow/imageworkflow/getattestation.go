package imageworkflow

import (
	"encoding/json"
	"errors"

	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetAttestation() workflow.Step {
	return workflow.Step{
		Name: "Get attestation",

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

			manifest, err := workflow.GetInput[any](ctx, "manifest", true)
			if err != nil {
				return nil, err
			}

			index, ok := manifest.(*oci.ImageIndex)
			if !ok {
				return nil, nil
			}

			if !index.HasAttestationManifest() {
				return nil, nil
			}

			// TODO: Instead of getting the attestations for all images (typically the
			// case for multi-arch images), we could use the host / node information
			// from the graph to only get data for the architectures in use
			attestations := make(map[string]oci.Attestation)
			for manifestDigest, attestationManifestDigest := range index.AttestationManifestDigest() {
				attestationManifest, err := registryClient.GetAttestationManifest(ctx, image, attestationManifestDigest)
				if err != nil {
					return nil, err
				}

				_, provenanceBlobDigest, ok := attestationManifest.ProvenanceDigest()
				if !ok {
					return nil, nil
				}

				blob, err := registryClient.GetBlob(ctx, image, provenanceBlobDigest, true)
				if err != nil {
					return nil, err
				}

				var attestation oci.Attestation
				if err := json.NewDecoder(blob).Decode(&attestation); err != nil {
					return nil, err
				}

				attestations[manifestDigest] = attestation
			}

			return workflow.SetOutput("attestations", attestations), nil
		},
	}
}
