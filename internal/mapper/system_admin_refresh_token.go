package mapper

import (
	"time"

	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
)

func ToSystemAdminRefreshToken(row pgdb.SystemAdminRefreshToken) *model.SystemAdminRefreshToken {
	var revokedAt *time.Time
	if row.RevokedAt.Valid {
		revokedAt = &row.RevokedAt.Time
	}

	return &model.SystemAdminRefreshToken{
		ID:              row.ID,
		AdminID:         row.AdminID,
		RefreshTokenHash: row.RefreshTokenHash,
		ExpiresAt:       row.ExpiresAt.Time,
		IsRevoked:       row.IsRevoked,
		RevokedAt:       revokedAt,
		RevokedReason:   row.RevokedReason.String,
		LastActiveAt:    row.LastActiveAt.Time,
		IpAddress:       row.IpAddress.String,
		UserAgent:       row.UserAgent.String,
		Metadata:        row.Metadata,
		BaseModel: model.BaseModel{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		},
	}
}
