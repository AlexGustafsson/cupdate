package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/registry"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetManifests() workflow.Step {
	return workflow.StepFunc("", "Get manifests", func(ctx workflow.Context) (map[string]any, error) {
		clientValue, err := ctx.Input("client")
		if err != nil {
			return nil, err
		}

		client := clientValue.(registry.Client)
		client.GetManifests(ctx)
	})
}
