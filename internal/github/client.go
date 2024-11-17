package github

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/registry/ghcr"
	"github.com/AlexGustafsson/cupdate/internal/registry/oci"
	"golang.org/x/net/html"
)

type Client struct {
	// Endpoint is the GitHub endpoint. Useful for using an enterprise instance,
	// for example. Defaults to "https://github.com".
	Endpoint string
	Client   *httputil.Client
}

func (c *Client) GetRelease(ctx context.Context, owner string, repository string, tag string) (*Release, error) {
	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = "https://github.com"
	}

	url := fmt.Sprintf("%s/%s/%s/releases/tag/%s", endpoint, url.PathEscape(owner), url.PathEscape(repository), tag)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
	} else if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status: %s - %d", res.Status, res.StatusCode)
	}

	release, err := parseRelease(res.Body)
	if err != nil {
		return nil, err
	}

	release.URL = url
	release.Owner = owner
	release.Repository = repository
	release.Tag = tag

	return release, nil
}

func (c *Client) GetDescription(ctx context.Context, owner string, repository string) (string, error) {
	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = "https://github.com"
	}

	url := fmt.Sprintf("%s/%s/%s", endpoint, url.PathEscape(owner), url.PathEscape(repository))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	res, err := c.Client.DoCached(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return "", nil
	} else if res.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status: %s - %d", res.Status, res.StatusCode)
	}

	return parseAbout(res.Body)
}

func (c *Client) GetPackage(ctx context.Context, reference oci.Reference) (*Package, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ghcr.PackagePath(reference), nil)
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
	} else if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status: %s - %d", res.Status, res.StatusCode)
	}

	owner, _, _ := strings.Cut(reference.Path, "/")
	return parsePackage(res.Body, owner)
}

func parsePackage(r io.Reader, owner string) (*Package, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	repositoryRef := match(node, func(node *html.Node) bool {
		return node.Data == "a" && strings.HasPrefix(strings.ToLower(attr(node, "href")), fmt.Sprintf("https://github.com/%s/", strings.ToLower(owner)))
	})
	if repositoryRef == nil {
		return nil, nil
	}
	_, _, repository, _, ok := ParseURL(attr(repositoryRef, "href"))
	if !ok {
		return nil, nil
	}

	hrefs := make([]string, 0)
	match(node, func(node *html.Node) bool {
		if node.Data == "a" && strings.HasPrefix(attr(node, "href"), fmt.Sprintf("/users/%s/packages/container/", owner)) {
			hrefs = append(hrefs, attr(node, "href"))
		}
		return false
	})
	if len(hrefs) == 0 {
		return nil, nil
	}

	// Map image tags to image ids
	tagToId := make(map[string]string)
	for _, href := range hrefs {
		path, args, ok := strings.Cut(href, "?")
		if !ok {
			continue
		}

		query, err := url.ParseQuery(args)
		if err != nil {
			continue
		}

		tag := query.Get("tag")
		if tag == "" {
			continue
		}

		i := strings.LastIndex(path, "/")
		if i < 0 {
			continue
		}

		id := path[i+1:]

		tagToId[tag] = id
	}

	tags := make([]PackageTag, 0)
	latestId, hasLatest := tagToId["latest"]
	for tag, id := range tagToId {
		isLatest := hasLatest && id == latestId
		tags = append(tags, PackageTag{Name: tag, Latest: isLatest})
	}

	return &Package{
		Owner:      owner,
		Repository: repository,
		Tags:       tags,
	}, nil
}

func parseAbout(r io.Reader) (string, error) {
	node, err := html.Parse(r)
	if err != nil {
		return "", err
	}

	// <div class="Layout-sidebar">
	sidebar := match(node, func(node *html.Node) bool {
		return node.Data == "div" &&
			attr(node, "class") == "Layout-sidebar"
	})
	if sidebar == nil {
		return "", nil
	}

	// <h2>About</h2>
	about := match(sidebar, func(node *html.Node) bool {
		return node.Data == "h2" && node.FirstChild != nil && node.FirstChild.Data == "About"
	})
	if about == nil {
		return "", nil
	}

	description := match(about, func(node *html.Node) bool {
		return node.Data == "p"
	})
	if node == nil {
		return "", nil
	}

	var buffer bytes.Buffer
	if err := html.Render(&buffer, description.FirstChild); err != nil {
		return "", err
	}

	return strings.TrimSpace(buffer.String()), nil
}

func parseRelease(r io.Reader) (*Release, error) {
	html, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	title, err := parseTitle(html)
	if err != nil {
		return nil, err
	}

	// <relative-time class="no-wrap" prefix="" datetime="2024-09-15T05:17:58Z" title="Sep 15, 2024 at 7:17 AM GMT+2">
	releaseTime := parseReleaseTime(html)

	releaseNotes, err := parseReleaseNotes(html)
	if err != nil {
		return nil, err
	}

	return &Release{
		Title:       title,
		Released:    releaseTime,
		Description: releaseNotes,
	}, nil
}

func parseTitle(node *html.Node) (string, error) {
	box := match(node, func(node *html.Node) bool {
		return node.Data == "div" && attr(node, "class") == "Box"
	})
	if box == nil {
		return "", nil
	}

	node = match(box, func(node *html.Node) bool {
		return node.Data == "h1"
	})
	if node == nil {
		return "", nil
	}

	var buffer bytes.Buffer
	if err := html.Render(&buffer, node.FirstChild); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// parseReleaseTime extracts the release time by performing a depth-first search
// for the first <relative-time> tag, which specifies the release time.
func parseReleaseTime(node *html.Node) time.Time {
	node = match(node, func(node *html.Node) bool {
		return node.Data == "relative-time"
	})

	if node != nil {
		datetime := attr(node, "datetime")
		if datetime != "" {
			time, err := time.Parse(time.RFC3339, datetime)
			if err == nil {
				return time
			}
		}
	}

	return time.Time{}
}

func parseReleaseNotes(node *html.Node) (string, error) {
	node = match(node, func(node *html.Node) bool {
		if node.Data != "div" {
			return false
		}

		return strings.Contains(attr(node, "class"), "markdown-body")
	})

	if node != nil {
		var buffer bytes.Buffer

		// The node we matched is a container of the data. We don't need that
		// container, so render each child on their own instead
		child := node.FirstChild
		for child != nil {
			if err := html.Render(&buffer, child); err != nil {
				return "", err
			}
			child = child.NextSibling
			if _, err := buffer.WriteRune('\n'); err != nil {
				return "", err
			}
		}

		return buffer.String(), nil
	}

	return "", nil
}

// match matches a node by performing a depth-first search.
func match(node *html.Node, matchFunc func(node *html.Node) bool) *html.Node {
	if matchFunc(node) {
		return node
	}

	// Check subtree
	if node.FirstChild != nil {
		result := match(node.FirstChild, matchFunc)
		if result != nil {
			return result
		}
	}

	// Go to next child
	if node.NextSibling != nil {
		result := match(node.NextSibling, matchFunc)
		if result != nil {
			return result
		}
	}

	return nil
}

func attr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}

	return ""
}
