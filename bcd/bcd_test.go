package bcd

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

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
	assert.Equal(t, union.OfError[Number](fmt.Errorf(numberDecimalsErrMsg, 17)), union.OfResultError(OfHex(Positive, 0x1, 17)))

	// Invalid digit
	assert.Equal(t, union.OfError[Number](fmt.Errorf(numberDigitsErrMsg, 0x1A)), union.OfResultError(OfHex(Positive, 0x1A, 0)))
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

func TestString_(t *testing.T) {
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
	assert.Equal(t, Zero, ofHexInternal(Positive, 0, 0).sign)
	assert.Equal(t, Zero, ofHexInternal(Zero, 0, 0).sign)
	assert.Equal(t, Zero, ofHexInternal(Negative, 0, 0).sign)

	assert.Equal(t, Positive, ofHexInternal(Positive, 1, 0).sign)
	assert.Equal(t, Zero, ofHexInternal(Positive, 0, 0).sign)
	assert.Equal(t, Negative, ofHexInternal(Negative, 1, 0).sign)
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
  // Both positive
  a,b := funcs.MustValue(OfString("5")), funcs.MustValue(OfString("4"))
  assert.Equal(t,union.OfResult(1), union.OfResultError(a.Cmp(b)))

  a,b = funcs.MustValue(OfString("5")), funcs.MustValue(OfString("5"))
  assert.Equal(t,union.OfResult(0), union.OfResultError(a.Cmp(b)))

  a,b = funcs.MustValue(OfString("5")), funcs.MustValue(OfString("6"))
  assert.Equal(t,union.OfResult(-1), union.OfResultError(a.Cmp(b)))

  // Both negative
  a,b = funcs.MustValue(OfString("-12.34")), funcs.MustValue(OfString("-12.35"))
  assert.Equal(t,union.OfResult(1), union.OfResultError(a.Cmp(b)))

  a,b = funcs.MustValue(OfString("-12.34")), funcs.MustValue(OfString("-12.34"))
  assert.Equal(t,union.OfResult(0), union.OfResultError(a.Cmp(b)))

  a,b = funcs.MustValue(OfString("-12.34")), funcs.MustValue(OfString("-12.33"))
  assert.Equal(t,union.OfResult(-1), union.OfResultError(a.Cmp(b)))

  // Different signs
  a,b = funcs.MustValue(OfString("5")), funcs.MustValue(OfString("-4"))
  assert.Equal(t,union.OfResult(1), union.OfResultError(a.Cmp(b)))

  a,b = funcs.MustValue(OfString("-5")), funcs.MustValue(OfString("4"))
  assert.Equal(t,union.OfResult(-1), union.OfResultError(a.Cmp(b)))

  // Compare to zero
  a,b = funcs.MustValue(OfString("5")), funcs.MustValue(OfString("0"))
  assert.Equal(t,union.OfResult(1), union.OfResultError(a.Cmp(b)))

  a,b = funcs.MustValue(OfString("0")), funcs.MustValue(OfString("5"))
  assert.Equal(t,union.OfResult(-1), union.OfResultError(a.Cmp(b)))

  a,b = funcs.MustValue(OfString("-5")), funcs.MustValue(OfString("0"))
  assert.Equal(t,union.OfResult(-1), union.OfResultError(a.Cmp(b)))

  a,b = funcs.MustValue(OfString("0")), funcs.MustValue(OfString("-5"))
  assert.Equal(t,union.OfResult(1), union.OfResultError(a.Cmp(b)))

  a,b = funcs.MustValue(OfString("0")), funcs.MustValue(OfString("0"))
  assert.Equal(t,union.OfResult(0), union.OfResultError(a.Cmp(b)))

  // Error - number of decimals differ
  a,b = funcs.MustValue(OfString("5.1")), funcs.MustValue(OfString("5.12"))
  assert.Equal(
    t,
    union.OfError[int](fmt.Errorf("Invalid Number pair: the number of decimals do not match (%d and %d)", 1, 2)),
    union.OfResultError(a.Cmp(b)),
  )
}

func TestAdd_(t *testing.T) {
  a,b := funcs.MustValue(OfString("1")), funcs.MustValue(OfString("2"))
  assert.Equal(t, union.OfResultError(OfString("3")), union.OfResultError(a.Add(b)))
  assert.Equal(t, a, funcs.MustValue(OfString("1")))
  assert.Equal(t, b, funcs.MustValue(OfString("2")))

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

  assert.Equal(
    t,
    fmt.Errorf("Overflow adding %s to %s", "9000000000000000", "1000000000000000"),
    union.OfResultError(funcs.MustValue(OfString("9000000000000000")).Add(funcs.MustValue(OfString("1000000000000000")))).Error(),
  )

  assert.Equal(
    t,
    fmt.Errorf("Underflow adding %s to %s", "-9000000000000000", "-1000000000000000"),
    union.OfResultError(funcs.MustValue(OfString("-9000000000000000")).Add(funcs.MustValue(OfString("-1000000000000000")))).Error(),
  )
}
