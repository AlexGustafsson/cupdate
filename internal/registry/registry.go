package registry

import (
	"context"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

type Registry interface {
	GetLatestVersion(ctx context.Context, image oci.Reference) (*Image, error)
}

type Image struct {
	Name      oci.Reference
	Published time.Time
	Digest    string
}

type Client interface {
	GetManifests(ctx context.Context, image oci.Reference) ([]oci.Manifest, error)
}
