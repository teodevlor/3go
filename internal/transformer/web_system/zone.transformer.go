package web_system

import (
	dto "go-structure/internal/dto/web_system"
	"go-structure/internal/repository/model"
)

// ToZoneItemDto chuyển model.Zone sang DTO dùng cho Get/List response.
func ToZoneItemDto(z *model.Zone) dto.ZoneItemDto {
	if z == nil {
		return dto.ZoneItemDto{}
	}
	return dto.ZoneItemDto{
		ID:              z.ID.String(),
		Code:            z.Code,
		Name:            z.Name,
		PriceMultiplier: z.PriceMultiplier,
		Polygon:         z.Polygon,
		IsActive:        z.IsActive,
	}
}
