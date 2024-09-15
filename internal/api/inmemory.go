package api

import (
	"context"
	"slices"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/models"
)

var _ API = (*InMemoryAPI)(nil)

type InMemoryAPI struct {
	Store *models.Store
}

func (a *InMemoryAPI) GetTags(ctx context.Context) ([]*models.Tag, error) {
	return a.Store.Tags, nil
}

func (a *InMemoryAPI) GetImages(ctx context.Context, tags []string, sort string, order string, page int64, limit int64) (*models.ImagePage, error) {
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

	filteredImages := filter(images, func(v *models.Image) bool {
		matched := 0
		for _, tag := range tags {
			if slices.Contains(v.Tags, tag) {
				matched++
			}
		}
		return matched == len(tags)
	})

	slices.SortFunc(filteredImages, func(a *models.Image, b *models.Image) int {
		valueA := ""
		valueB := ""
		if sort == "imageName" {
			valueA = a.Name
			valueB = b.Name
		} else if sort == "currentVersion" {
			// TODO: User version sort
			valueA = a.CurrentVersion
			valueB = b.CurrentVersion
		} else if sort == "latestVersion" {
			// TODO: User version sort
			valueA = a.LatestVersion
			valueB = b.LatestVersion
		}

		cmp := strings.Compare(valueA, valueB)
		if order == "desc" {
			return -cmp
		} else if order == "asc" {
			return cmp
		} else {
			// Default asc
			return cmp
		}
	})

	return &models.ImagePage{
		Images: filteredImages,
		Summary: &models.ImagePageSummary{
			Images:   len(a.Store.Images),
			Outdated: outdated,
			Pods:     pods,
		},
		Pagination: &models.PaginationMetadata{
			Total:    len(filteredImages),
			Page:     1,
			Size:     len(filteredImages),
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

func filter[T any](values []T, filterFunc func(T) bool) []T {
	new := make([]T, 0)
	for _, v := range values {
		if filterFunc(v) {
			new = append(new, v)
		}
	}
	return new
}
