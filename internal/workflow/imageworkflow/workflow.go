package imageworkflow

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/dockerhub"
	"github.com/AlexGustafsson/cupdate/internal/github"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func New(httpClient httputil.Requester, data *Data) workflow.Workflow {
	return workflow.Workflow{
		Name: "Process image",
		Jobs: []workflow.Job{
			{
				ID:   "oci",
				Name: "Get OCI information",
				Steps: []workflow.Step{
					SetupRegistryClient().
						WithID("registry").
						With("httpClient", httpClient).
						With("reference", data.ImageReference).
						With("registryAuth", data.RegistryAuth),
					GetManifest().
						WithID("manifest").
						With("registryClient", workflow.Ref{Key: "step.registry.client"}).
						With("reference", data.ImageReference),
					GetAnnotations().
						WithID("annotations").
						With("registryClient", workflow.Ref{Key: "step.registry.client"}).
						With("reference", data.ImageReference).
						With("manifest", workflow.Ref{Key: "step.manifest.manifest"}),
					GetLatestReference().
						WithID("latest").
						With("registryClient", workflow.Ref{Key: "step.registry.client"}).
						With("reference", data.ImageReference).
						With("graph", data.Graph),
					GetManifest().
						WithID("latest-manifest").
						WithCondition(workflow.ValueExists("step.latest.reference")).
						With("registryClient", workflow.Ref{Key: "step.registry.client"}).
						With("reference", workflow.Ref{Key: "step.latest.reference"}),
					GetAnnotations().
						WithID("latest-annotations").
						WithCondition(workflow.ValueExists("step.latest.reference")).
						With("registryClient", workflow.Ref{Key: "step.registry.client"}).
						With("reference", workflow.Ref{Key: "step.latest.reference"}).
						With("manifest", workflow.Ref{Key: "step.latest-manifest.manifest"}),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						domain, err := workflow.GetValue[string](ctx, "step.registry.domain")
						if err != nil {
							return nil, err
						}

						data.InsertLink(models.ImageLink{
							Type: "oci-registry",
							URL:  "https://" + domain,
						})

						annotations, err := workflow.GetValue[oci.Annotations](ctx, "step.annotations.annotations")
						if err != nil {
							return nil, err
						}

						if url := annotations.Source(); url != "" {
							data.InsertLink(models.ImageLink{
								Type: "svc",
								URL:  url,
							})
						}

						if url := annotations.URL(); url != "" {
							data.InsertLink(models.ImageLink{
								Type: "generic",
								URL:  url,
							})
						}

						if url := annotations.DocumentationURL(); url != "" {
							data.InsertLink(models.ImageLink{
								Type: "docs",
								URL:  url,
							})
						}

						if time := annotations.CreatedTime(); !time.IsZero() {
							data.Created = &time
						}

						reference, err := workflow.GetValue[*oci.Reference](ctx, "step.latest.reference")
						if err != nil {
							return nil, err
						}

						data.LatestReference = reference

						latestAnnotations, err := workflow.GetValue[oci.Annotations](ctx, "step.latest-annotations.annotations")
						if err != nil {
							return nil, err
						}

						if latestAnnotations != nil {
							time := latestAnnotations.CreatedTime()
							if !time.IsZero() {
								data.LatestCreated = &time
							}
						}

						return nil, nil
					}),
				},
			},
			{
				ID:        "attestations",
				Name:      "Get attestations",
				DependsOn: []string{"oci"},
				Steps: []workflow.Step{
					GetAttestationManifests().
						WithID("attestations").
						With("registryClient", workflow.Ref{Key: "job.oci.step.registry.client"}).
						With("reference", data.ImageReference).
						With("manifest", workflow.Ref{Key: "job.oci.step.manifest.manifest"}),
					GetProvenanceAttestations().
						WithID("provenance").
						With("registryClient", workflow.Ref{Key: "job.oci.step.registry.client"}).
						With("reference", data.ImageReference).
						With("manifests", workflow.Ref{Key: "step.attestations.manifests"}),
					GetSBOMAttestations().
						WithID("sbom").
						With("registryClient", workflow.Ref{Key: "job.oci.step.registry.client"}).
						With("reference", data.ImageReference).
						With("manifests", workflow.Ref{Key: "step.attestations.manifests"}),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						currentManifest, err := workflow.GetValue[any](ctx, "job.oci.step.manifest.manifest")
						if err != nil {
							return nil, err
						}

						currentIndexManifest, ok := currentManifest.(*oci.ImageIndex)
						if ok {
							if currentIndexManifest.HasAttestationManifest() {
								data.InsertTag("attestation")
							}
						}

						index, _ := currentManifest.(*oci.ImageIndex)

						provenanceAttestations, err := workflow.GetValue[map[string]oci.ProvenanceAttestation](ctx, "step.provenance.attestations")
						if err != nil {
							return nil, err
						}

						provenance := &models.ImageProvenance{
							BuildInfo: []models.ProvenanceBuildInfo{},
						}
						for imageDigest, attestation := range provenanceAttestations {
							buildInfo := models.ProvenanceBuildInfo{
								ImageDigest:     imageDigest,
								Source:          attestation.Source,
								SourceRevision:  attestation.SourceRevision,
								BuildStartedOn:  attestation.BuildStartedOn,
								BuildFinishedOn: attestation.BuildFinishedOn,
								Dockerfile:      attestation.Dockerfile,
							}

							var imageManifest *oci.ImageManifest
							if index != nil {
								for _, manifest := range index.Manifests {
									if manifest.Digest == imageDigest {
										imageManifest = &manifest
										break
									}
								}
							}

							if imageManifest != nil && imageManifest.Platform != nil {
								buildInfo.OperatingSystem = imageManifest.Platform.OS
								buildInfo.Architecture = imageManifest.Platform.Architecture
								buildInfo.ArchitectureVariant = imageManifest.Platform.Variant
							}

							provenance.BuildInfo = append(provenance.BuildInfo, buildInfo)
						}

						data.Provenance.OK = true
						if len(provenance.BuildInfo) > 0 {
							data.Provenance.Value = provenance
						}

						sbomAttestations, err := workflow.GetValue[map[string]oci.SBOMAttestation](ctx, "step.sbom.attestations")
						if err != nil {
							return nil, err
						}

						if len(sbomAttestations) > 0 {
							data.InsertTag("sbom")
						}

						sboms := &models.ImageSBOM{
							SBOM: []models.SBOM{},
						}
						for imageDigest, attestation := range sbomAttestations {
							sbom := models.SBOM{
								ImageDigest: imageDigest,
								Type:        string(attestation.Type),
								SBOM:        attestation.SBOM,
							}

							var imageManifest *oci.ImageManifest
							if index != nil {
								for _, manifest := range index.Manifests {
									if manifest.Digest == imageDigest {
										imageManifest = &manifest
										break
									}
								}
							}

							if imageManifest != nil && imageManifest.Platform != nil {
								sbom.OperatingSystem = imageManifest.Platform.OS
								sbom.Architecture = imageManifest.Platform.Architecture
								sbom.ArchitectureVariant = imageManifest.Platform.Variant
							}

							sboms.SBOM = append(sboms.SBOM, sbom)
						}

						data.SBOM.OK = true
						if len(sboms.SBOM) > 0 {
							data.SBOM.Value = sboms
						}

						return nil, nil
					}),
				},
			},
			{
				ID:        "docker",
				Name:      "Get Docker Hub information",
				DependsOn: []string{"oci"},
				// Only run for Docker images
				If: func(ctx workflow.Context) (bool, error) {
					domain, err := workflow.GetValue[string](ctx, "job.oci.step.registry.domain")
					if err != nil {
						return false, err
					}

					return domain == "docker.io", nil
				},
				Steps: []workflow.Step{
					GetDockerHubRepository().
						WithID("repository").
						With("httpClient", httpClient).
						With("reference", data.ImageReference),
					GetDockerHubRepositoryOwner().
						WithID("owner").
						With("httpClient", httpClient).
						With("repository", workflow.Ref{Key: "step.repository.repository"}),
					GetDockerHubVulnerabilities().
						WithID("vulnerabilities").
						With("httpClient", httpClient).
						With("reference", data.ImageReference).
						With("manifest", workflow.Ref{Key: "job.oci.step.manifest.manifest"}),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						repository, err := workflow.GetValue[*dockerhub.Repository](ctx, "step.repository.repository")
						if err != nil {
							return nil, err
						}

						data.Description = repository.Description
						data.FullDescription.OK = true
						data.FullDescription.Value = &models.ImageDescription{
							Markdown: repository.FullDescription,
						}

						vulnerabilities, err := workflow.GetValue[[]models.ImageVulnerability](ctx, "step.vulnerabilities.vulnerabilities")
						if err != nil {
							return nil, err
						}

						if len(vulnerabilities) > 0 {
							data.InsertVulnerabilities(vulnerabilities)
						}

						return nil, nil
					}),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						owner, err := workflow.GetValue[*dockerhub.Entity](ctx, "step.owner.owner")
						if err != nil {
							return nil, err
						}

						data.Image = owner.GravatarURL
						return nil, nil
					}),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						data.InsertLink(models.ImageLink{
							Type: "docker",
							URL:  dockerhub.RepositoryUIPath(data.ImageReference),
						})
						return nil, nil
					}),
				},
			},
			{
				ID:        "ghcr",
				Name:      "Get GHCR information",
				DependsOn: []string{"oci"},
				// Only run for GHCR images
				If: func(ctx workflow.Context) (bool, error) {
					domain, err := workflow.GetValue[string](ctx, "job.oci.step.registry.domain")
					if err != nil {
						return false, err
					}

					return domain == "ghcr.io", nil
				},
				Steps: []workflow.Step{
					GetGithubPackage().
						WithID("package").
						With("httpClient", httpClient).
						With("reference", data.ImageReference),
					GetGitHubDescription().
						WithID("description").
						With("httpClient", httpClient).
						With("owner", workflow.Ref{Key: "step.package.owner"}).
						With("repository", workflow.Ref{Key: "step.package.repository"}),
					GetGithubPackageREADME().
						WithID("readme").
						With("httpClient", httpClient).
						With("package", workflow.Ref{Key: "step.package.package"}),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						description, err := workflow.GetValue[string](ctx, "step.description.description")
						if err != nil {
							return nil, err
						}

						data.InsertLink(models.ImageLink{
							Type: "ghcr",
							URL:  github.PackageURL(data.ImageReference),
						})
						data.Description = description

						readme, err := workflow.GetValue[string](ctx, "step.readme.readme")
						if err != nil {
							return nil, err
						}

						data.FullDescription.OK = true
						if readme != "" {
							data.FullDescription.Value = &models.ImageDescription{
								HTML: readme,
							}
						}

						return nil, nil
					}),
				},
			},
			{
				ID:        "quay",
				Name:      "Get Quay information",
				DependsOn: []string{},
				// Only run for Quay images
				If: func(ctx workflow.Context) (bool, error) {
					return data.ImageReference.Domain == "quay.io", nil
				},
				Steps: []workflow.Step{
					GetQuayVulnerabilities().
						WithID("vulnerabilities").
						With("httpClient", httpClient).
						With("reference", data.ImageReference),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						vulnerabilities, err := workflow.GetValue[[]models.ImageVulnerability](ctx, "step.vulnerabilities.vulnerabilities")
						if err != nil {
							return nil, err
						}

						if len(vulnerabilities) > 0 {
							data.InsertVulnerabilities(vulnerabilities)
						}

						return nil, nil
					}),
				},
			},
			{
				ID:        "gitlab",
				Name:      "Get GitLab information",
				DependsOn: []string{"oci"},
				// Only run for GitLab images
				If: func(ctx workflow.Context) (bool, error) {
					domain, err := workflow.GetValue[string](ctx, "job.oci.step.registry.domain")
					if err != nil {
						return false, err
					}

					return domain == "registry.gitlab.com", nil
				},
				Steps: []workflow.Step{
					GetGitLabDescription().
						WithID("description").
						With("reference", data.ImageReference).
						With("httpClient", httpClient),
					GetGitLabRepositoryREADME().
						WithID("readme").
						With("reference", data.ImageReference).
						With("httpClient", httpClient),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						data.InsertLink(models.ImageLink{
							Type: "gitlab",
							URL:  "https://gitlab.com/" + data.ImageReference.Path,
						})

						description, err := workflow.GetValue[string](ctx, "step.description.description")
						if err != nil {
							return nil, err
						}
						data.Description = description

						readmeMime, err := workflow.GetValue[string](ctx, "step.readme.mime")
						if err != nil {
							return nil, err
						}

						html, err := workflow.GetValue[string](ctx, "step.readme.html")
						if err != nil {
							return nil, err
						}

						raw, err := workflow.GetValue[[]byte](ctx, "step.readme.raw")
						if err != nil {
							return nil, err
						}

						data.FullDescription.OK = true
						// Prefer markdown over pre-rendered HTML
						if readmeMime == "text/markdown" && raw != nil {
							data.FullDescription.Value = &models.ImageDescription{
								Markdown: string(raw),
							}
						} else if html != "" {
							data.FullDescription.Value = &models.ImageDescription{
								HTML: html,
							}
						}

						return nil, nil
					}),
				},
			},
			{
				ID:   "github",
				Name: "Get GitHub information",
				// Depend on whatever provides us with the latest image version
				// Implicitly depends on OCI
				DependsOn: []string{"docker", "ghcr", "gitlab"},
				// Only run for images with a reference to GitHub
				If: func(ctx workflow.Context) (bool, error) {
					if data.ImageReference.Domain == "ghcr.io" {
						return true, nil
					}

					annotations, err := workflow.GetValue[oci.Annotations](ctx, "job.oci.step.annotations.annotations")
					if err != nil {
						return false, err
					} else if annotations == nil {
						return false, nil
					}

					source := annotations.Source()
					if strings.HasPrefix(source, "https://github.com/") {
						return true, nil
					}

					return false, nil
				},
				Steps: []workflow.Step{
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						data.InsertTag("github")

						reference := data.ImageReference
						if data.LatestReference != nil {
							reference = *data.LatestReference
						}
						return workflow.SetOutput("reference", reference), nil
					}).WithID("init"),
					GetGitHubRepsitory().
						WithID("repository").
						With("annotations", workflow.Ref{Key: "job.oci.step.annotations.annotations"}).
						With("package", workflow.Ref{Key: "job.ghcr.step.package.package"}),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						endpoint, err := workflow.GetValue[string](ctx, "step.repository.endpoint")
						if err != nil {
							return nil, err
						}

						owner, err := workflow.GetValue[string](ctx, "step.repository.owner")
						if err != nil {
							return nil, err
						}

						repository, err := workflow.GetValue[string](ctx, "step.repository.name")
						if err != nil {
							return nil, err
						}

						data.InsertLink(models.ImageLink{
							Type: "github",
							URL:  fmt.Sprintf("%s/%s/%s", endpoint, url.PathEscape(owner), url.PathEscape(repository)),
						})

						return nil, nil
					}),
					// TODO: Get description if not found
					GetGitHubRelease().
						WithID("release").
						With("httpClient", httpClient).
						With("endpoint", workflow.Ref{Key: "step.repository.endpoint"}).
						With("owner", workflow.Ref{Key: "step.repository.owner"}).
						With("repository", workflow.Ref{Key: "step.repository.name"}).
						With("reference", workflow.Ref{Key: "step.init.reference"}),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						release, err := workflow.GetValue[*github.Release](ctx, "step.release.release")
						if err != nil {
							return nil, err
						}

						data.ReleaseNotes.OK = true
						if release == nil {
							return nil, nil
						}

						data.ReleaseNotes.Value = &models.ImageReleaseNotes{
							Title: release.Title,
							HTML:  release.Description,
						}
						return nil, nil
					}),
					GetGitHubAdvisoriesForRepository().
						WithID("vulnerabilities").
						With("httpClient", httpClient).
						With("reference", data.ImageReference).
						With("owner", workflow.Ref{Key: "step.repository.owner"}).
						With("repository", workflow.Ref{Key: "step.repository.name"}),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						vulnerabilities, err := workflow.GetValue[[]models.ImageVulnerability](ctx, "step.vulnerabilities.vulnerabilities")
						if err != nil {
							return nil, err
						}

						if len(vulnerabilities) > 0 {
							data.InsertVulnerabilities(vulnerabilities)
						}

						return nil, nil
					}),
				},
			},
			{
				ID:        "openssf",
				Name:      "Get OpenSSF Scorecard",
				DependsOn: []string{"github", "gitlab"},
				If: func(ctx workflow.Context) (bool, error) {
					githubRepository, err := workflow.GetValue[string](ctx, "job.github.step.repository.repository")
					if err != nil {
						return false, err
					}

					gitlabRepository := ""
					if data.ImageReference.Domain == "registry.gitlab.com" {
						// The repository path is <owner>/<group>/<project>
						parts := strings.Split(data.ImageReference.Path, "/")
						if len(parts) < 3 {
							return false, nil
						}

						gitlabRepository = "gitlab.com/" + strings.Join(parts[0:3], "/")
					}

					return githubRepository != "" || gitlabRepository != "", nil
				},
				Steps: []workflow.Step{
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						githubRepository, err := workflow.GetValue[string](ctx, "job.github.step.repository.repository")
						if err != nil {
							return nil, err
						}

						if githubRepository != "" {
							return workflow.SetOutput("repository", githubRepository), nil
						}

						if data.ImageReference.Domain == "registry.gitlab.com" {
							// The repository path is <owner>/<group>/<project>
							parts := strings.Split(data.ImageReference.Path, "/")
							if len(parts) < 3 {
								return nil, nil
							}

							gitlabRepository := "gitlab.com/" + strings.Join(parts[0:3], "/")
							return workflow.SetOutput("repository", gitlabRepository), nil
						}

						return nil, nil
					}).WithID("repository"),
					GetOpenSSFScorecard().
						WithID("scorecard").
						With("httpClient", httpClient).
						With("repository", workflow.Ref{Key: "step.repository.repository"}),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						scorecard, err := workflow.GetValue[*models.ImageScorecard](ctx, "step.scorecard.scorecard")
						if err != nil {
							return nil, err
						}

						data.Scorecard.OK = true
						if scorecard != nil {
							data.Scorecard.Value = scorecard
							data.InsertLink(models.ImageLink{
								Type: "openssf-scorecard",
								URL:  scorecard.ReportURL,
							})
						}

						return nil, nil
					}),
				},
			},
			{
				ID:        "sbom",
				Name:      "Scan SBOMs",
				DependsOn: []string{"attestations"},
				If: func(ctx workflow.Context) (bool, error) {
					sbom, err := workflow.GetValue[map[string]oci.SBOMAttestation](ctx, "job.attestations.step.sbom.attestations")
					if err != nil {
						return false, err
					}

					return sbom != nil, nil
				},
				Steps: []workflow.Step{
					ScanSBOM().
						WithID("vulnerabilities").
						With("attestations", workflow.Ref{Key: "job.attestations.step.sbom.attestations"}),
					workflow.Run(func(ctx workflow.Context) (workflow.Command, error) {
						vulnerabilities, err := workflow.GetValue[[]models.ImageVulnerability](ctx, "step.vulnerabilities.vulnerabilities")
						if err != nil {
							return nil, err
						}

						if len(vulnerabilities) > 0 {
							data.InsertVulnerabilities(vulnerabilities)
						}

						return nil, nil
					}),
				},
			},
		},
	}
}
