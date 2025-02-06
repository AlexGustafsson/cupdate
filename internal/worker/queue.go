package worker

import (
	"iter"
	"slices"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Queue[T comparable] struct {
	cond    *sync.Cond
	backlog []T
	tokens  int
	closed  bool

	burstGauge  prometheus.Gauge
	lengthGauge prometheus.Gauge
}

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

	q.lengthGauge.Set(float64(len(q.backlog)))
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

func (q *Queue[T]) Len() int {
	return len(q.backlog)
}

func (q *Queue[T]) AvailableBurst() int {
	return q.tokens
}

func (q *Queue[T]) Close() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.closed = true
	q.backlog = []T{}
	q.cond.Broadcast()
}

// Collect implements [prometheus.Collector].
func (w *Queue[T]) Collect(ch chan<- prometheus.Metric) {
	w.burstGauge.Collect(ch)
	w.lengthGauge.Collect(ch)
}

// Describe implements [prometheus.Collector].
func (w *Queue[T]) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(w, descs)
}
