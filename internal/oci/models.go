package oci

type DockerDistributionManifestListV2 struct {
	// 2
	SchemaVersion int `json:"schemaVersion"`
	// application/vnd.docker.distribution.manifest.list.v2+json
	MediaType string `json:"mediaType"`
	Manifests []struct {
		// application/vnd.docker.distribution.manifest.v2+json
		// application/vnd.docker.distribution.manifest.v1+json
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
		Platform  struct {
			Architecture string `json:"architecture"`
			OS           string `json:"os"`
			Variant      string `json:"variant"`
		} `json:"platform"`
	} `json:"manifests"`
}

type DockerDistributionManifestV2 struct {
	// 2
	SchemaVersion int `json:"schemaVersion"`
	// application/vnd.docker.distribution.manifest.v2+json
	MediaType string `json:"mediaType"`
	Config    struct {
		// application/vnd.docker.container.image.v1+json
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"config"`
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

type TagsPage struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
