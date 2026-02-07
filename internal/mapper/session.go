package mapper

import (
	"time"

	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
)

func ToSession(row pgdb.Session) *model.Session {
	var revokedAt *time.Time
	if row.RevokedAt.Valid {
		t := row.RevokedAt.Time
		revokedAt = &t
	}
	return &model.Session{
		ID:                 row.ID,
		AccountAppDeviceID: row.AccountAppDeviceID,
		RefreshTokenHash:   row.RefreshTokenHash,
		ExpiresAt:          row.ExpiresAt.Time,
		IsRevoked:          row.IsRevoked,
		RevokedAt:          revokedAt,
		RevokedReason:      row.RevokedReason.String,
		LastActiveAt:       row.LastActiveAt.Time,
		IpAddress:          row.IpAddress.String,
		UserAgent:          row.UserAgent.String,
		Metadata:           row.Metadata,
	}
}
