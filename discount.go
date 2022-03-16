package money

import (
	"github.com/shopspring/decimal"
)

type Discount struct {
	Gives    Money
	Receives Money
}

func (this Discount) To(to []*Money) {
	var (
		g, r       []*Money
		gpct, rpct decimal.Decimal
		gtot       = decimal.Zero
		rtot       = decimal.Zero
	)

	for _, v := range to {
		if v.IsPositive() {
			gtot = gtot.Add(v.Decimal)
			g = append(g, v)
		} else if v.IsNegative() {
			rtot = rtot.Add(v.Decimal)
			r = append(r, v)
		}
	}

	if !gtot.IsZero() {
		gpct, gtot = this.Gives.Decimal.Div(gtot), decimal.Zero
	}

	if !rtot.IsZero() {
		rpct, rtot = this.Receives.Decimal.Div(rtot), decimal.Zero
	}

	for i := range to {
		v := to[i].Decimal
		if v.IsZero() {
			continue
		}
		if v.IsPositive() {
			val := v.Mul(gpct).Truncate(2)
			v = v.Sub(val)
			gtot = gtot.Add(val)
			to[i].Decimal = v
		} else if v.IsNegative() {
			val := v.Mul(rpct).Truncate(2)
			v = v.Add(val)
			rtot = rtot.Add(val)
			to[i].Decimal = v
		}
	}

	if !gtot.IsZero() {
		var (
			diff    = this.Gives.Decimal.Sub(gtot)
			unit    = decimal.New(1, diff.Exponent())
			changed = true
		)

	gdiff_loop:
		for !diff.IsZero() && changed {
			changed = false

			for _, res := range g {
				res.Decimal = res.Decimal.Sub(unit)
				diff = diff.Sub(unit)
				if diff.IsZero() {
					break gdiff_loop
				}
				changed = true
			}
		}
	}

	if !rtot.IsZero() {
		var (
			diff    = this.Receives.Decimal.Sub(rtot)
			unit    = decimal.New(1, diff.Exponent())
			changed = true
		)

	rdiff_loop:
		for !diff.IsZero() && changed {
			changed = false

			for _, res := range r {
				res.Decimal = res.Decimal.Add(unit)
				diff = diff.Sub(unit)
				if diff.IsZero() {
					break rdiff_loop
				}
				changed = true
			}
		}
	}

	return
}
