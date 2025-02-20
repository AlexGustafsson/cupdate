package scorecard

import (
	"strings"
)

// RepositoryIsSupported returns true for supported repositories.
// Example:
//
//	RepositoryIsSupported("github.com/home-assistant/core")
//	RepositoryIsSupported("gitlab.com/baserow/baserow")
func RepositoryIsSupported(repository string) bool {
	host, _, ok := strings.Cut(repository, "/")
	if !ok {
		return false
	}

	if host != "github.com" && host != "gitlab.com" {
		return false
	}

	return true
}
