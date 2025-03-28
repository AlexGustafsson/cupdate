package oci

import (
	"encoding/json"
	"fmt"
)

var _ json.Unmarshaler = (*SBOMAttestation)(nil)

type SBOMType string

const SBOMTypeSPDX = "spdx"

// SBOMAttestation holds information gathered from an in-toto SBOM attestation
// document.
// SEE: https://github.com/in-toto/attestation.
// SEE: https://docs.docker.com/build/metadata/attestations/sbom/.
type SBOMAttestation struct {
	Type SBOMType
	SBOM string
}

func (a *SBOMAttestation) UnmarshalJSON(d []byte) error {
	var attestation struct {
		PredicateType string          `json:"predicateType"`
		Predicate     json.RawMessage `json:"predicate"`
	}
	if err := json.Unmarshal(d, &attestation); err != nil {
		return err
	}

	// NOTE: For now, we only support SPDX. We should look into supporting
	// CycloneDX as well, but they don't seem to be as prevalent in images
	if attestation.PredicateType != "https://spdx.dev/Document" {
		return fmt.Errorf("unsupported sbom attestation predicate type: %s", attestation.PredicateType)
	}

	// Without pretty printing the document, it will be intendented just as it's
	// written in its "envelope"
	sbom, err := json.MarshalIndent(attestation.Predicate, "", "  ")
	if err != nil {
		return err
	}

	res := SBOMAttestation{
		Type: SBOMTypeSPDX,
		SBOM: string(sbom),
	}

	*a = res

	return nil
}
