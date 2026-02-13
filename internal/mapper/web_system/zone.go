package mapper

import (
	"go-structure/internal/common"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
)

func ToZone(row *pgdb.SystemZone) *model.Zone {
	if row == nil {
		return nil
	}
	return &model.Zone{
		ID:              row.ID,
		Code:            row.Code,
		Name:            row.Name,
		Polygon:         "",
		PriceMultiplier: common.NumericToFloat64(row.PriceMultiplier),
		IsActive:        row.IsActive,
		BaseModel: model.BaseModel{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			DeletedAt: row.DeletedAt.Time,
		},
	}
}
