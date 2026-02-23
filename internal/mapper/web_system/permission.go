package mapper

import (
	pgdb "go-structure/internal/orm/db/postgres"
	websystem "go-structure/internal/repository/model/web_system"
)

func ToPermissionFromRow(row *pgdb.SystemPermission) *websystem.Permission {
	if row == nil {
		return nil
	}
	desc := ""
	if row.Description.Valid {
		desc = row.Description.String
	}
	code := ""
	if row.Code.Valid {
		code = row.Code.String
	}
	return &websystem.Permission{
		ID:          row.ID,
		Resource:    row.Resource,
		Action:      row.Action,
		Code:        code,
		Name:        row.Name,
		Description: desc,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}
