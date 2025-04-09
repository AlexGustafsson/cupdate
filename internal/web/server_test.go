package web

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	testCases := []struct {
		Path                string
		Accept              string
		FS                  fstest.MapFS
		ExpectedStatusCode  int
		ExpectedBody        string
		ExpectedContentType string
	}{
		{
			Path:   "http://example.com",
			Accept: "",
			FS: fstest.MapFS{
				"index.html": &fstest.MapFile{Data: []byte("<html></html>")},
			},
			ExpectedStatusCode:  http.StatusOK,
			ExpectedBody:        "<html></html>",
			ExpectedContentType: "text/html; charset=utf-8",
		},
		{
			Path:   "http://example.com/",
			Accept: "",
			FS: fstest.MapFS{
				"index.html": &fstest.MapFile{Data: []byte("<html></html>")},
			},
			ExpectedStatusCode:  http.StatusOK,
			ExpectedBody:        "<html></html>",
			ExpectedContentType: "text/html; charset=utf-8",
		},
		{
			Path:   "http://example.com/index.html",
			Accept: "",
			FS: fstest.MapFS{
				"index.html": &fstest.MapFile{Data: []byte("<html></html>")},
			},
			ExpectedStatusCode:  http.StatusMovedPermanently,
			ExpectedContentType: "",
		},
		{
			Path:   "http://example.com/index2.html",
			Accept: "text/html",
			FS: fstest.MapFS{
				"index2.html": &fstest.MapFile{Data: []byte("<html></html>")},
			},
			ExpectedStatusCode:  http.StatusOK,
			ExpectedBody:        "<html></html>",
			ExpectedContentType: "text/html; charset=utf-8",
		},
		{
			Path:   "http://example.com/not-found",
			Accept: "",
			FS: fstest.MapFS{
				"index.html": &fstest.MapFile{Data: []byte("<html></html>")},
			},
			ExpectedStatusCode:  http.StatusNotFound,
			ExpectedBody:        "404 page not found\n",
			ExpectedContentType: "text/plain; charset=utf-8",
		},
		{
			Path:   "http://example.com/not-found",
			Accept: "*/*",
			FS: fstest.MapFS{
				"index.html": &fstest.MapFile{Data: []byte("<html></html>")},
			},
			ExpectedStatusCode:  http.StatusOK,
			ExpectedBody:        "<html></html>",
			ExpectedContentType: "text/html; charset=utf-8",
		},
		{
			Path:   "http://example.com/not-found",
			Accept: "text/html",
			FS: fstest.MapFS{
				"index.html": &fstest.MapFile{Data: []byte("<html></html>")},
			},
			ExpectedStatusCode:  http.StatusOK,
			ExpectedBody:        "<html></html>",
			ExpectedContentType: "text/html; charset=utf-8",
		},
		{
			Path:                "http://example.com/not-found",
			Accept:              "text/html",
			FS:                  fstest.MapFS{},
			ExpectedStatusCode:  http.StatusNotFound,
			ExpectedBody:        "404 page not found\n",
			ExpectedContentType: "text/plain; charset=utf-8",
		},
		{
			Path:                "http://example.com/not-found/",
			Accept:              "text/html",
			FS:                  fstest.MapFS{},
			ExpectedStatusCode:  http.StatusNotFound,
			ExpectedBody:        "404 page not found\n",
			ExpectedContentType: "text/plain; charset=utf-8",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Path, func(t *testing.T) {
			server := NewServer(testCase.FS)

			req := httptest.NewRequest(http.MethodGet, testCase.Path, nil)
			if testCase.Accept != "" {
				req.Header.Set("Accept", testCase.Accept)
			}

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)
			res := w.Result()

			assert.Equal(t, testCase.ExpectedStatusCode, res.StatusCode)
			assert.Equal(t, testCase.ExpectedContentType, res.Header.Get("Content-Type"))
			if testCase.ExpectedBody != "" {
				body, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, testCase.ExpectedBody, string(body))
			}
		})
	}
}
