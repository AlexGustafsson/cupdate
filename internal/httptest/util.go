package httptest

import (
	"io"
)

var _ io.Reader = (*ErrorReader)(nil)

// ErrorReader is an io.Reader that always returns an error.
// Useful for testing network errors when reading HTTP response bodies.
type ErrorReader struct {
	Error error
}

// Read implements io.Reader.
func (e ErrorReader) Read(p []byte) (n int, err error) {
	return 0, e.Error
}
