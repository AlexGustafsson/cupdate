package imageworkflow

import (
	"fmt"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/semver"
	"github.com/AlexGustafsson/cupdate/internal/vulndb"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetGitHubAdvisoriesForRepository() workflow.Step {
	return workflow.Step{
		Name: "Get GitHub advisories for repository",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			httpClient, err := workflow.GetInput[*httputil.Client](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			reference, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			owner, err := workflow.GetInput[string](ctx, "owner", true)
			if err != nil {
				return nil, err
			}

			repository, err := workflow.GetInput[string](ctx, "repository", true)
			if err != nil {
				return nil, err
			}

			vulndb, err := vulndb.AutoFetchAndOpen(ctx, "vulndb.sqlite", httpClient, 24*time.Hour)
			if err != nil {
				return nil, err
			}

			advisories, err := vulndb.GetGitHubAdvisoriesForRepository(ctx, "https://github.com/"+owner+"/"+repository)
			if err != nil {
				return nil, err
			}

			vulnerabilities := make([]models.ImageVulnerability, 0)
			for _, advisory := range advisories {
				version, err := semver.ParseVersion(reference.Version())
				if err != nil {
					continue
				}

				introducedVersion, err := semver.ParseVersion(advisory.IntroducedVersion)
				if err != nil {
					continue
				}

				// Current version is lower than introduced version
				if version.Compare(introducedVersion) < 0 {
					continue
				}

				if advisory.FixedVersion != "" {
					fixedVersion, err := semver.ParseVersion(advisory.FixedVersion)
					if err != nil {
						continue
					}

					// Current version is higher than or equal to fixed version
					if version.Compare(fixedVersion) >= 0 {
						continue
					}
				}

				vulnerabilities = append(vulnerabilities, models.ImageVulnerability{
					Severity:  string(advisory.Severity),
					Authority: "GitHub Advisory Database",
					Links:     []string{fmt.Sprintf("https://github.com/advisories/%s", advisory.ID)},
				})
			}

			return workflow.SetOutput("vulnerabilities", vulnerabilities), nil
		},
	}
}
