package app_driver

import (
	"context"
	"errors"

	"go-structure/internal/helper/database"
	appdrivermapper "go-structure/internal/mapper/app_driver"
	pgdb "go-structure/orm/db/postgres"
	appdrivermodel "go-structure/internal/repository/model/app_driver"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	IDriverProfileRepository interface {
		Create(ctx context.Context, accountID uuid.UUID, fullName string) (*appdrivermodel.DriverProfile, error)
		GetByAccountID(ctx context.Context, accountID uuid.UUID) (*appdrivermodel.DriverProfile, error)
		GetByID(ctx context.Context, id uuid.UUID) (*appdrivermodel.DriverProfile, error)
		List(ctx context.Context, search, globalStatus string, limit, offset int32) ([]*appdrivermodel.DriverProfile, error)
		Count(ctx context.Context, search, globalStatus string) (int64, error)
		Update(ctx context.Context, arg pgdb.UpdateDriverProfileParams) (*appdrivermodel.DriverProfile, error)
		UpdateStatus(ctx context.Context, id uuid.UUID, status pgdb.DriverProfileStatus) (*appdrivermodel.DriverProfile, error)
		CreateStatusHistory(ctx context.Context, driverID uuid.UUID, fromStatus pgdb.NullDriverProfileStatus, toStatus pgdb.DriverProfileStatus, changedBy *uuid.UUID, reason *string) error
		Delete(ctx context.Context, id uuid.UUID) error
	}

	driverProfileRepository struct {
		pool *pgxpool.Pool
	}
)

func NewDriverProfileRepository(pool *pgxpool.Pool) IDriverProfileRepository {
	return &driverProfileRepository{pool: pool}
}

func (r *driverProfileRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *driverProfileRepository) Create(ctx context.Context, accountID uuid.UUID, fullName string) (*appdrivermodel.DriverProfile, error) {
	db := r.getDB(ctx)
	arg := pgdb.CreateDriverProfileParams{
		AccountID:    accountID,
		FullName:     fullName,
		GlobalStatus: pgdb.DriverProfileStatusPENDINGVERIFICATION,
	}
	row, err := db.CreateDriverProfile(ctx, arg)
	if err != nil {
		return nil, err
	}
	return appdrivermapper.ToDriverProfileFromRow(&row), nil
}

func (r *driverProfileRepository) Update(ctx context.Context, arg pgdb.UpdateDriverProfileParams) (*appdrivermodel.DriverProfile, error) {
	db := r.getDB(ctx)
	row, err := db.UpdateDriverProfile(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return appdrivermapper.ToDriverProfileFromRow(&row), nil
}

func (r *driverProfileRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status pgdb.DriverProfileStatus) (*appdrivermodel.DriverProfile, error) {
	db := r.getDB(ctx)
	arg := pgdb.UpdateDriverProfileParams{
		GlobalStatus: pgdb.NullDriverProfileStatus{
			DriverProfileStatus: status,
			Valid:               true,
		},
		ID: id,
	}
	row, err := db.UpdateDriverProfile(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return appdrivermapper.ToDriverProfileFromRow(&row), nil
}

func (r *driverProfileRepository) CreateStatusHistory(ctx context.Context, driverID uuid.UUID, fromStatus pgdb.NullDriverProfileStatus, toStatus pgdb.DriverProfileStatus, changedBy *uuid.UUID, reason *string) error {
	db := r.getDB(ctx)
	var changedByVal pgtype.UUID
	if changedBy != nil {
		changedByVal = pgtype.UUID{Bytes: *changedBy, Valid: true}
	}
	var reasonVal pgtype.Text
	if reason != nil {
		reasonVal = pgtype.Text{String: *reason, Valid: true}
	}
	_, err := db.CreateDriverProfileStatusHistory(ctx, pgdb.CreateDriverProfileStatusHistoryParams{
		DriverID:   driverID,
		FromStatus: fromStatus,
		ToStatus:   toStatus,
		ChangedBy:  changedByVal,
		Reason:     reasonVal,
	})
	return err
}

func (r *driverProfileRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID) (*appdrivermodel.DriverProfile, error) {
	db := r.getDB(ctx)
	row, err := db.GetDriverProfileByAccountID(ctx, accountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return appdrivermapper.ToDriverProfileFromRow(&row), nil
}

func (r *driverProfileRepository) GetByID(ctx context.Context, id uuid.UUID) (*appdrivermodel.DriverProfile, error) {
	db := r.getDB(ctx)
	row, err := db.GetDriverProfileByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return appdrivermapper.ToDriverProfileFromRow(&row), nil
}

func (r *driverProfileRepository) List(ctx context.Context, search, globalStatus string, limit, offset int32) ([]*appdrivermodel.DriverProfile, error) {
	db := r.getDB(ctx)
	rows, err := db.ListDriverProfiles(ctx, pgdb.ListDriverProfilesParams{
		Column1: search,
		Column2: globalStatus,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*appdrivermodel.DriverProfile, 0, len(rows))
	for i := range rows {
		out = append(out, appdrivermapper.ToDriverProfileFromRow(&rows[i]))
	}
	return out, nil
}

func (r *driverProfileRepository) Count(ctx context.Context, search, globalStatus string) (int64, error) {
	db := r.getDB(ctx)
	return db.CountDriverProfiles(ctx, pgdb.CountDriverProfilesParams{
		Column1: search,
		Column2: globalStatus,
	})
}

func (r *driverProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.DeleteDriverProfile(ctx, id)
}
