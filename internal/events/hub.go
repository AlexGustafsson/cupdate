package events

import (
	"context"
	"sync"
)

type Hub[T any] struct {
	mutex    sync.Mutex
	channels map[chan T]struct{}
}

func NewHub[T any]() *Hub[T] {
	return &Hub[T]{
		channels: make(map[chan T]struct{}, 0),
	}
}

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
