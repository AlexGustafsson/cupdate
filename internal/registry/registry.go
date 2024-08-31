package registry

import (
	"context"
	"time"
)

type Registry interface {
	GetLatestVersion(ctx context.Context, name string) (*Image, error)
	Get(ctx context.Context, name string, version string) (*Image, error)
}

type Image struct {
	Name         string
	Version      string
	Published    time.Time
	Digest       string
	ReleaseNotes string
}
