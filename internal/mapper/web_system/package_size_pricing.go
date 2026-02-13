package mapper

import (
	"go-structure/internal/common"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
	websystem "go-structure/internal/repository/model/web_system"
)

func ToPackageSizePricingFromRow(row *pgdb.SystemPackageSizePricing) *websystem.PackageSizePricing {
	if row == nil {
		return nil
	}
	return &websystem.PackageSizePricing{
		ID:          row.ID,
		ServiceID:   row.ServiceID,
		PackageSize: row.PackageSize,
		ExtraPrice:  common.NumericToFloat64(row.ExtraPrice),
		IsActive:    row.IsActive,
		BaseModel: model.BaseModel{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			DeletedAt: row.DeletedAt.Time,
		},
	}
}
