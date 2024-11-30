package imageworkflow

import (
	"fmt"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/registry/docker"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
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

			manifests, err := workflow.GetInput[[]oci.Manifest](ctx, "manifests", true)
			if err != nil {
				return nil, err
			}

			httpClient, err := workflow.GetInput[*httputil.Client](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			// NOTE: For now, only "library" images are supported as it's unclear how
			// the API works for other images
			if strings.Contains(reference.Name(), "/") {
				return nil, nil
			}

			// TODO: For now, use the first digest of a manifest
			digest := ""
			for _, manifest := range manifests {
				fmt.Printf("%+v\n", manifest)
				if manifest.Digest != "" {
					digest = manifest.Digest
					break
				}
			}

			// NOTE: For now, to not have to perform additional queries, only look up
			// manifests that include the digest upfront
			if digest == "" {
				return nil, nil
			}

			client := &docker.Client{
				Client: httpClient,
			}

			report, err := client.GetVulnerabilityReport(ctx, reference.Name(), digest)
			if err != nil {
				return nil, err
			}

			vulnerabilities := make([]models.ImageVulnerability, 0)

			if report != nil {
				for i := 0; i < report.Criticial; i++ {
					vulnerabilities = append(vulnerabilities, models.ImageVulnerability{
						Severity:  "criticial",
						Authority: "Docker Scout",
						Link:      docker.TagUIPath(reference, digest),
					})
				}

				for i := 0; i < report.High; i++ {
					vulnerabilities = append(vulnerabilities, models.ImageVulnerability{
						Severity:  "high",
						Authority: "Docker Scout",
						Link:      docker.TagUIPath(reference, digest),
					})
				}

				for i := 0; i < report.Medium; i++ {
					vulnerabilities = append(vulnerabilities, models.ImageVulnerability{
						Severity:  "medium",
						Authority: "Docker Scout",
						Link:      docker.TagUIPath(reference, digest),
					})
				}

				for i := 0; i < report.Low; i++ {
					vulnerabilities = append(vulnerabilities, models.ImageVulnerability{
						Severity:  "low",
						Authority: "Docker Scout",
						Link:      docker.TagUIPath(reference, digest),
					})
				}

				for i := 0; i < report.Unspecified; i++ {
					vulnerabilities = append(vulnerabilities, models.ImageVulnerability{
						Severity:  "unspecified",
						Authority: "Docker Scout",
						Link:      docker.TagUIPath(reference, digest),
					})
				}
			}

			return workflow.SetOutput("vulnerabilities", vulnerabilities), nil
		},
	}
}
