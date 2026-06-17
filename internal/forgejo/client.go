package forgejo

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/AlexGustafsson/cupdate/internal/htmlutil"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
)

type Client struct {
	Client httputil.Requester
}

// GetREADME retrieves the README at url - a repository's main page.
func (c *Client) GetREADME(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return "", err
	}

	document, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var document xml.Name
	if err := decoder.Decode(&document); err != nil {
		return "", err
	}

	fmt.Printf("%+v\n", document)

	readme := ""

	cleaned, err := htmlutil.ResolveReferences(readme, req.URL)
	if err == nil {
		readme = cleaned
	}

	return readme, nil
}
