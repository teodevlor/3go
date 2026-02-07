package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Setting struct {
	ID          uuid.UUID       `json:"id"`
	AccountID   uuid.UUID       `json:"account_id"`
	Key         string          `json:"key"`
	Value       json.RawMessage `json:"value"`
	Type        string          `json:"type"`
	Description string          `json:"description"`
	IsActive    bool            `json:"is_active"`
	Metadata    json.RawMessage `json:"metadata"`
	BaseModel
}
