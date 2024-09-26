package pipeline

import (
	"errors"
	"sync"
)

type Job[T any] interface {
	Execute(ctx Context[T]) error
}

type JobFunc[T any] func(ctx Context[T]) error

func (f JobFunc[T]) Execute(ctx Context[T]) error {
	return f(ctx)
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
