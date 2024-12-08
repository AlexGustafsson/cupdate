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
			httpClient, err := workflow.GetInput[*httputil.Client](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			reference, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			endpoint, err := workflow.GetInput[string](ctx, "endpoint", true)
			if err != nil {
				return nil, err
			}

			owner, err := workflow.GetInput[string](ctx, "owner", true)
			if err != nil {
				return nil, err
			}

			repository, err := workflow.GetInput[string](ctx, "repository", true)
			if err != nil {
				return nil, err
			}

			if !reference.HasTag {
				return nil, fmt.Errorf("cannot get GitHub release for image without tag")
			}

			client := &github.Client{
				Endpoint: endpoint,
				Client:   httpClient,
			}

			release, err := client.GetRelease(ctx, owner, repository, reference.Tag)
			if err != nil {
				return nil, err
			}

			// It's not uncommon for tags / releases to be prefixed with "v". If no
			// release was found for the verbatim release, also try with a "v" prefix
			if release == nil && !strings.HasPrefix(reference.Tag, "v") {
				release, err = client.GetRelease(ctx, owner, repository, "v"+reference.Tag)
				if err != nil {
					return nil, err
				}
			}

			return workflow.SetOutput("release", release), nil
		},
	}
}
