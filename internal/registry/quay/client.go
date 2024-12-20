package quay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/registry"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"github.com/AlexGustafsson/cupdate/internal/semver"
)

type Client struct {
	Client *httputil.Client
}

func (c *Client) GetLatestVersion(ctx context.Context, image oci.Reference) (*registry.Image, error) {
	if !image.HasTag {
		return nil, nil
	}

	// There's not going to be any latest version
	if image.Tag == "latest" {
		return nil, nil
	}

	currentVersion, err := semver.ParseVersion(image.Tag)
	if err != nil {
		return nil, fmt.Errorf("unsupported version: %w", err)
	} else if currentVersion == nil {
		return nil, fmt.Errorf("unsupported version")
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

	// TODO:
	// As we've sorted versions in released time, let's assume the first version
	// that is higher than ours, is the latest version. Might not be true if the
	// current version is 1.0.0, there have been a lot of nightlies or other types
	// of tags, so that the page contains only fix 1.0.1, but in reality 2.0.0 was
	// released a while ago and would be on the next page, would we be greedy.
	// Look at any large image with LTS, such as postgres, node.
	for _, tag := range result.Tags {
		if tag.Name == "" {
			continue
		}

		newVersion, err := semver.ParseVersion(tag.Name)
		if err != nil || newVersion == nil {
			continue
		}

		if currentVersion.Prerelease == "" && newVersion.Prerelease != "" {
			continue
		}

		if newVersion.IsCompatible(currentVersion) && newVersion.Compare(currentVersion) >= 0 {
			image.Tag = tag.Name

			time, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", tag.LastModified)
			return &registry.Image{
				Name:      image,
				Published: time,
				Digest:    tag.ManifestDigest,
			}, nil
		}
	}

	return nil, nil
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
