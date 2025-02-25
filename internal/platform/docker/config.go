package docker

// ConfigFile implements config file format used by Docker (and Kubernetes) for
// registry secrets and auth.
type ConfigFile struct {
	Auths       map[string]ConfigEntry `json:"auths"`
	HttpHeaders map[string]string      `json:"HttpHeaders,omitempty"`
}

// ConfigEntry is an entry for a specific pattern, as defined in a [ConfigFile].
type ConfigEntry struct {
	Auth     string `json:"auth,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}
