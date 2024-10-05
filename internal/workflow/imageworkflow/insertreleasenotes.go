package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func InsertReleaseNotes() workflow.Step {
	return workflow.Step{
		Name: "Insert release notes",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			releaseNotes, err := workflow.GetInput[*models.ImageReleaseNotes](ctx, "releaseNotes", true)
			if err != nil {
				return nil, err
			}

			data, err := workflow.GetInput[*Data](ctx, "data", true)
			if err != nil {
				return nil, err
			}

			data.Lock()
			defer data.Unlock()

			data.ReleaseNotes = releaseNotes
			return nil, nil
		},
	}
}
