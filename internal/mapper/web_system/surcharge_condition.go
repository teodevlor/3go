package mapper

import (
	pgdb "go-structure/internal/orm/db/postgres"
	websystem "go-structure/internal/repository/model/web_system"
)

func ToSurchargeConditionFromRow(row *pgdb.SystemSurchargeCondition) *websystem.SurchargeCondition {
	if row == nil {
		return nil
	}
	return &websystem.SurchargeCondition{
		ID:            row.ID,
		Code:          row.Code,
		ConditionType: row.ConditionType,
		Config:        row.Config,
		IsActive:      row.IsActive.Bool,
	}
}

