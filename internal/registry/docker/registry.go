package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AlexGustafsson/k8s-image-feed/internal/registry"
)

var _ registry.Registry = (*Registry)(nil)

type Registry struct {
	Client http.Client
}

// Get implements registry.Registry.
func (r *Registry) Get(ctx context.Context, name string, version string) (*registry.Image, error) {
	dockerName := name
	if !strings.Contains(dockerName, "/") {
		dockerName = "library/" + dockerName
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags/%s", dockerName, version), nil)
	if err != nil {
		return nil, err
	}

	res, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	var result Tag
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	fmt.Println(result.Name, result.Digest, result.LastUpdated)

	return &registry.Image{
		Name:         name,
		Version:      result.Name,
		Published:    result.LastUpdated,
		Digest:       result.Digest,
		ReleaseNotes: "",
	}, nil
}

func (r *Registry) GetRegistryToken(ctx context.Context, name string) (string, error) {
	dockerName := name
	if !strings.Contains(dockerName, "/") {
		dockerName = "library/" + url.PathEscape(dockerName)
	}

	// TODO: Registries expose the realm and scheme via Www-Authenticate if 403
	// is given
	u, err := url.Parse("https://auth.docker.io/token?service=registry.docker.io")
	if err != nil {
		return "", err
	}

	query := u.Query()
	query.Set("scope", fmt.Sprintf("repository:%s:pull", dockerName))
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}

	res, err := r.Client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %s", res.Status)
	}

	var result struct {
		Token     string    `json:"token"`
		ExpiresIn int       `json:"expires_in"`
		IssuedAt  time.Time `json:"issued_at"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Token, nil
}

func (r *Registry) GetManifests(ctx context.Context, name string, tag string) ([]Manifest, error) {
	dockerName := name
	if !strings.Contains(dockerName, "/") {
		dockerName = "library/" + url.PathEscape(dockerName)
	}

	// NOTE: If name contains a /, it is not path escaped as we can't easily tell
	// what part is the namespace or not
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://registry-1.docker.io/v2/%s/manifests/%s", dockerName, url.PathEscape(tag)), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", strings.Join([]string{
		"application/vnd.docker.distribution.manifest.list.v2+json",
		"application/vnd.docker.distribution.manifest.v2+json",
		"application/vnd.oci.image.manifest.v1+json",
		"application/vnd.oci.image.index.v1+json",
	}, ", "))

	token, err := r.GetRegistryToken(ctx, name)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	var result struct {
		SchemaVersion int        `json:"schemaVersion"`
		MediaType     string     `json:"mediaType"`
		Manifests     []Manifest `json:"manifests"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	// TODO: Handle content type switch
	if result.MediaType != "application/vnd.oci.image.index.v1+json" || result.SchemaVersion != 2 {
		return nil, fmt.Errorf("unsupported manifest type")
	}

	return result.Manifests, nil
}

// GetLatestVersion implements registry.Registry.
func (r *Registry) GetLatestVersion(ctx context.Context, name string) (*registry.Image, error) {
	dockerName := name
	if !strings.Contains(dockerName, "/") {
		dockerName = "library/" + url.PathEscape(dockerName)
	}

	u, err := url.Parse(fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags", dockerName))
	if err != nil {
		return nil, err
	}

	query := u.Query()
	query.Set("page_size", "25")
	query.Set("page", "1")
	query.Set("ordering", "last_updated")
	query.Set("name", "")
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	var result Page[Tag]
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	for _, tag := range result.Results {
		if tag.Name != "" {
			fmt.Println(tag.Name, tag.Digest, tag.LastUpdated)
		}
	}

	return nil, fmt.Errorf("not implemented")
}

func (r *Registry) GetRepository(ctx context.Context, owner string, name string) (*Repository, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/%s", url.PathEscape(owner), url.PathEscape(name)), nil)
	if err != nil {
		return nil, err
	}

	res, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	var result Repository
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type Page[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

type Tag struct {
	Creator             int       `json:"creator"`
	ID                  int       `json:"id"`
	LastUpdated         time.Time `json:"last_updated"`
	LastUpdater         int       `json:"last_updater"`
	LastUpdaterUsername string    `json:"last_updater_username"`
	Name                string    `json:"name"`
	Repository          int       `json:"repository"`
	FullSize            int       `json:"full_size"`
	V2                  bool      `json:"v2"`
	TagStatus           string    `json:"tag_status"`
	TagLastPulled       time.Time `json:"tag_last_pulled"`
	TagLastPushed       time.Time `json:"tag_last_pushed"`
	MediaType           string    `json:"media_type"`
	ContentType         string    `json:"content_type"`
	Digest              string    `json:"digest"`
	Images              []Image   `json:"images"`
}

type Image struct {
	Architecture string    `json:"architecture"`
	Features     string    `json:"features"`
	Variant      *string   `json:"variant"`
	Digest       string    `json:"digest"`
	OS           string    `json:"os"`
	OSFeatures   string    `json:"os_features"`
	OSVersion    *string   `json:"os_version"`
	Size         int       `json:"size"`
	Status       string    `json:"status"`
	LastPulled   time.Time `json:"last_pulled"`
	LastPushed   time.Time `json:"last_pushed"`
}

type Repository struct {
	User              string          `json:"user"`
	Name              string          `json:"name"`
	Namespace         string          `json:"namespace"`
	Type              string          `json:"repository_type"`
	Status            int             `json:"status"`
	StatusDescription string          `json:"status_description"`
	Description       string          `json:"description"`
	IsPrivate         bool            `json:"is_private"`
	IsAutomated       bool            `json:"is_automated"`
	StarCount         int             `json:"star_count"`
	PullCount         int             `json:"pull_count"`
	LastUpdated       time.Time       `json:"last_updated"`
	DateRegistered    time.Time       `json:"date_registered"`
	CollaboratorCount int             `json:"collaborator_count"`
	Affiliation       json.RawMessage `json:"affiliation"` // Unknown
	HubUser           string          `json:"hub_user"`
	HasStarred        bool            `json:"has_starred"`
	FullDescription   string          `json:"full_description"`
	Permissions       struct {
		Read  bool `json:"read"`
		Write bool `json:"write"`
		Admin bool `json:"admin"`
	} `json:"permissions"`
	MediaTypes   []string `json:"media_types"`
	ContentTypes []string `json:"content_types"`
	Categories   []struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"categories"`
	ImmutableTags      bool   `json:"immutable_tags"`
	ImmutableTagsRules string `json:"immutable_tags_rules"`
}

// TODO: Rearrange - OCI is the common package, then packages for docker, ghcr,
// google, ecr etc.?

type Manifest struct {
	Annotations map[string]string `json:"annotations"`
	Digest      string            `json:"digest"`
	MediaType   string            `json:"mediaType"`
	Platform    struct {
		Architecture string `json:"architecture"`
		OS           string `json:"os"`
	} `json:"platform"`
	Size int `json:"size"`
}
