package oci

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	Client     *http.Client
	Authorizer Authorizer
}

func (c *Client) GetManifests(ctx context.Context, image Reference) ([]Manifest, error) {
	// NOTE: It's rather unclear why we need to do this dance manually and why
	// docker.io simply doesn't just redirect us
	id := ""
	if image.HasTag {
		id = image.Tag
	} else if image.HasDigest {
		id = image.Digest
	} else {
		return nil, fmt.Errorf("unsupported reference type: must be tagged or digested")
	}
	domain := strings.Replace(image.Domain, "docker.io", "registry-1.docker.io/v2", 1)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/%s/manifests/%s", domain, image.Path, url.PathEscape(id)), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", strings.Join([]string{
		"application/vnd.docker.distribution.manifest.list.v2+json",
		"application/vnd.oci.image.index.v1+json",
		// These two formats never occur?
		// "application/vnd.docker.distribution.manifest.v2+json",
		// "application/vnd.oci.image.manifest.v1+json",
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

	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if res.StatusCode != http.StatusOK {
		fmt.Println(res.Header)
		x, _ := io.ReadAll(res.Body)
		fmt.Printf("%s\n", x)
		fmt.Printf("%s\n", res.Request.URL.String())
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		SchemaVersion int    `json:"schemaVersion"`
		MediaType     string `json:"mediaType"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	contentType := res.Header.Get("Content-Type")

	if result.MediaType == "application/vnd.docker.distribution.manifest.list.v2+json" && result.SchemaVersion == 2 {
		var result DockerDistributionManifestListV2
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		manifests := make([]Manifest, 0)
		for _, manifest := range result.Manifests {
			manifests = append(manifests, Manifest{
				SchemaVersion: manifest.SchemaVersion,
				MediaType:     manifest.MediaType,
				Annotations:   make(map[string]string),
			})
		}

		return manifests, nil
	}

	if result.MediaType == "application/vnd.oci.image.index.v1+json" && result.SchemaVersion == 2 {
		var result OCIImageIndexV1
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		manifests := make([]Manifest, 0)
		for _, manifest := range result.Manifests {
			manifests = append(manifests, Manifest{
				SchemaVersion: manifest.SchemaVersion,
				MediaType:     manifest.MediaType,
				Annotations:   manifest.Annotations,
			})
		}

		return manifests, nil
	}

	if contentType == "application/vnd.docker.distribution.manifest.v1+prettyjws" && result.SchemaVersion == 1 {
		var result DockerDistributionManifestV1
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		return make([]Manifest, 0), nil
	}

	return nil, fmt.Errorf("unsupported manifest type: %s (%s ,%d)", res.Header["Content-Type"], result.MediaType, result.SchemaVersion)
}
