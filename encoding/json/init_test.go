package json

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"
	"testing"

	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/union"
	"github.com/stretchr/testify/assert"
)

func TestInit_(t *testing.T) {
	var (
		v  Value
		z  = MustNumberToValue(0)
		ma map[string]any
		sa []any
		zs = NumberString("0")
	)

	// Object
	v = Value{}
	ma = map[string]any{"foo": 0}
	assert.Nil(t, conv.To(ma, &v))
	assert.Equal(t, Value{typ: Object, val: union.Of4T[map[string]Value, []Value, string, bool](map[string]Value{"foo": z})}, v)
	ma = nil
	assert.Nil(t, conv.To(v, &ma))
	assert.Equal(t, map[string]any{"foo": zs}, ma)

	assert.Nil(t, conv.To(map[string]Value{"foo": z}, &v))
	assert.Equal(t, Value{typ: Object, val: union.Of4T[map[string]Value, []Value, string, bool](map[string]Value{"foo": z})}, v)

	// Array
	v = Value{}
	sa = []any{0}
	assert.Nil(t, conv.To(sa, &v))
	assert.Equal(t, Value{typ: Array, val: union.Of4U[map[string]Value, []Value, string, bool]([]Value{z})}, v)
	sa = nil
	assert.Nil(t, conv.To(v, &sa))
	assert.Equal(t, []any{zs}, sa)

	assert.Nil(t, conv.To([]Value{z}, &v))
	assert.Equal(t, Value{typ: Array, val: union.Of4U[map[string]Value, []Value, string, bool]([]Value{z})}, v)

	// String
	v = Value{}
	assert.Nil(t, conv.To("foo", &v))
	assert.Equal(t, Value{typ: String, val: union.Of4V[map[string]Value, []Value, string, bool]("foo")}, v)
	var str string
	assert.Nil(t, conv.To(v, &str))
	assert.Equal(t, "foo", str)

	// NumberString
	v = Value{}
	assert.Nil(t, conv.To(NumberString("1"), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("1")}, v)
	var ns NumberString
	assert.Nil(t, conv.To(v, &ns))
	assert.Equal(t, NumberString("1"), ns)

	// Signed ints
	v = Value{}
	assert.Nil(t, conv.To(2, &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("2")}, v)
	var i int
	assert.Nil(t, conv.To(v, &i))
	assert.Equal(t, 2, i)

	v = Value{}
	assert.Nil(t, conv.To(int8(3), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("3")}, v)
	var i8 int8
	assert.Nil(t, conv.To(v, &i8))
	assert.Equal(t, int8(3), i8)

	v = Value{}
	assert.Nil(t, conv.To(int16(4), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("4")}, v)
	var i16 int16
	assert.Nil(t, conv.To(v, &i16))
	assert.Equal(t, int16(4), i16)

	v = Value{}
	assert.Nil(t, conv.To(int32(5), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("5")}, v)
	var i32 int32
	assert.Nil(t, conv.To(v, &i32))
	assert.Equal(t, int32(5), i32)

	v = Value{}
	assert.Nil(t, conv.To(int64(6), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("6")}, v)
	var i64 int64
	assert.Nil(t, conv.To(v, &i64))
	assert.Equal(t, int64(6), i64)

	// Unsigned ints
	v = Value{}
	assert.Nil(t, conv.To(uint(7), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("7")}, v)
	var ui uint
	assert.Nil(t, conv.To(v, &ui))
	assert.Equal(t, uint(7), ui)

	v = Value{}
	assert.Nil(t, conv.To(uint8(8), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("8")}, v)
	var ui8 uint8
	assert.Nil(t, conv.To(v, &ui8))
	assert.Equal(t, uint8(8), ui8)

	v = Value{}
	assert.Nil(t, conv.To(uint16(9), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("9")}, v)
	var ui16 uint16
	assert.Nil(t, conv.To(v, &ui16))
	assert.Equal(t, uint16(9), ui16)

	v = Value{}
	assert.Nil(t, conv.To(uint32(10), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("10")}, v)
	var ui32 uint32
	assert.Nil(t, conv.To(v, &ui32))
	assert.Equal(t, uint32(10), ui32)

	v = Value{}
	assert.Nil(t, conv.To(uint64(11), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("11")}, v)
	var ui64 uint64
	assert.Nil(t, conv.To(v, &ui64))
	assert.Equal(t, uint64(11), ui64)

	// Floats
	v = Value{}
	assert.Nil(t, conv.To(float32(12.25), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("12.25")}, v)
	var f32 float32
	assert.Nil(t, conv.To(v, &f32))
	assert.Equal(t, float32(12.25), f32)

	v = Value{}
	assert.Nil(t, conv.To(13.5, &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("13.5")}, v)
	var f64 float64
	assert.Nil(t, conv.To(v, &f64))
	assert.Equal(t, 13.5, f64)

	// Bigs
	v = Value{}
	var bf *big.Float
	assert.Nil(t, conv.To(14.75, &bf))
	assert.Nil(t, conv.To(bf, &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("14.75")}, v)
	bf = nil
	f64 = 0
	assert.Nil(t, conv.To(v, &bf))
	assert.NotNil(t, bf)
	assert.Nil(t, conv.To(bf, &f64))
	assert.Equal(t, 14.75, f64)

	v = Value{}
	var bi *big.Int
	conv.To(15, &bi)
	assert.Nil(t, conv.To(bi, &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("15")}, v)
	bi = nil
	i = 0
	assert.Nil(t, conv.To(v, &bi))
	assert.NotNil(t, bi)
	assert.Nil(t, conv.To(bi, &i))
	assert.Equal(t, 15, i)

	v = Value{}
	var br *big.Rat
	conv.To(16.25, &br)
	assert.Nil(t, conv.To(br, &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("16.25")}, v)
	br = nil
	f64 = 0
	assert.Nil(t, conv.To(v, &br))
	assert.NotNil(t, br)
	assert.Nil(t, conv.To(br, &f64))
	assert.Equal(t, 16.25, f64)
}
