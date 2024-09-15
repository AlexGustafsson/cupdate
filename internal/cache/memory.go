package cache

import (
	"context"
	"encoding/json"
	"time"
)

type inMemoryItem struct {
	created time.Time
	data    []byte
}

type InMemoryCache struct {
	// TODO: garbage collection
	items map[string]inMemoryItem
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		items: make(map[string]inMemoryItem),
	}
}

func (c *InMemoryCache) Has(ctx context.Context, key string, maxAge time.Duration) (bool, error) {
	_, ok := c.items[key]
	return ok, nil
}

func (c *InMemoryCache) Get(ctx context.Context, key string, maxAge time.Duration) ([]byte, error) {
	item, ok := c.items[key]
	if !ok {
		return nil, nil
	}

	if time.Since(item.created) > maxAge {
		return nil, nil
	}

	return item.data, nil
}

func (c *InMemoryCache) Set(ctx context.Context, key string, content []byte) error {
	c.items[key] = inMemoryItem{
		created: time.Now(),
		data:    content,
	}
	return nil
}

func (c *InMemoryCache) GetJSON(ctx context.Context, key string, v any, maxAge time.Duration) error {
	data, err := c.Get(ctx, key, maxAge)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func (c *InMemoryCache) SetJSON(ctx context.Context, key string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return c.Set(ctx, key, data)
}
