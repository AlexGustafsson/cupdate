package ghcr

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

func PackagePath(image oci.Reference) string {
	user, name, _ := strings.Cut(image.Path, "/")
	return fmt.Sprintf("https://github.com/users/%s/packages/container/package/%s", url.PathEscape(user), url.PathEscape(name))
}
