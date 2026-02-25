package appdriver

import (
	"time"

	"github.com/google/uuid"
)

type DriverDocument struct {
	ID             uuid.UUID  `json:"id"`
	DriverID       uuid.UUID  `json:"driver_id"`
	DocumentTypeID uuid.UUID  `json:"document_type_id"`
	FileUrl        string     `json:"file_url"`
	ExpireAt       *time.Time `json:"expire_at,omitempty"`
	Status         string     `json:"status"`
	RejectReason   *string    `json:"reject_reason,omitempty"`
	VerifiedAt     *time.Time `json:"verified_at,omitempty"`
	VerifiedBy     *uuid.UUID `json:"verified_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
