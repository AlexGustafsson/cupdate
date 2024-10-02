package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func InsertLatestVersion() workflow.Step {
	return workflow.Step{
		Name: "Insert latest version",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			reference, err := workflow.GetInput[*oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			data, err := workflow.GetInput[*Data](ctx, "data", true)
			if err != nil {
				return nil, err
			}

			data.Lock()
			defer data.Unlock()

			data.LatestVersion = reference
			return nil, nil
		},
	}
}
