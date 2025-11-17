package worker

import (
	"iter"
	"slices"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/oci"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueueBacklog(t *testing.T) {
	q := NewQueue[string](5, 0)
	defer q.Close()

	items := []string{
		"1",
		"2",
		"3",
		"4",
		"5",
	}

	q.PushBack(items...)

	assert.Equal(t, 5, q.Len())

	next, stop := iter.Pull(q.Pull())

	for i := 0; i < len(items); i++ {
		expected := items[i]
		actual, ok := next()
		assert.Equal(t, expected, actual)
		assert.True(t, ok)
	}

	stop()

	_, ok := next()
	assert.False(t, ok)
}

func TestQueuePushBack(t *testing.T) {
	q := NewQueue[string](5, 0)
	defer q.Close()

	items := []string{
		"1",
		"2",
		"3",
		"4",
		"5",
	}
	q.PushBack(items[2:]...)
	q.PushFront(items[:2]...)

	assert.Equal(t, 5, q.Len())

	next, stop := iter.Pull(q.Pull())

	for i := 0; i < len(items); i++ {
		expected := items[i]
		actual, ok := next()
		assert.Equal(t, expected, actual)
		assert.True(t, ok)
	}

	stop()

	_, ok := next()
	assert.False(t, ok)
}

func TestQueueNoBacklog(t *testing.T) {
	q := NewQueue[string](2, 0)

	items := make(chan string)
	go func() {
		for item := range q.Pull() {
			items <- item
		}
		close(items)
	}()

	q.PushBack("1", "2")

	assert.Equal(t, "1", <-items)
	assert.Equal(t, "2", <-items)

	q.Close()
	<-items

	assert.Equal(t, 0, q.Len())
}

func TestQueueMultipleConsumers(t *testing.T) {
	// By having 1 burst we ensure that a worker will only pull a single item each
	// time it's woken up, assuming it's quicker to process than the tick
	q := NewQueue[int](1, 1*time.Millisecond)
	defer q.Close()

	// Keep track of which worker handled with item
	items := make([]int, 5)
	for i := 0; i < 5; i++ {
		go func() {
			for item := range q.Pull() {
				items[item] = i
			}
		}()
	}

	q.PushBack(0, 1, 2, 3, 4)

	// Wait for all items to be processed
	<-time.After(1 * time.Second)

	slices.Sort(items)
	assert.NotEqual(t, items[0], items[4], "different workers handled requests")
}

func TestQueueClose(t *testing.T) {
	q := NewQueue[string](0, 0)

	next, stop := iter.Pull(q.Pull())
	defer stop()

	closed := make(chan struct{})
	go func() {
		next()
		close(closed)
	}()

	q.Close()
	<-closed

	assert.Panics(t, func() {
		q.PushBack("panic when closed")
	})
}

func TestQueueEmptiedOnClose(t *testing.T) {
	q := NewQueue[string](5, 0)

	q.PushBack("1", "2", "3", "4", "5")

	assert.Equal(t, 5, q.Len())
	q.Close()
	assert.Equal(t, 0, q.Len())
}

func TestQueueDeduplication(t *testing.T) {
	q := NewQueue[oci.Reference](2, 0)
	defer q.Close()

	assert.Equal(t, 0, q.Len())

	ref, err := oci.ParseReference("mongo:4")
	require.NoError(t, err)

	q.PushBack(ref)
	assert.Equal(t, 1, q.Len())

	q.PushBack(ref)
	assert.Equal(t, 1, q.Len())

	ref.Tag = "5"
	q.PushBack(ref)
	assert.Equal(t, 2, q.Len())
}
