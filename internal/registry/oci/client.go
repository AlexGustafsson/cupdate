package oci

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Authorizer interface {
	Authorize(context.Context, *http.Request) error
}

type AuthorizeFunc func(context.Context, *http.Request) error

func (f AuthorizeFunc) Authorize(ctx context.Context, req *http.Request) error {
	return f(ctx, req)
}

type AuthorizerToken string

func (s AuthorizerToken) Authorize(ctx context.Context, req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+string(s))
	return nil
}

type Client struct {
	// Registry is the registry endpoint, such as https://registry-1.docker.io/v2.
	Registry   string
	Client     *http.Client
	Authorizer Authorizer
}

func (c *Client) GetManifests(ctx context.Context, name string, tag string) ([]Manifest, error) {
	dockerName := name
	if !strings.Contains(dockerName, "/") {
		dockerName = "library/" + url.PathEscape(dockerName)
	}

	// NOTE: If name contains a /, it is not path escaped as we can't easily tell
	// what part is the namespace or not
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/manifests/%s", c.Registry, dockerName, url.PathEscape(tag)), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", strings.Join([]string{
		"application/vnd.docker.distribution.manifest.list.v2+json",
		"application/vnd.docker.distribution.manifest.v2+json",
		"application/vnd.oci.image.manifest.v1+json",
		"application/vnd.oci.image.index.v1+json",
	}, ", "))

	if c.Authorizer != nil {
		if err := c.Authorizer.Authorize(ctx, req); err != nil {
			return nil, err
		}
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
