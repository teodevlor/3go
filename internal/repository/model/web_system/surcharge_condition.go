package websystem

import "github.com/google/uuid"

type SurchargeCondition struct {
	ID            uuid.UUID `json:"id"`
	Code          string    `json:"code"`
	ConditionType string    `json:"condition_type"`
	Config        []byte    `json:"config"`
	IsActive      bool      `json:"is_active"`
}

