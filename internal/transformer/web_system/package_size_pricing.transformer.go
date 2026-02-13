package web_system

import (
	dto "go-structure/internal/dto/web_system"
	websystem "go-structure/internal/repository/model/web_system"
)

func ToPackageSizePricingItemDto(r *websystem.PackageSizePricing) dto.PackageSizePricingItemDto {
	if r == nil {
		return dto.PackageSizePricingItemDto{}
	}
	return dto.PackageSizePricingItemDto{
		ID:          r.ID.String(),
		ServiceID:   r.ServiceID.String(),
		PackageSize: r.PackageSize,
		ExtraPrice:  r.ExtraPrice,
		IsActive:    r.IsActive,
	}
}
