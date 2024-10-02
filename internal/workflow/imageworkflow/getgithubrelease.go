package imageworkflow

import (
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetGitHubRelease(manifestsKey string) workflow.Step {
	return workflow.StepFunc("", "Get GitHub release", func(ctx workflow.Context) (map[string]any, error) {
		manifests, ok := ctx.Output(manifestsKey)
		if !ok {
			return nil, fmt.Errorf("no manifests")
		}
	})
}
