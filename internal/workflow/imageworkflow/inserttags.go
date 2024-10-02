package imageworkflow

import (
	"slices"

	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func InsertTags(data *Data, f func(ctx workflow.Context) []string) workflow.Step {
	return workflow.StepFunc("", "Insert tags", func(ctx workflow.Context) (map[string]any, error) {
		data.Lock()
		defer data.Unlock()

		for _, tag := range f(ctx) {
			if !slices.Contains(data.Tags, tag) {
				data.Tags = append(data.Tags, tag)
			}
		}

		return nil, nil
	})
}
