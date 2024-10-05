package cache

import (
	"context"
	"errors"
	"io"
	"time"
)

// Entry is an io.ReadCloser.
// Entries may also confirm to [EntryInfo] if the data is available.
type Entry = io.ReadCloser

type EntryInfo interface {
	ModTime() time.Time
}

var (
	ErrNotExist = errors.New("entry does not exist")
)

type Cache interface {
	Stat(ctx context.Context, key string) (EntryInfo, bool, error)
	Get(ctx context.Context, key string) (Entry, error)
	Set(ctx context.Context, key string, reader io.Reader) error
	Unset(ctx context.Context, key string) error
}
