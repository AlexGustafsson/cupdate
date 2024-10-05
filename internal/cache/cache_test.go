package cache

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CacheTests(t *testing.T, cache Cache) {
	// Setting an entry works
	err := cache.Set(context.TODO(), "foo", strings.NewReader("bar"))
	require.NoError(t, err)

	// Retrieving an existing entry works
	entry, err := cache.Get(context.TODO(), "foo")
	require.NoError(t, err)

	// Content of entries persist
	content, err := io.ReadAll(entry)
	require.NoError(t, err)
	assert.Equal(t, []byte("bar"), content)

	// Entries' EntryInfo is retrievable
	info, exists, err := cache.Stat(context.TODO(), "foo")
	require.NoError(t, err)
	assert.True(t, exists)
	assert.NotNil(t, info)

	// Removing an entry works
	err = cache.Unset(context.TODO(), "foo")
	require.NoError(t, err)

	// Non-existing entries can be identified
	info, exists, err = cache.Stat(context.TODO(), "foo")
	require.NoError(t, err)
	assert.False(t, exists)
	assert.Nil(t, info)

	// Reading a non-existing entry does not work
	entry, err = cache.Get(context.TODO(), "foo")
	require.ErrorIs(t, err, ErrNotExist)
	assert.Nil(t, entry)
}
