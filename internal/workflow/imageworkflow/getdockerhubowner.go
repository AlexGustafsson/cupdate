package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/registry/docker"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetDockerHubOwner() workflow.Step {
	return workflow.Step{
		Name: "Get Docker Hub owner",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			httpClient, err := workflow.GetInput[*httputil.Client](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			repository, err := workflow.GetInput[*docker.Repository](ctx, "repository", true)
			if err != nil {
				return nil, err
			}

			client := &docker.Client{
				Client: httpClient,
			}

			entity, err := client.GetOrganizationOrUser(ctx, repository.Namespace)
			if err != nil {
				return nil, err
			}

			return workflow.SetOutput("owner", entity), nil
		},
	}
}
