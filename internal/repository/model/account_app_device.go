package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AccountAppDevice struct {
	ID         uuid.UUID       `json:"id"`
	AccountID  uuid.UUID       `json:"account_id"`
	DeviceID   uuid.UUID       `json:"device_id"`
	AppType    string          `json:"app_type"`
	FcmToken   string          `json:"fcm_token"`
	IsActive   bool            `json:"is_active"`
	LastUsedAt *time.Time      `json:"last_used_at"`
	Metadata   json.RawMessage `json:"metadata"`
	BaseModel
}
