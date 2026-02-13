package model

import (
	"github.com/google/uuid"
)

type SystemAdmin struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	Department   string    `json:"department"`
	IsActive     bool      `json:"is_active"`
	LastLoginAt  *string   `json:"last_login_at"`
	BaseModel
}
