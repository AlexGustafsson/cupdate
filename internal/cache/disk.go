package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"time"
)

var _ Entry = (*diskEntry)(nil)
var _ EntryInfo = (*diskEntry)(nil)

type diskEntry struct {
	io.ReadCloser
	fileInfo fs.FileInfo
}

func (e diskEntry) ModTime() time.Time {
	return e.fileInfo.ModTime()
}

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

func (c *DiskCache) Stat(ctx context.Context, key string) (EntryInfo, bool, error) {
	path := path.Join(c.directory, c.formatKey(key))
	fileInfo, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	if fileInfo.IsDir() {
		return nil, false, nil
	}

	return fileInfo, true, nil
}

func (c *DiskCache) Get(ctx context.Context, key string) (Entry, error) {
	path := path.Join(c.directory, c.formatKey(key))
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotExist
	} else if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, ErrNotExist
	}

	return diskEntry{
		ReadCloser: file,
		fileInfo:   fileInfo,
	}, nil
}

func (c *DiskCache) Set(ctx context.Context, key string, reader io.Reader) error {
	path := path.Join(c.directory, c.formatKey(key))
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	if _, err := io.Copy(file, reader); err != nil {
		return err
	}

	return nil
}

func (c *DiskCache) Unset(ctx context.Context, key string) error {
	path := path.Join(c.directory, c.formatKey(key))
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}

	return nil
}

func (c *DiskCache) formatKey(key string) string {
	digest := sha256.Sum256([]byte(key))
	return hex.EncodeToString(digest[:])
}
