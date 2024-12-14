package httputil

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveRequestURL(t *testing.T) {
	testCases := []struct {
		Name     string
		Request  *http.Request
		Expected *url.URL
	}{
		{
			Name: "localhost",
			Request: &http.Request{
				Host: "localhost:8080",
				URL: &url.URL{
					Path: "/api/v1/feed.rss",
				},
				Header: http.Header{},
			},
			Expected: &url.URL{
				Scheme: "http",
				Host:   "localhost:8080",
				Path:   "/api/v1/feed.rss",
			},
		},
		{
			Name: "proxied",
			Request: &http.Request{
				Host: "localhost:8080",
				URL: &url.URL{
					Path: "/api/v1/feed.rss",
				},
				Header: http.Header{
					"X-Forwarded-Host":  []string{"example.com"},
					"X-Forwarded-Proto": []string{"https"},
				},
			},
			Expected: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/api/v1/feed.rss",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual, err := ResolveRequestURL(testCase.Request)
			require.NoError(t, err)
			assert.Equal(t, testCase.Expected, actual)
		})
	}
}
