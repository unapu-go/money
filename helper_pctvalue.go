package money

type (
	PctValueD struct {
		Counts Decimal
		// Pct value between 0 to 1
		Pct Decimal
	}

	PctValue struct {
		Counts Decimal
		// Pct value between 0 to 1
		Pct  float32
		Data interface{}
	}

	PctValueResult struct {
		Value Money
		Data  interface{}
	}
)

func UnitValuesFromTotalD(total Money, refValueCount Decimal, otherValues ...*PctValueD) []Decimal {
	z := refValueCount
	for _, v := range otherValues {
		z = z.Add(v.Counts.Mul(v.Pct))
	}
	vRef := total.Decimal.Div(z)
	res := make([]Decimal, len(otherValues)+1)
	res[0] = vRef
	for i, v := range otherValues {
		res[i+1] = vRef.Mul(v.Pct)
	}
	return res
}

func UnitValuesFromTotal(total Money, values ...*PctValue) (_ []PctValueResult, remainder Money) {
	var (
		baseValueCount = values[0].Counts
		z              = baseValueCount
	)

	for _, v := range values[1:] {
		z = z.Add(v.Counts.Mul(DecimalFromFloat32(v.Pct)))
	}

	var (
		baseV = total.Decimal.Div(z).Truncate(2)
		tot   = baseV.Mul(baseValueCount)
		res   = make([]PctValueResult, len(values))
	)

	res[0].Data = values[0].Data
	res[0].Value.Decimal = baseV

	for i, v := range values[1:] {
		v2 := baseV.Mul(DecimalFromFloat32(v.Pct)).Truncate(2)
		res[i+1] = PctValueResult{Value: Money{Decimal: v2}, Data: v.Data}
		tot = tot.Add(v2.Mul(v.Counts))
	}

	remainder = Money{total.Decimal.Sub(tot)}

	if !remainder.IsZero() {
		d, r := remainder.Decimal.QuoRem(baseValueCount, 2)
		remainder.Decimal = r
		res[0].Value.Decimal = res[0].Value.Decimal.Add(d)
	}

	return res, remainder
}
