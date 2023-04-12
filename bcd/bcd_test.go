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

func TestOf_(t *testing.T) {
  // == Valid strings

  // -1
  assert.Equal(t, union.OfResult(Fixed{Negative, 0x1, 0}), union.OfResultError(Of("-1")))

  // 0
  assert.Equal(t, union.OfResult(Fixed{Zero, 0x0, 0}), union.OfResultError(Of("0")))

  // 1
  assert.Equal(t, union.OfResult(Fixed{Positive, 0x1, 0}), union.OfResultError(Of("1")))

  // 123
  assert.Equal(t, union.OfResult(Fixed{Positive, 0x123, 0}), union.OfResultError(Of("123")))

  // 0.456
  assert.Equal(t, union.OfResult(Fixed{Positive, 0x456, 3}), union.OfResultError(Of("0.456")))

  // 123.456
  assert.Equal(t, union.OfResult(Fixed{Positive, 0x123456, 3}), union.OfResultError(Of("123.456")))

  // == Invalid strings

  for _, s := range []string{"", ".", ".1", "a", "++1", "--1"} {
    assert.Equal(t, union.OfError[Fixed](fmt.Errorf(fixedErrMsg, s)), union.OfResultError(Of(s)))
  }
}
