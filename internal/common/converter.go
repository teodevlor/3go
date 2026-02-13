package common

import (
	"math/big"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

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

func Float64ToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(strconv.FormatFloat(f, 'f', -1, 64))
	return n
}
