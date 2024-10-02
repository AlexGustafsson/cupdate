package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func InsertLink() workflow.Step {
	return workflow.Step{
		Name: "Insert link",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			link, err := workflow.GetInput[models.ImageLink](ctx, "link", true)
			if err != nil {
				return nil, err
			}

			data, err := workflow.GetInput[*Data](ctx, "data", true)
			if err != nil {
				return nil, err
			}

			data.Lock()
			defer data.Unlock()

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
			return nil, nil
		},
	}
}
