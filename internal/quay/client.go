package quay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/AlexGustafsson/cupdate/internal/osv"
)

type Client struct {
	Client httputil.Requester
}

// GetVulnerabilities retrieves a vulnerability scan of a manifest referenced by
// its digest.
func (c *Client) GetVulnerabilities(ctx context.Context, reference oci.Reference) ([]osv.Vulnerability, error) {
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
					NamespaceName   string
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

	vulnerabilities := make([]osv.Vulnerability, 0)
	for _, feature := range result.Data.Layer.Features {
		for _, vulnerability := range feature.Vulnerabilities {
			severity := ""
			switch vulnerability.Severity {
			case "Critical":
				severity = "CRITICAL"
			case "High":
				severity = "HIGH"
			case "Medium":
				severity = "MODERATE"
			case "Low":
				severity = "LOW"
			}

			var references []osv.Reference
			if vulnerability.Link != "" {
				links := strings.Split(vulnerability.Link, " ")
				if len(links) > 0 {
					references = make([]osv.Reference, 0)
					for _, link := range links {
						references = append(references, osv.Reference{
							Type: osv.ReferenceTypeWeb,
							URL:  link,
						})
					}
				}
			}

			databaseSpecific := make(map[string]any)
			if severity != "" {
				databaseSpecific["severity"] = severity
			}

			name := feature.Name
			if feature.NamespaceName != "" {
				name = feature.NamespaceName + "/" + name
			}

			affected := []osv.Affected{
				{
					Package: &osv.AffectedPackage{
						Name: name,
					},
				},
			}

			vulnerabilities = append(vulnerabilities, osv.Vulnerability{
				ID:               vulnerability.Name,
				Modified:         time.Now(), // TODO: Is there a better time to use?
				Summary:          vulnerability.Description,
				Affected:         affected,
				References:       references,
				DatabaseSpecific: databaseSpecific,
			})
		}
	}

	return vulnerabilities, nil
}
