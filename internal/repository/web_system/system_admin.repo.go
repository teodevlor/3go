package repository

import (
	"context"
	"time"

	"go-structure/internal/helper/database"
	"go-structure/internal/mapper"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	ISystemAdminRepository interface {
		GetByEmail(ctx context.Context, email string) (*model.SystemAdmin, error)
		GetByID(ctx context.Context, id uuid.UUID) (*model.SystemAdmin, error)
		UpdateLastLoginAt(ctx context.Context, id uuid.UUID) error
	}

	systemAdminRepository struct {
		pool *pgxpool.Pool
	}
)

func NewSystemAdminRepository(pool *pgxpool.Pool) ISystemAdminRepository {
	return &systemAdminRepository{pool: pool}
}

func (r *systemAdminRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *systemAdminRepository) GetByEmail(ctx context.Context, email string) (*model.SystemAdmin, error) {
	db := r.getDB(ctx)
	row, err := db.GetSystemAdminByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return mapper.ToSystemAdmin(row), nil
}

func (r *systemAdminRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.SystemAdmin, error) {
	db := r.getDB(ctx)
	row, err := db.GetSystemAdminByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.ToSystemAdmin(row), nil
}

func (r *systemAdminRepository) UpdateLastLoginAt(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	now := time.Now()
	return db.UpdateSystemAdminLastLoginAt(ctx, pgdb.UpdateSystemAdminLastLoginAtParams{
		ID:          id,
		LastLoginAt: pgtype.Timestamptz{Time: now, Valid: true},
	})
}
