package util

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/bantling/micro/go/funcs"
	"github.com/stretchr/testify/assert"
)

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

	i := math.Inf(0)
	funcs.TryTo(
		func() {
			FloatToBigRat(i)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(errMsg, i, i, "*big.Rat"), e)
		},
	)
}

func TestBigIntToBigRat(t *testing.T) {
	assert.Equal(t, big.NewRat(1, 1), BigIntToBigRat(big.NewInt(1)))
	assert.Equal(t, big.NewRat(100_000, 1), BigIntToBigRat(big.NewInt(100_000)))
}

func TestStringToBigRat(t *testing.T) {
	assert.Equal(t, big.NewRat(1, 1), StringToBigRat("1/1"))
	assert.Equal(t, big.NewRat(125, 100), StringToBigRat("125/100"))

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
