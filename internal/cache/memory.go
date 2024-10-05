package cache

import (
	"bytes"
	"context"
	"io"
	"time"
)

var _ Entry = (*memoryEntry)(nil)
var _ EntryInfo = (*memoryEntry)(nil)

type memoryEntry struct {
	io.ReadCloser

	modTime time.Time
}

func (e memoryEntry) ModTime() time.Time {
	return e.modTime
}

type memoryEntryStore struct {
	data    []byte
	modTime time.Time
}

var _ Cache = (*InMemoryCache)(nil)

type InMemoryCache struct {
	// TODO: garbage collection?
	entries map[string]memoryEntryStore
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		entries: make(map[string]memoryEntryStore),
	}
}

func (c *InMemoryCache) Stat(ctx context.Context, key string) (EntryInfo, bool, error) {
	entry, ok := c.entries[key]
	if !ok {
		return nil, false, nil
	}

	return memoryEntry{modTime: entry.modTime}, true, nil
}

func (c *InMemoryCache) Get(ctx context.Context, key string) (Entry, error) {
	entry, ok := c.entries[key]
	if !ok {
		return nil, ErrNotExist
	}

	return memoryEntry{
		ReadCloser: io.NopCloser(bytes.NewReader(entry.data)),
		modTime:    entry.modTime,
	}, nil
}

func (c *InMemoryCache) Set(ctx context.Context, key string, r io.Reader) error {
	modTime := time.Now()
	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, r); err != nil {
		return err
	}

	c.entries[key] = memoryEntryStore{
		data:    buffer.Bytes(),
		modTime: modTime,
	}
	return nil
}

func (c *InMemoryCache) Unset(ctx context.Context, key string) error {
	delete(c.entries, key)
	return nil
}
