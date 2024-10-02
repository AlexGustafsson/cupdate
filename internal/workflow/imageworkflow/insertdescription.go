package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func InsertDescription(data *Data, f func(ctx workflow.Context) *models.ImageDescription) workflow.Step {
	return workflow.StepFunc("", "Insert description", func(ctx workflow.Context) (map[string]any, error) {
		data.Lock()
		defer data.Unlock()

		data.Description = f(ctx)

		return nil, nil
	})
}
