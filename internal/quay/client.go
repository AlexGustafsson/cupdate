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
	Client *httputil.Client
}

func (c *Client) GetScan(ctx context.Context, reference oci.Reference, digest string) (*Scan, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://quay.io/api/v1/repository/%s/manifest/%s/security?vulnerabilities=true", reference.Path, digest), nil)
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
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	var scan Scan
	if err := json.NewDecoder(res.Body).Decode(&scan); err != nil {
		return nil, err
	}

	return &scan, nil
}
