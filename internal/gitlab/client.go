package gitlab

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
)

var readmePathRegexp = regexp.MustCompile(`href="(.*?/blob/.*?)"`)

type Client struct {
	Client *httputil.Client
}

func (c *Client) GetRepositoryDescription(ctx context.Context, fullPath string) (string, error) {
	payload, err := json.Marshal(map[string]any{
		"operationName": "getProject",
		"variables": map[string]any{
			"fullPath": fullPath,
		},
		"query": `query getProject($fullPath: ID!) {
  project(fullPath: $fullPath) {
    description
  }
}`,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://gitlab.com/api/graphql", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// TODO: The cache doesn't understand graphql, so we can't cache this request
	res, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return "", err
	}

	var result struct {
		Errors []any `json:"error"`
		Data   struct {
			Project struct {
				Description string `json:"description"`
			} `json:"project"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	} else if len(result.Errors) > 0 {
		return "", fmt.Errorf("got unexpected graphql error: %v", result.Errors)
	}

	return result.Data.Project.Description, nil
}

func (c *Client) GetRepositoryREADMEBlob(ctx context.Context, fullPath string) (*Blob, error) {
	// https://gitlab.com/arm-research/smarter/smarter-device-manager/-/blob/master/README.md?format=json&viewer=rich
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://gitlab.com/"+fullPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Client.DoCached(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return nil, err
	}

	html, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// The API itself doesn't expose the README. The UI has some "startup" calls
	// defined that include the README, whichever file it is.

	readmeHref := readmePathRegexp.FindSubmatch(html)
	if readmeHref == nil {
		return nil, nil
	}

	return c.GetBlob(ctx, string(readmeHref[1]), true)
}

func (c *Client) GetBlob(ctx context.Context, href string, includeRaw bool) (*Blob, error) {
	u := url.URL{
		Scheme:   "https",
		Host:     "gitlab.com",
		Path:     href,
		RawQuery: "format=json",
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.DoCached(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return nil, err
	}

	var blob Blob
	if err := json.NewDecoder(res.Body).Decode(&blob); err != nil {
		return nil, err
	}

	if includeRaw {
		u := url.URL{
			Scheme: "https",
			Host:   "gitlab.com",
			Path:   blob.RawPath,
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}

		res, err := c.Client.DoCached(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
			return nil, err
		}

		raw, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		blob.Raw = raw
	}

	return &blob, nil
}

func (c *Client) GetRegistryToken(ctx context.Context, repository string) (string, error) {
	// TODO: Registries expose the realm and scheme via Www-Authenticate if 403
	// is given
	u, err := url.Parse("https://gitlab.com/jwt/auth?service=container_registry")
	if err != nil {
		return "", err
	}

	query := u.Query()
	query.Set("scope", fmt.Sprintf("repository:%s:pull", repository))
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

	if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return "", err
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

func (c *Client) HandleAuth(r *http.Request) error {
	name := oci.NameFromAPI(r.URL.Path)
	// lscr.io is a pseudo-registry that forwards to one of multiple backends,
	// among them registry.gitlab.com
	if (r.Host != "registry.gitlab.com" && r.Host != "lscr.io") || name == "" {
		return nil
	}

	token, err := c.GetRegistryToken(r.Context(), name)
	if err != nil {
		return err
	}

	r.Header.Set("Authorization", "Bearer "+token)

	return nil
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

	if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return nil, err
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

	if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return nil, err
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
