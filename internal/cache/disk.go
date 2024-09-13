package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path"
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

func (d *DiskCache) Has(ctx context.Context, key string) (bool, error) {
	path := path.Join(d.directory, d.formatKey(key))
	stat, err := os.Stat(path)
	if err == os.ErrNotExist {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return !stat.IsDir(), nil
}

func (d *DiskCache) Get(ctx context.Context, key string) ([]byte, error) {
	path := path.Join(d.directory, d.formatKey(key))
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (d *DiskCache) Set(ctx context.Context, key string, data []byte) error {
	path := path.Join(d.directory, d.formatKey(key))
	return os.WriteFile(path, data, 0600)
}

func (d *DiskCache) formatKey(key string) string {
	digest := sha256.Sum256([]byte(key))
	return hex.EncodeToString(digest[:])
}
