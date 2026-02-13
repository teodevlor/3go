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
	ILoginHistoryRepository interface {
		CreateLoginHistory(ctx context.Context, loginHistory *model.AppLoginHistory) (*model.AppLoginHistory, error)
	}

	loginHistoryRepository struct {
		pool *pgxpool.Pool
	}
)

func NewLoginHistoryRepository(pool *pgxpool.Pool) ILoginHistoryRepository {
	return &loginHistoryRepository{pool: pool}
}

func (r *loginHistoryRepository) CreateLoginHistory(ctx context.Context, loginHistory *model.AppLoginHistory) (*model.AppLoginHistory, error) {
	db := database.GetQueries(ctx, r.pool)

	params := pgdb.CreateLoginHistoryParams{
		AccountID: loginHistory.AccountID,
		DeviceID:  loginHistory.DeviceID,
		AppType:   loginHistory.AppType,
		Result:    loginHistory.Result,
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

	row, err := db.CreateLoginHistory(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapper.ToAppLoginHistory(row), nil
}
