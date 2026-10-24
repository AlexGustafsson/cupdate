package main

import (
	"context"

	"github.com/AlexGustafsson/cupdate/internal/otelutil"
)

func ConfigureOtel(ctx context.Context, config *Config) (func(context.Context) error, error) {
	if config.OTEL.Target != "" {
		return otelutil.Init(ctx, config.OTEL.Target, config.OTEL.Insecure)
	}

	return func(ctx context.Context) error { return nil }, nil
}
