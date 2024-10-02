package imageworkflow

import (
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/registry/docker"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func SetupRegistryClient(data *Data) workflow.Step {
	return workflow.StepFunc("", "Setup registry client", func(ctx workflow.Context) (map[string]any, error) {
		// TODO: Support other registries
		if data.ImageReference.Domain != "docker.io" {
			return nil, fmt.Errorf("unsupported registry domain: %s", data.ImageReference.Domain)
		}

		return map[string]any{
			"client": &docker.Client{},
			"domain": data.ImageReference.Domain,
		}, nil
	})
}
