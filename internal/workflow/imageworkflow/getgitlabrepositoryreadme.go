package imageworkflow

import (
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/gitlab"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetGitLabRepositoryREADME() workflow.Step {
	return workflow.Step{
		Name: "Get a repository's README from GitLab",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			httpClient, err := workflow.GetInput[httputil.Requester](ctx, "httpClient", true)
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

			blob, err := client.GetRepositoryREADMEBlob(ctx, fullPath)
			if err != nil {
				return nil, err
			} else if blob == nil {
				return nil, nil
			}

			return workflow.Batch(
				workflow.SetOutput("mime", blob.MimeType),
				workflow.SetOutput("raw", blob.Raw),
				workflow.SetOutput("html", blob.HTML),
			), nil
		},
	}
}
