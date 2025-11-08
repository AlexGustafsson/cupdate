package imageworkflow

import (
	"errors"
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/osv"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func ScanSBOM() workflow.Step {
	return workflow.Step{
		Name: "Scan SBOM",

		Main: func(ctx workflow.Context) (workflow.Command, error) {
			attestations, err := workflow.GetInput[map[string]oci.SBOMAttestation](ctx, "attestations", true)
			if err != nil {
				return nil, err
			}

			vulnerabilities := make([]osv.Vulnerability, 0)
			seen := make(map[string]struct{})

			for _, attestation := range attestations {
				// For now, we only support spdx
				vulns, err := osv.ScanSPDX(ctx, attestation.SBOM)
				if errors.Is(err, osv.ErrScannerNotFound) {
					return nil, fmt.Errorf("osv-scanner is not installed")
				} else if err != nil {
					return nil, err
				}

				for _, vulnerability := range vulns {
					// Assume same information available in each instance
					if _, ok := seen[vulnerability.ID]; ok {
						continue
					}

					seen[vulnerability.ID] = struct{}{}

					vulnerabilities = append(vulnerabilities, vulnerability)
				}
			}

			return workflow.SetOutput("vulnerabilities", vulnerabilities), nil
		},
	}
}
