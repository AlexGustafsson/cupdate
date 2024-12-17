package db

import (
	"context"
	"database/sql"
	_ "embed" // Embed SQL files
	"net/url"
	"strings"

	"github.com/AlexGustafsson/cupdate/tools/vulndb/internal/ossf"
	_ "modernc.org/sqlite"
)

//go:embed createTablesIfNotExist.sql
var createTablesIfNotExist string

type Conn struct {
	db *sql.DB
}

func Open(path string) (*Conn, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createTablesIfNotExist)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return &Conn{db: db}, nil
}

func (c *Conn) Insert(ctx context.Context, vuln ossf.OpenSourceVulnerability) error {
	statement, err := c.db.PrepareContext(ctx, `INSERT INTO github_advisories (id, repository, published, severity) VALUES (?, ?, ?, ?) ON CONFLICT DO NOTHING;`)
	if err != nil {
		return err
	}

	repository := ""
	for _, reference := range vuln.References {
		u, err := url.Parse(reference.Url)
		if err == nil {
			segments := len(u.Path) - len(strings.ReplaceAll(u.Path, "/", ""))
			if u.Host == "github.com" && segments == 2 {
				repository = reference.Url
				break
			}
		}
	}

	// No repository found
	if repository == "" {
		return nil
	}

	severity := ""
	if value, ok := vuln.DatabaseSpecific["severity"]; ok {
		severity = value.(string)
	}

	// TODO: Insert ranges, or duplicate for each range

	_, err = statement.ExecContext(ctx, vuln.ID, repository, vuln.Published, severity)
	if err != nil {
		return err
	}

	return nil
}

func (c *Conn) Close() error {
	return c.db.Close()
}
