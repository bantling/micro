package math

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	gomath "math"
	"math/big"
	"testing"

  "github.com/bantling/micro/tuple"
	"github.com/stretchr/testify/assert"
)

func TestAbs_(t *testing.T) {
	{
		// ==== int
		i := -1
		assert.Nil(t, Abs(&i))
		assert.Equal(t, 1, i)

		i = 2
		assert.Nil(t, Abs(&i))
		assert.Equal(t, 2, i)

		i = gomath.MinInt
		assert.Equal(t, fmt.Errorf("Absolute value error for %d: there is no corresponding positive value in type int", i), Abs(&i))
	}

	{
		// ==== int8
		i := int8(-1)
		assert.Nil(t, Abs(&i))
		assert.Equal(t, int8(1), i)

		i = 2
		assert.Nil(t, Abs(&i))
		assert.Equal(t, int8(2), i)

		i = gomath.MinInt8
		assert.Equal(t, fmt.Errorf("Absolute value error for %d: there is no corresponding positive value in type int8", i), Abs(&i))
	}

	{
		// ==== int16
		i := int16(-1)
		assert.Nil(t, Abs(&i))
		assert.Equal(t, int16(1), i)

		i = 2
		assert.Nil(t, Abs(&i))
		assert.Equal(t, int16(2), i)

		i = gomath.MinInt16
		assert.Equal(t, fmt.Errorf("Absolute value error for %d: there is no corresponding positive value in type int16", i), Abs(&i))
	}

	{
		// ==== int32
		i := int32(-1)
		assert.Nil(t, Abs(&i))
		assert.Equal(t, int32(1), i)

		i = 2
		assert.Nil(t, Abs(&i))
		assert.Equal(t, int32(2), i)

		i = gomath.MinInt32
		assert.Equal(t, fmt.Errorf("Absolute value error for %d: there is no corresponding positive value in type int32", i), Abs(&i))
	}

	{
		// ==== int64
		i := int64(-1)
		assert.Nil(t, Abs(&i))
		assert.Equal(t, int64(1), i)

		i = 2
		assert.Nil(t, Abs(&i))
		assert.Equal(t, int64(2), i)

		i = gomath.MinInt64
		assert.Equal(t, fmt.Errorf("Absolute value error for %d: there is no corresponding positive value in type int64", i), Abs(&i))
	}

	{
		// ==== uint
		i := uint(1)
		assert.Nil(t, Abs(&i))
		assert.Equal(t, uint(1), i)
	}

	{
		// ==== uint8
		i := uint8(2)
		assert.Nil(t, Abs(&i))
		assert.Equal(t, uint8(2), i)
	}

	{
		// ==== uint16
		i := uint16(3)
		assert.Nil(t, Abs(&i))
		assert.Equal(t, uint16(3), i)
	}

	{
		// ==== uint32
		i := uint32(4)
		assert.Nil(t, Abs(&i))
		assert.Equal(t, uint32(4), i)
	}

	{
		// ==== uint64
		i := uint64(5)
		assert.Nil(t, Abs(&i))
		assert.Equal(t, uint64(5), i)
	}

	{
		// ==== float32
		f := float32(-1.25)
		assert.Nil(t, Abs(&f))
		assert.Equal(t, float32(1.25), f)

		f = 2.5
		assert.Nil(t, Abs(&f))
		assert.Equal(t, float32(2.5), f)
	}

	{
		// ==== float64
		f := -1.25
		assert.Nil(t, Abs(&f))
		assert.Equal(t, 1.25, f)

		f = 2.5
		assert.Nil(t, Abs(&f))
		assert.Equal(t, 2.5, f)
	}

	{
		// ==== *big.Int
		i := big.NewInt(-1)
		assert.Nil(t, Abs(&i))
		assert.Equal(t, big.NewInt(1), i)

		i = big.NewInt(2)
		assert.Nil(t, Abs(&i))
		assert.Equal(t, big.NewInt(2), i)
	}

	{
		// ==== *big.Float
		f := big.NewFloat(-1.25)
		assert.Nil(t, Abs(&f))
		assert.Equal(t, big.NewFloat(1.25), f)

		f = big.NewFloat(2.5)
		assert.Nil(t, Abs(&f))
		assert.Equal(t, big.NewFloat(2.5), f)
	}

	{
		// ==== *big.Rat
		r := big.NewRat(-125, 100)
		assert.Nil(t, Abs(&r))
		assert.Equal(t, big.NewRat(125, 100), r)

		r = big.NewRat(250, 10)
		assert.Nil(t, Abs(&r))
		assert.Equal(t, big.NewRat(250, 10), r)
	}
}

func TestAddInt_(t *testing.T) {
	{
		var i int = 2
		assert.Nil(t, AddInt(1, &i))
		assert.Equal(t, 3, i)

		assert.Equal(t, OverflowErr, AddInt(gomath.MaxInt, &i))
		assert.Equal(t, gomath.MinInt+2, i)

		assert.Equal(t, UnderflowErr, AddInt(-4, &i))
		assert.Equal(t, gomath.MaxInt-1, i)
	}

	{
		var i int8 = 2
		assert.Nil(t, AddInt(1, &i))
		assert.Equal(t, int8(3), i)

		assert.Equal(t, OverflowErr, AddInt(gomath.MaxInt8, &i))
		assert.Equal(t, int8(gomath.MinInt8+2), i)

		assert.Equal(t, UnderflowErr, AddInt(-4, &i))
		assert.Equal(t, int8(gomath.MaxInt8-1), i)
	}

	{
		var i int16 = 2
		assert.Nil(t, AddInt(1, &i))
		assert.Equal(t, int16(3), i)

		assert.Equal(t, OverflowErr, AddInt(gomath.MaxInt16, &i))
		assert.Equal(t, int16(gomath.MinInt16+2), i)

		assert.Equal(t, UnderflowErr, AddInt(-4, &i))
		assert.Equal(t, int16(gomath.MaxInt16-1), i)
	}

	{
		var i int32 = 2
		assert.Nil(t, AddInt(1, &i))
		assert.Equal(t, int32(3), i)

		assert.Equal(t, OverflowErr, AddInt(gomath.MaxInt32, &i))
		assert.Equal(t, int32(gomath.MinInt32+2), i)

		assert.Equal(t, UnderflowErr, AddInt(-4, &i))
		assert.Equal(t, int32(gomath.MaxInt32-1), i)
	}

	{
		var i int64 = 2
		assert.Nil(t, AddInt(1, &i))
		assert.Equal(t, int64(3), i)

		assert.Equal(t, OverflowErr, AddInt(gomath.MaxInt64, &i))
		assert.Equal(t, int64(gomath.MinInt64+2), i)

		assert.Equal(t, UnderflowErr, AddInt(-4, &i))
		assert.Equal(t, int64(gomath.MaxInt64-1), i)
	}
}

func TestAddUint_(t *testing.T) {
	{
		var i uint = 2
		assert.Nil(t, AddUint(1, &i))
		assert.Equal(t, uint(3), i)

		assert.Equal(t, OverflowErr, AddUint(gomath.MaxUint, &i))
		assert.Equal(t, uint(2), i)
	}

	{
		var i uint8 = 2
		assert.Nil(t, AddUint(1, &i))
		assert.Equal(t, uint8(3), i)

		assert.Equal(t, OverflowErr, AddUint(gomath.MaxUint8, &i))
		assert.Equal(t, uint8(2), i)
	}

	{
		var i uint16 = 2
		assert.Nil(t, AddUint(1, &i))
		assert.Equal(t, uint16(3), i)

		assert.Equal(t, OverflowErr, AddUint(gomath.MaxUint16, &i))
		assert.Equal(t, uint16(2), i)
	}

	{
		var i uint32 = 2
		assert.Nil(t, AddUint(1, &i))
		assert.Equal(t, uint32(3), i)

		assert.Equal(t, OverflowErr, AddUint(gomath.MaxUint32, &i))
		assert.Equal(t, uint32(2), i)
	}

	{
		var i uint64 = 2
		assert.Nil(t, AddUint(1, &i))
		assert.Equal(t, uint64(3), i)

		assert.Equal(t, OverflowErr, AddUint(gomath.MaxUint64, &i))
		assert.Equal(t, uint64(2), i)
	}
}

func TestSubInt_(t *testing.T) {
	{
		var i int = 2
		assert.Nil(t, SubInt(1, &i))
		assert.Equal(t, -1, i)

		i = -1
		assert.Equal(t, OverflowErr, SubInt(gomath.MaxInt, &i))
		assert.Equal(t, gomath.MinInt, i)

		i = 1
		assert.Equal(t, UnderflowErr, SubInt(gomath.MinInt, &i))
		assert.Equal(t, gomath.MaxInt, i)
	}

	{
		var i int8 = 2
		assert.Nil(t, SubInt(1, &i))
		assert.Equal(t, int8(-1), i)

		i = -1
		assert.Equal(t, OverflowErr, SubInt(gomath.MaxInt8, &i))
		assert.Equal(t, int8(gomath.MinInt8), i)

		i = 1
		assert.Equal(t, UnderflowErr, SubInt(gomath.MinInt8, &i))
		assert.Equal(t, int8(gomath.MaxInt8), i)
	}

	{
		var i int16 = 2
		assert.Nil(t, SubInt(1, &i))
		assert.Equal(t, int16(-1), i)

		i = -1
		assert.Equal(t, OverflowErr, SubInt(gomath.MaxInt16, &i))
		assert.Equal(t, int16(gomath.MinInt16), i)

		i = 1
		assert.Equal(t, UnderflowErr, SubInt(gomath.MinInt16, &i))
		assert.Equal(t, int16(gomath.MaxInt16), i)
	}

	{
		var i int32 = 2
		assert.Nil(t, SubInt(1, &i))
		assert.Equal(t, int32(-1), i)

		i = -1
		assert.Equal(t, OverflowErr, SubInt(gomath.MaxInt32, &i))
		assert.Equal(t, int32(gomath.MinInt32), i)

		i = 1
		assert.Equal(t, UnderflowErr, SubInt(gomath.MinInt32, &i))
		assert.Equal(t, int32(gomath.MaxInt32), i)
	}

	{
		var i int64 = 2
		assert.Nil(t, SubInt(1, &i))
		assert.Equal(t, int64(-1), i)

		i = -1
		assert.Equal(t, OverflowErr, SubInt(gomath.MaxInt64, &i))
		assert.Equal(t, int64(gomath.MinInt64), i)

		i = 1
		assert.Equal(t, UnderflowErr, SubInt(gomath.MinInt64, &i))
		assert.Equal(t, int64(gomath.MaxInt64), i)
	}
}

func TestSubUint_(t *testing.T) {
	{
		var i uint = 1
		assert.Nil(t, SubUint(4, &i))
		assert.Equal(t, uint(3), i)

		i = 1
		assert.Equal(t, UnderflowErr, SubUint(0, &i))
		assert.Equal(t, uint(gomath.MaxUint), i)
	}

	{
		var i uint8 = 1
		assert.Nil(t, SubUint(4, &i))
		assert.Equal(t, uint8(3), i)

		i = 1
		assert.Equal(t, UnderflowErr, SubUint(0, &i))
		assert.Equal(t, uint8(gomath.MaxUint8), i)
	}

	{
		var i uint16 = 1
		assert.Nil(t, SubUint(4, &i))
		assert.Equal(t, uint16(3), i)

		i = 1
		assert.Equal(t, UnderflowErr, SubUint(0, &i))
		assert.Equal(t, uint16(gomath.MaxUint16), i)
	}

	{
		var i uint32 = 1
		assert.Nil(t, SubUint(4, &i))
		assert.Equal(t, uint32(3), i)

		i = 1
		assert.Equal(t, UnderflowErr, SubUint(0, &i))
		assert.Equal(t, uint32(gomath.MaxUint32), i)
	}

	{
		var i uint64 = 1
		assert.Nil(t, SubUint(4, &i))
		assert.Equal(t, uint64(3), i)

		i = 1
		assert.Equal(t, UnderflowErr, SubUint(0, &i))
		assert.Equal(t, uint64(gomath.MaxUint64), i)
	}
}

func TestMul_(t *testing.T) {
	{
		var i int = 2
		assert.Nil(t, Mul(5, &i))
		assert.Equal(t, 10, i)

		assert.Equal(t, OverflowErr, Mul(gomath.MaxInt, &i))
		assert.Equal(t, 10, i)

		assert.Equal(t, UnderflowErr, Mul(gomath.MinInt, &i))
		assert.Equal(t, 10, i)
	}

	{
		var i int8 = 2
		assert.Nil(t, Mul(5, &i))
		assert.Equal(t, int8(10), i)

		assert.Equal(t, OverflowErr, Mul(gomath.MaxInt8, &i))
		assert.Equal(t, int8(10), i)

		assert.Equal(t, UnderflowErr, Mul(gomath.MinInt8, &i))
		assert.Equal(t, int8(10), i)
	}

	{
		var i int16 = 2
		assert.Nil(t, Mul(5, &i))
		assert.Equal(t, int16(10), i)

		assert.Equal(t, OverflowErr, Mul(gomath.MaxInt16, &i))
		assert.Equal(t, int16(10), i)

		assert.Equal(t, UnderflowErr, Mul(gomath.MinInt16, &i))
		assert.Equal(t, int16(10), i)
	}

	{
		var i int32 = 2
		assert.Nil(t, Mul(5, &i))
		assert.Equal(t, int32(10), i)

		assert.Equal(t, OverflowErr, Mul(gomath.MaxInt32, &i))
		assert.Equal(t, int32(10), i)

		assert.Equal(t, UnderflowErr, Mul(gomath.MinInt32, &i))
		assert.Equal(t, int32(10), i)
	}

	{
		var i int64 = 2
		assert.Nil(t, Mul(5, &i))
		assert.Equal(t, int64(10), i)

		assert.Equal(t, OverflowErr, Mul(gomath.MaxInt64, &i))
		assert.Equal(t, int64(10), i)

		assert.Equal(t, UnderflowErr, Mul(gomath.MinInt64, &i))
		assert.Equal(t, int64(10), i)
	}

	{
		var i uint = 2
		assert.Nil(t, Mul(5, &i))
		assert.Equal(t, uint(10), i)

		assert.Equal(t, OverflowErr, Mul(gomath.MaxUint, &i))
		assert.Equal(t, uint(10), i)
	}

	{
		var i uint8 = 2
		assert.Nil(t, Mul(5, &i))
		assert.Equal(t, uint8(10), i)

		assert.Equal(t, OverflowErr, Mul(gomath.MaxUint8, &i))
		assert.Equal(t, uint8(10), i)
	}

	{
		var i uint16 = 2
		assert.Nil(t, Mul(5, &i))
		assert.Equal(t, uint16(10), i)

		assert.Equal(t, OverflowErr, Mul(gomath.MaxUint16, &i))
		assert.Equal(t, uint16(10), i)
	}

	{
		var i uint32 = 2
		assert.Nil(t, Mul(5, &i))
		assert.Equal(t, uint32(10), i)

		assert.Equal(t, OverflowErr, Mul(gomath.MaxUint32, &i))
		assert.Equal(t, uint32(10), i)
	}

	{
		var i uint64 = 2
		assert.Nil(t, Mul(5, &i))
		assert.Equal(t, uint64(10), i)

		assert.Equal(t, OverflowErr, Mul(gomath.MaxUint64, &i))
		assert.Equal(t, uint64(10), i)
	}
}

func TestDiv_(t *testing.T) {
	{
		// ==== float32
		var (
			de, dv, q float32
		)

		// 18 / 4 = 4.5
		de, dv = 18.0, 4.0
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, float32(4.5), q)

		// 17 / 4 = 4.25
		de, dv = 17.0, 4.0
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, float32(4.25), q)

		// non-zero / zero = infinity
		de, dv = 1.0, 0.0
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, float32(gomath.Inf(1)), q)

		// 0 / 0 = NaN
		de, dv = 0.0, 0.0
		assert.Nil(t, Div(de, dv, &q))
		assert.True(t, gomath.IsNaN(float64(q)))

		// infinity / non-infinity = infinity
		de, dv = float32(gomath.Inf(1)), 0.0
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, float32(gomath.Inf(1)), q)

		// infinity / infinity = NaN
		de, dv = float32(gomath.Inf(1)), float32(gomath.Inf(1))
		assert.Nil(t, Div(de, dv, &q))
		assert.True(t, gomath.IsNaN(float64(q)))
	}

	{
		// ==== float64
		var (
			de, dv, q float64
		)

		// 18 / 4 = 4.5
		de, dv = 18.0, 4.0
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, 4.5, q)

		// 17 / 4 = 4.25
		de, dv = 17.0, 4.0
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, 4.25, q)

		// non-zero / zero = infinity
		de, dv = 1.0, 0.0
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, gomath.Inf(1), q)

		// 0 / 0 = NaN
		de, dv = 0.0, 0.0
		assert.Nil(t, Div(de, dv, &q))
		assert.True(t, gomath.IsNaN(q))

		// infinity / non-infinity = infinity
		de, dv = gomath.Inf(1), 0.0
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, gomath.Inf(1), q)

		// infinity / infinity = NaN
		de, dv = gomath.Inf(1), gomath.Inf(1)
		assert.Nil(t, Div(de, dv, &q))
		assert.True(t, gomath.IsNaN(q))
	}

	{
		// ==== int
		var (
			de, dv, q int
		)

		// 18 / 4 =  4 r 2 = 5
		de, dv = 18, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, 5, q)

		// -18 / 4 = -4 r -2 = -5
		de, dv = -18, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, -5, q)

		// 18 / -4 = -4 r 2 = -5
		de, dv = 18, -4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, -5, q)

		// -18 / -4 = 4 r -2 = 5
		de, dv = -18, -4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, 5, q)

		// 17 / 4 = 4 r 1 = 4
		de, dv = 17, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, 4, q)

		// -17 / 4 = -4 r -1 = -4
		de, dv = -17, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, -4, q)

		// 17 / -4 = -4 r 1 = -4
		de, dv = 17, -4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, -4, q)

		// -17 / -4 = 4 r -1 = 4
		de, dv = -17, -4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, 4, q)

		// 18 / 5 = 3 r 3 = 4
		de, dv = 18, 5
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, 4, q)

		// -18 / 5 = -3 r -3 = -4
		de, dv = -18, 5
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, -4, q)

		// 18 / -5 = -3 r 3 = -4
		de, dv = 18, -5
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, -4, q)

		// -18 / -5 = 3 r -3 = 4
		de, dv = -18, -5
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, 4, q)

		// 17 / 5 = 3 r 2 = 3
		de, dv = 17, 5
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, 3, q)

		// -17 / 5 = -3 r -2 = -3
		de, dv = -17, 5
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, -3, q)

		// 17 / -5 = -3 r 2 = -3
		de, dv = 17, -5
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, -3, q)

		// -17 / -5 = 3 r -3 = 3
		de, dv = -17, -5
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, 3, q)

		// division by zero
		de, dv, q = 1, 0, 2
		assert.Equal(t, DivByZeroErr, Div(de, dv, &q))
		assert.Equal(t, 2, q)
	}

	{
		// ==== int8

		var (
			de, dv, q int8
		)

		// 18 / 4 =  4 r 2 = 5
		de, dv = 18, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int8(5), q)

		// 17 / 5 = 3 r 2 = 3
		de, dv = -17, 5
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int8(-3), q)

		// division by zero
		de, dv, q = 1, 0, 2
		assert.Equal(t, DivByZeroErr, Div(de, dv, &q))
		assert.Equal(t, int8(2), q)
	}

	{
		// ==== int16

		var (
			de, dv, q int16
		)

		// 18 / 4 =  4 r 2 = 5
		de, dv = 18, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int16(5), q)

		// 17 / 5 = 3 r 2 = 3
		de, dv = -17, 5
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int16(-3), q)

		// division by zero
		de, dv, q = 1, 0, 2
		assert.Equal(t, DivByZeroErr, Div(de, dv, &q))
		assert.Equal(t, int16(2), q)
	}

	{
		// ==== int32

		var (
			de, dv, q int32
		)

		// 18 / 4 =  4 r 2 = 5
		de, dv = 18, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int32(5), q)

		// 17 / 5 = 3 r 2 = 3
		de, dv = -17, 5
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int32(-3), q)

		// division by zero
		de, dv, q = 1, 0, 2
		assert.Equal(t, DivByZeroErr, Div(de, dv, &q))
		assert.Equal(t, int32(2), q)
	}

	{
		// ==== int64

		var (
			de, dv, q int64
		)

		// 18 / 4 =  4 r 2 = 5
		de, dv = 18, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(5), q)

		// 17 / 5 = 3 r 2 = 3
		de, dv = -17, 5
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(-3), q)

		// division by zero
		de, dv, q = 1, 0, 2
		assert.Equal(t, DivByZeroErr, Div(de, dv, &q))
		assert.Equal(t, int64(2), q)
	}

	{
		// ==== uint
		var (
			de, dv, q uint
		)

		// 18 / 4 =  4 r 2 = 5
		de, dv = 18, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, uint(5), q)

		// 17 / 4 = 4 r 1 = 4
		de, dv = 17, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, uint(4), q)

		// 18 / 5 = 3 r 3 = 4
		de, dv = 18, 5
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, uint(4), q)

		// 17 / 5 = 3 r 2 = 3
		de, dv = 17, 5
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, uint(3), q)

		// division by zero
		de, dv, q = 1, 0, 2
		assert.Equal(t, DivByZeroErr, Div(de, dv, &q))
		assert.Equal(t, uint(2), q)
	}

	{
		// ==== uint8
		var (
			de, dv, q uint8
		)

		// 18 / 4 =  4 r 2 = 5
		de, dv = 18, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, uint8(5), q)

		// 17 / 4 = 4 r 1 = 4
		de, dv = 17, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, uint8(4), q)

		// division by zero
		de, dv, q = 1, 0, 2
		assert.Equal(t, DivByZeroErr, Div(de, dv, &q))
		assert.Equal(t, uint8(2), q)
	}

	{
		// ==== uint16
		var (
			de, dv, q uint16
		)

		// 18 / 4 =  4 r 2 = 5
		de, dv = 18, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, uint16(5), q)

		// 17 / 4 = 4 r 1 = 4
		de, dv = 17, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, uint16(4), q)

		// division by zero
		de, dv, q = 1, 0, 2
		assert.Equal(t, DivByZeroErr, Div(de, dv, &q))
		assert.Equal(t, uint16(2), q)
	}

	{
		// ==== uint32
		var (
			de, dv, q uint32
		)

		// 18 / 4 =  4 r 2 = 5
		de, dv = 18, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, uint32(5), q)

		// 17 / 4 = 4 r 1 = 4
		de, dv = 17, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, uint32(4), q)

		// division by zero
		de, dv, q = 1, 0, 2
		assert.Equal(t, DivByZeroErr, Div(de, dv, &q))
		assert.Equal(t, uint32(2), q)
	}

	{
		// ==== uint64
		var (
			de, dv, q uint64
		)

		// 18 / 4 =  4 r 2 = 5
		de, dv = 18, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, uint64(5), q)

		// 17 / 4 = 4 r 1 = 4
		de, dv = 17, 4
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, uint64(4), q)

		// division by zero
		de, dv, q = 1, 0, 2
		assert.Equal(t, DivByZeroErr, Div(de, dv, &q))
		assert.Equal(t, uint64(2), q)
	}

	{
		// ==== *big.Int
		var (
			de, dv *big.Int
			q      = big.NewInt(0)
		)

		// 18 / 4 = 4 r 2 = 5
		de, dv = big.NewInt(18), big.NewInt(4)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(18), de.Int64())
		assert.Equal(t, int64(4), dv.Int64())
		assert.Equal(t, int64(5), q.Int64())

		// -18 / 4 = -4 r -2 = -5
		de, dv = big.NewInt(-18), big.NewInt(4)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(-18), de.Int64())
		assert.Equal(t, int64(4), dv.Int64())
		assert.Equal(t, int64(-5), q.Int64())

		// 18 / -4 = -4 r 2 = -5
		de, dv = big.NewInt(18), big.NewInt(-4)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(18), de.Int64())
		assert.Equal(t, int64(-4), dv.Int64())
		assert.Equal(t, int64(-5), q.Int64())

		// -18 / -4 = 4 r -2 = 5
		de, dv = big.NewInt(-18), big.NewInt(-4)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(-18), de.Int64())
		assert.Equal(t, int64(-4), dv.Int64())
		assert.Equal(t, int64(5), q.Int64())

		// 17 / 4 = 4 r 1 = 4
		de, dv = big.NewInt(17), big.NewInt(4)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(17), de.Int64())
		assert.Equal(t, int64(4), dv.Int64())
		assert.Equal(t, int64(4), q.Int64())

		// -17 / 4 = -4 r -1 = -4
		de, dv = big.NewInt(-17), big.NewInt(4)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(-17), de.Int64())
		assert.Equal(t, int64(4), dv.Int64())
		assert.Equal(t, int64(-4), q.Int64())

		// 17 / -4 = -4 r 1 = -4
		de, dv = big.NewInt(17), big.NewInt(-4)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(17), de.Int64())
		assert.Equal(t, int64(-4), dv.Int64())
		assert.Equal(t, int64(-4), q.Int64())

		// -17 / -4 = 4 r -1 = 4
		de, dv = big.NewInt(-17), big.NewInt(-4)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(-17), de.Int64())
		assert.Equal(t, int64(-4), dv.Int64())
		assert.Equal(t, int64(4), q.Int64())

		// 18 / 5 = 3 r 3 = 4
		de, dv = big.NewInt(18), big.NewInt(5)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(18), de.Int64())
		assert.Equal(t, int64(5), dv.Int64())
		assert.Equal(t, int64(4), q.Int64())

		// -18 / 5 = -3 r -3 = -4
		de, dv = big.NewInt(-18), big.NewInt(5)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(-18), de.Int64())
		assert.Equal(t, int64(5), dv.Int64())
		assert.Equal(t, int64(-4), q.Int64())

		// 18 / -5 = -3 r 3 = -4
		de, dv = big.NewInt(18), big.NewInt(-5)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(18), de.Int64())
		assert.Equal(t, int64(-5), dv.Int64())
		assert.Equal(t, int64(-4), q.Int64())

		// -18 / -5 = 3 r -3 = 4
		de, dv = big.NewInt(-18), big.NewInt(-5)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(-18), de.Int64())
		assert.Equal(t, int64(-5), dv.Int64())
		assert.Equal(t, int64(4), q.Int64())

		// 17 / 5 = 3 r 2 = 3
		de, dv = big.NewInt(17), big.NewInt(5)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(17), de.Int64())
		assert.Equal(t, int64(5), dv.Int64())
		assert.Equal(t, int64(3), q.Int64())

		// -17 / 5 = -3 r -2 = -3
		de, dv = big.NewInt(-17), big.NewInt(5)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(-17), de.Int64())
		assert.Equal(t, int64(5), dv.Int64())
		assert.Equal(t, int64(-3), q.Int64())

		// 17 / -5 = -3 r 2 = -3
		de, dv = big.NewInt(17), big.NewInt(-5)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(17), de.Int64())
		assert.Equal(t, int64(-5), dv.Int64())
		assert.Equal(t, int64(-3), q.Int64())

		// -17 / -5 = 3 r -2 = 3
		de, dv = big.NewInt(-17), big.NewInt(-5)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(-17), de.Int64())
		assert.Equal(t, int64(-5), dv.Int64())
		assert.Equal(t, int64(3), q.Int64())

		// division by zero
		de, dv = big.NewInt(1), big.NewInt(0)
		assert.Equal(t, DivByZeroErr, Div(de, dv, &q))
		assert.Equal(t, int64(1), de.Int64())
		assert.Equal(t, int64(0), dv.Int64())

		// nil quotient gets modified
		de, dv, q = big.NewInt(1), big.NewInt(1), nil
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, int64(1), de.Int64())
		assert.Equal(t, int64(1), dv.Int64())
	}

	{
		// ==== *big.Float
		var (
			de, dv *big.Float
			q      = big.NewFloat(0.0)
		)

		// 18 / 4 = 4.5
		de, dv = big.NewFloat(18.0), big.NewFloat(4.0)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, big.NewFloat(18.0), de)
		assert.Equal(t, big.NewFloat(4.0), dv)
		assert.Equal(t, big.NewFloat(4.5), q)

		// 17 / 4 = 4.25
		de, dv = big.NewFloat(17.0), big.NewFloat(4.0)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, big.NewFloat(17.0), de)
		assert.Equal(t, big.NewFloat(4.0), dv)
		assert.Equal(t, big.NewFloat(4.25), q)

		// non-zero / zero = infinity
		de, dv = big.NewFloat(1.0), big.NewFloat(0.0)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, big.NewFloat(1.0), de)
		assert.Equal(t, big.NewFloat(0.0), dv)
		assert.True(t, q.IsInf())

		// zero / zero = NaN
		de, dv, q = big.NewFloat(0.0), big.NewFloat(0.0), big.NewFloat(1.0)
		assert.Equal(t, big.ErrNaN{}, Div(de, dv, &q))
		assert.Equal(t, big.NewFloat(0.0), de)
		assert.Equal(t, big.NewFloat(0.0), dv)
		assert.Equal(t, big.NewFloat(1.0), q)

		// infinity / non-infinity = infinity
		de, dv = big.NewFloat(gomath.Inf(1)), big.NewFloat(0.0)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, big.NewFloat(gomath.Inf(1)), de)
		assert.Equal(t, big.NewFloat(0.0), dv)
		assert.True(t, q.IsInf())

		// infinity / infinity = NaN
		de, dv, q = big.NewFloat(gomath.Inf(1)), big.NewFloat(gomath.Inf(1)), big.NewFloat(1.0)
		assert.Equal(t, big.ErrNaN{}, Div(de, dv, &q))
		assert.Equal(t, big.NewFloat(gomath.Inf(1)), de)
		assert.Equal(t, big.NewFloat(gomath.Inf(1)), dv)
		assert.Equal(t, big.NewFloat(1.0), q)

		// nil quotient gets modified
		de, dv, q = big.NewFloat(1), big.NewFloat(1), nil
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, big.NewFloat(1.0), de)
		assert.Equal(t, big.NewFloat(1.0), dv)
		assert.Equal(t, big.NewFloat(1.0), q)
	}

	{
		// ==== *big.Rat
		var (
			de, dv *big.Rat
			q      = big.NewRat(0, 1)
		)

		// 18 / 4 = 4.5
		de, dv = big.NewRat(18, 1), big.NewRat(4, 1)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, big.NewRat(18, 1), de)
		assert.Equal(t, big.NewRat(4, 1), dv)
		assert.Equal(t, big.NewRat(45, 10), q)

		// 17 / 4 = 4.25
		de, dv = big.NewRat(17, 1), big.NewRat(4, 1)
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, big.NewRat(17, 1), de)
		assert.Equal(t, big.NewRat(4, 1), dv)
		assert.Equal(t, big.NewRat(425, 100), q)

		// division by zero
		de, dv, q = big.NewRat(1, 1), big.NewRat(0, 1), big.NewRat(2, 1)
		assert.Equal(t, DivByZeroErr, Div(de, dv, &q))
		assert.Equal(t, big.NewRat(1, 1), de)
		assert.Equal(t, big.NewRat(0, 1), dv)
		assert.Equal(t, big.NewRat(2, 1), q)

		// nil quotient gets modified
		de, dv, q = big.NewRat(1, 1), big.NewRat(1, 1), nil
		assert.Nil(t, Div(de, dv, &q))
		assert.Equal(t, big.NewRat(1, 1), de)
		assert.Equal(t, big.NewRat(1, 1), dv)
		assert.Equal(t, big.NewRat(1, 1), q)
	}
}

func TestDivBigOps_(t *testing.T) {
	de, dv := big.NewInt(18), big.NewInt(4)
	var q *big.Int
	assert.Nil(t, DivBigOps(de, dv, &q))
	assert.Equal(t, big.NewInt(5), q)

	dv = big.NewInt(0)
	assert.Equal(t, DivByZeroErr, DivBigOps(de, dv, &q))
	assert.Equal(t, big.NewInt(5), q)
}

type cmp int

func (t cmp) Cmp(o cmp) int {
	if t < o {
		return -1
	}

	if t == o {
		return 0
	}

	return 1
}

func TestMinMax_(t *testing.T) {
	// Ordered = int
	func() {
		i, j := 1, 2
		assert.Equal(t, i, MinOrdered(i, j))
		assert.Equal(t, j, MaxOrdered(i, j))

		i, j = 1, 1
		assert.Equal(t, i, MinOrdered(i, j))
		assert.Equal(t, j, MaxOrdered(i, j))

		i, j = 2, 1
		assert.Equal(t, j, MinOrdered(i, j))
		assert.Equal(t, i, MaxOrdered(i, j))
	}()

	// Ordered = string
	func() {
		i, j := "1", "2"
		assert.Equal(t, i, MinOrdered(i, j))
		assert.Equal(t, j, MaxOrdered(i, j))

		i, j = "1", "1"
		assert.Equal(t, i, MinOrdered(i, j))
		assert.Equal(t, j, MaxOrdered(i, j))

		i, j = "2", "1"
		assert.Equal(t, j, MinOrdered(i, j))
		assert.Equal(t, i, MaxOrdered(i, j))
	}()

	// Complex
	func() {
		i, j := 1+2i, 2+3i
		assert.Equal(t, i, MinComplex(i, j))
		assert.Equal(t, j, MaxComplex(i, j))

		i, j = 1+2i, 1+2i
		assert.Equal(t, i, MinComplex(i, j))
		assert.Equal(t, j, MaxComplex(i, j))

		i, j = 2+3i, 1+2i
		assert.Equal(t, j, MinComplex(i, j))
		assert.Equal(t, i, MaxComplex(i, j))
	}()

	// Cmp
	func() {
		i, j := cmp(1), cmp(2)
		assert.Equal(t, i, MinCmp(i, j))
		assert.Equal(t, j, MaxCmp(i, j))

		i, j = cmp(1), cmp(1)
		assert.Equal(t, i, MinCmp(i, j))
		assert.Equal(t, j, MaxCmp(i, j))

		i, j = cmp(2), cmp(1)
		assert.Equal(t, j, MinCmp(i, j))
		assert.Equal(t, i, MaxCmp(i, j))
	}()
}

func TestOfDecimal_(t *testing.T) {
  assert.Equal(t, tuple.Of2(Decimal{scale: 2, value: 100}, error(nil)), tuple.Of2(OfDecimal(100)))
  assert.Equal(t, tuple.Of2(Decimal{scale: 3, value: -1001}, error(nil)), tuple.Of2(OfDecimal(-1001, 3)))

  assert.Equal(t, Decimal{scale: 2, value: 100}, MustDecimal(100))
  assert.Equal(t, Decimal{scale: 3, value: -1001}, MustDecimal(-1001, 3))

  assert.Equal(
    t,
    tuple.Of2(Decimal{}, fmt.Errorf("The Decimal scale 19 is too large: the value must be <= 18")),
    tuple.Of2(OfDecimal(0, 19)),
  )
  assert.Equal(
    t,
    tuple.Of2(Decimal{}, fmt.Errorf("The Decimal value 1234567890123456789 is too large: the value must be <= 999_999_999_999_999_999")),
    tuple.Of2(OfDecimal(1234567890123456789)),
  )
  assert.Equal(
    t,
    tuple.Of2(Decimal{}, fmt.Errorf("The Decimal value -1234567890123456789 is too small: the value must be >= -999_999_999_999_999_999")),
    tuple.Of2(OfDecimal(-1234567890123456789)),
  )
}

func TestSign_(t *testing.T) {
  d := MustDecimal(0)
  assert.Equal(t, 0, d.Sign())

  d = MustDecimal(0, 5)
  assert.Equal(t, 0, d.Sign())

  d = MustDecimal(1)
  assert.Equal(t, 1, d.Sign())

  d = MustDecimal(-1)
  assert.Equal(t, -1, d.Sign())
}

func TestDecimalString_(t *testing.T) {
  assert.Equal(t, "123", MustDecimal(123, 0).String())
  assert.Equal(t, "-123", MustDecimal(-123, 0).String())
  assert.Equal(t, "12.3", MustDecimal(123, 1).String())
  assert.Equal(t, "1.23", MustDecimal(123, 2).String())
  assert.Equal(t, "0.123", MustDecimal(123, 3).String())
  assert.Equal(t, "0.0123", MustDecimal(123, 4).String())
  assert.Equal(t, "0.00123", MustDecimal(123, 5).String())
  assert.Equal(t, "-0.00123", MustDecimal(-123, 5).String())
}
