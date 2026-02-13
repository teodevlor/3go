package websystem

import (
	"go-structure/internal/repository/model"

	"github.com/google/uuid"
)

type PackageSizePricing struct {
	ID          uuid.UUID `json:"id"`
	ServiceID   uuid.UUID `json:"service_id"`
	PackageSize string    `json:"package_size"`
	ExtraPrice  float64   `json:"extra_price"`
	IsActive    bool      `json:"is_active"`
	model.BaseModel
}
