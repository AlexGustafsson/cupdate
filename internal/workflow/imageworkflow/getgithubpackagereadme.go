package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/github"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetGithubPackageREADME() workflow.Step {
	return workflow.Step{
		Name: "Get GitHub package README",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			httpClient, err := workflow.GetInput[httputil.Requester](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			client := &github.Client{
				Endpoint: "https://github.com",
				Client:   httpClient,
			}

			pkg, err := workflow.GetInput[*github.Package](ctx, "package", true)
			if err != nil {
				return nil, err
			}

			if pkg == nil {
				return nil, nil
			}

			if pkg.ReadmeURL == "" {
				return nil, nil
			}

			readme, err := client.GetREADME(ctx, pkg.ReadmeURL)
			if err != nil {
				return nil, err
			}
			return workflow.Batch(
				workflow.SetOutput("readme", string(readme)),
			), nil
		},
	}
}
