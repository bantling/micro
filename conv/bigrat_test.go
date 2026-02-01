package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"math"
	"math/big"
	"testing"

	"github.com/bantling/micro/funcs"
	"github.com/stretchr/testify/assert"
)

func TestIntToBigRat_(t *testing.T) {
	var o *big.Rat
	IntToBigRat(int8(1), &o)
	assert.Equal(t, big.NewRat(1, 1), o)

	IntToBigRat(100_000, &o)
	assert.Equal(t, big.NewRat(100_000, 1), o)
}

func TestUintToBigRat_(t *testing.T) {
	var o *big.Rat
	UintToBigRat(uint8(1), &o)
	assert.Equal(t, big.NewRat(1, 1), o)

	UintToBigRat(uint(100_000), &o)
	assert.Equal(t, big.NewRat(100_000, 1), o)
}

func TestFloatToBigRat_(t *testing.T) {
	var o *big.Rat
	assert.Nil(t, FloatToBigRat(float32(1.25), &o))
	assert.Equal(t, big.NewRat(125, 100), o)

	assert.Nil(t, FloatToBigRat(float64(2.5), &o))
	assert.Equal(t, big.NewRat(25, 10), o)

	i := math.Inf(-1)
	assert.Equal(t, "The float64 value of -Inf cannot be converted to *big.Rat", FloatToBigRat(i, &o).Error())

	i = math.Inf(1)
	assert.Equal(t, "The float64 value of +Inf cannot be converted to *big.Rat", FloatToBigRat(i, &o).Error())

	i = math.NaN()
	assert.Equal(t, "The float64 value of NaN cannot be converted to *big.Rat", FloatToBigRat(i, &o).Error())

	funcs.TryTo(
		func() {
			MustFloatToBigRat(i, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The float64 value of NaN cannot be converted to *big.Rat", e.(error).Error())
		},
	)
}

func TestBigIntToBigRat_(t *testing.T) {
	var o *big.Rat
	BigIntToBigRat(big.NewInt(1), &o)
	assert.Equal(t, big.NewRat(1, 1), o)

	BigIntToBigRat(big.NewInt(100_000), &o)
	assert.Equal(t, big.NewRat(100_000, 1), o)
}

func TestBigFloatToBigRat_(t *testing.T) {
	var o *big.Rat
	assert.Nil(t, BigFloatToBigRat(big.NewFloat(1), &o))
	assert.Equal(t, big.NewRat(1, 1), o)

	assert.Nil(t, BigFloatToBigRat(big.NewFloat(1.25), &o))
	assert.Equal(t, big.NewRat(125, 100), o)

	i := big.NewFloat(math.Inf(1))
	assert.Equal(t, "The *big.Float value of +Inf cannot be converted to *big.Rat", BigFloatToBigRat(i, &o).Error())

	i = big.NewFloat(math.Inf(-1))
	assert.Equal(t, "The *big.Float value of -Inf cannot be converted to *big.Rat", BigFloatToBigRat(i, &o).Error())

	funcs.TryTo(
		func() {
			MustBigFloatToBigRat(i, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The *big.Float value of -Inf cannot be converted to *big.Rat", e.(error).Error())
		},
	)
}

func TestStringToBigRat_(t *testing.T) {
	var o *big.Rat
	assert.Nil(t, StringToBigRat("1/1", &o))
	assert.Equal(t, big.NewRat(1, 1), o)

	assert.Nil(t, StringToBigRat("125/100", &o))
	assert.Equal(t, big.NewRat(125, 100), o)

	assert.Nil(t, StringToBigRat("1.25", &o))
	assert.Equal(t, big.NewRat(125, 100), o)

	assert.Equal(t, "The string value of 1.25p cannot be converted to *big.Rat", StringToBigRat("1.25p", &o).Error())
	assert.Equal(t, "The string value of +Inf cannot be converted to *big.Rat", StringToBigRat("+Inf", &o).Error())
	assert.Equal(t, "The string value of -Inf cannot be converted to *big.Rat", StringToBigRat("-Inf", &o).Error())
	assert.Equal(t, "The string value of NaN cannot be converted to *big.Rat", StringToBigRat("NaN", &o).Error())

	funcs.TryTo(
		func() {
			MustStringToBigRat("NaN", &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The string value of NaN cannot be converted to *big.Rat", e.(error).Error())
		},
	)
}

func TestFloatStringToBigRat_(t *testing.T) {
	var o *big.Rat
	assert.Nil(t, FloatStringToBigRat("1.25", &o))
	assert.Equal(t, big.NewRat(125, 100), o)

	assert.Equal(t, "The float string value of 125/100 cannot be converted to *big.Rat", FloatStringToBigRat("125/100", &o).Error())
	assert.Equal(t, "The float string value of +Inf cannot be converted to *big.Rat", FloatStringToBigRat("+Inf", &o).Error())
	assert.Equal(t, "The float string value of -Inf cannot be converted to *big.Rat", FloatStringToBigRat("-Inf", &o).Error())
	assert.Equal(t, "The float string value of NaN cannot be converted to *big.Rat", FloatStringToBigRat("NaN", &o).Error())

	funcs.TryTo(
		func() {
			MustFloatStringToBigRat("NaN", &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The float string value of NaN cannot be converted to *big.Rat", e.(error).Error())
		},
	)
}
