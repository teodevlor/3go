package mapper

import (
	"go-structure/internal/common"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
	websystem "go-structure/internal/repository/model/web_system"
)

func ToDistancePricingRuleFromRow(row *pgdb.SystemDistancePricingRule) *websystem.DistancePricingRule {
	if row == nil {
		return nil
	}
	return &websystem.DistancePricingRule{
		ID:         row.ID,
		ServiceID:  row.ServiceID,
		FromKm:     common.NumericToFloat64(row.FromKm),
		ToKm:       common.NumericToFloat64(row.ToKm),
		PricePerKm: common.NumericToFloat64(row.PricePerKm),
		IsActive:   row.IsActive,
		BaseModel: model.BaseModel{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			DeletedAt: row.DeletedAt.Time,
		},
	}
}
