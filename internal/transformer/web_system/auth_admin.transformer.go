package web_system

import (
	websystemdto "go-structure/internal/dto/web_system"
	"go-structure/internal/repository/model"
)

func ToAdminLoginResponseDto(accessToken string, refreshToken string, admin *model.SystemAdmin) websystemdto.AdminLoginResponseDto {
	return websystemdto.AdminLoginResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Admin: websystemdto.AdminProfileResponseDto{
			ID:         admin.ID,
			Email:      admin.Email,
			FullName:   admin.FullName,
			Department: admin.Department,
			IsActive:   admin.IsActive,
		},
	}
}

func ToAdminRefreshTokenResponseDto(accessToken string, refreshToken string) websystemdto.AdminRefreshTokenResponseDto {
	return websystemdto.AdminRefreshTokenResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
