package main

import (
	"context"

	"github.com/AlexGustafsson/cupdate/internal/otelutil"
)

func ConfigureOtel(ctx context.Context, config *Config) error {
	if config.OTEL.Target != "" {
		shutdown, err := otelutil.Init(ctx, config.OTEL.Target, config.OTEL.Insecure)
		if err != nil {
			return err
		}

		// TODO: This won't be invoked on exit, make it part of the shutdown
		// procedure
		go func() {
			<-ctx.Done()
			shutdown(ctx)
		}()
	}

	return nil
}
