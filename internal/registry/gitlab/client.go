package gitlab

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/registry"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

type Client struct {
	Client *httputil.Client
}

func (c *Client) GetRegistryToken(ctx context.Context, image oci.Reference) (string, error) {
	// TODO: Registries expose the realm and scheme via Www-Authenticate if 403
	// is given
	u, err := url.Parse("https://gitlab.com/jwt/auth?service=container_registry")
	if err != nil {
		return "", err
	}

	query := u.Query()
	query.Set("scope", fmt.Sprintf("repository:%s:pull", image.Path))
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}

	// TODO: The cache doesn't understand graphql, so we can't cache this request
	res, err := c.Client.Do(req)
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

func (c *Client) Authorize(ctx context.Context, image oci.Reference, req *http.Request) error {
	token, err := c.GetRegistryToken(ctx, image)
	if err != nil {
		return err
	}

	return oci.AuthorizerToken(token).Authorize(ctx, image, req)
}

func (c *Client) GetProjectContainerRepositories(ctx context.Context, fullPath string) ([]ContainerRepository, error) {
	payload, err := json.Marshal(map[string]any{
		"operationName": "getProjectContainerRepositories",
		"variables": map[string]any{
			"sort":     "UPDATED_DESC",
			"fullPath": fullPath,
			"first":    20,
		},
		"query": `query getProjectContainerRepositories($fullPath: ID!, $name: String, $first: Int, $last: Int, $after: String, $before: String, $sort: ContainerRepositorySort) {
  project(fullPath: $fullPath) {
    id
    containerRepositories(
      name: $name
      after: $after
      before: $before
      first: $first
      last: $last
      sort: $sort
    ) {
      nodes {
        id
        location
      }
    }
  }
}`,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://gitlab.com/api/graphql", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// TODO: The cache doesn't understand graphql, so we can't cache this request
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	var result struct {
		Errors []any `json:"error"`
		Data   struct {
			Project struct {
				ContainerRepositories struct {
					Nodes []ContainerRepository `json:"nodes"`
				} `json:"containerRepositories"`
			} `json:"project"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	} else if len(result.Errors) > 0 {
		return nil, fmt.Errorf("got unexpected graphql error: %v", result.Errors)
	}

	return result.Data.Project.ContainerRepositories.Nodes, nil
}
func (c *Client) GetProjectContainerRepositoryTags(ctx context.Context, id string) ([]ContainerRepositoryTag, error) {
	payload, err := json.Marshal(map[string]any{
		"operationName": "getContainerRepositoryTags",
		"variables": map[string]any{
			"referrers": true,
			"id":        id,
			"first":     20,
			"sort":      "PUBLISHED_AT_DESC",
		},
		"query": `query getContainerRepositoryTags($id: ContainerRepositoryID!, $first: Int, $last: Int, $after: String, $before: String, $name: String, $sort: ContainerRepositoryTagSort, $referrers: Boolean = false) {
  containerRepository(id: $id) {
    tags(
      after: $after
      before: $before
      first: $first
      last: $last
      name: $name
      sort: $sort
      referrers: $referrers
    ) {
      nodes {
        digest
        location
        path
        name
        createdAt
        publishedAt
      }
    }
  }
}`})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://gitlab.com/api/graphql", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	var result struct {
		Errors []any `json:"error"`
		Data   struct {
			ContainerRepository struct {
				Tags struct {
					Nodes []ContainerRepositoryTag `json:"nodes"`
				} `json:"tags"`
			} `json:"containerRepository"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	} else if len(result.Errors) > 0 {
		return nil, fmt.Errorf("got unexpected graphql error: %v", result.Errors)
	}

	return result.Data.ContainerRepository.Tags.Nodes, nil
}

func (c *Client) GetLatestVersion(ctx context.Context, image oci.Reference) (*registry.Image, error) {
	if !image.HasTag {
		return nil, nil
	}

	// There's not going to be any latest version
	if image.Tag == "latest" {
		return nil, nil
	}

	currentVersion, err := oci.ParseVersion(image.Tag)
	if err != nil {
		return nil, fmt.Errorf("unsupported version: %w", err)
	} else if currentVersion == nil {
		return nil, fmt.Errorf("unsupported version")
	}

	// The repository path is <owner>/<group>/<project>
	parts := strings.Split(image.Path, "/")
	if len(parts) < 3 {
		return nil, nil
	}

	fullPath := strings.Join(parts[0:3], "/")
	repositories, err := c.GetProjectContainerRepositories(ctx, fullPath)
	if err != nil {
		return nil, err
	}

	var repository *ContainerRepository
	for i := range repositories {
		if repositories[i].Location == image.Name() {
			r := repositories[i]
			repository = &r
			break
		}
	}

	if repository == nil {
		return nil, nil
	}

	tags, err := c.GetProjectContainerRepositoryTags(ctx, repository.ID)
	if err != nil {
		return nil, err
	}

	// TODO:
	// As we've sorted versions in released time, let's assume the first version
	// that is higher than ours, is the latest version. Might not be true if the
	// current version is 1.0.0, there have been a lot of nightlies or other types
	// of tags, so that the page contains only fix 1.0.1, but in reality 2.0.0 was
	// released a while ago and would be on the next page, would we be greedy.
	// Look at any large image with LTS, such as postgres, node.
	for _, tag := range tags {
		if tag.Name == "" {
			continue
		}

		newVersion, err := oci.ParseVersion(tag.Name)
		if err != nil || newVersion == nil {
			continue
		}

		if newVersion.IsCompatible(currentVersion) && newVersion.Compare(currentVersion) >= 0 {
			image.Tag = tag.Name

			return &registry.Image{
				Name:      image,
				Published: tag.PublishedAt,
				Digest:    tag.Digest,
			}, nil
		}
	}

	return nil, nil
}

type ContainerRepository struct {
	ID       string `json:"id"`
	Location string `json:"location"`
}

type ContainerRepositoryTag struct {
	Digest      string    `json:"digest"`
	Location    string    `json:"location"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
	PublishedAt time.Time `json:"publishedAt"`

	// ... unused fields
}
