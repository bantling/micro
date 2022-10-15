package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/bantling/micro/go/funcs"
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
	{
		var d int8
		IntToInt(1, &d)
		assert.Equal(t, int8(1), d)
	}

	{
		var s = math.MaxInt8
		var d int8
		IntToInt(s, &d)
		assert.Equal(t, int8(math.MaxInt8), d)
	}

	funcs.TryTo(
		func() {
			var d int8
			IntToInt(math.MinInt16, &d)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, math.MinInt16, fmt.Sprintf("%d", math.MinInt16), "int8"), e)
		},
	)

	funcs.TryTo(
		func() {
			var d int8
			IntToInt(math.MaxInt16, &d)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, math.MaxInt16, fmt.Sprintf("%d", math.MaxInt16), "int8"), e)
		},
	)
}

func TestIntToUint(t *testing.T) {
	{
		var d uint
		IntToUint(1, &d)
		assert.Equal(t, uint(1), d)
	}

	{
		var s = math.MaxUint16
		var d uint16
		IntToUint(s, &d)
		assert.Equal(t, uint16(math.MaxUint16), d)
	}

	funcs.TryTo(
		func() {
			var d uint8
			IntToUint(-1, &d)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, -1, "-1", "uint8"), e)
		},
	)

	funcs.TryTo(
		func() {
			var d uint8
			IntToUint(math.MaxUint16, &d)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, math.MaxUint16, fmt.Sprintf("%d", math.MaxUint16), "uint8"), e)
		},
	)
}

func TestUintToInt(t *testing.T) {
	{
		var d int
		UintToInt(uint(1), &d)
		assert.Equal(t, 1, d)
	}

	{
		var d int
		UintToInt(uint64(math.MaxInt16), &d)
		assert.Equal(t, math.MaxInt16, d)
	}

	funcs.TryTo(
		func() {
			var d int8
			UintToInt(uint(math.MaxInt16), &d)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, uint(0), fmt.Sprintf("%d", math.MaxInt16), "int8"), e)
		},
	)
}

func TestUintToUint(t *testing.T) {
	{
		var d uint
		UintToUint(uint(1), &d)
		assert.Equal(t, uint(1), d)
	}

	{
		var d uint32
		UintToUint(uint64(math.MaxUint16), &d)
		assert.Equal(t, uint32(math.MaxUint16), d)
	}

	funcs.TryTo(
		func() {
			var d uint8
			UintToUint(uint(math.MaxInt16), &d)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, uint(0), fmt.Sprintf("%d", math.MaxInt16), "uint8"), e)
		},
	)
}

func TestFloatToInt(t *testing.T) {
	var d int
	FloatToInt(float32(125), &d)
	assert.Equal(t, 125, d)

	funcs.TryTo(
		func() {
			FloatToInt(1.25, &d)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, 1.25, fmt.Sprintf("%f", 1.25), "int"), e)
		},
	)
}

func TestFloatToUint(t *testing.T) {
	var d uint
	FloatToUint(float32(125), &d)
	assert.Equal(t, uint(125), d)

	funcs.TryTo(
		func() {
			FloatToUint(1.25, &d)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, 1.25, fmt.Sprintf("%f", 1.25), "uint"), e)
		},
	)
}

func TestFloat64ToFloat32(t *testing.T) {
	assert.Equal(t, float32(1), Float64ToFloat32(1))
	assert.Equal(t, float32(math.SmallestNonzeroFloat32), Float64ToFloat32(math.SmallestNonzeroFloat32))
	assert.Equal(t, float32(math.MaxFloat32), Float64ToFloat32(math.MaxFloat32))
	assert.Equal(t, float32(math.Inf(-1)), Float64ToFloat32(math.Inf(-1)))
	assert.Equal(t, float32(math.Inf(1)), Float64ToFloat32(math.Inf(1)))

	f64 := float64(math.SmallestNonzeroFloat32) - 1
	funcs.TryTo(
		func() {
			Float64ToFloat32(f64)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, f64, fmt.Sprintf("%f", f64), "float32"), e)
		},
	)
}

// ==== ToInt64

func TestBigIntToInt64(t *testing.T) {
	assert.Equal(t, int64(1), BigIntToInt64(big.NewInt(1)))
	assert.Equal(t, int64(100_000), BigIntToInt64(big.NewInt(100_000)))

	funcs.TryTo(
		func() {
			BigIntToInt64(StringToBigInt("123456789012345678901"))
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, big.NewInt(0), "123456789012345678901", "int64"), e)
		},
	)
}

func TestBigFloatToInt64(t *testing.T) {
	assert.Equal(t, int64(1), BigFloatToInt64(big.NewFloat(1)))
	assert.Equal(t, int64(100_000), BigFloatToInt64(big.NewFloat(100_000)))

	funcs.TryTo(
		func() {
			BigFloatToInt64(big.NewFloat(1.25))
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, big.NewFloat(0), "1.25", "int64"), e)
		},
	)

	negInf := big.NewFloat(0).Quo(big.NewFloat(-1), big.NewFloat(0))
	funcs.TryTo(
		func() {
			BigFloatToInt64(negInf)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, negInf, negInf.String(), "int64"), e)
		},
	)

	posInf := big.NewFloat(0).Quo(big.NewFloat(1), big.NewFloat(0))
	funcs.TryTo(
		func() {
			BigFloatToInt64(posInf)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, posInf, posInf.String(), "int64"), e)
		},
	)
}

func TestBigRatToInt64(t *testing.T) {
	assert.Equal(t, int64(1), BigRatToInt64(big.NewRat(1, 1)))
	assert.Equal(t, int64(100_000), BigRatToInt64(big.NewRat(100_000, 1)))

	funcs.TryTo(
		func() {
			BigRatToInt64(big.NewRat(125, 100))
			assert.Fail(t, "Must die")
		},
		func(e any) {
			// 125/100 is reducced to 5/4
			assert.Equal(t, fmt.Errorf(errMsg, big.NewRat(1, 1), "5/4", "int64"), e)
		},
	)
}

func TestStringToInt64(t *testing.T) {
	assert.Equal(t, int64(1), StringToInt64("1"))
	assert.Equal(t, int64(100_000), StringToInt64("100000"))

	str := "123456789012345678901"
	funcs.TryTo(
		func() {
			StringToInt64(str)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, str, str, "int64"), e)
		},
	)
}

// ==== ToUint64

func TestBigIntToUint64(t *testing.T) {
	assert.Equal(t, uint64(1), BigIntToUint64(big.NewInt(1)))
	assert.Equal(t, uint64(100_000), BigIntToUint64(big.NewInt(100_000)))

	funcs.TryTo(
		func() {
			BigIntToUint64(StringToBigInt("123456789012345678901"))
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, big.NewInt(0), "123456789012345678901", "uint64"), e)
		},
	)
}

func TestBigFloatToUint64(t *testing.T) {
	assert.Equal(t, uint64(1), BigFloatToUint64(big.NewFloat(1)))
	assert.Equal(t, uint64(100_000), BigFloatToUint64(big.NewFloat(100_000)))

	funcs.TryTo(
		func() {
			BigFloatToUint64(big.NewFloat(1.25))
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, big.NewFloat(0), "1.25", "uint64"), e)
		},
	)
}

func TestBigRatToUint64(t *testing.T) {
	assert.Equal(t, uint64(1), BigRatToUint64(big.NewRat(1, 1)))
	assert.Equal(t, uint64(100_000), BigRatToUint64(big.NewRat(100_000, 1)))

	funcs.TryTo(
		func() {
			BigRatToUint64(big.NewRat(125, 100))
			assert.Fail(t, "Must die")
		},
		func(e any) {
			// 125/100 is reducced to 5/4
			assert.Equal(t, fmt.Errorf(errMsg, big.NewRat(1, 1), "5/4", "uint64"), e)
		},
	)
}

func TestStringToUint64(t *testing.T) {
	assert.Equal(t, uint64(1), StringToUint64("1"))
	assert.Equal(t, uint64(100_000), StringToUint64("100000"))

	str := "123456789012345678901"
	funcs.TryTo(
		func() {
			StringToUint64(str)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, str, str, "uint64"), e)
		},
	)
}

// ==== ToFloat32

func TestBigIntToFloat32(t *testing.T) {
	assert.Equal(t, float32(1), BigIntToFloat32(big.NewInt(1)))
	assert.Equal(t, float32(100_000), BigIntToFloat32(big.NewInt(100_000)))

	funcs.TryTo(
		func() {
			BigIntToFloat32(StringToBigInt("123456789012345678901"))
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, big.NewInt(0), "123456789012345678901", "float32"), e)
		},
	)
}

func TestBigFloatToFloat32(t *testing.T) {
	assert.Equal(t, float32(1), BigFloatToFloat32(big.NewFloat(1)))
	assert.Equal(t, float32(100_000), BigFloatToFloat32(big.NewFloat(100_000)))

	funcs.TryTo(
		func() {
			BigFloatToFloat32(StringToBigFloat("123456789012345678901"))
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, big.NewFloat(0), "123456789012345678901", "float32"), e)
		},
	)
}

func TestBigRatToFloat32(t *testing.T) {
	assert.Equal(t, float32(1), BigRatToFloat32(big.NewRat(1, 1)))
	assert.Equal(t, float32(1.25), BigRatToFloat32(big.NewRat(125, 100)))

	val := big.NewRat(math.MaxInt64, 100)
	funcs.TryTo(
		func() {
			BigRatToFloat32(val)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, val, val.String(), "float32"), e)
		},
	)
}

func TestStringToFloat32(t *testing.T) {
	assert.Equal(t, float32(1), StringToFloat32("1"))
	assert.Equal(t, float32(1.25), StringToFloat32("1.25"))

	str := "123456789012345678901"
	funcs.TryTo(
		func() {
			StringToFloat32(str)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, str, str, "float32"), e)
		},
	)
}

// ==== ToFloat64

func TestBigIntToFloat64(t *testing.T) {
	assert.Equal(t, float64(1), BigIntToFloat64(big.NewInt(1)))
	assert.Equal(t, float64(100_000), BigIntToFloat64(big.NewInt(100_000)))

	funcs.TryTo(
		func() {
			BigIntToFloat64(StringToBigInt("123456789012345678901"))
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, big.NewInt(0), "123456789012345678901", "float64"), e)
		},
	)
}

func TestBigFloatToFloat64(t *testing.T) {
	assert.Equal(t, float64(1), BigFloatToFloat64(big.NewFloat(1)))
	assert.Equal(t, float64(100_000), BigFloatToFloat64(big.NewFloat(100_000)))

	funcs.TryTo(
		func() {
			BigFloatToFloat64(StringToBigFloat("123456789012345678901"))
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, big.NewFloat(0), "123456789012345678901", "float64"), e)
		},
	)
}

func TestBigRatToFloat64(t *testing.T) {
	assert.Equal(t, float64(1), BigRatToFloat64(big.NewRat(1, 1)))
	assert.Equal(t, float64(1.25), BigRatToFloat64(big.NewRat(125, 100)))

	val := big.NewRat(math.MaxInt64, 100)
	funcs.TryTo(
		func() {
			BigRatToFloat64(val)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, val, val.String(), "float64"), e)
		},
	)
}

func TestStringToFloat64(t *testing.T) {
	assert.Equal(t, float64(1), StringToFloat64("1"))
	assert.Equal(t, float64(1.25), StringToFloat64("1.25"))

	str := "123456789012345678901"
	funcs.TryTo(
		func() {
			StringToFloat64(str)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, str, str, "float64"), e)
		},
	)
}

// ==== ToBigInt

func TestIntToBigInt(t *testing.T) {
	assert.Equal(t, int64(1), IntToBigInt(int8(1)).Int64())
	assert.Equal(t, int64(2), IntToBigInt(2).Int64())
}

func TestUintToBigInt(t *testing.T) {
	assert.Equal(t, uint64(1), UintToBigInt(uint8(1)).Uint64())
	assert.Equal(t, uint64(2), UintToBigInt(uint64(2)).Uint64())
}

func TestFloatToBigInt(t *testing.T) {
	assert.Equal(t, big.NewInt(1), FloatToBigInt(float32(1)))
	assert.Equal(t, big.NewInt(100_000), FloatToBigInt(float64(100_000)))

	funcs.TryTo(
		func() {
			FloatToBigInt(1.25)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, 1.25, fmt.Sprintf("%f", 1.25), "*big.Int"), e)
		},
	)
}

func TestBigFloatToBigInt(t *testing.T) {
	assert.Equal(t, big.NewInt(1), BigFloatToBigInt(big.NewFloat(1)))
	assert.Equal(t, big.NewInt(100_000), BigFloatToBigInt(big.NewFloat(100_000)))

	funcs.TryTo(
		func() {
			BigFloatToBigInt(big.NewFloat(1.25))
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, big.NewFloat(1.25), big.NewFloat(1.25).String(), "*big.Int"), e)
		},
	)
}

func TestBigRatToBigInt(t *testing.T) {
	assert.Equal(t, big.NewInt(1), BigRatToBigInt(big.NewRat(1, 1)))
	assert.Equal(t, big.NewInt(100_000), BigRatToBigInt(big.NewRat(100_000, 1)))

	funcs.TryTo(
		func() {
			BigRatToBigInt(big.NewRat(125, 100))
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, big.NewRat(125, 100), big.NewRat(125, 100).String(), "*big.Int"), e)
		},
	)
}

func TestStringToBigInt(t *testing.T) {
	assert.Equal(t, big.NewInt(1), StringToBigInt("1"))
	assert.Equal(t, big.NewInt(100_000), StringToBigInt("100000"))

	str := "1234567890123456789012"
	assert.Equal(t, str, StringToBigInt(str).String())

	funcs.TryTo(
		func() {
			StringToBigInt("1.25")
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, "1.25", "1.25", "*big.Int"), e)
		},
	)
}

// ==== ToBigFloat

func TestIntToBigFloat(t *testing.T) {
	assert.Equal(t, int64(1), funcs.FirstValue2(IntToBigFloat(int8(1)).Int64()))
	assert.Equal(t, int64(2), funcs.FirstValue2(IntToBigFloat(2).Int64()))
}

func TestUintToBigFloat(t *testing.T) {
	assert.Equal(t, uint64(1), funcs.FirstValue2(UintToBigFloat(uint8(1)).Uint64()))
	assert.Equal(t, uint64(2), funcs.FirstValue2(UintToBigFloat(uint64(2)).Uint64()))
}

func TestFloatToBigFloat(t *testing.T) {
	assert.Equal(t, float64(1.25), funcs.FirstValue2(FloatToBigFloat(float32(1.25)).Float64()))
	assert.Equal(t, float64(2.5), funcs.FirstValue2(FloatToBigFloat(float64(2.5)).Float64()))
}

func TestBigIntToBigFloat(t *testing.T) {
	assert.Equal(t, big.NewFloat(1), BigIntToBigFloat(big.NewInt(1)))
	assert.Equal(t, big.NewFloat(100_000), BigIntToBigFloat(big.NewInt(100_000)))

	str := "1234567890123456789012"
	assert.Equal(t, str, fmt.Sprintf("%.f", BigIntToBigFloat(StringToBigInt(str))))
}

func TestBigRatToBigFloat(t *testing.T) {
	assert.Equal(t, big.NewFloat(1.25), BigRatToBigFloat(big.NewRat(125, 100)))
	assert.Equal(t, big.NewFloat(2.5), BigRatToBigFloat(big.NewRat(25, 10)))
}

func TestStringToBigFloat(t *testing.T) {
	assert.Equal(t, big.NewFloat(1), StringToBigFloat("1"))
	assert.Equal(t, big.NewFloat(100000.25), StringToBigFloat("100000.25"))

	str := "1234567890123456789012"
	assert.Equal(t, str, fmt.Sprintf("%.f", StringToBigFloat(str)))

	funcs.TryTo(
		func() {
			StringToBigFloat("1.25p")
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, "1.25p", "1.25p", "*big.Float"), e)
		},
	)
}

// ==== ToBigRat

func TestIntToBigRat(t *testing.T) {
	assert.Equal(t, big.NewRat(1, 1), IntToBigRat(int8(1)))
	assert.Equal(t, big.NewRat(100_000, 1), IntToBigRat(100_000))
}

func TestUintToBigRat(t *testing.T) {
	assert.Equal(t, big.NewRat(1, 1), UintToBigRat(uint8(1)))
	assert.Equal(t, big.NewRat(100_000, 1), UintToBigRat(uint64(100_000)))
}

func TestFloatToBigRat(t *testing.T) {
	assert.Equal(t, big.NewRat(125, 100), FloatToBigRat(float32(1.25)))
	assert.Equal(t, big.NewRat(25, 10), FloatToBigRat(float64(2.5)))
	assert.Equal(t, big.NewRat(12345, 100), FloatToBigRat(float64(123.45)))

	i := math.Inf(1)
	funcs.TryTo(
		func() {
			FloatToBigRat(i)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, i, fmt.Sprintf("%f", i), "*big.Rat"), e)
		},
	)
}

func TestBigIntToBigRat(t *testing.T) {
	assert.Equal(t, big.NewRat(1, 1), BigIntToBigRat(big.NewInt(1)))
	assert.Equal(t, big.NewRat(100_000, 1), BigIntToBigRat(big.NewInt(100_000)))
}

func TestBigFloatToBigRat(t *testing.T) {
	assert.Equal(t, big.NewRat(1, 1), BigFloatToBigRat(big.NewFloat(1)))
	assert.Equal(t, big.NewRat(125, 100), BigFloatToBigRat(big.NewFloat(125.0/100.0)))

	bf := big.NewFloat(math.Inf(1))
	funcs.TryTo(
		func() {
			BigFloatToBigRat(bf)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, bf, bf.String(), "*big.Rat"), e)
		},
	)
}

func TestStringToBigRat(t *testing.T) {
	assert.Equal(t, big.NewRat(1, 1), StringToBigRat("1/1"))
	assert.Equal(t, big.NewRat(125, 100), StringToBigRat("125/100"))
	assert.Equal(t, big.NewRat(125, 100), StringToBigRat("1.25"))

	funcs.TryTo(
		func() {
			StringToBigRat("1.25p")
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, "1.25p", "1.25p", "*big.Rat"), e)
		},
	)
}

func TestFloatStringToBigRat(t *testing.T) {
	assert.Equal(t, big.NewRat(125, 100), FloatStringToBigRat("1.25"))

	funcs.TryTo(
		func() {
			FloatStringToBigRat("125/100")
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, "125/100", "125/100", "*big.Rat"), e)
		},
	)
}
