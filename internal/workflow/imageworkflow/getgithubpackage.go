package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/github"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetGithubPackage() workflow.Step {
	return workflow.Step{
		Name: "Get GitHub package",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			httpClient, err := workflow.GetInput[*httputil.Client](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			reference, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			client := &github.Client{
				Endpoint: "https://github.com",
				Client:   httpClient,
			}

			pkg, err := client.GetPackage(ctx, reference)
			if err != nil {
				return nil, err
			}

			return workflow.Batch(
				workflow.SetOutput("package", pkg),
				workflow.SetOutput("owner", pkg.Owner),
				workflow.SetOutput("repository", pkg.Repository),
			), nil
		},
	}
}
