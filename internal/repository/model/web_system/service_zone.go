package websystem

import (
	"go-structure/internal/repository/model"

	"github.com/google/uuid"
)

type ServiceZone struct {
	ID        uuid.UUID `json:"id"`
	ZoneID    uuid.UUID `json:"zone_id"`
	ServiceID uuid.UUID `json:"service_id"`
	model.BaseModel
}
