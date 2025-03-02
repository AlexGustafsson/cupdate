package vulndb

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

type Conn struct {
	db *sql.DB
}

func Open(uri string) (*Conn, error) {
	uri += "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(1000)&_time_format=sqlite&_pragma=query_only(true)"

	db, err := sql.Open("sqlite", uri)
	if err != nil {
		return nil, err
	}

	return &Conn{db: db}, nil
}

func (c *Conn) GetGitHubAdvisoriesForRepository(ctx context.Context, repository string) ([]GitHubAdvisory, error) {
	statement, err := c.db.PrepareContext(ctx, `SELECT id, repository, published, severity, introduced_version, fixed_version FROM github_advisories WHERE repository = ?;`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	res, err := statement.QueryContext(ctx, repository)
	if err != nil {
		return nil, err
	}

	advisories := make([]GitHubAdvisory, 0)
	for res.Next() {
		var advisory GitHubAdvisory
		var fixedVersion *string
		err := res.Scan(&advisory.ID, &advisory.Repository, &advisory.Published, &advisory.Severity, &advisory.IntroducedVersion, &fixedVersion)
		if err != nil {
			res.Close()
			return nil, err
		}
		if fixedVersion != nil {
			advisory.FixedVersion = *fixedVersion
		}
		switch advisory.Severity {
		case "CRITICAL":
			advisory.Severity = SeverityCritical
		case "HIGH":
			advisory.Severity = SeverityHigh
		case "MODERATE":
			advisory.Severity = SeverityMedium
		case "LOW":
			advisory.Severity = SeverityLow
		default:
			advisory.Severity = SeverityUnspecified
		}
		advisories = append(advisories, advisory)
	}
	res.Close()
	if err := res.Err(); err != nil {
		return nil, err
	}

	return advisories, nil
}

func (c *Conn) Close() error {
	return c.db.Close()
}
