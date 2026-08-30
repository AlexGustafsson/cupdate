package sql

import (
	"context"
	"database/sql"
)

func Scan[T Unmarshaler](ctx context.Context, db *sql.DB, query string, args ...any) ([]T, error) {
	statement, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	res, err := statement.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	rows := make([]T, 0)
	for res.Next() {
		var row T
		if err := res.Scan(&rows); err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}
	if err := res.Err(); err != nil {
		return nil, err
	}

	return rows, nil
}
