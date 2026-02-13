package repository

import (
	"context"

	"go-structure/internal/helper/database"
	pgdb "go-structure/internal/orm/db/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	IServiceZoneRepository interface {
		CreateServiceZones(ctx context.Context, serviceID uuid.UUID, zoneIDs []uuid.UUID) error
		ListServiceZoneIDsByServiceID(ctx context.Context, serviceID uuid.UUID) ([]uuid.UUID, error)
		DeleteServiceZonesByServiceID(ctx context.Context, serviceID uuid.UUID) error
	}

	serviceZoneRepository struct {
		pool *pgxpool.Pool
	}
)

func NewServiceZoneRepository(pool *pgxpool.Pool) IServiceZoneRepository {
	return &serviceZoneRepository{pool: pool}
}

func (r *serviceZoneRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *serviceZoneRepository) CreateServiceZones(ctx context.Context, serviceID uuid.UUID, zoneIDs []uuid.UUID) error {
	db := r.getDB(ctx)
	for _, zid := range zoneIDs {
		_, err := db.CreateServiceZone(ctx, pgdb.CreateServiceZoneParams{ZoneID: zid, ServiceID: serviceID})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *serviceZoneRepository) ListServiceZoneIDsByServiceID(ctx context.Context, serviceID uuid.UUID) ([]uuid.UUID, error) {
	db := r.getDB(ctx)
	rows, err := db.ListServiceZonesByServiceID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	ids := make([]uuid.UUID, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ZoneID)
	}
	return ids, nil
}

func (r *serviceZoneRepository) DeleteServiceZonesByServiceID(ctx context.Context, serviceID uuid.UUID) error {
	db := r.getDB(ctx)
	return db.DeleteServiceZonesByServiceID(ctx, serviceID)
}

