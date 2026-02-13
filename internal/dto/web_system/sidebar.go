package web_system

import (
	dto_common "go-structure/internal/dto/common"
	"time"
)

type SidebarBadgeDto struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type SidebarItemDto struct {
	ID                 string           `json:"id"`
	Label              string           `json:"label"`
	Icon               string           `json:"icon,omitempty"`
	Path               string           `json:"path,omitempty"`
	Order              int              `json:"order,omitempty"`
	Visible            bool             `json:"visible,omitempty"`
	Children           []SidebarItemDto `json:"children,omitempty"`
	PermissionRequired []string         `json:"permission_required,omitempty"`
	Badge              *SidebarBadgeDto `json:"badge,omitempty"`
	FeatureFlag        string           `json:"feature_flag,omitempty"`
}

type CreateSidebarRequestDto struct {
	Context     string           `json:"context" binding:"required,max=100"`
	Version     string           `json:"version" binding:"max=50"`
	GeneratedAt *time.Time       `json:"generated_at,omitempty"`
	Items       []SidebarItemDto `json:"items" binding:"required"`
}

type UpdateSidebarRequestDto struct {
	Context     string           `json:"context" binding:"required,max=100"`
	Version     string           `json:"version" binding:"max=50"`
	GeneratedAt *time.Time       `json:"generated_at,omitempty"`
	Items       []SidebarItemDto `json:"items" binding:"required"`
}

type SidebarResponseDto struct {
	ID          string           `json:"id"`
	Context     string           `json:"context"` // system | app_user...
	Version     string           `json:"version"`
	GeneratedAt *time.Time       `json:"generated_at,omitempty"`
	Items       []SidebarItemDto `json:"items"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type ListSidebarsResponseDto struct {
	Items      []SidebarResponseDto `json:"items"`
	Pagination dto_common.PaginationMeta `json:"pagination"`
}
