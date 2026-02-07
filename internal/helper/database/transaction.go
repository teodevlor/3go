package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionManager quản lý transactions
type TransactionManager interface {
	// WithTransaction thực thi function trong một transaction
	// Tự động commit nếu success, rollback nếu error hoặc panic
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	
	// GetPool trả về connection pool (để dùng cho non-transaction queries)
	GetPool() *pgxpool.Pool
}

// transactionManager implementation cho pgx
type transactionManager struct {
	pool *pgxpool.Pool
}

// NewTransactionManager tạo transaction manager mới
func NewTransactionManager(pool *pgxpool.Pool) TransactionManager {
	return &transactionManager{
		pool: pool,
	}
}

// WithTransaction thực thi function trong transaction với auto commit/rollback
func (tm *transactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// Bắt đầu transaction
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer để handle rollback trong mọi trường hợp
	defer func() {
		if p := recover(); p != nil {
			// Nếu có panic, rollback và re-panic
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	// Lưu transaction vào context để các repository có thể dùng
	txCtx := contextWithTransaction(ctx, tx)

	// Thực thi function
	if err := fn(txCtx); err != nil {
		// Nếu có error, rollback
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	// Nếu success, commit
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetPool trả về connection pool
func (tm *transactionManager) GetPool() *pgxpool.Pool {
	return tm.pool
}

// Context keys để lưu transaction
type contextKey string

const txContextKey contextKey = "pgx_transaction"

// contextWithTransaction thêm transaction vào context
func contextWithTransaction(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txContextKey, tx)
}

// TransactionFromContext lấy transaction từ context (nếu có)
func TransactionFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txContextKey).(pgx.Tx)
	return tx, ok
}
