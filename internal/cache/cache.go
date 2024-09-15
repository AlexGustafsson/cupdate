package cache

import (
	"context"
	"time"
)

type Cache interface {
	Has(ctx context.Context, key string, maxAge time.Duration) (bool, error)
	Get(ctx context.Context, key string, maxAge time.Duration) ([]byte, error)
	Set(ctx context.Context, key string, content []byte) error

	// TODO: Just implement a generic request cache instead. Cache 200 responses?
	// That way new versions can add fields without invalidating cache.
	GetJSON(ctx context.Context, key string, v any, maxAge time.Duration) error
	SetJSON(ctx context.Context, key string, v any) error
}
