package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==== ToString

func TestIntToString_(t *testing.T) {
	assert.Equal(t, IntToString(int8(1)), "1")
	assert.Equal(t, IntToString(int16(2)), "2")
	assert.Equal(t, IntToString(int32(3)), "3")
	assert.Equal(t, IntToString(int64(4)), "4")
	assert.Equal(t, IntToString(int(5)), "5")
}

func TestUintToString_(t *testing.T) {
	assert.Equal(t, UintToString(uint8(1)), "1")
	assert.Equal(t, UintToString(uint16(2)), "2")
	assert.Equal(t, UintToString(uint32(3)), "3")
	assert.Equal(t, UintToString(uint64(4)), "4")
	assert.Equal(t, UintToString(uint(5)), "5")
}

func TestFloatToString_(t *testing.T) {
	assert.Equal(t, "1.25", FloatToString(float32(1.25)))
	assert.Equal(t, "1.25", FloatToString(float64(1.25)))
	assert.Equal(t, "-Inf", FloatToString(float32(math.Inf(-1))))
	assert.Equal(t, "+Inf", FloatToString(math.Inf(1)))
	assert.Equal(t, "NaN", FloatToString(math.NaN()))
	assert.Equal(t, "-0", FloatToString(-1/math.Inf(1)))
}

func TestBigIntToString_(t *testing.T) {
	assert.Equal(t, "1234", BigIntToString(big.NewInt(1234)))
}

func TestBigFloatToString_(t *testing.T) {
	assert.Equal(t, "1234.5678", BigFloatToString(big.NewFloat(1234.5678)))
	assert.Equal(t, "-Inf", BigFloatToString(big.NewFloat(math.Inf(-1))))
	assert.Equal(t, "+Inf", BigFloatToString(big.NewFloat(math.Inf(1))))
	assert.Equal(t, "-0", BigFloatToString(big.NewFloat(-1/math.Inf(1))))
}

func TestBigRatToString_(t *testing.T) {
	assert.Equal(t, "5/4", BigRatToString(big.NewRat(125, 100)))
}

func TestBigRatToNormalizedString_(t *testing.T) {
	assert.Equal(t, "1234", BigRatToNormalizedString(big.NewRat(1234, 1)))
	assert.Equal(t, "1.25", BigRatToNormalizedString(big.NewRat(125, 100)))
}

func TestToString_(t *testing.T) {
	assert.Equal(t, "1", ToString(int(1)))
	assert.Equal(t, "2", ToString(int8(2)))
	assert.Equal(t, "3", ToString(int16(3)))
	assert.Equal(t, "4", ToString(int32(4)))
	assert.Equal(t, "5", ToString(int64(5)))

	assert.Equal(t, "1", ToString(uint(1)))
	assert.Equal(t, "2", ToString(uint8(2)))
	assert.Equal(t, "3", ToString(uint16(3)))
	assert.Equal(t, "4", ToString(uint32(4)))
	assert.Equal(t, "5", ToString(uint64(5)))

	assert.Equal(t, "1.25", ToString(float32(1.25)))
	assert.Equal(t, "2.75", ToString(float64(2.75)))

	assert.Equal(t, "1", ToString(big.NewInt(1)))
	assert.Equal(t, "1.25", ToString(big.NewFloat(1.25)))
	assert.Equal(t, "2.75", ToString(big.NewRat(275, 100)))
}

// ==== int/uint to int/uint, float to int, float64 to float32

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

func TestIntToInt_(t *testing.T) {
	var d int8
	assert.Nil(t, IntToInt(1, &d))
	assert.Equal(t, int8(1), d)

	var s = math.MaxInt8
	assert.Nil(t, IntToInt(s, &d))
	assert.Equal(t, int8(math.MaxInt8), d)

	assert.Equal(t, "The int value of 32767 cannot be converted to int8", IntToInt(math.MaxInt16, &d).Error())
}

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
}

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
}

func TestFloatToUint_(t *testing.T) {
	var d uint
	assert.Nil(t, FloatToUint(float32(125), &d))
	assert.Equal(t, uint(125), d)

	assert.Equal(t, "The float64 value of 1.25 cannot be converted to uint", FloatToUint(1.25, &d).Error())
	assert.Equal(t, "The float64 value of -Inf cannot be converted to uint", FloatToUint(math.Inf(-1), &d).Error())
	assert.Equal(t, "The float64 value of +Inf cannot be converted to uint", FloatToUint(math.Inf(1), &d).Error())
	assert.Equal(t, "The float64 value of NaN cannot be converted to uint", FloatToUint(math.NaN(), &d).Error())
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
}

func TestBigRatToInt64_(t *testing.T) {
	var o int64
	assert.Nil(t, BigRatToInt64(big.NewRat(1, 1), &o))
	assert.Equal(t, int64(1), o)

	assert.Nil(t, BigRatToInt64(big.NewRat(100_000, 1), &o))
	assert.Equal(t, int64(100_000), o)

	assert.Equal(t, "The *big.Rat value of 5/4 cannot be converted to int64", BigRatToInt64(big.NewRat(125, 100), &o).Error())
}

func TestStringToInt64_(t *testing.T) {
	var o int64
	assert.Nil(t, StringToInt64("1", &o))
	assert.Equal(t, int64(1), o)

	assert.Nil(t, StringToInt64("100000", &o))
	assert.Equal(t, int64(100_000), o)

	str := "123456789012345678901"
	assert.Equal(t, "The string value of 123456789012345678901 cannot be converted to int64", StringToInt64(str, &o).Error())
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
}

func TestBigRatToUint64_(t *testing.T) {
	var o uint64
	assert.Nil(t, BigRatToUint64(big.NewRat(1, 1), &o))
	assert.Equal(t, uint64(1), o)

	assert.Nil(t, BigRatToUint64(big.NewRat(100_000, 1), &o))
	assert.Equal(t, uint64(100_000), o)

	assert.Equal(t, "The *big.Rat value of 5/4 cannot be converted to uint64", BigRatToUint64(big.NewRat(125, 100), &o).Error())
}

func TestStringToUint64_(t *testing.T) {
	var o uint64
	assert.Nil(t, StringToUint64("1", &o))
	assert.Equal(t, uint64(1), o)

	assert.Nil(t, StringToUint64("100000", &o))
	assert.Equal(t, uint64(100_000), o)

	str := "123456789012345678901"
	assert.Equal(t, "The string value of 123456789012345678901 cannot be converted to uint64", StringToUint64(str, &o).Error())
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
}

func TestBigRatToFloat32_(t *testing.T) {
	var o float32
	assert.Nil(t, BigRatToFloat32(big.NewRat(1, 1), &o))
	assert.Equal(t, float32(1), o)

	assert.Nil(t, BigRatToFloat32(big.NewRat(125, 100), &o))
	assert.Equal(t, float32(1.25), o)

	i := big.NewRat(math.MaxInt64, 100)
	assert.Equal(t, "The *big.Rat value of 9223372036854775807/100 cannot be converted to float32", BigRatToFloat32(i, &o).Error())
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
}

func TestBigRatToFloat64_(t *testing.T) {
	var o float64
	assert.Nil(t, BigRatToFloat64(big.NewRat(1, 1), &o))
	assert.Equal(t, float64(1), o)

	assert.Nil(t, BigRatToFloat64(big.NewRat(125, 100), &o))
	assert.Equal(t, float64(1.25), o)

	i := big.NewRat(math.MaxInt64, 100)
	assert.Equal(t, "The *big.Rat value of 9223372036854775807/100 cannot be converted to float64", BigRatToFloat64(i, &o).Error())
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
}

// ==== ToBigInt

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
}

func TestBigRatToBigInt_(t *testing.T) {
	var o *big.Int
	assert.Nil(t, BigRatToBigInt(big.NewRat(1, 1), &o))
	assert.Equal(t, big.NewInt(1), o)

	assert.Nil(t, BigRatToBigInt(big.NewRat(100_000, 1), &o))
	assert.Equal(t, big.NewInt(100_000), o)

	assert.Equal(t, "The *big.Rat value of 5/4 cannot be converted to *big.Int", BigRatToBigInt(big.NewRat(125, 100), &o).Error())
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
}

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
}

// ==== ToBigRat

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
}

func TestFloatStringToBigRat_(t *testing.T) {
	var o *big.Rat
	assert.Nil(t, FloatStringToBigRat("1.25", &o))
	assert.Equal(t, big.NewRat(125, 100), o)

	assert.Equal(t, "The float string value of 125/100 cannot be converted to *big.Rat", FloatStringToBigRat("125/100", &o).Error())
	assert.Equal(t, "The float string value of +Inf cannot be converted to *big.Rat", FloatStringToBigRat("+Inf", &o).Error())
	assert.Equal(t, "The float string value of -Inf cannot be converted to *big.Rat", FloatStringToBigRat("-Inf", &o).Error())
	assert.Equal(t, "The float string value of NaN cannot be converted to *big.Rat", FloatStringToBigRat("NaN", &o).Error())
}

func TestTo_(t *testing.T) {
	// == int
	{
		var i int

		// ints
		assert.Nil(t, To(-1, &i))
		assert.Equal(t, -1, i)

		assert.Nil(t, To(int8(-2), &i))
		assert.Equal(t, -2, i)

		assert.Nil(t, To(int16(-3), &i))
		assert.Equal(t, -3, i)

		assert.Nil(t, To(int32(-4), &i))
		assert.Equal(t, -4, i)

		assert.Nil(t, To(int64(-5), &i))
		assert.Equal(t, -5, i)

		// uints
		assert.Nil(t, To(uint(1), &i))
		assert.Equal(t, 1, i)

		assert.Nil(t, To(uint8(2), &i))
		assert.Equal(t, 2, i)

		assert.Nil(t, To(uint16(3), &i))
		assert.Equal(t, 3, i)

		assert.Nil(t, To(uint32(4), &i))
		assert.Equal(t, 4, i)

		assert.Nil(t, To(uint64(5), &i))
		assert.Equal(t, 5, i)

		// floats
		assert.Nil(t, To(float32(1), &i))
		assert.Equal(t, 1, i)

		assert.Nil(t, To(2.0, &i))
		assert.Equal(t, 2, i)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &i))
		assert.Equal(t, 1, i)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(2))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 18446744073709551614 cannot be converted to int64"), To(bi, &i))

		assert.Nil(t, To(big.NewFloat(2), &i))
		assert.Equal(t, 2, i)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to int64"), To(big.NewFloat(1.25), &i))

		assert.Nil(t, To(big.NewRat(3, 1), &i))
		assert.Equal(t, 3, i)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to int64"), To(big.NewRat(5, 4), &i))

		// string
		assert.Nil(t, To("1", &i))
		assert.Equal(t, 1, i)

    assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "int64"), To("a", &i))
		assert.Equal(t, 1, i)
	}

	// == int8
	{
		var i8 int8

		// ints
		assert.Nil(t, To(-1, &i8))
		assert.Equal(t, int8(-1), i8)

		assert.Nil(t, To(int8(-2), &i8))
		assert.Equal(t, int8(-2), i8)

		assert.Nil(t, To(int16(-3), &i8))
		assert.Equal(t, int8(-3), i8)

		assert.Nil(t, To(int32(-4), &i8))
		assert.Equal(t, int8(-4), i8)

		assert.Nil(t, To(int64(-5), &i8))
		assert.Equal(t, int8(-5), i8)

		// uints
		assert.Nil(t, To(uint(1), &i8))
		assert.Equal(t, int8(1), i8)

		assert.Nil(t, To(uint8(2), &i8))
		assert.Equal(t, int8(2), i8)

		assert.Nil(t, To(uint16(3), &i8))
		assert.Equal(t, int8(3), i8)

		assert.Nil(t, To(uint32(4), &i8))
		assert.Equal(t, int8(4), i8)

		assert.Nil(t, To(uint64(5), &i8))
		assert.Equal(t, int8(5), i8)

		// floats
		assert.Nil(t, To(float32(1), &i8))
		assert.Equal(t, int8(1), i8)

		assert.Nil(t, To(2.0, &i8))
		assert.Equal(t, int8(2), i8)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &i8))
		assert.Equal(t, int8(1), i8)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(2))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 18446744073709551614 cannot be converted to int64"), To(bi, &i8))

		assert.Nil(t, To(big.NewFloat(2), &i8))
		assert.Equal(t, int8(2), i8)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to int64"), To(big.NewFloat(1.25), &i8))

		assert.Nil(t, To(big.NewRat(3, 1), &i8))
		assert.Equal(t, int8(3), i8)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to int64"), To(big.NewRat(5, 4), &i8))

		// string
		assert.Nil(t, To("1", &i8))
		assert.Equal(t, int8(1), i8)

    assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "int64"), To("a", &i8))
		assert.Equal(t, int8(1), i8)
	}

	// == int16
	{
		var i16 int16

		// ints
		assert.Nil(t, To(-1, &i16))
		assert.Equal(t, int16(-1), i16)

		assert.Nil(t, To(int8(-2), &i16))
		assert.Equal(t, int16(-2), i16)

		assert.Nil(t, To(int16(-3), &i16))
		assert.Equal(t, int16(-3), i16)

		assert.Nil(t, To(int32(-4), &i16))
		assert.Equal(t, int16(-4), i16)

		assert.Nil(t, To(int64(-5), &i16))
		assert.Equal(t, int16(-5), i16)

		// uints
		assert.Nil(t, To(uint(1), &i16))
		assert.Equal(t, int16(1), i16)

		assert.Nil(t, To(uint8(2), &i16))
		assert.Equal(t, int16(2), i16)

		assert.Nil(t, To(uint16(3), &i16))
		assert.Equal(t, int16(3), i16)

		assert.Nil(t, To(uint32(4), &i16))
		assert.Equal(t, int16(4), i16)

		assert.Nil(t, To(uint64(5), &i16))
		assert.Equal(t, int16(5), i16)

		// floats
		assert.Nil(t, To(float32(1), &i16))
		assert.Equal(t, int16(1), i16)

		assert.Nil(t, To(2.0, &i16))
		assert.Equal(t, int16(2), i16)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &i16))
		assert.Equal(t, int16(1), i16)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(2))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 18446744073709551614 cannot be converted to int64"), To(bi, &i16))

		assert.Nil(t, To(big.NewFloat(2), &i16))
		assert.Equal(t, int16(2), i16)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to int64"), To(big.NewFloat(1.25), &i16))

		assert.Nil(t, To(big.NewRat(3, 1), &i16))
		assert.Equal(t, int16(3), i16)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to int64"), To(big.NewRat(5, 4), &i16))

		// string
		assert.Nil(t, To("1", &i16))
		assert.Equal(t, int16(1), i16)

    assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "int64"), To("a", &i16))
		assert.Equal(t, int16(1), i16)
	}

	// == int32
	{
		var i32 int32

		// ints
		assert.Nil(t, To(-1, &i32))
		assert.Equal(t, int32(-1), i32)

		assert.Nil(t, To(int8(-2), &i32))
		assert.Equal(t, int32(-2), i32)

		assert.Nil(t, To(int16(-3), &i32))
		assert.Equal(t, int32(-3), i32)

		assert.Nil(t, To(int32(-4), &i32))
		assert.Equal(t, int32(-4), i32)

		assert.Nil(t, To(int64(-5), &i32))
		assert.Equal(t, int32(-5), i32)

		// uints
		assert.Nil(t, To(uint(1), &i32))
		assert.Equal(t, int32(1), i32)

		assert.Nil(t, To(uint8(2), &i32))
		assert.Equal(t, int32(2), i32)

		assert.Nil(t, To(uint16(3), &i32))
		assert.Equal(t, int32(3), i32)

		assert.Nil(t, To(uint32(4), &i32))
		assert.Equal(t, int32(4), i32)

		assert.Nil(t, To(uint64(5), &i32))
		assert.Equal(t, int32(5), i32)

		// floats
		assert.Nil(t, To(float32(1), &i32))
		assert.Equal(t, int32(1), i32)

		assert.Nil(t, To(2.0, &i32))
		assert.Equal(t, int32(2), i32)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &i32))
		assert.Equal(t, int32(1), i32)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(2))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 18446744073709551614 cannot be converted to int64"), To(bi, &i32))

		assert.Nil(t, To(big.NewFloat(2), &i32))
		assert.Equal(t, int32(2), i32)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to int64"), To(big.NewFloat(1.25), &i32))

		assert.Nil(t, To(big.NewRat(3, 1), &i32))
		assert.Equal(t, int32(3), i32)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to int64"), To(big.NewRat(5, 4), &i32))

		// string
		assert.Nil(t, To("1", &i32))
		assert.Equal(t, int32(1), i32)

    assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "int64"), To("a", &i32))
		assert.Equal(t, int32(1), i32)
	}

	// == int64
	{
		var i64 int64

		// ints
		assert.Nil(t, To(-1, &i64))
		assert.Equal(t, int64(-1), i64)

		assert.Nil(t, To(int8(-2), &i64))
		assert.Equal(t, int64(-2), i64)

		assert.Nil(t, To(int16(-3), &i64))
		assert.Equal(t, int64(-3), i64)

		assert.Nil(t, To(int32(-4), &i64))
		assert.Equal(t, int64(-4), i64)

		assert.Nil(t, To(int64(-5), &i64))
		assert.Equal(t, int64(-5), i64)

		// uints
		assert.Nil(t, To(uint(1), &i64))
		assert.Equal(t, int64(1), i64)

		assert.Nil(t, To(uint8(2), &i64))
		assert.Equal(t, int64(2), i64)

		assert.Nil(t, To(uint16(3), &i64))
		assert.Equal(t, int64(3), i64)

		assert.Nil(t, To(uint32(4), &i64))
		assert.Equal(t, int64(4), i64)

		assert.Nil(t, To(uint64(5), &i64))
		assert.Equal(t, int64(5), i64)

		// floats
		assert.Nil(t, To(float32(1), &i64))
		assert.Equal(t, int64(1), i64)

		assert.Nil(t, To(2.0, &i64))
		assert.Equal(t, int64(2), i64)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &i64))
		assert.Equal(t, int64(1), i64)

		assert.Nil(t, To(big.NewFloat(2), &i64))
		assert.Equal(t, int64(2), i64)

		assert.Nil(t, To(big.NewRat(3, 1), &i64))
		assert.Equal(t, int64(3), i64)

		// string
		assert.Nil(t, To("1", &i64))
		assert.Equal(t, int64(1), i64)

    assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "int64"), To("a", &i64))
		assert.Equal(t, int64(0), i64)
	}

	// == uint
	{
		var ui uint

		// ints
		assert.Nil(t, To(1, &ui))
		assert.Equal(t, uint(1), ui)

		assert.Nil(t, To(int8(2), &ui))
		assert.Equal(t, uint(2), ui)

		assert.Nil(t, To(int16(3), &ui))
		assert.Equal(t, uint(3), ui)

		assert.Nil(t, To(int32(4), &ui))
		assert.Equal(t, uint(4), ui)

		assert.Nil(t, To(int64(5), &ui))
		assert.Equal(t, uint(5), ui)

		// uints
		assert.Nil(t, To(uint(1), &ui))
		assert.Equal(t, uint(1), ui)

		assert.Nil(t, To(uint8(2), &ui))
		assert.Equal(t, uint(2), ui)

		assert.Nil(t, To(uint16(3), &ui))
		assert.Equal(t, uint(3), ui)

		assert.Nil(t, To(uint32(4), &ui))
		assert.Equal(t, uint(4), ui)

		assert.Nil(t, To(uint64(5), &ui))
		assert.Equal(t, uint(5), ui)

		// floats
		assert.Nil(t, To(float32(1), &ui))
		assert.Equal(t, uint(1), ui)

		assert.Nil(t, To(2.0, &ui))
		assert.Equal(t, uint(2), ui)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &ui))
		assert.Equal(t, uint(1), ui)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(4))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 36893488147419103228 cannot be converted to uint64"), To(bi, &ui))

		assert.Nil(t, To(big.NewFloat(2), &ui))
		assert.Equal(t, uint(2), ui)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to uint64"), To(big.NewFloat(1.25), &ui))

		assert.Nil(t, To(big.NewRat(3, 1), &ui))
		assert.Equal(t, uint(3), ui)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to uint64"), To(big.NewRat(5, 4), &ui))

		// string
		assert.Nil(t, To("1", &ui))
		assert.Equal(t, uint(1), ui)

    assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "uint64"), To("a", &ui))
		assert.Equal(t, uint(1), ui)
	}

	// == uint8
	{
		var ui8 uint8

		// ints
		assert.Nil(t, To(1, &ui8))
		assert.Equal(t, uint8(1), ui8)

		assert.Nil(t, To(int8(2), &ui8))
		assert.Equal(t, uint8(2), ui8)

		assert.Nil(t, To(int16(3), &ui8))
		assert.Equal(t, uint8(3), ui8)

		assert.Nil(t, To(int32(4), &ui8))
		assert.Equal(t, uint8(4), ui8)

		assert.Nil(t, To(int64(5), &ui8))
		assert.Equal(t, uint8(5), ui8)

		// uints
		assert.Nil(t, To(uint(1), &ui8))
		assert.Equal(t, uint8(1), ui8)

		assert.Nil(t, To(uint8(2), &ui8))
		assert.Equal(t, uint8(2), ui8)

		assert.Nil(t, To(uint16(3), &ui8))
		assert.Equal(t, uint8(3), ui8)

		assert.Nil(t, To(uint32(4), &ui8))
		assert.Equal(t, uint8(4), ui8)

		assert.Nil(t, To(uint64(5), &ui8))
		assert.Equal(t, uint8(5), ui8)

		// floats
		assert.Nil(t, To(float32(1), &ui8))
		assert.Equal(t, uint8(1), ui8)

		assert.Nil(t, To(2.0, &ui8))
		assert.Equal(t, uint8(2), ui8)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &ui8))
		assert.Equal(t, uint8(1), ui8)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(4))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 36893488147419103228 cannot be converted to uint64"), To(bi, &ui8))

		assert.Nil(t, To(big.NewFloat(2), &ui8))
		assert.Equal(t, uint8(2), ui8)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to uint64"), To(big.NewFloat(1.25), &ui8))

		assert.Nil(t, To(big.NewRat(3, 1), &ui8))
		assert.Equal(t, uint8(3), ui8)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to uint64"), To(big.NewRat(5, 4), &ui8))

		// string
		assert.Nil(t, To("1", &ui8))
		assert.Equal(t, uint8(1), ui8)

    assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "uint64"), To("a", &ui8))
		assert.Equal(t, uint8(1), ui8)
	}

	// == uint16
	{
		var ui16 uint16

		// ints
		assert.Nil(t, To(1, &ui16))
		assert.Equal(t, uint16(1), ui16)

		assert.Nil(t, To(int8(2), &ui16))
		assert.Equal(t, uint16(2), ui16)

		assert.Nil(t, To(int16(3), &ui16))
		assert.Equal(t, uint16(3), ui16)

		assert.Nil(t, To(int32(4), &ui16))
		assert.Equal(t, uint16(4), ui16)

		assert.Nil(t, To(int64(5), &ui16))
		assert.Equal(t, uint16(5), ui16)

		// uints
		assert.Nil(t, To(uint(1), &ui16))
		assert.Equal(t, uint16(1), ui16)

		assert.Nil(t, To(uint8(2), &ui16))
		assert.Equal(t, uint16(2), ui16)

		assert.Nil(t, To(uint16(3), &ui16))
		assert.Equal(t, uint16(3), ui16)

		assert.Nil(t, To(uint32(4), &ui16))
		assert.Equal(t, uint16(4), ui16)

		assert.Nil(t, To(uint64(5), &ui16))
		assert.Equal(t, uint16(5), ui16)

		// floats
		assert.Nil(t, To(float32(1), &ui16))
		assert.Equal(t, uint16(1), ui16)

		assert.Nil(t, To(2.0, &ui16))
		assert.Equal(t, uint16(2), ui16)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &ui16))
		assert.Equal(t, uint16(1), ui16)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(4))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 36893488147419103228 cannot be converted to uint64"), To(bi, &ui16))

		assert.Nil(t, To(big.NewFloat(2), &ui16))
		assert.Equal(t, uint16(2), ui16)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to uint64"), To(big.NewFloat(1.25), &ui16))

		assert.Nil(t, To(big.NewRat(3, 1), &ui16))
		assert.Equal(t, uint16(3), ui16)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to uint64"), To(big.NewRat(5, 4), &ui16))

		// string
		assert.Nil(t, To("1", &ui16))
		assert.Equal(t, uint16(1), ui16)

    assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "uint64"), To("a", &ui16))
		assert.Equal(t, uint16(1), ui16)
	}

	// == uint32
	{
		var ui32 uint32

		// ints
		assert.Nil(t, To(1, &ui32))
		assert.Equal(t, uint32(1), ui32)

		assert.Nil(t, To(int8(2), &ui32))
		assert.Equal(t, uint32(2), ui32)

		assert.Nil(t, To(int16(3), &ui32))
		assert.Equal(t, uint32(3), ui32)

		assert.Nil(t, To(int32(4), &ui32))
		assert.Equal(t, uint32(4), ui32)

		assert.Nil(t, To(int64(5), &ui32))
		assert.Equal(t, uint32(5), ui32)

		// uints
		assert.Nil(t, To(uint(1), &ui32))
		assert.Equal(t, uint32(1), ui32)

		assert.Nil(t, To(uint8(2), &ui32))
		assert.Equal(t, uint32(2), ui32)

		assert.Nil(t, To(uint16(3), &ui32))
		assert.Equal(t, uint32(3), ui32)

		assert.Nil(t, To(uint32(4), &ui32))
		assert.Equal(t, uint32(4), ui32)

		assert.Nil(t, To(uint64(5), &ui32))
		assert.Equal(t, uint32(5), ui32)

		// floats
		assert.Nil(t, To(float32(1), &ui32))
		assert.Equal(t, uint32(1), ui32)

		assert.Nil(t, To(2.0, &ui32))
		assert.Equal(t, uint32(2), ui32)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &ui32))
		assert.Equal(t, uint32(1), ui32)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(4))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 36893488147419103228 cannot be converted to uint64"), To(bi, &ui32))

		assert.Nil(t, To(big.NewFloat(2), &ui32))
		assert.Equal(t, uint32(2), ui32)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to uint64"), To(big.NewFloat(1.25), &ui32))

		assert.Nil(t, To(big.NewRat(3, 1), &ui32))
		assert.Equal(t, uint32(3), ui32)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to uint64"), To(big.NewRat(5, 4), &ui32))

		// string
		assert.Nil(t, To("1", &ui32))
		assert.Equal(t, uint32(1), ui32)

    assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "uint64"), To("a", &ui32))
		assert.Equal(t, uint32(1), ui32)
	}

	// == uint64
	{
		var ui64 uint64

		// ints
		assert.Nil(t, To(1, &ui64))
		assert.Equal(t, uint64(1), ui64)

		assert.Nil(t, To(int8(2), &ui64))
		assert.Equal(t, uint64(2), ui64)

		assert.Nil(t, To(int16(3), &ui64))
		assert.Equal(t, uint64(3), ui64)

		assert.Nil(t, To(int32(4), &ui64))
		assert.Equal(t, uint64(4), ui64)

		assert.Nil(t, To(int64(5), &ui64))
		assert.Equal(t, uint64(5), ui64)

		// uints
		assert.Nil(t, To(uint(1), &ui64))
		assert.Equal(t, uint64(1), ui64)

		assert.Nil(t, To(uint8(2), &ui64))
		assert.Equal(t, uint64(2), ui64)

		assert.Nil(t, To(uint16(3), &ui64))
		assert.Equal(t, uint64(3), ui64)

		assert.Nil(t, To(uint32(4), &ui64))
		assert.Equal(t, uint64(4), ui64)

		assert.Nil(t, To(uint64(5), &ui64))
		assert.Equal(t, uint64(5), ui64)

		// floats
		assert.Nil(t, To(float32(1), &ui64))
		assert.Equal(t, uint64(1), ui64)

		assert.Nil(t, To(2.0, &ui64))
		assert.Equal(t, uint64(2), ui64)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &ui64))
		assert.Equal(t, uint64(1), ui64)

		assert.Nil(t, To(big.NewFloat(2), &ui64))
		assert.Equal(t, uint64(2), ui64)

		assert.Nil(t, To(big.NewRat(3, 1), &ui64))
		assert.Equal(t, uint64(3), ui64)

		// string
		assert.Nil(t, To("1", &ui64))
		assert.Equal(t, uint64(1), ui64)

    assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "uint64"), To("a", &ui64))
		assert.Equal(t, uint64(0), ui64)
	}

	// == float32
	{
		var f32 float32

		// ints
		assert.Nil(t, To(1, &f32))
		assert.Equal(t, float32(1), f32)

		assert.Nil(t, To(int8(2), &f32))
		assert.Equal(t, float32(2), f32)

		assert.Nil(t, To(int16(3), &f32))
		assert.Equal(t, float32(3), f32)

		assert.Nil(t, To(int32(4), &f32))
		assert.Equal(t, float32(4), f32)

		assert.Nil(t, To(int64(5), &f32))
		assert.Equal(t, float32(5), f32)

		// uints
		assert.Nil(t, To(uint(1), &f32))
		assert.Equal(t, float32(1), f32)

		assert.Nil(t, To(uint8(2), &f32))
		assert.Equal(t, float32(2), f32)

		assert.Nil(t, To(uint16(3), &f32))
		assert.Equal(t, float32(3), f32)

		assert.Nil(t, To(uint32(4), &f32))
		assert.Equal(t, float32(4), f32)

		assert.Nil(t, To(uint64(5), &f32))
		assert.Equal(t, float32(5), f32)

		// floats
		assert.Nil(t, To(float32(1.25), &f32))
		assert.Equal(t, float32(1.25), f32)

		assert.Nil(t, To(2.5, &f32))
		assert.Equal(t, float32(2.5), f32)

		// *bigs
		assert.Nil(t, To(big.NewInt(1), &f32))
		assert.Equal(t, float32(1), f32)
		assert.Equal(t, fmt.Errorf("The *big.Int value of 9223372036854775807 cannot be converted to float64"), To(big.NewInt(math.MaxInt64), &f32))

		assert.Nil(t, To(big.NewFloat(1.25), &f32))
		assert.Equal(t, float32(1.25), f32)

		bf := big.NewFloat(0)
		IntToBigFloat(math.MaxInt64, &bf)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 9223372036854775807 cannot be converted to float64"), To(bf, &f32))

		assert.Nil(t, To(big.NewRat(250, 100), &f32))
		assert.Equal(t, float32(2.5), f32)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 9223372036854775807/1 cannot be converted to float64"), To(big.NewRat(math.MaxInt64, 1), &f32))

		// string
		assert.Nil(t, To("1.25", &f32))
		assert.Equal(t, float32(1.25), f32)

    assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "float32"), To("a", &f32))
		assert.Equal(t, float32(1.25), f32)
	}

	// == float64
	{
		var f64 float64

		// ints
		assert.Nil(t, To(1, &f64))
		assert.Equal(t, 1.0, f64)

		assert.Nil(t, To(int8(2), &f64))
		assert.Equal(t, 2.0, f64)

		assert.Nil(t, To(int16(3), &f64))
		assert.Equal(t, 3.0, f64)

		assert.Nil(t, To(int32(4), &f64))
		assert.Equal(t, 4.0, f64)

		assert.Nil(t, To(int64(5), &f64))
		assert.Equal(t, 5.0, f64)

		// uints
		assert.Nil(t, To(uint(1), &f64))
		assert.Equal(t, 1.0, f64)

		assert.Nil(t, To(uint8(2), &f64))
		assert.Equal(t, 2.0, f64)

		assert.Nil(t, To(uint16(3), &f64))
		assert.Equal(t, 3.0, f64)

		assert.Nil(t, To(uint32(4), &f64))
		assert.Equal(t, 4.0, f64)

		assert.Nil(t, To(uint64(5), &f64))
		assert.Equal(t, 5.0, f64)

		// floats
		assert.Nil(t, To(float32(1.25), &f64))
		assert.Equal(t, 1.25, f64)

		assert.Nil(t, To(2.5, &f64))
		assert.Equal(t, 2.5, f64)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &f64))
		assert.Equal(t, 1.0, f64)

		assert.Nil(t, To(big.NewFloat(1.25), &f64))
		assert.Equal(t, 1.25, f64)

		assert.Nil(t, To(big.NewRat(250, 100), &f64))
		assert.Equal(t, 2.5, f64)

		// string
		assert.Nil(t, To("1.25", &f64))
		assert.Equal(t, 1.25, f64)

    assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "float64"), To("a", &f64))
		assert.Equal(t, 1.25, f64)
	}

	// == *big.Int
	{
		var bi *big.Int

		// ints
		assert.Nil(t, To(1, &bi))
		assert.Equal(t, big.NewInt(1), bi)

		assert.Nil(t, To(int8(2), &bi))
		assert.Equal(t, big.NewInt(2), bi)

		assert.Nil(t, To(int16(3), &bi))
		assert.Equal(t, big.NewInt(3), bi)

		assert.Nil(t, To(int32(4), &bi))
		assert.Equal(t, big.NewInt(4), bi)

		assert.Nil(t, To(int64(5), &bi))
		assert.Equal(t, big.NewInt(5), bi)

		// uints
		assert.Nil(t, To(uint(1), &bi))
		assert.Equal(t, big.NewInt(1), bi)

		assert.Nil(t, To(uint8(2), &bi))
		assert.Equal(t, big.NewInt(2), bi)

		assert.Nil(t, To(uint16(3), &bi))
		assert.Equal(t, big.NewInt(3), bi)

		assert.Nil(t, To(uint32(4), &bi))
		assert.Equal(t, big.NewInt(4), bi)

		assert.Nil(t, To(uint64(5), &bi))
		assert.Equal(t, big.NewInt(5), bi)

		// floats
		assert.Nil(t, To(float32(125), &bi))
		assert.Equal(t, big.NewInt(125), bi)

		assert.Nil(t, To(25.0, &bi))
		assert.Equal(t, big.NewInt(25), bi)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &bi))
		assert.Equal(t, big.NewInt(1), bi)

		assert.Nil(t, To(big.NewFloat(125), &bi))
		assert.Equal(t, big.NewInt(125), bi)

		assert.Nil(t, To(big.NewRat(250, 1), &bi))
		assert.Equal(t, big.NewInt(250), bi)

		// string
		assert.Nil(t, To("1", &bi))
		assert.Equal(t, big.NewInt(1), bi)

    assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "*big.Int"), To("a", &bi))
		assert.Equal(t, big.NewInt(0), bi)
	}

	// == *big.Float
	{
		var bf *big.Float

		// ints
		assert.Nil(t, To(1, &bf))
		cmp := big.NewFloat(1)
		cmp.SetPrec(4)
		assert.Equal(t, cmp, bf)

		assert.Nil(t, To(int8(2), &bf))
		cmp = big.NewFloat(2)
		cmp.SetPrec(4)
		assert.Equal(t, cmp, bf)

		assert.Nil(t, To(int16(3), &bf))
		cmp = big.NewFloat(3)
		cmp.SetPrec(4)
		assert.Equal(t, cmp, bf)

		assert.Nil(t, To(int32(4), &bf))
		cmp = big.NewFloat(4)
		cmp.SetPrec(4)
		assert.Equal(t, cmp, bf)

		assert.Nil(t, To(int64(5), &bf))
		cmp = big.NewFloat(5)
		cmp.SetPrec(4)
		assert.Equal(t, cmp, bf)

		// uints
		assert.Nil(t, To(uint(1), &bf))
		assert.Equal(t, big.NewFloat(1), bf)

		assert.Nil(t, To(uint8(2), &bf))
		assert.Equal(t, big.NewFloat(2), bf)

		assert.Nil(t, To(uint16(3), &bf))
		assert.Equal(t, big.NewFloat(3), bf)

		assert.Nil(t, To(uint32(4), &bf))
		assert.Equal(t, big.NewFloat(4), bf)

		assert.Nil(t, To(uint64(5), &bf))
		assert.Equal(t, big.NewFloat(5), bf)

		// floats
		assert.Nil(t, To(float32(1.25), &bf))
		assert.Equal(t, big.NewFloat(1.25), bf)

		assert.Nil(t, To(2.5, &bf))
		assert.Equal(t, big.NewFloat(2.5), bf)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &bf))
		assert.Equal(t, big.NewFloat(1), bf)

		assert.Nil(t, To(big.NewFloat(1.25), &bf))
		assert.Equal(t, big.NewFloat(1.25), bf)

		assert.Nil(t, To(big.NewRat(250, 100), &bf))
		assert.Equal(t, big.NewFloat(2.5), bf)

		// string
		assert.Nil(t, To("1.25", &bf))
		assert.Equal(t, big.NewFloat(1.25), bf)

    assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "*big.Float"), To("a", &bf))
		assert.Equal(t, (*big.Float)(nil), bf)
	}

	// == *big.Rat
	{
		var br *big.Rat

		// ints
		assert.Nil(t, To(1, &br))
		assert.Equal(t, big.NewRat(1, 1), br)

		assert.Nil(t, To(int8(2), &br))
		assert.Equal(t, big.NewRat(2, 1), br)

		assert.Nil(t, To(int16(3), &br))
		assert.Equal(t, big.NewRat(3, 1), br)

		assert.Nil(t, To(int32(4), &br))
		assert.Equal(t, big.NewRat(4, 1), br)

		assert.Nil(t, To(int64(5), &br))
		assert.Equal(t, big.NewRat(5, 1), br)

		// uints
		assert.Nil(t, To(uint(1), &br))
		assert.Equal(t, big.NewRat(1, 1), br)

		assert.Nil(t, To(uint8(2), &br))
		assert.Equal(t, big.NewRat(2, 1), br)

		assert.Nil(t, To(uint16(3), &br))
		assert.Equal(t, big.NewRat(3, 1), br)

		assert.Nil(t, To(uint32(4), &br))
		assert.Equal(t, big.NewRat(4, 1), br)

		assert.Nil(t, To(uint64(5), &br))
		assert.Equal(t, big.NewRat(5, 1), br)

		// floats
		assert.Nil(t, To(float32(1.25), &br))
		assert.Equal(t, big.NewRat(125, 100), br)

		assert.Nil(t, To(2.5, &br))
		assert.Equal(t, big.NewRat(25, 10), br)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &br))
		assert.Equal(t, big.NewRat(1, 1), br)

		assert.Nil(t, To(big.NewFloat(1.25), &br))
		assert.Equal(t, big.NewRat(125, 100), br)

		assert.Nil(t, To(big.NewRat(25, 10), &br))
		assert.Equal(t, big.NewRat(25, 10), br)

		// string
		assert.Nil(t, To("5/4", &br))
		assert.Equal(t, big.NewRat(5, 4), br)

    assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "*big.Rat"), To("a", &br))
		assert.Equal(t, (*big.Rat)(nil), br)
	}

	// == string
	{
		var s string

		// ints
		assert.Nil(t, To(1, &s))
		assert.Equal(t, "1", s)

		assert.Nil(t, To(int8(2), &s))
		assert.Equal(t, "2", s)

		assert.Nil(t, To(int16(3), &s))
		assert.Equal(t, "3", s)

		assert.Nil(t, To(int32(4), &s))
		assert.Equal(t, "4", s)

		assert.Nil(t, To(int64(5), &s))
		assert.Equal(t, "5", s)

		// uints
		assert.Nil(t, To(uint(1), &s))
		assert.Equal(t, "1", s)

		assert.Nil(t, To(uint8(2), &s))
		assert.Equal(t, "2", s)

		assert.Nil(t, To(uint16(3), &s))
		assert.Equal(t, "3", s)

		assert.Nil(t, To(uint32(4), &s))
		assert.Equal(t, "4", s)

		assert.Nil(t, To(uint64(5), &s))
		assert.Equal(t, "5", s)

		// floats
		assert.Nil(t, To(float32(1.25), &s))
		assert.Equal(t, "1.25", s)

		assert.Nil(t, To(2.5, &s))
		assert.Equal(t, "2.5", s)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &s))
		assert.Equal(t, "1", s)

		assert.Nil(t, To(big.NewFloat(1.25), &s))
		assert.Equal(t, "1.25", s)

		assert.Nil(t, To(big.NewRat(25, 10), &s))
		assert.Equal(t, "5/2", s)

		// string
		assert.Nil(t, To("foo", &s))
		assert.Equal(t, "foo", s)
	}

	{
		// byte to rune, which is really uint8 to int32
		// verify subtypes are handled correctly
		var r rune
		assert.Nil(t, To(byte('A'), &r))
		assert.Equal(t, 'A', r)
	}
}

func TestToBigOps_(t *testing.T) {
	{
		var bi *big.Int
		assert.Nil(t, ToBigOps(1, &bi))
		assert.Equal(t, big.NewInt(1), bi)

		assert.Nil(t, ToBigOps(bi, &bi))
		assert.Equal(t, big.NewInt(1), bi)

		// byte to *big.Int, which is relly uint8 to *big.Int
		// verify subtypes are handled correctly
		assert.Nil(t, To(byte('A'), &bi))
		assert.Equal(t, big.NewInt('A'), bi)
	}

	{
		var bf *big.Float
		assert.Nil(t, ToBigOps(2, &bf))
		cmp := big.NewFloat(2)
		cmp.SetPrec(uint(math.Ceil(1 * log2Of10)))
		assert.Equal(t, cmp, bf)

		assert.Nil(t, ToBigOps(bf, &bf))
		assert.Equal(t, cmp, bf)
	}

	{
		var br *big.Rat
		assert.Nil(t, ToBigOps(3, &br))
		assert.Equal(t, big.NewRat(3, 1), br)

		assert.Nil(t, ToBigOps(br, &br))
		assert.Equal(t, big.NewRat(3, 1), br)
	}
}
