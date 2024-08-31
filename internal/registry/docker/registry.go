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

// GetLatestVersion implements registry.Registry.
func (r *Registry) GetLatestVersion(ctx context.Context, name string) (*registry.Image, error) {
	dockerName := name
	if !strings.Contains(dockerName, "/") {
		dockerName = "library/" + dockerName
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
