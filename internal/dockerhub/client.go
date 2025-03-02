package dockerhub

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
)

type Client struct {
	Client httputil.Requester
}

func (c *Client) GetRepository(ctx context.Context, image oci.Reference) (*Repository, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://hub.docker.com/v2/repositories/"+url.PathEscape(image.Path), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.DoCached(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return nil, err
	}

	var result Repository
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOrganizationOrUser retrieves information about a Docker Hub user or
// organization by name.
func (c *Client) GetOrganizationOrUser(ctx context.Context, organizationOrUser string) (*Entity, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://hub.docker.com/v2/orgs/"+url.PathEscape(organizationOrUser), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.DoCached(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return nil, err
	}

	var result Entity
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetVulnerabilityReport retrieves a Docker Scout vulnerability report for a
// repository and image digest.
func (c *Client) GetVulnerabilityReport(ctx context.Context, repo string, digest string) (*VulnerabilityReport, error) {
	body, err := json.Marshal(map[string]any{
		"query": "query imageSummariesByDigest($v1:Context!,$v2:[String!]!,$v3:ScRepositoryInput){imageSummariesByDigest(context:$v1,digests:$v2,repository:$v3){digest,sbomState,vulnerabilityReport{critical,high,medium,low,unspecified,total}}}",
		"variables": map[string]any{
			"v1": map[string]any{},
			"v2": []string{
				digest,
			},
			"v3": map[string]any{
				"hostName": "hub.docker.com",
				"repoName": repo,
			},
		},
		"operationName": "imageSummariesByDigest",
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.dso.docker.com/v1/graphql", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// TODO: The cache doesn't understand graphql, so we can't cache this request
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			ImageSummariesByDigest []struct {
				Digest              string               `json:"digest"`
				SBOMStatae          string               `json:"sbomState"`
				VulnerabilityReport *VulnerabilityReport `json:"vulnerabilityReport"`
			} `json:"imageSummariesByDigest"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Data.ImageSummariesByDigest) != 1 {
		return nil, nil
	}

	return result.Data.ImageSummariesByDigest[0].VulnerabilityReport, nil
}
