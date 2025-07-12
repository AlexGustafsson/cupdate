package api

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"net/http"

	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/AlexGustafsson/cupdate/internal/oci"
)

type LogoProxy interface {
	ServeLogo(http.ResponseWriter, *http.Request, oci.Reference) error
}

var _ LogoProxy = (*LogoFSProxy)(nil)

type LogoFSProxy struct {
	FS fs.FS
}

func (p *LogoFSProxy) ServeLogo(w http.ResponseWriter, r *http.Request, reference oci.Reference) error {
	extensions := []string{".png", ".jpg", ".jpeg", ".svg", ".webp"}
	mimeTypes := []string{"image/png", "image/jpeg", "image/jpeg", "image/svg+xml", "image/webp"}

	var file fs.File
	var mimeType string
	for i, extension := range extensions {
		f, err := p.FS.Open(reference.Name() + extension)
		if err == nil {
			file = f
			mimeType = mimeTypes[i]
			break
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	if file == nil {
		return ErrNotFound
	}

	// Ask user agents to cache images found on disk for a few minutes
	w.Header().Set("Cache-Control", "max-age=300")
	w.Header().Set("Content-Type", mimeType)
	io.Copy(w, file)
	return nil
}

var _ LogoProxy = (*LogoHTTPProxy)(nil)

type LogoHTTPProxy struct {
	Client *httputil.Client
	GetURL func(context.Context, string) (string, error)
}

func (p *LogoHTTPProxy) ServeLogo(w http.ResponseWriter, r *http.Request, reference oci.Reference) error {
	url, err := p.GetURL(r.Context(), reference.String())
	if err != nil {
		return err
	}

	if url == "" {
		return ErrNotFound
	}

	// Redirect the user agents to the external URL, let the remote server handle
	// cache
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
	return nil
}

var _ LogoProxy = (*CompoundProxy)(nil)

type CompoundProxy struct {
	Proxies []LogoProxy
}

func (p CompoundProxy) ServeLogo(w http.ResponseWriter, r *http.Request, reference oci.Reference) error {
	for _, proxy := range p.Proxies {
		if err := proxy.ServeLogo(w, r, reference); err == ErrNotFound {
			continue
		} else if err != nil {
			return err
		}

		return nil
	}

	return ErrNotFound
}
