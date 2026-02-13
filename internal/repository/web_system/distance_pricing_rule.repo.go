package repository

import (
	"context"

	"go-structure/internal/helper/database"
	webmapper "go-structure/internal/mapper/web_system"
	pgdb "go-structure/internal/orm/db/postgres"
	websystem "go-structure/internal/repository/model/web_system"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	IDistancePricingRuleRepository interface {
		Create(ctx context.Context, arg pgdb.CreateDistancePricingRuleParams) (*websystem.DistancePricingRule, error)
		GetByID(ctx context.Context, id uuid.UUID) (*websystem.DistancePricingRule, error)
		List(ctx context.Context, serviceID *uuid.UUID) ([]*websystem.DistancePricingRule, error)
		Update(ctx context.Context, arg pgdb.UpdateDistancePricingRuleParams) (*websystem.DistancePricingRule, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	distancePricingRuleRepository struct {
		pool *pgxpool.Pool
	}
)

func NewDistancePricingRuleRepository(pool *pgxpool.Pool) IDistancePricingRuleRepository {
	return &distancePricingRuleRepository{pool: pool}
}

func (dpr *distancePricingRuleRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, dpr.pool)
}

func (dpr *distancePricingRuleRepository) Create(ctx context.Context, arg pgdb.CreateDistancePricingRuleParams) (*websystem.DistancePricingRule, error) {
	db := dpr.getDB(ctx)
	row, err := db.CreateDistancePricingRule(ctx, arg)
	if err != nil {
		return nil, err
	}
	return webmapper.ToDistancePricingRuleFromRow(&row), nil
}

func (dpr *distancePricingRuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*websystem.DistancePricingRule, error) {
	db := dpr.getDB(ctx)
	row, err := db.GetDistancePricingRuleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return webmapper.ToDistancePricingRuleFromRow(&row), nil
}

func (dpr *distancePricingRuleRepository) List(ctx context.Context, serviceID *uuid.UUID) ([]*websystem.DistancePricingRule, error) {
	db := dpr.getDB(ctx)
	var filterID pgtype.UUID
	if serviceID != nil {
		filterID = pgtype.UUID{
			Bytes: *serviceID,
			Valid: true,
		}
	}
	rows, err := db.ListDistancePricingRules(ctx, filterID)
	if err != nil {
		return nil, err
	}
	out := make([]*websystem.DistancePricingRule, 0, len(rows))
	for i := range rows {
		out = append(out, webmapper.ToDistancePricingRuleFromRow(&rows[i]))
	}
	return out, nil
}

func (dpr *distancePricingRuleRepository) Update(ctx context.Context, arg pgdb.UpdateDistancePricingRuleParams) (*websystem.DistancePricingRule, error) {
	db := dpr.getDB(ctx)
	row, err := db.UpdateDistancePricingRule(ctx, arg)
	if err != nil {
		return nil, err
	}
	return webmapper.ToDistancePricingRuleFromRow(&row), nil
}

func (dpr *distancePricingRuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := dpr.getDB(ctx)
	return db.DeleteDistancePricingRule(ctx, id)
}
