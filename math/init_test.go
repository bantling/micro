package math

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"
	goreflect "reflect"
	"testing"

	"github.com/bantling/micro/conv"
	"github.com/stretchr/testify/assert"
)

func TestInit_(t *testing.T) {
	// ==== To Decimal

	var v Decimal

	// Signed ints
	assert.Nil(t, conv.To(1, &v))
	assert.Equal(t, Decimal{value: 1, scale: 0}, v)

	assert.Nil(t, conv.To(int8(2), &v))
	assert.Equal(t, Decimal{value: 2, scale: 0}, v)

	assert.Nil(t, conv.To(int16(3), &v))
	assert.Equal(t, Decimal{value: 3, scale: 0}, v)

	assert.Nil(t, conv.To(int32(4), &v))
	assert.Equal(t, Decimal{value: 4, scale: 0}, v)

	assert.Nil(t, conv.To(int64(5), &v))
	assert.Equal(t, Decimal{value: 5, scale: 0}, v)

	// Unigned ints
	assert.Nil(t, conv.To(uint(6), &v))
	assert.Equal(t, Decimal{value: 6, scale: 0}, v)

	assert.Nil(t, conv.To(uint8(7), &v))
	assert.Equal(t, Decimal{value: 7, scale: 0}, v)

	assert.Nil(t, conv.To(uint16(8), &v))
	assert.Equal(t, Decimal{value: 8, scale: 0}, v)

	assert.Nil(t, conv.To(uint32(9), &v))
	assert.Equal(t, Decimal{value: 9, scale: 0}, v)

	assert.Nil(t, conv.To(uint64(10), &v))
	assert.Equal(t, Decimal{value: 10, scale: 0}, v)

	if goreflect.TypeOf(uint(0)).Bits() == 64 {
		var ui uint = math.MaxUint
		assert.Equal(t, fmt.Errorf("The uint value of %d cannot be converted to int64", ui), conv.To(ui, &v))
		assert.Equal(t, Decimal{value: 10, scale: 0}, v)
	}

	var ui64 uint64 = math.MaxUint64
	assert.Equal(t, fmt.Errorf("The uint64 value of %d cannot be converted to int64", ui64), conv.To(ui64, &v))
	assert.Equal(t, Decimal{value: 10, scale: 0}, v)

	// *big.Int and *big.Rat
	var bi *big.Int
	conv.To(11, &bi)
	assert.Nil(t, conv.To(bi, &v))
	assert.Equal(t, Decimal{value: 11, scale: 0}, v)

	conv.To(ui64, &bi)
	assert.Equal(t, fmt.Errorf("The *big.Int value of %d cannot be converted to int64", ui64), conv.To(bi, &v))
	assert.Equal(t, Decimal{value: 11, scale: 0}, v)

	var br *big.Rat
	conv.To(12, &br)
	assert.Nil(t, conv.To(br, &v))
	assert.Equal(t, Decimal{value: 12, scale: 0}, v)

	conv.To(ui64, &br)
	assert.Equal(t, fmt.Errorf("The *big.Rat value of %d/1 cannot be converted to int64", ui64), conv.To(br, &v))
	assert.Equal(t, Decimal{value: 12, scale: 0}, v)

	// String
	conv.To("123.456", &v)
	assert.Equal(t, Decimal{value: 123456, scale: 3}, v)

	assert.Equal(t, fmt.Errorf("The string value 1234567890123456789 is not a valid decimal string"), conv.To("1234567890123456789", &v))
	assert.Equal(t, Decimal{value: 0, scale: 0}, v)

	// From Decimal

	// *big.Int and *big.Rat
	assert.Nil(t, conv.To(MustDecimal(123, 0), &bi))
	assert.Equal(t, int64(123), bi.Int64())

	assert.Equal(t, fmt.Errorf("The decimal value 12.3 cannot be converted to a *big.Int"), conv.To(MustDecimal(123, 1), &bi))

	assert.Nil(t, conv.To(MustDecimal(123, 1), &br))
	assert.Equal(t, "123/10", br.String())

	// string
	var str string
	assert.Nil(t, conv.To(MustDecimal(-123, 2), &str))
	assert.Equal(t, "-1.23", str)
}
