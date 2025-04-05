package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/httputil"
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

			vulnerabilities, err := client.GetVulnerabilities(ctx, reference)
			if err != nil {
				return nil, err
			}

			return workflow.SetOutput("vulnerabilities", vulnerabilities), nil
		},
	}
}
