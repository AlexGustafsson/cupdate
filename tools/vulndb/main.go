package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/AlexGustafsson/cupdate/internal/osv"
	"github.com/AlexGustafsson/cupdate/tools/vulndb/internal/db"
	"github.com/AlexGustafsson/cupdate/tools/vulndb/internal/git"
	"github.com/AlexGustafsson/cupdate/tools/vulndb/internal/oci"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	slog.Info("Starting vulndb collector")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		<-signals
		slog.InfoContext(ctx, "Caught signal, exiting gracefully")
		cancel()
	}()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "Fatal error", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	githubToken := os.Getenv("INPUT_GITHUB_TOKEN")
	githubActor := os.Getenv("INPUT_GITHUB_ACTOR")

	if githubToken == "" || githubActor == "" {
		return fmt.Errorf("missing required input(s)")
	}

	workdirParent, err := os.MkdirTemp(os.TempDir(), "cupdate-vulndb-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(workdirParent)

	workdir := filepath.Join(workdirParent, "advisory-database")

	slog.DebugContext(ctx, "Performing shallow clone of GitHub's advisory database")
	err = git.ShallowClone(context.Background(), "https://github.com/github/advisory-database", workdir, "advisories/github-reviewed/2024", "advisories/github-reviewed/2025")
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	slog.DebugContext(ctx, "Creating database")
	db, err := db.Open("vulndb.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	slog.DebugContext(ctx, "Inserting advisories in database")
	err = filepath.WalkDir(workdir, func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) == ".json" {
			file, err := os.Open(path)
			if err != nil {
				return err
			}

			var vuln osv.Vulnerability
			if err := json.NewDecoder(file).Decode(&vuln); err != nil {
				return err
			}

			return db.Insert(ctx, vuln)
		}

		return nil
	})

	if err := db.Close(); err != nil {
		slog.ErrorContext(ctx, "Failed to close database", slog.Any("error", db))
	}

	if err != nil {
		return err
	}

	slog.DebugContext(ctx, "Pushing artifact")
	if err := oci.PushArtifact(ctx, "vulndb.sqlite", githubActor, githubToken); err != nil {
		return err
	}

	slog.InfoContext(ctx, "Successfully pushed artifact")
	return nil
}
