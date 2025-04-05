package imageworkflow

import (
	"os"
	"path/filepath"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/semver"
	"github.com/AlexGustafsson/cupdate/internal/vulndb"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetGitHubAdvisoriesForRepository() workflow.Step {
	return workflow.Step{
		Name: "Get GitHub advisories for repository",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			httpClient, err := workflow.GetInput[httputil.Requester](ctx, "httpClient", true)
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

			version, err := semver.ParseVersion(reference.Version())
			if err != nil {
				// We won't be able to compare versions correctly
				return nil, nil
			}

			vulndb, err := vulndb.AutoFetchAndOpen(ctx, filepath.Join(os.TempDir(), "vulndb.sqlite"), httpClient, 24*time.Hour)
			if err != nil {
				return nil, err
			}
			defer vulndb.Close()

			vulnerabilities, err := vulndb.GetGitHubAdvisoriesForRepository(ctx, "https://github.com/"+owner+"/"+repository, version)
			if err != nil {
				return nil, err
			}

			return workflow.SetOutput("vulnerabilities", vulnerabilities), nil
		},
	}
}
