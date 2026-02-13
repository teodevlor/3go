package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SystemLoginHistory struct {
	ID            uuid.UUID       `json:"id"`
	AdminID       uuid.UUID       `json:"admin_id"`
	LoginAt       time.Time       `json:"login_at"`
	Result        string          `json:"result"`
	FailureReason string          `json:"failure_reason"`
	IpAddress     string          `json:"ip_address"`
	UserAgent     string          `json:"user_agent"`
	Location      string          `json:"location"`
	Metadata      json.RawMessage `json:"metadata"`
	CreatedAt     time.Time       `json:"created_at"`
}
