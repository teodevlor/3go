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
	ISurchargeRuleRepository interface {
		Create(ctx context.Context, arg pgdb.CreateSurchargeRuleParams) (*websystem.SurchargeRule, error)
		GetByID(ctx context.Context, id uuid.UUID) (*websystem.SurchargeRule, error)
		List(ctx context.Context, serviceID, zoneID *uuid.UUID) ([]*websystem.SurchargeRule, error)
		Update(ctx context.Context, arg pgdb.UpdateSurchargeRuleParams) (*websystem.SurchargeRule, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	surchargeRuleRepository struct {
		pool *pgxpool.Pool
	}
)

func NewSurchargeRuleRepository(pool *pgxpool.Pool) ISurchargeRuleRepository {
	return &surchargeRuleRepository{pool: pool}
}

func (r *surchargeRuleRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *surchargeRuleRepository) Create(ctx context.Context, arg pgdb.CreateSurchargeRuleParams) (*websystem.SurchargeRule, error) {
	db := r.getDB(ctx)
	row, err := db.CreateSurchargeRule(ctx, arg)
	if err != nil {
		return nil, err
	}
	return webmapper.ToSurchargeRuleFromRow(&row), nil
}

func (r *surchargeRuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*websystem.SurchargeRule, error) {
	db := r.getDB(ctx)
	row, err := db.GetSurchargeRuleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return webmapper.ToSurchargeRuleFromRow(&row), nil
}

func (r *surchargeRuleRepository) List(ctx context.Context, serviceID, zoneID *uuid.UUID) ([]*websystem.SurchargeRule, error) {
	db := r.getDB(ctx)
	var svcID, zID pgtype.UUID
	if serviceID != nil {
		svcID = pgtype.UUID{Bytes: *serviceID, Valid: true}
	}
	if zoneID != nil {
		zID = pgtype.UUID{Bytes: *zoneID, Valid: true}
	}
	rows, err := db.ListSurchargeRules(ctx, pgdb.ListSurchargeRulesParams{ServiceID: svcID, ZoneID: zID})
	if err != nil {
		return nil, err
	}
	out := make([]*websystem.SurchargeRule, 0, len(rows))
	for i := range rows {
		out = append(out, webmapper.ToSurchargeRuleFromRow(&rows[i]))
	}
	return out, nil
}

func (r *surchargeRuleRepository) Update(ctx context.Context, arg pgdb.UpdateSurchargeRuleParams) (*websystem.SurchargeRule, error) {
	db := r.getDB(ctx)
	row, err := db.UpdateSurchargeRule(ctx, arg)
	if err != nil {
		return nil, err
	}
	return webmapper.ToSurchargeRuleFromRow(&row), nil
}

func (r *surchargeRuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.DeleteSurchargeRule(ctx, id)
}
