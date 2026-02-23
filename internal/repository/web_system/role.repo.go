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
	IRoleRepository interface {
		Create(ctx context.Context, arg pgdb.CreateRoleParams) (*websystem.Role, error)
		GetByID(ctx context.Context, id uuid.UUID) (*websystem.Role, error)
		GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*websystem.Role, error)
		GetByCode(ctx context.Context, code string) (*websystem.Role, error)
		List(ctx context.Context, search string, limit, offset int32) ([]*websystem.Role, error)
		Count(ctx context.Context, search string) (int64, error)
		Update(ctx context.Context, arg pgdb.UpdateRoleParams) (*websystem.Role, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	roleRepository struct {
		pool *pgxpool.Pool
	}
)

func NewRoleRepository(pool *pgxpool.Pool) IRoleRepository {
	return &roleRepository{pool: pool}
}

func (r *roleRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *roleRepository) Create(ctx context.Context, arg pgdb.CreateRoleParams) (*websystem.Role, error) {
	db := r.getDB(ctx)
	row, err := db.CreateRole(ctx, arg)
	if err != nil {
		return nil, err
	}
	return webmapper.ToRoleFromRow(&row), nil
}

func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*websystem.Role, error) {
	db := r.getDB(ctx)
	row, err := db.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return webmapper.ToRoleFromRow(&row), nil
}

func (r *roleRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*websystem.Role, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	db := r.getDB(ctx)
	rows, err := db.GetRolesByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	out := make([]*websystem.Role, 0, len(rows))
	for i := range rows {
		out = append(out, webmapper.ToRoleFromRow(&rows[i]))
	}
	return out, nil
}

func (r *roleRepository) GetByCode(ctx context.Context, code string) (*websystem.Role, error) {
	db := r.getDB(ctx)
	row, err := db.GetRoleByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return webmapper.ToRoleFromRow(&row), nil
}

func (r *roleRepository) List(ctx context.Context, search string, limit, offset int32) ([]*websystem.Role, error) {
	db := r.getDB(ctx)
	rows, err := db.ListRoles(ctx, pgdb.ListRolesParams{Column1: search, Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]*websystem.Role, 0, len(rows))
	for i := range rows {
		out = append(out, webmapper.ToRoleFromRow(&rows[i]))
	}
	return out, nil
}

func (r *roleRepository) Count(ctx context.Context, search string) (int64, error) {
	db := r.getDB(ctx)
	return db.CountRoles(ctx, search)
}

func (r *roleRepository) Update(ctx context.Context, arg pgdb.UpdateRoleParams) (*websystem.Role, error) {
	db := r.getDB(ctx)
	row, err := db.UpdateRole(ctx, arg)
	if err != nil {
		return nil, err
	}
	return webmapper.ToRoleFromRow(&row), nil
}

func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.DeleteRole(ctx, id)
}
