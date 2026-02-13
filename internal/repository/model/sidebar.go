package model

import (
	"time"

	"github.com/google/uuid"
)

type Sidebar struct {
	ID          uuid.UUID `json:"id"`
	Context     string    `json:"context"`
	Version     string    `json:"version"`
	GeneratedAt time.Time `json:"generated_at"`
	Items       []byte    `json:"items"` // JSON
	BaseModel
}
