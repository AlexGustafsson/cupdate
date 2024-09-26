package pipeline

import (
	"context"

	"github.com/AlexGustafsson/cupdate/internal/cache"
)

type Pipeline[T any] struct {
	cache cache.Cache
	job   Job[T]
}

func New[T any](cache cache.Cache, job Job[T]) *Pipeline[T] {
	return &Pipeline[T]{
		cache: cache,
		job:   job,
	}
}

func (p *Pipeline[T]) Run(ctx context.Context, data T) (T, error) {
	jobCtx := newContext[T](ctx, p.cache, data)
	return jobCtx.Data, p.job.Execute(jobCtx)
}

// IDEA: Auto dependencies
// Remove series, parallel. Start each job in a goroutine of its own.
// For inputs / outputs use sync.Cond so that each job can start and wait for
// their dependencies to be available. That way, the tree / flow doesn't need to
// be modeled, the pipeline author can just define all jobs, rely on the
// compiler to complain on missing outputs and then reference all jobs in the
// pipeline, one by one.
// The downside of course is the potentially high number of goroutines for
// larger jobs. But I think that realistically, we'll probably be hitting limits
// such as Docker Hub's rate limit and have to run the pipeline per image once
// every now and then anyway - so a goroutine or ten won't make that much of a
// difference. Perhaps, the issue is long-running jobs (i.e. a WaitForLimit)
// that could leave a lot of goroutines waiting for something we know will take
// a while.
// TODO: How do we handle cases when a job fails? Stop all jobs?
