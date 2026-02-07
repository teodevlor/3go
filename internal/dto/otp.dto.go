package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateOTPRequestData struct {
	Target     string
	OtpCode    string
	Purpose    string
	MaxAttempt int
	ExpiresAt  time.Time
}

type ActiveOTPResponseData struct {
	ID           uuid.UUID
	OtpCode      string
	AttemptCount int32
	MaxAttempt   int32
}

type OTPResendRequestDto struct {
	Target  string `json:"target" binding:"required"`
	Purpose string `json:"purpose" binding:"required,oneof=register forgot-password"`
}

type OTPResendResponseDto struct {
	UserMessage string `json:"user_message" example:"Gửi lại mã OTP thành công"`
}
