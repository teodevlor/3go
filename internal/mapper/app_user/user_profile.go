package appuser

import (
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
)

func ToUserProfile(row pgdb.UserProfile) *model.UserProfile {
	return &model.UserProfile{
		ID:        row.ID,
		AccountID: row.AccountID,
		FullName:  row.FullName,
		AvatarURL: row.AvatarUrl.String,
		IsActive:  row.IsActive.Bool,
		Metadata:  row.Metadata,
		BaseModel: model.BaseModel{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			DeletedAt: row.DeletedAt.Time,
		},
	}
}
