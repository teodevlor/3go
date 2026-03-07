package app_driver

import (
	dto_common "go-structure/internal/dto/common"
)

type (
	CreateDriverDocumentTypeRequestDto struct {
		Code              string  `json:"code" binding:"required,max=100"`
		Name              string  `json:"name" binding:"required,max=255"`
		Description       string  `json:"description"`
		IsRequired        bool    `json:"is_required"`
		RequireExpireDate bool    `json:"require_expire_date"`
		ServiceID         *string `json:"service_id"` // nil = áp dụng cho mọi service
		IsActive          bool    `json:"is_active"`
	}

	CreateDriverDocumentTypeResponseDto struct {
		ID                string  `json:"id"`
		Code              string  `json:"code"`
		Name              string  `json:"name"`
		Description       string  `json:"description"`
		IsRequired        bool    `json:"is_required"`
		RequireExpireDate bool    `json:"require_expire_date"`
		ServiceID         *string `json:"service_id"` // null là áp dụng document cho mọi dịch vụ
		IsActive          bool    `json:"is_active"`
	}

	UpdateDriverDocumentTypeRequestDto struct {
		Code              string  `json:"code" binding:"required,max=100"`
		Name              string  `json:"name" binding:"required,max=255"`
		Description       string  `json:"description"`
		IsRequired        bool    `json:"is_required"`
		RequireExpireDate bool    `json:"require_expire_date"`
		ServiceID         *string `json:"service_id"`
		IsActive          bool    `json:"is_active"`
	}

	UpdateDriverDocumentTypeResponseDto struct {
		ID                string  `json:"id"`
		Code              string  `json:"code"`
		Name              string  `json:"name"`
		Description       string  `json:"description"`
		IsRequired        bool    `json:"is_required"`
		RequireExpireDate bool    `json:"require_expire_date"`
		ServiceID         *string `json:"service_id"` // null là áp dụng document cho mọi dịch vụ
		IsActive          bool    `json:"is_active"`
	}

	DriverDocumentTypeItemDto struct {
		ID                string  `json:"id"`
		Code              string  `json:"code"`
		Name              string  `json:"name"`
		Description       string  `json:"description"`
		IsRequired        bool    `json:"is_required"`
		RequireExpireDate bool    `json:"require_expire_date"`
		ServiceID         *string `json:"service_id"` // null là áp dụng document cho mọi dịch vụ
		IsActive          bool    `json:"is_active"`
	}

	ListDriverDocumentTypesResponseDto struct {
		Items      []DriverDocumentTypeItemDto `json:"items"`
		Pagination dto_common.PaginationMeta   `json:"pagination"`
	}

	RequiredDriverDocumentTypesResponseDto struct {
		Items []DriverDocumentTypeItemDto `json:"items"`
	}
)
