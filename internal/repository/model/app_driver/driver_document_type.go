package appdriver

import (
	"time"

	"github.com/google/uuid"
)

type DriverDocumentType struct {
	ID                uuid.UUID  `json:"id"`
	Code              string     `json:"code"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	IsRequired        bool       `json:"is_required"`
	RequireExpireDate bool       `json:"require_expire_date"`
	ServiceID         *uuid.UUID `json:"service_id"` // nil = áp dụng cho mọi service
	IsActive          bool       `json:"is_active"`
	CreatedAt         time.Time  `json:"created_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}
