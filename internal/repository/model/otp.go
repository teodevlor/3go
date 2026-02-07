package model

import (
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	ID           uuid.UUID
	Target       string
	OtpCode      string
	Purpose      string
	AttemptCount int
	MaxAttempt   int
	ExpiresAt    time.Time
	UsedAt       time.Time
	Status       string
	Metadata     []byte
	BaseModel
}
