package registry

import (
	"time"

	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

type Image struct {
	Name      oci.Reference
	Published time.Time
	Digest    string
}
