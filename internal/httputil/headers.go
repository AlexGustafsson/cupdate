package httputil

import (
	"fmt"
	"net/url"
	"strconv"
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

// ParseWWWAuthenticateHeader parses a Www-Authenticate header.
//
// SEE: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/WWW-Authenticate.
func ParseWWWAuthenticateHeader(header string) (string, map[string]string, error) {
	var scheme string
	params := make(map[string]string)

	state := "scheme"
	paramKey := ""
	paramValue := ""
	gotParamDelimiter := false
	for i, c := range header {
		isAlpha := c >= 'a' && c <= 'z'
		isNumeric := c >= '0' && c <= '0'
		isAlphaNumeric := isAlpha || isNumeric || c == '-'
		isEnd := i == len(header)-1

		switch state {
		case "scheme":
			if c == ' ' {
				if isEnd {
					return "", nil, fmt.Errorf("httputil: invalid Www-Authenticate trailing whitespace")
				} else {
					state = "paramKey"
				}
			} else {
				scheme += string(c)
				if isEnd {
					state = "end"
				}
			}
		case "paramKey":
			// Consume optional whitespace after params delimiter
			if gotParamDelimiter && c == ' ' {
				continue
			} else {
				gotParamDelimiter = false
			}

			if c == '=' {
				state = "paramValue"
			} else if paramKey == "" && isAlpha {
				paramKey += string(c)
			} else if paramKey != "" && isAlphaNumeric {
				paramKey += string(c)
			} else {
				return "", nil, fmt.Errorf("httputil: invalid Www-Authenticate header param key")
			}
		case "paramValue":
			if paramValue == "" && c == '"' {
				// OK
			} else if paramValue != "" && c == '"' {
				params[paramKey] = paramValue
				paramKey = ""
				paramValue = ""
				if isEnd {
					state = "end"
				} else {
					state = "paramDelimiter"
				}
			} else {
				paramValue += string(c)
			}
		case "paramDelimiter":
			if c == ',' {
				gotParamDelimiter = true
				state = "paramKey"
			}
		default:
			return "", nil, fmt.Errorf("httputil: invalid Www-Authenticate state")
		}
	}

	if state != "end" {
		return "", nil, fmt.Errorf("httputil: invalid Www-Authenticate state")
	}

	return scheme, params, nil
}

// Accepts returns the most highly desired mime type supported by the server.
// Ignores parameters like charsets.
func Accepts(header string, mimeTypes ...string) string {
	matched := ""
	matchedWeight := 0.0

	entries := strings.Split(strings.TrimSpace(header), ",")
	for _, entry := range entries {
		directives := strings.Split(strings.TrimSpace(entry), ";")

		pattern := directives[0]
		wantedLeft, wantedRight, ok := strings.Cut(pattern, "/")
		if !ok {
			// Malformed
			continue
		}

		weight := 1.0
		for _, directive := range directives[1:] {
			k, v, ok := strings.Cut(strings.TrimSpace(directive), "=")
			if !ok {
				continue
			}

			if k == "q" {
				var err error
				weight, err = strconv.ParseFloat(v, 32)
				if err != nil {
					continue
				}
				break
			}
		}

		if weight <= matchedWeight {
			continue
		}

		for _, mimeType := range mimeTypes {
			actualLeft, actualRight, ok := strings.Cut(mimeType, "/")
			if !ok {
				// Malformed
				continue
			}

			matches := (wantedLeft == "*" || actualLeft == wantedLeft) && (wantedRight == "*" || actualRight == wantedRight)
			if !matches {
				continue
			}

			matched = mimeType
			matchedWeight = weight
			break
		}
	}

	return matched
}
