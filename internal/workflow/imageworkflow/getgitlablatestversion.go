package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/registry/gitlab"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetGitLabLatestVersion() workflow.Step {
	return workflow.Step{
		Name: "Get latest version from GitLab",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			httpClient, err := workflow.GetInput[*httputil.Client](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			reference, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			client := &gitlab.Client{Client: httpClient}

			image, err := client.GetLatestVersion(ctx, reference)
			if err != nil {
				return nil, err
			}

			if image == nil {
				return workflow.SetOutput("reference", (*oci.Reference)(nil)), nil
			}

			return workflow.SetOutput("reference", &image.Name), nil
		},
	}
}