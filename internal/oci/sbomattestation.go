package oci

import (
	"encoding/json"
)

var _ json.Unmarshaler = (*SBOMAttestation)(nil)

// SBOMAttestation holds information gathered from an in-toto SBOM attestation
// document.
// SEE: https://github.com/in-toto/attestation.
// SEE: https://docs.docker.com/build/metadata/attestations/sbom/.
type SBOMAttestation struct {
	PredicateType string
	SBOM          string
}

func (a *SBOMAttestation) UnmarshalJSON(d []byte) error {
	var attestation struct {
		PredicateType string          `json:"predicateType"`
		Predicate     json.RawMessage `json:"predicate"`
	}
	if err := json.Unmarshal(d, &attestation); err != nil {
		return err
	}

	// Without pretty printing the document, it will be intendented just as it's
	// written in its "envelope"
	sbom, err := json.MarshalIndent(attestation.Predicate, "", "  ")
	if err != nil {
		return err
	}

	res := SBOMAttestation{
		PredicateType: attestation.PredicateType,
		SBOM:          string(sbom),
	}

	*a = res

	return nil
}
