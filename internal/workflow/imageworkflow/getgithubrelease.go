package imageworkflow

import (
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetGitHubRelease(manifestsKey string) workflow.Step {
	return workflow.Step{
		Name: "Get GitHub release",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			manifests, err := workflow.GetInput[[]oci.Manifest](ctx, "manifests", true)
			if err != nil {
				return nil, err
			}

			if manifests == nil {
				return nil, fmt.Errorf("no manifests found")
			}

			// TODO: Find github links in manifests?

			// client := &github.Client{}

			// client.GetRelease(ctx, )

			return nil, fmt.Errorf("not implemented")
		},
	}
}
