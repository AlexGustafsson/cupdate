package api

import (
	"context"
	"errors"

	"github.com/AlexGustafsson/cupdate/internal/models"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)

type API interface {
	GetTags(ctx context.Context) ([]*models.Tag, error)
	GetImages(ctx context.Context, tags []string, sort string, asc bool, desc bool, page int64, limit int64) (*models.ImagePage, error)
	GetImage(ctx context.Context, name string, version string) (*models.Image, error)
	GetImageDescription(ctx context.Context, name string, version string) (*models.ImageDescription, error)
	GetImageReleaseNotes(ctx context.Context, name string, version string) (*models.ImageReleaseNotes, error)
	GetImageGraph(ctx context.Context, name string, version string) (*models.Graph, error)
}
