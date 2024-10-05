package cache

import "testing"

func TestInMemoryCache(t *testing.T) {
	CacheTests(t, NewInMemoryCache())
}
