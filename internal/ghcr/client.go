package ghcr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
)

type Client struct {
	Client        *httputil.Client
	TokenAuthFunc func(*http.Request) error
}

// GetRegistryToken returns a token for use with GHCR with pull permissions on
// the specified repository.
func (c *Client) GetRegistryToken(ctx context.Context, repository string) (string, error) {
	// TODO: Registries expose the realm and scheme via Www-Authenticate if 403
	// is given
	u, err := url.Parse("https://ghcr.io/token?service=ghcr.io")
	if err != nil {
		return "", err
	}

	query := u.Query()
	query.Set("scope", fmt.Sprintf("repository:%s:pull", repository))
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}

	if f := c.TokenAuthFunc; f != nil {
		if err := f(req); err != nil {
			return "", err
		}
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}

	if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return "", err
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Token, nil
}

// HandleAuth authenticates a request to the GHCR registry.
func (c *Client) HandleAuth(r *http.Request) error {
	name := oci.NameFromAPI(r.URL.Path)
	// lscr.io is a pseudo-registry that forwards to one of multiple backends,
	// among them ghcr.io
	if (r.Host != "ghcr.io" && r.Host != "lscr.io") || name == "" {
		return nil
	}

	token, err := c.GetRegistryToken(r.Context(), name)
	if err != nil {
		return err
	}

	r.Header.Set("Authorization", "Bearer "+token)

	return nil
}
