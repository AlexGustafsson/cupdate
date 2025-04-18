package imageworkflow

import (
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/dockerhub"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetDockerHubVulnerabilities() workflow.Step {
	return workflow.Step{
		Name: "Get Docker Hub vulnerabilities",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			reference, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			manifest, err := workflow.GetInput[any](ctx, "manifest", true)
			if err != nil {
				return nil, err
			}

			httpClient, err := workflow.GetInput[httputil.Requester](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			// NOTE: For now, only "library" images are supported as it's unclear how
			// the API works for other images
			if strings.Contains(reference.Name(), "/") {
				return nil, nil
			}

			// TODO: For now, use the first digest of a manifest
			var digest string
			switch m := manifest.(type) {
			case *oci.ImageManifest:
				digest = m.Digest
			case *oci.ImageIndex:
				for _, m := range m.Manifests {
					if m.Digest != "" {
						digest = m.Digest
						break
					}
				}
			}

			// NOTE: For now, to not have to perform additional queries, only look up
			// manifests that include the digest upfront
			if digest == "" {
				return nil, nil
			}

			client := &dockerhub.Client{
				Client: httpClient,
			}

			vulnerabilities, err := client.GetVulnerabilities(ctx, reference.Name(), digest)
			if err != nil {
				return nil, err
			}

			return workflow.SetOutput("vulnerabilities", vulnerabilities), nil
		},
	}
}
