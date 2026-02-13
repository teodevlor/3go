package model

import (
	"time"

	"github.com/google/uuid"
)

type SystemAdminRefreshToken struct {
	ID              uuid.UUID `json:"id"`
	AdminID         uuid.UUID `json:"admin_id"`
	RefreshTokenHash string    `json:"refresh_token_hash"`
	ExpiresAt       time.Time `json:"expires_at"`
	IsRevoked       bool      `json:"is_revoked"`
	RevokedAt       *time.Time `json:"revoked_at"`
	RevokedReason   string    `json:"revoked_reason"`
	LastActiveAt    time.Time `json:"last_active_at"`
	IpAddress       string    `json:"ip_address"`
	UserAgent       string    `json:"user_agent"`
	Metadata        []byte    `json:"metadata"`
	BaseModel
}
