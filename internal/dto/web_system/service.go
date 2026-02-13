package web_system

import (
	dto_common "go-structure/internal/dto/common"
)

type CreateServiceRequestDto struct {
	Code      string   `json:"code" binding:"required,max=100"`
	Name      string   `json:"name" binding:"required,max=255"`
	BasePrice float64  `json:"base_price" binding:"required,gte=0"`
	MinPrice  float64  `json:"min_price" binding:"required,gte=0"`
	IsActive  bool     `json:"is_active"`
	ZoneIDs   []string `json:"zone_ids"` // Danh sách ID zone mà service hoạt động
}

type CreateServiceResponseDto struct {
	ID        string   `json:"id"`
	Code      string   `json:"code"`
	Name      string   `json:"name"`
	BasePrice float64  `json:"base_price"`
	MinPrice  float64  `json:"min_price"`
	IsActive  bool     `json:"is_active"`
	ZoneIDs   []string `json:"zone_ids"`
}

type UpdateServiceRequestDto struct {
	Code      string   `json:"code" binding:"required,max=100"`
	Name      string   `json:"name" binding:"required,max=255"`
	BasePrice float64  `json:"base_price" binding:"required,gte=0"`
	MinPrice  float64  `json:"min_price" binding:"required,gte=0"`
	IsActive  bool     `json:"is_active"`
	ZoneIDs   []string `json:"zone_ids"` // Danh sách zone hoạt động (cập nhật thay thế)
}

type ServiceItemDto struct {
	ID        string   `json:"id"`
	Code      string   `json:"code"`
	Name      string   `json:"name"`
	BasePrice float64  `json:"base_price"`
	MinPrice  float64  `json:"min_price"`
	IsActive  bool     `json:"is_active"`
	ZoneIDs   []string `json:"zone_ids"`
}

type ListServicesResponseDto struct {
	Items      []ServiceItemDto           `json:"items"`
	Pagination dto_common.PaginationMeta `json:"pagination"`
}
