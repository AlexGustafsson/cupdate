package ratelimit

import (
	"context"
	"iter"
	"time"
)

// Channel limits thate rate of ch using a leaky bucket. Returns an iterator of
// the items sent to the channel as fast as the rate allows.
//
// The bucket will contain (and starts with) burst tokens. Every tick the bucket
// is filled with another token.
func Channel[T any](ctx context.Context, burst int, tick time.Duration, ch <-chan T) iter.Seq2[T, int] {
	ticker := time.NewTicker(tick)

	bucket := make(chan struct{}, burst)

	// Fill the bucket
	for i := 0; i < burst; i++ {
		bucket <- struct{}{}
	}

	// Fill the bucket over time
	go func() {
		defer close(bucket)

		for range ticker.C {
			select {
			case bucket <- struct{}{}:
			default:
				// Bucket is full
			}
		}
	}()

	return func(yield func(T, int) bool) {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				// Deadline exceeded / context canceled
				return
			case _, ok := <-bucket:
				if !ok {
					// Channel closed
					return
				}

				select {
				case <-ctx.Done():
					// Deadline exceeded / context canceled
					return
				case v, ok := <-ch:
					if !ok {
						// Channel closed
						return
					}

					if !yield(v, len(bucket)) {
						// Loop exited
						return
					}
				}
			}
		}
	}
}
