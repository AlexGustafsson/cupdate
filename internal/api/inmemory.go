package api

import (
	"context"
	"slices"

	"github.com/AlexGustafsson/cupdate/internal/models"
)

var _ API = (*InMemoryAPI)(nil)

type InMemoryAPI struct {
	Store *models.Store
}

func (a *InMemoryAPI) GetTags(ctx context.Context) ([]*models.Tag, error) {
	return a.Store.Tags, nil
}

func (a *InMemoryAPI) GetImages(ctx context.Context, tags []string, sort string, asc bool, desc bool, page int64, limit int64) (*models.ImagePage, error) {
	images := a.Store.Images

	outdated := 0
	pods := 0
	for _, image := range images {
		if slices.Contains(image.Tags, "pod") {
			pods++
		}
		if slices.Contains(image.Tags, "outdated") {
			outdated++
		}
	}

	return &models.ImagePage{
		Images: images,
		Summary: &models.ImagePageSummary{
			Images:   len(a.Store.Images),
			Outdated: outdated,
			Pods:     pods,
		},
		Pagination: &models.PaginationMetadata{
			Total:    len(images),
			Page:     1,
			Size:     len(images),
			Next:     "",
			Previous: "",
		},
	}, nil
}

func (a *InMemoryAPI) GetImage(ctx context.Context, name string, version string) (*models.Image, error) {
	if name == "" || version == "" {
		return nil, ErrBadRequest
	}

	for _, image := range a.Store.Images {
		if image.Name == name && image.CurrentVersion == version {
			return image, nil
		}
	}

	return nil, ErrNotFound
}

func (a *InMemoryAPI) GetImageDescription(ctx context.Context, name string, version string) (*models.ImageDescription, error) {
	if name == "" || version == "" {
		return nil, ErrBadRequest
	}

	result, ok := a.Store.Descriptions[name+":"+version]
	if !ok {
		return nil, ErrNotFound
	}

	return result, nil
}

func (a *InMemoryAPI) GetImageReleaseNotes(ctx context.Context, name string, version string) (*models.ImageReleaseNotes, error) {
	if name == "" || version == "" {
		return nil, ErrBadRequest
	}

	result, ok := a.Store.ReleaseNotes[name+":"+version]
	if !ok {
		return nil, ErrNotFound
	}

	return result, nil
}

func (a *InMemoryAPI) GetImageGraph(ctx context.Context, name string, version string) (*models.Graph, error) {
	if name == "" || version == "" {
		return nil, ErrBadRequest
	}

	result, ok := a.Store.Graphs[name+":"+version]
	if !ok {
		return nil, ErrNotFound
	}

	return result, nil
}
