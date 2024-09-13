package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/AlexGustafsson/cupdate/internal/api"
)

func main() {
	apiServer := api.NewServer(mockAPI)

	mux := http.NewServeMux()

	mux.Handle("/api/v1/", apiServer)

	if err := http.ListenAndServe(":8080", apiServer); err != nil {
		slog.Error("Failed to serve", slog.Any("error", err))
		os.Exit(1)
	}
}
