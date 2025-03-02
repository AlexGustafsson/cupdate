package httputil

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLinkHeader(t *testing.T) {
	testCases := []struct {
		Origin   *url.URL
		Header   string
		Expected []Link
		Error    bool
	}{
		{
			Origin: mustParseURL(t, "https://example.com"),
			Header: `</uri-reference>; param1=value1; param2="value2"`,
			Expected: []Link{
				{
					URL: mustParseURL(t, "https://example.com/uri-reference"),
					Params: map[string]string{
						"param1": "value1",
						"param2": "value2",
					},
				},
			},
		},
		{
			Origin: mustParseURL(t, "https://example.com"),
			Header: `</uri-reference>`,
			Expected: []Link{
				{
					URL:    mustParseURL(t, "https://example.com/uri-reference"),
					Params: map[string]string{},
				},
			},
		},
		{
			Origin: mustParseURL(t, "https://example.com"),
			Header: `<https://example.com>; rel="preconnect"`,
			Expected: []Link{
				{
					URL: mustParseURL(t, "https://example.com"),
					Params: map[string]string{
						"rel": "preconnect",
					},
				},
			},
		},
		{
			Origin: mustParseURL(t, "https://example.com"),
			Header: `<https://example.com/%E8%8B%97%E6%9D%A1>; rel="preconnect"`,
			Expected: []Link{
				{
					URL: mustParseURL(t, "https://example.com/%E8%8B%97%E6%9D%A1"),
					Params: map[string]string{
						"rel": "preconnect",
					},
				},
			},
		},
		{
			Origin: mustParseURL(t, "https://example.com"),
			Header: `</style.css>; rel=preload; as=style; fetchpriority="high"`,
			Expected: []Link{
				{
					URL: mustParseURL(t, "https://example.com/style.css"),
					Params: map[string]string{
						"rel":           "preload",
						"as":            "style",
						"fetchpriority": "high",
					},
				},
			},
		},
		{
			Origin: mustParseURL(t, "https://example.com"),
			Header: `<https://one.example.com>; rel="preconnect", <https://two.example.com>; rel="preconnect", <https://three.example.com>; rel="preconnect"`,
			Expected: []Link{
				{
					URL: mustParseURL(t, "https://one.example.com"),
					Params: map[string]string{
						"rel": "preconnect",
					},
				},
				{
					URL: mustParseURL(t, "https://two.example.com"),
					Params: map[string]string{
						"rel": "preconnect",
					},
				},
				{
					URL: mustParseURL(t, "https://three.example.com"),
					Params: map[string]string{
						"rel": "preconnect",
					},
				},
			},
		},
		{
			Origin: mustParseURL(t, "https://example.com"),
			Header: `<https://api.example.com/issues?page=2>; rel="prev", <https://api.example.com/issues?page=4>; rel="next", <https://api.example.com/issues?page=10>; rel="last", <https://api.example.com/issues?page=1>; rel="first"`,
			Expected: []Link{
				{
					URL: mustParseURL(t, "https://api.example.com/issues?page=2"),
					Params: map[string]string{
						"rel": "prev",
					},
				},
				{
					URL: mustParseURL(t, "https://api.example.com/issues?page=4"),
					Params: map[string]string{
						"rel": "next",
					},
				},
				{
					URL: mustParseURL(t, "https://api.example.com/issues?page=10"),
					Params: map[string]string{
						"rel": "last",
					},
				},
				{
					URL: mustParseURL(t, "https://api.example.com/issues?page=1"),
					Params: map[string]string{
						"rel": "first",
					},
				},
			},
		},
		{
			Origin: mustParseURL(t, "https://example.com"),
			Header: `https://bad.example; rel="preconnect"`,
			Error:  true,
		},
		{
			Origin: mustParseURL(t, "https://example.com"),
			Header: `<https://example.com/苗条>; rel="preconnect"`,
			Error:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Header, func(t *testing.T) {
			actual, err := ParseLinkHeader(testCase.Origin, testCase.Header)
			if testCase.Error {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.Expected, actual)
			}
		})
	}
}

func TestParseWWWAuthenticateHeader(t *testing.T) {
	testCases := []struct {
		Header         string
		ExpectedScheme string
		ExpectedParams map[string]string
		Error          bool
	}{
		{
			Header:         `Basic realm="Dev", charset="UTF-8"`,
			ExpectedScheme: "Basic",
			ExpectedParams: map[string]string{
				"realm":   "Dev",
				"charset": "UTF-8",
			},
			Error: false,
		},
		{
			Header:         `Basic realm="Dev",charset="UTF-8"`,
			ExpectedScheme: "Basic",
			ExpectedParams: map[string]string{
				"realm":   "Dev",
				"charset": "UTF-8",
			},
			Error: false,
		},
		{
			Header:         `Basic realm="Dev",charset="ASCII",charset="UTF-8"`,
			ExpectedScheme: "Basic",
			ExpectedParams: map[string]string{
				"realm":   "Dev",
				"charset": "UTF-8",
			},
			Error: false,
		},
		{
			Header:         `Bearer realm="https://auth.docker.io/token",service="registry.docker.io",scope="registry:catalog:*",error="insufficient_scope"`,
			ExpectedScheme: "Bearer",
			ExpectedParams: map[string]string{
				"realm":   "https://auth.docker.io/token",
				"service": "registry.docker.io",
				"scope":   "registry:catalog:*",
				"error":   "insufficient_scope",
			},
			Error: false,
		},
		{
			Header: `Basic realm="Dev" charset="ASCII" charset="UTF-8"`,
			Error:  true,
		},
		{
			Header: `Basic realm="Dev" `,
			Error:  true,
		},
		{
			Header: `Basic realm="Dev",charset="UTF-8",`,
			Error:  true,
		},
		{
			Header: `Basic realm:"Dev"`,
			Error:  true,
		},
		{
			Header:         `Basic`,
			ExpectedScheme: "Basic",
			ExpectedParams: map[string]string{},
			Error:          false,
		},
		{
			Header: `Basic `,
			Error:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Header, func(t *testing.T) {
			scheme, params, err := ParseWWWAuthenticateHeader(testCase.Header)
			if testCase.Error {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, testCase.ExpectedScheme, scheme)
			assert.Equal(t, testCase.ExpectedParams, params)
		})
	}
}

func mustParseURL(t *testing.T, u string) *url.URL {
	v, err := url.Parse(u)
	require.NoError(t, err)

	return v
}
