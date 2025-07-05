package kubernetes

import (
	"time"
)

// Debounce debounces messages sent to ch, returning at most one message per
// interval. The latest message sent to ch during a debounced time window, if
// any, will be sent the next window.
func Debounce[T any](ch <-chan T, interval time.Duration) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		var debouncedValue *T
		for {
			// For each time window, either send the latest debounced value, or wait
			// for one to be sent on ch
			var value T
			var ok bool
			if debouncedValue == nil {
				value, ok = <-ch
			} else {
				value = *debouncedValue
				debouncedValue = nil
				ok = true
			}
			if !ok {
				return
			}

			out <- value

			// During the debounced window, keep track of the latest value received
			debounce := time.After(interval)
		debounce:
			for {
				select {
				case v, ok := <-ch:
					if !ok {
						return
					}

					debouncedValue = &v
				case <-debounce:
					break debounce
				}
			}
		}
	}()

	return out
}
