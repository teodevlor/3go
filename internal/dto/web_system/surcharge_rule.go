package web_system

import (
	"encoding/json"

	dto_common "go-structure/internal/dto/common"
)

type (
	CreateSurchargeRuleRequestDto struct {
		ServiceID string          `json:"service_id" binding:"required,uuid"`
		ZoneID    string          `json:"zone_id" binding:"required,uuid"`
		Amount    float64         `json:"amount" binding:"required,gte=0"`
		Unit      string          `json:"unit" binding:"required"` // 'percent' | 'fixed'
		Priority  int             `json:"priority" binding:"gte=0"`
		Condition json.RawMessage `json:"condition"` // TODO: sử dụng system_surcharge_rule_conditions
		IsActive  bool            `json:"is_active"`
	}

	UpdateSurchargeRuleRequestDto struct {
		ServiceID string          `json:"service_id" binding:"required,uuid"`
		ZoneID    string          `json:"zone_id" binding:"required,uuid"`
		Amount    float64         `json:"amount" binding:"required,gte=0"`
		Unit      string          `json:"unit" binding:"required"`
		Priority  int             `json:"priority" binding:"gte=0"`
		Condition json.RawMessage `json:"condition"` // TODO: sử dụng system_surcharge_rule_conditions
		IsActive  bool            `json:"is_active"`
	}

	SurchargeRuleItemDto struct {
		ID        string          `json:"id"`
		ServiceID string          `json:"service_id"`
		ZoneID    string          `json:"zone_id"`
		Amount    float64         `json:"amount"`
		Unit      string          `json:"unit"`
		Priority  int             `json:"priority"`
		Condition json.RawMessage `json:"condition"` // TODO: sử dụng system_surcharge_rule_conditions
		IsActive  bool            `json:"is_active"`
	}

	ListSurchargeRulesResponseDto struct {
		Items      []SurchargeRuleItemDto  `json:"items"`
		Pagination dto_common.PaginationMeta `json:"pagination"`
	}
)
