package quay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
)

type Client struct {
	Client httputil.Requester
}

// GetVulnerabilities retrieves a vulnerability scan of a manifest referenced by
// its digest.
func (c *Client) GetVulnerabilities(ctx context.Context, reference oci.Reference) ([]Vulnerability, error) {
	if reference.Digest == "" {
		return nil, fmt.Errorf("reference has no digest")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://quay.io/api/v1/repository/%s/manifest/%s/security?vulnerabilities=true", reference.Path, reference.Digest), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	res, err := c.Client.DoCached(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return nil, err
	}

	var result struct {
		Status string
		Data   *struct {
			Layer struct {
				Features []struct {
					Name            string
					AddedBy         string
					Version         string
					Vulnerabilities []struct {
						Severity    string
						Link        string
						Description string
						Name        string
					}
				}
			}
		}
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Status != "scanned" || result.Data == nil {
		return nil, nil
	}

	vulnerabilities := make([]Vulnerability, 0)
	for _, feature := range result.Data.Layer.Features {
		for _, vulnerability := range feature.Vulnerabilities {
			severity := VulnerabilitySeverityUnspecified
			switch vulnerability.Severity {
			case "Critical":
				severity = VulnerabilitySeverityCritical
			case "High":
				severity = VulnerabilitySeverityHigh
			case "Medium":
				severity = VulnerabilitySeverityMedium
			case "Low":
				severity = VulnerabilitySeverityLow
			}

			vulnerabilities = append(vulnerabilities, Vulnerability{
				Name:           vulnerability.Name,
				Description:    vulnerability.Description,
				Links:          strings.Split(vulnerability.Link, " "),
				FeatureName:    feature.Name,
				FeatureVersion: feature.Version,
				Layer:          feature.AddedBy,
				Severity:       severity,
			})
		}
	}

	return vulnerabilities, nil
}
