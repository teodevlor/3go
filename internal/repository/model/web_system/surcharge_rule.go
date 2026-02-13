package websystem

import (
	"go-structure/internal/repository/model"

	"github.com/google/uuid"
)

type SurchargeRule struct {
	ID            uuid.UUID `json:"id"`
	ServiceID     uuid.UUID `json:"service_id"`
	ZoneID        uuid.UUID `json:"zone_id"`
	SurchargeType string    `json:"surcharge_type"`
	Amount        float64   `json:"amount"`
	Unit          string    `json:"unit"`
	Condition     []byte    `json:"condition"`
	IsActive      bool      `json:"is_active"`
	model.BaseModel
}
