package oci

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
)

type Client struct {
	Client   httputil.Requester
	AuthFunc func(*http.Request) error
}

// GetManifestBlob downloads a manifest from an OCI registry.
// SEE: https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pull
func (c *Client) GetManifestBlob(ctx context.Context, ref Reference) (Blob, error) {
	ref = c.rewriteReference(ref)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/v2/%s/manifests/%s", ref.Domain, ref.Path, ref.Reference()), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", strings.Join([]string{
		"application/vnd.oci.image.index.v1+json",
		"application/vnd.oci.image.manifest.v1+json",
		"application/vnd.docker.distribution.manifest.list.v2+json",
		"application/vnd.docker.distribution.manifest.v2+json",
		"application/vnd.docker.distribution.manifest.v1+prettyjws",
		"application/vnd.docker.distribution.manifest.v1+json",
	}, ", "))

	if f := c.AuthFunc; f != nil {
		if err := f(req); err != nil {
			return nil, err
		}
	}

	res, err := c.Client.DoCached(req)
	if err != nil {
		return nil, err
	}

	if err := assertStatusCode(res, http.StatusOK); err != nil {
		return nil, err
	}

	return blobFromResponse(res), nil
}

// GetManifest downloads a [Manifest] or a [ManifestIndex] from an OCI registry.
// SEE: https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pull
func (c *Client) GetManifest(ctx context.Context, ref Reference) (any, error) {
	ref = c.rewriteReference(ref)

	blob, err := c.GetManifestBlob(ctx, ref)
	if err != nil {
		return nil, err
	}

	return manifestFromBlob(blob)
}

// GetBlob downloads a blob from an OCI registry.
// SEE: https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pulling-blobs
func (c *Client) GetBlob(ctx context.Context, ref Reference, digest string, cache bool) (Blob, error) {
	ref = c.rewriteReference(ref)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/v2/%s/blobs/%s", ref.Domain, ref.Path, digest), nil)
	if err != nil {
		return nil, err
	}

	if f := c.AuthFunc; f != nil {
		if err := f(req); err != nil {
			return nil, err
		}
	}

	req.Header.Set("Accept", strings.Join([]string{
		"application/vnd.docker.distribution.manifest.list.v2+json",
		"application/vnd.oci.image.index.v1+json",
		"application/vnd.docker.distribution.manifest.v2+json",
		"application/vnd.oci.image.manifest.v1+json",
	}, ", "))

	var res *http.Response
	if cache {
		res, err = c.Client.DoCached(req)
	} else {
		res, err = c.Client.DoCached(req)
	}
	if err != nil {
		return nil, err
	}

	if err := assertStatusCode(res, http.StatusOK); err != nil {
		return nil, err
	}

	return blobFromResponse(res), nil
}

// HeadBlob gets information about a blob from an OCI registry.
// SEE: https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pulling-blobs
func (c *Client) HeadBlob(ctx context.Context, ref Reference, digest string) (*BlobInfo, error) {
	ref = c.rewriteReference(ref)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, fmt.Sprintf("https://%s/v2/%s/blobs/%s", ref.Domain, ref.Path, digest), nil)
	if err != nil {
		return nil, err
	}

	if f := c.AuthFunc; f != nil {
		if err := f(req); err != nil {
			return nil, err
		}
	}

	res, err := c.Client.DoCached(req)
	if err != nil {
		return nil, err
	}

	if err := assertStatusCode(res, http.StatusOK); err != nil {
		return nil, err
	}

	info := blobInfoFromResponse(res)
	return &info, nil
}

type GetAnnotationsOptions struct {
	Manifests    []ImageManifest
	Digest       string
	Architecture string
	OS           string
	Variant      string
}

// GetAnnotations tries to identify annotations for the reference.
// Fetches manifests as necessary.
// To narrow down the search and to avoid unnecessary fetches, specify the
// available options.
func (c *Client) GetAnnotations(ctx context.Context, ref Reference, options *GetAnnotationsOptions) (Annotations, error) {
	if options == nil {
		options = &GetAnnotationsOptions{}
	}

	manifests := options.Manifests
	if manifests == nil {
		manifestOrIndex, err := c.GetManifest(ctx, ref)
		if err != nil {
			return nil, err
		}

		switch m := manifestOrIndex.(type) {
		case *ImageManifest:
			manifests = []ImageManifest{*m}
		case *ImageIndex:
			manifests = m.Manifests
		}
	}

	if options.Digest != "" {
		manifests = slices.DeleteFunc(manifests, func(m ImageManifest) bool {
			return m.Digest != options.Digest
		})
	}

	if options.Architecture != "" {
		manifests = slices.DeleteFunc(manifests, func(m ImageManifest) bool {
			return m.Platform == nil || m.Platform.Architecture != options.Architecture
		})
	}

	if options.OS != "" {
		manifests = slices.DeleteFunc(manifests, func(m ImageManifest) bool {
			return m.Platform == nil || m.Platform.OS != options.OS
		})
	}

	if options.Variant != "" {
		manifests = slices.DeleteFunc(manifests, func(m ImageManifest) bool {
			return m.Platform == nil || m.Platform.Variant != options.Variant
		})
	}

	if len(manifests) == 0 {
		return nil, nil
	}

	// Pick the first manifest
	manifest := manifests[0]

	// Resolve the reference to its digest
	ref.HasTag = false
	ref.Tag = ""
	ref.HasDigest = true
	ref.Digest = manifest.Digest

	// Fetch the manifest
	manifestBlob, err := c.GetManifestBlob(ctx, ref)
	if err != nil {
		return nil, err
	}
	defer manifestBlob.Close()

	var manifestContent struct {
		Config struct {
			Digest string `json:"digest"`
		} `json:"config"`
		Annotations Annotations `json:"annotations,omitempty"`
	}
	if err := json.NewDecoder(manifestBlob).Decode(&manifestContent); err != nil {
		return nil, err
	}

	annotations := manifestContent.Annotations

	// The blob was probably not a manifest blob but could still contain
	// annotations if it was a OCI manifest
	if manifestContent.Config.Digest == "" {
		return annotations, nil
	}

	// Get the config itself
	configBlob, err := c.GetBlob(ctx, ref, manifestContent.Config.Digest, true)
	if err != nil {
		return nil, err
	}
	defer configBlob.Close()

	// For now, only two formats are known to support annotations in config:
	// application/vnd.docker.container.image.v1+json
	// application/vnd.oci.image.config.v1+json
	// But the content types don't seem to be returned by all servers, so just try
	// to parse it anyways

	var configContent struct {
		Config struct {
			Labels map[string]string `json:"Labels"`
		} `json:"config"`
	}
	if err := json.NewDecoder(configBlob).Decode(&configContent); err != nil {
		return nil, err
	}

	return annotations.Merge(configContent.Config.Labels), nil
}

type GetTagsOptions struct {
	// Last is the name of the last tag of the previous page. Used for pagination.
	Last string
	// Count is the number of tags to return.
	// The server might not respect the choice.
	Count int
	// AllPages determines if the pagination is automatically handled to return
	// all available tags.
	AllPages bool
}

// GetTags retrieves available tags stored in a registry for a specific image.
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

func (c *Client) getTags(ctx context.Context, ref Reference, options *GetTagsOptions) ([]string, *url.URL, string, error) {
	ref = c.rewriteReference(ref)

	u, err := url.Parse(fmt.Sprintf("https://%s/v2/%s/tags/list", ref.Domain, ref.Path))
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

	if f := c.AuthFunc; f != nil {
		if err := f(req); err != nil {
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
	} else if err := assertStatusCode(res, http.StatusOK); err != nil {
		return nil, nil, "", err
	}

	var page struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}
	if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
		return nil, nil, "", err
	}

	tags := page.Tags

	return tags, req.URL, res.Header.Get("Link"), nil
}

// rewriteReference rewrites ref to handle caveats like Docker's registry not
// dealing with redirects.
func (c *Client) rewriteReference(ref Reference) Reference {
	// NOTE: It's rather unclear why we need to do this dance manually and why
	// docker.io simply doesn't just redirect us
	if ref.Domain == "docker.io" {
		ref.Domain = "registry-1.docker.io"
	}

	return ref
}
