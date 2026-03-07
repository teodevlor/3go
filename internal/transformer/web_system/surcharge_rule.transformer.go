package web_system

import (
	dto "go-structure/internal/dto/web_system"
	"go-structure/internal/repository/model"
	websystem "go-structure/internal/repository/model/web_system"

	"github.com/google/uuid"
)

func ToSurchargeRuleItemDto(r *websystem.SurchargeRule) dto.SurchargeRuleItemDto {
	if r == nil {
		return dto.SurchargeRuleItemDto{}
	}
	return dto.SurchargeRuleItemDto{
		ID:           r.ID.String(),
		Amount:       r.Amount,
		Unit:         r.Unit,
		Priority:     int(r.Priority),
		IsActive:     r.IsActive,
		ConditionIDs: []string{},
		Conditions:   []dto.SurchargeConditionItemDto{},
	}
}

func ToSurchargeRuleItemDtoWithConditions(
	r *websystem.SurchargeRule,
	conditionIDs []uuid.UUID,
	conditions []*websystem.SurchargeCondition,
	service *websystem.Service,
	zone *model.Zone,
) dto.SurchargeRuleItemDto {
	item := ToSurchargeRuleItemDto(r)

	if len(conditionIDs) > 0 {
		item.ConditionIDs = make([]string, 0, len(conditionIDs))
		for _, id := range conditionIDs {
			item.ConditionIDs = append(item.ConditionIDs, id.String())
		}
	}

	if len(conditions) > 0 {
		item.Conditions = make([]dto.SurchargeConditionItemDto, 0, len(conditions))
		for _, c := range conditions {
			if c == nil {
				continue
			}
			item.Conditions = append(item.Conditions, ToSurchargeConditionItemDto(c))
		}
	}

	if service != nil {
		item.Service = &dto.SurchargeRuleServiceDto{
			ID:   service.ID.String(),
			Name: service.Name,
			Code: service.Code,
		}
	}

	if zone != nil {
		item.Zone = &dto.SurchargeRuleZoneDto{
			ID:   zone.ID.String(),
			Name: zone.Name,
			Code: zone.Code,
		}
	}

	return item
}
