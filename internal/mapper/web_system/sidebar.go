package mapper

import (
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
)

func ToSidebar(row pgdb.SystemSidebar) *model.Sidebar {
	return &model.Sidebar{
		ID:          row.ID,
		Context:     row.Context,
		Version:     row.Version,
		GeneratedAt: row.GeneratedAt.Time,
		Items:       row.Items,
		BaseModel: model.BaseModel{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			DeletedAt: row.DeletedAt.Time,
		},
	}
}
