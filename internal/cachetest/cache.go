package cachetest

import (
	"path"
	"testing"

	"github.com/AlexGustafsson/cupdate/internal/cache"
	"github.com/stretchr/testify/require"
)

// NewCache returns an initialized [cache.Cache].
func NewCache(t *testing.T) cache.Cache {
	cache, err := cache.NewDiskCache(path.Join(t.TempDir(), "testcache"))
	require.NoError(t, err)

	return cache
}
