package httputil

import (
	"fmt"
	"net/url"
	"strings"
)

type Link struct {
	URL    *url.URL
	Params map[string]string
}

// ParseLinkHeader parses a Link header.
//
// SEE: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Link.
func ParseLinkHeader(origin *url.URL, header string) ([]Link, error) {
	// NOTE: For now, we don't support commas in params
	parts := strings.Split(header, ", ")

	links := make([]Link, 0)
	for _, part := range parts {
		linkParts := strings.Split(part, "; ")

		if len(linkParts[0]) < 2 {
			return nil, fmt.Errorf("invalid link header")
		}

		if linkParts[0][0] != '<' || linkParts[0][len(linkParts[0])-1] != '>' {
			return nil, fmt.Errorf("invalid link header")
		}

		urlString := linkParts[0][1 : len(linkParts[0])-1]

		for _, char := range urlString {
			// Unescaped character
			if char > 255 {
				return nil, fmt.Errorf("invalid link header")
			}
		}

		url, err := origin.Parse(urlString)
		if err != nil {
			return nil, fmt.Errorf("invalid link header: %w", err)
		}

		params := make(map[string]string)
		for _, part := range linkParts[1:] {
			k, v, ok := strings.Cut(part, "=")
			if !ok {
				return nil, fmt.Errorf("invalid link header")
			}

			// Unquote
			if len(v) > 2 {
				if v[0] == '"' && v[len(v)-1] == '"' {
					v = v[1 : len(v)-1]
				}
			}

			params[k] = v
		}

		links = append(links, Link{
			URL:    url,
			Params: params,
		})
	}

	return links, nil
}
