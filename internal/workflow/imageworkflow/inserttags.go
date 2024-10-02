package imageworkflow

import (
	"slices"

	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func InsertTag() workflow.Step {
	return workflow.Step{
		Name: "Insert tag",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			tag, err := workflow.GetInput[string](ctx, "tag", true)
			if err != nil {
				return nil, err
			}

			data, err := workflow.GetInput[*Data](ctx, "data", true)
			if err != nil {
				return nil, err
			}

			data.Lock()
			defer data.Unlock()

			if !slices.Contains(data.Tags, tag) {
				data.Tags = append(data.Tags, tag)
			}
			return nil, nil
		},
	}
}
