package imageworkflow

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/openssf/scorecard"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetOpenSSFScorecard() workflow.Step {
	return workflow.Step{
		Name: "Get OpenSSF scorecard for repository",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			httpClient, err := workflow.GetInput[*httputil.Client](ctx, "httpClient", true)
			if err != nil {
				return nil, err
			}

			repository, err := workflow.GetInput[string](ctx, "repository", true)
			if err != nil {
				return nil, err
			}

			fmt.Println(repository)

			if !scorecard.RepositoryIsSupported(repository) {
				return nil, nil
			}

			client := scorecard.Client{
				Client: httpClient,
			}

			scorecard, err := client.GetScorecard(ctx, repository)
			if err != nil {
				return nil, err
			}

			if scorecard == nil {
				return nil, nil
			}

			reportTime, err := scorecard.Time()
			if err != nil {
				slog.WarnContext(ctx, "Failed to parse date scorecard report was gathered", slog.Any("error", err))
				return nil, nil
			}

			// Only use scores that were made within the last two months
			if time.Since(reportTime) > 2*30*24*time.Hour {
				return nil, nil
			}

			var risk models.ImageScorecardRisk
			if scorecard.Score <= 2.5 {
				risk = models.ImageScorecardRiskCritical
			} else if scorecard.Score <= 5 {
				risk = models.ImageScorecardRiskHigh
			} else if scorecard.Score <= 7.5 {
				risk = models.ImageScorecardRiskMedium
			} else {
				risk = models.ImageScorecardRiskLow
			}

			result := &models.ImageScorecard{
				ReportURL:  "https://scorecard.dev/viewer/?uri=" + repository,
				Score:      scorecard.Score,
				Risk:       risk,
				GenerateAt: reportTime,
			}

			return workflow.SetOutput("scorecard", result), nil
		},
	}
}
