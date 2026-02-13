package web_system

import (
	dto "go-structure/internal/dto/web_system"
	websystem "go-structure/internal/repository/model/web_system"
)

func ToDistancePricingRuleItemDto(r *websystem.DistancePricingRule) dto.DistancePricingRuleItemDto {
	if r == nil {
		return dto.DistancePricingRuleItemDto{}
	}
	return dto.DistancePricingRuleItemDto{
		ID:         r.ID.String(),
		ServiceID:  r.ServiceID.String(),
		FromKm:     r.FromKm,
		ToKm:       r.ToKm,
		PricePerKm: r.PricePerKm,
		IsActive:   r.IsActive,
	}
}
