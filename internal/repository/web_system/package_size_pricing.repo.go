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
	IPackageSizePricingRepository interface {
		Create(ctx context.Context, arg pgdb.CreatePackageSizePricingParams) (*websystem.PackageSizePricing, error)
		GetByID(ctx context.Context, id uuid.UUID) (*websystem.PackageSizePricing, error)
		List(ctx context.Context, serviceID *uuid.UUID) ([]*websystem.PackageSizePricing, error)
		Update(ctx context.Context, arg pgdb.UpdatePackageSizePricingParams) (*websystem.PackageSizePricing, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	packageSizePricingRepository struct {
		pool *pgxpool.Pool
	}
)

func NewPackageSizePricingRepository(pool *pgxpool.Pool) IPackageSizePricingRepository {
	return &packageSizePricingRepository{pool: pool}
}

func (r *packageSizePricingRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *packageSizePricingRepository) Create(ctx context.Context, arg pgdb.CreatePackageSizePricingParams) (*websystem.PackageSizePricing, error) {
	db := r.getDB(ctx)
	row, err := db.CreatePackageSizePricing(ctx, arg)
	if err != nil {
		return nil, err
	}
	return webmapper.ToPackageSizePricingFromRow(&row), nil
}

func (r *packageSizePricingRepository) GetByID(ctx context.Context, id uuid.UUID) (*websystem.PackageSizePricing, error) {
	db := r.getDB(ctx)
	row, err := db.GetPackageSizePricingByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return webmapper.ToPackageSizePricingFromRow(&row), nil
}

func (r *packageSizePricingRepository) List(ctx context.Context, serviceID *uuid.UUID) ([]*websystem.PackageSizePricing, error) {
	db := r.getDB(ctx)
	var filterID pgtype.UUID
	if serviceID != nil {
		filterID = pgtype.UUID{Bytes: *serviceID, Valid: true}
	}
	rows, err := db.ListPackageSizePricings(ctx, filterID)
	if err != nil {
		return nil, err
	}
	out := make([]*websystem.PackageSizePricing, 0, len(rows))
	for i := range rows {
		out = append(out, webmapper.ToPackageSizePricingFromRow(&rows[i]))
	}
	return out, nil
}

func (r *packageSizePricingRepository) Update(ctx context.Context, arg pgdb.UpdatePackageSizePricingParams) (*websystem.PackageSizePricing, error) {
	db := r.getDB(ctx)
	row, err := db.UpdatePackageSizePricing(ctx, arg)
	if err != nil {
		return nil, err
	}
	return webmapper.ToPackageSizePricingFromRow(&row), nil
}

func (r *packageSizePricingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.DeletePackageSizePricing(ctx, id)
}
