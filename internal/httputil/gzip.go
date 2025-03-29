package httputil

import (
	"compress/gzip"
	"net/http"
)

var _ http.ResponseWriter = (*GzipWriter)(nil)

// GzipWriter is an [http.ResponseWriter] that gzips the payload.
// Not safe for parallel use.
// Caller must set the correct content encoding header and close the writer.
type GzipWriter struct {
	http.ResponseWriter

	gzip *gzip.Writer
}

// Write implements [http.ResponseWriter].
func (g *GzipWriter) Write(p []byte) (int, error) {
	if g.gzip == nil {
		g.gzip = gzip.NewWriter(g.ResponseWriter)
	}

	return g.gzip.Write(p)
}

// Close implements [io.Closer].
func (g *GzipWriter) Close() error {
	if g.gzip == nil {
		return nil
	}

	return g.gzip.Close()
}
