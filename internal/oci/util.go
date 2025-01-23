package oci

import (
	"io"
	"strings"
)

// NameFromAPI returns the OCI name based on the distribution spec API endpoint.
// Assumes name has at least two components.
// SEE: https://github.com/opencontainers/distribution-spec/blob/main/spec.md#endpoints
func NameFromAPI(path string) string {
	// /v2/<name>/[blobs,manifests,tags,referrers
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return ""
	}

	if parts[0] != "" || parts[1] != "v2" {
		return ""
	}

	components := 2
loop:
	for i := 4; i < len(parts); i++ {
		switch parts[i] {
		case "blobs", "manifests", "tags", "referrers":
			break loop
		default:
			components++
		}
	}

	return strings.Join(parts[2:2+components], "/")
}

var _ io.Reader = (*teeReadCloser)(nil)
var _ io.Closer = (*teeReadCloser)(nil)

type teeReadCloser struct {
	reader io.Reader
	closer io.Closer
	writer io.Writer
}

func newTeeReadCloser(r io.ReadCloser, w io.Writer) *teeReadCloser {
	return &teeReadCloser{
		reader: io.TeeReader(r, w),
		closer: r,
		writer: w,
	}
}

// Read implements io.Reader.
func (t *teeReadCloser) Read(p []byte) (n int, err error) {
	return t.reader.Read(p)
}

// Close implements io.Closer.
func (t *teeReadCloser) Close() error {
	return t.closer.Close()
}
