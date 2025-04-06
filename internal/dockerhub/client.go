package dockerhub

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
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/osv"
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

// GetVulnerabilities retrieves a Docker Scout vulnerability report for a
// repository and image digest.
// Returns nil if the results are inconclusive or an SBOM was not found.
func (c *Client) GetVulnerabilities(ctx context.Context, repo string, digest string) ([]osv.Vulnerability, error) {
	body, err := json.Marshal(map[string]any{
		"query": `query imagePackagesForImageCoords($v1:Context!,$v2:IpImagePackagesForImageCoordsQuery!){imagePackagesForImageCoords(context:$v1,query:$v2){digest,sbomState,imagePackages{packages{package{vulnerabilities{sourceId,description,url,cvss{score,severity}}}}}}}`,
		"variables": map[string]any{
			"v1": map[string]any{},
			"v2": map[string]any{
				"digest":          digest,
				"hostName":        "hub.docker.com",
				"repoName":        repo,
				"includeExcepted": true,
			},
		},
		"operationName": "imagePackagesForImageCoords",
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
			ImagePackagesForImageCoords struct {
				Digest        string `json:"digest"`
				SBOMState     string `json:"sbomState"`
				ImagePackages struct {
					Packages []struct {
						Package struct {
							Vulnerabilities []struct {
								SourceID    string `json:"sourceId"`
								Description string `json:"description"`
								URL         string `json:"url"`
								CVSS        struct {
									Score    *float32 `json:"score"`
									Severity string   `json:"severity"`
								} `json:"cvss"`
							} `json:"vulnerabilities"`
						} `json:"package"`
					} `json:"packages"`
				} `json:"imagePackages"`
			} `json:"imagePackagesForImageCoords"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Data.ImagePackagesForImageCoords.SBOMState != "INDEXED" {
		return nil, nil
	}

	vulnerabilities := make([]osv.Vulnerability, 0)
	seen := make(map[string]struct{})
	for _, pkg := range result.Data.ImagePackagesForImageCoords.ImagePackages.Packages {
		for _, vulnerability := range pkg.Package.Vulnerabilities {
			// Sanity check
			if vulnerability.SourceID == "" {
				continue
			}

			// Assume same data in all entries
			if _, ok := seen[vulnerability.SourceID]; ok {
				continue
			}

			severity := ""
			switch strings.ToLower(vulnerability.CVSS.Severity) {
			case "critical":
				severity = "CRITICAL"
			case "high":
				severity = "HIGH"
			case "moderate", "medium":
				severity = "MODERATE"
			case "low":
				severity = "LOW"
			}

			databaseSpecific := map[string]any{}
			if severity != "" {
				databaseSpecific["severity"] = severity
			}

			var severities []osv.Severity
			if vulnerability.CVSS.Score != nil {
				severities = []osv.Severity{
					{
						// NOTE: Assumed version
						Type:  "CVSS_V3",
						Score: fmt.Sprintf("%0.2f", float64(*vulnerability.CVSS.Score)),
					},
				}
			}

			vulnerabilities = append(vulnerabilities, osv.Vulnerability{
				ID:       vulnerability.SourceID,
				Modified: time.Now(), // TODO: Is there a better time to use?
				Summary:  vulnerability.Description,
				References: []osv.Reference{
					{
						Type: osv.ReferenceTypeWeb,
						URL:  vulnerability.URL,
					},
				},
				Severities:       severities,
				DatabaseSpecific: databaseSpecific,
			})
			seen[vulnerability.SourceID] = struct{}{}
		}
	}

	return vulnerabilities, nil
}
