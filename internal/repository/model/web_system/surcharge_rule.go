package websystem

import (
	"fmt"

	"go-structure/internal/common"
	"go-structure/internal/repository/model"

	"github.com/google/uuid"
)

const (
	SurchargeUnitPercent = common.SurchargeUnitPercent
	SurchargeUnitFixed   = common.SurchargeUnitFixed
)

type SurchargeRule struct {
	ID        uuid.UUID `json:"id"`
	ServiceID uuid.UUID `json:"service_id"`
	ZoneID    uuid.UUID `json:"zone_id"`
	Amount    float64   `json:"amount"`
	Unit      string    `json:"unit"`
	IsActive  bool      `json:"is_active"`
	Priority  int32     `json:"priority"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedBy uuid.UUID `json:"updated_by"`
	model.BaseModel
}

func (surchargeRule *SurchargeRule) ValidateSurchargeRule() error {
	if surchargeRule == nil {
		return fmt.Errorf("surcharge rule is nil")
	}
	if surchargeRule.Amount < 0 {
		return fmt.Errorf("amount phải >= 0")
	}
	if surchargeRule.Unit == "" {
		return fmt.Errorf("unit không được để trống")
	}
	if surchargeRule.Unit != SurchargeUnitPercent && surchargeRule.Unit != SurchargeUnitFixed {
		return fmt.Errorf("unit phải là %q hoặc %q", SurchargeUnitPercent, SurchargeUnitFixed)
	}
	if surchargeRule.Priority < 0 {
		return fmt.Errorf("priority phải >= 0")
	}
	return nil
}
