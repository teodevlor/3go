package websystem

type GeoJSONPolygon struct {
	Type        string        `json:"type" binding:"required,eq=Polygon"`
	Coordinates [][][]float64 `json:"coordinates" binding:"required,min=1"`
}

type CreateZoneRequestDto struct {
	Code            string         `json:"code" binding:"required,max=100"`
	Name            string         `json:"name" binding:"required,max=255"`
	PriceMultiplier float64        `json:"price_multiplier" binding:"required,gte=0"`
	IsActive        bool           `json:"is_active"`
	Polygon         GeoJSONPolygon `json:"polygon" binding:"required"`
}
