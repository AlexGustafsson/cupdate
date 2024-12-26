package events

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHub(t *testing.T) {
	hub := NewHub[string]()

	ctx1, cancel1 := context.WithCancel(context.TODO())
	events1 := hub.Subscribe(ctx1)

	ctx2, cancel2 := context.WithCancel(context.TODO())
	events2 := hub.Subscribe(ctx2)

	// Events are broadcast
	go hub.Broadcast(context.TODO(), "event 1")
	assert.Equal(t, "event 1", <-events1)
	assert.Equal(t, "event 1", <-events2)

	// Channels are closed on context cancellation
	go func() {
		<-time.After(200 * time.Millisecond)
		cancel1()
	}()
	v, ok := <-events1
	assert.Equal(t, "", v)
	assert.False(t, ok)

	// Non-closed channels are unaffected
	go hub.Broadcast(context.TODO(), "event 2")
	assert.Equal(t, "event 2", <-events2)

	cancel2()
}
