package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-structure/config"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlDBHelper struct {
	db *sql.DB
}

func NewMySQLDB(cfg config.DatabaseMySQL) (DBHelper, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
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

	return &mysqlDBHelper{db: db}, nil
}

func (m *mysqlDBHelper) DB() *sql.DB {
	return m.db
}

func (m *mysqlDBHelper) BeginTransaction(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return m.db.BeginTx(ctx, opts)
}

func (m *mysqlDBHelper) Close() error {
	return m.db.Close()
}
