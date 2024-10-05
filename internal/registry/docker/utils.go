package docker

import (
	"net/url"

	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

func RepositoryPath(image oci.Reference) string {
	return "https://hub.docker.com/v2/repositories/" + url.PathEscape(image.Path)
}
