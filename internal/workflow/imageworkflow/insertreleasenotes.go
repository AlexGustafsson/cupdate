package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func InsertReleaseNotes(data *Data, f func(ctx workflow.Context) *models.ImageReleaseNotes) workflow.Step {
	return workflow.StepFunc("", "Insert description", func(ctx workflow.Context) (map[string]any, error) {
		data.Lock()
		defer data.Unlock()

		data.ReleaseNotes = f(ctx)

		return nil, nil
	})
}
