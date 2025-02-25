package web

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

//go:embed public/*
var public embed.FS

// MustNewServer calls [NewServer] and panics on any error.
func MustNewServer() http.Handler {
	handler, err := NewServer()
	if err != nil {
		panic(err)
	}
	return handler
}

// NewServer creates a new server for the web content.
func NewServer() (http.Handler, error) {
	public, err := fs.Sub(fs.FS(public), "public")
	if err != nil {
		return nil, err
	}

	fileServer := http.FileServer(http.FS(public))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		file, err := public.Open(strings.TrimPrefix(path, "/"))
		if os.IsNotExist(err) {
			// Serve HTML as fallback
			indexFile, err := public.Open("index.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			data, err := io.ReadAll(indexFile)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file.Close()

		// Serve a public file
		fileServer.ServeHTTP(w, r)
	}), nil
}
