package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/registry/gitlab"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetGitLabTags() workflow.Step {
	return workflow.Step{
		Name: "Get tags from GitLab",
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

			tags, err := client.GetTags(ctx, reference)
			if err != nil {
				return nil, err
			}

			return workflow.SetOutput("tags", tags), nil
		},
	}
}
