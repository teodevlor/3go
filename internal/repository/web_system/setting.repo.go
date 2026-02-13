package repository

import (
	"context"
	"errors"

	"go-structure/internal/helper/database"
	"go-structure/internal/mapper"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrSettingNotFound = errors.New("không tìm thấy cấu hình")
)

type (
	ISettingRepository interface {
		GetSettingByKey(ctx context.Context, key string) (*model.Setting, error)
	}

	settingRepository struct {
		pool *pgxpool.Pool
	}
)

func NewSettingRepository(pool *pgxpool.Pool) ISettingRepository {
	return &settingRepository{pool: pool}
}

func (r *settingRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *settingRepository) GetSettingByKey(ctx context.Context, key string) (*model.Setting, error) {
	db := r.getDB(ctx)
	row, err := db.GetSettingByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSettingNotFound
		}
		return nil, err
	}
	return mapper.ToSetting(row), nil
}
