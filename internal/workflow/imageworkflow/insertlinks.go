package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func InsertLinks(data *Data, f func(ctx workflow.Context) []models.ImageLink) workflow.Step {
	return workflow.StepFunc("", "Insert links", func(ctx workflow.Context) (map[string]any, error) {
		data.Lock()
		defer data.Unlock()

		for _, link := range f(ctx) {
			exists := false
			for _, other := range data.Links {
				if link.Type == other.Type && link.URL == other.URL {
					exists = true
					break
				}
			}

			if !exists {
				data.Links = append(data.Links, link)
			}
		}

		return nil, nil
	})
}
