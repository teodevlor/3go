package mapper

import (
	"go-structure/internal/common"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
	websystem "go-structure/internal/repository/model/web_system"
)

func ToSurchargeRuleFromRow(row *pgdb.SystemSurchargeRule) *websystem.SurchargeRule {
	if row == nil {
		return nil
	}
	return &websystem.SurchargeRule{
		ID:            row.ID,
		ServiceID:     row.ServiceID,
		ZoneID:        row.ZoneID,
		SurchargeType: row.SurchargeType,
		Amount:        common.NumericToFloat64(row.Amount),
		Unit:          row.Unit,
		Condition:     row.Condition,
		IsActive:      row.IsActive,
		BaseModel: model.BaseModel{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			DeletedAt: row.DeletedAt.Time,
		},
	}
}
