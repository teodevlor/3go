package database

import (
	"context"
	"database/sql"
)

type DBHelper interface {
	DB() *sql.DB
	BeginTransaction(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Close() error
}
