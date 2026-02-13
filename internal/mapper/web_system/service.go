package mapper

import (
	"go-structure/internal/common"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
	websystem "go-structure/internal/repository/model/web_system"
)

func ToServiceFromRow(row *pgdb.SystemService) *websystem.Service {
	if row == nil {
		return nil
	}
	return &websystem.Service{
		ID:        row.ID,
		Code:      row.Code,
		Name:      row.Name,
		BasePrice: common.NumericToFloat64(row.BasePrice),
		MinPrice:  common.NumericToFloat64(row.MinPrice),
		IsActive:  row.IsActive,
		BaseModel: model.BaseModel{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			DeletedAt: row.DeletedAt.Time,
		},
	}
}

func ToServiceZoneFromRow(row *pgdb.SystemServiceZone) *websystem.ServiceZone {
	if row == nil {
		return nil
	}
	return &websystem.ServiceZone{
		ID:        row.ID,
		ZoneID:    row.ZoneID,
		ServiceID: row.ServiceID,
		BaseModel: model.BaseModel{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			DeletedAt: row.DeletedAt.Time,
		},
	}
}
