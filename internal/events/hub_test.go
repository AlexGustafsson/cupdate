package events

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHub(t *testing.T) {
	hub := NewHub[string]()

	ctx1, cancel1 := context.WithCancel(context.TODO())
	events1 := hub.Subscribe(ctx1)

	ctx2, cancel2 := context.WithCancel(context.TODO())
	events2 := hub.Subscribe(ctx2)

	// Events are broadcast
	go hub.Broadcast(context.TODO(), "Hello, World!")
	assert.Equal(t, "Hello, World!", <-events1)
	assert.Equal(t, "Hello, World!", <-events2)

	// Channels are closed on context cancellation
	go func() {
		_, ok := <-events1
		assert.False(t, ok)
	}()
	cancel1()

	// Non-closed channels are unaffected
	go hub.Broadcast(context.TODO(), "Hello, World!")
	assert.Equal(t, "Hello, World!", <-events2)

	cancel2()
}
