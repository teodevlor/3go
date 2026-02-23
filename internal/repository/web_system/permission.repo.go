package repository

import (
	"context"

	"go-structure/internal/helper/database"
	webmapper "go-structure/internal/mapper/web_system"
	pgdb "go-structure/internal/orm/db/postgres"
	websystem "go-structure/internal/repository/model/web_system"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	IPermissionRepository interface {
		Create(ctx context.Context, arg pgdb.CreatePermissionParams) (*websystem.Permission, error)
		GetByID(ctx context.Context, id uuid.UUID) (*websystem.Permission, error)
		GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*websystem.Permission, error)
		GetByCode(ctx context.Context, code string) (*websystem.Permission, error)
		List(ctx context.Context, search string, limit, offset int32) ([]*websystem.Permission, error)
		Count(ctx context.Context, search string) (int64, error)
		Update(ctx context.Context, arg pgdb.UpdatePermissionParams) (*websystem.Permission, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	permissionRepository struct {
		pool *pgxpool.Pool
	}
)

func NewPermissionRepository(pool *pgxpool.Pool) IPermissionRepository {
	return &permissionRepository{pool: pool}
}

func (r *permissionRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *permissionRepository) Create(ctx context.Context, arg pgdb.CreatePermissionParams) (*websystem.Permission, error) {
	db := r.getDB(ctx)
	row, err := db.CreatePermission(ctx, arg)
	if err != nil {
		return nil, err
	}
	return webmapper.ToPermissionFromRow(&row), nil
}

func (r *permissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*websystem.Permission, error) {
	db := r.getDB(ctx)
	row, err := db.GetPermissionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return webmapper.ToPermissionFromRow(&row), nil
}

func (r *permissionRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*websystem.Permission, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	db := r.getDB(ctx)
	rows, err := db.GetPermissionsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	out := make([]*websystem.Permission, 0, len(rows))
	for i := range rows {
		out = append(out, webmapper.ToPermissionFromRow(&rows[i]))
	}
	return out, nil
}

func (r *permissionRepository) GetByCode(ctx context.Context, code string) (*websystem.Permission, error) {
	db := r.getDB(ctx)
	codeParam := pgtype.Text{String: code, Valid: true}
	row, err := db.GetPermissionByCode(ctx, codeParam)
	if err != nil {
		return nil, err
	}
	return webmapper.ToPermissionFromRow(&row), nil
}

func (r *permissionRepository) List(ctx context.Context, search string, limit, offset int32) ([]*websystem.Permission, error) {
	db := r.getDB(ctx)
	rows, err := db.ListPermissions(ctx, pgdb.ListPermissionsParams{Column1: search, Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]*websystem.Permission, 0, len(rows))
	for i := range rows {
		out = append(out, webmapper.ToPermissionFromRow(&rows[i]))
	}
	return out, nil
}

func (r *permissionRepository) Count(ctx context.Context, search string) (int64, error) {
	db := r.getDB(ctx)
	return db.CountPermissions(ctx, search)
}

func (r *permissionRepository) Update(ctx context.Context, arg pgdb.UpdatePermissionParams) (*websystem.Permission, error) {
	db := r.getDB(ctx)
	row, err := db.UpdatePermission(ctx, arg)
	if err != nil {
		return nil, err
	}
	return webmapper.ToPermissionFromRow(&row), nil
}

func (r *permissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.DeletePermission(ctx, id)
}
