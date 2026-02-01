package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/bantling/micro/funcs"
	"github.com/stretchr/testify/assert"
)

// ==== ToBigFloat

func TestIntToBigFloat_(t *testing.T) {
	var o *big.Float
	IntToBigFloat(int8(1), &o)
	cmp := big.NewFloat(1)
	cmp.SetPrec(uint(math.Ceil(1 * log2Of10)))
	assert.Equal(t, cmp, o)

	IntToBigFloat(100_000, &o)
	cmp = big.NewFloat(100_000)
	cmp.SetPrec(uint(math.Ceil(6 * log2Of10)))
	assert.Equal(t, cmp, o)
}

func TestUintToBigFloat_(t *testing.T) {
	var o *big.Float
	UintToBigFloat(uint8(1), &o)
	assert.Equal(t, big.NewFloat(1), o)

	UintToBigFloat(uint(100_000), &o)
	assert.Equal(t, big.NewFloat(100_000), o)
}

func TestFloatToBigFloat_(t *testing.T) {
	var o *big.Float
	assert.Nil(t, FloatToBigFloat(float32(1.25), &o))
	assert.Equal(t, big.NewFloat(1.25), o)

	assert.Nil(t, FloatToBigFloat(float64(100_000), &o))
	assert.Equal(t, big.NewFloat(100_000), o)

	assert.Nil(t, FloatToBigFloat(math.Inf(1), &o))
	assert.Equal(t, big.NewFloat(math.Inf(1)), o)

	assert.Nil(t, FloatToBigFloat(math.Inf(-1), &o))
	assert.Equal(t, big.NewFloat(math.Inf(-1)), o)

	assert.Equal(t, "The float64 value of NaN cannot be converted to *big.Float", FloatToBigFloat(math.NaN(), &o).Error())

	funcs.TryTo(
		func() {
			MustFloatToBigFloat(math.NaN(), &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The float64 value of NaN cannot be converted to *big.Float", e.(error).Error())
		},
	)
}

func TestBigIntToBigFloat_(t *testing.T) {
	var o *big.Float
	BigIntToBigFloat(big.NewInt(1), &o)
	assert.Equal(t, big.NewFloat(1), o)

	BigIntToBigFloat(big.NewInt(100_000), &o)
	assert.Equal(t, big.NewFloat(100_000), o)

	var (
		str   = "1234567890123456789012"
		inter *big.Int
	)
	StringToBigInt(str, &inter)
	BigIntToBigFloat(inter, &o)
	assert.Equal(t, str, fmt.Sprintf("%.f", o))
}

func TestBigRatToBigFloat_(t *testing.T) {
	var o *big.Float
	BigRatToBigFloat(big.NewRat(125, 100), &o)
	assert.Equal(t, big.NewFloat(1.25), o)

	BigRatToBigFloat(big.NewRat(25, 10), &o)
	assert.Equal(t, big.NewFloat(2.5), o)
}

func TestStringToBigFloat_(t *testing.T) {
	var o *big.Float
	assert.Nil(t, StringToBigFloat("1", &o))
	assert.Equal(t, big.NewFloat(1), o)

	assert.Nil(t, StringToBigFloat("100000.25", &o))
	assert.Equal(t, big.NewFloat(100000.25), o)

	str := "1234567890123456789012"
	assert.Nil(t, StringToBigFloat(str, &o))
	assert.Equal(t, str, fmt.Sprintf("%.f", o))

	assert.Nil(t, StringToBigFloat("+Inf", &o))
	assert.Equal(t, big.NewFloat(math.Inf(1)), o)

	assert.Nil(t, StringToBigFloat("-Inf", &o))
	assert.Equal(t, big.NewFloat(math.Inf(-1)), o)

	assert.Equal(t, "The string value of 1.25p cannot be converted to *big.Float", StringToBigFloat("1.25p", &o).Error())
	assert.Equal(t, "The string value of NaN cannot be converted to *big.Float", StringToBigFloat("NaN", &o).Error())

	funcs.TryTo(
		func() {
			MustStringToBigFloat("NaN", &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The string value of NaN cannot be converted to *big.Float", e.(error).Error())
		},
	)
}
