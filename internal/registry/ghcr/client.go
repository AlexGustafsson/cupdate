package ghcr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

type Client struct {
	Client *httputil.Client
}

func (c *Client) GetRegistryToken(ctx context.Context, image oci.Reference) (string, error) {
	// TODO: Registries expose the realm and scheme via Www-Authenticate if 403
	// is given
	u, err := url.Parse("https://ghcr.io/token?service=ghcr.io")
	if err != nil {
		return "", err
	}

	query := u.Query()
	query.Set("scope", fmt.Sprintf("repository:%s:pull", image.Path))
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %s", res.Status)
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Token, nil
}

func (c *Client) GetManifests(ctx context.Context, image oci.Reference) ([]oci.Manifest, error) {
	token, err := c.GetRegistryToken(ctx, image)
	if err != nil {
		return nil, err
	}

	ociClient := &oci.Client{
		Client: c.Client,
		// TODO: Cache token
		Authorizer: oci.AuthorizerToken(token),
	}

	return ociClient.GetManifests(ctx, image)
}
