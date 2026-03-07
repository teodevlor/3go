package websystem

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-structure/internal/common"
)

// condition type
const (
	ConditionTypeTimeWindow = common.ConditionTypeTimeWindow
	ConditionTypeWeather    = common.ConditionTypeWeather
	ConditionTypeTraffic    = common.ConditionTypeTraffic
	ConditionTypeHoliday    = common.ConditionTypeHoliday
)

type (
	timeWindowConfig struct {
		From string   `json:"from"`
		To   string   `json:"to"`
		Days []string `json:"days"`
	}
	weatherConfig struct {
		Rain   *bool          `json:"rain,omitempty"`
		RainMM *weatherMetric `json:"rain_mm,omitempty"`
	}
	weatherMetric struct {
		Operator string  `json:"operator"`
		Value    float64 `json:"value"`
	}
	trafficConfig struct {
		Level string `json:"level"`
	}
	holidayConfig struct {
		Dates []string `json:"dates"`
	}
)

func ValidateConditionConfig(conditionType string, raw []byte) error {
	switch conditionType {
	case ConditionTypeTimeWindow:
		var cfg timeWindowConfig
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return err
		}
		if cfg.From == "" || cfg.To == "" {
			return fmt.Errorf("from and to are required")
		}
		if !isValidTimeHHMM(cfg.From) || !isValidTimeHHMM(cfg.To) {
			return fmt.Errorf("from and to must be in HH:mm format")
		}
		if len(cfg.Days) == 0 {
			return fmt.Errorf("days is required")
		}
		for _, d := range cfg.Days {
			if !isValidDayOfWeek(d) {
				return fmt.Errorf("invalid day value: %s", d)
			}
		}
		return nil
	case ConditionTypeWeather:
		var cfg weatherConfig
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return err
		}
		if cfg.Rain == nil && cfg.RainMM == nil {
			return fmt.Errorf("either rain or rain_mm must be provided")
		}
		if cfg.RainMM != nil {
			if !isValidWeatherOperator(cfg.RainMM.Operator) {
				return fmt.Errorf("invalid weather operator: %s", cfg.RainMM.Operator)
			}
		}
		return nil
	case ConditionTypeTraffic:
		var cfg trafficConfig
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return err
		}
		if !isValidTrafficLevel(cfg.Level) {
			return fmt.Errorf("invalid traffic level: %s", cfg.Level)
		}
		return nil
	case ConditionTypeHoliday:
		var cfg holidayConfig
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return err
		}
		if len(cfg.Dates) == 0 {
			return fmt.Errorf("dates không được rỗng, cần ít nhất một ngày (định dạng YYYY-MM-DD)")
		}
		for i, d := range cfg.Dates {
			if !isValidDateYYYYMMDD(d) {
				return fmt.Errorf("ngày không hợp lệ tại vị trí %d: %q, cần định dạng YYYY-MM-DD", i+1, d)
			}
		}
		return nil
	default:
		return fmt.Errorf("invalid condition type: %s", conditionType)
	}
}

func isValidTimeHHMM(s string) bool {
	_, err := time.Parse("15:04", s)
	return err == nil
}

func isValidDateYYYYMMDD(s string) bool {
	_, err := time.Parse(common.YYYYMMDDFormat, s)
	return err == nil
}

func isValidDayOfWeek(s string) bool {
	switch strings.ToLower(s) {
	case "mon", "tue", "wed", "thu", "fri", "sat", "sun":
		return true
	default:
		return false
	}
}

func isValidWeatherOperator(operator string) bool {
	switch operator {
	case "=", ">", ">=", "<", "<=":
		return true
	default:
		return false
	}
}

func isValidTrafficLevel(level string) bool {
	switch strings.ToLower(level) {
	case "low", "medium", "high":
		return true
	default:
		return false
	}
}
