package pipeline

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cache"
)

var _ context.Context = Context[any]{}

type Context[T any] struct {
	Data T

	ctx     context.Context
	cache   cache.Cache
	outputs map[string]any
	mutex   *sync.RWMutex
}

func newContext[T any](ctx context.Context, cache cache.Cache, data T) Context[T] {
	return Context[T]{
		Data: data,

		ctx:     ctx,
		cache:   cache,
		outputs: make(map[string]any),
		mutex:   &sync.RWMutex{},
	}
}

// Deadline implements context.Context.
func (c Context[T]) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

// Done implements context.Context.
func (c Context[T]) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err implements context.Context.
func (c Context[T]) Err() error {
	return c.ctx.Err()
}

// Value implements context.Context[T].
func (c Context[T]) Value(key any) any {
	return c.ctx.Value(key)
}

// Lock locks for writing.
func (c Context[T]) Lock() {
	c.mutex.Lock()
}

// RLock locks for reading.
func (c Context[T]) RLock() {
	c.mutex.RLock()
}

// Unlock unlocks for writing.
func (c Context[T]) Unlock() {
	c.mutex.Unlock()
}

// RUnlock undoes a single RLock call.
func (c Context[T]) RUnlock() {
	c.mutex.RUnlock()
}

func (c Context[T]) Cache() cache.Cache {
	return c.cache
}

func (c Context[T]) SetOutput(key string, value any) {
	c.outputs[key] = value
}

func (c Context[T]) GetOutput(key string) (value any, ok bool) {
	value, ok = c.outputs[key]
	return
}

func (c Context[T]) MustGetOutput(key string) (value any) {
	var ok bool
	value, ok = c.outputs[key]
	if !ok {
		panic("job is missing expected output: " + key)
	}
	return
}

type Job[T any] interface {
	Execute(ctx Context[T]) error
}

type Series[T any] []Job[T]

func (j Series[T]) Execute(ctx Context[T]) error {
	for _, job := range j {
		if err := job.Execute(ctx); err != nil {
			return err
		}
	}
	return nil
}

type Parallel[T any] []Job[T]

func (j Parallel[T]) Execute(ctx Context[T]) error {
	errs := make([]error, len(j))
	for i, job := range j {
		errs[i] = job.Execute(ctx)
	}

	var wg sync.WaitGroup
	for i, job := range j {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[i] = job.Execute(ctx)
		}()
	}

	return errors.Join(errs...)
}
