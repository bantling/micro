package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"math"
	"math/big"
	"testing"

	"github.com/bantling/micro/funcs"
	"github.com/stretchr/testify/assert"
)

func TestIntToBigInt_(t *testing.T) {
	var o *big.Int
	IntToBigInt(int8(1), &o)
	assert.Equal(t, big.NewInt(1), o)

	IntToBigInt(100_000, &o)
	assert.Equal(t, big.NewInt(100_000), o)
}

func TestUintToBigInt_(t *testing.T) {
	var o *big.Int
	UintToBigInt(uint8(1), &o)
	assert.Equal(t, big.NewInt(1), o)

	UintToBigInt(uint(100_000), &o)
	assert.Equal(t, big.NewInt(100_000), o)
}

func TestFloatToBigInt_(t *testing.T) {
	var o *big.Int
	assert.Nil(t, FloatToBigInt(float32(1), &o))
	assert.Equal(t, big.NewInt(1), o)

	assert.Nil(t, FloatToBigInt(float32(100_000), &o))
	assert.Equal(t, big.NewInt(100_000), o)

	assert.Equal(t, "The float64 value of 1.25 cannot be converted to *big.Int", FloatToBigInt(1.25, &o).Error())
	assert.Equal(t, "The float64 value of +Inf cannot be converted to *big.Int", FloatToBigInt(math.Inf(1), &o).Error())
	assert.Equal(t, "The float64 value of -Inf cannot be converted to *big.Int", FloatToBigInt(math.Inf(-1), &o).Error())
	assert.Equal(t, "The float64 value of NaN cannot be converted to *big.Int", FloatToBigInt(math.NaN(), &o).Error())

	funcs.TryTo(
		func() {
			MustFloatToBigInt(math.NaN(), &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The float64 value of NaN cannot be converted to *big.Int", e.(error).Error())
		},
	)
}

func TestBigFloatToBigInt_(t *testing.T) {
	var o *big.Int
	assert.Nil(t, BigFloatToBigInt(big.NewFloat(1), &o))
	assert.Equal(t, big.NewInt(1), o)

	assert.Nil(t, BigFloatToBigInt(big.NewFloat(100_000), &o))
	assert.Equal(t, big.NewInt(100_000), o)

	assert.Equal(t, "The *big.Float value of 1.25 cannot be converted to *big.Int", BigFloatToBigInt(big.NewFloat(1.25), &o).Error())
	assert.Equal(t, "The *big.Float value of +Inf cannot be converted to *big.Int", BigFloatToBigInt(big.NewFloat(math.Inf(1)), &o).Error())
	assert.Equal(t, "The *big.Float value of -Inf cannot be converted to *big.Int", BigFloatToBigInt(big.NewFloat(math.Inf(-1)), &o).Error())

	funcs.TryTo(
		func() {
			MustBigFloatToBigInt(big.NewFloat(math.Inf(1)), &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The *big.Float value of +Inf cannot be converted to *big.Int", e.(error).Error())
		},
	)
}

func TestBigRatToBigInt_(t *testing.T) {
	var o *big.Int
	assert.Nil(t, BigRatToBigInt(big.NewRat(1, 1), &o))
	assert.Equal(t, big.NewInt(1), o)

	assert.Nil(t, BigRatToBigInt(big.NewRat(100_000, 1), &o))
	assert.Equal(t, big.NewInt(100_000), o)

	assert.Equal(t, "The *big.Rat value of 5/4 cannot be converted to *big.Int", BigRatToBigInt(big.NewRat(125, 100), &o).Error())

	funcs.TryTo(
		func() {
			MustBigRatToBigInt(big.NewRat(125, 100), &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The *big.Rat value of 5/4 cannot be converted to *big.Int", e.(error).Error())
		},
	)
}

func TestStringToBigInt_(t *testing.T) {
	var o *big.Int
	assert.Nil(t, StringToBigInt("1", &o))
	assert.Equal(t, big.NewInt(1), o)

	assert.Nil(t, StringToBigInt("100000", &o))
	assert.Equal(t, big.NewInt(100_000), o)

	str := "1234567890123456789012"
	assert.Nil(t, StringToBigInt(str, &o))
	assert.Equal(t, str, o.String())

	assert.Equal(t, "The string value of 1.25 cannot be converted to *big.Int", StringToBigInt("1.25", &o).Error())

	funcs.TryTo(
		func() {
			MustStringToBigInt("1.25", &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The string value of 1.25 cannot be converted to *big.Int", e.(error).Error())
		},
	)
}
