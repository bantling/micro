package math

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/tuple"
	"github.com/stretchr/testify/assert"
)

func TestOfDecimal_(t *testing.T) {
	assert.Equal(t, tuple.Of2(Decimal{scale: 2, value: 100}, error(nil)), tuple.Of2(OfDecimal(100)))
	assert.Equal(t, tuple.Of2(Decimal{scale: 3, value: -1001}, error(nil)), tuple.Of2(OfDecimal(-1001, 3)))

	assert.Equal(t, Decimal{scale: 2, value: 100}, MustDecimal(100))
	assert.Equal(t, Decimal{scale: 3, value: -1001}, MustDecimal(-1001, 3))

	assert.Equal(
		t,
		tuple.Of2(Decimal{}, fmt.Errorf("The Decimal scale 19 is too large: the value must be <= 18")),
		tuple.Of2(OfDecimal(0, 19)),
	)
	assert.Equal(
		t,
		tuple.Of2(Decimal{}, fmt.Errorf("The Decimal value 1234567890123456789 is too large: the value must be <= 999_999_999_999_999_999")),
		tuple.Of2(OfDecimal(1234567890123456789)),
	)
	assert.Equal(
		t,
		tuple.Of2(Decimal{}, fmt.Errorf("The Decimal value -1234567890123456789 is too small: the value must be >= -999_999_999_999_999_999")),
		tuple.Of2(OfDecimal(-1234567890123456789)),
	)
}

func TestStringToDecimal_(t *testing.T) {
	assert.Equal(t, tuple.Of2(Decimal{scale: 0, value: 100}, error(nil)), tuple.Of2(StringToDecimal("100")))
	assert.Equal(t, tuple.Of2(Decimal{scale: 3, value: -1001}, error(nil)), tuple.Of2(StringToDecimal("-1.001")))

	assert.Equal(
		t,
		tuple.Of2(Decimal{}, fmt.Errorf("The string value 1234567890123456789 is not a valid decimal string")),
		tuple.Of2(StringToDecimal("1234567890123456789")),
	)
	assert.Equal(
		t,
		tuple.Of2(Decimal{}, fmt.Errorf("The string value -1234567890123456789 is not a valid decimal string")),
		tuple.Of2(StringToDecimal("-1234567890123456789")),
	)
}

func TestDecimalString_(t *testing.T) {
	assert.Equal(t, "123", MustDecimal(123, 0).String())
	assert.Equal(t, "-123", MustDecimal(-123, 0).String())
	assert.Equal(t, "12.3", MustDecimal(123, 1).String())
	assert.Equal(t, "1.23", MustDecimal(123, 2).String())
	assert.Equal(t, "0.123", MustDecimal(123, 3).String())
	assert.Equal(t, "0.0123", MustDecimal(123, 4).String())
	assert.Equal(t, "0.00123", MustDecimal(123, 5).String())
	assert.Equal(t, "-0.00123", MustDecimal(-123, 5).String())
}

func TestDecimalPrecision_(t *testing.T) {
	assert.Equal(t, 3, MustDecimal(123, 0).Precision())
	assert.Equal(t, 3, MustDecimal(-123, 0).Precision())
	assert.Equal(t, 5, MustDecimal(12300).Precision())
	assert.Equal(t, 5, MustDecimal(-12345).Precision())
}

func TestDecimalScale_(t *testing.T) {
	assert.Equal(t, uint(0), MustDecimal(123, 0).Scale())
	assert.Equal(t, uint(0), MustDecimal(-123, 0).Scale())
	assert.Equal(t, uint(2), MustDecimal(12300).Scale())
	assert.Equal(t, uint(2), MustDecimal(-12345).Scale())
}

func TestDecimalSign_(t *testing.T) {
	d := MustDecimal(0)
	assert.Equal(t, 0, d.Sign())

	d = MustDecimal(0, 5)
	assert.Equal(t, 0, d.Sign())

	d = MustDecimal(1)
	assert.Equal(t, 1, d.Sign())

	d = MustDecimal(-1)
	assert.Equal(t, -1, d.Sign())
}

func TestAdjustDecimalScale_(t *testing.T) {
	d1, d2 := MustDecimal(15, 1), MustDecimal(125, 1)
	funcs.Must(AdjustDecimalScale(&d1, &d2))
	assert.Equal(t, MustDecimal(15, 1), d1)
	assert.Equal(t, MustDecimal(125, 1), d2)

	d1, d2 = MustDecimal(15, 1), MustDecimal(125, 2)
	funcs.Must(AdjustDecimalScale(&d1, &d2))
	assert.Equal(t, MustDecimal(150, 2), d1)
	assert.Equal(t, MustDecimal(125, 2), d2)

	d1, d2 = MustDecimal(15, 1), MustDecimal(100_000_000_000_000_000, 0)
	funcs.Must(AdjustDecimalScale(&d1, &d2))
	assert.Equal(t, MustDecimal(2, 0), d1)
	assert.Equal(t, MustDecimal(100_000_000_000_000_000, 0), d2)

	d1, d2 = MustDecimal(154, 2), MustDecimal(100_000_000_000_000_000, 0)
	funcs.Must(AdjustDecimalScale(&d1, &d2))
	assert.Equal(t, MustDecimal(2, 0), d1)
	assert.Equal(t, MustDecimal(100_000_000_000_000_000, 0), d2)

	d1, d2 = MustDecimal(149, 2), MustDecimal(100_000_000_000_000_000, 0)
	funcs.Must(AdjustDecimalScale(&d1, &d2))
	assert.Equal(t, MustDecimal(1, 0), d1)
	assert.Equal(t, MustDecimal(100_000_000_000_000_000, 0), d2)

	d1, d2 = MustDecimal(144, 2), MustDecimal(100_000_000_000_000_000, 0)
	funcs.Must(AdjustDecimalScale(&d1, &d2))
	assert.Equal(t, MustDecimal(1, 0), d1)
	assert.Equal(t, MustDecimal(100_000_000_000_000_000, 0), d2)

	d1, d2 = MustDecimal(-154, 2), MustDecimal(100_000_000_000_000_000, 0)
	funcs.Must(AdjustDecimalScale(&d1, &d2))
	assert.Equal(t, MustDecimal(-2, 0), d1)
	assert.Equal(t, MustDecimal(100_000_000_000_000_000, 0), d2)

	d1, d2 = MustDecimal(15, 1), MustDecimal(999_999_999_999_999_995, 0)
	assert.Equal(t, fmt.Errorf("The decimal value 999999999999999995 is too large to round up"), AdjustDecimalScale(&d1, &d2))
	assert.Equal(t, MustDecimal(15, 1), d1)
	assert.Equal(t, MustDecimal(999_999_999_999_999_995, 0), d2)

	d1, d2 = MustDecimal(15, 1), MustDecimal(-999_999_999_999_999_995, 0)
	assert.Equal(t, fmt.Errorf("The decimal value -999999999999999995 is too small to round down"), AdjustDecimalScale(&d1, &d2))
	assert.Equal(t, MustDecimal(15, 1), d1)
	assert.Equal(t, MustDecimal(-999_999_999_999_999_995, 0), d2)
}

func TestAdjustDecimalFormat_(t *testing.T) {
	d1, d2 := MustDecimal(1, 0), MustDecimal(2, 0)
	assert.Equal(t, tuple.Of2("/1", "/2"), tuple.Of2(AdjustDecimalFormat(d1, d2)))

	d1, d2 = MustDecimal(12, 0), MustDecimal(1, 0)
	assert.Equal(t, tuple.Of2("/12", "/01"), tuple.Of2(AdjustDecimalFormat(d1, d2)))

	d1, d2 = MustDecimal(-1, 0), MustDecimal(12, 0)
	assert.Equal(t, tuple.Of2("-01", "/12"), tuple.Of2(AdjustDecimalFormat(d1, d2)))

	d1, d2 = MustDecimal(11, 0), MustDecimal(-12, 0)
	assert.Equal(t, tuple.Of2("/11", "-12"), tuple.Of2(AdjustDecimalFormat(d1, d2)))

	d1, d2 = MustDecimal(-123, 1), MustDecimal(-234, 1)
	assert.Equal(t, tuple.Of2("-12.3", "-23.4"), tuple.Of2(AdjustDecimalFormat(d1, d2)))

	d1, d2 = MustDecimal(123, 2), MustDecimal(2345)
	assert.Equal(t, tuple.Of2("/01.23", "/23.45"), tuple.Of2(AdjustDecimalFormat(d1, d2)))
}

func TestDecimalCmp_(t *testing.T) {
	assert.Equal(t, -1, MustDecimal(1).Cmp(MustDecimal(2)))
	assert.Equal(t, 0, MustDecimal(2).Cmp(MustDecimal(2)))
	assert.Equal(t, 1, MustDecimal(2).Cmp(MustDecimal(1)))

	assert.Equal(t, -1, MustDecimal(123, 2).Cmp(MustDecimal(123, 1)))
	assert.Equal(t, 0, MustDecimal(123, 1).Cmp(MustDecimal(123, 1)))
	assert.Equal(t, 1, MustDecimal(123, 1).Cmp(MustDecimal(123, 2)))

	assert.Equal(t, -1, MustDecimal(-1).Cmp(MustDecimal(2)))
	assert.Equal(t, 0, MustDecimal(-1).Cmp(MustDecimal(-1)))
	assert.Equal(t, 1, MustDecimal(2).Cmp(MustDecimal(-1)))

	assert.Equal(t, -1, MustDecimal(-2).Cmp(MustDecimal(-1)))
	assert.Equal(t, 0, MustDecimal(-2).Cmp(MustDecimal(-2)))
	assert.Equal(t, 1, MustDecimal(-1).Cmp(MustDecimal(-2)))
}

func TestDecimalNegate_(t *testing.T) {
	assert.Equal(t, MustDecimal(-5), MustDecimal(5).Negate())
	assert.Equal(t, MustDecimal(0), MustDecimal(0).Negate())
	assert.Equal(t, MustDecimal(5), MustDecimal(-5).Negate())
}

func TestDecimalAdd_(t *testing.T) {
	//   0.010
	// + 0.001
	// = 0.011
	d1, d2 := MustDecimal(1, 2), MustDecimal(1, 3)
	assert.Equal(t, tuple.Of2(MustDecimal(11, 3), error(nil)), tuple.Of2(d1.Add(d2)))

	//   -0.010
	// + -0.001
	// = -0.011
	d1, d2 = MustDecimal(-1, 2), MustDecimal(-1, 3)
	assert.Equal(t, tuple.Of2(MustDecimal(-11, 3), error(nil)), tuple.Of2(d1.Add(d2)))

	//    0.010
	// + -0.001
	// =  0.009
	d1, d2 = MustDecimal(1, 2), MustDecimal(-1, 3)
	assert.Equal(t, MustDecimal(9, 3), d1.MustAdd(d2))

	//   -0.010
	// +  0.001
	// = -0.009
	d1, d2 = MustDecimal(-1, 2), MustDecimal(1, 3)
	assert.Equal(t, MustDecimal(-9, 3), d1.MustAdd(d2))

	// Scale error
	d1, d2 = MustDecimal(15, 1), MustDecimal(999_999_999_999_999_995, 0)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal value 999999999999999995 is too large to round up")), tuple.Of2(d1.Add(d2)))

	// Overflow
	d1, d2 = MustDecimal(999_999_999_999_999_995, 2), MustDecimal(5, 2)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 9999999999999999.95 + 0.05 overflowed")), tuple.Of2(d1.Add(d2)))

	// Overflow
	d1, d2 = MustDecimal(999_999_999_999_999_999, 2), MustDecimal(999_999_999_999_999_999, 2)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 9999999999999999.99 + 9999999999999999.99 overflowed")), tuple.Of2(d1.Add(d2)))

	// Underflow
	d1, d2 = MustDecimal(-999_999_999_999_999_995, 2), MustDecimal(-5, 2)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation -9999999999999999.95 + -0.05 underflowed")), tuple.Of2(d1.Add(d2)))

	// Underflow
	d1, d2 = MustDecimal(-999_999_999_999_999_999, 2), MustDecimal(-999_999_999_999_999_999, 2)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation -9999999999999999.99 + -9999999999999999.99 underflowed")), tuple.Of2(d1.Add(d2)))
}

func TestDecimalSub_(t *testing.T) {
	//   0.010
	// - 0.001
	// = 0.011
	d1, d2 := MustDecimal(1, 2), MustDecimal(1, 3)
	assert.Equal(t, tuple.Of2(MustDecimal(9, 3), error(nil)), tuple.Of2(d1.Sub(d2)))

	//   -0.010
	// - -0.001
	// = -0.009
	d1, d2 = MustDecimal(-1, 2), MustDecimal(-1, 3)
	assert.Equal(t, tuple.Of2(MustDecimal(-9, 3), error(nil)), tuple.Of2(d1.Sub(d2)))

	//    0.010
	// - -0.001
	// =  0.011
	d1, d2 = MustDecimal(1, 2), MustDecimal(-1, 3)
	assert.Equal(t, MustDecimal(11, 3), d1.MustSub(d2))

	//   -0.010
	// -  0.001
	// = -0.011
	d1, d2 = MustDecimal(-1, 2), MustDecimal(1, 3)
	assert.Equal(t, MustDecimal(-11, 3), d1.MustSub(d2))

	// Scale error
	d1, d2 = MustDecimal(15, 1), MustDecimal(999_999_999_999_999_995, 0)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal value -999999999999999995 is too small to round down")), tuple.Of2(d1.Sub(d2)))

	// Overflow
	d1, d2 = MustDecimal(999_999_999_999_999_995, 2), MustDecimal(-5, 2)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 9999999999999999.95 - -0.05 overflowed")), tuple.Of2(d1.Sub(d2)))

	// Overflow
	d1, d2 = MustDecimal(999_999_999_999_999_999, 2), MustDecimal(-999_999_999_999_999_999, 2)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 9999999999999999.99 - -9999999999999999.99 overflowed")), tuple.Of2(d1.Sub(d2)))

	// Underflow
	d1, d2 = MustDecimal(-999_999_999_999_999_995, 2), MustDecimal(5, 2)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation -9999999999999999.95 - 0.05 underflowed")), tuple.Of2(d1.Sub(d2)))

	// Underflow
	d1, d2 = MustDecimal(-999_999_999_999_999_999, 2), MustDecimal(999_999_999_999_999_999, 2)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation -9999999999999999.99 - 9999999999999999.99 underflowed")), tuple.Of2(d1.Sub(d2)))
}

func TestDecimalMul_(t *testing.T) {
	// 1.5 * 2.5 = 3.75
	d1, d2 := MustDecimal(15, 1), MustDecimal(25, 1)
	assert.Equal(t, tuple.Of2(MustDecimal(375, 2), error(nil)), tuple.Of2(d1.Mul(d2)))

	// -1.5 * 2.5 = -3.75
	d1, d2 = MustDecimal(-15, 1), MustDecimal(25, 1)
	assert.Equal(t, tuple.Of2(MustDecimal(-375, 2), error(nil)), tuple.Of2(d1.Mul(d2)))

	// 1.5 * -2.5 = -3.75
	d1, d2 = MustDecimal(15, 1), MustDecimal(-25, 1)
	assert.Equal(t, MustDecimal(-375, 2), d1.MustMul(d2))

	// -1.5 * -2.5 = 3.75
	d1, d2 = MustDecimal(-15, 1), MustDecimal(-25, 1)
	assert.Equal(t, MustDecimal(375, 2), d1.MustMul(d2))

	// Overflow
	// - Within bounds of binary, but beyond bounds of 18 decimals
	d1, d2 = MustDecimal(999_999_999_999_999_999, 0), MustDecimal(2, 1)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 999999999999999999 * 0.2 overflowed")), tuple.Of2(d1.Mul(d2)))

	// - Beyond bounds of binary, but only a little
	d1, d2 = MustDecimal(999_999_999_999_999_999, 0), MustDecimal(16, 0)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 999999999999999999 * 16 overflowed")), tuple.Of2(d1.Mul(d2)))

	// - Way beyond bounds of binary
	d1, d2 = MustDecimal(999_999_999_999_999_999, 0), MustDecimal(999_999_999_999_999_999, 1)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 999999999999999999 * 99999999999999999.9 overflowed")), tuple.Of2(d1.Mul(d2)))
}

func TestDecimalDivIntQuoRem_(t *testing.T) {
	// 100.00 / 3 = 33.33 r 00.01
	de, dv := MustDecimal(100_00), uint(3)
	assert.Equal(t, tuple.Of3(MustDecimal(3333), MustDecimal(1), error(nil)), tuple.Of3(de.DivIntQuoRem(dv)))

	// 100.00 / 100 = 1.00 r 0.00
	de, dv = MustDecimal(100_00), uint(100)
	assert.Equal(t, tuple.Of2(MustDecimal(100), MustDecimal(0)), tuple.Of2(de.MustDivIntQuoRem(dv)))

	// Division by zero
	assert.Equal(t, tuple.Of3(Decimal{}, Decimal{}, fmt.Errorf("The decimal calculation 100.00 / 0 is not allowed")), tuple.Of3(de.DivIntQuoRem(0)))

	// Division by divisor > dividend
	assert.Equal(t, tuple.Of3(Decimal{}, Decimal{}, fmt.Errorf("The decimal calculation 100.00 / 101 is not allowed, the divisor is larger than the dividend")), tuple.Of3(de.DivIntQuoRem(101)))
}

func TestDecimalDivIntAdd_(t *testing.T) {
	// 100.00 / 3 = [33.34, 33.34, 33.33]
	de, dv := MustDecimal(100_00), uint(3)
	assert.Equal(t, tuple.Of2([]Decimal{MustDecimal(3334), MustDecimal(3333), MustDecimal(3333)}, error(nil)), tuple.Of2(de.DivIntAdd(dv)))

	// 100.00 / 5 = [20.00, 20.00, 20.00, 20.00, 20.00]
	de, dv = MustDecimal(100_00), uint(5)
	assert.Equal(t, []Decimal{MustDecimal(2000), MustDecimal(2000), MustDecimal(2000), MustDecimal(2000), MustDecimal(2000)}, de.MustDivIntAdd(dv))

	// 100.00 / 100 = [1.00 repeated 10 times]
	de, dv = MustDecimal(100_00), uint(100)
	add := make([]Decimal, 100)
	for i := 0; i < 100; i++ {
		add[i] = MustDecimal(100)
	}
	assert.Equal(t, add, de.MustDivIntAdd(dv))

	// // Division by zero
	assert.Equal(t, tuple.Of2([]Decimal(nil), fmt.Errorf("The decimal calculation 100.00 / 0 is not allowed")), tuple.Of2(de.DivIntAdd(0)))

	// Division by divisor > dividend
	assert.Equal(t, tuple.Of2([]Decimal(nil), fmt.Errorf("The decimal calculation 100.00 / 101 is not allowed, the divisor is larger than the dividend")), tuple.Of2(de.DivIntAdd(101)))
}
