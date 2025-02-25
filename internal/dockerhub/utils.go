package dockerhub

import (
	"net/url"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/oci"
)

// RepositoryUIPath returns a URL to an image's repository page on Docker Hub.
func RepositoryUIPath(image oci.Reference) string {
	owner, name, _ := strings.Cut(image.Path, "/")
	if owner == "library" {
		return "https://hub.docker.com/_/" + url.PathEscape(name)
	}

	return "https://hub.docker.com/r/" + url.PathEscape(owner) + "/" + url.PathEscape(name)
}

// TagUIPath returns a URL to an image's tag page on Docker Hub.
func TagUIPath(image oci.Reference, digest string) string {
	owner, name, _ := strings.Cut(image.Path, "/")
	return "https://hub.docker.com/layers/" + url.PathEscape(owner) + "/" + url.PathEscape(name) + "/" + url.PathEscape(image.Tag) + "/images/" + url.PathEscape(strings.ReplaceAll(digest, ":", "-"))
}
