package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/registry"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
)

type Client struct {
	Client *httputil.Client
}

func (c *Client) GetRegistryToken(ctx context.Context, image oci.Reference) (string, error) {
	// TODO: Registries expose the realm and scheme via Www-Authenticate if 403
	// is given
	u, err := url.Parse("https://auth.docker.io/token?service=registry.docker.io")
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
		Token     string    `json:"token"`
		ExpiresIn int       `json:"expires_in"`
		IssuedAt  time.Time `json:"issued_at"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Token, nil
}

func (c *Client) Authorize(ctx context.Context, image oci.Reference, req *http.Request) error {
	token, err := c.GetRegistryToken(ctx, image)
	if err != nil {
		return err
	}

	return oci.AuthorizerToken(token).Authorize(ctx, image, req)
}

func (c *Client) GetLatestVersion(ctx context.Context, image oci.Reference) (*registry.Image, error) {
	if !image.HasTag {
		return nil, nil
	}

	// There's not going to be any latest version
	if image.Tag == "latest" {
		return nil, nil
	}

	currentVersion, err := oci.ParseVersion(image.Tag)
	if err != nil {
		return nil, fmt.Errorf("unsupported version: %w", err)
	} else if currentVersion == nil {
		return nil, fmt.Errorf("unsupported version")
	}

	u, err := url.Parse(fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags", image.Path))
	if err != nil {
		return nil, err
	}

	query := u.Query()
	query.Set("page_size", "25")
	query.Set("page", "1")
	query.Set("ordering", "last_updated")
	query.Set("name", "")
	u.RawQuery = query.Encode()

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

	var result Page[Tag]
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	// TODO:
	// As we've sorted versions in released time, let's assume the first version
	// that is higher than ours, is the latest version. Might not be true if the
	// current version is 1.0.0, there have been a lot of nightlies or other types
	// of tags, so that the page contains only fix 1.0.1, but in reality 2.0.0 was
	// released a while ago and would be on the next page, would we be greedy.
	// Look at any large image with LTS, such as postgres, node, calico/node.
	// We could try to take care of "common" edge cases such as calico/node's
	// v3.29.0-arm64, v3.29.0-amd64 etc where v3.29.0 is the tagged version and
	// the suffix is only the platform. If we were to trim those before comparing,
	// we would potentially find a (probable) latest version faster, but we could
	// also end up including "different" tags, if for example the user was using
	// a tag with an actual suffix, then the suffix was important. Could probably
	// be handled by only checking if there is a suffix in the user's image. Over
	// time, if the project is very active, I guess we should hit a correct
	// version eventually, if cupdate checks often enough
	for _, tag := range result.Results {
		if tag.Name == "" {
			continue
		}

		newVersion, err := oci.ParseVersion(tag.Name)
		if err != nil || newVersion == nil {
			continue
		}

		if newVersion.IsCompatible(currentVersion) && newVersion.Compare(currentVersion) >= 0 {
			image.Tag = tag.Name

			return &registry.Image{
				Name:      image,
				Published: tag.TagLastPushed,
				Digest:    tag.Digest,
			}, nil
		}
	}

	return nil, nil
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
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	var result Repository
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

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
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	var result Entity
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetVulnerabilityReport(ctx context.Context, repo string, digest string) (*VulnerabilityReport, error) {
	body, err := json.Marshal(map[string]any{
		"query": "query imageSummariesByDigest($v1:Context!,$v2:[String!]!,$v3:ScRepositoryInput){imageSummariesByDigest(context:$v1,digests:$v2,repository:$v3){digest,sbomState,vulnerabilityReport{critical,high,medium,low,unspecified,total}}}",
		"variables": map[string]any{
			"v1": map[string]any{},
			"v2": []string{
				digest,
			},
			"v3": map[string]any{
				"hostName": "hub.docker.com",
				"repoName": repo,
			},
		},
		"operationName": "imageSummariesByDigest",
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.dso.docker.com/v1/graphql", bytes.NewReader(body))
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
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	var result struct {
		Data struct {
			ImageSummariesByDigest []struct {
				Digest              string               `json:"digest"`
				SBOMStatae          string               `json:"sbomState"`
				VulnerabilityReport *VulnerabilityReport `json:"vulnerabilityReport"`
			} `json:"imageSummariesByDigest"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Data.ImageSummariesByDigest) != 1 {
		return nil, nil
	}

	return result.Data.ImageSummariesByDigest[0].VulnerabilityReport, nil
}

type Page[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

type Tag struct {
	Creator             int       `json:"creator"`
	ID                  int       `json:"id"`
	LastUpdated         time.Time `json:"last_updated"`
	LastUpdater         int       `json:"last_updater"`
	LastUpdaterUsername string    `json:"last_updater_username"`
	Name                string    `json:"name"`
	Repository          int       `json:"repository"`
	FullSize            int       `json:"full_size"`
	V2                  bool      `json:"v2"`
	TagStatus           string    `json:"tag_status"`
	TagLastPulled       time.Time `json:"tag_last_pulled"`
	TagLastPushed       time.Time `json:"tag_last_pushed"`
	MediaType           string    `json:"media_type"`
	ContentType         string    `json:"content_type"`
	Digest              string    `json:"digest"`
	Images              []Image   `json:"images"`
}

type Image struct {
	Architecture string    `json:"architecture"`
	Features     string    `json:"features"`
	Variant      *string   `json:"variant"`
	Digest       string    `json:"digest"`
	OS           string    `json:"os"`
	OSFeatures   string    `json:"os_features"`
	OSVersion    *string   `json:"os_version"`
	Size         int       `json:"size"`
	Status       string    `json:"status"`
	LastPulled   time.Time `json:"last_pulled"`
	LastPushed   time.Time `json:"last_pushed"`
}

type Repository struct {
	User              string          `json:"user"`
	Name              string          `json:"name"`
	Namespace         string          `json:"namespace"`
	Type              string          `json:"repository_type"`
	Status            int             `json:"status"`
	StatusDescription string          `json:"status_description"`
	Description       string          `json:"description"`
	IsPrivate         bool            `json:"is_private"`
	IsAutomated       bool            `json:"is_automated"`
	StarCount         int             `json:"star_count"`
	PullCount         int             `json:"pull_count"`
	LastUpdated       time.Time       `json:"last_updated"`
	DateRegistered    time.Time       `json:"date_registered"`
	CollaboratorCount int             `json:"collaborator_count"`
	Affiliation       json.RawMessage `json:"affiliation"` // Unknown
	HubUser           string          `json:"hub_user"`
	HasStarred        bool            `json:"has_starred"`
	FullDescription   string          `json:"full_description"`
	Permissions       struct {
		Read  bool `json:"read"`
		Write bool `json:"write"`
		Admin bool `json:"admin"`
	} `json:"permissions"`
	MediaTypes   []string `json:"media_types"`
	ContentTypes []string `json:"content_types"`
	Categories   []struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"categories"`
	ImmutableTags      bool   `json:"immutable_tags"`
	ImmutableTagsRules string `json:"immutable_tags_rules"`
}

type Entity struct {
	ID               string    `json:"id"`
	UUID             string    `json:"uuid,omitempty"`
	OrganizationName string    `json:"orgname"`
	Username         string    `json:"username,omitempty"`
	FullName         string    `json:"full_name"`
	Location         string    `json:"location"`
	Company          string    `json:"company"`
	ProfileURL       string    `json:"profile_url"`
	DateJoined       time.Time `json:"date_joined"`
	GravatarURL      string    `json:"gravatar_url"`
	GravatarEmail    string    `json:"gravatar_email"`
	Type             string    `json:"type"`
	Badge            string    `json:"badge,omitempty"`
}

type VulnerabilityReport struct {
	Critical    int `json:"critical"`
	High        int `json:"high"`
	Medium      int `json:"medium"`
	Low         int `json:"low"`
	Unspecified int `json:"unspecified"`
	Total       int `json:"total"`
}
