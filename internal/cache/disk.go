package cache

import (
	"context"
	"log/slog"
	"strconv"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

var _ Cache = (*DiskCache)(nil)

var (
	bucketNameCache = []byte("cachev1")
	bucketNameTTL   = []byte("ttlv1")
)

// DiskCache implements [Cache] by storing values to disk.
type DiskCache struct {
	db     *bolt.DB
	cancel chan struct{}
	wg     sync.WaitGroup
}

// NewDiskCache returns a [DiskCache] that stores its data in a file at the
// specified path.
func NewDiskCache(path string) (*DiskCache, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketNameCache); err != nil {
			return err
		}

		if _, err := tx.CreateBucketIfNotExists(bucketNameTTL); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}

	cache := &DiskCache{
		db:     db,
		cancel: make(chan struct{}),
	}

	cache.wg.Add(1)
	go func() {
		defer cache.wg.Done()

		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				slog.Debug("Invalidating expired cache entries")
				removed, err := cache.DeleteExpiredEntries(time.Now())
				if err != nil {
					slog.Warn("Failed to invalidate expired cache entries", slog.Any("error", err))
					continue
				}

				if removed > 0 {
					slog.Debug("Invalidated expired cache entries", slog.Int("entries", removed))
				}
			case <-cache.cancel:
				slog.Debug("Cache is closing - stopping cache invalidation")
				return
			}
		}
	}()

	return cache, nil
}

// Get implements [Cache].
func (d *DiskCache) Get(ctx context.Context, key string) ([]byte, error) {
	var data []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		cacheBucket := tx.Bucket(bucketNameCache)
		ttlBucket := tx.Bucket(bucketNameTTL)

		entry := cacheBucket.Get([]byte(key))
		if entry == nil {
			return ErrNotExist
		}

		// Make sure the entry is not expired and just haven't been invalidated yet
		ttl := ttlBucket.Get([]byte(key))
		if ttl != nil {
			expires, err := bytesToTime(ttl)
			if err != nil {
				return err
			}

			if expires.Before(time.Now()) {
				return ErrNotExist
			}
		}

		data = make([]byte, len(entry))
		copy(data, entry)
		return nil
	})
	return data, err
}

// Set implements [Cache].
func (d *DiskCache) Set(ctx context.Context, key string, data []byte, options *SetEntryOptions) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		// Insert the data
		cacheBucket := tx.Bucket(bucketNameCache)
		if err := cacheBucket.Put([]byte(key), data); err != nil {
			return err
		}

		// Map the key to its expiry time. Note that this implementation doesn't
		// support multiple keys expiring at the exact same time. But let's be real,
		// for our use case (caching HTTP requests in cupdate), we're extremely
		// unlikely to ever run into this issue. Let's just keep the code simple and
		// easy to maintain instead
		if options != nil && !options.Expires.IsZero() {
			ttlBucket := tx.Bucket(bucketNameTTL)
			if err := ttlBucket.Put([]byte(key), timeToBytes(options.Expires)); err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete implements [Cache].
func (d *DiskCache) Delete(ctx context.Context, key string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		cacheBucket := tx.Bucket(bucketNameCache)
		ttlBucket := tx.Bucket(bucketNameTTL)

		if err := cacheBucket.Delete([]byte(key)); err != nil {
			return err
		}

		if err := ttlBucket.Delete([]byte(key)); err != nil {
			return err
		}

		return nil
	})
}

// DeleteExpiredEntries deletes all expired entries.
func (d *DiskCache) DeleteExpiredEntries(maxAge time.Time) (int, error) {
	removed := 0
	err := d.db.Update(func(tx *bolt.Tx) error {
		cacheBucket := tx.Bucket(bucketNameCache)
		ttlBucket := tx.Bucket(bucketNameTTL)

		// In our use case, we'll have so few entries that we might as well be naive
		// and loop through them all when cleaning up in order to make sure we don't
		// have any edge cases like a key having multiple TTL entries if set more
		// than once
		keysToDelete := make([][]byte, 0)
		ttlBucket.ForEach(func(k, v []byte) error {
			expires, err := bytesToTime(v)
			if err != nil {
				return err
			}

			if expires.Before(maxAge) {
				keysToDelete = append(keysToDelete, k)
			}

			return nil
		})

		for _, key := range keysToDelete {
			if err := cacheBucket.Delete(key); err != nil {
				return err
			}

			if err := ttlBucket.Delete(key); err != nil {
				return err
			}

			removed++
		}

		return nil
	})

	return removed, err
}

// Close closes the cache.
func (d *DiskCache) Close() error {
	slog.Debug("Closing cache")
	close(d.cancel)
	slog.Debug("Waiting for invalidation worker to complete")
	d.wg.Wait()
	slog.Debug("Closing boltdb")
	return d.db.Close()
}

// timeToBytes converts time into a sortable byte slice.
// The format is Unix nano seconds (int64) and thus have some range limitations.
func timeToBytes(time time.Time) []byte {
	return []byte(strconv.FormatInt(time.UnixNano(), 10))
}

// bytesToTime returns a time as returned by timesToBytes.
func bytesToTime(v []byte) (time.Time, error) {
	i, err := strconv.ParseInt(string(v), 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, i), nil
}
