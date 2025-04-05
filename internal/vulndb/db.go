package vulndb

import (
	"context"
	"database/sql"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/osv"
	"github.com/AlexGustafsson/cupdate/internal/semver"
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

func (c *Conn) GetGitHubAdvisoriesForRepository(ctx context.Context, repository string, version *semver.Version) ([]osv.Vulnerability, error) {
	statement, err := c.db.PrepareContext(ctx, `SELECT id, published, severity, introduced_version, fixed_version FROM github_advisories WHERE repository = ?;`)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	res, err := statement.QueryContext(ctx, repository)
	if err != nil {
		return nil, err
	}

	vulnerabilities := make([]osv.Vulnerability, 0)
	for res.Next() {
		vulnerability := osv.Vulnerability{
			DatabaseSpecific: map[string]any{},
		}

		var severity string
		var introducedVersionString string
		var fixedVersionString string
		err := res.Scan(&vulnerability.ID, &vulnerability.Published, &severity, &introducedVersionString, &fixedVersionString)
		if err != nil {
			res.Close()
			return nil, err
		}

		// TODO: Handle in database?
		introducedVersion, err := semver.ParseVersion(introducedVersionString)
		if err != nil {
			continue
		}

		// Current version is lower than introduced version
		if version.Compare(introducedVersion) < 0 {
			continue
		}

		// TODO: Handle in database?
		if fixedVersionString != "" {
			fixedVersion, err := semver.ParseVersion(fixedVersionString)
			if err != nil {
				continue
			}

			// Current version is higher than or equal to fixed version
			if version.Compare(fixedVersion) >= 0 {
				continue
			}
		}

		vulnerability.DatabaseSpecific["severity"] = severity
		if vulnerability.Published == nil {
			vulnerability.Modified = time.Now()
		} else {
			vulnerability.Modified = *vulnerability.Published
		}

		vulnerabilities = append(vulnerabilities, vulnerability)
	}
	res.Close()
	if err := res.Err(); err != nil {
		return nil, err
	}

	return vulnerabilities, nil
}

func (c *Conn) Close() error {
	return c.db.Close()
}
