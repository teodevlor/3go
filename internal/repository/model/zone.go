package model

import (
	"github.com/google/uuid"
)

type Zone struct {
	ID              uuid.UUID `json:"id"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	Polygon         string    `json:"polygon"`
	PriceMultiplier float64   `json:"price_multiplier"`
	IsActive        bool      `json:"is_active"`
	BaseModel
}
