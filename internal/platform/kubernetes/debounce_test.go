package kubernetes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDebounce(t *testing.T) {
	incoming := make(chan time.Time)

	outgoing := Debounce(incoming, 1*time.Second)

	// The first message is passed immediately
	in := time.Now()
	incoming <- in
	<-outgoing
	assert.LessOrEqual(t, time.Since(in), 200*time.Millisecond)

	// Messages sent during the debounce window are consumed without locking
	incoming <- time.Now()
	incoming <- time.Now()
	incoming <- time.Now()
	incoming <- time.Now()
	incoming <- time.Now()

	// The last message sent in a time window is passed in the next
	in = time.Now()
	incoming <- in
	<-outgoing
	assert.GreaterOrEqual(t, time.Since(in), 1*time.Second)
}
