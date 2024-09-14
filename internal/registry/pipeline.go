package registry

import (
	"context"
	"reflect"

	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/models"
)

// TODO: I don't quite like the name. It's a pipeline taking source data,
// filtering them, deduplicating and then enriching them?
type Pipeline struct {
	cache cache.Cache
}

func NewPipeline(cache cache.Cache) *Pipeline {
	return &Pipeline{
		cache: cache,
	}
}

func (p *Pipeline) Run(ctx context.Context, store *models.UnprocessedStore) (*models.Store, error) {
	// Deduplicate tags
	tags := deduplicate(store.Tags, key[*models.Tag]("Name"), first)

	// Deduplicate images
	images := deduplicate(store.Images, func(v *models.Image) string {
		return v.Name + ":" + v.CurrentVersion
	}, func(t []*models.Image) *models.Image {
		i := &models.Image{
			Name:           t[0].Name,
			CurrentVersion: t[0].CurrentVersion,
			LatestVersion:  t[0].LatestVersion,
			Tags:           []string{},
			Links:          []*models.ImageLink{},
			Image:          "",
		}

		for _, t := range t {
			i.Tags = append(i.Tags, t.Tags...)
		}
		i.Tags = deduplicate(i.Tags, identity, first)

		return i
	})

	// Deduplicate graphs
	graphs := make(map[string]*models.Graph)
	for k, entries := range store.Graphs {
		g := &models.Graph{
			Root: &models.GraphNode{
				Domain:  entries[0].Root.Domain,
				Type:    entries[0].Root.Type,
				Name:    entries[0].Root.Name,
				Parents: make([]*models.GraphNode, 0),
			},
		}

		for _, e := range entries {
			g.Root.Parents = append(g.Root.Parents, e.Root.Parents...)
		}

		graphs[k] = g
	}

	// TODO: Populate other data

	return &models.Store{
		Tags:   tags,
		Images: images,
		Graphs: graphs,
	}, nil
}

func first[T any](values []T) T {
	return values[0]
}

func key[T any](k string) func(T) string {
	return func(t T) string {
		v := reflect.ValueOf(t)
		return reflect.Indirect(v).FieldByName(k).String()
	}
}

func identity(k string) string {
	return k
}

func deduplicate[T any](values []T, keyFunc func(T) string, mergeFunc func([]T) T) []T {
	valuesByKey := make(map[string][]T)
	for _, v := range values {
		key := keyFunc(v)
		entries, ok := valuesByKey[key]
		if !ok {
			entries = make([]T, 0)
		}
		entries = append(entries, v)
		valuesByKey[key] = entries
	}

	result := make([]T, 0)
	for _, entries := range valuesByKey {
		result = append(result, mergeFunc(entries))
	}

	return result
}
