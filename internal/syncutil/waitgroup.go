package syncutil

import (
	"context"
	"sync"
)

// WaitContext calls [sync.WaitGroup.Wait], returning when it completes or the
// context is canceled.
func WaitContext(ctx context.Context, wg *sync.WaitGroup) error {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		return nil
	}
}
