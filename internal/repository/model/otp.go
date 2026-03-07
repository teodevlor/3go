package model

import (
	"time"

	"go-structure/internal/common"

	"github.com/google/uuid"
)

const (
	OTPStatusActive  = common.OTPStatusActive
	OTPStatusUsed    = common.OTPStatusUsed
	OTPStatusExpired = common.OTPStatusExpired
	OTPStatusLocked  = common.OTPStatusLocked
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

func (otp *OTP) IsExpired() bool {
	return time.Now().After(otp.ExpiresAt)
}

func (otp *OTP) IsLocked() bool {
	return otp.Status == OTPStatusLocked
}

func (otp *OTP) IsUsed() bool {
	return otp.Status == OTPStatusUsed
}

func (otp *OTP) IsActive() bool {
	return otp.Status == OTPStatusActive && !otp.IsExpired()
}

func (otp *OTP) IsCodeMatched(code string) bool {
	return otp.OtpCode == code
}

// WillExceedMaxAttemptAfterThis check if the OTP will exceed the max attempt after this attempt
func (otp *OTP) WillExceedMaxAttemptAfterThis() bool {
	return otp.AttemptCount+1 >= otp.MaxAttempt
}
