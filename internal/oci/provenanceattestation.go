package oci

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

var _ json.Unmarshaler = (*ProvenanceAttestation)(nil)

// ProvenanceAttestation holds information gathered from an in-toto provenance
// attestation document containing moby buildkit metadata.
// SEE: https://github.com/in-toto/attestation.
// SEE: https://docs.docker.com/build/metadata/attestations/slsa-provenance/.
type ProvenanceAttestation struct {
	BuildStartedOn  time.Time
	BuildFinishedOn time.Time
	// Source is the VCS source containing the code.
	Source string
	// SourceRevision is the revision (typically sha-1) of the version built.
	SourceRevision string
	// Dockerfile contains the full Dockerfile of the image, if possible.
	Dockerfile string
}

func (a *ProvenanceAttestation) UnmarshalJSON(d []byte) error {
	var attestation struct {
		Predicate struct {
			Metadata struct {
				BuildStartedOn   time.Time `json:"buildStartedOn"`
				BuildFinishedOn  time.Time `json:"buildFinishedOn"`
				BuildKitMetadata *struct {
					VCS struct {
						Source   string `json:"source"`
						Revision string `json:"revision"`
					} `json:"vcs"`
					Source struct {
						Infos []struct {
							Language string `json:"language"`
							Data     string `json:"data"`
						} `json:"infos"`
					} `json:"source"`
				} `json:"https://mobyproject.org/buildkit@v1#metadata"`
			} `json:"metadata"`
		} `json:"predicate"`
	}
	if err := json.Unmarshal(d, &attestation); err != nil {
		return err
	}

	res := ProvenanceAttestation{
		BuildStartedOn:  attestation.Predicate.Metadata.BuildStartedOn,
		BuildFinishedOn: attestation.Predicate.Metadata.BuildFinishedOn,
	}

	if meta := attestation.Predicate.Metadata.BuildKitMetadata; meta != nil {
		res.Source = meta.VCS.Source
		res.SourceRevision = meta.VCS.Revision

		for _, infos := range meta.Source.Infos {
			if infos.Language == "Dockerfile" {
				content, err := base64.StdEncoding.DecodeString(infos.Data)
				if err != nil {
					continue
				}

				res.Dockerfile = string(content)
				break
			}
		}
	}

	*a = res

	return nil
}
