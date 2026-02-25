package app_driver

import (
	dto "go-structure/internal/dto/app_driver"
	appdrivermodel "go-structure/internal/repository/model/app_driver"
)

func ToDriverDocumentTypeItemDto(m *appdrivermodel.DriverDocumentType) dto.DriverDocumentTypeItemDto {
	if m == nil {
		return dto.DriverDocumentTypeItemDto{}
	}
	var serviceID *string
	if m.ServiceID != nil {
		s := m.ServiceID.String()
		serviceID = &s
	}
	return dto.DriverDocumentTypeItemDto{
		ID:                m.ID.String(),
		Code:              m.Code,
		Name:              m.Name,
		Description:       m.Description,
		IsRequired:        m.IsRequired,
		RequireExpireDate: m.RequireExpireDate,
		ServiceID:         serviceID,
		IsActive:          m.IsActive,
	}
}

func ToCreateDriverDocumentTypeResponseDto(m *appdrivermodel.DriverDocumentType) dto.CreateDriverDocumentTypeResponseDto {
	if m == nil {
		return dto.CreateDriverDocumentTypeResponseDto{}
	}
	var serviceID *string
	if m.ServiceID != nil {
		s := m.ServiceID.String()
		serviceID = &s
	}
	return dto.CreateDriverDocumentTypeResponseDto{
		ID:                m.ID.String(),
		Code:              m.Code,
		Name:              m.Name,
		Description:       m.Description,
		IsRequired:        m.IsRequired,
		RequireExpireDate: m.RequireExpireDate,
		ServiceID:         serviceID,
		IsActive:          m.IsActive,
	}
}
