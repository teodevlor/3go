package web_system

import (
	dto "go-structure/internal/dto/web_system"
	websystem "go-structure/internal/repository/model/web_system"
)

func ToPermissionItemDto(p *websystem.Permission) dto.PermissionItemDto {
	if p == nil {
		return dto.PermissionItemDto{}
	}
	return dto.PermissionItemDto{
		ID:          p.ID.String(),
		Resource:    p.Resource,
		Action:      p.Action,
		Code:        p.Code,
		Name:        p.Name,
		Description: p.Description,
	}
}
