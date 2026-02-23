package web_system

import (
	dto "go-structure/internal/dto/web_system"
	"go-structure/internal/repository/model"
)

func ToAdminItemDto(admin *model.SystemAdmin, roles []dto.RoleItemDto) dto.AdminItemDto {
	if admin == nil {
		return dto.AdminItemDto{}
	}
	if roles == nil {
		roles = []dto.RoleItemDto{}
	}
	return dto.AdminItemDto{
		ID:          admin.ID.String(),
		Email:       admin.Email,
		FullName:    admin.FullName,
		Department:  admin.Department,
		IsActive:    admin.IsActive,
		LastLoginAt: admin.LastLoginAt,
		Roles:       roles,
	}
}
