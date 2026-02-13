package web_system

import "github.com/google/uuid"

type (
	AdminLoginRequestDto struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	AdminLoginResponseDto struct {
		AccessToken  string                  `json:"access_token" example:"access_token"`
		RefreshToken string                  `json:"refresh_token" example:"refresh_token"`
		Admin        AdminProfileResponseDto `json:"admin"`
	}

	AdminProfileResponseDto struct {
		ID         uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
		Email      string    `json:"email" example:"admin@example.com"`
		FullName   string    `json:"full_name" example:"Nguyễn Văn A"`
		Department string    `json:"department" example:"admin"`
		IsActive   bool      `json:"is_active" example:"true"`
	}

	AdminRefreshTokenRequestDto struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	AdminRefreshTokenResponseDto struct {
		AccessToken  string `json:"access_token" example:"new_access_token"`
		RefreshToken string `json:"refresh_token" example:"new_refresh_token"`
	}
)
