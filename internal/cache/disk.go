package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path"
	"time"
)

var _ Cache = (*DiskCache)(nil)

type DiskCache struct {
	directory string
}

func NewDiskCache(directory string) (*DiskCache, error) {
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return nil, err
	}

	return &DiskCache{
		directory: directory,
	}, nil
}

func (c *DiskCache) Has(ctx context.Context, key string, maxAge time.Duration) (bool, error) {
	path := path.Join(c.directory, c.formatKey(key))
	stat, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if stat.IsDir() {
		return false, nil
	}

	if maxAge > 0 && time.Since(stat.ModTime()) > maxAge {
		return false, nil
	}

	return !stat.IsDir(), nil
}

func (c *DiskCache) Get(ctx context.Context, key string, maxAge time.Duration) ([]byte, error) {
	exists, err := c.Has(ctx, key, maxAge)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	path := path.Join(c.directory, c.formatKey(key))
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (c *DiskCache) Set(ctx context.Context, key string, data []byte) error {
	path := path.Join(c.directory, c.formatKey(key))
	return os.WriteFile(path, data, 0600)
}

func (c *DiskCache) GetJSON(ctx context.Context, key string, v any, maxAge time.Duration) error {
	data, err := c.Get(ctx, key, maxAge)
	if err != nil {
		return err
	}

	if data == nil {
		return nil
	}

	return json.Unmarshal(data, v)
}

func (c *DiskCache) SetJSON(ctx context.Context, key string, v any) error {
	if v == nil {
		// TODO: Remove instead?
		return nil
	}

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return c.Set(ctx, key, data)
}

func (c *DiskCache) formatKey(key string) string {
	digest := sha256.Sum256([]byte(key))
	return hex.EncodeToString(digest[:])
}
