package oci

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"net/http"
)

type BlobInfo struct {
	// ContentType is the blob's content type.
	ContentType string
	// ContentLength is the size of the blob as reported by the server, if
	// reported at all. SHOULD be reported by all servers.
	ContentLength int64
	// Digest is the digest of the blob as reported by the server, if reported at
	// all. SHOULD be reported by all servers. If it is reported MUST match the
	// blob's actual digest.
	Digest string
}

func blobInfoFromResponse(res *http.Response) BlobInfo {
	return BlobInfo{
		ContentType:   res.Header.Get("Content-Type"),
		ContentLength: res.ContentLength,
		Digest:        res.Header.Get("Docker-Content-Digest"),
	}
}

type Blob interface {
	io.Reader
	io.Closer
	// Digest returns the blob's digest so far.
	// Valid once the entire blob has been read.
	// It is the caller's responsibility to ensure the entire blob has been read.
	Digest() string
	Info() BlobInfo
}

var _ Blob = (*blob)(nil)

type blob struct {
	reader io.ReadCloser
	hash   hash.Hash
	info   BlobInfo
}

func blobFromResponse(res *http.Response) *blob {
	return newBlobResponse(res.Body, blobInfoFromResponse(res))
}

func newBlobResponse(reader io.ReadCloser, info BlobInfo) *blob {
	hash := sha256.New()
	return &blob{
		reader: newTeeReadCloser(reader, hash),
		hash:   hash,
		info:   info,
	}
}

// Read implements Blob.
func (b *blob) Read(p []byte) (n int, err error) {
	return b.reader.Read(p)
}

// Close implements Blob.
func (b *blob) Close() error {
	return b.reader.Close()
}

// Digest implements Blob.
func (b *blob) Digest() string {
	return "sha256:" + hex.EncodeToString(b.hash.Sum(nil))
}

// Digest implements Blob.
func (b *blob) Info() BlobInfo {
	return b.info
}
