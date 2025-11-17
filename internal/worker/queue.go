package worker

import (
	"iter"
	"slices"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = (*Queue[any])(nil)

// Queue is a queue data structure implementation useful for tracking jobs.
// It enables multiple producers to queue items for multiple consumers.
// Consumers will sleep whilst they're waiting for jobs.
// The implementation also provides token bucket-based rate limiting.
type Queue[T comparable] struct {
	cond    *sync.Cond
	backlog []T
	tokens  int
	closed  bool

	burstGauge  prometheus.Gauge
	lengthGauge prometheus.Gauge
}

// NewQueue returns a new [Queue].
// An initialized queue should be closed by calling [Queue.Close]
//
//   - burst controls the target number of items that can be taken from the
//     queue.
//   - tick is the target mean time between items being taken from the queue.
func NewQueue[T comparable](burst int, tick time.Duration) *Queue[T] {
	q := &Queue[T]{
		cond:    sync.NewCond(&sync.Mutex{}),
		backlog: make([]T, 0),
		tokens:  burst,

		burstGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "cupdate",
			Subsystem: "worker",
			Name:      "available_burst",
		}),
		lengthGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "cupdate",
			Subsystem: "worker",
			Name:      "queue_length",
		}),
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
					q.burstGauge.Set(float64(tokens))
					q.cond.Broadcast()
				}
			}
		}()
	}

	return q
}

// PushFront puts one or more unique items at the front of the queue.
// Panics if the queue is closed.
func (q *Queue[T]) PushFront(items ...T) {
	q.setFunc(func(backlog []T) []T {
		// Add unique items from the backlog to the end of the list of new items
		for _, item := range backlog {
			if !slices.Contains(items, item) {
				items = append(items, item)
			}
		}
		return items
	})
}

// PushBack puts one or more unique items at the back of the queue.
// Panics if the queue is closed.
func (q *Queue[T]) PushBack(items ...T) {
	q.setFunc(func(backlog []T) []T {
		// Add unique items to the backlog
		for _, item := range items {
			if !slices.Contains(backlog, item) {
				backlog = append(backlog, item)
			}
		}
		return backlog
	})
}

func (q *Queue[T]) setFunc(f func([]T) []T) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if q.closed {
		panic("worker: push on closed queue")
	}

	q.backlog = f(q.backlog)

	q.lengthGauge.Set(float64(len(q.backlog)))
	q.cond.Broadcast()
}

// Pull returns an iterator which will feed a consumer with items.
// The iterator will wait for new items if the queue is emptied and is only
// closed when the queue itself is closed, or when the consumer decides to stop
// processing items.
func (q *Queue[T]) Pull() iter.Seq[T] {
	return func(yield func(T) bool) {
		q.cond.L.Lock()
		defer q.cond.L.Unlock()

		for !q.closed {
			// Try to pick an item from the queue
			if len(q.backlog) > 0 && q.tokens > 0 {
				q.tokens--
				q.burstGauge.Set(float64(q.tokens))
				q.lengthGauge.Dec()
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

// Len returns the size of the queue.
func (q *Queue[T]) Len() int {
	return len(q.backlog)
}

// AvailableBurst returns the number of items that could currently be taken from
// the queue without being rate limited.
func (q *Queue[T]) AvailableBurst() int {
	return q.tokens
}

// Close closes the queue and frees resources.
// Closing a queue will wake any waiting consumer and exit their pull loops.
// Continuing to use a closed queue will panic.
func (q *Queue[T]) Close() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.closed = true
	q.backlog = []T{}
	q.cond.Broadcast()
}

// Collect implements prometheus.Collector.
func (w *Queue[T]) Collect(ch chan<- prometheus.Metric) {
	w.burstGauge.Collect(ch)
	w.lengthGauge.Collect(ch)
}

// Describe implements prometheus.Collector.
func (w *Queue[T]) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(w, descs)
}
