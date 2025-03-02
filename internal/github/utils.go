package github

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/oci"
)

// PackageURL returns the URL to a GitHub package, assuming the reference is a
// valid image hosted on GHCR.
func PackageURL(reference oci.Reference) string {
	user, name, _ := strings.Cut(reference.Path, "/")
	return fmt.Sprintf("https://github.com/users/%s/packages/container/package/%s", url.PathEscape(user), url.PathEscape(name))
}
