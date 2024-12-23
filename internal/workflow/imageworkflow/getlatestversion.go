package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/semver"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetLatestReference() workflow.Step {
	return workflow.Step{
		Name: "Get latest reference",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			reference, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			registryClient, err := workflow.GetInput[*oci.Client](ctx, "registryClient", true)
			if err != nil {
				return nil, err
			}

			tags, err := registryClient.GetTags(ctx, reference, &oci.GetTagsOptions{
				AllPages: true,
			})
			if err != nil {
				return nil, err
			}

			var latestReference *oci.Reference
			if tags != nil && reference.Tag != "" {
				// We only want to specify a latest reference when we're certain of it,
				// for example, when it has been seen in the list of tags
				latest, ok := semver.LatestOpinionatedVersionString(reference.Tag, tags)
				if ok {
					l := reference
					l.Tag = latest
					latestReference = &l
				}
			}

			return workflow.SetOutput("reference", latestReference), nil
		},
	}
}
