package imageworkflow

import (
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func SetupRegistryClient() workflow.Step {
	return workflow.Step{
		Name: "Setup registry client",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			httpClient, err := workflow.GetInput[httputil.Requester](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			image, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			registryAuth, err := workflow.GetInput[*httputil.AuthMux](ctx, "registryAuth", true)
			if err != nil {
				return nil, err
			}

			client := &oci.Client{
				Client:   httpClient,
				AuthFunc: registryAuth.HandleAuth,
			}

			return workflow.Batch(
				workflow.SetOutput("client", client),
				workflow.SetOutput("domain", image.Domain),
			), nil
		},
	}
}
