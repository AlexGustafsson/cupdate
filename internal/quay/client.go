package quay

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"slices"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
)

var _ oci.Authorizer = (*Client)(nil)

type Client struct {
	Client *httputil.Client
}

func (c *Client) AuthorizeOCIRequest(ctx context.Context, image oci.Reference, req *http.Request) error {
	// NOOP
	return nil
}

func (c *Client) GetTags(ctx context.Context, image oci.Reference) ([]string, error) {
	if !image.HasTag {
		return nil, nil
	}

	// There's not going to be any latest version
	if image.Tag == "latest" {
		return nil, nil
	}

	// NOTE: The quay UI themselves go through all pages immediately
	u, err := url.Parse(fmt.Sprintf("https://quay.io/api/v1/repository/%s/tag/", image.Path))
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("limit", "100")
	query.Set("page", "1")
	query.Set("onlyActiveTags", "true")
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.DoCached(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	var result TagPage
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	tags := make(map[string]struct{})
	for _, tag := range result.Tags {
		tags[tag.Name] = struct{}{}
	}

	return slices.Collect(maps.Keys(tags)), nil
}

type TagPage struct {
	Page          int   `json:"page"`
	HasAdditional bool  `json:"has_additional"`
	Tags          []Tag `json:"tags"`
}

type Tag struct {
	Name           string `json:"name"`
	Reversion      bool   `json:"reversion"`
	StartTimestamp int    `json:"start_ts"`
	ManifestDigest string `json:"manifest_digest"`
	IsManifestList bool   `json:"is_manifest_list"`
	Size           int    `json:"size"`
	LastModified   string `json:"last_modified"`
}
