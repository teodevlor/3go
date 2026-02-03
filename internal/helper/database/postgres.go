package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-structure/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgresDBHelper struct {
	db *sql.DB
}

func NewPostgresDB(cfg config.DatabasePostgre) (DBHelper, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &postgresDBHelper{db: db}, nil
}

func (p *postgresDBHelper) DB() *sql.DB {
	return p.db
}

func (p *postgresDBHelper) BeginTransaction(
	ctx context.Context,
	opts *sql.TxOptions,
) (*sql.Tx, error) {
	return p.db.BeginTx(ctx, opts)
}

func (p *postgresDBHelper) Close() error {
	return p.db.Close()
}
