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
		ExpiresIn    int64                   `json:"expires_in" example:"900"`     // Số giây access token còn hiệu lực
		ExpiredAt    string                  `json:"expired_at" example:"2025-02-24T12:30:00Z"` // Thời điểm access token hết hạn (RFC3339)
		Admin        AdminProfileResponseDto `json:"admin"`
	}

	AdminRoleSimpleDto struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	// AdminPermissionSimpleDto permission trong response login (chỉ id, code, name).
	AdminPermissionSimpleDto struct {
		ID   string `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	}

	AdminProfileResponseDto struct {
		ID          uuid.UUID                   `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
		Email       string                      `json:"email" example:"admin@example.com"`
		FullName    string                      `json:"full_name" example:"Nguyễn Văn A"`
		Department  string                      `json:"department" example:"admin"`
		IsActive    bool                        `json:"is_active" example:"true"`
		Roles       []AdminRoleSimpleDto        `json:"roles"`
		Permissions []AdminPermissionSimpleDto  `json:"permissions"`
	}

	AdminRefreshTokenRequestDto struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	AdminRefreshTokenResponseDto struct {
		AccessToken  string `json:"access_token" example:"new_access_token"`
		RefreshToken string `json:"refresh_token" example:"new_refresh_token"`
	}
)
