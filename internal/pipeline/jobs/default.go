package jobs

import "github.com/AlexGustafsson/cupdate/internal/pipeline"

func DefaultJobs() pipeline.Job[ImageData] {
	setupRegistryClient := SetupRegistryClient()
	getManifests := GetImageManifests(setupRegistryClient.Output.Client)

	return pipeline.Series[ImageData]{
		setupRegistryClient,
		getManifests,
		// pipeline.Parallel{
		// 	EnrichLinks(getManifests.Outputs.Manifests),
		// 	GetDescription(setupRegistryClient.Outputs.Client),
		// 	GetLatestImageVersion(),
		// 	GetGitHubRelease(getManifests.Outputs.Manifests),
		// },
	}
}
