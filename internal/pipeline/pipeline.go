package pipeline

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/registry/docker"
)

// TODO: I don't quite like the name. It's a pipeline taking source data,
// filtering them, deduplicating and then enriching them?
type Pipeline struct {
	cache cache.Cache
}

func New(cache cache.Cache) *Pipeline {
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

	newStore := &models.Store{
		Tags:         tags,
		Images:       images,
		Descriptions: make(map[string]*models.ImageDescription),
		ReleaseNotes: make(map[string]*models.ImageReleaseNotes),
		Graphs:       graphs,
	}

	if err := p.EnrichFromManifests(ctx, newStore); err != nil {
		slog.Error("Failed to enrich images from manifests", slog.Any("error", err))
	}

	if err := p.EnrichFromDockerHub(ctx, newStore); err != nil {
		slog.Error("Failed to enrich images from manifests", slog.Any("error", err))
	}

	return newStore, nil
}

func (p *Pipeline) EnrichFromManifests(ctx context.Context, store *models.Store) error {
	// Temp before cache
	fetched := 0
	for _, image := range store.Images {
		if fetched > 0 {
			continue
		}
		registry := ""
		name := ""

		parts := strings.Split(image.Name, "/")
		if strings.Contains(parts[0], ".") {
			registry = parts[0]
			name = strings.Join(parts[1:], "/")
		} else {
			registry = "docker.io"
			name = strings.Join(parts[0:], "/")
		}

		// Sanity check
		if registry == "" || name == "" {
			panic("invalid state - registry or name is empty")
		}

		log := slog.With(slog.String("registry", registry), slog.String("image", name), slog.String("version", image.CurrentVersion))

		// TODO: Support other registries
		if registry != "docker.io" {
			log.Warn("Skipping unsupported registry")
			return nil
		}

		log.Debug("Fetching annotations")

		c := docker.Client{}
		// TODO: Cache
		manifests, err := c.GetManifests(ctx, name, image.CurrentVersion)
		if err != nil {
			slog.Error("Failed to get manifests", slog.Any("error", err))
			return err
		}

		if len(manifests) == 0 {
			slog.Error("Got zero manifests", slog.Any("error", err))
			return nil
		}

		source := manifests[0].SourceAnnotation()
		if source != "" {
			// TODO: Identify different sources (GitHub etc.)
			image.Links = append(image.Links, &models.ImageLink{
				Type: "git",
				URL:  source,
			})
			fetched++
		}
	}

	return nil
}

func (p *Pipeline) EnrichFromDockerHub(ctx context.Context, store *models.Store) error {
	// Temp before cache
	fetched := 0
	for _, image := range store.Images {
		if fetched > 0 {
			continue
		}
		registry := ""
		name := ""

		parts := strings.Split(image.Name, "/")
		if strings.Contains(parts[0], ".") {
			registry = parts[0]
			name = strings.Join(parts[1:], "/")
		} else {
			registry = "docker.io"
			name = strings.Join(parts[0:], "/")
		}

		// Sanity check
		if registry == "" || name == "" {
			panic("invalid state - registry or name is empty")
		}

		log := slog.With(slog.String("registry", registry), slog.String("image", name), slog.String("version", image.CurrentVersion))

		if registry != "docker.io" {
			log.Warn("Skipping unsupported registry")
			return nil
		}

		log.Debug("Fetching repository")

		c := docker.Client{}
		// TODO: Cache
		repository, err := c.GetRepository(ctx, name)
		if err != nil {
			slog.Error("Failed to get Docker Hub  repository", slog.Any("error", err))
			return err
		}

		owner := repository.Namespace
		if owner == "library" {
			owner = "_"
		}
		image.Links = append(image.Links, &models.ImageLink{
			Type: "docker",
			URL:  fmt.Sprintf("https://hub.docker.com/%s/%s", owner, repository.Name),
		})

		image.Description = repository.Description
		store.Descriptions[image.Name+":"+image.CurrentVersion] = &models.ImageDescription{
			Markdown: repository.FullDescription,
		}

		// TODO: add tags from Docker Hub categories?
		fetched++
	}

	return nil
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
