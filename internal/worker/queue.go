package worker

import (
	"iter"
	"slices"
	"sync"
	"time"
)

type Queue[T comparable] struct {
	cond    *sync.Cond
	backlog []T
	tokens  int
	closed  bool
}

func NewQueue[T comparable](burst int, tick time.Duration) *Queue[T] {
	q := &Queue[T]{
		cond:    sync.NewCond(&sync.Mutex{}),
		backlog: make([]T, 0),
		tokens:  burst,
	}

	// Fill the available backlog over time
	if tick > 0 {
		go func() {
			ticker := time.NewTicker(tick)

			for range ticker.C {
				q.cond.L.Lock()
				tokens := min(burst, q.tokens+1)
				changed := tokens != q.tokens
				q.tokens = tokens
				q.cond.L.Unlock()

				if changed {
					q.cond.Broadcast()
				}
			}
		}()
	}

	return q
}

func (q *Queue[T]) Push(items ...T) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if q.closed {
		panic("worker: push on closed queue")
	}

	// Push unique items to the backlog
	for _, item := range items {
		if !slices.Contains(q.backlog, item) {
			q.backlog = append(q.backlog, item)
		}
	}

	q.cond.Broadcast()
}

func (q *Queue[T]) Pull() iter.Seq[T] {
	return func(yield func(T) bool) {
		q.cond.L.Lock()
		defer q.cond.L.Unlock()

		for !q.closed {
			// Try to pick an item from the queue
			if len(q.backlog) > 0 && q.tokens > 0 {
				q.tokens--
				ok := yield(q.backlog[0])
				q.backlog = q.backlog[1:]
				if !ok {
					return
				}

				continue
			}

			// There's no backlog, wait for a wake up
			q.cond.Wait()
		}
	}
}

func (q *Queue[T]) Len() int {
	return len(q.backlog)
}

func (q *Queue[T]) Close() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.closed = true
	q.backlog = []T{}
	q.cond.Broadcast()
}
