package dto

import "github.com/google/uuid"

type CreateOTPAuditRequestData struct {
	OTPId         uuid.UUID
	Target        string
	Purpose       string
	AttemptNumber int
	Result        string
	FailureReason string
	IPAddress     string
	UserAgent     string
	Metadata      []byte
}
