package imageworkflow

import (
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/gitlab"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetGitLabDescription() workflow.Step {
	return workflow.Step{
		Name: "Get a repository's description from GitLab",
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

			// The repository path is <owner>/<group>/<project>
			parts := strings.Split(reference.Path, "/")
			if len(parts) < 3 {
				return nil, nil
			}

			fullPath := strings.Join(parts[0:3], "/")

			description, err := client.GetRepositoryDescription(ctx, fullPath)
			if err != nil {
				return nil, err
			}

			return workflow.SetOutput("description", description), nil
		},
	}
}
