package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"math"
	"math/big"
	"testing"

	"github.com/bantling/micro/funcs"
	"github.com/stretchr/testify/assert"
)

func TestNumBits_(t *testing.T) {
	assert.Equal(t, 8, NumBits(int8(0)))
	assert.Equal(t, 8, NumBits(uint8(0)))

	assert.Equal(t, 16, NumBits(int16(0)))
	assert.Equal(t, 16, NumBits(uint16(0)))

	assert.Equal(t, 32, NumBits(int32(0)))
	assert.Equal(t, 32, NumBits(uint32(0)))

	assert.Equal(t, 64, NumBits(int64(0)))
	assert.Equal(t, 64, NumBits(uint64(0)))

	assert.Equal(t, 32, NumBits(float32(0)))
	assert.Equal(t, 64, NumBits(float64(0)))
}

// ==== ToInt

func TestIntToInt_(t *testing.T) {
	var d int8
	assert.Nil(t, IntToInt(1, &d))
	assert.Equal(t, int8(1), d)

	var s = math.MaxInt8
	assert.Nil(t, IntToInt(s, &d))
	assert.Equal(t, int8(math.MaxInt8), d)

	assert.Equal(t, "The int value of 32767 cannot be converted to int8", IntToInt(math.MaxInt16, &d).Error())

	MustIntToInt(2, &d)
	assert.Equal(t, int8(2), d)

	funcs.TryTo(
		func() {
			MustIntToInt(math.MaxInt16, &d)
			assert.Fail(t, "Never execute")
		},
		func(e any) { assert.Equal(t, "The int value of 32767 cannot be converted to int8", e.(error).Error()) },
	)
}

func TestUintToInt_(t *testing.T) {
	{
		var d int
		assert.Nil(t, UintToInt(uint(1), &d))
		assert.Equal(t, 1, d)
	}

	{
		var d int
		assert.Nil(t, UintToInt(uint64(math.MaxInt16), &d))
		assert.Equal(t, math.MaxInt16, d)
	}

	{
		var d int8
		assert.Equal(t, "The uint value of 32767 cannot be converted to int8", UintToInt(uint(math.MaxInt16), &d).Error())
	}

	{
		var d int8
		funcs.TryTo(
			func() {
				MustUintToInt(uint64(math.MaxUint16), &d)
				assert.Fail(t, "Never execute")
			},
			func(e any) {
				assert.Equal(t, "The uint64 value of 65535 cannot be converted to int8", e.(error).Error())
			},
		)
	}
}

func TestFloatToInt_(t *testing.T) {
	var d int
	assert.Nil(t, FloatToInt(float32(125), &d))
	assert.Equal(t, 125, d)

	assert.Equal(t, "The float64 value of 1.25 cannot be converted to int", FloatToInt(1.25, &d).Error())
	assert.Equal(t, "The float64 value of -Inf cannot be converted to int", FloatToInt(math.Inf(-1), &d).Error())
	assert.Equal(t, "The float64 value of +Inf cannot be converted to int", FloatToInt(math.Inf(1), &d).Error())
	assert.Equal(t, "The float64 value of NaN cannot be converted to int", FloatToInt(math.NaN(), &d).Error())

	funcs.TryTo(
		func() {
			MustFloatToInt(1.25, &d)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The float64 value of 1.25 cannot be converted to int", e.(error).Error())
		},
	)
}

// ==== ToUint

func TestIntToUint_(t *testing.T) {
	{
		var d uint
		assert.Nil(t, IntToUint(1, &d))
		assert.Equal(t, uint(1), d)
	}

	{
		var d uint16
		var s = math.MaxUint16
		assert.Nil(t, IntToUint(s, &d))
		assert.Equal(t, uint16(math.MaxUint16), d)
	}

	{
		var d uint8
		assert.Equal(t, "The int value of -1 cannot be converted to uint8", IntToUint(-1, &d).Error())
		assert.Equal(t, "The int value of 65535 cannot be converted to uint8", IntToUint(math.MaxUint16, &d).Error())
	}

	{
		var d uint8

		funcs.TryTo(
			func() {
				MustIntToUint(-1, &d)
				assert.Fail(t, "Never execute")
			},
			func(e any) { assert.Equal(t, "The int value of -1 cannot be converted to uint8", e.(error).Error()) },
		)
	}
}

func TestUintToUint_(t *testing.T) {
	{
		var d uint
		assert.Nil(t, UintToUint(uint(1), &d))
		assert.Equal(t, uint(1), d)
	}

	{
		var d uint32
		assert.Nil(t, UintToUint(uint64(math.MaxUint16), &d))
		assert.Equal(t, uint32(math.MaxUint16), d)
	}

	{
		var d uint8
		assert.Equal(t, "The uint value of 32767 cannot be converted to uint8", UintToUint(uint(math.MaxInt16), &d).Error())
	}

	{
		var d uint8
		funcs.TryTo(
			func() {
				MustUintToUint(uint64(math.MaxUint16), &d)
				assert.Fail(t, "Never execute")
			},
			func(e any) {
				assert.Equal(t, "The uint64 value of 65535 cannot be converted to uint8", e.(error).Error())
			},
		)
	}
}

func TestFloatToUint_(t *testing.T) {
	var d uint
	assert.Nil(t, FloatToUint(float32(125), &d))
	assert.Equal(t, uint(125), d)

	assert.Equal(t, "The float64 value of 1.25 cannot be converted to uint", FloatToUint(1.25, &d).Error())
	assert.Equal(t, "The float64 value of -Inf cannot be converted to uint", FloatToUint(math.Inf(-1), &d).Error())
	assert.Equal(t, "The float64 value of +Inf cannot be converted to uint", FloatToUint(math.Inf(1), &d).Error())
	assert.Equal(t, "The float64 value of NaN cannot be converted to uint", FloatToUint(math.NaN(), &d).Error())

	funcs.TryTo(
		func() {
			MustFloatToUint(1.25, &d)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The float64 value of 1.25 cannot be converted to uint", e.(error).Error())
		},
	)
}

// ==== ToInt64

func TestBigIntToInt64_(t *testing.T) {
	var o int64
	assert.Nil(t, BigIntToInt64(big.NewInt(1), &o))
	assert.Equal(t, int64(1), o)

	assert.Nil(t, BigIntToInt64(big.NewInt(100_000), &o))
	assert.Equal(t, int64(100_000), o)

	var (
		str   = "123456789012345678901"
		inter *big.Int
	)
	StringToBigInt(str, &inter)
	assert.Equal(t, "The *big.Int value of 123456789012345678901 cannot be converted to int64", BigIntToInt64(inter, &o).Error())

	funcs.TryTo(
		func() {
			MustBigIntToInt64(inter, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The *big.Int value of 123456789012345678901 cannot be converted to int64", e.(error).Error())
		},
	)
}

func TestBigFloatToInt64_(t *testing.T) {
	var o int64
	assert.Nil(t, BigFloatToInt64(big.NewFloat(1), &o))
	assert.Equal(t, int64(1), o)

	assert.Nil(t, BigFloatToInt64(big.NewFloat(100_000), &o))
	assert.Equal(t, int64(100_000), o)

	assert.Equal(t, "The *big.Float value of 1.25 cannot be converted to int64", BigFloatToInt64(big.NewFloat(1.25), &o).Error())

	negInf := big.NewFloat(0).Quo(big.NewFloat(-1), big.NewFloat(0))
	assert.Equal(t, "The *big.Float value of -Inf cannot be converted to int64", BigFloatToInt64(negInf, &o).Error())

	posInf := big.NewFloat(0).Quo(big.NewFloat(1), big.NewFloat(0))
	assert.Equal(t, "The *big.Float value of +Inf cannot be converted to int64", BigFloatToInt64(posInf, &o).Error())

	funcs.TryTo(
		func() {
			MustBigFloatToInt64(posInf, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The *big.Float value of +Inf cannot be converted to int64", e.(error).Error())
		},
	)
}

func TestBigRatToInt64_(t *testing.T) {
	var o int64
	assert.Nil(t, BigRatToInt64(big.NewRat(1, 1), &o))
	assert.Equal(t, int64(1), o)

	assert.Nil(t, BigRatToInt64(big.NewRat(100_000, 1), &o))
	assert.Equal(t, int64(100_000), o)

	assert.Equal(t, "The *big.Rat value of 5/4 cannot be converted to int64", BigRatToInt64(big.NewRat(125, 100), &o).Error())

	funcs.TryTo(
		func() {
			MustBigRatToInt64(big.NewRat(125, 100), &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The *big.Rat value of 5/4 cannot be converted to int64", e.(error).Error())
		},
	)
}

func TestStringToInt64_(t *testing.T) {
	var o int64
	assert.Nil(t, StringToInt64("1", &o))
	assert.Equal(t, int64(1), o)

	assert.Nil(t, StringToInt64("100000", &o))
	assert.Equal(t, int64(100_000), o)

	str := "123456789012345678901"
	assert.Equal(t, "The string value of 123456789012345678901 cannot be converted to int64", StringToInt64(str, &o).Error())

	funcs.TryTo(
		func() {
			MustStringToInt64(str, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The string value of 123456789012345678901 cannot be converted to int64", e.(error).Error())
		},
	)
}

// ==== ToUint64

func TestBigIntToUint64_(t *testing.T) {
	var o uint64
	assert.Nil(t, BigIntToUint64(big.NewInt(1), &o))
	assert.Equal(t, uint64(1), o)

	assert.Nil(t, BigIntToUint64(big.NewInt(100_000), &o))
	assert.Equal(t, uint64(100_000), o)

	var inter *big.Int
	StringToBigInt("123456789012345678901", &inter)
	assert.Equal(t, "The *big.Int value of 123456789012345678901 cannot be converted to uint64", BigIntToUint64(inter, &o).Error())

	funcs.TryTo(
		func() {
			MustBigIntToUint64(inter, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The *big.Int value of 123456789012345678901 cannot be converted to uint64", e.(error).Error())
		},
	)
}

func TestBigFloatToUint64_(t *testing.T) {
	var o uint64
	assert.Nil(t, BigFloatToUint64(big.NewFloat(1), &o))
	assert.Equal(t, uint64(1), o)

	assert.Nil(t, BigFloatToUint64(big.NewFloat(100_000), &o))
	assert.Equal(t, uint64(100_000), o)

	assert.Equal(t, "The *big.Float value of 1.25 cannot be converted to uint64", BigFloatToUint64(big.NewFloat(1.25), &o).Error())
	assert.Equal(t, "The *big.Float value of -Inf cannot be converted to uint64", BigFloatToUint64(big.NewFloat(math.Inf(-1)), &o).Error())
	assert.Equal(t, "The *big.Float value of +Inf cannot be converted to uint64", BigFloatToUint64(big.NewFloat(math.Inf(1)), &o).Error())

	funcs.TryTo(
		func() {
			MustBigFloatToUint64(big.NewFloat(math.Inf(1)), &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The *big.Float value of +Inf cannot be converted to uint64", e.(error).Error())
		},
	)
}

func TestBigRatToUint64_(t *testing.T) {
	var o uint64
	assert.Nil(t, BigRatToUint64(big.NewRat(1, 1), &o))
	assert.Equal(t, uint64(1), o)

	assert.Nil(t, BigRatToUint64(big.NewRat(100_000, 1), &o))
	assert.Equal(t, uint64(100_000), o)

	assert.Equal(t, "The *big.Rat value of 5/4 cannot be converted to uint64", BigRatToUint64(big.NewRat(125, 100), &o).Error())

	funcs.TryTo(
		func() {
			MustBigRatToUint64(big.NewRat(125, 100), &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The *big.Rat value of 5/4 cannot be converted to uint64", e.(error).Error())
		},
	)
}

func TestStringToUint64_(t *testing.T) {
	var o uint64
	assert.Nil(t, StringToUint64("1", &o))
	assert.Equal(t, uint64(1), o)

	assert.Nil(t, StringToUint64("100000", &o))
	assert.Equal(t, uint64(100_000), o)

	str := "123456789012345678901"
	assert.Equal(t, "The string value of 123456789012345678901 cannot be converted to uint64", StringToUint64(str, &o).Error())

	funcs.TryTo(
		func() {
			MustStringToUint64(str, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The string value of 123456789012345678901 cannot be converted to uint64", e.(error).Error())
		},
	)
}
