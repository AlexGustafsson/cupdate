package web

import (
	"embed"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
)

//go:embed public/*
var public embed.FS

// MustNewEmbeddedServer calls [NewEmbeddedServer] and panics on any error.
func MustNewEmbeddedServer() http.Handler {
	handler, err := NewEmbeddedServer()
	if err != nil {
		panic(err)
	}
	return handler
}

// NewEmbeddedServer creates a new SPA web server for the embedded web content.
func NewEmbeddedServer() (http.Handler, error) {
	public, err := fs.Sub(fs.FS(public), "public")
	if err != nil {
		return nil, err
	}

	return NewServer(public), nil
}

// NewServer creates a new SPA web server, serving files from fileServer.
func NewServer(public fs.FS) http.Handler {
	fileServer := http.FileServerFS(public)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gzip := &httputil.GzipWriter{ResponseWriter: w}
			defer gzip.Close()

			w = gzip
		}

		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		// If the file does not exist, serve index.html as fallback if the user
		// agent wants HTTP
		accept := httputil.Accepts(r.Header.Get("Accept"), "text/html")
		if accept == "text/html" {
			file, err := public.Open(strings.TrimPrefix(path, "/"))
			if errors.Is(err, fs.ErrNotExist) {
				indexFile, err := public.Open("index.html")
				if errors.Is(err, fs.ErrNotExist) {
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					w.Write([]byte("404 page not found\n"))
					return
				} else if err != nil {
					slog.ErrorContext(r.Context(), "Failed to open fallback file", slog.Any("error", err))
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				data, err := io.ReadAll(indexFile)
				if err != nil {
					slog.ErrorContext(r.Context(), "Failed to read fallback file", slog.Any("error", err))
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				w.Write(data)
				return
			} else if err != nil {
				slog.ErrorContext(r.Context(), "Failed to open fallback file", slog.Any("error", err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			file.Close()
		}

		// Add cache header for (assumed) immutable assets
		if strings.HasPrefix(path, "/assets/") {
			// SEE: https://web.dev/articles/http-cache#versioned-urls
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}

		fileServer.ServeHTTP(w, r)
	})
}
