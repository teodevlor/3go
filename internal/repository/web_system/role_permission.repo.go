package repository

import (
	"context"

	"go-structure/internal/helper/database"
	pgdb "go-structure/internal/orm/db/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	IRolePermissionRepository interface {
		CreateRolePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
		DeleteByRoleID(ctx context.Context, roleID uuid.UUID) error
		GetPermissionIDsByRoleID(ctx context.Context, roleID uuid.UUID) ([]uuid.UUID, error)
		GetPermissionIDsByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) (map[uuid.UUID][]uuid.UUID, error)
	}

	rolePermissionRepository struct {
		pool *pgxpool.Pool
	}
)

func NewRolePermissionRepository(pool *pgxpool.Pool) IRolePermissionRepository {
	return &rolePermissionRepository{pool: pool}
}

func (r *rolePermissionRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *rolePermissionRepository) CreateRolePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	if len(permissionIDs) == 0 {
		return nil
	}
	db := r.getDB(ctx)
	return db.CreateRolePermissionsBatch(ctx, pgdb.CreateRolePermissionsBatchParams{
		RoleID:  roleID,
		Column2: permissionIDs,
	})
}

func (r *rolePermissionRepository) DeleteByRoleID(ctx context.Context, roleID uuid.UUID) error {
	db := r.getDB(ctx)
	return db.DeleteRolePermissionsByRoleID(ctx, roleID)
}

func (r *rolePermissionRepository) GetPermissionIDsByRoleID(ctx context.Context, roleID uuid.UUID) ([]uuid.UUID, error) {
	db := r.getDB(ctx)
	return db.GetPermissionIDsByRoleID(ctx, roleID)
}

func (r *rolePermissionRepository) GetPermissionIDsByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) (map[uuid.UUID][]uuid.UUID, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}
	db := r.getDB(ctx)
	pairs, err := db.GetRolePermissionPairsByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, err
	}
	out := make(map[uuid.UUID][]uuid.UUID)
	for _, p := range pairs {
		out[p.RoleID] = append(out[p.RoleID], p.PermissionID)
	}
	return out, nil
}
