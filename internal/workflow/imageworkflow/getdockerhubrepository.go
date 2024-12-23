package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/dockerhub"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetDockerHubRepository() workflow.Step {
	return workflow.Step{
		Name: "Get Docker Hub repository",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			reference, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			httpClient, err := workflow.GetInput[*httputil.Client](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			client := &dockerhub.Client{
				Client: httpClient,
			}

			repository, err := client.GetRepository(ctx, reference)
			if err != nil {
				return nil, err
			}

			return workflow.SetOutput("repository", repository), nil
		},
	}
}
