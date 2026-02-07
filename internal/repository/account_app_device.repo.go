package repository

import (
	"context"
	"time"

	"go-structure/internal/helper/database"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IAccountAppDeviceRepository interface {
	CreateAccountAppDevice(ctx context.Context, a *model.AccountAppDevice) (*model.AccountAppDevice, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.AccountAppDevice, error)
	GetByAccountDeviceAndAppType(ctx context.Context, accountID, deviceID uuid.UUID, appType string) (*model.AccountAppDevice, error)
	UpdateAccountAppDevice(ctx context.Context, a *model.AccountAppDevice) (*model.AccountAppDevice, error)
}

type accountAppDeviceRepository struct {
	pool *pgxpool.Pool
}

func NewAccountAppDeviceRepository(pool *pgxpool.Pool) IAccountAppDeviceRepository {
	return &accountAppDeviceRepository{pool: pool}
}

func (r *accountAppDeviceRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *accountAppDeviceRepository) CreateAccountAppDevice(ctx context.Context, a *model.AccountAppDevice) (*model.AccountAppDevice, error) {
	db := r.getDB(ctx)

	row, err := db.CreateAccountAppDevice(ctx, pgdb.CreateAccountAppDeviceParams{
		AccountID: a.AccountID,
		DeviceID:  a.DeviceID,
		AppType:   a.AppType,
		FcmToken: pgtype.Text{
			String: a.FcmToken,
			Valid:  a.FcmToken != "",
		},
		IsActive: a.IsActive,
		Metadata: a.Metadata,
	})
	if err != nil {
		return nil, err
	}

	return mapAccountAppDevice(row), nil
}

func (r *accountAppDeviceRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.AccountAppDevice, error) {
	db := r.getDB(ctx)
	row, err := db.GetAccountAppDeviceByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return mapAccountAppDevice(row), nil
}

func (r *accountAppDeviceRepository) GetByAccountDeviceAndAppType(ctx context.Context, accountID, deviceID uuid.UUID, appType string) (*model.AccountAppDevice, error) {
	db := r.getDB(ctx)
	row, err := db.GetAccountAppDevice(ctx, pgdb.GetAccountAppDeviceParams{
		AccountID: accountID,
		DeviceID:  deviceID,
		AppType:   appType,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return mapAccountAppDevice(row), nil
}

func (r *accountAppDeviceRepository) UpdateAccountAppDevice(ctx context.Context, a *model.AccountAppDevice) (*model.AccountAppDevice, error) {
	db := r.getDB(ctx)

	var lastUsedAt pgtype.Timestamptz
	if a.LastUsedAt != nil {
		lastUsedAt = pgtype.Timestamptz{
			Time:  *a.LastUsedAt,
			Valid: true,
		}
	}

	row, err := db.UpdateAccountAppDevice(ctx, pgdb.UpdateAccountAppDeviceParams{
		ID: a.ID,
		FcmToken: pgtype.Text{
			String: a.FcmToken,
			Valid:  a.FcmToken != "",
		},
		IsActive: pgtype.Bool{
			Bool:  a.IsActive,
			Valid: true,
		},
		LastUsedAt: lastUsedAt,
		Metadata:   a.Metadata,
	})
	if err != nil {
		return nil, err
	}

	return mapAccountAppDevice(row), nil
}

func mapAccountAppDevice(row pgdb.AccountAppDevice) *model.AccountAppDevice {
	var lastUsed *time.Time
	if row.LastUsedAt.Valid {
		t := row.LastUsedAt.Time
		lastUsed = &t
	}

	return &model.AccountAppDevice{
		ID:        row.ID,
		AccountID: row.AccountID,
		DeviceID:  row.DeviceID,
		AppType:   row.AppType,
		FcmToken:  row.FcmToken.String,
		IsActive:  row.IsActive,
		LastUsedAt: lastUsed,
		Metadata:  row.Metadata,
		BaseModel: model.BaseModel{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			DeletedAt: row.DeletedAt.Time,
		},
	}
}

