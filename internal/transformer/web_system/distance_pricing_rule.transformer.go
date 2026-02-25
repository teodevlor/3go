package web_system

import (
	dto "go-structure/internal/dto/web_system"
	websystem "go-structure/internal/repository/model/web_system"
)

// ToDistancePricingRuleItemDtoWithService map rule sang DTO và gắn thêm service (name, code) nếu có.
func ToDistancePricingRuleItemDtoWithService(r *websystem.DistancePricingRule, svc *websystem.Service) dto.DistancePricingRuleItemDto {
	if r == nil {
		return dto.DistancePricingRuleItemDto{}
	}
	out := dto.DistancePricingRuleItemDto{
		ID:         r.ID.String(),
		ServiceID:  r.ServiceID.String(),
		FromKm:     r.FromKm,
		ToKm:       r.ToKm,
		PricePerKm: r.PricePerKm,
		IsActive:   r.IsActive,
	}
	if svc != nil {
		out.Service = &dto.DistancePricingRuleServiceDto{
			Name: svc.Name,
			Code: svc.Code,
		}
	}
	return out
}
