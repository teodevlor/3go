package repository

import (
	"context"

	"go-structure/internal/helper/database"
	"go-structure/internal/mapper"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	ISystemLoginHistoryRepository interface {
		CreateSystemLoginHistory(ctx context.Context, loginHistory *model.SystemLoginHistory) (*model.SystemLoginHistory, error)
	}

	systemLoginHistoryRepository struct {
		pool *pgxpool.Pool
	}
)

func NewSystemLoginHistoryRepository(pool *pgxpool.Pool) ISystemLoginHistoryRepository {
	return &systemLoginHistoryRepository{pool: pool}
}

func (r *systemLoginHistoryRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *systemLoginHistoryRepository) CreateSystemLoginHistory(ctx context.Context, loginHistory *model.SystemLoginHistory) (*model.SystemLoginHistory, error) {
	db := r.getDB(ctx)

	params := pgdb.CreateSystemLoginHistoryParams{
		AdminID: loginHistory.AdminID,
		Result:  loginHistory.Result,
		FailureReason: pgtype.Text{
			String: loginHistory.FailureReason,
			Valid:  loginHistory.FailureReason != "",
		},
		IpAddress: pgtype.Text{
			String: loginHistory.IpAddress,
			Valid:  loginHistory.IpAddress != "",
		},
		UserAgent: pgtype.Text{
			String: loginHistory.UserAgent,
			Valid:  loginHistory.UserAgent != "",
		},
		Location: pgtype.Text{
			String: loginHistory.Location,
			Valid:  loginHistory.Location != "",
		},
		Metadata: loginHistory.Metadata,
	}

	row, err := db.CreateSystemLoginHistory(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapper.ToSystemLoginHistory(row), nil
}
