package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/github"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetGitHubDescription() workflow.Step {
	return workflow.Step{
		Name: "Get GitHub description",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			httpClient, err := workflow.GetInput[httputil.Requester](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			owner, err := workflow.GetInput[string](ctx, "owner", true)
			if err != nil {
				return nil, err
			}

			repository, err := workflow.GetInput[string](ctx, "repository", true)
			if err != nil {
				return nil, err
			}

			client := &github.Client{
				Endpoint: "https://github.com",
				Client:   httpClient,
			}

			description, err := client.GetDescription(ctx, owner, repository)
			if err != nil {
				return nil, err
			}

			return workflow.SetOutput("description", description), nil
		},
	}
}
