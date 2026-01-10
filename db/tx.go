package db

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

func Tx(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error) error {
	if db == nil {
		return errors.New("db not configured")
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
