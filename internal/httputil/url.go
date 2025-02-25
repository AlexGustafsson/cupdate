package httputil

import (
	"net/http"
	"net/url"
)

// ResolveRequestURL resolves the URL in request using common headers like
// X-Forwarded-Host and X-Forwarded-Proto.
func ResolveRequestURL(r *http.Request) (*url.URL, error) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	host := r.Host

	if header := r.Header.Get("X-Forwarded-Host"); header != "" {
		host = header
	}

	if header := r.Header.Get("X-Forwarded-Proto"); header != "" {
		if header == "http" || header == "https" {
			scheme = header
		}
	}

	base, err := url.Parse(scheme + "://" + host)
	if err != nil {
		return nil, err
	}

	return base.ResolveReference(r.URL), nil
}
