package appuser

import "github.com/google/uuid"

type (
	UserRegisterRequestDto struct {
		Phone    string `json:"phone" binding:"required"`
		FullName string `json:"full_name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	UserRegisterResponseDto struct {
		ID        uuid.UUID `json:"id"`
		Phone     string    `json:"phone"`
		FullName  string    `json:"full_name"`
		CreatedAt string    `json:"created_at"`
		UpdatedAt string    `json:"updated_at"`
	}
)
