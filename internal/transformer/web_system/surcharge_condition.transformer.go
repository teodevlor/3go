package web_system

import (
	"encoding/json"

	dto "go-structure/internal/dto/web_system"
	websystem "go-structure/internal/repository/model/web_system"
)

func ToSurchargeConditionItemDto(c *websystem.SurchargeCondition) dto.SurchargeConditionItemDto {
	if c == nil {
		return dto.SurchargeConditionItemDto{}
	}
	return dto.SurchargeConditionItemDto{
		ID:            c.ID.String(),
		Code:          c.Code,
		ConditionType: c.ConditionType,
		Config:        json.RawMessage(c.Config),
		IsActive:      c.IsActive,
	}
}

