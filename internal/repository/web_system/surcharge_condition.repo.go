package repository

import (
	"context"

	"go-structure/internal/helper/database"
	webmapper "go-structure/internal/mapper/web_system"
	pgdb "go-structure/internal/orm/db/postgres"
	websystem "go-structure/internal/repository/model/web_system"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	ISurchargeConditionRepository interface {
		Create(ctx context.Context, arg pgdb.CreateSurchargeConditionParams) (*websystem.SurchargeCondition, error)
		GetByID(ctx context.Context, id uuid.UUID) (*websystem.SurchargeCondition, error)
		List(ctx context.Context) ([]*websystem.SurchargeCondition, error)
		GetByCode(ctx context.Context, code string) (*websystem.SurchargeCondition, error)
		Update(ctx context.Context, arg pgdb.UpdateSurchargeConditionParams) (*websystem.SurchargeCondition, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	surchargeConditionRepository struct {
		pool *pgxpool.Pool
	}
)

func NewSurchargeConditionRepository(pool *pgxpool.Pool) ISurchargeConditionRepository {
	return &surchargeConditionRepository{pool: pool}
}

func (r *surchargeConditionRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *surchargeConditionRepository) Create(ctx context.Context, arg pgdb.CreateSurchargeConditionParams) (*websystem.SurchargeCondition, error) {
	db := r.getDB(ctx)
	row, err := db.CreateSurchargeCondition(ctx, arg)
	if err != nil {
		return nil, err
	}
	return webmapper.ToSurchargeConditionFromRow(&row), nil
}

func (r *surchargeConditionRepository) GetByID(ctx context.Context, id uuid.UUID) (*websystem.SurchargeCondition, error) {
	db := r.getDB(ctx)
	row, err := db.GetSurchargeConditionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return webmapper.ToSurchargeConditionFromRow(&row), nil
}

func (r *surchargeConditionRepository) List(ctx context.Context) ([]*websystem.SurchargeCondition, error) {
	db := r.getDB(ctx)
	rows, err := db.ListSurchargeConditions(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*websystem.SurchargeCondition, 0, len(rows))
	for i := range rows {
		out = append(out, webmapper.ToSurchargeConditionFromRow(&rows[i]))
	}
	return out, nil
}

func (r *surchargeConditionRepository) Update(ctx context.Context, arg pgdb.UpdateSurchargeConditionParams) (*websystem.SurchargeCondition, error) {
	db := r.getDB(ctx)
	row, err := db.UpdateSurchargeCondition(ctx, arg)
	if err != nil {
		return nil, err
	}
	return webmapper.ToSurchargeConditionFromRow(&row), nil
}

func (r *surchargeConditionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.DeleteSurchargeCondition(ctx, id)
}

func (r *surchargeConditionRepository) GetByCode(ctx context.Context, code string) (*websystem.SurchargeCondition, error) {
	db := r.getDB(ctx)
	row, err := db.GetSurchargeConditionByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return webmapper.ToSurchargeConditionFromRow(&row), nil
}
