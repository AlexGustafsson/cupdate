package oci

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
)

type Authorizer interface {
	AuthorizeOCIRequest(context.Context, Reference, *http.Request) error
}

type AuthorizeFunc func(context.Context, Reference, *http.Request) error

func (f AuthorizeFunc) AuthorizeOCIRequest(ctx context.Context, image Reference, req *http.Request) error {
	return f(ctx, image, req)
}

type AuthorizerToken string

func (s AuthorizerToken) AuthorizeOCIRequest(ctx context.Context, image Reference, req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+string(s))
	return nil
}

type Client struct {
	Client     *httputil.Client
	Authorizer Authorizer
}

func (c *Client) GetManifests(ctx context.Context, image Reference) ([]Manifest, error) {
	id := ""
	if image.HasTag {
		id = image.Tag
	} else if image.HasDigest {
		id = image.Digest
	} else {
		return nil, fmt.Errorf("unsupported reference type: must be tagged or digested")
	}

	// NOTE: It's rather unclear why we need to do this dance manually and why
	// docker.io simply doesn't just redirect us
	domain := strings.Replace(image.Domain, "docker.io", "registry-1.docker.io", 1)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/v2/%s/manifests/%s", domain, image.Path, url.PathEscape(id)), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", strings.Join([]string{
		"application/vnd.docker.distribution.manifest.list.v2+json",
		"application/vnd.oci.image.index.v1+json",
		"application/vnd.docker.distribution.manifest.v2+json",
	}, ", "))

	if c.Authorizer != nil {
		if err := c.Authorizer.AuthorizeOCIRequest(ctx, image, req); err != nil {
			return nil, err
		}
	}

	res, err := c.Client.DoCached(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if res.StatusCode != http.StatusOK {
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
			// The manifest's schema version always seems to be unset in this case,
			// fall back to use the parent manifest's version
			schemaVersion := manifest.SchemaVersion
			if schemaVersion == 0 {
				schemaVersion = result.SchemaVersion
			}

			var platform *Platform
			if manifest.Platform.Architecture != "" {
				platform = &Platform{
					Architecture: manifest.Platform.Architecture,
					OS:           manifest.Platform.OS,
					Variant:      manifest.Platform.Variant,
				}
			}

			manifests = append(manifests, Manifest{
				SchemaVersion: schemaVersion,
				MediaType:     manifest.MediaType,
				Digest:        manifest.Digest,
				Platform:      platform,
			})
		}

		return manifests, nil
	}

	if result.MediaType == "application/vnd.docker.distribution.manifest.v2+json" && result.SchemaVersion == 2 {
		var manifest DockerDistributionManifestV2
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		var platform *Platform
		if manifest.Platform.Architecture != "" {
			platform = &Platform{
				Architecture: manifest.Platform.Architecture,
				OS:           manifest.Platform.OS,
				Variant:      manifest.Platform.Variant,
			}
		}

		return []Manifest{
			{
				SchemaVersion: manifest.SchemaVersion,
				MediaType:     manifest.MediaType,
				Digest:        manifest.Digest,
				Platform:      platform,
			},
		}, nil
	}

	if result.MediaType == "application/vnd.oci.image.index.v1+json" && result.SchemaVersion == 2 {
		var result OCIImageIndexV1
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		manifests := make([]Manifest, 0)
		for _, manifest := range result.Manifests {
			// The manifest's schema version always seems to be unset in this case,
			// fall back to use the parent manifest's version
			schemaVersion := manifest.SchemaVersion
			if schemaVersion == 0 {
				schemaVersion = result.SchemaVersion
			}

			var platform *Platform
			if manifest.Platform.Architecture != "" {
				platform = &Platform{
					Architecture: manifest.Platform.Architecture,
					OS:           manifest.Platform.OS,
					Variant:      manifest.Platform.Variant,
				}
			}

			manifests = append(manifests, Manifest{
				SchemaVersion: schemaVersion,
				MediaType:     manifest.MediaType,
				Annotations:   manifest.Annotations,
				Digest:        manifest.Digest,
				Platform:      platform,
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

	return nil, fmt.Errorf("unsupported manifest type: %s (%s, %d)", res.Header["Content-Type"], result.MediaType, result.SchemaVersion)
}

func (c *Client) GetManifest(ctx context.Context, image Reference, digest string) ([]byte, error) {
	// NOTE: It's rather unclear why we need to do this dance manually and why
	// docker.io simply doesn't just redirect us
	domain := strings.Replace(image.Domain, "docker.io", "registry-1.docker.io", 1)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/v2/%s/manifests/%s", domain, image.Path, digest), nil)
	if err != nil {
		return nil, err
	}

	if c.Authorizer != nil {
		if err := c.Authorizer.AuthorizeOCIRequest(ctx, image, req); err != nil {
			return nil, err
		}
	}

	req.Header.Set("Accept", strings.Join([]string{
		"application/vnd.docker.distribution.manifest.list.v2+json",
		"application/vnd.oci.image.index.v1+json",
		"application/vnd.docker.distribution.manifest.v2+json",
		"application/vnd.oci.image.manifest.v1+json",
	}, ", "))

	res, err := c.Client.DoCached(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	return io.ReadAll(res.Body)
}

func (c *Client) GetBlob(ctx context.Context, image Reference, digest string) ([]byte, error) {
	// NOTE: It's rather unclear why we need to do this dance manually and why
	// docker.io simply doesn't just redirect us
	domain := strings.Replace(image.Domain, "docker.io", "registry-1.docker.io", 1)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/v2/%s/blobs/%s", domain, image.Path, digest), nil)
	if err != nil {
		return nil, err
	}

	if c.Authorizer != nil {
		if err := c.Authorizer.AuthorizeOCIRequest(ctx, image, req); err != nil {
			return nil, err
		}
	}

	res, err := c.Client.DoCached(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	return io.ReadAll(res.Body)
}

type GetAnnotationsOptions struct {
	Manifests    []Manifest
	Digest       string
	Architecture string
	OS           string
	Variant      string
}

// GetAnnotations tries to identify annotations for the image.
// Fetches manifests as necessary.
// To narrow down the search and to avoid unnecessary fetches, specify the
// available options.
func (c *Client) GetAnnotations(ctx context.Context, image Reference, options *GetAnnotationsOptions) (Annotations, error) {
	if options == nil {
		options = &GetAnnotationsOptions{}
	}

	var manifests []Manifest
	if options.Manifests == nil {
		var err error
		manifests, err = c.GetManifests(ctx, image)
		if err != nil {
			return nil, err
		}
	} else {
		manifests = make([]Manifest, len(options.Manifests))
		copy(manifests, options.Manifests)
	}

	if options.Digest != "" {
		manifests = slices.DeleteFunc(manifests, func(m Manifest) bool {
			return m.Digest != options.Digest
		})
	}

	if options.Architecture != "" {
		manifests = slices.DeleteFunc(manifests, func(m Manifest) bool {
			return m.Platform == nil || m.Platform.Architecture != options.Architecture
		})
	}

	if options.OS != "" {
		manifests = slices.DeleteFunc(manifests, func(m Manifest) bool {
			return m.Platform == nil || m.Platform.OS != options.OS
		})
	}

	if options.Variant != "" {
		manifests = slices.DeleteFunc(manifests, func(m Manifest) bool {
			return m.Platform == nil || m.Platform.Variant != options.Variant
		})
	}

	if len(manifests) == 0 {
		return nil, nil
	}

	// Pick the first manifest
	manifest := manifests[0]

	// If the manifest doesn't have a digest, it cannot be fetched
	if manifest.Digest == "" {
		return nil, nil
	}

	// Fetch the blob for the manifest (trying to get labels from the config)
	blob, err := c.GetManifest(ctx, image, manifest.Digest)
	if err != nil {
		return nil, err
	} else if blob == nil {
		return nil, fmt.Errorf("manifest blob not found")
	}

	var manifestBlob struct {
		Config struct {
			Digest string `json:"digest"`
		} `json:"config"`
	}
	if err := json.Unmarshal(blob, &manifestBlob); err != nil {
		return nil, err
	}

	// The blob was probably not a manifest blob
	if manifestBlob.Config.Digest == "" {
		return nil, nil
	}

	blob, err = c.GetBlob(ctx, image, manifestBlob.Config.Digest)
	if err != nil {
		return nil, err
	} else if blob == nil {
		return nil, fmt.Errorf("manifest config blob not found")
	}

	var configBlob struct {
		Config struct {
			Labels map[string]string `json:"Labels"`
		} `json:"config"`
	}
	if err := json.Unmarshal(blob, &configBlob); err != nil {
		return nil, err
	}

	return configBlob.Config.Labels, nil
}

type GetTagsOptions struct {
	Last     string
	Count    int
	AllPages bool
}

func (c *Client) GetTags(ctx context.Context, image Reference, options *GetTagsOptions) ([]string, error) {
	tags, origin, linkHeader, err := c.getTags(ctx, image, options)
	if err != nil {
		return nil, err
	}

	allTags := append([]string{}, tags...)

	// Follow pagination
	for linkHeader != "" && options != nil && options.AllPages {
		links, err := httputil.ParseLinkHeader(origin, linkHeader)
		if err != nil {
			return nil, err
		}

		var next *httputil.Link
		for _, link := range links {
			if link.Params["rel"] == "next" {
				next = &link
				break
			}
		}
		if next == nil {
			break
		}

		// As a precaution, don't leave the origin
		if next.URL.Host != origin.Host {
			return nil, fmt.Errorf("refusing to follow link to other origin")
		}

		query := next.URL.Query()

		options := GetTagsOptions{}

		if query.Has("n") {
			n, err := strconv.ParseInt(query.Get("n"), 10, 32)
			if err != nil {
				return nil, err
			}

			options.Count = int(n)
		}

		options.Last = query.Get("last")

		tags, _, linkHeader, err = c.getTags(ctx, image, &options)
		if err != nil {
			return nil, err
		}

		allTags = append(allTags, tags...)
	}

	return allTags, nil
}

func (c *Client) getTags(ctx context.Context, image Reference, options *GetTagsOptions) ([]string, *url.URL, string, error) {
	// NOTE: It's rather unclear why we need to do this dance manually and why
	// docker.io simply doesn't just redirect us
	domain := strings.Replace(image.Domain, "docker.io", "registry-1.docker.io", 1)

	u, err := url.Parse(fmt.Sprintf("https://%s/v2/%s/tags/list", domain, image.Path))
	if err != nil {
		return nil, nil, "", err
	}

	query := make(url.Values)
	if options != nil && options.Last != "" {
		query.Set("last", options.Last)
	}
	if options != nil && options.Count > 0 {
		query.Set("n", strconv.FormatInt(int64(options.Count), 10))
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, nil, "", err
	}

	req.Header.Set("Accept", "application/json")

	if c.Authorizer != nil {
		if err := c.Authorizer.AuthorizeOCIRequest(ctx, image, req); err != nil {
			return nil, nil, "", err
		}
	}

	res, err := c.Client.DoCached(req)
	if err != nil {
		return nil, nil, "", err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, nil, "", nil
	} else if res.StatusCode != http.StatusOK {
		return nil, nil, "", fmt.Errorf("unexpected status code: %s", res.Status)
	}

	var page TagsPage
	if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
		return nil, nil, "", err
	}

	tags := page.Tags

	return tags, req.URL, res.Header.Get("Link"), nil
}
