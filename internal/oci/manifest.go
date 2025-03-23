package oci

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ImageManifest represents an image manifest.
// This is an abstraction for all known image manifest formats (be it legacy
// Docker manifests or standard OCI manifests).
type ImageManifest struct {
	// ContentType is the MIME type as returned by the server providing the
	// manifest.
	ContentType string
	// SchemaVersion specifies the image manifest schema version.
	SchemaVersion int
	// MediaType is the MIME type of the image manifest.
	MediaType string
	// Platform optionally holds details about the platform the image supports.
	Platform *Platform
	// Digest is the digest of the index, including the "sha256:" prefix.
	Digest string
	// Annotations contains user-defined labels of the manifest.
	Annotations Annotations
}

// ImageIndex represents an group of images manifests.
// This is an abstraction for all known image index formats (be it legacy Docker
// manifests or standard OCI manifests).
type ImageIndex struct {
	// ContentType is the MIME type as returned by the server providing the
	// index.
	ContentType string
	// SchemaVersion specifies the image manifest schema version.
	SchemaVersion int
	// MediaType is the MIME type of the image index.
	MediaType string
	// Manifests contains the manifests provided by the index. Note that these may
	// or may not hold the same information as if each manifest was retrieved
	// individually.
	Manifests []ImageManifest
	// Digest is the digest of the index, including the "sha256:" prefix.
	Digest string
	// Annotations contains user-defined labels of the index.
	Annotations Annotations
}

// AttestationManifestDigests returns the digests for attestation manifests
// contained in the index, mapped by the manifest digest the attestation is for.
// SEE: https://docs.docker.com/build/metadata/attestations/attestation-storage/#attestation-manifest-descriptor.
func (i *ImageIndex) AttestationManifestDigest() map[string]string {
	digests := make(map[string]string)
	for _, manifest := range i.Manifests {
		dockerReferenceType := manifest.Annotations.DockerReferenceType()
		dockerReferenceDigest := manifest.Annotations.DockerReferenceDigest()
		if manifest.MediaType == "application/vnd.oci.image.manifest.v1+json" && dockerReferenceType == "attestation-manifest" {
			digests[dockerReferenceDigest] = manifest.Digest
		}
	}

	return digests
}

// AttestationManifestDigest returns whether or not the index contains an
// attestation manifest.
// SEE: https://docs.docker.com/build/metadata/attestations/attestation-storage/#attestation-manifest-descriptor.
func (i *ImageIndex) HasAttestationManifest() bool {
	for _, manifest := range i.Manifests {
		dockerReferenceType := manifest.Annotations.DockerReferenceType()
		if manifest.MediaType == "application/vnd.oci.image.manifest.v1+json" && dockerReferenceType == "attestation-manifest" {
			return true
		}
	}

	return false
}

// AttestationManifest represent an attestation image manifest.
type AttestationManifest struct {
	Layers []AttestationManifestLayer `json:"layers"`
}

// ProvenanceDigest returns the in-toto predicate type and digest of the first
// layer containing provenance.
func (a *AttestationManifest) ProvenanceDigest() (string, string, bool) {
	for _, layer := range a.Layers {
		if layer.MediaType != "application/vnd.in-toto+json" {
			continue
		}

		predicateType := layer.Annotations.InTotoPredicateType()
		if strings.HasPrefix(predicateType, "https://slsa.dev/provenance/") {
			return predicateType, layer.Digest, true
		}
	}

	return "", "", false
}

// SBOMDigest returns the in-toto predicate type and digest of the first layer
// containing a (well-known type of) SBOM.
func (a *AttestationManifest) SBOMDigest() (string, string, bool) {
	for _, layer := range a.Layers {
		if layer.MediaType != "application/vnd.in-toto+json" {
			continue
		}

		predicateType := layer.Annotations.InTotoPredicateType()
		switch predicateType {
		case "https://spdx.dev/Document":
			return predicateType, layer.Digest, true
		}
	}

	return "", "", false
}

// AttestationManifestLayer represents a layer entry in an attestation image
// manifest.
type AttestationManifestLayer struct {
	MediaType   string      `json:"mediaType"`
	Digest      string      `json:"digest"`
	Size        int         `json:"size"`
	Annotations Annotations `json:"annotations"`
}

type Platform struct {
	// OS is the operating system supported by the manifest.
	OS string
	// Architecture is the architecture supported by the manifest.
	Architecture string
	// Variant is the architecture variant supported by the manifest.
	// Typically a value such as "v8" for ARM images.
	Variant string
}

func manifestFromBlob(blob Blob) (any, error) {
	defer blob.Close()

	contentType := blob.Info().ContentType

	switch contentType {
	// Docker Image Manifest Version 2, Schema 1 is a deprecated manifest format
	// used by Docker engine since v1.3.0 (2024-10-04).
	// SEE: https://github.com/openshift/docker-distribution/blob/master/docs/spec/manifest-v2-1.md
	case "application/vnd.docker.distribution.manifest.v1+json",
		"application/vnd.docker.distribution.manifest.v1+prettyjws":
		var manifest struct {
			SchemaVersion int    `json:"schemaVersion"`
			Name          string `json:"name"`
			Tag           string `json:"tag"`
			Architecture  string `json:"architecture"`
		}

		if err := json.NewDecoder(blob).Decode(&manifest); err != nil {
			return nil, err
		}

		if manifest.SchemaVersion != 1 {
			return nil, fmt.Errorf("unsupported docker image manifest version 2 schema version")
		}

		var platform *Platform
		if manifest.Architecture != "" {
			platform = &Platform{
				Architecture: manifest.Architecture,
			}
		}

		digest := blob.Digest()

		return &ImageManifest{
			ContentType:   contentType,
			SchemaVersion: 1,
			Platform:      platform,
			Digest:        digest,
			Annotations:   make(Annotations),
		}, nil
	// Docker Image Manifest Version 2, Schema 2 is a manifest format used by
	// Docker.
	// SEE: https://github.com/openshift/docker-distribution/blob/master/docs/spec/manifest-v2-2.md
	case "application/vnd.docker.distribution.manifest.v2+json":
		var manifest struct {
			SchemaVersion int    `json:"schemaVersion"`
			MediaType     string `json:"mediaType"`
		}

		if err := json.NewDecoder(blob).Decode(&manifest); err != nil {
			return nil, err
		}

		if manifest.SchemaVersion != 2 {
			return nil, fmt.Errorf("unsupported docker image manifest version 2 schema version")
		}

		if manifest.MediaType != "application/vnd.docker.distribution.manifest.v2+json" {
			return nil, fmt.Errorf("unsupported docker image manifest version 2 schema 2 media type")
		}

		digest := blob.Digest()

		return &ImageManifest{
			ContentType:   contentType,
			SchemaVersion: 2,
			MediaType:     "application/vnd.docker.distribution.manifest.v2+json",
			Digest:        digest,
			Annotations:   make(Annotations),
		}, nil
	// Docker Image Manifest List Version 2, Schema 2 is a "fat" manifest format
	// used by Docker.
	// SEE: https://github.com/openshift/docker-distribution/blob/master/docs/spec/manifest-v2-2.md
	case "application/vnd.docker.distribution.manifest.list.v2+json":
		var manifest struct {
			SchemaVersion int    `json:"schemaVersion"`
			MediaType     string `json:"mediaType"`
			Manifests     []struct {
				MediaType string `json:"mediaType"`
				Digest    string `json:"digest"`
				Size      int64  `json:"size"`
				Platform  *struct {
					Architecture string `json:"architecture"`
					OS           string `json:"os"`
				} `json:"platform"`
			} `json:"manifests"`
		}

		if err := json.NewDecoder(blob).Decode(&manifest); err != nil {
			return nil, err
		}

		if manifest.SchemaVersion != 2 {
			return nil, fmt.Errorf("unsupported docker image manifest list version 2 schema version")
		}

		if manifest.MediaType != "application/vnd.docker.distribution.manifest.list.v2+json" {
			return nil, fmt.Errorf("unsupported docker image manifest list version 2 schema 2 media type")
		}

		digest := blob.Digest()

		manifests := make([]ImageManifest, len(manifest.Manifests))
		for i, manifest := range manifest.Manifests {
			var platform *Platform
			if manifest.Platform != nil {
				platform = &Platform{
					Architecture: manifest.Platform.Architecture,
					OS:           manifest.Platform.OS,
				}
			}

			manifests[i] = ImageManifest{
				MediaType:   manifest.MediaType,
				Digest:      manifest.Digest,
				Platform:    platform,
				Annotations: make(Annotations),
			}
		}

		return &ImageIndex{
			ContentType:   contentType,
			SchemaVersion: 2,
			MediaType:     "application/vnd.docker.distribution.manifest.list.v2+json",
			Manifests:     manifests,
			Digest:        digest,
			Annotations:   make(Annotations),
		}, nil
	// OCI Image Manifest v1 is a standardized image format.
	// SEE: https://github.com/opencontainers/image-spec/blob/main/manifest.md
	case "application/vnd.oci.image.manifest.v1+json":
		var manifest struct {
			SchemaVersion int               `json:"schemaVersion"`
			MediaType     string            `json:"mediaType"`
			Annotations   map[string]string `json:"annotations,omitempty"`
		}

		if err := json.NewDecoder(blob).Decode(&manifest); err != nil {
			return nil, err
		}

		if manifest.SchemaVersion != 2 {
			return nil, fmt.Errorf("unsupported oci image manifest version 2 schema version")
		}

		if manifest.MediaType != "application/vnd.oci.image.manifest.v1+json" {
			return nil, fmt.Errorf("unsupported oci image manifest version 2 schema 2 media type")
		}

		digest := blob.Digest()

		var annotations Annotations
		if manifest.Annotations == nil {
			annotations = make(Annotations)
		} else {
			annotations = manifest.Annotations
		}

		return &ImageManifest{
			ContentType:   contentType,
			SchemaVersion: 2,
			MediaType:     "application/vnd.oci.image.manifest.v1+json",
			Digest:        digest,
			Annotations:   annotations,
		}, nil
	// OCI Image Manifest List v1 is a standardized "fat" image format.
	// SEE: https://github.com/opencontainers/image-spec/blob/main/image-index.md
	case "application/vnd.oci.image.index.v1+json":
		var manifest struct {
			SchemaVersion int    `json:"schemaVersion"`
			MediaType     string `json:"mediaType"`
			Manifests     []struct {
				MediaType string `json:"mediaType"`
				Platform  *struct {
					Architecture string `json:"architecture"`
					OS           string `json:"os"`
					Variant      string `json:"variant"`
				} `json:"platform"`
				Digest      string            `json:"digest"`
				Annotations map[string]string `json:"annotations,omitempty"`
			} `json:"manifests"`
			Annotations map[string]string `json:"annotations,omitempty"`
		}

		if err := json.NewDecoder(blob).Decode(&manifest); err != nil {
			return nil, err
		}

		if manifest.SchemaVersion != 2 {
			return nil, fmt.Errorf("unsupported oci image list manifest version 2 schema version")
		}

		if manifest.MediaType != "application/vnd.oci.image.index.v1+json" {
			return nil, fmt.Errorf("unsupported oci image manifest version 2 schema 2 media type")
		}

		digest := blob.Digest()

		var annotations Annotations
		if manifest.Annotations == nil {
			annotations = make(Annotations)
		} else {
			annotations = manifest.Annotations
		}

		manifests := make([]ImageManifest, len(manifest.Manifests))
		for i, manifest := range manifest.Manifests {
			var platform *Platform
			if manifest.Platform != nil {
				platform = &Platform{
					Architecture: manifest.Platform.Architecture,
					OS:           manifest.Platform.OS,
					Variant:      manifest.Platform.Variant,
				}
			}

			var annotations Annotations
			if manifest.Annotations == nil {
				annotations = make(Annotations)
			} else {
				annotations = manifest.Annotations
			}

			manifests[i] = ImageManifest{
				MediaType:   manifest.MediaType,
				Platform:    platform,
				Digest:      manifest.Digest,
				Annotations: annotations,
			}
		}

		return &ImageIndex{
			ContentType:   contentType,
			SchemaVersion: 2,
			MediaType:     "application/vnd.oci.image.index.v1+json",
			Manifests:     manifests,
			Digest:        digest,
			Annotations:   annotations,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported manifest content type")
	}
}

// ManifestsMaybeEqual returns true if the manifests may be equal when resolved
// on the (optionally) specified platform.
func ManifestsMaybeEqual(a any, b any, platform *Platform) bool {
	digestA, manifestsA := normalizeManifest(a)
	digestB, manifestsB := normalizeManifest(b)

	// The manifests are reported by the server as being equal
	if digestA == digestB {
		return true
	}

	// Filter out manifests for the platform (or do nothing if it's not provided)
	manifestsA = filterManifestsByPlatform(manifestsA, platform)
	manifestsB = filterManifestsByPlatform(manifestsB, platform)

	// As we can't be sure what manifest will be used by the underlying engine,
	// play it safe and assume they're equal if any image is referenced in both
	for _, manifestA := range manifestsA {
		for _, manifestB := range manifestsB {
			if manifestA.Digest == manifestB.Digest {
				return true
			}
		}
	}

	return false
}

func normalizeManifest(manifest any) (string, []ImageManifest) {
	if manifest == nil {
		return "", nil
	}

	switch m := manifest.(type) {
	case *ImageIndex:
		return m.Digest, m.Manifests
	case *ImageManifest:
		return m.Digest, []ImageManifest{*m}
	}

	return "", nil
}

// filterManifestsByPlatform filters manifests based on matching platform.
func filterManifestsByPlatform(manifests []ImageManifest, platform *Platform) []ImageManifest {
	filtered := make([]ImageManifest, 0)
	for _, manifest := range manifests {
		if platform == nil {
			filtered = append(filtered, manifest)
		}

		if manifest.Platform == nil || platform == nil {
			continue
		}

		// Sometimes image authors set the fields to "unknown". Normalize such cases
		// by clearing them
		if manifest.Platform.Architecture == "unknown" {
			manifest.Platform.Architecture = ""
		}
		if manifest.Platform.OS == "unknown" {
			manifest.Platform.OS = ""
		}
		if manifest.Platform.Variant == "unknown" {
			manifest.Platform.Variant = ""
		}

		architectureMatches := platform.OS == "" || manifest.Platform.OS == platform.OS
		osMatches := platform.Architecture == "" || manifest.Platform.Architecture == platform.Architecture
		variantMatches := platform.Variant == "" || manifest.Platform.Variant == platform.Variant

		if architectureMatches && osMatches && variantMatches {
			filtered = append(filtered, manifest)
		}
	}
	return filtered
}
