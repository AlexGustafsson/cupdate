package registry

import (
	"context"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/distribution/reference"
)

type Registry interface {
	GetLatestVersion(ctx context.Context, image reference.NamedTagged) (*Image, error)
}

type Image struct {
	Name      reference.NamedTagged
	Published time.Time
	Digest    string
}

type Client interface {
	GetManifests(ctx context.Context, image reference.Named) ([]oci.Manifest, error)
}
