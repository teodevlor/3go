package app_driver

import (
	"time"

	"go-structure/internal/common"
)

const (
	DriverDocumentStatusPENDING  = common.DriverDocumentStatusPending
	DriverDocumentStatusAPPROVED = common.DriverDocumentStatusApproved
	DriverDocumentStatusREJECTED = common.DriverDocumentStatusRejected
)

type (
	CreateDriverDocumentRequestDto struct {
		DriverID       string  `json:"driver_id" binding:"required"`
		DocumentTypeID string  `json:"document_type_id" binding:"required"`
		FileUrl        string  `json:"file_url" binding:"required"`
		ExpireAt       *string `json:"expire_at"`
	}

	CreateDriverDocumentResponseDto struct {
		ID             string     `json:"id"`
		DriverID       string     `json:"driver_id"`
		DocumentTypeID string     `json:"document_type_id"`
		FileUrl        string     `json:"file_url"`
		ExpireAt       *string    `json:"expire_at,omitempty"`
		Status         string     `json:"status"`
		RejectReason   *string    `json:"reject_reason,omitempty"`
		VerifiedAt     *time.Time `json:"verified_at,omitempty"`
		VerifiedBy     *string    `json:"verified_by,omitempty"`
		CreatedAt      time.Time  `json:"created_at"`
		UpdatedAt      time.Time  `json:"updated_at"`
	}

	CreateDriverDocumentItemDto struct {
		DocumentTypeID string  `json:"document_type_id" binding:"required"`
		FileUrl        string  `json:"file_url" binding:"required"`
		ExpireAt       *string `json:"expire_at"`
	}

	BulkCreateDriverDocumentsRequestDto struct {
		DriverID string                        `json:"driver_id" binding:"required"`
		Items    []CreateDriverDocumentItemDto `json:"items" binding:"required,min=1,dive"`
	}

	BulkCreateDriverDocumentsResponseDto struct {
		Items []CreateDriverDocumentResponseDto `json:"items"`
	}

	DriverDocumentItemDto struct {
		ID             string                    `json:"id"`
		DriverID       string                    `json:"driver_id"`
		DocumentTypeID string                    `json:"document_type_id"`
		DocumentType   *DriverDocumentTypeItemDto `json:"document_type,omitempty"`
		FileUrl        string                    `json:"file_url"`
		ExpireAt       *string                   `json:"expire_at,omitempty"`
		Status         string                    `json:"status"`
		RejectReason   *string                   `json:"reject_reason,omitempty"`
		VerifiedAt     *time.Time                `json:"verified_at,omitempty"`
		VerifiedBy     *string                   `json:"verified_by,omitempty"`
		CreatedAt      time.Time                 `json:"created_at"`
		UpdatedAt      time.Time                 `json:"updated_at"`
	}

	UpdateDriverDocumentRequestDto struct {
		FileUrl      string  `json:"file_url" binding:"required"`
		ExpireAt     *string `json:"expire_at"`
		Status       string  `json:"status" binding:"required,oneof=PENDING APPROVED REJECTED"`
		RejectReason *string `json:"reject_reason"`
	}

	UpdateDriverDocumentStatusRequestDto struct {
		Status string `json:"status" binding:"required,oneof=PENDING APPROVED REJECTED"`
	}

	BulkUpdateDriverDocumentsRequestDto struct {
		Items []BulkUpdateDriverDocumentItemDto `json:"items" binding:"required,min=1,dive"`
	}

	BulkUpdateDriverDocumentItemDto struct {
		ID           string  `json:"id" binding:"required"`
		FileUrl      *string `json:"file_url"`
		ExpireAt     *string `json:"expire_at"`
		Status       *string `json:"status"` // PENDING, APPROVED, REJECTED
		RejectReason *string `json:"reject_reason"`
	}

	BulkUpdateDriverDocumentsResponseDto struct {
		Items []DriverDocumentItemDto `json:"items"`
	}

	ListDriverDocumentsResponseDto struct {
		Items []DriverDocumentItemDto `json:"items"`
	}
)
