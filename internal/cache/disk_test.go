package cache

import (
	"context"
	"log/slog"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiskCache(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	// Create a cache
	cache, err := NewDiskCache(path.Join(t.TempDir(), "cache.boltdb"))
	require.NoError(t, err)

	// Ensure that a non-existing item returns the expexted error
	data, err := cache.Get(context.TODO(), "foo")
	assert.Nil(t, data)
	assert.Equal(t, ErrNotExist, err)

	// Create an entry
	err = cache.Set(context.TODO(), "foo", []byte("bar"), nil)
	require.NoError(t, err)

	// Ensure it can be read
	data, err = cache.Get(context.TODO(), "foo")
	assert.Equal(t, []byte("bar"), data)
	assert.Equal(t, nil, err)

	// Replace the entry, with a TTL
	err = cache.Set(context.TODO(), "foo", []byte("bar"), &SetEntryOptions{Expires: time.Now().Add(-5 * time.Minute)})
	require.NoError(t, err)

	// Make sure it cannot be read, even though it hasn't been removed yet
	data, err = cache.Get(context.TODO(), "foo")
	assert.Nil(t, data)
	assert.Equal(t, ErrNotExist, err)

	// Remove expired entries
	removed, err := cache.DeleteExpiredEntries(time.Now())
	require.NoError(t, err)
	assert.Equal(t, 1, removed)

	// Make sure it cannot be read
	data, err = cache.Get(context.TODO(), "foo")
	assert.Nil(t, data)
	assert.Equal(t, ErrNotExist, err)

	require.NoError(t, cache.Close())
}
