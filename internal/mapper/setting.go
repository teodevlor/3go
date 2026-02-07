package mapper

import (
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
)

func ToSetting(row pgdb.Setting) *model.Setting {
	settingType := ""
	if t, ok := row.Type.(string); ok {
		settingType = t
	}

	return &model.Setting{
		ID:          row.ID,
		AccountID:   row.AccountID,
		Key:         row.Key,
		Value:       row.Value,
		Type:        settingType,
		Description: row.Description.String,
		IsActive:    row.IsActive.Bool,
		Metadata:    row.Metadata,
		BaseModel: model.BaseModel{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			DeletedAt: row.DeletedAt.Time,
		},
	}
}
