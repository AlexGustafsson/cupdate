package sql

import (
	"context"
	"database/sql"
)

type CRUDDB[T any] struct {
	db    *sql.DB
	table string
}

func CRUD[T any](db *sql.DB, table string) *CRUDDB[T] {
	// Identify primary key, panic on error

	return &CRUDDB[T]{
		db:    db,
		table: table,
	}
}

func (d *CRUDDB[T]) Create(ctx context.Context) error {

}

func (d *CRUDDB[T]) Read(ctx context.Context) (*T, error) {

}

func (d *CRUDDB[T]) Update(ctx context.Context) error {

}

func (d *CRUDDB[T]) Delete(ctx context.Context, id string) error {
	statement, err := d.db.PrepareContext(ctx, `DELETE FROM ? WHERE ? = ?;`)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.ExecContext(ctx, id)
	return err
}
