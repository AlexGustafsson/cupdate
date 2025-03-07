package scorecard

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
)

type Client struct {
	Client httputil.Requester
}

// GetScorecard retrieves a scorecard report for a repository.
// A repository is formated like so: <hostname>/<path>.
// See also [RepositoryIsSupported].
//
// Examples:
//
//	c.GetScorecard(ctx, "github.com/homeassistant/core")
//	c.GetScorecard(ctx, "gitlab.com/baserow/baserow")
func (c *Client) GetScorecard(ctx context.Context, repository string) (*Scorecard, error) {
	if !RepositoryIsSupported(repository) {
		return nil, fmt.Errorf("unsupported repository")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.scorecard.dev/projects/"+repository, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.DoCached(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if err := httputil.AssertStatusCode(res, http.StatusOK); err != nil {
		return nil, err
	}

	var scorecard Scorecard
	if err := json.NewDecoder(res.Body).Decode(&scorecard); err != nil {
		return nil, err
	}

	return &scorecard, nil
}
