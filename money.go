package money

import (
	"database/sql/driver"
	"reflect"

	"github.com/moisespsena-go/aorm"
	"github.com/shopspring/decimal"
)

type (
	Money struct {
		Decimal
	}

	key string

	Signal int8
)

const (
	cfgKey key = "money"

	NEGATIVE Signal = -1
	ALL      Signal = 0
	POSITIVE Signal = 1
)

var (
	Zero = Money{Decimal: decimal.Zero}
	Cent = Money{Decimal: decimal.NewFromInt(1).Div(decimal.NewFromInt(100))}
)

func NewZero() *Money {
	return &Money{Decimal: decimal.Zero}
}

func New(decimal ...Decimal) Money {
	if len(decimal) == 0 {
		return Money{}
	}
	return Money{Decimal: decimal[0]}
}

func (Money) SqlZero() aorm.RawQuery {
	return "0"
}

func (this Money) ToDecimal() Decimal {
	return this.Decimal
}

func (this Money) Module() Money {
	if this.IsNegative() {
		return this.Neg()
	}
	return this
}

func (this Money) IsZero() bool {
	return this.Decimal.IsZero()
}

func (this *Money) IsNil() bool {
	if this == nil {
		return true
	}

	return reflect.ValueOf(&this.Decimal).Elem().FieldByName("value").IsNil()
}

func (this *Money) PrimaryGoValue() interface{} {
	return this.Float()
}

func (this Money) Equal(to Money) bool {
	return this.Decimal.Equal(to.Decimal)
}

func (this Money) Neg() Money {
	return New(this.Decimal.Neg())
}

func (this Money) Sub(v Money) Money {
	return New(this.Decimal.Sub(v.Decimal))
}

func (this Money) Add(v Money) Money {
	return New(this.Decimal.Add(v.Decimal))
}

func (this Money) AddF(v float64) Money {
	return New(this.Decimal.Add(decimal.NewFromFloat(v)))
}

func (this Money) Mul(v Money) Money {
	return New(this.Decimal.Mul(v.Decimal))
}

func (this Money) MulF(v float64) Money {
	return New(this.Decimal.Mul(decimal.NewFromFloat(v)))
}

func (this Money) Div(v Money) Money {
	return New(this.Decimal.Div(v.Decimal))
}

func (this Money) DivF(v float64) Money {
	return New(this.Decimal.Div(decimal.NewFromFloat(v)))
}

func (this Money) DiscountsOfPctF(pct float32) Money {
	if pct == 0 {
		return this
	}
	this.Decimal = this.Decimal.Sub(this.Decimal.Mul(decimal.NewFromFloat32(pct)).Div(decimal.NewFromInt32(100)))
	return this
}

func (this Money) DiscountsOfPctD(pct Decimal) Money {
	if pct.IsZero() {
		return this
	}
	this.Decimal = this.Decimal.Sub(this.Decimal.Mul(pct).Div(decimal.NewFromInt32(100)))
	return this
}

func (this Money) PercentualF(pct float32) Money {
	if pct == 0 {
		return this
	}
	this.Decimal = this.Decimal.Mul(decimal.NewFromFloat32(pct)).Div(decimal.NewFromInt32(100))
	return this
}

func (this Money) PercentualD(pct Decimal) Money {
	if pct.IsZero() {
		return this
	}
	this.Decimal = this.Decimal.Mul(pct).Div(decimal.NewFromInt32(100))
	return this
}

func (this Money) Fix() Money {
	this.Decimal = this.Decimal.Truncate(2)
	return this
}

func (this *Money) PFix() {
	this.Decimal = this.Decimal.Truncate(2)
}

func (this *Money) Do(f func(v Money) Money, v Money) {
	this.Decimal = f(v).Decimal
}

func (this *Money) DoF(f func(v float64) Money, v float64) {
	this.Decimal = f(v).Decimal
}

func (this *Money) DoD(f func(v Decimal) Money, v Decimal) {
	this.Decimal = f(v).Decimal
}

// Scan implements the sql.Scanner interface for database deserialization.
func (d *Money) Scan(value interface{}) error {
	if value == nil {
		d.Decimal = Decimal{}
		return nil
	}
	switch t := value.(type) {
	case Money:
		*d = t
		return nil
	default:
		return d.Decimal.Scan(value)
	}
}

func (this *Money) Value() (driver.Value, error) {
	if this == nil {
		return nil, nil
	}
	return this.Decimal.Truncate(2).Value()
}

func (this *Money) SetZero() {
	*this = Money{}
}

func (this *Money) Float() float64 {
	v, _ := this.Decimal.Float64()
	return v
}

func (this *Money) Int() int64 {
	return this.BigInt().Int64()
}

func (this *Money) Parcel(qnt uint) (parts []Money) {
	parts = make([]Money, qnt)
	val := this.Decimal.Div(decimal.NewFromInt(int64(qnt))).Truncate(2)

	for i := range parts {
		parts[i] = Money{Decimal: val}
	}

	mod := this.Decimal.Sub(val.Mul(decimal.NewFromInt(int64(qnt))))

	if !mod.IsZero() {
		// one decimal unit
		unit := decimal.New(1, mod.Exponent())

		for i := range parts[0:int(mod.Coefficient().Int64())] {
			parts[i].Decimal = parts[i].Decimal.Add(unit)
		}
	}
	return
}

func FromFloat(f float64) Money {
	D := decimal.NewFromFloat(f).Truncate(2)
	return Money{Decimal: D}
}

func FromInt(v int64, exp ...int32) Money {
	if len(exp) > 0 {
		return Money{Decimal: decimal.New(v, exp[0])}
	}
	return Money{Decimal: decimal.NewFromInt(v)}
}

func FromUint(v uint64) Money {
	return Money{DecimalFromUint(v)}
}

func toFixed(num float64, precision int32) float64 {
	if num == 0 {
		return 0
	}
	num, _ = decimal.NewFromFloat(num).Truncate(precision).Float64()
	return num
}

type Config struct {
	Sig Signal
}
