package model

import (
	"github.com/google/uuid"
)

type OTPAudit struct {
	ID            uuid.UUID
	OTPId         uuid.UUID
	Target        string
	Purpose       string
	AttemptNumber int
	Result        string
	FailureReason string
	IPAddress     string
	UserAgent     string
	Metadata      []byte
	BaseModel
}
