package docker

type ConfigFile struct {
	Auths       map[string]ConfigEntry `json:"auths"`
	HttpHeaders map[string]string      `json:"HttpHeaders,omitempty"`
}

type ConfigEntry struct {
	Auth     string `json:"auth,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}
