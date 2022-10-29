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

func TestIntToString(t *testing.T) {
	assert.Equal(t, IntToString(int8(1)), "1")
	assert.Equal(t, IntToString(int16(2)), "2")
	assert.Equal(t, IntToString(int32(3)), "3")
	assert.Equal(t, IntToString(int64(4)), "4")
	assert.Equal(t, IntToString(int(5)), "5")
}

func TestUintToString(t *testing.T) {
	assert.Equal(t, UintToString(uint8(1)), "1")
	assert.Equal(t, UintToString(uint16(2)), "2")
	assert.Equal(t, UintToString(uint32(3)), "3")
	assert.Equal(t, UintToString(uint64(4)), "4")
	assert.Equal(t, UintToString(uint(5)), "5")
}

func TestFloatToString(t *testing.T) {
	assert.Equal(t, "1.25", FloatToString(float32(1.25)))
	assert.Equal(t, "1.25", FloatToString(float64(1.25)))
	assert.Equal(t, "-Inf", FloatToString(float32(math.Inf(-1))))
	assert.Equal(t, "+Inf", FloatToString(math.Inf(1)))
}

func TestBigIntToString(t *testing.T) {
	assert.Equal(t, "1234", BigIntToString(big.NewInt(1234)))
}

func TestBigFloatToString(t *testing.T) {
	assert.Equal(t, "1234.5678", BigFloatToString(big.NewFloat(1234.5678)))
}

func TestBigRatToString(t *testing.T) {
	assert.Equal(t, "5/4", BigRatToString(big.NewRat(125, 100)))
}

func TestBigRatToNormalizedString(t *testing.T) {
	assert.Equal(t, "1234", BigRatToNormalizedString(big.NewRat(1234, 1)))
	assert.Equal(t, "1.25", BigRatToNormalizedString(big.NewRat(125, 100)))
}

func TestToString(t *testing.T) {
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

func TestNumBits(t *testing.T) {
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

func TestIntToInt(t *testing.T) {
	var d int8
	assert.Nil(t, IntToInt(1, &d))
	assert.Equal(t, int8(1), d)

	var s = math.MaxInt8
	assert.Nil(t, IntToInt(s, &d))
	assert.Equal(t, int8(math.MaxInt8), d)

	assert.Equal(t, fmt.Errorf(errMsg, math.MinInt16, fmt.Sprintf("%d", math.MinInt16), "int8"), IntToInt(math.MinInt16, &d))
	assert.Equal(t, fmt.Errorf(errMsg, math.MaxInt16, fmt.Sprintf("%d", math.MaxInt16), "int8"), IntToInt(math.MaxInt16, &d))
}

func TestIntToUint(t *testing.T) {
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
		assert.Equal(t, fmt.Errorf(errMsg, -1, "-1", "uint8"), IntToUint(-1, &d))
		assert.Equal(t, fmt.Errorf(errMsg, math.MaxUint16, fmt.Sprintf("%d", math.MaxUint16), "uint8"), IntToUint(math.MaxUint16, &d))
	}
}

func TestUintToInt(t *testing.T) {
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
		assert.Equal(t, fmt.Errorf(errMsg, uint(0), fmt.Sprintf("%d", math.MaxInt16), "int8"), UintToInt(uint(math.MaxInt16), &d))
	}
}

func TestUintToUint(t *testing.T) {
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
		assert.Equal(t, fmt.Errorf(errMsg, uint(0), fmt.Sprintf("%d", math.MaxInt16), "uint8"), UintToUint(uint(math.MaxInt16), &d))
	}
}

func TestIntToFloat(t *testing.T) {
	// The following code tries int values that start at maximum value a float32 can hold, and continue for 8 values after
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
			assert.Equal(t, fmt.Errorf(errMsg, ival, fmt.Sprintf("%d", ival), "float32"), IntToFloat(ival, &fval))
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
			assert.Equal(t, fmt.Errorf(errMsg, ival, fmt.Sprintf("%d", ival), "float64"), IntToFloat(ival, &fval))
		}
	}
}

func TestFloatToInt(t *testing.T) {
	var d int
	assert.Nil(t, FloatToInt(float32(125), &d))
	assert.Equal(t, 125, d)

	assert.Equal(t, fmt.Errorf(errMsg, 1.25, fmt.Sprintf("%f", 1.25), "int"), FloatToInt(1.25, &d))
}

func TestFloatToUint(t *testing.T) {
	var d uint
	assert.Nil(t, FloatToUint(float32(125), &d))
	assert.Equal(t, uint(125), d)

	assert.Equal(t, fmt.Errorf(errMsg, 1.25, fmt.Sprintf("%f", 1.25), "uint"), FloatToUint(1.25, &d))
}

func TestFloat64ToFloat32(t *testing.T) {
	var (
		i float64
		o float32
	)
	assert.Nil(t, Float64ToFloat32(i, &o))
	assert.Equal(t, float32(0), o)

	i = 1
	assert.Nil(t, Float64ToFloat32(i, &o))
	assert.Equal(t, float32(1), o)

	assert.Nil(t, Float64ToFloat32(math.SmallestNonzeroFloat32, &o))
	assert.Equal(t, float32(math.SmallestNonzeroFloat32), o)

	assert.Nil(t, Float64ToFloat32(math.MaxFloat32, &o))
	assert.Equal(t, float32(math.MaxFloat32), o)

	assert.Nil(t, Float64ToFloat32(math.Inf(-1), &o))
	assert.Equal(t, float32(math.Inf(-1)), o)

	assert.Nil(t, Float64ToFloat32(math.Inf(1), &o))
	assert.Equal(t, float32(math.Inf(1)), o)

	i = float64(math.SmallestNonzeroFloat32) - 1
	assert.Equal(t, fmt.Errorf(errMsg, i, fmt.Sprintf("%f", i), "float32"), Float64ToFloat32(i, &o))
}

// ==== ToInt64

func TestBigIntToInt64(t *testing.T) {
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
	assert.Equal(t, fmt.Errorf(errMsg, big.NewInt(0), str, "int64"), BigIntToInt64(inter, &o))
}

func TestBigFloatToInt64(t *testing.T) {
	var o int64
	assert.Nil(t, BigFloatToInt64(big.NewFloat(1), &o))
	assert.Equal(t, int64(1), o)

	assert.Nil(t, BigFloatToInt64(big.NewFloat(100_000), &o))
	assert.Equal(t, int64(100_000), o)

	assert.Equal(t, fmt.Errorf(errMsg, big.NewFloat(0), "1.25", "int64"), BigFloatToInt64(big.NewFloat(1.25), &o))

	negInf := big.NewFloat(0).Quo(big.NewFloat(-1), big.NewFloat(0))
	assert.Equal(t, fmt.Errorf(errMsg, negInf, negInf.String(), "int64"), BigFloatToInt64(negInf, &o))

	posInf := big.NewFloat(0).Quo(big.NewFloat(1), big.NewFloat(0))
	assert.Equal(t, fmt.Errorf(errMsg, posInf, posInf.String(), "int64"), BigFloatToInt64(posInf, &o))
}

func TestBigRatToInt64(t *testing.T) {
	var o int64
	assert.Nil(t, BigRatToInt64(big.NewRat(1, 1), &o))
	assert.Equal(t, int64(1), o)

	assert.Nil(t, BigRatToInt64(big.NewRat(100_000, 1), &o))
	assert.Equal(t, int64(100_000), o)

	assert.Equal(t, fmt.Errorf(errMsg, big.NewRat(1, 1), "5/4", "int64"), BigRatToInt64(big.NewRat(125, 100), &o))
}

func TestStringToInt64(t *testing.T) {
	var o int64
	assert.Nil(t, StringToInt64("1", &o))
	assert.Equal(t, int64(1), o)

	assert.Nil(t, StringToInt64("100000", &o))
	assert.Equal(t, int64(100_000), o)

	str := "123456789012345678901"
	assert.Equal(t, fmt.Errorf(errMsg, str, str, "int64"), StringToInt64(str, &o))
}

// ==== ToUint64

func TestBigIntToUint64(t *testing.T) {
	var o uint64
	assert.Nil(t, BigIntToUint64(big.NewInt(1), &o))
	assert.Equal(t, uint64(1), o)

	assert.Nil(t, BigIntToUint64(big.NewInt(100_000), &o))
	assert.Equal(t, uint64(100_000), o)

	var inter *big.Int
	StringToBigInt("123456789012345678901", &inter)
	assert.Equal(t, fmt.Errorf(errMsg, big.NewInt(0), "123456789012345678901", "uint64"), BigIntToUint64(inter, &o))
}

func TestBigFloatToUint64(t *testing.T) {
	var o uint64
	assert.Nil(t, BigFloatToUint64(big.NewFloat(1), &o))
	assert.Equal(t, uint64(1), o)

	assert.Nil(t, BigFloatToUint64(big.NewFloat(100_000), &o))
	assert.Equal(t, uint64(100_000), o)

	assert.Equal(t, fmt.Errorf(errMsg, big.NewFloat(0), "1.25", "uint64"), BigFloatToUint64(big.NewFloat(1.25), &o))
}

func TestBigRatToUint64(t *testing.T) {
	var o uint64
	assert.Nil(t, BigRatToUint64(big.NewRat(1, 1), &o))
	assert.Equal(t, uint64(1), o)

	assert.Nil(t, BigRatToUint64(big.NewRat(100_000, 1), &o))
	assert.Equal(t, uint64(100_000), o)

	assert.Equal(t, fmt.Errorf(errMsg, big.NewRat(1, 1), "5/4", "uint64"), BigRatToUint64(big.NewRat(125, 100), &o))
}

func TestStringToUint64(t *testing.T) {
	var o uint64
	assert.Nil(t, StringToUint64("1", &o))
	assert.Equal(t, uint64(1), o)

	assert.Nil(t, StringToUint64("100000", &o))
	assert.Equal(t, uint64(100_000), o)

	str := "123456789012345678901"
	assert.Equal(t, fmt.Errorf(errMsg, str, str, "uint64"), StringToUint64(str, &o))
}

// ==== ToFloat32

func TestBigIntToFloat32(t *testing.T) {
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
	assert.Equal(t, fmt.Errorf(errMsg, big.NewInt(0), str, "float32"), BigIntToFloat32(inter, &o))
}

func TestBigFloatToFloat32(t *testing.T) {
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
	assert.Equal(t, fmt.Errorf(errMsg, big.NewFloat(0), str, "float32"), BigFloatToFloat32(inter, &o))
}

func TestBigRatToFloat32(t *testing.T) {
	var o float32
	assert.Nil(t, BigRatToFloat32(big.NewRat(1, 1), &o))
	assert.Equal(t, float32(1), o)

	assert.Nil(t, BigRatToFloat32(big.NewRat(125, 100), &o))
	assert.Equal(t, float32(1.25), o)

	i := big.NewRat(math.MaxInt64, 100)
	assert.Equal(t, fmt.Errorf(errMsg, i, i.String(), "float32"), BigRatToFloat32(i, &o))
}

func TestStringToFloat32(t *testing.T) {
	var o float32
	assert.Nil(t, StringToFloat32("1", &o))
	assert.Equal(t, float32(1), o)

	assert.Nil(t, StringToFloat32("1.25", &o))
	assert.Equal(t, float32(1.25), o)

	str := "123456789012345678901"
	assert.Equal(t, fmt.Errorf(errMsg, str, str, "float32"), StringToFloat32(str, &o))
}

// ==== ToFloat64

func TestBigIntToFloat64(t *testing.T) {
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
	assert.Equal(t, fmt.Errorf(errMsg, big.NewInt(0), str, "float64"), BigIntToFloat64(inter, &o))
}

func TestBigFloatToFloat64(t *testing.T) {
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
	assert.Equal(t, fmt.Errorf(errMsg, big.NewFloat(0), str, "float64"), BigFloatToFloat64(inter, &o))
}

func TestBigRatToFloat64(t *testing.T) {
	var o float64
	assert.Nil(t, BigRatToFloat64(big.NewRat(1, 1), &o))
	assert.Equal(t, float64(1), o)

	assert.Nil(t, BigRatToFloat64(big.NewRat(125, 100), &o))
	assert.Equal(t, float64(1.25), o)

	i := big.NewRat(math.MaxInt64, 100)
	assert.Equal(t, fmt.Errorf(errMsg, i, i.String(), "float64"), BigRatToFloat64(i, &o))
}

func TestStringToFloat64(t *testing.T) {
	var o float64
	assert.Nil(t, StringToFloat64("1", &o))
	assert.Equal(t, float64(1), o)

	assert.Nil(t, StringToFloat64("1.25", &o))
	assert.Equal(t, float64(1.25), o)

	str := "123456789012345678901"
	assert.Equal(t, fmt.Errorf(errMsg, str, str, "float64"), StringToFloat64(str, &o))
}

// ==== ToBigInt

func TestIntToBigInt(t *testing.T) {
	var o *big.Int
	IntToBigInt(int8(1), &o)
	assert.Equal(t, big.NewInt(1), o)

	IntToBigInt(100_000, &o)
	assert.Equal(t, big.NewInt(100_000), o)
}

func TestUintToBigInt(t *testing.T) {
	var o *big.Int
	UintToBigInt(uint8(1), &o)
	assert.Equal(t, big.NewInt(1), o)

	UintToBigInt(uint(100_000), &o)
	assert.Equal(t, big.NewInt(100_000), o)
}

func TestFloatToBigInt(t *testing.T) {
	var o *big.Int
	assert.Nil(t, FloatToBigInt(float32(1), &o))
	assert.Equal(t, big.NewInt(1), o)

	assert.Nil(t, FloatToBigInt(float32(100_000), &o))
	assert.Equal(t, big.NewInt(100_000), o)

	assert.Equal(t, fmt.Errorf(errMsg, 1.25, fmt.Sprintf("%f", 1.25), "*big.Int"), FloatToBigInt(1.25, &o))
}

func TestBigFloatToBigInt(t *testing.T) {
	var o *big.Int
	assert.Nil(t, BigFloatToBigInt(big.NewFloat(1), &o))
	assert.Equal(t, big.NewInt(1), o)

	assert.Nil(t, BigFloatToBigInt(big.NewFloat(100_000), &o))
	assert.Equal(t, big.NewInt(100_000), o)

	assert.Equal(t, fmt.Errorf(errMsg, big.NewFloat(1.25), big.NewFloat(1.25).String(), "*big.Int"), BigFloatToBigInt(big.NewFloat(1.25), &o))
}

func TestBigRatToBigInt(t *testing.T) {
	var o *big.Int
	assert.Nil(t, BigRatToBigInt(big.NewRat(1, 1), &o))
	assert.Equal(t, big.NewInt(1), o)

	assert.Nil(t, BigRatToBigInt(big.NewRat(100_000, 1), &o))
	assert.Equal(t, big.NewInt(100_000), o)

	assert.Equal(t, fmt.Errorf(errMsg, big.NewRat(125, 100), big.NewRat(125, 100).String(), "*big.Int"), BigRatToBigInt(big.NewRat(125, 100), &o))
}

func TestStringToBigInt(t *testing.T) {
	var o *big.Int
	assert.Nil(t, StringToBigInt("1", &o))
	assert.Equal(t, big.NewInt(1), o)

	assert.Nil(t, StringToBigInt("100000", &o))
	assert.Equal(t, big.NewInt(100_000), o)

	str := "1234567890123456789012"
	assert.Nil(t, StringToBigInt(str, &o))
	assert.Equal(t, str, o.String())

	assert.Equal(t, fmt.Errorf(errMsg, "1.25", "1.25", "*big.Int"), StringToBigInt("1.25", &o))
}

// ==== ToBigFloat

func TestIntToBigFloat(t *testing.T) {
	var o *big.Float
	IntToBigFloat(int8(1), &o)
	assert.Equal(t, big.NewFloat(1), o)

	IntToBigFloat(100_000, &o)
	assert.Equal(t, big.NewFloat(100_000), o)
}

func TestUintToBigFloat(t *testing.T) {
	var o *big.Float
	UintToBigFloat(uint8(1), &o)
	assert.Equal(t, big.NewFloat(1), o)

	UintToBigFloat(uint(100_000), &o)
	assert.Equal(t, big.NewFloat(100_000), o)
}

func TestFloatToBigFloat(t *testing.T) {
	var o *big.Float
	FloatToBigFloat(float32(1.25), &o)
	assert.Equal(t, big.NewFloat(1.25), o)

	FloatToBigFloat(float64(100_000), &o)
	assert.Equal(t, big.NewFloat(100_000), o)
}

func TestBigIntToBigFloat(t *testing.T) {
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

func TestBigRatToBigFloat(t *testing.T) {
	var o *big.Float
	BigRatToBigFloat(big.NewRat(125, 100), &o)
	assert.Equal(t, big.NewFloat(1.25), o)

	BigRatToBigFloat(big.NewRat(25, 10), &o)
	assert.Equal(t, big.NewFloat(2.5), o)
}

func TestStringToBigFloat(t *testing.T) {
	var o *big.Float
	assert.Nil(t, StringToBigFloat("1", &o))
	assert.Equal(t, big.NewFloat(1), o)

	assert.Nil(t, StringToBigFloat("100000.25", &o))
	assert.Equal(t, big.NewFloat(100000.25), o)

	str := "1234567890123456789012"
	assert.Nil(t, StringToBigFloat(str, &o))
	assert.Equal(t, str, fmt.Sprintf("%.f", o))

	assert.Equal(t, fmt.Errorf(errMsg, "1.25p", "1.25p", "*big.Float"), StringToBigFloat("1.25p", &o))
}

// ==== ToBigRat

func TestIntToBigRat(t *testing.T) {
	var o *big.Rat
	IntToBigRat(int8(1), &o)
	assert.Equal(t, big.NewRat(1, 1), o)

	IntToBigRat(100_000, &o)
	assert.Equal(t, big.NewRat(100_000, 1), o)
}

func TestUintToBigRat(t *testing.T) {
	var o *big.Rat
	UintToBigRat(uint8(1), &o)
	assert.Equal(t, big.NewRat(1, 1), o)

	UintToBigRat(uint(100_000), &o)
	assert.Equal(t, big.NewRat(100_000, 1), o)
}

func TestFloatToBigRat(t *testing.T) {
	var o *big.Rat
	assert.Nil(t, FloatToBigRat(float32(1.25), &o))
	assert.Equal(t, big.NewRat(125, 100), o)

	assert.Nil(t, FloatToBigRat(float64(2.5), &o))
	assert.Equal(t, big.NewRat(25, 10), o)

	i := math.Inf(-1)
	assert.Equal(t, fmt.Errorf(errMsg, i, fmt.Sprintf("%f", i), "*big.Rat"), FloatToBigRat(i, &o))

	i = math.Inf(1)
	assert.Equal(t, fmt.Errorf(errMsg, i, fmt.Sprintf("%f", i), "*big.Rat"), FloatToBigRat(i, &o))
}

func TestBigIntToBigRat(t *testing.T) {
	var o *big.Rat
	BigIntToBigRat(big.NewInt(1), &o)
	assert.Equal(t, big.NewRat(1, 1), o)

	BigIntToBigRat(big.NewInt(100_000), &o)
	assert.Equal(t, big.NewRat(100_000, 1), o)
}

func TestBigFloatToBigRat(t *testing.T) {
	var o *big.Rat
	assert.Nil(t, BigFloatToBigRat(big.NewFloat(1), &o))
	assert.Equal(t, big.NewRat(1, 1), o)

	assert.Nil(t, BigFloatToBigRat(big.NewFloat(1.25), &o))
	assert.Equal(t, big.NewRat(125, 100), o)

	i := big.NewFloat(math.Inf(-1))
	assert.Equal(t, fmt.Errorf(errMsg, i, i.String(), "*big.Rat"), BigFloatToBigRat(i, &o))

	i = big.NewFloat(math.Inf(1))
	assert.Equal(t, fmt.Errorf(errMsg, i, i.String(), "*big.Rat"), BigFloatToBigRat(i, &o))
}

func TestStringToBigRat(t *testing.T) {
	var o *big.Rat
	assert.Nil(t, StringToBigRat("1/1", &o))
	assert.Equal(t, big.NewRat(1, 1), o)

	assert.Nil(t, StringToBigRat("125/100", &o))
	assert.Equal(t, big.NewRat(125, 100), o)

	assert.Nil(t, StringToBigRat("1.25", &o))
	assert.Equal(t, big.NewRat(125, 100), o)

	assert.Equal(t, fmt.Errorf(errMsg, "1.25p", "1.25p", "*big.Rat"), StringToBigRat("1.25p", &o))
}

func TestFloatStringToBigRat(t *testing.T) {
	var o *big.Rat
	assert.Nil(t, FloatStringToBigRat("1.25", &o))
	assert.Equal(t, big.NewRat(125, 100), o)

	assert.Equal(t, fmt.Errorf(errMsg, "125/100", "125/100", "*big.Rat"), FloatStringToBigRat("125/100", &o))
}
