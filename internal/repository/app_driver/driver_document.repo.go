package app_driver

import (
	"context"

	"go-structure/internal/helper/database"
	appdrivermapper "go-structure/internal/mapper/app_driver"
	pgdb "go-structure/internal/orm/db/postgres"
	appdrivermodel "go-structure/internal/repository/model/app_driver"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	IDriverDocumentRepository interface {
		Create(ctx context.Context, arg pgdb.CreateDriverDocumentParams) (*appdrivermodel.DriverDocument, error)
		GetByID(ctx context.Context, id uuid.UUID) (*appdrivermodel.DriverDocument, error)
		ListByDriverID(ctx context.Context, driverID uuid.UUID) ([]*appdrivermodel.DriverDocument, error)
		Update(ctx context.Context, arg pgdb.UpdateDriverDocumentParams) (*appdrivermodel.DriverDocument, error)
		UpdatePartial(ctx context.Context, arg pgdb.UpdateDriverDocumentPartialParams) (*appdrivermodel.DriverDocument, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	driverDocumentRepository struct {
		pool *pgxpool.Pool
	}
)

func NewDriverDocumentRepository(pool *pgxpool.Pool) IDriverDocumentRepository {
	return &driverDocumentRepository{pool: pool}
}

func (r *driverDocumentRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *driverDocumentRepository) Create(ctx context.Context, arg pgdb.CreateDriverDocumentParams) (*appdrivermodel.DriverDocument, error) {
	db := r.getDB(ctx)
	row, err := db.CreateDriverDocument(ctx, arg)
	if err != nil {
		return nil, err
	}
	return appdrivermapper.ToDriverDocumentFromRow(&row), nil
}

func (r *driverDocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*appdrivermodel.DriverDocument, error) {
	db := r.getDB(ctx)
	row, err := db.GetDriverDocumentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return appdrivermapper.ToDriverDocumentFromRow(&row), nil
}

func (r *driverDocumentRepository) ListByDriverID(ctx context.Context, driverID uuid.UUID) ([]*appdrivermodel.DriverDocument, error) {
	db := r.getDB(ctx)
	rows, err := db.ListDriverDocumentsByDriverID(ctx, driverID)
	if err != nil {
		return nil, err
	}
	out := make([]*appdrivermodel.DriverDocument, 0, len(rows))
	for i := range rows {
		out = append(out, appdrivermapper.ToDriverDocumentFromRow(&rows[i]))
	}
	return out, nil
}

func (r *driverDocumentRepository) Update(ctx context.Context, arg pgdb.UpdateDriverDocumentParams) (*appdrivermodel.DriverDocument, error) {
	db := r.getDB(ctx)
	row, err := db.UpdateDriverDocument(ctx, arg)
	if err != nil {
		return nil, err
	}
	return appdrivermapper.ToDriverDocumentFromRow(&row), nil
}

func (r *driverDocumentRepository) UpdatePartial(ctx context.Context, arg pgdb.UpdateDriverDocumentPartialParams) (*appdrivermodel.DriverDocument, error) {
	db := r.getDB(ctx)
	row, err := db.UpdateDriverDocumentPartial(ctx, arg)
	if err != nil {
		return nil, err
	}
	return appdrivermapper.ToDriverDocumentFromRow(&row), nil
}

func (r *driverDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.DeleteDriverDocument(ctx, id)
}
