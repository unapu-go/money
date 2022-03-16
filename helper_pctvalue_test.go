package money

import (
	"reflect"
	"testing"
)

func TestUnitValuesFromTotal(t *testing.T) {
	type args struct {
		total         Money
		refValueCount uint64
		otherValues   []*PctValue
	}
	tests := []struct {
		name          string
		args          args
		want          []Money
		wantTotal     Money
		wantRemainder Money
	}{
		{"", args{FromUint(100), 5, []*PctValue{{3, .8}, {2, .7}}},
			[]Money{FromFloat(11.37), FromFloat(9.08), FromFloat(7.95)}, FromFloat(99.99), FromFloat(0.01)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotRemainder := UnitValuesFromTotal(tt.args.total, tt.args.refValueCount, tt.args.otherValues...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnitValuesFromTotal() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(gotRemainder.Decimal.String(), tt.wantRemainder.Decimal.String()) {
				t.Errorf("UnitValuesFromTotal() gotRemainder = %v, want %v", gotRemainder, tt.wantRemainder)
			}
			var tot = got[0].Mul(FromUint(tt.args.refValueCount))
			for i, v := range got[1:] {
				tot.Decimal = tot.Decimal.Add(v.Decimal.Mul(DecimalFromUint(tt.args.otherValues[i].Counts)))
			}
			if !tot.Equal(tt.wantTotal) {
				t.Errorf("UnitValuesFromTotal() total: got = %v, want %v", tot, tt.wantTotal)
			}
			if !tot.Add(tt.wantRemainder).Equal(tt.args.total) {
				t.Errorf("UnitValuesFromTotal() remainder: got = %v, want %v", gotRemainder, tt.wantRemainder)
			}
		})
	}
}

func TestUnitValuesFromTotalD(t *testing.T) {
	type args struct {
		total         Money
		refValueCount Decimal
		otherValues   []*PctValueD
	}
	tests := []struct {
		name string
		args args
		want []Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnitValuesFromTotalD(tt.args.total, tt.args.refValueCount, tt.args.otherValues...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnitValuesFromTotalD() = %v, want %v", got, tt.want)
			}
		})
	}
}
