package api

import "context"

var _ API = (*InMemoryAPI)(nil)

type InMemoryAPI struct {
	Tags   []Tag
	Images []Image
	// Descriptions is mapped by name:version
	Descriptions map[string]*ImageDescription
	// ReleaseNotes is mapped by name:version
	ReleaseNotes map[string]*ImageReleaseNotes
	// Graphs is mapped by name:version
	Graphs map[string]*Graph
}

func (a *InMemoryAPI) GetTags(ctx context.Context) ([]Tag, error) {
	return a.Tags, nil
}

func (a *InMemoryAPI) GetImages(ctx context.Context, tags []string, sort string, asc bool, desc bool, page int64, limit int64) (*ImagePage, error) {
	images := a.Images

	return &ImagePage{
		Images: images,
		Summary: ImagePageSummary{
			Images: len(a.Images),
		},
		Pagination: PaginationMetadata{
			Total:    len(images),
			Page:     1,
			Size:     len(images),
			Next:     "",
			Previous: "",
		},
	}, nil
}

func (a *InMemoryAPI) GetImage(ctx context.Context, name string, version string) (*Image, error) {
	if name == "" || version == "" {
		return nil, ErrBadRequest
	}

	for _, image := range a.Images {
		if image.Name == name && image.CurrentVersion == version {
			return &image, nil
		}
	}

	return nil, ErrNotFound
}

func (a *InMemoryAPI) GetImageDescription(ctx context.Context, name string, version string) (*ImageDescription, error) {
	if name == "" || version == "" {
		return nil, ErrBadRequest
	}

	result, ok := a.Descriptions[name+":"+version]
	if !ok {
		return nil, ErrNotFound
	}

	return result, nil
}

func (a *InMemoryAPI) GetImageReleaseNotes(ctx context.Context, name string, version string) (*ImageReleaseNotes, error) {
	if name == "" || version == "" {
		return nil, ErrBadRequest
	}

	result, ok := a.ReleaseNotes[name+":"+version]
	if !ok {
		return nil, ErrNotFound
	}

	return result, nil
}

func (a *InMemoryAPI) GetImageGraph(ctx context.Context, name string, version string) (*Graph, error) {
	if name == "" || version == "" {
		return nil, ErrBadRequest
	}

	result, ok := a.Graphs[name+":"+version]
	if !ok {
		return nil, ErrNotFound
	}

	return result, nil
}
