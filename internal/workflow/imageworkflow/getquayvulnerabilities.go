package imageworkflow

import (
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/quay"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetQuayVulnerabilities() workflow.Step {
	return workflow.Step{
		Name: "Get Quay vulnerabilities",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			reference, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			httpClient, err := workflow.GetInput[*httputil.Client](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			// NOTE: For now, to not have to perform additional queries, only look up
			// manifests that include the digest upfront
			if reference.Digest == "" {
				return nil, nil
			}

			client := &quay.Client{
				Client: httpClient,
			}

			scan, err := client.GetScan(ctx, reference)
			if err != nil {
				return nil, err
			}

			vulnerabilities := make([]models.ImageVulnerability, 0)

			if scan != nil && scan.Status == quay.ScanStatusScanned && scan.Data != nil {
				for _, feature := range scan.Data.Layer.Features {
					for _, vulnerability := range feature.Vulnerabilities {
						vulnerabilities = append(vulnerabilities, models.ImageVulnerability{
							Severity:    strings.Replace(strings.ToLower(string(vulnerability.Severity)), "unknown", "unspecified", 1),
							Authority:   "Quay",
							Description: vulnerability.Description,
							Links:       strings.Split(vulnerability.Link, " "),
						})
					}
				}
			}

			return workflow.SetOutput("vulnerabilities", vulnerabilities), nil
		},
	}
}
