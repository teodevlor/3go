package database

import (
	"context"

	pgdb "go-structure/internal/orm/db/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GetQueryExecutor trả về DBTX executor phù hợp
// Nếu context có transaction thì dùng tx, không thì dùng pool
func GetQueryExecutor(ctx context.Context, pool *pgxpool.Pool) pgdb.DBTX {
	// Kiểm tra xem context có transaction không
	if tx, ok := TransactionFromContext(ctx); ok {
		return tx
	}
	// Không có transaction, dùng pool
	return pool
}

// GetQueries trả về *pgdb.Queries với executor phù hợp
// Tự động detect transaction từ context
func GetQueries(ctx context.Context, pool *pgxpool.Pool) *pgdb.Queries {
	executor := GetQueryExecutor(ctx, pool)
	return pgdb.New(executor)
}
