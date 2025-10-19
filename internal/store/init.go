package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"slices"

	_ "embed" // Embed SQL files
)

const Revision = 6

//go:embed migrations
var migrations embed.FS

//go:embed init.sql
var initsql string

// Initialize initializes a database, migrating and creating tables as
// necessary.
func Initialize(ctx context.Context, uri string) error {
	// Use WAL to allow multiple readers
	uri += "?_pragma=journal_mode(WAL)&_time_format=sqlite"

	db, err := sql.Open("sqlite", uri)
	if err != nil {
		return err
	}
	defer db.Close()

	tables, err := getTables(ctx, db)
	if err != nil {
		return err
	}

	if len(tables) == 0 {
		slog.InfoContext(ctx, "Initializing store")
		if err := initialize(ctx, db); err != nil {
			slog.ErrorContext(ctx, "Failed to initialize store", slog.Any("error", err))
			return err
		}
		slog.DebugContext(ctx, "Successfully initialized store")
	}

	sourceRevision := 0
	// TODO: In v1, remove this check as we can assume the revision table is
	// part of the core set of tables
	if slices.Contains(tables, "revision") {
		sourceRevision, err = getStoreRevision(ctx, db)
		if err != nil {
			return err
		}
	}

	slog.InfoContext(ctx, "Migrating store", slog.Int("sourceRevision", sourceRevision), slog.Int("targetRevision", Revision))
	if err := migrate(ctx, db, sourceRevision, Revision); err != nil {
		slog.ErrorContext(ctx, "Failed to migrate store", slog.Any("error", err))
		return err
	}
	slog.DebugContext(ctx, "Successfully migrated store", slog.Int("sourceRevision", sourceRevision), slog.Int("targetRevision", Revision))

	return nil
}

func getTables(ctx context.Context, db *sql.DB) ([]string, error) {
	res, err := db.QueryContext(ctx, "SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	tables := make([]string, 0)
	for res.Next() {
		var name string
		if err := res.Scan(&name); err != nil {
			return nil, err
		}

		tables = append(tables, name)
	}
	if err := res.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func getStoreRevision(ctx context.Context, db *sql.DB) (int, error) {
	res, err := db.QueryContext(ctx, "SELECT revision FROM revision;")
	if err != nil {
		return 0, err
	}
	defer res.Close()

	found := res.Next()
	if err := res.Err(); err != nil {
		return 0, err
	}
	if !found {
		return 0, nil
	}

	var revision int
	if err := res.Scan(&revision); err != nil {
		return 0, err
	}

	return revision, nil
}

func migrate(ctx context.Context, db *sql.DB, sourceRevision int, targetRevision int) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i := sourceRevision; i < targetRevision; i++ {
		log := slog.With(slog.Int("sourceRevision", i), slog.Int("targetRevision", i+1))

		log.DebugContext(ctx, "Running migration")
		script, err := migrations.ReadFile(fmt.Sprintf("migrations/%d.sql", i))
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, string(script))
		if err != nil {
			return err
		}

		log.DebugContext(ctx, "Successfully ran migration")
	}

	return tx.Commit()
}

func initialize(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	slog.DebugContext(ctx, "Running init script")

	_, err = tx.ExecContext(ctx, initsql)
	if err != nil {
		return err
	}

	slog.DebugContext(ctx, "Successfully ran init script")

	return tx.Commit()
}
