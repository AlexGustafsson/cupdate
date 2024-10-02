package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func InsertLatestVersion(data *Data, f func(ctx workflow.Context) *oci.Reference) workflow.Step {
	return workflow.StepFunc("", "Insert latest version", func(ctx workflow.Context) (map[string]any, error) {
		data.Lock()
		defer data.Unlock()

		data.LatestVersion = f(ctx)

		return nil, nil
	})
}
