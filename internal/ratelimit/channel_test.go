package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	ch := make(chan time.Time, 10)
	go func() {
		// Fill the channel
		for i := 0; i < 10; i++ {
			ch <- time.Now()
		}

		ch <- time.Now()
		ch <- time.Now()
		close(ch)
	}()

	expectedDiffs := []time.Duration{
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		1 * time.Second,
		2 * time.Second,
	}

	epsilon := 2 * time.Millisecond

	i := 0
	start := time.Now()
	for range Channel(context.TODO(), 10, 1*time.Second, ch) {
		// Expect the first ten items to be handled instantenously (burst)
		assert.WithinDuration(t, time.Now(), start, expectedDiffs[i]+epsilon)

		i++
	}

	assert.Equal(t, 12, i)
}
