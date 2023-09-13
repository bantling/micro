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
	var v, z Value
	z = MustNumberToValue(0)

	// Object
	assert.Nil(t, conv.To(map[string]any{"foo": 0}, &v))
	assert.Equal(t, Value{typ: Object, val: union.Of4T[map[string]Value, []Value, string, bool](map[string]Value{"foo": z})}, v)
	assert.Nil(t, conv.To(map[string]Value{"foo": z}, &v))
	assert.Equal(t, Value{typ: Object, val: union.Of4T[map[string]Value, []Value, string, bool](map[string]Value{"foo": z})}, v)

	// Array
	assert.Nil(t, conv.To([]any{0}, &v))
	assert.Equal(t, Value{typ: Array, val: union.Of4U[map[string]Value, []Value, string, bool]([]Value{z})}, v)
	assert.Nil(t, conv.To([]Value{z}, &v))
	assert.Equal(t, Value{typ: Array, val: union.Of4U[map[string]Value, []Value, string, bool]([]Value{z})}, v)

	// String
	assert.Nil(t, conv.To("foo", &v))
	assert.Equal(t, Value{typ: String, val: union.Of4V[map[string]Value, []Value, string, bool]("foo")}, v)

	// NumberString
	assert.Nil(t, conv.To(NumberString("1"), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("1")}, v)

	// Signed ints
	assert.Nil(t, conv.To(2, &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("2")}, v)

	assert.Nil(t, conv.To(int8(3), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("3")}, v)

	assert.Nil(t, conv.To(int16(4), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("4")}, v)

	assert.Nil(t, conv.To(int32(5), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("5")}, v)

	assert.Nil(t, conv.To(int64(6), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("6")}, v)

	// Unsigned ints
	assert.Nil(t, conv.To(uint(7), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("7")}, v)

	assert.Nil(t, conv.To(uint8(8), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("8")}, v)

	assert.Nil(t, conv.To(uint16(9), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("9")}, v)

	assert.Nil(t, conv.To(uint32(10), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("10")}, v)

	assert.Nil(t, conv.To(uint64(11), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("11")}, v)

	// Floats
	assert.Nil(t, conv.To(float32(12.25), &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("12.25")}, v)

	assert.Nil(t, conv.To(13.5, &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("13.5")}, v)

	// Bigs
	var bf *big.Float
	conv.To(14.75, &bf)
	assert.Nil(t, conv.To(bf, &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("14.75")}, v)

	var bi *big.Int
	conv.To(15, &bi)
	assert.Nil(t, conv.To(bi, &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("15")}, v)

	var br *big.Rat
	conv.To(16.25, &br)
	assert.Nil(t, conv.To(br, &v))
	assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("16.25")}, v)
}
