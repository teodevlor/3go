package repository

import (
	"context"
	"strconv"

	"go-structure/internal/common"
	"go-structure/internal/helper/database"
	webmapper "go-structure/internal/mapper/web_system"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	IZoneRepository interface {
		CreateZone(ctx context.Context, zone *model.Zone) (*model.Zone, error)
		GetZoneByID(ctx context.Context, id uuid.UUID) (*model.Zone, error)
		GetZoneByCode(ctx context.Context, code string) (*model.Zone, error)
		ListZones(ctx context.Context, search string, limit, offset int32) ([]*model.Zone, error)
		CountZones(ctx context.Context, search string) (int64, error)
		UpdateZone(ctx context.Context, zone *model.Zone) (*model.Zone, error)
		DeleteZone(ctx context.Context, id uuid.UUID) error
	}

	zoneRepository struct {
		pool *pgxpool.Pool
	}
)

func NewZoneRepository(pool *pgxpool.Pool) IZoneRepository {
	return &zoneRepository{pool: pool}
}

func (r *zoneRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *zoneRepository) CreateZone(ctx context.Context, zone *model.Zone) (*model.Zone, error) {
	db := r.getDB(ctx)

	var priceMultiplier pgtype.Numeric
	if err := priceMultiplier.Scan(strconv.FormatFloat(zone.PriceMultiplier, 'f', -1, 64)); err != nil {
		return nil, err
	}

	row, err := db.CreateZone(ctx, pgdb.CreateZoneParams{
		Code:              zone.Code,
		Name:              zone.Name,
		StGeomfromgeojson: zone.Polygon,
		PriceMultiplier:   priceMultiplier,
		IsActive:          zone.IsActive,
	})
	if err != nil {
		return nil, err
	}
	return webmapper.ToZone(&row), nil
}

func (r *zoneRepository) GetZoneByID(ctx context.Context, id uuid.UUID) (*model.Zone, error) {
	db := r.getDB(ctx)
	row, err := db.GetZoneByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// GetZoneByIDRow có đầy đủ polygon GeoJSON, nhưng mapper zone hiện
	// đang chuẩn hoá theo style ToDistancePricingRuleFromRow (SystemZone).
	// Ở đây vẫn có thể map thủ công hoặc sau này tách mapper riêng nếu cần polygon.
	return &model.Zone{
		ID:              row.ID,
		Code:            row.Code,
		Name:            row.Name,
		Polygon:         row.PolygonGeojson,
		PriceMultiplier: common.NumericToFloat64(row.PriceMultiplier),
		IsActive:        row.IsActive,
		BaseModel: model.BaseModel{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			DeletedAt: row.DeletedAt.Time,
		},
	}, nil
}

func (r *zoneRepository) GetZoneByCode(ctx context.Context, code string) (*model.Zone, error) {
	db := r.getDB(ctx)
	row, err := db.GetZoneByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return &model.Zone{
		ID:              row.ID,
		Code:            row.Code,
		Name:            row.Name,
		Polygon:         row.PolygonGeojson,
		PriceMultiplier: common.NumericToFloat64(row.PriceMultiplier),
		IsActive:        row.IsActive,
		BaseModel: model.BaseModel{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			DeletedAt: row.DeletedAt.Time,
		},
	}, nil
}

func (r *zoneRepository) ListZones(ctx context.Context, search string, limit, offset int32) ([]*model.Zone, error) {
	db := r.getDB(ctx)
	rows, err := db.ListZones(ctx, pgdb.ListZonesParams{Column1: search, Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]*model.Zone, 0, len(rows))
	for _, row := range rows {
		out = append(out, &model.Zone{
			ID:              row.ID,
			Code:            row.Code,
			Name:            row.Name,
			Polygon:         row.PolygonGeojson,
			PriceMultiplier: common.NumericToFloat64(row.PriceMultiplier),
			IsActive:        row.IsActive,
			BaseModel: model.BaseModel{
				CreatedAt: row.CreatedAt.Time,
				UpdatedAt: row.UpdatedAt.Time,
				DeletedAt: row.DeletedAt.Time,
			},
		})
	}
	return out, nil
}

func (r *zoneRepository) CountZones(ctx context.Context, search string) (int64, error) {
	db := r.getDB(ctx)
	return db.CountZones(ctx, search)
}

func (r *zoneRepository) UpdateZone(ctx context.Context, zone *model.Zone) (*model.Zone, error) {
	db := r.getDB(ctx)
	var priceMultiplier pgtype.Numeric
	if err := priceMultiplier.Scan(strconv.FormatFloat(zone.PriceMultiplier, 'f', -1, 64)); err != nil {
		return nil, err
	}
	row, err := db.UpdateZone(ctx, pgdb.UpdateZoneParams{
		ID:                zone.ID,
		Code:              zone.Code,
		Name:              zone.Name,
		StGeomfromgeojson: zone.Polygon,
		PriceMultiplier:   priceMultiplier,
		IsActive:          zone.IsActive,
	})
	if err != nil {
		return nil, err
	}
	return webmapper.ToZone(&row), nil
}

func (r *zoneRepository) DeleteZone(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.DeleteZone(ctx, id)
}
