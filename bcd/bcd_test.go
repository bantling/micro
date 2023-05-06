package bcd

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/union"
	"github.com/stretchr/testify/assert"
)

func TestSignValues(t *testing.T) {
	var s Sign
	assert.Equal(t, Zero, s)

	assert.Equal(t, -1, int(Negative))
	assert.Equal(t, 0, int(Zero))
	assert.Equal(t, 1, int(Positive))
}

func TestOfSign_(t *testing.T) {
	assert.Equal(t, Negative, OfSign("-"))
	assert.Equal(t, Positive, OfSign(""))
	assert.Equal(t, Positive, OfSign("+"))
}

func TestSignString_(t *testing.T) {
	assert.Equal(t, "-", Negative.String())
	assert.Equal(t, "", Zero.String())
	assert.Equal(t, "", Positive.String())
}

func TestSignNegate_(t *testing.T) {
	assert.Equal(t, Positive, Negative.Negate())
	assert.Equal(t, Zero, Zero.Negate())
	assert.Equal(t, Negative, Positive.Negate())
}

func TestOfHexInternal_(t *testing.T) {
	// -1
	assert.Equal(t, Number{Negative, 0x1, 0}, ofHexInternal(Negative, 0x1, 0))

	// 0
	assert.Equal(t, Number{Zero, 0x0, 0}, ofHexInternal(Negative, 0x0, 0))

	// 1
	assert.Equal(t, Number{Positive, 0x1, 0}, ofHexInternal(Positive, 0x1, 0))

	// 0.456
	assert.Equal(t, Number{Positive, 0x456, 3}, ofHexInternal(Positive, 0x456, 3))

	// 123.456
	assert.Equal(t, Number{Positive, 0x123456, 3}, ofHexInternal(Positive, 0x123_456, 3))
}

func TestOfHex_(t *testing.T) {
	// -1
	assert.Equal(t, union.OfResult(Number{Negative, 0x1, 0}), union.OfResultError(OfHex(Negative, 0x1, 0)))

	// 0
	assert.Equal(t, union.OfResult(Number{Zero, 0x0, 0}), union.OfResultError(OfHex(Negative, 0x0, 0)))

	// 1
	assert.Equal(t, union.OfResult(Number{Positive, 0x1, 0}), union.OfResultError(OfHex(Positive, 0x1, 0)))

	// 0.456
	assert.Equal(t, union.OfResult(Number{Positive, 0x456, 3}), union.OfResultError(OfHex(Positive, 0x456, 3)))

	// 123.456
	assert.Equal(t, union.OfResult(Number{Positive, 0x123456, 3}), union.OfResultError(OfHex(Positive, 0x123_456, 3)))

	// Invalid number of decimals
	assert.Equal(t, union.OfError[Number](fmt.Errorf("Invalid number of decimals 17: the valid range is [0 .. 16]")), union.OfResultError(OfHex(Positive, 0x1, 17)))

	// Invalid digit
	assert.Equal(t, union.OfError[Number](fmt.Errorf(`Invalid Number "0x1A": the value must contain only decimal digits for each hex group`)), union.OfResultError(OfHex(Positive, 0x1A, 0)))
}

func TestOfString_(t *testing.T) {
	// -1
	assert.Equal(t, union.OfResult(Number{Negative, 0x1, 0}), union.OfResultError(OfString("-1")))

	// 0
	assert.Equal(t, union.OfResult(Number{Zero, 0x0, 0}), union.OfResultError(OfString("0")))

	// 1
	assert.Equal(t, union.OfResult(Number{Positive, 0x1, 0}), union.OfResultError(OfString("1")))

	// 123
	assert.Equal(t, union.OfResult(Number{Positive, 0x123, 0}), union.OfResultError(OfString("123")))

	// 0.456
	assert.Equal(t, union.OfResult(Number{Positive, 0x456, 3}), union.OfResultError(OfString("0.456")))

	// 123.456
	assert.Equal(t, union.OfResult(Number{Positive, 0x123456, 3}), union.OfResultError(OfString("123.456")))

	// Invalid strings
	for _, s := range []string{"", ".", ".1", "a", "++1", "--1"} {
		assert.Equal(t, union.OfError[Number](fmt.Errorf(numberStringErrMsg, s)), union.OfResultError(OfString(s)))
	}
}

func TestNumberString_(t *testing.T) {
	// 0
	assert.Equal(t, "0", funcs.MustValue(OfString("0")).String())

	// 1
	assert.Equal(t, "1", funcs.MustValue(OfString("1")).String())

	// -1
	assert.Equal(t, "-1", funcs.MustValue(OfString("-1")).String())

	// 0.1
	assert.Equal(t, "0.1", funcs.MustValue(OfString("0.1")).String())

	// 0.123
	assert.Equal(t, "0.123", funcs.MustValue(OfString("0.123")).String())

	// 1.23
	assert.Equal(t, "1.23", funcs.MustValue(OfString("1.23")).String())

	// 12.3
	assert.Equal(t, "12.3", funcs.MustValue(OfString("12.3")).String())

	// 123
	assert.Equal(t, "123", funcs.MustValue(OfString("123")).String())

	// 123.4
	assert.Equal(t, "123.4", funcs.MustValue(OfString("123.4")).String())

	// 123.45
	assert.Equal(t, "123.45", funcs.MustValue(OfString("123.45")).String())

	// 123.456
	assert.Equal(t, "123.456", funcs.MustValue(OfString("123.456")).String())
}

func TestAdjustToZero_(t *testing.T) {
	var n Number

	n = Number{Negative, 0, 0}
	n.AdjustToZero()
	assert.Equal(t, Zero, n.sign)

	n = Number{Zero, 0, 0}
	n.AdjustToZero()
	assert.Equal(t, Zero, n.sign)

	n = Number{Positive, 0, 0}
	n.AdjustToZero()
	assert.Equal(t, Zero, n.sign)

	n = Number{Negative, 1, 0}
	n.AdjustToZero()
	assert.Equal(t, Negative, n.sign)

	n = Number{Zero, 1, 0}
	n.AdjustToZero()
	assert.Equal(t, Zero, n.sign)

	n = Number{Positive, 1, 0}
	n.AdjustToZero()
	assert.Equal(t, Positive, n.sign)
}

func TestAdjustedToPositive_(t *testing.T) {
	assert.Equal(t, Positive, ofHexInternal(Positive, 1, 0).AdjustedToPositive())
	assert.Equal(t, Positive, ofHexInternal(Zero, 0, 0).AdjustedToPositive())
	assert.Equal(t, Positive, ofHexInternal(Positive, 1, 0).AdjustedToPositive())
}

func TestConvertDecimals_(t *testing.T) {
	n := funcs.MustValue(OfHex(Positive, 0x12_345, 3))

	// Same decimals = no op
	orig := n
	assert.Nil(t, n.ConvertDecimals(3))
	assert.Equal(t, orig, n)

	// Fewer decimals = rounding: 12.345 => 12.35
	n = funcs.MustValue(OfHex(Positive, 0x12_345, 3))
	assert.Nil(t, n.ConvertDecimals(2))
	assert.Equal(t, funcs.MustValue(OfHex(Positive, 0x12_35, 2)), n)

	// Fewer decimals = rounding: 12.345 => 12.4
	n = funcs.MustValue(OfHex(Positive, 0x12_345, 3))
	assert.Nil(t, n.ConvertDecimals(1))
	assert.Equal(t, funcs.MustValue(OfHex(Positive, 0x12_4, 1)), n)

	// Fewer decimals = rounding: 12.345 => 12
	n = funcs.MustValue(OfHex(Positive, 0x12_345, 3))
	assert.Nil(t, n.ConvertDecimals(0))
	assert.Equal(t, funcs.MustValue(OfHex(Positive, 0x12, 0)), n)

	// Fewer decimals = rounding: 12.567 => 12.57
	n = funcs.MustValue(OfHex(Positive, 0x12_567, 3))
	assert.Nil(t, n.ConvertDecimals(2))
	assert.Equal(t, funcs.MustValue(OfHex(Positive, 0x12_57, 2)), n)

	// Fewer decimals = rounding: 12.567 => 12.6
	n = funcs.MustValue(OfHex(Positive, 0x12_567, 3))
	assert.Nil(t, n.ConvertDecimals(1))
	assert.Equal(t, funcs.MustValue(OfHex(Positive, 0x12_6, 1)), n)

	// Fewer decimals = rounding: 12.567 => 13
	n = funcs.MustValue(OfHex(Positive, 0x12_567, 3))
	assert.Nil(t, n.ConvertDecimals(0))
	assert.Equal(t, funcs.MustValue(OfHex(Positive, 0x13, 0)), n)

	// Fewer decimals = rounding: 999_999_999_999_999.9 => 1_000_000_000_000_000
	n = funcs.MustValue(OfHex(Positive, 0x999_999_999_999_999_9, 1))
	assert.Nil(t, n.ConvertDecimals(0))
	assert.Equal(t, funcs.MustValue(OfHex(Positive, 0x1_000_000_000_000_000, 0)), n)

	// More decimals = trailing zeros: 12.345 => 12.3450
	n = funcs.MustValue(OfHex(Positive, 0x12_345, 3))
	assert.Nil(t, n.ConvertDecimals(4))
	assert.Equal(t, funcs.MustValue(OfHex(Positive, 0x12_3450, 4)), n)

	// More decimals = trailing zeros: 999_999_999_999_999 => 9_999_999_999_999_990
	n = funcs.MustValue(OfHex(Positive, 0x999_999_999_999_999, 0))
	assert.Nil(t, n.ConvertDecimals(1))
	assert.Equal(t, funcs.MustValue(OfHex(Positive, 0x9_999_999_999_999_990, 1)), n)

	// Error: decimals > 16
	n.ConvertDecimals(17)
	assert.Equal(t, fmt.Errorf("Invalid number of decimals 17: the valid range is [0 .. 16]"), n.ConvertDecimals(17))

	// Error: more decimals loses a significant digit
	n = funcs.MustValue(OfHex(Positive, 0x9_999_999_999_999_999, 0))
	assert.Equal(t, fmt.Errorf("Cannot convert 9999999999999999 to 1 decimal(s), as significant leading digits would be lost"), n.ConvertDecimals(1))
}

func TestCmp_(t *testing.T) {
	// a == b
	a, b := funcs.MustValue(OfString("5")), funcs.MustValue(OfString("5"))
	assert.Equal(t, 0, a.Cmp(b))

	// a.sign < b.sign (- < 0)
	a, b = funcs.MustValue(OfString("-5")), funcs.MustValue(OfString("0"))
	assert.Equal(t, -1, a.Cmp(b))

	// a.sign < b.sign (0 < +)
	a, b = funcs.MustValue(OfString("0")), funcs.MustValue(OfString("5"))
	assert.Equal(t, -1, a.Cmp(b))

	// a.sign > b.sign (0 > -)
	a, b = funcs.MustValue(OfString("0")), funcs.MustValue(OfString("-5"))
	assert.Equal(t, 1, a.Cmp(b))

	// a.sign > b.sign (+ > 0)
	a, b = funcs.MustValue(OfString("5")), funcs.MustValue(OfString("0"))
	assert.Equal(t, 1, a.Cmp(b))

  // a.sign == b.sign, a.decimals == b.decimals

	// a.digits < b.digits, a +
	a, b = funcs.MustValue(OfString("5")), funcs.MustValue(OfString("9"))
	assert.Equal(t, -1, a.Cmp(b))

	// a.digits < b.digits, a -
	a, b = funcs.MustValue(OfString("-5")), funcs.MustValue(OfString("-9"))
	assert.Equal(t, +1, a.Cmp(b))

	// a.digits > b.digits, a +
	a, b = funcs.MustValue(OfString("9")), funcs.MustValue(OfString("5"))
	assert.Equal(t, 1, a.Cmp(b))

	// a.digits > b.digits, a -
	a, b = funcs.MustValue(OfString("-9")), funcs.MustValue(OfString("-5"))
	assert.Equal(t, -1, a.Cmp(b))

  // a.sign == b.sign, a.decimals != b.decimals, a == 0 and/or b == 0

  // a == 0, b == 0
	a, b = funcs.MustValue(OfString("0")), funcs.MustValue(OfString("0.0"))
	assert.Equal(t, 0, a.Cmp(b))

  // a != 0, b == 0, a +
	a, b = funcs.MustValue(OfString("1.0")), funcs.MustValue(OfString("0"))
	assert.Equal(t, 1, a.Cmp(b))

  // a != 0, b == 0, a -
	a, b = funcs.MustValue(OfString("-1.0")), funcs.MustValue(OfString("0"))
	assert.Equal(t, -1, a.Cmp(b))

  // a == 0, b != 0, b +
	a, b = funcs.MustValue(OfString("0")), funcs.MustValue(OfString("1.0"))
	assert.Equal(t, -1, a.Cmp(b))

  // a == 0, b != 0, b -
	a, b = funcs.MustValue(OfString("0")), funcs.MustValue(OfString("-1.0"))
	assert.Equal(t, 1, a.Cmp(b))

	// a.sign == b.sign, a.decimals != b.decimals, a.digits != 0, b.digits != 0

  // a int < b int, +
  a, b = funcs.MustValue(OfString("1.0")), funcs.MustValue(OfString("2.00"))
  assert.Equal(t, -1, a.Cmp(b))

  // a int < b int, -
  a, b = funcs.MustValue(OfString("-1.0")), funcs.MustValue(OfString("-2.00"))
  assert.Equal(t, 1, a.Cmp(b))

  // a int > b int, +
  a, b = funcs.MustValue(OfString("2.0")), funcs.MustValue(OfString("1.00"))
  assert.Equal(t, -1, a.Cmp(b))

  // a int > b int, -
  a, b = funcs.MustValue(OfString("-2.0")), funcs.MustValue(OfString("-1.00"))
  assert.Equal(t, 1, a.Cmp(b))

	// a.sign == b.sign, a.decimals != b.decimals, a.digits != 0, b.digits != 0, a int == b int, a.decimals > b.decimals

  // a frac < b frac, +
  a, b = funcs.MustValue(OfString("1.00")), funcs.MustValue(OfString("1.1"))
  assert.Equal(t, -1, a.Cmp(b))

  // a frac < b frac, -
  a, b = funcs.MustValue(OfString("-1.00")), funcs.MustValue(OfString("-1.1"))
  assert.Equal(t, 1, a.Cmp(b))

  // a frac > b frac, +
  a, b = funcs.MustValue(OfString("1.01")), funcs.MustValue(OfString("1.0"))
  assert.Equal(t, 1, a.Cmp(b))

  // a frac > b frac, -
  a, b = funcs.MustValue(OfString("-1.01")), funcs.MustValue(OfString("-1.0"))
  assert.Equal(t, -1, a.Cmp(b))

	// a.sign == b.sign, a.decimals != b.decimals, a.digits != 0, b.digits != 0, a int == b int, a.decimals < b.decimals

  // a frac < b frac, +
  a, b = funcs.MustValue(OfString("1.0")), funcs.MustValue(OfString("1.01"))
  assert.Equal(t, -1, a.Cmp(b))

  // a frac < b frac, -
  a, b = funcs.MustValue(OfString("-1.0")), funcs.MustValue(OfString("-1.01"))
  assert.Equal(t, 1, a.Cmp(b))

  // a frac > b frac, +
  a, b = funcs.MustValue(OfString("1.01")), funcs.MustValue(OfString("1.0"))
  assert.Equal(t, 1, a.Cmp(b))

  // a frac > b frac, -
  a, b = funcs.MustValue(OfString("-1.01")), funcs.MustValue(OfString("-1.0"))
  assert.Equal(t, -1, a.Cmp(b))

  // a frac = b frac, +
  a, b = funcs.MustValue(OfString("1.10")), funcs.MustValue(OfString("1.1"))
  assert.Equal(t, 0, a.Cmp(b))

  // All numbers of decimal positions

  // 1 < 1.2
	a, b = funcs.MustValue(OfString("1")), funcs.MustValue(OfString("1.2"))
	assert.Equal(t, -1, a.Cmp(b))

  // 1.2 > 1
	a, b = funcs.MustValue(OfString("1.2")), funcs.MustValue(OfString("1"))
	assert.Equal(t, 1, a.Cmp(b))

  // 12 < 12.3
	a, b = funcs.MustValue(OfString("12")), funcs.MustValue(OfString("12.3"))
	assert.Equal(t, -1, a.Cmp(b))

  // 12.3 > 12
	a, b = funcs.MustValue(OfString("12.3")), funcs.MustValue(OfString("12"))
	assert.Equal(t, 1, a.Cmp(b))

  // 123 < 123.4
	a, b = funcs.MustValue(OfString("123")), funcs.MustValue(OfString("123.4"))
	assert.Equal(t, -1, a.Cmp(b))

  // 123.4 > 123
	a, b = funcs.MustValue(OfString("123.4")), funcs.MustValue(OfString("123"))
	assert.Equal(t, 1, a.Cmp(b))

  // 1234 < 1234.5
	a, b = funcs.MustValue(OfString("1234")), funcs.MustValue(OfString("1234.5"))
	assert.Equal(t, -1, a.Cmp(b))

  // 1234.5 > 1234
	a, b = funcs.MustValue(OfString("1234.5")), funcs.MustValue(OfString("1234"))
	assert.Equal(t, 1, a.Cmp(b))

  // 12345 < 12345.6
	a, b = funcs.MustValue(OfString("12345")), funcs.MustValue(OfString("12345.6"))
	assert.Equal(t, -1, a.Cmp(b))

  // 12345.6 > 12345
	a, b = funcs.MustValue(OfString("12345.6")), funcs.MustValue(OfString("12345"))
	assert.Equal(t, 1, a.Cmp(b))

  // 123456 < 123456.7
	a, b = funcs.MustValue(OfString("123456")), funcs.MustValue(OfString("123456.7"))
	assert.Equal(t, -1, a.Cmp(b))

  // 123456.7 > 123456
	a, b = funcs.MustValue(OfString("123456.7")), funcs.MustValue(OfString("123456"))
	assert.Equal(t, 1, a.Cmp(b))

  // 1234567 < 1234567.8
	a, b = funcs.MustValue(OfString("1234567")), funcs.MustValue(OfString("1234567.8"))
	assert.Equal(t, -1, a.Cmp(b))

  // 1234567.8 > 1234567
	a, b = funcs.MustValue(OfString("1234567.8")), funcs.MustValue(OfString("1234567"))
	assert.Equal(t, 1, a.Cmp(b))

  // 12345678 < 12345678.9
	a, b = funcs.MustValue(OfString("12345678")), funcs.MustValue(OfString("12345678.9"))
	assert.Equal(t, -1, a.Cmp(b))

  // 12345678.9 > 12345678
	a, b = funcs.MustValue(OfString("12345678.9")), funcs.MustValue(OfString("12345678"))
	assert.Equal(t, 1, a.Cmp(b))

  // 123456789 < 123456789.1
	a, b = funcs.MustValue(OfString("123456789")), funcs.MustValue(OfString("123456789.1"))
	assert.Equal(t, -1, a.Cmp(b))

  // 123456789.1 > 123456789
	a, b = funcs.MustValue(OfString("123456789.1")), funcs.MustValue(OfString("123456789"))
	assert.Equal(t, 1, a.Cmp(b))

  // 1234567890 < 1234567890.1
	a, b = funcs.MustValue(OfString("1234567890")), funcs.MustValue(OfString("1234567890.1"))
	assert.Equal(t, -1, a.Cmp(b))

  // 1234567890.1 > 1234567890
	a, b = funcs.MustValue(OfString("1234567890.1")), funcs.MustValue(OfString("1234567890"))
	assert.Equal(t, 1, a.Cmp(b))

  // 12345678901 < 12345678901.2
	a, b = funcs.MustValue(OfString("12345678901")), funcs.MustValue(OfString("12345678901.2"))
	assert.Equal(t, -1, a.Cmp(b))

  // 12345678901.2 > 12345678901
	a, b = funcs.MustValue(OfString("12345678901.2")), funcs.MustValue(OfString("12345678901"))
	assert.Equal(t, 1, a.Cmp(b))

  // 123456789012 < 123456789012.3
	a, b = funcs.MustValue(OfString("123456789012")), funcs.MustValue(OfString("123456789012.3"))
	assert.Equal(t, -1, a.Cmp(b))

  // 123456789012.3 > 123456789012
	a, b = funcs.MustValue(OfString("123456789012.3")), funcs.MustValue(OfString("123456789012"))
	assert.Equal(t, 1, a.Cmp(b))

  // 1234567890123 < 1234567890123.4
	a, b = funcs.MustValue(OfString("1234567890123")), funcs.MustValue(OfString("1234567890123.4"))
	assert.Equal(t, -1, a.Cmp(b))

  // 1234567890123.4 > 1234567890123
	a, b = funcs.MustValue(OfString("1234567890123.4")), funcs.MustValue(OfString("1234567890123"))
	assert.Equal(t, 1, a.Cmp(b))

  // 12345678901234 < 12345678901234.5
	a, b = funcs.MustValue(OfString("1234567890124")), funcs.MustValue(OfString("12345678901234.5"))
	assert.Equal(t, -1, a.Cmp(b))

  // 12345678901234.5 > 12345678901234
	a, b = funcs.MustValue(OfString("12345678901234.5")), funcs.MustValue(OfString("12345678901234"))
	assert.Equal(t, 1, a.Cmp(b))

  // 123456789012345 < 123456789012345.6
	a, b = funcs.MustValue(OfString("12345678901245")), funcs.MustValue(OfString("123456789012345.6"))
	assert.Equal(t, -1, a.Cmp(b))

  // 123456789012345.6 > 123456789012345
	a, b = funcs.MustValue(OfString("123456789012345.6")), funcs.MustValue(OfString("123456789012345"))
	assert.Equal(t, 1, a.Cmp(b))

  // 1234567890123456 < 1234567890123457
	a, b = funcs.MustValue(OfString("123456789012456")), funcs.MustValue(OfString("1234567890123457"))
	assert.Equal(t, -1, a.Cmp(b))

  // 1234567890123457 > 1234567890123456
	a, b = funcs.MustValue(OfString("1234567890123457")), funcs.MustValue(OfString("1234567890123456"))
	assert.Equal(t, 1, a.Cmp(b))
}

func TestAdd_(t *testing.T) {
	a, b := funcs.MustValue(OfString("9")), funcs.MustValue(OfString("5"))
	assert.Equal(t, union.OfResultError(OfString("14")), union.OfResultError(a.Add(b)))
	assert.Equal(t, a, funcs.MustValue(OfString("9")))
	assert.Equal(t, b, funcs.MustValue(OfString("5")))

	assert.Equal(
		t,
		union.OfResultError(OfString("14")),
		union.OfResultError(funcs.MustValue(OfString("5")).Add(funcs.MustValue(OfString("9")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("-14")),
		union.OfResultError(funcs.MustValue(OfString("-9")).Add(funcs.MustValue(OfString("-5")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("-14")),
		union.OfResultError(funcs.MustValue(OfString("-5")).Add(funcs.MustValue(OfString("-9")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("4")),
		union.OfResultError(funcs.MustValue(OfString("9")).Add(funcs.MustValue(OfString("-5")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("-4")),
		union.OfResultError(funcs.MustValue(OfString("5")).Add(funcs.MustValue(OfString("-9")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("-4")),
		union.OfResultError(funcs.MustValue(OfString("-9")).Add(funcs.MustValue(OfString("5")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("4")),
		union.OfResultError(funcs.MustValue(OfString("-5")).Add(funcs.MustValue(OfString("9")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("24.6912")),
		union.OfResultError(funcs.MustValue(OfString("12.3456")).Add(funcs.MustValue(OfString("12.3456")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("-24.6912")),
		union.OfResultError(funcs.MustValue(OfString("-12.3456")).Add(funcs.MustValue(OfString("-12.3456")))),
	)

	// Try adding 0 - 999 with 0 - 999
	var (
		istr, jstr, kstr string
		c                Number
	)
	for i := 0; i <= 999; i++ {
		for j := 0; j <= 999; j++ {
			conv.To(i, &istr)
			conv.To(j, &jstr)
			conv.To(i+j, &kstr)
			a, b, c = funcs.MustValue(OfString(istr)), funcs.MustValue(OfString(jstr)), funcs.MustValue(OfString(kstr))
			assert.Equal(t, union.OfResult(c), union.OfResultError(a.Add(b)))
		}
	}

	assert.Equal(
		t,
		fmt.Errorf("Invalid Number pair: the number of decimals do not match (1 and 2)"),
		union.OfResultError(funcs.MustValue(OfString("1.1")).Add(funcs.MustValue(OfString("1.12")))).Error(),
	)

	assert.Equal(
		t,
		fmt.Errorf("Overflow adding 9000000000000000 to 1000000000000000"),
		union.OfResultError(funcs.MustValue(OfString("9000000000000000")).Add(funcs.MustValue(OfString("1000000000000000")))).Error(),
	)

	assert.Equal(
		t,
		fmt.Errorf("Underflow adding -9000000000000000 to -1000000000000000"),
		union.OfResultError(funcs.MustValue(OfString("-9000000000000000")).Add(funcs.MustValue(OfString("-1000000000000000")))).Error(),
	)
}

func TestSub_(t *testing.T) {
	a, b := funcs.MustValue(OfString("9")), funcs.MustValue(OfString("5"))
	assert.Equal(t, union.OfResultError(OfString("4")), union.OfResultError(a.Sub(b)))
	assert.Equal(t, a, funcs.MustValue(OfString("9")))
	assert.Equal(t, b, funcs.MustValue(OfString("5")))

	assert.Equal(
		t,
		union.OfResultError(OfString("-4")),
		union.OfResultError(funcs.MustValue(OfString("5")).Sub(funcs.MustValue(OfString("9")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("-4")),
		union.OfResultError(funcs.MustValue(OfString("-9")).Sub(funcs.MustValue(OfString("-5")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("4")),
		union.OfResultError(funcs.MustValue(OfString("-5")).Sub(funcs.MustValue(OfString("-9")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("14")),
		union.OfResultError(funcs.MustValue(OfString("9")).Sub(funcs.MustValue(OfString("-5")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("14")),
		union.OfResultError(funcs.MustValue(OfString("5")).Sub(funcs.MustValue(OfString("-9")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("-14")),
		union.OfResultError(funcs.MustValue(OfString("-9")).Sub(funcs.MustValue(OfString("5")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("-14")),
		union.OfResultError(funcs.MustValue(OfString("-5")).Sub(funcs.MustValue(OfString("9")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("24.6912")),
		union.OfResultError(funcs.MustValue(OfString("12.3456")).Sub(funcs.MustValue(OfString("-12.3456")))),
	)

	assert.Equal(
		t,
		union.OfResultError(OfString("-24.6912")),
		union.OfResultError(funcs.MustValue(OfString("-12.3456")).Sub(funcs.MustValue(OfString("12.3456")))),
	)

	// Try subtracting 0 - 999 from 0 - 999
	var (
		istr, jstr, kstr string
		c                Number
	)
	for i := 0; i <= 999; i++ {
		for j := 0; j <= 999; j++ {
			conv.To(i, &istr)
			conv.To(j, &jstr)
			conv.To(i-j, &kstr)
			a, b, c = funcs.MustValue(OfString(istr)), funcs.MustValue(OfString(jstr)), funcs.MustValue(OfString(kstr))
			assert.Equal(t, union.OfResult(c), union.OfResultError(a.Sub(b)))
		}
	}

	// Decimals mismatch
	assert.Equal(
		t,
		fmt.Errorf("Invalid Number pair: the number of decimals do not match (1 and 2)"),
		union.OfResultError(funcs.MustValue(OfString("1.1")).Sub(funcs.MustValue(OfString("1.12")))).Error(),
	)

	// Underflow subtracting positive from negative
	assert.Equal(
		t,
		fmt.Errorf("Underflow subtracting 0.1 from -999999999999999.9"),
		union.OfResultError(funcs.MustValue(OfString("-999999999999999.9")).Sub(funcs.MustValue(OfString("0.1")))).Error(),
	)
}
