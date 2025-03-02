package imageworkflow

import (
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

			httpClient, err := workflow.GetInput[httputil.Requester](ctx, "httpClient", true)
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

			scan, err := client.GetVulnerabilities(ctx, reference)
			if err != nil {
				return nil, err
			}

			vulnerabilities := make([]models.ImageVulnerability, 0)
			for _, vulnerability := range scan {
				vulnerabilities = append(vulnerabilities, models.ImageVulnerability{
					Severity:    string(vulnerability.Severity),
					Authority:   "Quay",
					Description: vulnerability.Description,
					Links:       vulnerability.Links,
				})
			}

			return workflow.SetOutput("vulnerabilities", vulnerabilities), nil
		},
	}
}
