package docker

type Version struct {
	// Version is the version of Docker.
	Version      string `json:"Version"`
	OS           string `json:"Os"`
	Architecture string `json:"Arch"`
	// MinimumAPIVersion is the minimum support API version.
	MinimumAPIVersion string `json:"MinAPIVersion"`
	// APIVersion is the current API version.
	APIVersion string `json:"ApiVersion"`
}
