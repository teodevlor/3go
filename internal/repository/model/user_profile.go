package model

import (
	"github.com/google/uuid"
)

type UserProfile struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	FullName  string
	AvatarURL string
	IsActive  bool
	Metadata  []byte
	CreatedAt string
	UpdatedAt string
}
