package imageworkflow

import (
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/registry/quay"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetQuayLatestVersion() workflow.Step {
	return workflow.Step{
		Name: "Get latest version from Quay",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			httpClient, err := workflow.GetInput[*httputil.Client](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			reference, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			client := &quay.Client{Client: httpClient}

			image, err := client.GetLatestVersion(ctx, reference)
			if err != nil {
				return nil, err
			}

			fmt.Println(image)
			fmt.Println(image)

			if image == nil {
				return workflow.SetOutput("reference", (*oci.Reference)(nil)), nil
			}

			return workflow.SetOutput("reference", &image.Name), nil
		},
	}
}
