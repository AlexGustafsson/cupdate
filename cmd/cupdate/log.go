package main

import (
	"log/slog"
	"os"

	"github.com/AlexGustafsson/cupdate/internal/slogutil"
)

func InitDefaultLogger() {
	options := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slogutil.NewHandler(os.Stderr, options)

	logger := slog.New(handler).
		With(slog.String("service.version", Version)).
		With(slog.String("service.name", "cupdate"))

	slog.SetDefault(logger)
}
