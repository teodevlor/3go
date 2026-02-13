package web_system

import (
	dto "go-structure/internal/dto/web_system"
	websystem "go-structure/internal/repository/model/web_system"
)

func ToSurchargeRuleItemDto(r *websystem.SurchargeRule) dto.SurchargeRuleItemDto {
	if r == nil {
		return dto.SurchargeRuleItemDto{}
	}
	return dto.SurchargeRuleItemDto{
		ID:            r.ID.String(),
		ServiceID:     r.ServiceID.String(),
		ZoneID:        r.ZoneID.String(),
		SurchargeType: r.SurchargeType,
		Amount:        r.Amount,
		Unit:          r.Unit,
		Condition:     r.Condition,
		IsActive:      r.IsActive,
	}
}
