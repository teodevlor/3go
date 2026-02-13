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
	IServiceRepository interface {
		CreateService(ctx context.Context, arg pgdb.CreateServiceParams) (*websystem.Service, error)
		GetServiceByID(ctx context.Context, id uuid.UUID) (*websystem.Service, error)
		GetServiceByCode(ctx context.Context, code string) (*websystem.Service, error)
		ListServices(ctx context.Context, search string, limit, offset int32) ([]*websystem.Service, error)
		CountServices(ctx context.Context, search string) (int64, error)
		UpdateService(ctx context.Context, arg pgdb.UpdateServiceParams) (*websystem.Service, error)
		DeleteService(ctx context.Context, id uuid.UUID) error
	}

	serviceRepository struct {
		pool *pgxpool.Pool
	}
)

func NewServiceRepository(pool *pgxpool.Pool) IServiceRepository {
	return &serviceRepository{pool: pool}
}

func (sr *serviceRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, sr.pool)
}

func (sr *serviceRepository) CreateService(ctx context.Context, arg pgdb.CreateServiceParams) (*websystem.Service, error) {
	db := sr.getDB(ctx)
	row, err := db.CreateService(ctx, arg)
	if err != nil {
		return nil, err
	}
	return webmapper.ToServiceFromRow(&row), nil
}

func (sr *serviceRepository) GetServiceByID(ctx context.Context, id uuid.UUID) (*websystem.Service, error) {
	db := sr.getDB(ctx)
	row, err := db.GetServiceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return webmapper.ToServiceFromRow(&row), nil
}

func (sr *serviceRepository) GetServiceByCode(ctx context.Context, code string) (*websystem.Service, error) {
	db := sr.getDB(ctx)
	row, err := db.GetServiceByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return webmapper.ToServiceFromRow(&row), nil
}

func (sr *serviceRepository) ListServices(ctx context.Context, search string, limit, offset int32) ([]*websystem.Service, error) {
	db := sr.getDB(ctx)
	rows, err := db.ListServices(ctx, pgdb.ListServicesParams{Column1: search, Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]*websystem.Service, 0, len(rows))
	for i := range rows {
		out = append(out, webmapper.ToServiceFromRow(&rows[i]))
	}
	return out, nil
}

func (sr *serviceRepository) CountServices(ctx context.Context, search string) (int64, error) {
	db := sr.getDB(ctx)
	return db.CountServices(ctx, search)
}

func (sr *serviceRepository) UpdateService(ctx context.Context, arg pgdb.UpdateServiceParams) (*websystem.Service, error) {
	db := sr.getDB(ctx)
	row, err := db.UpdateService(ctx, arg)
	if err != nil {
		return nil, err
	}
	return webmapper.ToServiceFromRow(&row), nil
}

func (sr *serviceRepository) DeleteService(ctx context.Context, id uuid.UUID) error {
	db := sr.getDB(ctx)
	return db.DeleteService(ctx, id)
}
