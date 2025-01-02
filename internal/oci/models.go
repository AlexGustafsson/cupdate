package oci

import (
	"strings"
	"time"
)

type DockerDistributionManifestListV2 struct {
	// 2
	SchemaVersion int `json:"schemaVersion"`
	// application/vnd.docker.distribution.manifest.list.v2+json
	MediaType string                         `json:"mediaType"`
	Manifests []DockerDistributionManifestV2 `json:"manifests"`
}

type DockerDistributionManifestV2 struct {
	// 2
	SchemaVersion int `json:"schemaVersion"`
	// application/vnd.docker.distribution.manifest.v2+json
	MediaType string `json:"mediaType"`
	Digest    string `json:"digest"`
	Platform  struct {
		Architecture string `json:"architecture"`
		OS           string `json:"os"`
		Variant      string `json:"variant"`
	} `json:"platform"`
	Size int `json:"size"`
}

type OCIImageIndexV1 struct {
	// 2
	SchemaVersion int `json:"schemaVersion"`
	// application/vnd.oci.image.index.v1+json
	MediaType string               `json:"mediaType"`
	Manifests []OCIImageManifestV1 `json:"manifests"`
}

type OCIImageManifestV1 struct {
	// 2
	SchemaVersion int `json:"schemaVersion"`
	// application/vnd.oci.image.manifest.v1+json
	MediaType   string            `json:"mediaType"`
	Annotations map[string]string `json:"annotations"`
	Digest      string            `json:"digest"`
	Platform    struct {
		Architecture string `json:"architecture"`
		OS           string `json:"os"`
		Variant      string `json:"variant"`
	} `json:"platform"`
}

// application/vnd.docker.distribution.manifest.v1+prettyjws
type DockerDistributionManifestV1 struct {
	// 1
	SchemaVersion int    `json:"schemaVersion"`
	Name          string `json:"name"`
	Tag           string `json:"tag"`
	Architecture  string `json:"architecture"`
}

type Manifest struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	// Annotations contains manifest / image annotations. Nil if none were found.
	// Note the even if annotations were found at the top level, they might not
	// match the annotations / label of the image itself.
	Annotations Annotations `json:"annotations"`
	Digest      string      `json:"digest"`
	Platform    *Platform   `json:"platform,omitempty"`
}

type Platform struct {
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	Variant      string `json:"variant"`
}

type Annotations map[string]string

func (a Annotations) Source() string {
	if a == nil {
		return ""
	}

	s := a["org.opencontainers.image.source"]
	if s == "" {
		s = a["org.label-schema.vcs-url"]
	}

	return s
}

func (a Annotations) CreatedTime() time.Time {
	s := a["org.opencontainers.image.created"]
	if s == "" {
		s = a["org.label-schema.build-date"]
	}

	// Golang's RFC3339 format requires a "T", but the spec allows a space
	s = strings.Replace(s, " ", "T", 1)
	time, _ := time.Parse(time.RFC3339, s)
	return time
}

func (a Annotations) URL() string {
	return a["org.opencontainers.image.url"]
}

func (a Annotations) DocumentationURL() string {
	return a["org.opencontainers.image.documentation"]
}

type TagsPage struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
