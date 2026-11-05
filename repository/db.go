package repository

import (
	"context"
	"errors"

	"inventory-tracker/apperr"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBTX is the slice of *pgxpool.Pool the repositories actually use. Depending on
// an interface rather than the concrete pool is what lets a test substitute a
// fake connection without a database.
type DBTX interface {
	Close()
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// translate converts storage errors into domain errors. Postgres codes and pgx
// sentinels stop at this boundary: no layer above the repository should need to
// import pgx to understand what went wrong.
func translate(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return apperr.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return apperr.ErrDuplicate
		case "23503": // foreign_key_violation
			return apperr.ErrReferenced
		}
	}
	return err
}
