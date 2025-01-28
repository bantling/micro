package math

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/bantling/micro/tuple"
	"github.com/bantling/micro/union"
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
	assert.Equal(t, Decimal{scale: 0, value: 0}, MustStringToDecimal("0"))
	assert.Equal(t, tuple.Of2(Decimal{scale: 0, value: 100}, error(nil)), tuple.Of2(StringToDecimal("100")))
	assert.Equal(t, tuple.Of2(Decimal{scale: 3, value: -1001}, error(nil)), tuple.Of2(StringToDecimal("-1.001")))

	assert.Equal(
		t,
		tuple.Of2(Decimal{}, fmt.Errorf("The string value x.5 is not a valid decimal string")),
		tuple.Of2(StringToDecimal("x.5")),
	)
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
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, MustDecimal(15, 1), d1)
	assert.Equal(t, MustDecimal(125, 1), d2)

	d1, d2 = MustDecimal(15, 1), MustDecimal(125, 2)
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, MustDecimal(150, 2), d1)
	assert.Equal(t, MustDecimal(125, 2), d2)

	d1, d2 = MustDecimal(15, 1), MustDecimal(100_000_000_000_000_000, 0)
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, MustDecimal(2, 0), d1)
	assert.Equal(t, MustDecimal(100_000_000_000_000_000, 0), d2)

	d1, d2 = MustDecimal(154, 2), MustDecimal(100_000_000_000_000_000, 0)
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, MustDecimal(2, 0), d1)
	assert.Equal(t, MustDecimal(100_000_000_000_000_000, 0), d2)

	d1, d2 = MustDecimal(149, 2), MustDecimal(100_000_000_000_000_000, 0)
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, MustDecimal(1, 0), d1)
	assert.Equal(t, MustDecimal(100_000_000_000_000_000, 0), d2)

	d1, d2 = MustDecimal(144, 2), MustDecimal(100_000_000_000_000_000, 0)
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, MustDecimal(1, 0), d1)
	assert.Equal(t, MustDecimal(100_000_000_000_000_000, 0), d2)

	d1, d2 = MustDecimal(-154, 2), MustDecimal(100_000_000_000_000_000, 0)
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, MustDecimal(-2, 0), d1)
	assert.Equal(t, MustDecimal(100_000_000_000_000_000, 0), d2)

	// 1.5 and 199_999_999_999_999_995
	// the second cannot be multiplied by 10, so the 1.5 is rounded to 2
	d1, d2 = MustDecimal(15, 1), MustDecimal(199_999_999_999_999_995, 0)
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, MustDecimal(2, 0), d1)
	assert.Equal(t, MustDecimal(199_999_999_999_999_995, 0), d2)

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
	assert.Equal(t, tuple.Of2(" 1", " 2"), tuple.Of2(AdjustDecimalFormat(d1, d2)))

	d1, d2 = MustDecimal(12, 0), MustDecimal(1, 0)
	assert.Equal(t, tuple.Of2(" 12", " 01"), tuple.Of2(AdjustDecimalFormat(d1, d2)))

	d1, d2 = MustDecimal(-1, 0), MustDecimal(12, 0)
	assert.Equal(t, tuple.Of2("-01", " 12"), tuple.Of2(AdjustDecimalFormat(d1, d2)))

	d1, d2 = MustDecimal(11, 0), MustDecimal(-12, 0)
	assert.Equal(t, tuple.Of2(" 11", "-12"), tuple.Of2(AdjustDecimalFormat(d1, d2)))

	d1, d2 = MustDecimal(-123, 1), MustDecimal(-234, 1)
	assert.Equal(t, tuple.Of2("-12.3", "-23.4"), tuple.Of2(AdjustDecimalFormat(d1, d2)))

	d1, d2 = MustDecimal(123, 2), MustDecimal(2345)
	assert.Equal(t, tuple.Of2(" 01.23", " 23.45"), tuple.Of2(AdjustDecimalFormat(d1, d2)))
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

func TestMagnitudeLessThanOne_(t *testing.T) {
	// scale = 0
	assert.True(t, Decimal{value: 0, scale: 0}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1, scale: 0}.MagnitudeLessThanOne())

	// scale = 1
	assert.True(t, Decimal{value: -9, scale: 1}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 10, scale: 1}.MagnitudeLessThanOne())

	// scale = 2
	assert.True(t, Decimal{value: 99, scale: 2}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: -100, scale: 2}.MagnitudeLessThanOne())

	// scale = 3
	assert.True(t, Decimal{value: 999, scale: 3}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_000, scale: 3}.MagnitudeLessThanOne())

	// scale = 4
	assert.True(t, Decimal{value: 9_999, scale: 4}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 10_000, scale: 4}.MagnitudeLessThanOne())

	// scale = 5
	assert.True(t, Decimal{value: 99_999, scale: 5}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 100_000, scale: 5}.MagnitudeLessThanOne())

	// scale = 6
	assert.True(t, Decimal{value: 999_999, scale: 6}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_000_000, scale: 6}.MagnitudeLessThanOne())

	// scale = 7
	assert.True(t, Decimal{value: 9_999_999, scale: 7}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 10_000_000, scale: 7}.MagnitudeLessThanOne())

	// scale = 8
	assert.True(t, Decimal{value: 99_999_999, scale: 8}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 100_000_000, scale: 8}.MagnitudeLessThanOne())

	// scale = 9
	assert.True(t, Decimal{value: 999_999_999, scale: 9}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_000_000_000, scale: 9}.MagnitudeLessThanOne())

	// scale = 10
	assert.True(t, Decimal{value: 9_999_999_999, scale: 10}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 10_000_000_000, scale: 10}.MagnitudeLessThanOne())

	// scale = 11
	assert.True(t, Decimal{value: 99_999_999_999, scale: 11}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 100_000_000_000, scale: 11}.MagnitudeLessThanOne())

	// scale = 12
	assert.True(t, Decimal{value: 999_999_999_999, scale: 12}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_000_000_000_000, scale: 12}.MagnitudeLessThanOne())

	// scale = 13
	assert.True(t, Decimal{value: 9_999_999_999_999, scale: 13}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 10_000_000_000_000, scale: 13}.MagnitudeLessThanOne())

	// scale = 14
	assert.True(t, Decimal{value: 99_999_999_999_999, scale: 14}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 100_000_000_000_000, scale: 14}.MagnitudeLessThanOne())

	// scale = 15
	assert.True(t, Decimal{value: 999_999_999_999_999, scale: 15}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_000_000_000_000_000, scale: 15}.MagnitudeLessThanOne())

	// scale = 16
	assert.True(t, Decimal{value: 9_999_999_999_999_999, scale: 16}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 10_000_000_000_000_000, scale: 16}.MagnitudeLessThanOne())

	// scale = 17
	assert.True(t, Decimal{value: 99_999_999_999_999_999, scale: 17}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 100_000_000_000_000_000, scale: 17}.MagnitudeLessThanOne())

	// scale = 18 (this isn't quite realistic, the constructors would not allow a 19-digit value)
	assert.True(t, Decimal{value: 999_999_999_999_999_999, scale: 18}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_000_000_000_000_000_000, scale: 18}.MagnitudeLessThanOne())
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

	// Over/underflow check does not result in division by zero error
	d1, d2 = MustDecimal(1, 0), MustDecimal(0, 0)
	assert.Equal(t, MustDecimal(0, 0), d1.MustMul(d2))

	// Overflow
	// - Within bounds of signed 64 bit int, but beyond bounds of 18 decimals
	d1, d2 = MustDecimal(999_999_999_999_999_999, 0), MustDecimal(2, 1)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 999999999999999999 * 0.2 overflowed")), tuple.Of2(d1.Mul(d2)))

	// - Beyond bounds of signed 64 bit int, but only a little
	d1, d2 = MustDecimal(999_999_999_999_999_999, 0), MustDecimal(16, 0)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 999999999999999999 * 16 overflowed")), tuple.Of2(d1.Mul(d2)))

	// - Way beyond bounds of signed 64 bit int
	d1, d2 = MustDecimal(999_999_999_999_999_999, 0), MustDecimal(999_999_999_999_999_999, 1)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 999999999999999999 * 99999999999999999.9 overflowed")), tuple.Of2(d1.Mul(d2)))

	// Underflow
	// - Within bounds of signed 64 bit int, but beyond bounds of 18 decimals
	d1, d2 = MustDecimal(-999_999_999_999_999_999, 0), MustDecimal(2, 1)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation -999999999999999999 * 0.2 underflowed")), tuple.Of2(d1.Mul(d2)))

	// - Beyond bounds of signed 64 bit int, but only a little
	d1, d2 = MustDecimal(999_999_999_999_999_999, 0), MustDecimal(-16, 0)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 999999999999999999 * -16 underflowed")), tuple.Of2(d1.Mul(d2)))

	// - Way beyond bounds of signed 64 bit int
	d1, d2 = MustDecimal(-999_999_999_999_999_999, 0), MustDecimal(999_999_999_999_999_999, 1)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation -999999999999999999 * 99999999999999999.9 underflowed")), tuple.Of2(d1.Mul(d2)))
}

func TestDecimalDivIntQuoRem_(t *testing.T) {
	// 100.00 / 3 = 33.33 r 00.01
	de, dv := MustDecimal(100_00), uint(3)
	assert.Equal(t, tuple.Of3(MustDecimal(3333), MustDecimal(1), error(nil)), tuple.Of3(de.DivIntQuoRem(dv)))

	// -100.00 / 3 = -33.33 r -00.01
	de, dv = MustDecimal(-100_00), uint(3)
	assert.Equal(t, tuple.Of3(MustDecimal(-3333), MustDecimal(-1), error(nil)), tuple.Of3(de.DivIntQuoRem(dv)))

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

func TestDecimalDiv_(t *testing.T) {
	// 	// 5000 / 200 = 25
	// 	de, dv := MustDecimal(5000, 0), MustDecimal(200, 0)
	// 	assert.Equal(t, union.OfResult(MustDecimal(25, 0)), union.OfResultError(de.Div(dv)))
	//
	// 	// 500.0 / 200 = 2.5
	// 	de, dv = MustDecimal(5000, 1), MustDecimal(200, 0)
	// 	assert.Equal(t, union.OfResult(MustDecimal(2_5, 1)), union.OfResultError(de.Div(dv)))
	//
	// 	// 500.0 / 2.00 = 250
	// 	de, dv = MustDecimal(5000, 1), MustDecimal(200, 2)
	// 	assert.Equal(t, union.OfResult(MustDecimal(250, 0)), union.OfResultError(de.Div(dv)))
	//
	// 	// 500.1 / 2.00 = 250.05
	// 	de, dv = MustDecimal(5001, 1), MustDecimal(200, 2)
	// 	assert.Equal(t, union.OfResult(MustDecimal(250_05, 2)), union.OfResultError(de.Div(dv)))
	//
	// 	// 5001 / 200 = 25.005
	// 	de, dv = MustDecimal(5001, 0), MustDecimal(200, 0)
	// 	assert.Equal(t, union.OfResult(MustDecimal(25_005, 3)), union.OfResultError(de.Div(dv)))
	//
	// 	// 5001 / -200 = -25.005
	// 	de, dv = MustDecimal(5001, 0), MustDecimal(-200, 0)
	// 	assert.Equal(t, union.OfResult(MustDecimal(-25_005, 3)), union.OfResultError(de.Div(dv)))
	//
	// 	// -5001 / 200 = -25.005
	// 	de, dv = MustDecimal(-5001, 0), MustDecimal(200, 0)
	// 	assert.Equal(t, union.OfResult(MustDecimal(-25_005, 3)), union.OfResultError(de.Div(dv)))
	//
	// 	// -5001 / -200 = 25.005
	// 	de, dv = MustDecimal(-5001, 0), MustDecimal(-200, 0)
	// 	assert.Equal(t, union.OfResult(MustDecimal(25_005, 3)), union.OfResultError(de.Div(dv)))
	//
	// 	// -500.1 / 200 = -2.5005
	// 	de, dv = MustDecimal(-5001, 1), MustDecimal(200, 0)
	// 	assert.Equal(t, union.OfResult(MustDecimal(-2_5005, 4)), union.OfResultError(de.Div(dv)))

	// 5.123 / 0.021 = 243.952380952380952
	// 2439523809523809520
	de, dv := MustDecimal(5123, 3), MustDecimal(21, 3)
	assert.Equal(t, union.OfResult(MustDecimal(243_952380952380952, 15)), union.OfResultError(de.Div(dv)))

	// // 1.03075 / 0.25 = 4.123
	// de, dv = MustDecimal(1_03075, 5), MustDecimal(25, 2)
	// assert.Equal(t, union.OfResult(MustDecimal(4_123, 3)), union.OfResultError(de.Div(dv)))
	//
	// // 1_234_567_890_123_456.78 / 2.5 = 493_827_156_049_382.712
	// de, dv = MustDecimal(1_234_567_890_123_456_78, 2), MustDecimal(2_5, 1)
	// assert.Equal(t, union.OfResult(MustDecimal(493_827_156_049_382_712, 3)), union.OfResultError(de.Div(dv)))
	//
	// // 1_234_567_890_123_456.78 / 0.25 = 4_938_271_560_493_827.12
	// de, dv = MustDecimal(1_234_567_890_123_456_78, 2), MustDecimal(25, 2)
	// assert.Equal(t, union.OfResult(MustDecimal(4_938_271_560_493_827_12, 2)), union.OfResultError(de.Div(dv)))
	//
	// // 1_234_567_890_123_456.78 / 0.00025 = overflow
	// de, dv = MustDecimal(1_234_567_890_123_456_78, 2), MustDecimal(25, 5)
	// assert.Equal(t, union.OfError[Decimal](fmt.Errorf(errDecimalOverflowMsg, de, "/", dv)), union.OfResultError(de.Div(dv)))
	//
	// // 1 / 100_000_000_000_000_000 = 0.000_000_000_000_000_01
	// de, dv = MustDecimal(1, 0), MustDecimal(100_000_000_000_000_000, 0)
	// assert.Equal(t, union.OfResult(MustDecimal(1, 17)), union.OfResultError(de.Div(dv)))
	//
	// // 1 / 200_000_000_000_000_000 = overflow
	// de, dv = MustDecimal(1, 0), MustDecimal(200_000_000_000_000_000, 0)
	// assert.Equal(t, union.OfError[Decimal](fmt.Errorf(errDecimalOverflowMsg, de, "/", dv)), union.OfResultError(de.Div(dv)))
}
