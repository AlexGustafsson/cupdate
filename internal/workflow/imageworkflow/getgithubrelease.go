package imageworkflow

import (
	"fmt"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/github"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetGitHubRelease() workflow.Step {
	return workflow.Step{
		Name: "Get GitHub release",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			manifests, err := workflow.GetInput[[]oci.Manifest](ctx, "manifests", true)
			if err != nil {
				return nil, err
			}

			httpClient, err := workflow.GetInput[*httputil.Client](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			reference, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			if manifests == nil {
				return nil, fmt.Errorf("no manifests found")
			}

			if !reference.HasTag {
				return nil, fmt.Errorf("cannot get GitHub release for image without tag")
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
			if !ok {
				return nil, fmt.Errorf("no GitHub reference found")
			}

			client := &github.Client{
				Endpoint: endpoint,
				Client:   httpClient,
			}

			release, err := client.GetRelease(ctx, owner, repository, reference.Tag)
			if err != nil {
				return nil, err
			}

			return workflow.SetOutput("release", release), nil
		},
	}
}
