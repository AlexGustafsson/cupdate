package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func InsertDescription() workflow.Step {
	return workflow.Step{
		Name: "Insert description",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			description, err := workflow.GetInput[*models.ImageDescription](ctx, "description", true)
			if err != nil {
				return nil, err
			}

			data, err := workflow.GetInput[*Data](ctx, "data", true)
			if err != nil {
				return nil, err
			}

			data.Lock()
			defer data.Unlock()

			data.Description = description
			return nil, nil
		},
	}
}
