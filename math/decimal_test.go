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
	assert.Equal(t, tuple.Of2(Decimal{scale: 0, value: 1, denormalized: false}, error(nil)), tuple.Of2(OfDecimal(1_00, 2)))
	assert.Equal(t, tuple.Of2(Decimal{scale: 3, value: -1_001, denormalized: false}, error(nil)), tuple.Of2(OfDecimal(-1_001, 3)))

	assert.Equal(t, Decimal{scale: 0, value: 1, denormalized: false}, MustDecimal(1_00, 2))
	assert.Equal(t, Decimal{scale: 3, value: -1001, denormalized: false}, MustDecimal(-1_001, 3))

	assert.Equal(
		t,
		tuple.Of2(Decimal{}, fmt.Errorf("The Decimal scale 19 is too large: the value must be <= 18")),
		tuple.Of2(OfDecimal(0, 19)),
	)
	assert.Equal(
		t,
		tuple.Of2(Decimal{}, fmt.Errorf("The Decimal value 1234567890123456789 is too large: the value must be <= 999_999_999_999_999_999")),
		tuple.Of2(OfDecimal(12345678901234567_89, 2)),
	)
	assert.Equal(
		t,
		tuple.Of2(Decimal{}, fmt.Errorf("The Decimal value -1234567890123456789 is too small: the value must be >= -999_999_999_999_999_999")),
		tuple.Of2(OfDecimal(-12345678901234567_89, 2)),
	)
}

func TestStringToDecimal_(t *testing.T) {
	assert.Equal(t, Decimal{scale: 0, value: 0, denormalized: false}, MustStringToDecimal("0"))
	assert.Equal(t, tuple.Of2(Decimal{scale: 0, value: 100, denormalized: false}, error(nil)), tuple.Of2(StringToDecimal("100")))
	assert.Equal(t, tuple.Of2(Decimal{scale: 3, value: -1_001, denormalized: false}, error(nil)), tuple.Of2(StringToDecimal("-1.001")))

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
	assert.Equal(t, "12.3", MustDecimal(12_3, 1).String())
	assert.Equal(t, "1.23", MustDecimal(1_23, 2).String())
	assert.Equal(t, "0.123", MustDecimal(123, 3).String())
	assert.Equal(t, "0.0123", MustDecimal(123, 4).String())
	assert.Equal(t, "0.00123", MustDecimal(123, 5).String())
	assert.Equal(t, "-0.00123", MustDecimal(-123, 5).String())
}

func TestDecimalPrecision_(t *testing.T) {
	assert.Equal(t, 3, MustDecimal(123, 0).Precision())
	assert.Equal(t, 3, MustDecimal(-123, 0).Precision())
	assert.Equal(t, 3, MustDecimal(123_00, 2).Precision())
	assert.Equal(t, 5, MustDecimal(-123_45, 2).Precision())
}

func TestDecimalScale_(t *testing.T) {
	assert.Equal(t, uint(0), MustDecimal(123, 0).Scale())
	assert.Equal(t, uint(0), MustDecimal(-123, 0).Scale())
	assert.Equal(t, uint(0), MustDecimal(123_00, 2).Scale())
	assert.Equal(t, uint(2), MustDecimal(-123_45, 2).Scale())
}

func TestDecimalNormalized_(t *testing.T) {
	assert.True(t, union.OfResultError(OfDecimal(123, 0)).Get().Normalized())
	assert.False(t, union.OfResultError(OfDecimal(123, 0, false)).Get().Normalized())

	assert.True(t, MustDecimal(123, 0).Normalized())
	assert.False(t, MustDecimal(123, 0, false).Normalized())
}

func TestDecimalSign_(t *testing.T) {
	d := MustDecimal(0, 2)
	assert.Equal(t, 0, d.Sign())

	d = MustDecimal(0, 5)
	assert.Equal(t, 0, d.Sign())

	d = MustDecimal(1, 2)
	assert.Equal(t, 1, d.Sign())

	d = MustDecimal(-1, 2)
	assert.Equal(t, -1, d.Sign())
}

func TestAdjustDecimalScale_(t *testing.T) {
	d1, d2 := MustDecimal(1_5, 1), MustDecimal(12_5, 1)
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, MustDecimal(1_5, 1), d1)
	assert.Equal(t, MustDecimal(12_5, 1), d2)

	d1, d2 = MustDecimal(1_5, 1), MustDecimal(1_25, 2)
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, Decimal{value: 1_50, scale: 2}, d1)
	assert.Equal(t, MustDecimal(1_25, 2), d2)

	d1, d2 = MustDecimal(1_5, 1), MustDecimal(100_000_000_000_000_000, 0)
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, MustDecimal(2, 0), d1)
	assert.Equal(t, MustDecimal(100_000_000_000_000_000, 0), d2)

	d1, d2 = MustDecimal(1_54, 2), MustDecimal(100_000_000_000_000_000, 0)
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, MustDecimal(2, 0), d1)
	assert.Equal(t, MustDecimal(100_000_000_000_000_000, 0), d2)

	d1, d2 = MustDecimal(1_49, 2), MustDecimal(100_000_000_000_000_000, 0)
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, MustDecimal(1, 0), d1)
	assert.Equal(t, MustDecimal(100_000_000_000_000_000, 0), d2)

	d1, d2 = MustDecimal(1_44, 2), MustDecimal(100_000_000_000_000_000, 0)
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, MustDecimal(1, 0), d1)
	assert.Equal(t, MustDecimal(100_000_000_000_000_000, 0), d2)

	d1, d2 = MustDecimal(-1_54, 2), MustDecimal(100_000_000_000_000_000, 0)
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, MustDecimal(-2, 0), d1)
	assert.Equal(t, MustDecimal(100_000_000_000_000_000, 0), d2)

	// 1.5 and 199_999_999_999_999_995
	// the second cannot be multiplied by 10, so the 1.5 is rounded to 2
	d1, d2 = MustDecimal(1_5, 1), MustDecimal(199_999_999_999_999_995, 0)
	MustAdjustDecimalScale(&d1, &d2)
	assert.Equal(t, MustDecimal(2, 0), d1)
	assert.Equal(t, MustDecimal(199_999_999_999_999_995, 0), d2)

	d1, d2 = MustDecimal(1_5, 1), MustDecimal(999_999_999_999_999_995, 0)
	assert.Equal(t, fmt.Errorf("The decimal value 999999999999999995 is too large to round up"), AdjustDecimalScale(&d1, &d2))
	assert.Equal(t, MustDecimal(1_5, 1), d1)
	assert.Equal(t, MustDecimal(999_999_999_999_999_995, 0), d2)

	d1, d2 = MustDecimal(1_5, 1), MustDecimal(-999_999_999_999_999_995, 0)
	assert.Equal(t, fmt.Errorf("The decimal value -999999999999999995 is too small to round down"), AdjustDecimalScale(&d1, &d2))
	assert.Equal(t, MustDecimal(1_5, 1), d1)
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

	d1, d2 = MustDecimal(-12_3, 1), MustDecimal(-23_4, 1)
	assert.Equal(t, tuple.Of2("-12.3", "-23.4"), tuple.Of2(AdjustDecimalFormat(d1, d2)))

	d1, d2 = MustDecimal(1_23, 2), MustDecimal(23_45, 2)
	assert.Equal(t, tuple.Of2(" 01.23", " 23.45"), tuple.Of2(AdjustDecimalFormat(d1, d2)))
}

func TestDecimalCmp_(t *testing.T) {
	assert.Equal(t, -1, MustDecimal(1, 2).Cmp(MustDecimal(2, 2)))
	assert.Equal(t, 0, MustDecimal(2, 2).Cmp(MustDecimal(2, 2)))
	assert.Equal(t, 1, MustDecimal(2, 2).Cmp(MustDecimal(1, 2)))

	assert.Equal(t, -1, MustDecimal(1_23, 2).Cmp(MustDecimal(12_3, 1)))
	assert.Equal(t, 0, MustDecimal(12_3, 1).Cmp(MustDecimal(12_3, 1)))
	assert.Equal(t, 1, MustDecimal(12_3, 1).Cmp(MustDecimal(1_23, 2)))

	assert.Equal(t, -1, MustDecimal(-1, 2).Cmp(MustDecimal(2, 2)))
	assert.Equal(t, 0, MustDecimal(-1, 2).Cmp(MustDecimal(-1, 2)))
	assert.Equal(t, 1, MustDecimal(2, 2).Cmp(MustDecimal(-1, 2)))

	assert.Equal(t, -1, MustDecimal(-2, 2).Cmp(MustDecimal(-1, 2)))
	assert.Equal(t, 0, MustDecimal(-2, 2).Cmp(MustDecimal(-2, 2)))
	assert.Equal(t, 1, MustDecimal(-1, 2).Cmp(MustDecimal(-2, 2)))
}

func TestDecimalNegate_(t *testing.T) {
	assert.Equal(t, MustDecimal(-5, 2), MustDecimal(5, 2).Negate())
	assert.Equal(t, MustDecimal(0, 2), MustDecimal(0, 2).Negate())
	assert.Equal(t, MustDecimal(5, 2), MustDecimal(-5, 2).Negate())
}

func TestDecimalMagnitudeLessThanOne_(t *testing.T) {
	// scale = 0
	assert.True(t, Decimal{value: 0, scale: 0, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1, scale: 0, denormalized: false}.MagnitudeLessThanOne())

	// scale = 1
	assert.True(t, Decimal{value: -9, scale: 1, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_0, scale: 1, denormalized: false}.MagnitudeLessThanOne())

	// scale = 2
	assert.True(t, Decimal{value: 99, scale: 2, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: -1_00, scale: 2, denormalized: false}.MagnitudeLessThanOne())

	// scale = 3
	assert.True(t, Decimal{value: 999, scale: 3, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_000, scale: 3, denormalized: false}.MagnitudeLessThanOne())

	// scale = 4
	assert.True(t, Decimal{value: 9_999, scale: 4, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_0_000, scale: 4, denormalized: false}.MagnitudeLessThanOne())

	// scale = 5
	assert.True(t, Decimal{value: 99_999, scale: 5, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_00_000, scale: 5, denormalized: false}.MagnitudeLessThanOne())

	// scale = 6
	assert.True(t, Decimal{value: 999_999, scale: 6, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_000_000, scale: 6, denormalized: false}.MagnitudeLessThanOne())

	// scale = 7
	assert.True(t, Decimal{value: 9_999_999, scale: 7, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_0_000_000, scale: 7, denormalized: false}.MagnitudeLessThanOne())

	// scale = 8
	assert.True(t, Decimal{value: 99_999_999, scale: 8, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_00_000_000, scale: 8, denormalized: false}.MagnitudeLessThanOne())

	// scale = 9
	assert.True(t, Decimal{value: 999_999_999, scale: 9, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_000_000_000, scale: 9, denormalized: false}.MagnitudeLessThanOne())

	// scale = 10
	assert.True(t, Decimal{value: 9_999_999_999, scale: 10, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_0_000_000_000, scale: 10, denormalized: false}.MagnitudeLessThanOne())

	// scale = 11
	assert.True(t, Decimal{value: 99_999_999_999, scale: 11, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_00_000_000_000, scale: 11, denormalized: false}.MagnitudeLessThanOne())

	// scale = 12
	assert.True(t, Decimal{value: 999_999_999_999, scale: 12, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_000_000_000_000, scale: 12, denormalized: false}.MagnitudeLessThanOne())

	// scale = 13
	assert.True(t, Decimal{value: 9_999_999_999_999, scale: 13, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_0_000_000_000_000, scale: 13, denormalized: false}.MagnitudeLessThanOne())

	// scale = 14
	assert.True(t, Decimal{value: 99_999_999_999_999, scale: 14, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_00_000_000_000_000, scale: 14, denormalized: false}.MagnitudeLessThanOne())

	// scale = 15
	assert.True(t, Decimal{value: 999_999_999_999_999, scale: 15, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_000_000_000_000_000, scale: 15, denormalized: false}.MagnitudeLessThanOne())

	// scale = 16
	assert.True(t, Decimal{value: 9_999_999_999_999_999, scale: 16, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_0_000_000_000_000_000, scale: 16, denormalized: false}.MagnitudeLessThanOne())

	// scale = 17
	assert.True(t, Decimal{value: 99_999_999_999_999_999, scale: 17, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_00_000_000_000_000_000, scale: 17, denormalized: false}.MagnitudeLessThanOne())

	// scale = 18 (this isn't quite realistic, the constructors would not allow a 19-digit value)
	assert.True(t, Decimal{value: 999_999_999_999_999_999, scale: 18, denormalized: false}.MagnitudeLessThanOne())
	assert.False(t, Decimal{value: 1_000_000_000_000_000_000, scale: 18, denormalized: false}.MagnitudeLessThanOne())
}

func TestDecimalNormalize_(t *testing.T) {
	// No trailing zeros
	d := MustDecimal(1, 0)
	assert.Equal(t, MustDecimal(1, 0), d)

	d = MustDecimal(-1, 0)
	assert.Equal(t, MustDecimal(-1, 0), d)

	d = MustDecimal(12, 1)
	assert.Equal(t, MustDecimal(12, 1), d)

	d = MustDecimal(123, 2)
	assert.Equal(t, MustDecimal(123, 2), d)

	// One trailing zero
	d = MustDecimal(10, 1)
	assert.Equal(t, MustDecimal(1, 0), d)

	d = MustDecimal(-10, 1)
	assert.Equal(t, MustDecimal(-1, 0), d)

	d = MustDecimal(120, 2)
	assert.Equal(t, MustDecimal(12, 1), d)

	d = MustDecimal(-1230, 3)
	assert.Equal(t, MustDecimal(-123, 2), d)

	// Two trailing zeros
	d = MustDecimal(100, 2)
	assert.Equal(t, MustDecimal(1, 0), d)

	d = MustDecimal(-100, 2)
	assert.Equal(t, MustDecimal(-1, 0), d)

	d = MustDecimal(1200, 3)
	assert.Equal(t, MustDecimal(12, 1), d)

	d = MustDecimal(12300, 4)
	assert.Equal(t, MustDecimal(123, 2), d)
}

func TestDecimalAdd_(t *testing.T) {
	//   0.01
	// + 0.001
	// = 0.011
	d1, d2 := MustDecimal(1, 2), MustDecimal(1, 3)
	assert.Equal(t, tuple.Of2(MustDecimal(11, 3), error(nil)), tuple.Of2(d1.Add(d2)))

	//   0.01
	// + 0.0010
	// = 0.011
	d1, d2 = MustDecimal(1, 2), MustDecimal(10, 4)
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

    // Normalization
	//   0.010
	// + 0.020
	// = 0.03
	d1, d2 = MustDecimal(-1, 2), MustDecimal(1, 3)
	assert.Equal(t, MustDecimal(-9, 3), d1.MustAdd(d2))

	// Scale error
	d1, d2 = MustDecimal(1_5, 1), MustDecimal(999_999_999_999_999_995, 0)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal value 999999999999999995 is too large to round up")), tuple.Of2(d1.Add(d2)))

	// Overflow
	d1, d2 = MustDecimal(9_999_999_999_999_999_95, 2), MustDecimal(5, 2)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 9999999999999999.95 + 0.05 overflowed")), tuple.Of2(d1.Add(d2)))

	// Overflow
	d1, d2 = MustDecimal(9_999_999_999_999_999_99, 2), MustDecimal(9_999_999_999_999_999_99, 2)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 9999999999999999.99 + 9999999999999999.99 overflowed")), tuple.Of2(d1.Add(d2)))

	// Underflow
	d1, d2 = MustDecimal(-9_999_999_999_999_999_95, 2), MustDecimal(-5, 2)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation -9999999999999999.95 + -0.05 underflowed")), tuple.Of2(d1.Add(d2)))

	// Underflow
	d1, d2 = MustDecimal(-9_999_999_999_999_999_99, 2), MustDecimal(-9_999_999_999_999_999_99, 2)
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
	d1, d2 = MustDecimal(1_5, 1), MustDecimal(999_999_999_999_999_995, 0)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal value -999999999999999995 is too small to round down")), tuple.Of2(d1.Sub(d2)))

	// Overflow
	d1, d2 = MustDecimal(9_999_999_999_999_999_95, 2), MustDecimal(-5, 2)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 9999999999999999.95 - -0.05 overflowed")), tuple.Of2(d1.Sub(d2)))

	// Overflow
	d1, d2 = MustDecimal(9_999_999_999_999_999_99, 2), MustDecimal(-9_999_999_999_999_999_99, 2)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 9999999999999999.99 - -9999999999999999.99 overflowed")), tuple.Of2(d1.Sub(d2)))

	// Underflow
	d1, d2 = MustDecimal(-9_999_999_999_999_999_95, 2), MustDecimal(5, 2)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation -9999999999999999.95 - 0.05 underflowed")), tuple.Of2(d1.Sub(d2)))

	// Underflow
	d1, d2 = MustDecimal(-9_999_999_999_999_999_99, 2), MustDecimal(9_999_999_999_999_999_99, 2)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation -9999999999999999.99 - 9999999999999999.99 underflowed")), tuple.Of2(d1.Sub(d2)))
}

func TestDecimalMul_(t *testing.T) {
	// 1.5 * 2.5 = 3.75
	d1, d2 := MustDecimal(1_5, 1), MustDecimal(2_5, 1)
	assert.Equal(t, tuple.Of2(MustDecimal(3_75, 2), error(nil)), tuple.Of2(d1.Mul(d2)))

	// -1.5 * 2.5 = -3.75
	d1, d2 = MustDecimal(-1_5, 1), MustDecimal(2_5, 1)
	assert.Equal(t, tuple.Of2(MustDecimal(-3_75, 2), error(nil)), tuple.Of2(d1.Mul(d2)))

	// 1.5 * -2.5 = -3.75
	d1, d2 = MustDecimal(1_5, 1), MustDecimal(-2_5, 1)
	assert.Equal(t, MustDecimal(-3_75, 2), d1.MustMul(d2))

	// -1.5 * -2.5 = 3.75
	d1, d2 = MustDecimal(-1_5, 1), MustDecimal(-2_5, 1)
	assert.Equal(t, MustDecimal(3_75, 2), d1.MustMul(d2))

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
	d1, d2 = MustDecimal(999_999_999_999_999_999, 0), MustDecimal(99_999_999_999_999_999_9, 1)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 999999999999999999 * 99999999999999999.9 overflowed")), tuple.Of2(d1.Mul(d2)))

	// Underflow
	// - Within bounds of signed 64 bit int, but beyond bounds of 18 decimals
	d1, d2 = MustDecimal(-999_999_999_999_999_999, 0), MustDecimal(2, 1)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation -999999999999999999 * 0.2 underflowed")), tuple.Of2(d1.Mul(d2)))

	// - Beyond bounds of signed 64 bit int, but only a little
	d1, d2 = MustDecimal(999_999_999_999_999_999, 0), MustDecimal(-16, 0)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation 999999999999999999 * -16 underflowed")), tuple.Of2(d1.Mul(d2)))

	// - Way beyond bounds of signed 64 bit int
	d1, d2 = MustDecimal(-999_999_999_999_999_999, 0), MustDecimal(99_999_999_999_999_999_9, 1)
	assert.Equal(t, tuple.Of2(Decimal{}, fmt.Errorf("The decimal calculation -999999999999999999 * 99999999999999999.9 underflowed")), tuple.Of2(d1.Mul(d2)))
}

func TestDecimalDivIntQuoRem_(t *testing.T) {
	// 100.00 / 3 = 33.33 r 00.01
	de, dv := MustDecimal(100_00, 2, false), uint(3)
	assert.Equal(t, tuple.Of3(MustDecimal(33_33, 2, false), MustDecimal(1, 2, false), error(nil)), tuple.Of3(de.DivIntQuoRem(dv)))

	// -100.00 / 3 = -33.33 r -00.01
	de, dv = MustDecimal(-100_00, 2, false), uint(3)
	assert.Equal(t, tuple.Of3(MustDecimal(-33_33, 2, false), MustDecimal(-1, 2, false), error(nil)), tuple.Of3(de.DivIntQuoRem(dv)))

	// 100.00 / 100 = 1.00 r 0.00
	de, dv = MustDecimal(100_00, 2, false), uint(100)
	assert.Equal(t, tuple.Of2(MustDecimal(1_00, 2, false), MustDecimal(0, 2, false)), tuple.Of2(de.MustDivIntQuoRem(dv)))

	// Division by zero
	assert.Equal(t, tuple.Of3(Decimal{}, Decimal{}, fmt.Errorf("The decimal calculation 100.00 / 0 is not allowed")), tuple.Of3(de.DivIntQuoRem(0)))

	// Division by divisor > dividend
	assert.Equal(t, tuple.Of3(Decimal{}, Decimal{}, fmt.Errorf("The decimal calculation 100.00 / 101 is not allowed, the divisor is larger than the dividend")), tuple.Of3(de.DivIntQuoRem(101)))
}

func TestDecimalDivIntAdd_(t *testing.T) {
	// 100.00 / 3 = [33.34, 33.34, 33.33]
	de, dv := MustDecimal(100_00, 2, false), uint(3)
	assert.Equal(t, tuple.Of2([]Decimal{MustDecimal(33_34, 2, false), MustDecimal(33_33, 2, false), MustDecimal(33_33, 2, false)}, error(nil)), tuple.Of2(de.DivIntAdd(dv)))

	// 100.00 / 5 = [20.00, 20.00, 20.00, 20.00, 20.00]
	de, dv = MustDecimal(100_00, 2, false), uint(5)
	assert.Equal(t, []Decimal{MustDecimal(20_00, 2, false), MustDecimal(20_00, 2, false), MustDecimal(20_00, 2, false), MustDecimal(20_00, 2, false), MustDecimal(20_00, 2, false)}, de.MustDivIntAdd(dv))

	// 100.00 / 100 = [1.00 repeated 100 times]
	de, dv = MustDecimal(100_00, 2, false), uint(100)
	add := make([]Decimal, 100)
	for i := 0; i < 100; i++ {
		add[i] = MustDecimal(1_00, 2, false)
	}
	assert.Equal(t, add, de.MustDivIntAdd(dv))

	// // Division by zero
	assert.Equal(t, tuple.Of2([]Decimal(nil), fmt.Errorf("The decimal calculation 100.00 / 0 is not allowed")), tuple.Of2(de.DivIntAdd(0)))

	// Division by divisor > dividend
	assert.Equal(t, tuple.Of2([]Decimal(nil), fmt.Errorf("The decimal calculation 100.00 / 101 is not allowed, the divisor is larger than the dividend")), tuple.Of2(de.DivIntAdd(101)))
}

func TestDecimalDiv_(t *testing.T) {
	// 1. 5000 / 200 = 25
	de, dv := MustDecimal(5_000, 0), MustDecimal(200, 0)
	assert.Equal(t, MustDecimal(25, 0), de.MustDiv(dv))

	// 2. 500.0 / 200 = 2.5
	de, dv = MustDecimal(500_0, 1), MustDecimal(200, 0)
	assert.Equal(t, MustDecimal(2_5, 1), de.MustDiv(dv))

	// 3. 500.0 / 2.00 = 250
	de, dv = MustDecimal(500_0, 1), MustDecimal(2_00, 2)
	assert.Equal(t, MustDecimal(250, 0), de.MustDiv(dv))

	// 4. 500.1 / 2.00 = 250.05
	de, dv = MustDecimal(500_1, 1), MustDecimal(2_00, 2)
	assert.Equal(t, MustDecimal(250_05, 2), de.MustDiv(dv))

	// 5. 5001 / 200 = 25.005
	de, dv = MustDecimal(5_001, 0), MustDecimal(200, 0)
	assert.Equal(t, MustDecimal(25_005, 3), de.MustDiv(dv))

	// 6. 5001 / -200 = -25.005
	de, dv = MustDecimal(5_001, 0), MustDecimal(-200, 0)
	assert.Equal(t, MustDecimal(-25_005, 3), de.MustDiv(dv))

	// 7. -5001 / 200 = -25.005
	de, dv = MustDecimal(-5_001, 0), MustDecimal(200, 0)
	assert.Equal(t, MustDecimal(-25_005, 3), de.MustDiv(dv))

	// 8. -5001 / -200 = 25.005
	de, dv = MustDecimal(-5_001, 0), MustDecimal(-200, 0)
	assert.Equal(t, MustDecimal(25_005, 3), de.MustDiv(dv))

	// 9. -500.1 / 200 = -2.5005
	de, dv = MustDecimal(-500_1, 1), MustDecimal(200, 0)
	assert.Equal(t, MustDecimal(-2_5005, 4), de.MustDiv(dv))

	// 	10. 5.123 / 0.021 = 243.952380952380952
	de, dv = MustDecimal(5_123, 3), MustDecimal(21, 3)
	assert.Equal(t, MustDecimal(243_952_380_952_380_952, 15), de.MustDiv(dv))

	// 11. 5 / 9 = 0.555555555555555556
	de, dv = MustDecimal(5, 0), MustDecimal(9, 0)
	assert.Equal(t, MustDecimal(555_555_555_555_555_556, 18), de.MustDiv(dv))

	// 12. 1.03075 / 0.25 = 4.123
	de, dv = MustDecimal(1_030_75, 5), MustDecimal(25, 2)
	assert.Equal(t, MustDecimal(4_123, 3), de.MustDiv(dv))

	// 13. 1_234_567_890_123_456.78 / 2.5 = 493_827_156_049_382.712
	de, dv = MustDecimal(1_234_567_890_123_456_78, 2), MustDecimal(2_5, 1)
	assert.Equal(t, MustDecimal(493_827_156_049_382_712, 3), de.MustDiv(dv))

	// 14. 1_234_567_890_123_456.78 / 0.25 = 4_938_271_560_493_827.12
	de, dv = MustDecimal(1_234_567_890_123_456_78, 2), MustDecimal(25, 2)
	assert.Equal(t, MustDecimal(4_938_271_560_493_827_12, 2), de.MustDiv(dv))

	// 15. 1_234_567_890_123_456.78 / 0.00025 = overflow
	de, dv = MustDecimal(1_234_567_890_123_456_78, 2), MustDecimal(25, 5)
	assert.Equal(t, union.OfError[Decimal](fmt.Errorf(errDecimalOverflowMsg, de, "/", dv)), union.OfResultError(de.Div(dv)))

	// 16. 1 / 100_000_000_000_000_000 = 0.000_000_000_000_000_01
	de, dv = MustDecimal(1, 0), MustDecimal(100_000_000_000_000_000, 0)
	assert.Equal(t, MustDecimal(1, 17), de.MustDiv(dv))

	// 17. 1 / 200_000_000_000_000_000 = overflow
	de, dv = MustDecimal(1, 0), MustDecimal(200_000_000_000_000_000, 0)
	assert.Equal(t, union.OfError[Decimal](fmt.Errorf(errDecimalOverflowMsg, de, "/", dv)), union.OfResultError(de.Div(dv)))

	// 18. 100_000_000_000_000_000 / 0.1 = overflow
	de, dv = MustDecimal(100_000_000_000_000_000, 0), MustDecimal(1, 1)
	assert.Equal(t, union.OfError[Decimal](fmt.Errorf(errDecimalOverflowMsg, de, "/", dv)), union.OfResultError(de.Div(dv)))
}
