package ghcr

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/oci"
)

func PackagePath(reference oci.Reference) string {
	user, name, _ := strings.Cut(reference.Path, "/")
	return fmt.Sprintf("https://github.com/users/%s/packages/container/package/%s", url.PathEscape(user), url.PathEscape(name))
}
