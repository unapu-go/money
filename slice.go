package money

import "github.com/shopspring/decimal"

type Slice []Money

// DiscontsPct Apply percentual disconts to credits (positive values).
func (this Slice) DiscontsPct(pct decimal.Decimal) (res Slice, sumDiscont Money) {
	if this == nil {
		return
	}
	res = make(Slice, len(this))
	copy(res, this)

	var (
		diff decimal.Decimal
	)

	for i := range res {
		if res[i].IsPositive() {
			diff = res[i].Decimal.Mul(pct)
			res[i].Decimal = res[i].Decimal.Sub(diff)
			sumDiscont.Decimal = sumDiscont.Decimal.Add(diff)
		}
	}

	return
}
