package websystem

import (
	"go-structure/internal/repository/model"

	"github.com/google/uuid"
)

type SurchargeRule struct {
	ID        uuid.UUID `json:"id"`
	ServiceID uuid.UUID `json:"service_id"`
	ZoneID    uuid.UUID `json:"zone_id"`
	Amount    float64   `json:"amount"`
	Unit      string    `json:"unit"`
	IsActive  bool      `json:"is_active"`
	Priority  int32     `json:"priority"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
	model.BaseModel
}
