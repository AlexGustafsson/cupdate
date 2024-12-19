package registry

import (
	"time"

	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

// TODO: Remove the Published and Digest fields as the oci job to get a manifest
// can just return the digest and created time?
type Image struct {
	Name      oci.Reference
	Published time.Time
	Digest    string
}
