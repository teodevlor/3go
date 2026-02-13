package mapper

import (
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
)

func ToSystemAdmin(row pgdb.SystemAdmin) *model.SystemAdmin {
	var lastLoginAt *string
	if row.LastLoginAt.Valid {
		lastLoginAtStr := row.LastLoginAt.Time.Format("2006-01-02T15:04:05Z07:00")
		lastLoginAt = &lastLoginAtStr
	}

	return &model.SystemAdmin{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		FullName:     row.FullName.String,
		Department:   row.Department,
		IsActive:     row.IsActive.Bool,
		LastLoginAt:  lastLoginAt,
		BaseModel: model.BaseModel{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		},
	}
}
