package web_system

import "go-structure/internal/dto/common"

const (
	DepartmentEmployee = "employee"
	DepartmentAdmin    = "admin"
	DepartmentSeller   = "seller"
	DepartmentMarketer = "marketer"
)

var ValidDepartments = []string{DepartmentEmployee, DepartmentAdmin, DepartmentSeller, DepartmentMarketer}

type (
	CreateAdminRequestDto struct {
		Email      string   `json:"email" binding:"required,email,max=255"`
		Password   string   `json:"password" binding:"required,min=6"`
		FullName   string   `json:"full_name" binding:"max=255"`
		Department string   `json:"department" binding:"required,oneof=employee admin seller marketer"`
		IsActive   bool     `json:"is_active"`
		RoleIDs    []string `json:"role_ids"`
	}

	UpdateAdminRequestDto struct {
		Email      *string   `json:"email" binding:"required,email,max=255"`
		FullName   *string   `json:"full_name" binding:"required,max=255"`
		Department *string   `json:"department" binding:"required,oneof=employee admin seller marketer"`
		IsActive   *bool     `json:"is_active" binding:"required"`
		RoleIDs    *[]string `json:"role_ids" binding:"required"`
	}
)

type (
	AdminItemDto struct {
		ID          string        `json:"id"`
		Email       string        `json:"email"`
		FullName    string        `json:"full_name"`
		Department  string        `json:"department"`
		IsActive    bool          `json:"is_active"`
		LastLoginAt *string       `json:"last_login_at,omitempty"`
		Roles       []RoleItemDto `json:"roles"`
	}

	ListAdminsResponseDto struct {
		Items      []AdminItemDto        `json:"items"`
		Pagination common.PaginationMeta `json:"pagination"`
	}
)
