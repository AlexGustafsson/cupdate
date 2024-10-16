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
			// If not, try to find references to GitHub and go from there
			manifests, err := workflow.GetInput[[]oci.Manifest](ctx, "manifests", true)
			if err != nil {
				return nil, err
			}

			if manifests == nil {
				return nil, fmt.Errorf("no manifests found")
			}

			var endpoint, owner, repository string
			var ok bool
			for _, manifest := range manifests {
				if strings.Contains(manifest.SourceAnnotation(), "github.com") {
					endpoint, owner, repository, _, ok = github.ParseURL(manifest.SourceAnnotation())
					// NOTE: Only support github.com for now
					if ok && endpoint == "https://github.com" {
						break
					}
				}
			}
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
