package model

import (
	"github.com/google/uuid"
)

type DistancePricingRules struct {
	ID         uuid.UUID `json:"id"`
	ServiceID  uuid.UUID `json:"service_id"`
	FromKm     float64   `json:"from_km"`
	ToKm       float64   `json:"to_km"`
	PricePerKm float64   `json:"price_per_km"`
	IsActive   bool      `json:"is_active"`
	BaseModel
}
