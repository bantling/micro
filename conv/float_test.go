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

// ==== ToFloat

func TestIntToFloat_(t *testing.T) {
	// The following code tries int values that start at the maximum value a float32 can hold, and continue for 8 values after
	// i := 0xFFFFFF
	// bits := fmt.Sprintf("%b", i)
	// fmt.Printf("%d, %s, %d\n", i, bits, len(bits))
	//
	// for j := 0; j <= 8; j++ {
	//   v := i + j
	//   fmt.Printf("%d, %b, %.f, %t\n", v, v, float32(v), v == int(float32(v)))
	// }
	//
	// 16777215, 111111111111111111111111, 24
	// 16777215, 111111111111111111111111, 16777215, true
	// 16777216, 1000000000000000000000000, 16777216, true
	// 16777217, 1000000000000000000000001, 16777216, false
	// 16777218, 1000000000000000000000010, 16777218, true
	// 16777219, 1000000000000000000000011, 16777220, false
	// 16777220, 1000000000000000000000100, 16777220, true
	// 16777221, 1000000000000000000000101, 16777220, false
	// 16777222, 1000000000000000000000110, 16777222, true
	// 16777223, 1000000000000000000000111, 16777224, false
	//
	// Notice that the initial value is 24 bits, not 23 as stated maximum number of bits for float32 mantissa.
	// That's because IEE 754 floating point numbers have one implicit bit of precision that is not stored.

	{
		goodCases := []int32{
			16777215,
			16777216,
			16777218,
			16777220,
			16777222,
		}

		badCases := []int32{
			16777217,
			16777219,
			16777221,
			16777223,
		}

		var fval float32
		for _, ival := range goodCases {
			assert.Nil(t, IntToFloat(ival, &fval))
			assert.Equal(t, ival, int32(fval))
		}

		for _, ival := range badCases {
			assert.Equal(t, fmt.Sprintf("The int32 value of %d cannot be converted to float32", ival), IntToFloat(ival, &fval).Error())
		}

		{
			var d float32
			funcs.TryTo(
				func() {
					MustIntToFloat(math.MaxUint32, &d)
					assert.Fail(t, "Never execute")
				},
				func(e any) {
					assert.Equal(t, "The int value of 4294967295 cannot be converted to float32", e.(error).Error())
				},
			)
		}
	}

	// The following code tries int values that start at maximum value a float64 can hold, and continue for 8 values after
	// i := 0x1FFFFFFFFFFFFF
	// bits := fmt.Sprintf("%b", i)
	// fmt.Printf("%d, %s, %d\n", i, bits, len(bits))
	//
	// for j := 0; j <= 8; j++ {
	// 	v := i + j
	// 	fmt.Printf("%d, %b, %.f, %t\n", v, v, float64(v), v == int(float64(v)))
	// }
	//
	// 9007199254740991, 11111111111111111111111111111111111111111111111111111, 53
	// 9007199254740991, 11111111111111111111111111111111111111111111111111111, 9007199254740991, true
	// 9007199254740992, 100000000000000000000000000000000000000000000000000000, 9007199254740992, true
	// 9007199254740993, 100000000000000000000000000000000000000000000000000001, 9007199254740992, false
	// 9007199254740994, 100000000000000000000000000000000000000000000000000010, 9007199254740994, true
	// 9007199254740995, 100000000000000000000000000000000000000000000000000011, 9007199254740996, false
	// 9007199254740996, 100000000000000000000000000000000000000000000000000100, 9007199254740996, true
	// 9007199254740997, 100000000000000000000000000000000000000000000000000101, 9007199254740996, false
	// 9007199254740998, 100000000000000000000000000000000000000000000000000110, 9007199254740998, true
	// 9007199254740999, 100000000000000000000000000000000000000000000000000111, 9007199254741000, false
	//
	// Notice that the initial value is 53 bits, not 52 as stated maximum number of bits for float64 mantissa.
	// That's because IEE 754 floating point numbers have one implicit bit of precision that is not stored.

	{
		goodCases := []int64{
			9007199254740991,
			9007199254740992,
			9007199254740994,
			9007199254740996,
			9007199254740998,
		}

		badCases := []int64{
			9007199254740993,
			9007199254740995,
			9007199254740997,
			9007199254740999,
		}

		var fval float64
		for _, ival := range goodCases {
			assert.Nil(t, IntToFloat(ival, &fval))
			assert.Equal(t, ival, int64(fval))
		}

		for _, ival := range badCases {
			assert.Equal(t, fmt.Sprintf("The int64 value of %d cannot be converted to float64", ival), IntToFloat(ival, &fval).Error())
		}

		{
			var d float64
			funcs.TryTo(
				func() {
					MustIntToFloat(int64((1<<60)+1), &d)
					assert.Fail(t, "Never execute")
				},
				func(e any) {
					assert.Equal(t, "The int64 value of 1152921504606846977 cannot be converted to float64", e.(error).Error())
				},
			)
		}
	}
}

func TestFloatToFloat_(t *testing.T) {
	var (
		i float64
		o float32
	)
	assert.Nil(t, FloatToFloat(i, &o))
	assert.Equal(t, float32(0), o)

	i = 1
	assert.Nil(t, FloatToFloat(i, &o))
	assert.Equal(t, float32(1), o)

	i = -1
	assert.Nil(t, FloatToFloat(i, &o))
	assert.Equal(t, float32(-1), o)

	assert.Nil(t, FloatToFloat(math.SmallestNonzeroFloat32, &o))
	assert.Equal(t, float32(math.SmallestNonzeroFloat32), o)

	assert.Nil(t, FloatToFloat(-math.SmallestNonzeroFloat32, &o))
	assert.Equal(t, float32(-math.SmallestNonzeroFloat32), o)

	assert.Nil(t, FloatToFloat(math.MaxFloat32, &o))
	assert.Equal(t, float32(math.MaxFloat32), o)

	assert.Nil(t, FloatToFloat(math.Inf(-1), &o))
	assert.Equal(t, float32(math.Inf(-1)), o)

	assert.Nil(t, FloatToFloat(math.Inf(1), &o))
	assert.Equal(t, float32(math.Inf(1)), o)

	assert.Nil(t, FloatToFloat(math.NaN(), &o))
	assert.True(t, math.IsNaN(float64(o)))

	i = math.SmallestNonzeroFloat64
	assert.Equal(t, "The float64 value of 5e-324 cannot be converted to float32", FloatToFloat(i, &o).Error())

	i = -math.SmallestNonzeroFloat64
	assert.Equal(t, "The float64 value of -5e-324 cannot be converted to float32", FloatToFloat(i, &o).Error())

	funcs.TryTo(
		func() {
			MustFloatToFloat(i, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The float64 value of -5e-324 cannot be converted to float32", e.(error).Error())
		},
	)
}

// ==== ToFloat32

func TestBigIntToFloat32_(t *testing.T) {
	var o float32
	assert.Nil(t, BigIntToFloat32(big.NewInt(1), &o))
	assert.Equal(t, float32(1), o)

	assert.Nil(t, BigIntToFloat32(big.NewInt(100_000), &o))
	assert.Equal(t, float32(100_000), o)

	var (
		str   = "123456789012345678901"
		inter *big.Int
	)
	StringToBigInt(str, &inter)
	assert.Equal(t, "The *big.Int value of 123456789012345678901 cannot be converted to float32", BigIntToFloat32(inter, &o).Error())

	funcs.TryTo(
		func() {
			MustBigIntToFloat32(inter, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The *big.Int value of 123456789012345678901 cannot be converted to float32", e.(error).Error())
		},
	)
}

func TestBigFloatToFloat32_(t *testing.T) {
	var o float32
	assert.Nil(t, BigFloatToFloat32(big.NewFloat(1), &o))
	assert.Equal(t, float32(1), o)

	assert.Nil(t, BigFloatToFloat32(big.NewFloat(100_000), &o))
	assert.Equal(t, float32(100_000), o)

	var (
		str   = "123456789012345678901"
		inter *big.Float
	)
	StringToBigFloat(str, &inter)
	assert.Equal(t, "The *big.Float value of 123456789012345678901 cannot be converted to float32", BigFloatToFloat32(inter, &o).Error())

	funcs.TryTo(
		func() {
			MustBigFloatToFloat32(inter, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The *big.Float value of 123456789012345678901 cannot be converted to float32", e.(error).Error())
		},
	)
}

func TestBigRatToFloat32_(t *testing.T) {
	var o float32
	assert.Nil(t, BigRatToFloat32(big.NewRat(1, 1), &o))
	assert.Equal(t, float32(1), o)

	assert.Nil(t, BigRatToFloat32(big.NewRat(125, 100), &o))
	assert.Equal(t, float32(1.25), o)

	i := big.NewRat(math.MaxInt64, 100)
	assert.Equal(t, "The *big.Rat value of 9223372036854775807/100 cannot be converted to float32", BigRatToFloat32(i, &o).Error())

	funcs.TryTo(
		func() {
			MustBigRatToFloat32(i, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The *big.Rat value of 9223372036854775807/100 cannot be converted to float32", e.(error).Error())
		},
	)
}

func TestStringToFloat32_(t *testing.T) {
	var o float32
	assert.Nil(t, StringToFloat32("1", &o))
	assert.Equal(t, float32(1), o)

	assert.Nil(t, StringToFloat32("1.25", &o))
	assert.Equal(t, float32(1.25), o)

	assert.Nil(t, StringToFloat32("-Inf", &o))
	assert.Equal(t, float32(math.Inf(-1)), o)

	assert.Nil(t, StringToFloat32("+Inf", &o))
	assert.Equal(t, float32(math.Inf(1)), o)

	assert.Nil(t, StringToFloat32("NaN", &o))
	assert.True(t, math.IsNaN(float64(o)))

	str := "123456789012345678901"
	assert.Equal(t, "The string value of 123456789012345678901 cannot be converted to float32", StringToFloat32(str, &o).Error())

	funcs.TryTo(
		func() {
			MustStringToFloat32(str, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The string value of 123456789012345678901 cannot be converted to float32", e.(error).Error())
		},
	)
}

// ==== ToFloat64

func TestBigIntToFloat64_(t *testing.T) {
	var o float64
	assert.Nil(t, BigIntToFloat64(big.NewInt(1), &o))
	assert.Equal(t, float64(1), o)

	assert.Nil(t, BigIntToFloat64(big.NewInt(100_000), &o))
	assert.Equal(t, float64(100_000), o)

	var (
		str   = "123456789012345678901"
		inter *big.Int
	)
	StringToBigInt(str, &inter)
	assert.Equal(t, "The *big.Int value of 123456789012345678901 cannot be converted to float64", BigIntToFloat64(inter, &o).Error())

	funcs.TryTo(
		func() {
			MustBigIntToFloat64(inter, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The *big.Int value of 123456789012345678901 cannot be converted to float64", e.(error).Error())
		},
	)
}

func TestBigFloatToFloat64_(t *testing.T) {
	var o float64
	assert.Nil(t, BigFloatToFloat64(big.NewFloat(1), &o))
	assert.Equal(t, float64(1), o)
	assert.Nil(t, BigFloatToFloat64(big.NewFloat(100_000), &o))
	assert.Equal(t, float64(100_000), o)

	var (
		str   = "123456789012345678901"
		inter *big.Float
	)
	StringToBigFloat(str, &inter)
	assert.Equal(t, "The *big.Float value of 123456789012345678901 cannot be converted to float64", BigFloatToFloat64(inter, &o).Error())

	funcs.TryTo(
		func() {
			MustBigFloatToFloat64(inter, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The *big.Float value of 123456789012345678901 cannot be converted to float64", e.(error).Error())
		},
	)
}

func TestBigRatToFloat64_(t *testing.T) {
	var o float64
	assert.Nil(t, BigRatToFloat64(big.NewRat(1, 1), &o))
	assert.Equal(t, float64(1), o)

	assert.Nil(t, BigRatToFloat64(big.NewRat(125, 100), &o))
	assert.Equal(t, float64(1.25), o)

	i := big.NewRat(math.MaxInt64, 100)
	assert.Equal(t, "The *big.Rat value of 9223372036854775807/100 cannot be converted to float64", BigRatToFloat64(i, &o).Error())

	funcs.TryTo(
		func() {
			MustBigRatToFloat64(i, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The *big.Rat value of 9223372036854775807/100 cannot be converted to float64", e.(error).Error())
		},
	)
}

func TestStringToFloat64_(t *testing.T) {
	var o float64
	assert.Nil(t, StringToFloat64("1", &o))
	assert.Equal(t, float64(1), o)

	assert.Nil(t, StringToFloat64("1.25", &o))
	assert.Equal(t, float64(1.25), o)

	assert.Nil(t, StringToFloat64("-Inf", &o))
	assert.Equal(t, math.Inf(-1), o)

	assert.Nil(t, StringToFloat64("+Inf", &o))
	assert.Equal(t, math.Inf(1), o)

	assert.Nil(t, StringToFloat64("NaN", &o))
	assert.True(t, math.IsNaN(o))

	str := "123456789012345678901"
	assert.Equal(t, "The string value of 123456789012345678901 cannot be converted to float64", StringToFloat64(str, &o).Error())

	funcs.TryTo(
		func() {
			MustStringToFloat64(str, &o)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, "The string value of 123456789012345678901 cannot be converted to float64", e.(error).Error())
		},
	)
}
