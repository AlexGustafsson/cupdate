package oci

import (
	"encoding/json"
	"fmt"
)

type ImageManifest struct {
	ContentType   string
	SchemaVersion int
	MediaType     string
	Platform      *Platform
	Digest        string
	Annotations   Annotations
}

type ImageIndex struct {
	ContentType   string
	SchemaVersion int
	MediaType     string
	Manifests     []ImageManifest
	Digest        string
	Annotations   Annotations
}

type Platform struct {
	OS           string
	Architecture string
	Variant      string
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
