package repository

import (
	"context"

	"go-structure/internal/helper/database"
	pgdb "go-structure/internal/orm/db/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	ISystemAdminRoleRepository interface {
		SetAdminRoles(ctx context.Context, adminID uuid.UUID, roleIDs []uuid.UUID, assignedBy *uuid.UUID) error
		GetRoleIDsByAdminID(ctx context.Context, adminID uuid.UUID) ([]uuid.UUID, error)
	}

	systemAdminRoleRepository struct {
		pool *pgxpool.Pool
	}
)

func NewSystemAdminRoleRepository(pool *pgxpool.Pool) ISystemAdminRoleRepository {
	return &systemAdminRoleRepository{pool: pool}
}

func (r *systemAdminRoleRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *systemAdminRoleRepository) SetAdminRoles(ctx context.Context, adminID uuid.UUID, roleIDs []uuid.UUID, assignedBy *uuid.UUID) error {
	db := r.getDB(ctx)
	if err := db.DeleteAdminRolesByAdminID(ctx, adminID); err != nil {
		return err
	}
	var ab pgtype.UUID
	if assignedBy != nil {
		ab = pgtype.UUID{Bytes: *assignedBy, Valid: true}
	}
	for _, roleID := range roleIDs {
		if err := db.InsertAdminRole(ctx, pgdb.InsertAdminRoleParams{
			AdminID:    adminID,
			RoleID:     roleID,
			AssignedBy: ab,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *systemAdminRoleRepository) GetRoleIDsByAdminID(ctx context.Context, adminID uuid.UUID) ([]uuid.UUID, error) {
	db := r.getDB(ctx)
	return db.GetRoleIDsByAdminID(ctx, adminID)
}
