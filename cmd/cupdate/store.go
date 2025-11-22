package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/store"
)

func HandleWorkflowCleanup(ctx context.Context, config *Config, writeStore *store.Store) {
	ticker := time.NewTicker(config.Workflow.CleanupInterval)
	defer ticker.Stop()
	defer slog.Debug("Closed workflow cleanup")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			slog.Debug("Cleaning up old workflow runs")
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			removed, err := writeStore.DeleteWorkflowRuns(ctx, time.Now().Add(-config.Workflow.CleanupMaxAge))
			cancel()
			if err == nil {
				slog.Debug("Cleaned up old workflow runs successfully", slog.Int64("removed", removed))
			} else {
				slog.Error("Failed to clean up old workflow runs", slog.Any("error", err))
			}
		}
	}
}
