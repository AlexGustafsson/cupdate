package api

import (
	"context"
	"errors"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)

type API interface {
	GetTags(ctx context.Context) ([]Tag, error)
	GetImages(ctx context.Context, tags []string, sort string, asc bool, desc bool, page int64, limit int64) (*ImagePage, error)
	GetImage(ctx context.Context, name string, version string) (*Image, error)
	GetImageDescription(ctx context.Context, name string, version string) (*ImageDescription, error)
	GetImageReleaseNotes(ctx context.Context, name string, version string) (*ImageReleaseNotes, error)
	GetImageGraph(ctx context.Context, name string, version string) (*Graph, error)
}
