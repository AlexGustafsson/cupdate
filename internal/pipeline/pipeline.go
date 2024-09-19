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
