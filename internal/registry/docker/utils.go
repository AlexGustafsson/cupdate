package docker

import (
	"net/url"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

func RepositoryUIPath(image oci.Reference) string {
	owner, name, _ := strings.Cut(image.Path, "/")
	if owner == "library" {
		return "https://hub.docker.com/_/" + url.PathEscape(name)
	}

	return "https://hub.docker.com/r/" + url.PathEscape(owner) + "/" + url.PathEscape(name)
}
