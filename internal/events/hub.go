package events

import (
	"context"
	"sync"
)

// Hub is an event hub allowing broadcasting of events to subscribers.
type Hub[T any] struct {
	mutex    sync.Mutex
	channels map[chan T]struct{}
}

func NewHub[T any]() *Hub[T] {
	return &Hub[T]{
		channels: make(map[chan T]struct{}),
	}
}

// Subscribe to events.
// Returns a channel which will receives all broadcast events.
// The channel is closed whenever the context expires.
// Subscribers should receive events from the channel in a timely manner.
func (h *Hub[T]) Subscribe(ctx context.Context) <-chan T {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	ch := make(chan T)

	go func() {
		<-ctx.Done()
		h.mutex.Lock()
		defer h.mutex.Unlock()

		delete(h.channels, ch)

		close(ch)
	}()

	h.channels[ch] = struct{}{}

	return ch
}

// Broadcast an event to all subscribers synchronously.
// Returns an error if the context expires before all subscribers were invoked.
func (h *Hub[T]) Broadcast(ctx context.Context, event T) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for ch := range h.channels {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- event:
		}
	}

	return nil
}
