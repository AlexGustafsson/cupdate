package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/AlexGustafsson/cupdate/internal/api"
	"github.com/AlexGustafsson/cupdate/internal/cache"
)

func main() {
	cache, err := cache.NewDiskCache("./cache")
	if err != nil {
		slog.Error("Failed to serve", slog.Any("error", err))
		os.Exit(1)
	}

	apiServer := api.NewServer(mockAPI)

	mux := http.NewServeMux()

	mux.Handle("/api/v1/", apiServer)

	if err := http.ListenAndServe(":8080", apiServer); err != nil {
		slog.Error("Failed to serve", slog.Any("error", err))
		os.Exit(1)
	}
}
