package repository

import (
	"context"
	"errors"
	"time"

	"go-structure/internal/helper/database"
	"go-structure/internal/mapper"
	"go-structure/internal/repository/model"
	pgdb "go-structure/orm/db/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	ISystemAdminRepository interface {
		GetByEmail(ctx context.Context, email string) (*model.SystemAdmin, error)
		GetByID(ctx context.Context, id uuid.UUID) (*model.SystemAdmin, error)
		Create(ctx context.Context, arg pgdb.CreateSystemAdminParams) (*model.SystemAdmin, error)
		Update(ctx context.Context, arg pgdb.UpdateSystemAdminParams) (*model.SystemAdmin, error)
		List(ctx context.Context, search string, limit, offset int32) ([]*model.SystemAdmin, error)
		Count(ctx context.Context, search string) (int64, error)
		Delete(ctx context.Context, id uuid.UUID) error
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToSystemAdmin(row), nil
}

func (r *systemAdminRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.SystemAdmin, error) {
	db := r.getDB(ctx)
	row, err := db.GetSystemAdminByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ToSystemAdmin(row), nil
}

func (r *systemAdminRepository) Create(ctx context.Context, arg pgdb.CreateSystemAdminParams) (*model.SystemAdmin, error) {
	db := r.getDB(ctx)
	row, err := db.CreateSystemAdmin(ctx, arg)
	if err != nil {
		return nil, err
	}
	admin := mapper.ToSystemAdmin(row)
	return admin, nil
}

func (r *systemAdminRepository) Update(ctx context.Context, arg pgdb.UpdateSystemAdminParams) (*model.SystemAdmin, error) {
	db := r.getDB(ctx)
	row, err := db.UpdateSystemAdmin(ctx, arg)
	if err != nil {
		return nil, err
	}
	admin := mapper.ToSystemAdmin(row)
	return admin, nil
}

func (r *systemAdminRepository) List(ctx context.Context, search string, limit, offset int32) ([]*model.SystemAdmin, error) {
	db := r.getDB(ctx)
	rows, err := db.ListSystemAdmins(ctx, pgdb.ListSystemAdminsParams{Column1: search, Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]*model.SystemAdmin, 0, len(rows))
	for i := range rows {
		out = append(out, mapper.ToSystemAdmin(rows[i]))
	}
	return out, nil
}

func (r *systemAdminRepository) Count(ctx context.Context, search string) (int64, error) {
	db := r.getDB(ctx)
	return db.CountSystemAdmins(ctx, search)
}

func (r *systemAdminRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.DeleteSystemAdmin(ctx, id)
}

func (r *systemAdminRepository) UpdateLastLoginAt(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	now := time.Now()
	return db.UpdateSystemAdminLastLoginAt(ctx, pgdb.UpdateSystemAdminLastLoginAtParams{
		ID:          id,
		LastLoginAt: pgtype.Timestamptz{Time: now, Valid: true},
	})
}
