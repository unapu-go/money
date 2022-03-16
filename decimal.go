package money

import (
	"math/big"

	"github.com/shopspring/decimal"
)

type Decimal = decimal.Decimal

var (
	DecimalFromFloat32 = decimal.NewFromFloat32
	DecimalFromFloat   = decimal.NewFromFloat
)

func DecimalFromUint(v uint64) Decimal {
	return decimal.NewFromBigInt(new(big.Int).SetUint64(v), 0)
}
