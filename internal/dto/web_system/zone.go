package web_system

import (
	"encoding/json"
	"errors"

	dto_common "go-structure/internal/dto/common"
)

// UnmarshalJSON chấp nhận cả hai cách gửi:
//   - Chuẩn: [ [ [lng,lat], [lng,lat], ... ] ] (mảng ring, mỗi ring là mảng điểm)
//   - Một ring phẳng: [ [lng,lat], [lng,lat], ... ] (tự bọc thành một ring)
type PolygonCoordinates [][][]float64

func (c *PolygonCoordinates) UnmarshalJSON(data []byte) error {
	var raw json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var standard [][][]float64
	if err := json.Unmarshal(raw, &standard); err == nil && len(standard) > 0 {
		*c = standard
		return nil
	}

	var singleRing [][]float64
	if err := json.Unmarshal(raw, &singleRing); err != nil {
		return err
	}
	*c = [][][]float64{singleRing}
	return nil
}

type GeoJSONPolygon struct {
	Type        string             `json:"type" binding:"required,eq=Polygon"`
	Coordinates PolygonCoordinates `json:"coordinates" binding:"required,min=1"`
}

func (p GeoJSONPolygon) Validate() error {
	if p.Type != "Polygon" {
		return errors.New("polygon.type must be Polygon")
	}
	if len(p.Coordinates) == 0 {
		return errors.New("polygon.coordinates must have at least one ring")
	}
	ring := p.Coordinates[0]
	if len(ring) < 4 {
		return errors.New("polygon ring must have at least 4 points")
	}
	first, last := ring[0], ring[len(ring)-1]
	if len(first) != 2 || len(last) != 2 {
		return errors.New("each coordinate must have [lng, lat]")
	}
	if first[0] != last[0] || first[1] != last[1] {
		return errors.New("polygon ring must be closed (first point equals last point)")
	}
	for _, point := range ring {
		if len(point) != 2 {
			return errors.New("each coordinate must have exactly 2 values [lng, lat]")
		}
		lng, lat := point[0], point[1]
		if lng < -180 || lng > 180 {
			return errors.New("longitude must be in range [-180, 180]")
		}
		if lat < -90 || lat > 90 {
			return errors.New("latitude must be in range [-90, 90]")
		}
	}
	return nil
}

type CreateZoneRequestDto struct {
	Code            string         `json:"code" binding:"required,max=100"`
	Name            string         `json:"name" binding:"required,max=255"`
	PriceMultiplier float64        `json:"price_multiplier" binding:"required,gte=0"`
	IsActive        bool           `json:"is_active"`
	Polygon         GeoJSONPolygon `json:"polygon" binding:"required"`
}

type CreateZoneResponseDto struct {
	ID              string  `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	PriceMultiplier float64 `json:"price_multiplier"`
	Polygon         string  `json:"polygon"`
	IsActive        bool    `json:"is_active"`
}

type UpdateZoneRequestDto struct {
	Code            string         `json:"code" binding:"required,max=100"`
	Name            string         `json:"name" binding:"required,max=255"`
	PriceMultiplier float64        `json:"price_multiplier" binding:"required,gte=0"`
	IsActive        bool           `json:"is_active"`
	Polygon         GeoJSONPolygon `json:"polygon" binding:"required"`
}

type ZoneItemDto struct {
	ID              string  `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	PriceMultiplier float64 `json:"price_multiplier"`
	Polygon         string  `json:"polygon"`
	IsActive        bool    `json:"is_active"`
}

type ListZonesResponseDto struct {
	Items      []ZoneItemDto  `json:"items"`
	Pagination dto_common.PaginationMeta `json:"pagination"`
}
