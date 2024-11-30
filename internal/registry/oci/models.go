package oci

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
	SchemaVersion int               `json:"schemaVersion"`
	MediaType     string            `json:"mediaType"`
	Annotations   map[string]string `json:"annotations"`
	Digest        string            `json:"digest"`
}

func (m Manifest) SourceAnnotation() string {
	s := m.Annotations["org.opencontainers.image.source"]
	if s == "" {
		s = m.Annotations["org.label-schema.vcs-url"]
	}
	return s
}

func (m Manifest) RevisionAnnotation() string {
	return m.Annotations["org.opencontainers.image.revision"]
}
