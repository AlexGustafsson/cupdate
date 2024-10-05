package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiskCache(t *testing.T) {
	cache, err := NewDiskCache(t.TempDir())
	require.NoError(t, err)

	CacheTests(t, cache)
}
