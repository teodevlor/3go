package repository

import (
	"context"

	"go-structure/internal/helper/database"
	"go-structure/internal/mapper"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IDeviceRepository interface {
	CreateDevice(ctx context.Context, device *model.Device) (*model.Device, error)
	GetDeviceByID(ctx context.Context, id string) (*model.Device, error)
	GetDeviceByUID(ctx context.Context, uid string) (*model.Device, error)
	UpdateDevice(ctx context.Context, device *model.Device) (*model.Device, error)
	DeleteDevice(ctx context.Context, id string) error
}

type deviceRepository struct {
	pool *pgxpool.Pool
}

func NewDeviceRepository(pool *pgxpool.Pool) IDeviceRepository {
	return &deviceRepository{pool: pool}
}

func (r *deviceRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *deviceRepository) CreateDevice(ctx context.Context, device *model.Device) (*model.Device, error) {
	db := r.getDB(ctx)

	row, err := db.CreateDevice(ctx, pgdb.CreateDeviceParams{
		DeviceUid: device.DeviceUID,
		Platform:  device.Platform,
		DeviceName: pgtype.Text{
			String: device.DeviceName,
			Valid:  device.DeviceName != "",
		},
		OsVersion: pgtype.Text{
			String: device.OsVersion,
			Valid:  device.OsVersion != "",
		},
		AppVersion: pgtype.Text{
			String: device.AppVersion,
			Valid:  device.AppVersion != "",
		},
		Metadata: device.Metadata,
	})
	if err != nil {
		return nil, err
	}

	return mapper.ToDevice(row), nil
}

func (r *deviceRepository) GetDeviceByID(ctx context.Context, id string) (*model.Device, error) {
	db := r.getDB(ctx)
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	row, err := db.GetDeviceByID(ctx, uuidID)
	if err != nil {
		return nil, err
	}

	return mapper.ToDevice(row), nil
}

func (r *deviceRepository) GetDeviceByUID(ctx context.Context, uid string) (*model.Device, error) {
	db := r.getDB(ctx)
	row, err := db.GetDeviceByUID(ctx, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToDevice(row), nil
}

func (r *deviceRepository) UpdateDevice(ctx context.Context, device *model.Device) (*model.Device, error) {
	db := r.getDB(ctx)

	row, err := db.UpdateDevice(ctx, pgdb.UpdateDeviceParams{
		ID: device.ID,
		DeviceName: pgtype.Text{
			String: device.DeviceName,
			Valid:  device.DeviceName != "",
		},
		OsVersion: pgtype.Text{
			String: device.OsVersion,
			Valid:  device.OsVersion != "",
		},
		AppVersion: pgtype.Text{
			String: device.AppVersion,
			Valid:  device.AppVersion != "",
		},
		Metadata: device.Metadata,
	})
	if err != nil {
		return nil, err
	}

	return mapper.ToDevice(row), nil
}

func (r *deviceRepository) DeleteDevice(ctx context.Context, id string) error {
	db := r.getDB(ctx)
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return db.DeleteDevice(ctx, uuidID)
}
