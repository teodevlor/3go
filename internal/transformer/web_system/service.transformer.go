package web_system

import (
	dto "go-structure/internal/dto/web_system"
	websystem "go-structure/internal/repository/model/web_system"
)

func ToServiceItemDto(s *websystem.Service, zoneIDs []string) dto.ServiceItemDto {
	if s == nil {
		return dto.ServiceItemDto{}
	}
	return dto.ServiceItemDto{
		ID:        s.ID.String(),
		Code:      s.Code,
		Name:      s.Name,
		BasePrice: s.BasePrice,
		MinPrice:  s.MinPrice,
		IsActive:  s.IsActive,
		ZoneIDs:   zoneIDs,
	}
}
