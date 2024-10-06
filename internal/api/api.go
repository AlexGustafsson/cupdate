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
	GetTags(ctx context.Context) ([]models.Tag, error)
	ListImages(ctx context.Context, tags []string, order string, page int, limit int) (*models.ImagePage, error)
	GetImage(ctx context.Context, reference string) (*models.Image, error)
	GetImageDescription(ctx context.Context, reference string) (*models.ImageDescription, error)
	GetImageReleaseNotes(ctx context.Context, reference string) (*models.ImageReleaseNotes, error)
	GetImageGraph(ctx context.Context, reference string) (*models.Graph, error)
}
