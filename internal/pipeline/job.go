package pipeline

import (
	"errors"
	"log/slog"
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
		slog.Debug("Running sequential job")
		err := job.Execute(ctx)
		if err == nil {
			slog.Debug("Sequential job completed successfully")
		} else {
			slog.Debug("Sequential job failed", slog.Any("error", err))
			return err
		}
	}
	return nil
}

type Parallel[T any] []Job[T]

func (j Parallel[T]) Execute(ctx Context[T]) error {
	var mutex sync.Mutex
	errs := make([]error, len(j))
	for i, job := range j {
		errs[i] = job.Execute(ctx)
	}

	var wg sync.WaitGroup
	for i, job := range j {
		wg.Add(1)
		go func() {
			defer wg.Done()
			slog.Debug("Running parallel job")
			err := job.Execute(ctx)
			if err == nil {
				slog.Debug("Parallel job completed successfully")
			} else {
				slog.Debug("Parallel job failed", slog.Any("error", err))
			}
			mutex.Lock()
			errs[i] = err
			defer mutex.Unlock()
		}()
	}

	return errors.Join(errs...)
}
