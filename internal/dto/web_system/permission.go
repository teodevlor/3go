package web_system

import "go-structure/internal/dto/common"

type (
	CreatePermissionRequestDto struct {
		Resource    string `json:"resource" binding:"required,max=100"`
		Action      string `json:"action" binding:"required,max=100"`
		Name        string `json:"name" binding:"required,max=255"`
		Description string `json:"description"`
	}

	UpdatePermissionRequestDto struct {
		Resource    *string `json:"resource" binding:"required,max=100"`
		Action      *string `json:"action" binding:"required,max=100"`
		Name        *string `json:"name" binding:"required,max=255"`
		Description *string `json:"description" binding:"required"`
	}
)

type (
	PermissionItemDto struct {
		ID          string `json:"id"`
		Resource    string `json:"resource"`
		Action      string `json:"action"`
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	ListPermissionsResponseDto struct {
		Items      []PermissionItemDto   `json:"items"`
		Pagination common.PaginationMeta `json:"pagination"`
	}
)
