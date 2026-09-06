package handler

import (
	"context"
	"inventory-tracker/repository"
	service "inventory-tracker/services"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	// "github.com/jackc/pgx/v5/pgxpool"
)

type DBTX interface {
	Close()
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type Handler struct {
	// DB *pgxpool.Pool //FOR DB HITS
	DB      repository.DBTX //FOR DB TESTING
	Service service.AssetService
}
