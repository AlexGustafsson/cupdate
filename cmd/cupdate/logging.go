package main

import (
	"log/slog"
	"os"

	"github.com/AlexGustafsson/cupdate/internal/slogutil"
)

func ConfigureLogging(config *Config) {
	logLevel := slog.LevelError
	if config != nil {
		logLevel = config.LogLevel()
	}

	handler := slogutil.NewHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: logLevel},
	)

	logger := slog.New(handler).
		With(slog.String("service.version", Version)).
		With(slog.String("service.name", "cupdate"))

	slog.SetDefault(logger)
}
