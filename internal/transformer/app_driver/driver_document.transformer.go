package app_driver

import (
	"go-structure/internal/common"
	dto "go-structure/internal/dto/app_driver"
	appdrivermodel "go-structure/internal/repository/model/app_driver"
)

func ToDriverDocumentItemDto(m *appdrivermodel.DriverDocument) dto.DriverDocumentItemDto {
	if m == nil {
		return dto.DriverDocumentItemDto{}
	}
	var verifiedBy *string
	var expireAt *string
	if m.VerifiedBy != nil {
		s := m.VerifiedBy.String()
		verifiedBy = &s
	}
	if m.ExpireAt != nil {
		s := common.FormatDateToYYYYMMDD(*m.ExpireAt)
		expireAt = &s
	}
	return dto.DriverDocumentItemDto{
		ID:             m.ID.String(),
		DriverID:       m.DriverID.String(),
		DocumentTypeID: m.DocumentTypeID.String(),
		FileUrl:        m.FileUrl,
		ExpireAt:       expireAt,
		Status:         m.Status,
		RejectReason:   m.RejectReason,
		VerifiedAt:     m.VerifiedAt,
		VerifiedBy:     verifiedBy,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func ToCreateDriverDocumentResponseDto(m *appdrivermodel.DriverDocument) dto.CreateDriverDocumentResponseDto {
	if m == nil {
		return dto.CreateDriverDocumentResponseDto{}
	}
	var verifiedBy *string
	var expireAt *string
	if m.VerifiedBy != nil {
		s := m.VerifiedBy.String()
		verifiedBy = &s
	}
	if m.ExpireAt != nil {
		s := common.FormatDateToYYYYMMDD(*m.ExpireAt)
		expireAt = &s
	}
	return dto.CreateDriverDocumentResponseDto{
		ID:             m.ID.String(),
		DriverID:       m.DriverID.String(),
		DocumentTypeID: m.DocumentTypeID.String(),
		FileUrl:        m.FileUrl,
		ExpireAt:       expireAt,
		Status:         m.Status,
		RejectReason:   m.RejectReason,
		VerifiedAt:     m.VerifiedAt,
		VerifiedBy:     verifiedBy,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
