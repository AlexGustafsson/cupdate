package imageworkflow

import (
	"errors"
	"maps"
	"os"
	"slices"

	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
	"github.com/google/osv-scanner/v2/pkg/osvscanner"
)

func ScanSBOM() workflow.Step {
	return workflow.Step{
		Name: "Scan SBOM",

		Main: func(ctx workflow.Context) (workflow.Command, error) {
			attestations, err := workflow.GetInput[map[string]oci.SBOMAttestation](ctx, "attestations", true)
			if err != nil {
				return nil, err
			}

			vulnerabilities := make([]models.ImageVulnerability, 0)
			seen := make(map[string]struct{})
			for _, attestation := range attestations {
				// NOTE: For now, we only support SPDX
				file, err := os.CreateTemp("", "cupdate-scan-sbom-*.spdx.json")
				if err != nil {
					return nil, err
				}

				_, err = file.WriteString(attestation.SBOM)
				file.Close()
				if err != nil {
					os.Remove(file.Name())
					return nil, err
				}

				actions := osvscanner.ScannerActions{
					SBOMPaths: []string{file.Name()},
				}

				results, err := osvscanner.DoScan(actions)
				os.Remove(file.Name())
				if err != nil && !errors.Is(err, osvscanner.ErrVulnerabilitiesFound) {
					return nil, err
				}

				for _, result := range results.Results {
					for _, pkg := range result.Packages {
						for _, vulnerability := range pkg.Vulnerabilities {
							_, alreadySeen := seen[vulnerability.ID]
							if !alreadySeen {
								seen[vulnerability.ID] = struct{}{}

								links := make(map[string]struct{}, 0)
								for _, reference := range vulnerability.References {
									links[reference.URL] = struct{}{}
								}

								vulnerabilities = append(vulnerabilities, models.ImageVulnerability{
									ID: vulnerability.ID,
									// TODO: For now I haven't seen that many cases of severity
									// being specified...
									Severity:    models.SeverityUnspecified,
									Authority:   "OSV",
									Description: vulnerability.Summary,
									Links:       slices.Collect(maps.Keys(links)),
								})
							}
						}
					}
				}

			}

			return workflow.SetOutput("vulnerabilities", vulnerabilities), nil
		},
	}
}
