package imageworkflow

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/github"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetGitHubRepsitory() workflow.Step {
	return workflow.Step{
		Name: "Get GitHub repository",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			source, err := workflow.GetInput[string](ctx, "repository", true)
			if err != nil {
				return nil, err
			}

			pkg, err := workflow.GetInput[*github.Package](ctx, "package", false)
			if err != nil {
				return nil, err
			}

			// Find the repository's URL. Prefer the GHCR package's repository, if set
			// otherwise try to use the URL previously identified
			if pkg != nil {
				source = fmt.Sprintf("https://github.com/%s/%s", url.PathEscape(pkg.Owner), url.PathEscape(pkg.Repository))
			}

			if !strings.Contains(source, "://github.com/") {
				return nil, fmt.Errorf("no valid GitHub reference found")
			}

			var endpoint, owner, repository string
			endpoint, owner, repository, _, ok := github.ParseURL(source)
			// NOTE: Only support github.com for now
			if ok {
				return workflow.Batch(
					workflow.SetOutput("endpoint", endpoint),
					workflow.SetOutput("owner", owner),
					workflow.SetOutput("name", repository),
					workflow.SetOutput("repository", fmt.Sprintf("%s/%s/%s", strings.TrimPrefix(endpoint, "https://"), owner, repository)),
				), nil
			}

			return nil, fmt.Errorf("no GitHub references found")
		},
	}
}
