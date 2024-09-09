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
	"github.com/AlexGustafsson/k8s-image-feed/internal/registry/oci"
)

type Client struct {
	Client *http.Client
}

func (c *Client) Get(ctx context.Context, name string, version string) (*registry.Image, error) {
	dockerName := name
	if !strings.Contains(dockerName, "/") {
		dockerName = "library/" + dockerName
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags/%s", dockerName, version), nil)
	if err != nil {
		return nil, err
	}

	client := c.Client
	if client == nil {
		client = http.DefaultClient
	}

	res, err := client.Do(req)
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

func (c *Client) GetRegistryToken(ctx context.Context, name string) (string, error) {
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

	client := c.Client
	if client == nil {
		client = http.DefaultClient
	}

	res, err := client.Do(req)
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

func (c *Client) GetManifests(ctx context.Context, name string, tag string) ([]oci.Manifest, error) {
	token, err := c.GetRegistryToken(ctx, name)
	if err != nil {
		return nil, err
	}

	ociClient := &oci.Client{
		Registry: "https://registry-1.docker.io/v2",
		Client:   c.Client,
		// TODO: Cache token
		Authorizer: oci.AuthorizerToken(token),
	}

	return ociClient.GetManifests(ctx, name, tag)
}

func (c *Client) GetLatestVersion(ctx context.Context, name string, currentTag string) (*registry.Image, error) {
	currentVersion, err := oci.ParseVersion(currentTag)
	if err != nil || currentVersion == nil {
		return nil, fmt.Errorf("unsupported version: %s", err)
	}

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

	client := c.Client
	if client == nil {
		client = http.DefaultClient
	}

	res, err := client.Do(req)
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

	// TODO:
	// As we've sorted versions in released time, let's assume the first version
	// that is higher than ours, is the latest version. Might not be true if the
	// current version is 1.0.0, there have been a lot of nightlies or other types
	// of tags, so that the page contains only fix 1.0.1, but in reality 2.0.0 was
	// released a while ago and would be on the next page, would we be greedy.
	// Look at any large image with LTS, such as postgres, node.
	for _, tag := range result.Results {
		if tag.Name == "" {
			continue
		}

		newVersion, err := oci.ParseVersion(tag.Name)
		if err != nil || newVersion == nil {
			continue
		}

		if newVersion.IsCompatible(currentVersion) && newVersion.Compare(currentVersion) >= 0 {
			return &registry.Image{
				Name:         name,
				Version:      tag.Name,
				Published:    tag.TagLastPushed,
				Digest:       tag.Digest,
				ReleaseNotes: "", // TODO
			}, nil
		}
	}

	return nil, nil
}

func (c *Client) GetRepository(ctx context.Context, owner string, name string) (*Repository, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/%s", url.PathEscape(owner), url.PathEscape(name)), nil)
	if err != nil {
		return nil, err
	}

	client := c.Client
	if client == nil {
		client = http.DefaultClient
	}

	res, err := client.Do(req)
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
