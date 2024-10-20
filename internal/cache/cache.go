package cache

import (
	"context"
	"errors"
	"time"
)

type SetEntryOptions struct {
	// Expires optionally controls when an entry should be invalidated.
	Expires time.Time
}

var (
	// ErrNotExist is returned if the entry does not exist.
	ErrNotExist = errors.New("entry does not exist")
)

// Cache stores key-value pairs for quick retrieval.
// It is safe to use a Cache implementation from multiple goroutines
// simultaneously.
type Cache interface {
	// Get retrieves an entry. If the entry does not exist, ErrNotExist is
	// returned.
	Get(ctx context.Context, key string) ([]byte, error)
	// Set inserts an entry.
	Set(ctx context.Context, key string, data []byte, options *SetEntryOptions) error
	// Delete deletes an entry. If the entry does not exist, nothing is done and a
	// nil error is returned.
	Delete(ctx context.Context, key string) error
}
