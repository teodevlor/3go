package web_system

import "encoding/json"

type (
	CreateSurchargeRuleRequestDto struct {
		ServiceID     string          `json:"service_id" binding:"required,uuid"`
		ZoneID        string          `json:"zone_id" binding:"required,uuid"`
		SurchargeType string          `json:"surcharge_type" binding:"required"`
		Amount        float64         `json:"amount" binding:"required,gte=0"`
		Unit          string          `json:"unit" binding:"required"` // 'percent' | 'fixed'
		Condition     json.RawMessage `json:"condition"`
		IsActive      bool            `json:"is_active"`
	}

	UpdateSurchargeRuleRequestDto struct {
		ServiceID     string          `json:"service_id" binding:"required,uuid"`
		ZoneID        string          `json:"zone_id" binding:"required,uuid"`
		SurchargeType string          `json:"surcharge_type" binding:"required"`
		Amount        float64         `json:"amount" binding:"required,gte=0"`
		Unit          string          `json:"unit" binding:"required"`
		Condition     json.RawMessage `json:"condition"`
		IsActive      bool            `json:"is_active"`
	}

	SurchargeRuleItemDto struct {
		ID            string          `json:"id"`
		ServiceID     string          `json:"service_id"`
		ZoneID        string          `json:"zone_id"`
		SurchargeType string          `json:"surcharge_type"`
		Amount        float64         `json:"amount"`
		Unit          string          `json:"unit"`
		Condition     json.RawMessage `json:"condition"`
		IsActive      bool            `json:"is_active"`
	}
)
