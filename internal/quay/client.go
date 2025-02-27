package quay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
)

type Client struct {
	Client httputil.Requester
}

// GetScan retrieves a scan of a manifest referenced by its digest.
func (c *Client) GetScan(ctx context.Context, reference oci.Reference) (*Scan, error) {
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

	var scan Scan
	if err := json.NewDecoder(res.Body).Decode(&scan); err != nil {
		return nil, err
	}

	return &scan, nil
}
