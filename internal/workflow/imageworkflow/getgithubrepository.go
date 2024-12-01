package imageworkflow

import (
	"fmt"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/github"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetGitHubRepsitory() workflow.Step {
	return workflow.Step{
		Name: "Get GitHub repository",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			annotations, err := workflow.GetValue[oci.Annotations](ctx, "step.annotations.annotations")
			if err != nil {
				return nil, err
			}

			source := annotations.Source()
			if !strings.Contains(source, "://github.com/") {
				return nil, fmt.Errorf("no GitHub references found")
			}

			var endpoint, owner, repository string
			endpoint, owner, repository, _, ok := github.ParseURL(source)
			// NOTE: Only support github.com for now
			if ok {
				return workflow.Batch(
					workflow.SetOutput("endpoint", endpoint),
					workflow.SetOutput("owner", owner),
					workflow.SetOutput("name", repository),
				), nil
			}

			return nil, fmt.Errorf("no GitHub references found")
		},
	}
}
