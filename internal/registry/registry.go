package registry

import (
	"context"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/distribution/reference"
)

type Registry interface {
	GetLatestVersion(ctx context.Context, name string, currentTag string) (*Image, error)
	Get(ctx context.Context, name string, version string) (*Image, error)
}

type Image struct {
	Name         string
	Version      string
	Published    time.Time
	Digest       string
	ReleaseNotes string
}

type Client interface {
	GetManifests(ctx context.Context, image reference.Named) ([]oci.Manifest, error)
}
