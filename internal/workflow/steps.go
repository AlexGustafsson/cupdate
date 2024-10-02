package workflow

import (
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cache"
)

// Cache caches will populate the valueKey value, if it's found in cache.
// When the job is complete, the valueKey value will be cached.
// Inputs:
//   - cache
//   - cacheKey
//   - valueKey
//   - maxAge (defaults to 24h)
//
// Outputs:
//   - hit
//   - miss
func Cache[T any]() Step {
	// cache cache.Cache, cacheKey string, valueKey string, maxAge time.Duration
	return Step{
		Name: "Cache",
		// Main populates the value with the value from cache.
		Main: func(ctx Context) (Command, error) {
			cache, err := GetInput[cache.Cache](ctx, "cache", true)
			if err != nil {
				return nil, err
			}

			cacheKey, err := GetInput[string](ctx, "cacheKey", true)
			if err != nil {
				return nil, err
			}

			valueKey, err := GetInput[string](ctx, "valueKey", true)
			if err != nil {
				return nil, err
			}

			maxAge, err := GetInput[time.Duration](ctx, "maxAge", false)
			if err != nil {
				return nil, err
			}
			if maxAge == 0 {
				maxAge = 24 * time.Hour
			}

			// TODO: Problem here - GetJSON is badly implemented - Get returns nil,nil
			// on no cache found - GetJSON doesn't use that fact. Work around that for
			// now, but introduce a potential timing issue
			hit, err := cache.Has(ctx, cacheKey, maxAge)
			if err != nil {
				return nil, err
			}

			if !hit {
				return Batch(
					SetOutput("hit", false),
					SetOutput("miss", true),
				), nil
			}

			var v T
			if err := cache.GetJSON(ctx, cacheKey, &v, maxAge); err != nil {
				return nil, err
			}

			return Batch(
				SetValue(valueKey, v),
				SetOutput("hit", true),
				SetOutput("miss", false),
			), nil
		},
		Post: func(ctx Context) error {
			cache, err := GetInput[cache.Cache](ctx, "cache", true)
			if err != nil {
				return err
			}

			cacheKey, err := GetInput[string](ctx, "cacheKey", true)
			if err != nil {
				return err
			}

			valueKey, err := GetInput[string](ctx, "valueKey", true)
			if err != nil {
				return err
			}

			value, ok := ctx.Values[valueKey]
			if ok {
				return cache.SetJSON(ctx, cacheKey, value)
			}

			return nil
		},
	}
}

func StoreValue() Step {
	return Step{
		Name: "Store value",
		Main: func(ctx Context) (Command, error) {
			name, err := GetInput[string](ctx, "name", true)
			if err != nil {
				return nil, err
			}

			value, err := GetAnyInput(ctx, "value", true)
			if err != nil {
				return nil, err
			}

			return SetValue(name, value), nil
		},
	}
}
