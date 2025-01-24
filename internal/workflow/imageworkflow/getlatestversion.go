package imageworkflow

import (
	"log/slog"

	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/semver"
	"github.com/AlexGustafsson/cupdate/internal/workflow"
)

func GetLatestReference() workflow.Step {
	return workflow.Step{
		Name: "Get latest reference",
		Main: func(ctx workflow.Context) (workflow.Command, error) {
			reference, err := workflow.GetInput[oci.Reference](ctx, "reference", true)
			if err != nil {
				return nil, err
			}

			registryClient, err := workflow.GetInput[*oci.Client](ctx, "registryClient", true)
			if err != nil {
				return nil, err
			}

			// If the tag is not adhering to semver, try to identify whether or not
			// the underlying manifest has changed if we have a digest to compare to
			if reference.HasTag && reference.HasDigest {
				_, err := semver.ParseVersion(reference.Tag)
				isSemver := err == nil
				if !isSemver {
					// Remove the digest and only look up the tag
					ref := reference
					ref.HasDigest = false
					ref.Digest = ""
					manifest, err := registryClient.GetManifest(ctx, ref)
					if err != nil {
						return nil, err
					}

					switch m := manifest.(type) {
					// If we got a single image, we can be pretty sure it's what's
					// actually used by the runtime
					// TODO: We don't know if this is just one manifest from the same fat
					// manifest already in use by the user (just look up the reference's
					// manifest and see if its either the same, or in the case of index,
					// contains this manifest?) That way we always try to use the fat
					// manifest if possible, probably what's used by most users / systems
					// anyway?
					case *oci.ImageManifest:
						ref := reference
						ref.HasDigest = true
						ref.Digest = m.Digest
						slog.Debug("Identified fixed tag and found an image manifest", slog.String("reference", reference.String()), slog.String("latestReference", ref.String()))
						return workflow.SetOutput("reference", &ref), nil
					case *oci.ImageIndex:
						slog.Debug("Identified fixed tag and found an image index", slog.String("reference", reference.String()))
						// We don't know if this is just the fat manifest for the same
						// manifest the user is already using, look it up
						currentDigestFound := false
						for _, manifest := range m.Manifests {
							if manifest.Digest == reference.Digest {
								currentDigestFound = true
								break
							}
						}
						if currentDigestFound {
							slog.Debug("Fixed tag in use was found in index", slog.String("reference", reference.String()))
							return workflow.SetOutput("reference", &reference), nil
						}

						// If we got a "fat" manifest, we don't necessarily know what digest
						// is actually used by the runtime as it can vary by platform /
						// architecture, but the version itself can still be referred to by
						// the fat manifest, use it
						ref := reference
						ref.HasDigest = true
						ref.Digest = m.Digest
						slog.Debug("Fixed tag in use was not found in index, assuming index is an updated version", slog.String("reference", reference.String()), slog.String("latestReference", ref.String()))
						return workflow.SetOutput("reference", &ref), nil
					}
				}
			}

			tags, err := registryClient.GetTags(ctx, reference, &oci.GetTagsOptions{
				AllPages: true,
			})
			if err != nil {
				return nil, err
			}

			var latestReference *oci.Reference
			if tags != nil && reference.Tag != "" {
				// We only want to specify a latest reference when we're certain of it,
				// for example, when it has been seen in the list of tags
				latest, ok := semver.LatestOpinionatedVersionString(reference.Tag, tags)
				if ok {
					l := reference
					l.HasTag = true
					l.Tag = latest
					// Remove the digest of the latest reference in all cases as we don't
					// know the new image's digest
					l.HasDigest = false
					l.Digest = ""

					// Try to find the actual digest
					manifest, err := registryClient.GetManifest(ctx, l)
					if err == nil {
						switch m := manifest.(type) {
						case *oci.ImageManifest:
							l.HasDigest = true
							l.Digest = m.Digest
						case *oci.ImageIndex:
							l.HasDigest = true
							l.Digest = m.Digest
						}
					} else {
						slog.Warn("Failed to look up the manifest of the latest reference, falling back to no digest", slog.Any("error", err))
					}

					latestReference = &l
				}
			}

			return workflow.SetOutput("reference", latestReference), nil
		},
	}
}
