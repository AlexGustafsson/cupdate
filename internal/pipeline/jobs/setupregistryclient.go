package jobs

import (
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/pipeline"
	"github.com/AlexGustafsson/cupdate/internal/registry/docker"
)

type SetupRegistryClientJob struct {
	Output SetupRegistryClientJobOutput
}

type SetupRegistryClientJobOutput struct {
	Client string
}

func SetupRegistryClient() SetupRegistryClientJob {
	return SetupRegistryClientJob{
		Output: SetupRegistryClientJobOutput{
			Client: "setup-registry-client/client",
		},
	}
}

func (j SetupRegistryClientJob) Execute(ctx pipeline.Context[ImageData]) error {
	ctx.Lock()
	defer ctx.Unlock()

	// TODO: Support other registries
	if ctx.Data.ImageReference.Domain != "docker.io" {
		return fmt.Errorf("unsupported registry domain: %s", ctx.Data.ImageReference.Domain)
	}

	ctx.SetOutput(j.Output.Client, &docker.Client{})

	return nil
}
