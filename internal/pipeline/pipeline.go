package pipeline

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/AlexGustafsson/cupdate/internal/github"
	"github.com/AlexGustafsson/cupdate/internal/models"
	"github.com/AlexGustafsson/cupdate/internal/registry"
	"github.com/AlexGustafsson/cupdate/internal/registry/docker"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
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

	if err := p.EnrichFromGitHub(ctx, newStore); err != nil {
		slog.Error("Failed to enrich images from GitHub", slog.Any("error", err))
	}

	return newStore, nil
}

func (p *Pipeline) EnrichFromManifests(ctx context.Context, store *models.Store) error {
	for _, image := range store.Images {
		if err := p.enrichImageFromManifests(ctx, image); err != nil {
			slog.Error("Failed to enrich from manifest", slog.Any("error", err))
			// Fallthrough
		}
	}

	return nil
}

func (p *Pipeline) enrichImageFromManifests(ctx context.Context, image *models.Image) error {
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

	cacheKey := fmt.Sprintf("v1/docker/manifests/%s/%s", name, image.CurrentVersion)
	var manifests []oci.Manifest
	if err := p.cache.GetJSON(ctx, cacheKey, &manifests, 24*time.Hour); err != nil {
		slog.Error("Failed to get cache", slog.Any("error", err))
		// Fallthrough
	}

	if manifests == nil {
		log.Debug("Fetching annotations")

		c := docker.Client{}
		// TODO: Cache
		var err error
		manifests, err = c.GetManifests(ctx, name, image.CurrentVersion)
		if err != nil {
			slog.Error("Failed to get manifests", slog.Any("error", err))
			return err
		}

		if err := p.cache.SetJSON(ctx, cacheKey, &manifests); err != nil {
			slog.Error("Failed to set cache", slog.Any("error", err))
			// Fallthrough
		}
	}

	if manifests == nil {
		slog.Error("Image manifests not found")
		return nil
	}

	if len(manifests) == 0 {
		slog.Error("Got zero manifests")
		return nil
	}

	// TODO: Support custom overrides. Very few images seem to actually use
	// these annotations...
	source := manifests[0].SourceAnnotation()
	if source != "" {
		// TODO: Identify different sources (GitHub etc.)
		image.Links = append(image.Links, &models.ImageLink{
			Type: "git",
			URL:  source,
		})
	}

	return nil
}

func (p *Pipeline) EnrichFromDockerHub(ctx context.Context, store *models.Store) error {
	for _, image := range store.Images {
		if err := p.enrichDescriptionFromDockerHub(ctx, image, store); err != nil {
			slog.Error("Failed to enrich description from Docker Hub", slog.Any("error", err))
			// Fallthrough
		}

		if err := p.enrichImageReleaseFromRegistry(ctx, image); err != nil {
			slog.Error("Failed to enrich image release from registry", slog.Any("error", err))
			// Fallthrough
		}

		// TODO: add tags from Docker Hub categories?
	}

	return nil
}

func (p *Pipeline) enrichDescriptionFromDockerHub(ctx context.Context, image *models.Image, store *models.Store) error {
	registryName := ""
	name := ""

	parts := strings.Split(image.Name, "/")
	if strings.Contains(parts[0], ".") {
		registryName = parts[0]
		name = strings.Join(parts[1:], "/")
	} else {
		registryName = "docker.io"
		name = strings.Join(parts[0:], "/")
	}

	// Sanity check
	if registryName == "" || name == "" {
		panic("invalid state - registry or name is empty")
	}

	log := slog.With(slog.String("registry", registryName), slog.String("image", name), slog.String("version", image.CurrentVersion))

	if registryName != "docker.io" {
		log.Warn("Skipping unsupported registry", slog.String("registry", registryName), slog.String("image", name), slog.String("version", image.CurrentVersion))
		return nil
	}

	cacheKey := fmt.Sprintf("v1/docker/repositories/%s/%s", name, image.CurrentVersion)
	var repository *docker.Repository
	if err := p.cache.GetJSON(ctx, cacheKey, &repository, 24*time.Hour); err != nil {
		slog.Error("Failed to get cache", slog.Any("error", err))
		// Fallthrough
	}

	if repository == nil {
		log.Debug("Fetching repository")

		c := docker.Client{}
		var err error
		// TODO: A repository can hold multiple images (like yooooomi/your_spotify_your-spotify-frontend)
		// support this by actually reading the manifest to see what the image's
		// namespace is? Then I guess we would need to support looking at different
		// names in the GetLatestVersion function as not all images are for the
		// current version, despite sharing repository?
		repository, err = c.GetRepository(ctx, name)
		if err != nil {
			slog.Error("Failed to get Docker Hub repository", slog.Any("error", err))
			return nil
		}

		if repository != nil {
			if err := p.cache.SetJSON(ctx, cacheKey, &repository); err != nil {
				slog.Error("Failed to set cache", slog.Any("error", err))
				// Fallthrough
			}
		}
	}

	if repository == nil {
		log.Warn("Repository not found", slog.String("image", name))
		return nil
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

	return nil
}

func (p *Pipeline) enrichImageReleaseFromRegistry(ctx context.Context, image *models.Image) error {
	registryName := ""
	name := ""

	parts := strings.Split(image.Name, "/")
	if strings.Contains(parts[0], ".") {
		registryName = parts[0]
		name = strings.Join(parts[1:], "/")
	} else {
		registryName = "docker.io"
		name = strings.Join(parts[0:], "/")
	}

	// Sanity check
	if registryName == "" || name == "" {
		panic("invalid state - registry or name is empty")
	}

	log := slog.With(slog.String("registry", registryName), slog.String("image", name), slog.String("version", image.CurrentVersion))

	if registryName != "docker.io" {
		log.Warn("Skipping unsupported registry", slog.String("registry", registryName), slog.String("image", name), slog.String("version", image.CurrentVersion))
		return nil
	}

	cacheKey := fmt.Sprintf("v1/docker/latest-versions/%s/%s", name, image.CurrentVersion)
	var newImage *registry.Image
	if err := p.cache.GetJSON(ctx, cacheKey, &newImage, 24*time.Hour); err != nil {
		slog.Error("Failed to get cache", slog.Any("error", err))
		// Fallthrough
	}

	if newImage == nil {
		log.Debug("Fetching image releases")

		c := docker.Client{}
		var err error
		newImage, err = c.GetLatestVersion(ctx, name, image.CurrentVersion)
		if err != nil {
			slog.Error("Failed to get Docker Hub repository", slog.Any("error", err))
			return nil
		}

		if newImage != nil {
			if err := p.cache.SetJSON(ctx, cacheKey, &newImage); err != nil {
				slog.Error("Failed to set cache", slog.Any("error", err))
				// Fallthrough
			}
		}
	}

	if newImage == nil {
		log.Warn("No new image found", slog.String("image", name))
		return nil
	}

	if newImage != nil {
		image.LatestVersion = newImage.Version
	}

	return nil
}

func (p *Pipeline) EnrichFromGitHub(ctx context.Context, store *models.Store) error {
	// Temp before cache
	fetched := 0
	for _, image := range store.Images {
		if fetched > 0 {
			continue
		}

		// TODO: Don't rely on links for this?
		source := ""
		for _, link := range image.Links {
			if strings.HasPrefix(link.URL, "https://github.com/") {
				source = link.URL
				break
			}
		}

		if source == "" {
			continue
		}

		// TODO: Support other hosts?
		_, owner, repository, _, ok := github.ParseURL(source)
		if !ok {
			slog.Warn("Failed to parse GitHub URL", slog.String("url", source))
			continue
		}

		c := &github.Client{}
		release, err := c.GetRelease(ctx, owner, repository, image.LatestVersion)
		if err != nil {
			slog.Error("Failed to get GitHub release", slog.Any("error", err))
		}
		if release == nil {
			continue
		}

		image.Links = append(image.Links, &models.ImageLink{
			Type: "github-release",
			URL:  release.URL,
		})

		// NOTE: The current version is used as an identifier as we only ever show
		// images as being of the current version, with information that can
		// reference a new version
		notes := &models.ImageReleaseNotes{
			Title: release.Title,
			HTML:  release.Description,
		}
		if !release.Released.IsZero() {
			notes.Released = release.Released.UTC().Format(time.RFC3339)
		}
		store.ReleaseNotes[image.Name+":"+image.CurrentVersion] = notes

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
