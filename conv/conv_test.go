package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"
	goreflect "reflect"
	"testing"

	"github.com/bantling/micro/funcs"
	"github.com/stretchr/testify/assert"
)

func TestTo_(t *testing.T) {
	// Nil target error
	{
		assert.Equal(t, fmt.Errorf("The target value of type *int cannot be nil"), To(0, (*int)(nil)))
	}

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

		{
			// big to big of same type, input is nil
			var bi *big.Int
			var bo *big.Int = big.NewInt(0)
			assert.Nil(t, To(bi, &bo))
			assert.Nil(t, bo)

			// big to big of same type, and are same pointer
			bi = big.NewInt(1)
			bo = bi
			assert.Nil(t, To(bi, &bo))
			assert.False(t, bi == bo)
			assert.Equal(t, big.NewInt(1), bi)
			assert.Equal(t, big.NewInt(1), bo)

			// big to big of same type, output is nil
			bi = big.NewInt(1)
			bo = nil
			assert.Nil(t, To(bi, &bo))
			assert.False(t, bi == bo)
			assert.Equal(t, big.NewInt(1), bi)
			assert.Equal(t, big.NewInt(1), bo)

			bi = big.NewInt(math.MaxInt64)
			bi = bi.Mul(bi, big.NewInt(2))
			assert.Equal(t, fmt.Errorf("The *big.Int value of 18446744073709551614 cannot be converted to int64"), To(bi, &i))

			funcs.TryTo(
				func() {
					MustTo(bi, &i)
					assert.Fail(t, "Never execute")
				},
				func(e any) {
					assert.Equal(t, "The *big.Int value of 18446744073709551614 cannot be converted to int64", e.(error).Error())
				},
			)
		}

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
		bisrc := big.NewInt(1)
		assert.Nil(t, To(bisrc, &bi))
		assert.False(t, bisrc == bi)
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

		bfsrc := big.NewFloat(1.25)
		assert.Nil(t, To(bfsrc, &bf))
		assert.False(t, bfsrc == bf)
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

		brsrc := big.NewRat(25, 10)
		assert.Nil(t, To(brsrc, &br))
		assert.False(t, brsrc == br)
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

	// source type = target type (int -> int)
	{
		var o int
		assert.Nil(t, To(1, &o))
		assert.Equal(t, 1, o)
	}

	{
		// byte to rune, which is really uint8 to int32
		// it is not a subtype, reflection sees uint8 and int32
		var r rune
		assert.Nil(t, To(byte('A'), &r))
		assert.Equal(t, 'A', r)
	}

	// Subtypes where no conversion exists, base types are the same
	{
		type foo int
		type bar int
		var b bar
		assert.Nil(t, To(foo(1), &b))
		assert.Equal(t, bar(1), b)
	}

	// Subtypes where no conversion exists, base types are different
	{
		type foo uint
		type bar int
		var b bar
		assert.Nil(t, To(foo(1), &b))
		assert.Equal(t, bar(1), b)
	}
}

func TestAnyTo_(t *testing.T) {
	var o int

	// int
	assert.Nil(t, AnyTo(1, &o))
	assert.Equal(t, 1, o)

	// int8
	assert.Nil(t, AnyTo(int8(2), &o))
	assert.Equal(t, 2, o)

	// int16
	assert.Nil(t, AnyTo(int16(3), &o))
	assert.Equal(t, 3, o)

	// int32
	assert.Nil(t, AnyTo(int32(4), &o))
	assert.Equal(t, 4, o)

	// int64
	assert.Nil(t, AnyTo(int64(5), &o))
	assert.Equal(t, 5, o)

	// uint
	assert.Nil(t, AnyTo(uint(6), &o))
	assert.Equal(t, 6, o)

	// uint8
	assert.Nil(t, AnyTo(uint8(7), &o))
	assert.Equal(t, 7, o)

	// uint16
	assert.Nil(t, AnyTo(uint16(8), &o))
	assert.Equal(t, 8, o)

	// uint32
	assert.Nil(t, AnyTo(uint32(9), &o))
	assert.Equal(t, 9, o)

	// uint64
	assert.Nil(t, AnyTo(uint64(10), &o))
	assert.Equal(t, 10, o)

	// float32
	assert.Nil(t, AnyTo(float32(11), &o))
	assert.Equal(t, 11, o)

	// float64
	assert.Nil(t, AnyTo(float64(12), &o))
	assert.Equal(t, 12, o)

	// string
	assert.Nil(t, AnyTo("13", &o))
	assert.Equal(t, 13, o)

	// *big.Int
	assert.Nil(t, AnyTo(big.NewInt(14), &o))
	assert.Equal(t, 14, o)

	// *big.Float
	assert.Nil(t, AnyTo(big.NewFloat(15), &o))
	assert.Equal(t, 15, o)

	// *big.Rat
	assert.Nil(t, AnyTo(big.NewRat(16, 1), &o))
	assert.Equal(t, 16, o)

	// Any other value is an error
	assert.Equal(
		t,
		fmt.Errorf("AnyTo cannot convert the input type struct { Name string }"),
		AnyTo(struct{ Name string }{Name: "Foo"}, &o),
	)
}

func TestToBigOps_(t *testing.T) {
	// Target cannot be nil
	assert.Equal(t, fmt.Errorf("The target value of type **big.Int cannot be nil"), ToBigOps(1, (**big.Int)(nil)))

	{
		// int to *big.Int
		var bi *big.Int
		assert.Nil(t, ToBigOps(1, &bi))
		assert.Equal(t, big.NewInt(1), bi)

		// copy nil *big.Int to non-nil *big.Int
		bi = nil
		var nbi *big.Int = big.NewInt(2)
		assert.Nil(t, ToBigOps(bi, &nbi))
		assert.Nil(t, nbi)

		// copy *big.Int to *big.Int with same pointer
		bi = big.NewInt(2)
		nbi = bi
		assert.Nil(t, ToBigOps(bi, &nbi))
		assert.False(t, nbi == bi)
		assert.Equal(t, big.NewInt(2), nbi)

		// copy *big.Int to nil *big.Int
		bi = big.NewInt(3)
		nbi = nil
		assert.Nil(t, ToBigOps(bi, &nbi))
		assert.False(t, nbi == bi)
		assert.Equal(t, big.NewInt(3), nbi)

		// byte to *big.Int, which is really uint8 to *big.Int
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

	{
		var br *big.Rat
		assert.Nil(t, ToBigOps(3, &br))
		assert.Equal(t, big.NewRat(3, 1), br)

		assert.Nil(t, ToBigOps(br, &br))
		assert.Equal(t, big.NewRat(3, 1), br)
	}

	funcs.TryTo(
		func() {
			var br *big.Int
			MustToBigOps("a", &br)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, fmt.Sprintf("The string value of a cannot be converted to *big.Int"), e.(error).Error())
		},
	)
}

type Foo struct{}

func TestReflectTo_(t *testing.T) {
	{
		// int to string
		var str string
		assert.Nil(t, ReflectTo(goreflect.ValueOf(1), goreflect.ValueOf(&str)))
		assert.Equal(t, "1", str)
	}

	{
		// int to struct string field
		type Foo struct {
			Str string
		}

		var f = Foo{}
		assert.Nil(t, ReflectTo(goreflect.ValueOf(2), goreflect.ValueOf(&f).Elem().FieldByName("Str").Addr()))
		assert.Equal(t, "2", f.Str)

		// cannot convert from Invalid to *Foo (failure in ReflectTo)
		assert.Equal(t, errReflectToInvalidSrc, ReflectTo(goreflect.Value{}, goreflect.ValueOf(&f)))

		// cannot convert from int to Invalid (failure in ReflectTo)
		assert.Equal(t, errReflectToInvalidTgt, ReflectTo(goreflect.ValueOf(4), goreflect.Value{}))

		// cannot convert to a non-pointer
		assert.Equal(t, fmt.Errorf("ReflectTo target must be be a pointer"), ReflectTo(goreflect.ValueOf(5), goreflect.ValueOf(6)))

		// cannot convert to a nil pointer
		assert.Equal(t, fmt.Errorf("The target value of type *int cannot be nil"), ReflectTo(goreflect.ValueOf(7), goreflect.ValueOf((*int)(nil))))

		// cannot convert to big value
		assert.Equal(t, fmt.Errorf("ReflectTo target must be be a pointer"), ReflectTo(goreflect.ValueOf(8), goreflect.ValueOf(big.Int{})))

		// cannot convert to big pointer
		assert.Equal(t, fmt.Errorf("The target value of type *big.Int is invalid: big types have to be a **"), ReflectTo(goreflect.ValueOf(9), goreflect.ValueOf(big.NewInt(9))))

		// cannot convert a src type with multiple pointers (failure in LookupConversion)
		assert.Equal(t, fmt.Errorf("There is no conversion function from **int to conv.Foo"), ReflectTo(goreflect.ValueOf((**int)(nil)), goreflect.ValueOf(&f)))

		// non-existent conversion from int to struct (LookupConversion returns a nil conversion func)
		assert.Equal(t, fmt.Errorf("There is no conversion function from int to conv.Foo"), ReflectTo(goreflect.ValueOf(3), goreflect.ValueOf(&f)))
	}
}

func TestMustReflectTo_(t *testing.T) {
	var s string
	MustReflectTo(goreflect.ValueOf(1), goreflect.ValueOf(&s))
	assert.Equal(t, "1", s)

	funcs.TryTo(
		func() {
			var i int
			MustReflectTo(goreflect.ValueOf("a"), goreflect.ValueOf(&i))
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, fmt.Sprintf("The string value of a cannot be converted to int64"), e.(error).Error())
		},
	)
}
