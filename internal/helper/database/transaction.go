package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	TransactionManager interface {
		GetPool() *pgxpool.Pool
	}
	transactionManager struct {
		pool *pgxpool.Pool
	}

	contextKey string
)

const (
	txContextKey contextKey = "pgx_transaction"
)

func NewTransactionManager(pool *pgxpool.Pool) TransactionManager {
	return &transactionManager{
		pool: pool,
	}
}

func (tm *transactionManager) GetPool() *pgxpool.Pool {
	return tm.pool
}

func WithTransaction[T any](
	tm TransactionManager,
	ctx context.Context,
	fn func(ctx context.Context) (T, error),
) (T, error) {
	var zero T

	tx, err := tm.GetPool().Begin(ctx)
	if err != nil {
		return zero, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	txCtx := contextWithTransaction(ctx, tx)

	result, err := fn(txCtx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return zero, fmt.Errorf("failed to rollback transaction: %v (original error: %w)", rbErr, err)
		}
		return zero, err
	}

	if err := tx.Commit(ctx); err != nil {
		return zero, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}

func contextWithTransaction(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txContextKey, tx)
}

func TransactionFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txContextKey).(pgx.Tx)
	return tx, ok
}
