package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func insertLinks(ctx workflow.Context, links []models.ImageLink) (workflow.Command, error) {
	data, err := workflow.GetInput[*Data](ctx, "data", true)
	if err != nil {
		return nil, err
	}

	data.Lock()
	defer data.Unlock()

	for _, link := range links {
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

}

func InsertLinks() workflow.Step {
	return workflow.Step{
		Name: "Insert links",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			links, err := workflow.GetInput[[]models.ImageLink](ctx, "links", true)
			if err != nil {
				return nil, err
			}

			return insertLinks(ctx, links)
		},
	}
}

func InsertLink() workflow.Step {
	return workflow.Step{
		Name: "Insert link",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			link, err := workflow.GetInput[models.ImageLink](ctx, "link", true)
			if err != nil {
				return nil, err
			}

			return insertLinks(ctx, []models.ImageLink{link})
		},
	}
}
