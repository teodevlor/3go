package web_system

type CreateDistancePricingRuleRequestDto struct {
	ServiceID  string  `json:"service_id" binding:"required,uuid"`
	FromKm     float64 `json:"from_km" binding:"required,gte=0"`
	ToKm       float64 `json:"to_km" binding:"required,gte=0"`
	PricePerKm float64 `json:"price_per_km" binding:"required,gte=0"`
	IsActive   bool    `json:"is_active"`
}

type CreateDistancePricingRuleResponseDto struct {
	ID         string  `json:"id"`
	ServiceID  string  `json:"service_id"`
	FromKm     float64 `json:"from_km"`
	ToKm       float64 `json:"to_km"`
	PricePerKm float64 `json:"price_per_km"`
	IsActive   bool    `json:"is_active"`
}

type UpdateDistancePricingRuleRequestDto struct {
	ServiceID  string  `json:"service_id" binding:"required,uuid"`
	FromKm     float64 `json:"from_km" binding:"required,gte=0"`
	ToKm       float64 `json:"to_km" binding:"required,gte=0"`
	PricePerKm float64 `json:"price_per_km" binding:"required,gte=0"`
	IsActive   bool    `json:"is_active"`
}

type DistancePricingRuleItemDto struct {
	ID         string  `json:"id"`
	ServiceID  string  `json:"service_id"`
	FromKm     float64 `json:"from_km"`
	ToKm       float64 `json:"to_km"`
	PricePerKm float64 `json:"price_per_km"`
	IsActive   bool    `json:"is_active"`
}
