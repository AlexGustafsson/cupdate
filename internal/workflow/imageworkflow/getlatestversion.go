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

			var latestReference *oci.Reference

			_, err = semver.ParseVersion(reference.Tag)
			isSemver := err == nil

			// If the tag is adhering to semver, try to identify the latest available
			// tag
			if isSemver {
				tags, err := registryClient.GetTags(ctx, reference, &oci.GetTagsOptions{
					AllPages: true,
				})
				if err != nil {
					return nil, err
				}

				// We only want to specify a latest reference when we're certain of it,
				// for example, when it has been seen in the list of tags
				latest, ok := semver.LatestOpinionatedVersionString(reference.Tag, tags)
				if ok {
					l := reference

					// Set the latest tag
					l.HasTag = true
					l.Tag = latest

					// Remove the digest of the latest reference in all cases as we don't
					// know the new image's digest (yet)
					l.HasDigest = false
					l.Digest = ""

					// Try to find the digest
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

			// As the latest reference could very well match the same semantic version
			// as the current, always compare the current manifest with the latest to
			// see if the underlying manifests have changed
			if reference.HasTag && reference.HasDigest {
				currentManifest, err := registryClient.GetManifest(ctx, reference)
				if err != nil {
					return nil, err
				}

				var latestManifest any
				if latestReference == nil {
					// Remove the digest and only look up the latest information for the
					// current tag instead
					ref := reference
					ref.HasDigest = false
					ref.Digest = ""

					latestManifest, err = registryClient.GetManifest(ctx, ref)
					if err != nil {
						return nil, err
					}

					ref.HasDigest = true
					switch m := latestManifest.(type) {
					case *oci.ImageIndex:
						ref.Digest = m.Digest
					case *oci.ImageManifest:
						ref.Digest = m.Digest
					}
					latestReference = &ref
				} else {
					latestManifest, err = registryClient.GetManifest(ctx, *latestReference)
					if err != nil {
						return nil, err
					}
				}

				// Doing an as good job as we can, try to see if the manifest in use is
				// equal to the new one, if it is, just assume the current reference is
				// the latest
				// TODO: Specify an actual platform
				maybeEqual := oci.ManifestsMaybeEqual(currentManifest, latestManifest, nil)
				if maybeEqual {
					latestReference = &reference
				}
			}

			return workflow.SetOutput("reference", latestReference), nil
		},
	}
}
