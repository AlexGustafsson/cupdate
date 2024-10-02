package jobs

import (
	"log/slog"
	"strings"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/github"
	"github.com/AlexGustafsson/cupdate/internal/pipeline"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

type GetGitHubReleaseJob struct {
	Output GetGitHubReleaseJobOutput

	manifests string
}

type GetGitHubReleaseJobOutput struct {
	Release string
}

func GetGitHubRelease(manifests string) *GetGitHubReleaseJob {
	return &GetGitHubReleaseJob{
		Output: GetGitHubReleaseJobOutput{
			Release: "get-github-release/release",
		},
		manifests: manifests,
	}
}

func (j GetGitHubReleaseJob) Execute(ctx pipeline.Context[ImageData]) error {
	ctx.Lock()
	defer ctx.Unlock()

	log := slog.With(slog.String("imageReference", ctx.Data.ImageReference.String()))

	var manifests []oci.Manifest
	if !ctx.GetOutput(j.manifests, &manifests) {
		log.Debug("Skipping GitHub release - no manifests found")
		return nil
	}

	if ctx.Data.LatestVersion == nil {
		log.Debug("Skipping GitHub release - latest version not found")
		return nil
	}

	if !ctx.Data.LatestVersion.HasTag {
		log.Debug("Skipping GitHub release - latest version isn't tagged")
		return nil
	}

	cacheKey := "pipeline/get-github-release-v1/" + ctx.Data.ImageReference.String()
	var release *github.Release
	if err := ctx.Cache().GetJSON(ctx, cacheKey, &release, 24*time.Hour); err != nil {
		log.Error("Failed to get cache", slog.Any("error", err))
		// Fallthrough
	}

	if release == nil {
		log.Debug("Fetching felease")

		// Look through all manifests and find a valid reference to GitHub
		var endpoint, owner, repository string
		for _, manifest := range manifests {
			source := manifest.SourceAnnotation()
			if source == "" {
				continue
			}

			// TODO: Support enterprise or other hosts?
			if !strings.Contains(source, "github.com") {
				continue
			}

			var ok bool
			endpoint, owner, repository, _, ok = github.ParseURL(source)
			if !ok {
				continue
			}
			break
		}

		if endpoint == "" || owner == "" || repository == "" {
			log.Debug("No GitHub URL found")
			return nil
		}

		var err error
		client := &github.Client{}
		release, err = client.GetRelease(ctx, owner, repository, ctx.Data.LatestVersion.Tag)
		if err != nil {
			log.Error("Failed to get release", slog.Any("error", err))
			return err
		}

		if err := ctx.Cache().SetJSON(ctx, cacheKey, &release); err != nil {
			log.Error("Failed to set cache", slog.Any("error", err))
			// Fallthrough
		}
	}

	if release == nil {
		log.Info("GitHub release not found")
		return nil
	}

	ctx.SetOutput(j.Output.Release, release)

	return nil
}
