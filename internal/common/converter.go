package common

import (
	"math/big"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// NumericToFloat64 chuyển kiểu pgtype.Numeric (số thập phân trong PostgreSQL)
// sang float64 trong Go, trả về 0 nếu giá trị không hợp lệ.
func NumericToFloat64(n pgtype.Numeric) float64 {
	if n.NaN || n.Int == nil {
		return 0
	}

	r := new(big.Rat).SetInt(n.Int)

	if n.Exp < 0 {
		denom := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-n.Exp)), nil)
		r = r.Quo(r, new(big.Rat).SetInt(denom))
	} else if n.Exp > 0 {
		mult := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n.Exp)), nil)
		r = r.Mul(r, new(big.Rat).SetInt(mult))
	}

	f, _ := r.Float64()
	return f
}

// Float64ToNumeric chuyển float64 trong Go sang pgtype.Numeric
// để lưu trữ ngược lại vào PostgreSQL.
func Float64ToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(strconv.FormatFloat(f, 'f', -1, 64))
	return n
}

// NullableTextToString chuyển pgtype.Text có thể NULL sang string.
// Nếu giá trị NULL (không hợp lệ) sẽ trả về chuỗi rỗng.
func NullableTextToString(v pgtype.Text) string {
	if !v.Valid {
		return ""
	}
	return v.String
}

// NullableInt4ToInt32 chuyển pgtype.Int4 có thể NULL sang int32.
// Nếu giá trị NULL sẽ trả về 0.
func NullableInt4ToInt32(v pgtype.Int4) int32 {
	if !v.Valid {
		return 0
	}
	return v.Int32
}

// NullableDateToTimePtr chuyển pgtype.Date có thể NULL sang *time.Time.
// Nếu giá trị NULL sẽ trả về nil, ngược lại trả về con trỏ tới time.Time.
func NullableDateToTimePtr(v pgtype.Date) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
}

// NullableTimestamptzToTimePtr chuyển pgtype.Timestamptz có thể NULL sang *time.Time.
// Nếu giá trị NULL sẽ trả về nil, ngược lại trả về con trỏ tới time.Time.
func NullableTimestamptzToTimePtr(v pgtype.Timestamptz) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
}

// NullableNumericToFloat64 chuyển pgtype.Numeric có thể NULL sang float64.
// Nếu giá trị NULL hoặc không convert được sẽ trả về 0.
func NullableNumericToFloat64(v pgtype.Numeric) float64 {
	if !v.Valid {
		return 0
	}
	f, err := v.Float64Value()
	if err != nil || !f.Valid {
		return 0
	}
	return f.Float64
}

// NullableTextToStringPtr chuyển pgtype.Text có thể NULL sang *string.
// Nếu giá trị NULL sẽ trả về nil, ngược lại trả về con trỏ tới string.
func NullableTextToStringPtr(v pgtype.Text) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}

// NullableUUIDToUUIDPtr chuyển pgtype.UUID có thể NULL sang *uuid.UUID.
// Nếu giá trị NULL sẽ trả về nil, ngược lại trả về con trỏ tới uuid.UUID.
func NullableUUIDToUUIDPtr(v pgtype.UUID) *uuid.UUID {
	if !v.Valid {
		return nil
	}
	uid := uuid.UUID(v.Bytes)
	return &uid
}
