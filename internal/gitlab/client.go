package gitlab

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
)

var readmePathRegexp = regexp.MustCompile(`href="(.*?/blob/.*?)"`)

type Client struct {
	Client *httputil.Client
}

func (c *Client) GetRepositoryDescription(ctx context.Context, fullPath string) (string, error) {
	payload, err := json.Marshal(map[string]any{
		"operationName": "getProject",
		"variables": map[string]any{
			"fullPath": fullPath,
		},
		"query": `query getProject($fullPath: ID!) {
  project(fullPath: $fullPath) {
    description
  }
}`,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://gitlab.com/api/graphql", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// TODO: The cache doesn't understand graphql, so we can't cache this request
	res, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %s", res.Status)
	}

	var result struct {
		Errors []any `json:"error"`
		Data   struct {
			Project struct {
				Description string `json:"description"`
			} `json:"project"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	} else if len(result.Errors) > 0 {
		return "", fmt.Errorf("got unexpected graphql error: %v", result.Errors)
	}

	return result.Data.Project.Description, nil
}

func (c *Client) GetRepositoryREADMEBlob(ctx context.Context, fullPath string) (*Blob, error) {
	// https://gitlab.com/arm-research/smarter/smarter-device-manager/-/blob/master/README.md?format=json&viewer=rich
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://gitlab.com/"+fullPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Client.DoCached(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	html, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// The API itself doesn't expose the README. The UI has some "startup" calls
	// defined that include the README, whichever file it is.

	readmeHref := readmePathRegexp.FindSubmatch(html)
	if readmeHref == nil {
		return nil, nil
	}

	return c.GetBlob(ctx, string(readmeHref[1]), true)
}

func (c *Client) GetBlob(ctx context.Context, href string, includeRaw bool) (*Blob, error) {
	u := url.URL{
		Scheme:   "https",
		Host:     "gitlab.com",
		Path:     href,
		RawQuery: "format=json",
	}

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

	var blob Blob
	if err := json.NewDecoder(res.Body).Decode(&blob); err != nil {
		return nil, err
	}

	if includeRaw {
		u := url.URL{
			Scheme: "https",
			Host:   "gitlab.com",
			Path:   blob.RawPath,
		}

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

		raw, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		blob.Raw = raw
	}

	return &blob, nil
}

type Blob struct {
	ID            string `json:"id"`
	Path          string `json:"path"`
	Name          string `json:"name"`
	Extension     string `json:"extension"`
	Size          int    `json:"size"`
	MimeType      string `json:"mime_type"`
	Binary        bool   `json:"binary"`
	RawPath       string `json:"raw_path"`
	BlamePath     string `json:"blame_path"`
	CommitsPath   string `json:"commits_path"`
	TreePath      string `json:"tree_path"`
	Permalink     string `json:"permalink"`
	LastCommitSHA string `json:"last_commit_sha"`
	HTML          string `json:"html"`
	Raw           []byte `json:"-"`

	// ... some available fields are ignored
}
