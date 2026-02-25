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
	IDriverDocumentTypeRepository interface {
		Create(ctx context.Context, arg pgdb.CreateDriverDocumentTypeParams) (*appdrivermodel.DriverDocumentType, error)
		GetByID(ctx context.Context, id uuid.UUID) (*appdrivermodel.DriverDocumentType, error)
		GetByCodeGlobal(ctx context.Context, code string) (*appdrivermodel.DriverDocumentType, error)
		GetByCodeAndServiceID(ctx context.Context, code string, serviceID uuid.UUID) (*appdrivermodel.DriverDocumentType, error)
		List(ctx context.Context, search string, limit, offset int32) ([]*appdrivermodel.DriverDocumentType, error)
		ListByServiceID(ctx context.Context, serviceID uuid.UUID, search string, limit, offset int32) ([]*appdrivermodel.DriverDocumentType, error)
		GetRequiredByServiceID(ctx context.Context, serviceID uuid.UUID) ([]*appdrivermodel.DriverDocumentType, error)
		Count(ctx context.Context, search string) (int64, error)
		CountByServiceID(ctx context.Context, serviceID uuid.UUID, search string) (int64, error)
		Update(ctx context.Context, arg pgdb.UpdateDriverDocumentTypeParams) (*appdrivermodel.DriverDocumentType, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	driverDocumentTypeRepository struct {
		pool *pgxpool.Pool
	}
)

func NewDriverDocumentTypeRepository(pool *pgxpool.Pool) IDriverDocumentTypeRepository {
	return &driverDocumentTypeRepository{pool: pool}
}

func (r *driverDocumentTypeRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *driverDocumentTypeRepository) Create(ctx context.Context, arg pgdb.CreateDriverDocumentTypeParams) (*appdrivermodel.DriverDocumentType, error) {
	db := r.getDB(ctx)
	row, err := db.CreateDriverDocumentType(ctx, arg)
	if err != nil {
		return nil, err
	}
	return appdrivermapper.ToDriverDocumentTypeFromRow(&row), nil
}

func (r *driverDocumentTypeRepository) GetByID(ctx context.Context, id uuid.UUID) (*appdrivermodel.DriverDocumentType, error) {
	db := r.getDB(ctx)
	row, err := db.GetDriverDocumentTypeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return appdrivermapper.ToDriverDocumentTypeFromRow(&row), nil
}

func (r *driverDocumentTypeRepository) GetByCodeGlobal(ctx context.Context, code string) (*appdrivermodel.DriverDocumentType, error) {
	db := r.getDB(ctx)
	row, err := db.GetDriverDocumentTypeByCodeGlobal(ctx, code)
	if err != nil {
		return nil, err
	}
	return appdrivermapper.ToDriverDocumentTypeFromRow(&row), nil
}

func (r *driverDocumentTypeRepository) GetByCodeAndServiceID(ctx context.Context, code string, serviceID uuid.UUID) (*appdrivermodel.DriverDocumentType, error) {
	db := r.getDB(ctx)
	arg := pgdb.GetDriverDocumentTypeByCodeAndServiceIDParams{
		Code:      code,
		ServiceID: &serviceID,
	}
	row, err := db.GetDriverDocumentTypeByCodeAndServiceID(ctx, arg)
	if err != nil {
		return nil, err
	}
	return appdrivermapper.ToDriverDocumentTypeFromRow(&row), nil
}

func (r *driverDocumentTypeRepository) List(ctx context.Context, search string, limit, offset int32) ([]*appdrivermodel.DriverDocumentType, error) {
	db := r.getDB(ctx)
	rows, err := db.ListDriverDocumentTypes(ctx, pgdb.ListDriverDocumentTypesParams{
		Column1: search,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*appdrivermodel.DriverDocumentType, 0, len(rows))
	for i := range rows {
		out = append(out, appdrivermapper.ToDriverDocumentTypeFromRow(&rows[i]))
	}
	return out, nil
}

func (r *driverDocumentTypeRepository) ListByServiceID(ctx context.Context, serviceID uuid.UUID, search string, limit, offset int32) ([]*appdrivermodel.DriverDocumentType, error) {
	db := r.getDB(ctx)
	arg := pgdb.ListDriverDocumentTypesByServiceIDParams{
		ServiceID: &serviceID,
		Column2:   search,
		Limit:     limit,
		Offset:    offset,
	}
	rows, err := db.ListDriverDocumentTypesByServiceID(ctx, arg)
	if err != nil {
		return nil, err
	}
	out := make([]*appdrivermodel.DriverDocumentType, 0, len(rows))
	for i := range rows {
		out = append(out, appdrivermapper.ToDriverDocumentTypeFromRow(&rows[i]))
	}
	return out, nil
}

func (r *driverDocumentTypeRepository) GetRequiredByServiceID(ctx context.Context, serviceID uuid.UUID) ([]*appdrivermodel.DriverDocumentType, error) {
	db := r.getDB(ctx)
	rows, err := db.ListRequiredDriverDocumentTypesByServiceID(ctx, &serviceID)
	if err != nil {
		return nil, err
	}
	out := make([]*appdrivermodel.DriverDocumentType, 0, len(rows))
	for i := range rows {
		out = append(out, appdrivermapper.ToDriverDocumentTypeFromRow(&rows[i]))
	}
	return out, nil
}

func (r *driverDocumentTypeRepository) Count(ctx context.Context, search string) (int64, error) {
	db := r.getDB(ctx)
	return db.CountDriverDocumentTypes(ctx, search)
}

func (r *driverDocumentTypeRepository) CountByServiceID(ctx context.Context, serviceID uuid.UUID, search string) (int64, error) {
	db := r.getDB(ctx)
	arg := pgdb.CountDriverDocumentTypesByServiceIDParams{
		ServiceID: &serviceID,
		Column2:   search,
	}
	return db.CountDriverDocumentTypesByServiceID(ctx, arg)
}

func (r *driverDocumentTypeRepository) Update(ctx context.Context, arg pgdb.UpdateDriverDocumentTypeParams) (*appdrivermodel.DriverDocumentType, error) {
	db := r.getDB(ctx)
	row, err := db.UpdateDriverDocumentType(ctx, arg)
	if err != nil {
		return nil, err
	}
	return appdrivermapper.ToDriverDocumentTypeFromRow(&row), nil
}

func (r *driverDocumentTypeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.DeleteDriverDocumentType(ctx, id)
}
