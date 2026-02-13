package websystem

import (
	"go-structure/internal/repository/model"

	"github.com/google/uuid"
)

type Service struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	BasePrice float64   `json:"base_price"`
	MinPrice  float64   `json:"min_price"`
	IsActive  bool      `json:"is_active"`
	model.BaseModel
}
