package web_system

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	dto_common "go-structure/internal/dto/common"
)

type (
	CreateSurchargeConditionRequestDto struct {
		Code          string          `json:"code" binding:"required"`
		ConditionType string          `json:"condition_type" binding:"required,oneof=time_window weather traffic holiday"`
		Config        json.RawMessage `json:"config" binding:"required"`
		IsActive      bool            `json:"is_active"`
	}

	UpdateSurchargeConditionRequestDto struct {
		Code          string          `json:"code" binding:"required"`
		ConditionType string          `json:"condition_type" binding:"required,oneof=time_window weather traffic holiday"`
		Config        json.RawMessage `json:"config" binding:"required"`
		IsActive      bool            `json:"is_active"`
	}

	SurchargeConditionItemDto struct {
		ID            string          `json:"id"`
		Code          string          `json:"code"`
		ConditionType string          `json:"condition_type"`
		Config        json.RawMessage `json:"config"`
		IsActive      bool            `json:"is_active"`
	}

	ListSurchargeConditionsResponseDto struct {
		Items      []SurchargeConditionItemDto `json:"items"`
		Pagination dto_common.PaginationMeta   `json:"pagination"`
	}
)

type ConditionType string

const (
	ConditionTimeWindow ConditionType = "time_window"
	ConditionWeather    ConditionType = "weather"
	ConditionTraffic    ConditionType = "traffic"
	ConditionHoliday    ConditionType = "holiday"
)

type TimeWindowConfig struct {
	From string   `json:"from"` // HH:mm
	To   string   `json:"to"`   // HH:mm
	Days []string `json:"days"` // mon-sun
}

const TimeWindowTemplate = `{
	"from": "17:00",
	"to": "19:00",
	"days": ["mon","tue","wed","thu","fri"]
}`

type WeatherConfig struct {
	Rain   *bool          `json:"rain,omitempty"`
	RainMM *WeatherMetric `json:"rain_mm,omitempty"`
}

type WeatherMetric struct {
	Operator string  `json:"operator"` // =, >, >=, <, <=
	Value    float64 `json:"value"`
}

const WeatherRainTemplate = `{
	"rain": true
}`

const WeatherRainMMTemplate = `{
	"rain_mm": {
		"operator": ">=",
		"value": 10
	}
}`

type TrafficConfig struct {
	Level string `json:"level"` // low, medium, high
}

const TrafficTemplate = `{
	"level": "high"
}`

type HolidayConfig struct {
	HolidayCode string `json:"holiday_code"`
}

const HolidayTemplate = `{
	"holiday_code": "TET"
}`

func ValidateConditionConfig(t ConditionType, raw []byte) error {
	switch t {

	case ConditionTimeWindow:
		var cfg TimeWindowConfig
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

	case ConditionWeather:
		var cfg WeatherConfig
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

	case ConditionTraffic:
		var cfg TrafficConfig
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return err
		}

		if !isValidTrafficLevel(cfg.Level) {
			return fmt.Errorf("invalid traffic level: %s", cfg.Level)
		}

		return nil

	case ConditionHoliday:
		var cfg HolidayConfig
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return err
		}

		if cfg.HolidayCode == "" {
			return fmt.Errorf("holiday_code is required")
		}

		return nil

	default:
		return fmt.Errorf("invalid condition type")
	}
}

func isValidTimeHHMM(s string) bool {
	_, err := time.Parse("15:04", s)
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

func isValidWeatherOperator(op string) bool {
	switch op {
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
