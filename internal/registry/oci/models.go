package oci

type Manifest struct {
	Annotations map[string]string `json:"annotations"`
	Digest      string            `json:"digest"`
	MediaType   string            `json:"mediaType"`
	Platform    struct {
		Architecture string `json:"architecture"`
		OS           string `json:"os"`
	} `json:"platform"`
	Size int `json:"size"`
}
