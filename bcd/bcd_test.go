package bcd

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

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

func TestSignAdjust_(t *testing.T) {
  s := Negative
  s.Adjust(0)
  assert.Equal(t, Zero, s)

  s = Zero
  s.Adjust(0)
  assert.Equal(t, Zero, s)

  s = Positive
  s.Adjust(0)
  assert.Equal(t, Zero, s)

  s = Negative
  s.Adjust(1)
  assert.Equal(t, Negative, s)

  s = Positive
  s.Adjust(1)
  assert.Equal(t, Positive, s)
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
  assert.Equal(t, union.OfError[Number](fmt.Errorf(fixedNumberDecimalsErrMsg, 17)), union.OfResultError(OfHex(Positive, 0x1, 17)))

  // Invalid digit
  assert.Equal(t, union.OfError[Number](fmt.Errorf(fixedNumberDigitsErrMsg, 0x1A)), union.OfResultError(OfHex(Positive, 0x1A)))
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
    assert.Equal(t, union.OfError[Number](fmt.Errorf(fixedStringErrMsg, s)), union.OfResultError(OfString(s)))
  }
}
