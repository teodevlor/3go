package web_system

import "go-structure/internal/dto/common"

type (
	CreateRoleRequestDto struct {
		Code          string   `json:"code" binding:"required,max=64"`
		Name          string   `json:"name" binding:"required,max=255"`
		Description   string   `json:"description"`
		IsActive      bool     `json:"is_active"`
		PermissionIDs []string `json:"permission_ids"`
	}

	UpdateRoleRequestDto struct {
		Code          *string   `json:"code" binding:"required,max=64"`
		Name          *string   `json:"name" binding:"required,max=255"`
		Description   *string   `json:"description" binding:"required"`
		IsActive      *bool     `json:"is_active" binding:"required"`
		PermissionIDs *[]string `json:"permission_ids" binding:"required"`
	}

	RoleItemDto struct {
		ID            string   `json:"id"`
		Code          string   `json:"code"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		IsActive      bool     `json:"is_active"`
		PermissionIDs []string `json:"permission_ids"`
	}

	ListRolesResponseDto struct {
		Items      []RoleItemDto         `json:"items"`
		Pagination common.PaginationMeta `json:"pagination"`
	}
)
