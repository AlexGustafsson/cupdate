package workflow

import (
	"log/slog"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cache"
)

func Cache(cache cache.Cache, key string, f func(ctx Context) (any, bool)) Step {
	return StepFunc("", "Cache "+key, func(ctx Context) (map[string]any, error) {
		value, ok := f(ctx)
		if !ok {
			slog.Debug("Skipping cache as the value to cache was not found")
			return nil, nil
		}

		return nil, cache.SetJSON(ctx, key, value)
	})
}

func UnlessCached(cache cache.Cache, key string, maxAge time.Duration, step Step) Step {
	return conditionalStep{
		Step: step,
		shouldRun: func(ctx Context) (bool, error) {
			return cache.Has(ctx, key, maxAge)
		},
	}
}
