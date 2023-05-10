
package bcd

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/tuple"
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

func TestStateString_(t *testing.T) {
  assert.Equal(t, "Normal", Normal.String())
  assert.Equal(t, "Overflow", Overflow.String())
  assert.Equal(t, "Underflow", Underflow.String())
}

func TestOfHexInternal_(t *testing.T) {
	// -1
	assert.Equal(t, Number{Negative, 0x1, 0, Normal, ""}, ofHexInternal(Negative, 0x1, 0))

	// 0
	assert.Equal(t, Number{Zero, 0x0, 0, Normal, ""}, ofHexInternal(Negative, 0x0, 0))

	// 1
	assert.Equal(t, Number{Positive, 0x1, 0, Normal, ""}, ofHexInternal(Positive, 0x1, 0))

	// 0.456
	assert.Equal(t, Number{Positive, 0x456, 3, Normal, ""}, ofHexInternal(Positive, 0x456, 3))

	// 123.456
	assert.Equal(t, Number{Positive, 0x123456, 3, Normal, ""}, ofHexInternal(Positive, 0x123_456, 3))

  // Verify adjusted sign
	assert.Equal(t, Zero, ofHexInternal(Negative, 0, 0).sign)
	assert.Equal(t, Zero, ofHexInternal(Zero, 0, 0).sign)
	assert.Equal(t, Zero, ofHexInternal(Positive, 0, 0).sign)
	assert.Equal(t, Negative, ofHexInternal(Negative, 1, 0).sign)
	assert.Equal(t, Positive, ofHexInternal(Zero, 1, 0).sign)
	assert.Equal(t, Positive, ofHexInternal(Positive, 1, 0).sign)
}

func TestOfHex_(t *testing.T) {
	// -1
	assert.Equal(t, union.OfResult(Number{Negative, 0x1, 0, Normal, ""}), union.OfResultError(OfHex(Negative, 0x1, 0)))

	// 0
	assert.Equal(t, Number{Zero, 0x0, 0, Normal, ""}, MustHex(Negative, 0x0, 0))

	// 1
	assert.Equal(t, Number{Positive, 0x1, 0, Normal, ""}, MustHex(Positive, 0x1, 0))

	// 0.456
	assert.Equal(t, Number{Positive, 0x456, 3, Normal, ""}, MustHex(Positive, 0x456, 3))

	// 123.456
	assert.Equal(t, Number{Positive, 0x123456, 3, Normal, ""}, MustHex(Positive, 0x123_456, 3))

	// Invalid number of decimals
	assert.Equal(t, union.OfError[Number](fmt.Errorf("Invalid number of decimals 17: the valid range is [0 .. 16]")), union.OfResultError(OfHex(Positive, 0x1, 17)))

	// Invalid digit
	assert.Equal(t, union.OfError[Number](fmt.Errorf(`Invalid Number "0x1A": the value must contain only decimal digits for each hex group`)), union.OfResultError(OfHex(Positive, 0x1A, 0)))
}

func TestOfString_(t *testing.T) {
	// -1
	assert.Equal(t, Number{Negative, 0x1, 0, Normal, ""}, MustString("-1"))

	// 0
	assert.Equal(t, Number{Zero, 0x0, 0, Normal, ""}, MustString("0"))

	// 1
	assert.Equal(t, Number{Positive, 0x1, 0, Normal, ""}, MustString("1"))

	// 123
	assert.Equal(t, Number{Positive, 0x123, 0, Normal, ""}, MustString("123"))

	// 0.456
	assert.Equal(t, Number{Positive, 0x456, 3, Normal, ""}, MustString("0.456"))

	// 123.456
	assert.Equal(t, Number{Positive, 0x123456, 3, Normal, ""}, MustString("123.456"))

	// Invalid strings
	for _, s := range []string{"", ".", ".1", "a", "++1", "--1"} {
		assert.Equal(t, union.OfError[Number](fmt.Errorf(numberStringErrMsg, s)), union.OfResultError(OfString(s)))
	}
}

func TestNumberString_(t *testing.T) {
	// 0
	assert.Equal(t, "0", MustString("0").String())

	// 1
	assert.Equal(t, "1", MustString("1").String())

	// -1
	assert.Equal(t, "-1", MustString("-1").String())

	// 0.1
	assert.Equal(t, "0.1", MustString("0.1").String())

	// 0.123
	assert.Equal(t, "0.123", MustString("0.123").String())

	// 1.23
	assert.Equal(t, "1.23", MustString("1.23").String())

	// 12.3
	assert.Equal(t, "12.3", MustString("12.3").String())

	// 123
	assert.Equal(t, "123", MustString("123").String())

	// 123.4
	assert.Equal(t, "123.4", MustString("123.4").String())

	// 123.45
	assert.Equal(t, "123.45", MustString("123.45").String())

	// 123.456
	assert.Equal(t, "123.456", MustString("123.456").String())
}

func TestAdjustedToPositive_(t *testing.T) {
	assert.Equal(t, Positive, ofHexInternal(Positive, 1, 0).AdjustedToPositive())
	assert.Equal(t, Positive, ofHexInternal(Zero, 0, 0).AdjustedToPositive())
	assert.Equal(t, Positive, ofHexInternal(Positive, 1, 0).AdjustedToPositive())
}

func TestNumberNegate_(t *testing.T) {
  assert.Equal(t, MustString("1"), MustString("-1").Negate())
  assert.Equal(t, MustString("0"), MustString("0").Negate())
  assert.Equal(t, MustString("-1"), MustString("1").Negate())
}

func TestNumberState_(t *testing.T) {
  n := Number{Zero, 0, 0, Normal, ""}
  assert.True(t, n.IsNormal())
  assert.Equal(t, Normal, n.State())
  assert.Equal(t, "", n.StateMsg())

  n = Number{Positive, 1, 0, Overflow, "Overflowed"}
  assert.False(t, n.IsNormal())
  assert.Equal(t, Overflow, n.State())
  assert.Equal(t, "Overflowed", n.StateMsg())

  n = Number{Positive, 1, 0, Underflow, "Underflowed"}
  assert.False(t, n.IsNormal())
  assert.Equal(t, Underflow, n.State())
  assert.Equal(t, "Underflowed", n.StateMsg())
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

func TestAlignDecimals_(t *testing.T) {
  // Same decimals
  a, b := MustString("1.23"), MustString("2.34")
  assert.Equal(t, tuple.Of2(a, b), tuple.Of2(alignDecimals(a, b)))

  // a decimals > b
  a, b = MustString("1.234"), MustString("2.34")
  assert.Equal(t, tuple.Of2(a, MustString("2.340")), tuple.Of2(alignDecimals(a, b)))

  // a decimals > b, extending b would lose precision, requiring rounding of a
  a, b = MustString("1.234567890123456"), MustString("23.45678901234567")
  assert.Equal(t, tuple.Of2(MustString("1.23456789012346"), b), tuple.Of2(alignDecimals(a, b)))

  // a < b
  a, b = MustString("1.23"), MustString("2.345")
  assert.Equal(t, tuple.Of2(MustString("1.230"), b), tuple.Of2(alignDecimals(a, b)))

  // a decimals < b, extending a would lose precision, requiring rounding of b
  a, b = MustString("12.34567890123456"), MustString("2.345678901234567")
  assert.Equal(t, tuple.Of2(a, MustString("2.34567890123457")), tuple.Of2(alignDecimals(a, b)))
}

func TestCmp_(t *testing.T) {
	// a == b
	a, b := MustString("5"), MustString("5")
	assert.Equal(t, 0, a.Cmp(b))

	// a.sign < b.sign (- < 0)
	a, b = MustString("-5"), MustString("0")
	assert.Equal(t, -1, a.Cmp(b))

	// a.sign < b.sign (0 < +)
	a, b = MustString("0"), MustString("5")
	assert.Equal(t, -1, a.Cmp(b))

	// a.sign > b.sign (0 > -)
	a, b = MustString("0"), MustString("-5")
	assert.Equal(t, 1, a.Cmp(b))

	// a.sign > b.sign (+ > 0)
	a, b = MustString("5"), MustString("0")
	assert.Equal(t, 1, a.Cmp(b))

  // a.sign == b.sign, a.decimals == b.decimals

	// a.digits < b.digits, a +
	a, b = MustString("5"), MustString("9")
	assert.Equal(t, -1, a.Cmp(b))

	// a.digits < b.digits, a -
	a, b = MustString("-5"), MustString("-9")
	assert.Equal(t, +1, a.Cmp(b))

	// a.digits > b.digits, a +
	a, b = MustString("9"), MustString("5")
	assert.Equal(t, 1, a.Cmp(b))

	// a.digits > b.digits, a -
	a, b = MustString("-9"), MustString("-5")
	assert.Equal(t, -1, a.Cmp(b))

  // a.sign == b.sign, a.decimals != b.decimals, a == 0 and/or b == 0

  // a == 0, b == 0
	a, b = MustString("0"), MustString("0.0")
	assert.Equal(t, 0, a.Cmp(b))

  // a != 0, b == 0, a +
	a, b = MustString("1.0"), MustString("0")
	assert.Equal(t, 1, a.Cmp(b))

  // a != 0, b == 0, a -
	a, b = MustString("-1.0"), MustString("0")
	assert.Equal(t, -1, a.Cmp(b))

  // a == 0, b != 0, b +
	a, b = MustString("0"), MustString("1.0")
	assert.Equal(t, -1, a.Cmp(b))

  // a == 0, b != 0, b -
	a, b = MustString("0"), MustString("-1.0")
	assert.Equal(t, 1, a.Cmp(b))

	// a.sign == b.sign, a.decimals != b.decimals, a.digits != 0, b.digits != 0

  // a int < b int, +
  a, b = MustString("1.0"), MustString("2.00")
  assert.Equal(t, -1, a.Cmp(b))

  // a int < b int, -
  a, b = MustString("-1.0"), MustString("-2.00")
  assert.Equal(t, 1, a.Cmp(b))

  // a int > b int, +
  a, b = MustString("2.0"), MustString("1.00")
  assert.Equal(t, -1, a.Cmp(b))

  // a int > b int, -
  a, b = MustString("-2.0"), MustString("-1.00")
  assert.Equal(t, 1, a.Cmp(b))

	// a.sign == b.sign, a.decimals != b.decimals, a.digits != 0, b.digits != 0, a int == b int, a.decimals > b.decimals

  // a frac < b frac, +
  a, b = MustString("1.00"), MustString("1.1")
  assert.Equal(t, -1, a.Cmp(b))

  // a frac < b frac, -
  a, b = MustString("-1.00"), MustString("-1.1")
  assert.Equal(t, 1, a.Cmp(b))

  // a frac > b frac, +
  a, b = MustString("1.01"), MustString("1.0")
  assert.Equal(t, 1, a.Cmp(b))

  // a frac > b frac, -
  a, b = MustString("-1.01"), MustString("-1.0")
  assert.Equal(t, -1, a.Cmp(b))

	// a.sign == b.sign, a.decimals != b.decimals, a.digits != 0, b.digits != 0, a int == b int, a.decimals < b.decimals

  // a frac < b frac, +
  a, b = MustString("1.0"), MustString("1.01")
  assert.Equal(t, -1, a.Cmp(b))

  // a frac < b frac, -
  a, b = MustString("-1.0"), MustString("-1.01")
  assert.Equal(t, 1, a.Cmp(b))

  // a frac > b frac, +
  a, b = MustString("1.01"), MustString("1.0")
  assert.Equal(t, 1, a.Cmp(b))

  // a frac > b frac, -
  a, b = MustString("-1.01"), MustString("-1.0")
  assert.Equal(t, -1, a.Cmp(b))

  // a frac = b frac, +
  a, b = MustString("1.10"), MustString("1.1")
  assert.Equal(t, 0, a.Cmp(b))

  // All numbers of decimal positions

  // 1 < 1.2
	a, b = MustString("1"), MustString("1.2")
	assert.Equal(t, -1, a.Cmp(b))

  // 1.2 > 1
	a, b = MustString("1.2"), MustString("1")
	assert.Equal(t, 1, a.Cmp(b))

  // 12 < 12.3
	a, b = MustString("12"), MustString("12.3")
	assert.Equal(t, -1, a.Cmp(b))

  // 12.3 > 12
	a, b = MustString("12.3"), MustString("12")
	assert.Equal(t, 1, a.Cmp(b))

  // 123 < 123.4
	a, b = MustString("123"), MustString("123.4")
	assert.Equal(t, -1, a.Cmp(b))

  // 123.4 > 123
	a, b = MustString("123.4"), MustString("123")
	assert.Equal(t, 1, a.Cmp(b))

  // 1234 < 1234.5
	a, b = MustString("1234"), MustString("1234.5")
	assert.Equal(t, -1, a.Cmp(b))

  // 1234.5 > 1234
	a, b = MustString("1234.5"), MustString("1234")
	assert.Equal(t, 1, a.Cmp(b))

  // 12345 < 12345.6
	a, b = MustString("12345"), MustString("12345.6")
	assert.Equal(t, -1, a.Cmp(b))

  // 12345.6 > 12345
	a, b = MustString("12345.6"), MustString("12345")
	assert.Equal(t, 1, a.Cmp(b))

  // 123456 < 123456.7
	a, b = MustString("123456"), MustString("123456.7")
	assert.Equal(t, -1, a.Cmp(b))

  // 123456.7 > 123456
	a, b = MustString("123456.7"), MustString("123456")
	assert.Equal(t, 1, a.Cmp(b))

  // 1234567 < 1234567.8
	a, b = MustString("1234567"), MustString("1234567.8")
	assert.Equal(t, -1, a.Cmp(b))

  // 1234567.8 > 1234567
	a, b = MustString("1234567.8"), MustString("1234567")
	assert.Equal(t, 1, a.Cmp(b))

  // 12345678 < 12345678.9
	a, b = MustString("12345678"), MustString("12345678.9")
	assert.Equal(t, -1, a.Cmp(b))

  // 12345678.9 > 12345678
	a, b = MustString("12345678.9"), MustString("12345678")
	assert.Equal(t, 1, a.Cmp(b))

  // 123456789 < 123456789.1
	a, b = MustString("123456789"), MustString("123456789.1")
	assert.Equal(t, -1, a.Cmp(b))

  // 123456789.1 > 123456789
	a, b = MustString("123456789.1"), MustString("123456789")
	assert.Equal(t, 1, a.Cmp(b))

  // 1234567890 < 1234567890.1
	a, b = MustString("1234567890"), MustString("1234567890.1")
	assert.Equal(t, -1, a.Cmp(b))

  // 1234567890.1 > 1234567890
	a, b = MustString("1234567890.1"), MustString("1234567890")
	assert.Equal(t, 1, a.Cmp(b))

  // 12345678901 < 12345678901.2
	a, b = MustString("12345678901"), MustString("12345678901.2")
	assert.Equal(t, -1, a.Cmp(b))

  // 12345678901.2 > 12345678901
	a, b = MustString("12345678901.2"), MustString("12345678901")
	assert.Equal(t, 1, a.Cmp(b))

  // 123456789012 < 123456789012.3
	a, b = MustString("123456789012"), MustString("123456789012.3")
	assert.Equal(t, -1, a.Cmp(b))

  // 123456789012.3 > 123456789012
	a, b = MustString("123456789012.3"), MustString("123456789012")
	assert.Equal(t, 1, a.Cmp(b))

  // 1234567890123 < 1234567890123.4
	a, b = MustString("1234567890123"), MustString("1234567890123.4")
	assert.Equal(t, -1, a.Cmp(b))

  // 1234567890123.4 > 1234567890123
	a, b = MustString("1234567890123.4"), MustString("1234567890123")
	assert.Equal(t, 1, a.Cmp(b))

  // 12345678901234 < 12345678901234.5
	a, b = MustString("1234567890124"), MustString("12345678901234.5")
	assert.Equal(t, -1, a.Cmp(b))

  // 12345678901234.5 > 12345678901234
	a, b = MustString("12345678901234.5"), MustString("12345678901234")
	assert.Equal(t, 1, a.Cmp(b))

  // 123456789012345 < 123456789012345.6
	a, b = MustString("12345678901245"), MustString("123456789012345.6")
	assert.Equal(t, -1, a.Cmp(b))

  // 123456789012345.6 > 123456789012345
	a, b = MustString("123456789012345.6"), MustString("123456789012345")
	assert.Equal(t, 1, a.Cmp(b))

  // 1234567890123456 < 1234567890123457
	a, b = MustString("123456789012456"), MustString("1234567890123457")
	assert.Equal(t, -1, a.Cmp(b))

  // 1234567890123457 > 1234567890123456
	a, b = MustString("1234567890123457"), MustString("1234567890123456")
	assert.Equal(t, 1, a.Cmp(b))
}

func TestAdd_(t *testing.T) {
	a, b := MustString("9"), MustString("5")
	assert.Equal(t, union.OfResultError(OfString("14")), union.OfResult(a.Add(b)))

  a, b = MustString("5"), MustString("9")
	assert.Equal(t, union.OfResultError(OfString("14")), union.OfResult(a.Add(b)))

  a, b = MustString("-9"), MustString("-5")
	assert.Equal(t, union.OfResultError(OfString("-14")), union.OfResult(a.Add(b)))

  a, b = MustString("-5"), MustString("-9")
	assert.Equal(t, union.OfResultError(OfString("-14")), union.OfResult(a.Add(b)))

  a, b = MustString("9"), MustString("-5")
	assert.Equal(t, union.OfResultError(OfString("4")), union.OfResult(a.Add(b)))

  a, b = MustString("5"), MustString("-9")
	assert.Equal(t, union.OfResultError(OfString("-4")), union.OfResult(a.Add(b)))

  a, b = MustString("-9"), MustString("5")
	assert.Equal(t, union.OfResultError(OfString("-4")), union.OfResult(a.Add(b)))

  a, b = MustString("-5"), MustString("9")
	assert.Equal(t, union.OfResultError(OfString("4")), union.OfResult(a.Add(b)))

  a, b = MustString("12.3456"), MustString("12.3456")
	assert.Equal(t, union.OfResultError(OfString("24.6912")), union.OfResult(a.Add(b)))

  a, b = MustString("-12.3456"), MustString("-12.3456")
	assert.Equal(t, union.OfResultError(OfString("-24.6912")), union.OfResult(a.Add(b)))

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
			a, b, c = MustString(istr), MustString(jstr), MustString(kstr)
			assert.Equal(t, union.OfResult(c), union.OfResult(a.Add(b)))
		}
	}

  // Overflow
  v9, v1 := MustString("9000000000000000"), MustString("1000000000000000")
  a = v9.Add(v1)
	assert.Equal(t, Number{Positive, 0x9000000000000000, 0, Overflow, "Overflow adding 1000000000000000 to 9000000000000000"}, a)

  // Overflow + Normal stays Overflow with original message
  b = a.Add(v1)
  assert.True(t, a == b)

  // Normal + Overflow stays Overflow
  a = v9.Add(a)
  assert.Equal(t, Number{Positive, 0x9000000000000000, 0, Overflow, "Overflow adding Overflow to 9000000000000000"}, a)

  // Underflow
  v9, v1 = v9.Negate(), v1.Negate()
  a = v9.Add(v1)
	assert.Equal(t, Number{Negative, 0x9000000000000000, 0, Underflow, "Underflow adding -1000000000000000 to -9000000000000000"}, a)

  // Underflow + Normal stays Underflow with original message
  b = a.Add(v1)
  assert.True(t, a == b)

  // Normal + Underflow stays Underflow
  a = v9.Add(a)
  assert.Equal(t, Number{Negative, 0x9000000000000000, 0, Underflow, "Underflow adding Underflow to -9000000000000000"}, a)
}

func TestSub_(t *testing.T) {
	a, b := MustString("9"), MustString("5")
	assert.Equal(t, MustString("4"), a.Sub(b))

  a, b = MustString("5"), MustString("9")
	assert.Equal(t, MustString("-4"), a.Sub(b))

  a, b = MustString("-9"), MustString("-5")
	assert.Equal(t, MustString("-4"), a.Sub(b))

  a, b = MustString("-5"), MustString("-9")
	assert.Equal(t, MustString("4"), a.Sub(b))

  a, b = MustString("9"), MustString("-5")
	assert.Equal(t, MustString("14"), a.Sub(b))

  a, b = MustString("5"), MustString("-9")
	assert.Equal(t, MustString("14"), a.Sub(b))

  a, b = MustString("-9"), MustString("5")
	assert.Equal(t, MustString("-14"), a.Sub(b))

  a, b = MustString("-5"), MustString("9")
	assert.Equal(t, MustString("-14"), a.Sub(b))

  a, b = MustString("12.3456"), MustString("-12.3456")
	assert.Equal(t, MustString("24.6912"), a.Sub(b))

  a, b = MustString("-12.3456"), MustString("12.3456")
	assert.Equal(t, MustString("-24.6912"), a.Sub(b))

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
			a, b, c = MustString(istr), MustString(jstr), MustString(kstr)
			assert.Equal(t, c, a.Sub(b))
		}
	}

  // Overflow
  v9, v1 := MustString("9000000000000000"), MustString("-1000000000000000")
  a = v9.Sub(v1)
	assert.Equal(t, Number{Positive, 0x9000000000000000, 0, Overflow, "Overflow subtracting -1000000000000000 from 9000000000000000"}, a)

  // Overflow - Normal stays Overflow with original message
  b = a.Sub(v1)
  assert.True(t, a == b)

  // Normal - Overflow stays Overflow
  a = v9.Sub(a)
  assert.Equal(t, Number{Positive, 0x9000000000000000, 0, Overflow, "Overflow subtracting Overflow from 9000000000000000"}, a)

  // Underflow
  v9, v1 = v9.Negate(), v1.Negate()
  a = v9.Sub(v1)
	assert.Equal(t, Number{Negative, 0x9000000000000000, 0, Underflow, "Underflow subtracting 1000000000000000 from -9000000000000000"}, a)

  // Underflow - Normal stays Underflow with original message
  b = a.Sub(v1)
  assert.True(t, a == b)

  // Normal - Underflow stays Underflow
  a = v9.Sub(a)
  assert.Equal(t, Number{Negative, 0x9000000000000000, 0, Underflow, "Underflow subtracting Underflow from -9000000000000000"}, a)
}
