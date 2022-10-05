package json

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/go/funcs"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType(t *testing.T) {
	val := JSONValue{typ: Object, value: map[string]JSONValue{}}
	assert.Equal(t, Object, val.Type())
}

func TestAsMap(t *testing.T) {
	val := JSONValue{typ: Object, value: map[string]JSONValue{}}
	assert.Equal(t, val.value, val.AsMap())

	val = NullValue
	funcs.TryTo(
		func() {
			val.AsMap()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotObject, e)
		},
	)
}

func TestAsSlice(t *testing.T) {
	val := JSONValue{typ: Array, value: []JSONValue{}}
	assert.Equal(t, val.value, val.AsSlice())

	val = NullValue
	funcs.TryTo(
		func() {
			val.AsSlice()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotArray, e)
		},
	)
}

func TestAsString(t *testing.T) {
	val := JSONValue{typ: String, value: ""}
	assert.Equal(t, val.value, val.AsString())

	val = JSONValue{typ: Number, value: big.NewFloat(0)}
	assert.Equal(t, "0", val.AsString())

	val = JSONValue{typ: Boolean, value: true}
	assert.Equal(t, "true", val.AsString())

	val = NullValue
	funcs.TryTo(
		func() {
			val.AsString()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotStringable, e)
		},
	)
}

func TestAsNumber(t *testing.T) {
	val := JSONValue{typ: Number, value: big.NewFloat(0)}
	assert.Equal(t, val.value, val.AsNumber())

	val = NullValue
	funcs.TryTo(
		func() {
			val.AsNumber()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotNumber, e)
		},
	)
}

func TestAsBoolean(t *testing.T) {
	val := JSONValue{typ: Boolean, value: true}
	assert.Equal(t, val.value, val.AsBoolean())

	val = NullValue
	funcs.TryTo(
		func() {
			val.AsBoolean()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotBoolean, e)
		},
	)
}

func TestIsNull(t *testing.T) {
	val := NullValue
	assert.True(t, val.IsNull())

	val = JSONValue{typ: Boolean, value: true}
	assert.False(t, val.IsNull())
}

func TestVisit(t *testing.T) {
	// Object {}
	var (
		val    = JSONValue{typ: Object, value: map[string]JSONValue{}}
		noconv = func(jv JSONValue) any { return jv }
	)
	assert.Equal(t, map[string]any{}, val.Visit(noconv))

	// Object {"foo": "bar"}
	str := JSONValue{typ: String, value: "bar"}
	val = JSONValue{typ: Object, value: map[string]JSONValue{"foo": str}}
	assert.Equal(t, map[string]any{"foo": str}, val.Visit(noconv))

	// Array []
	val = JSONValue{typ: Array, value: []JSONValue{}}
	assert.Equal(t, []any{}, val.Visit(noconv))

	// Array ["bar"]
	val = JSONValue{typ: Array, value: []JSONValue{str}}
	assert.Equal(t, []any{str}, val.Visit(noconv))

	// String "bar"
	assert.Equal(t, str, str.Visit(noconv))

	// Number 0
	val = JSONValue{typ: Number, value: big.NewFloat(0)}
	assert.Equal(t, val, val.Visit(noconv))

	// Boolean true
	val = JSONValue{typ: Boolean, value: true}
	assert.Equal(t, val, val.Visit(noconv))

	// Null
	val = NullValue
	assert.Equal(t, val, val.Visit(noconv))
}

func TestDefaultVisitor(t *testing.T) {
	// Object {}
	val := JSONValue{typ: Object, value: map[string]JSONValue{}}
	assert.Equal(t, map[string]any{}, val.Visit(DefaultVisitor))

	// Object {"foo": "bar"}
	str := JSONValue{typ: String, value: "bar"}
	val = JSONValue{typ: Object, value: map[string]JSONValue{"foo": str}}
	assert.Equal(t, map[string]any{"foo": "bar"}, val.Visit(DefaultVisitor))

	// Array []
	val = JSONValue{typ: Array, value: []JSONValue{}}
	assert.Equal(t, []any{}, val.Visit(DefaultVisitor))

	// Array ["bar"]
	val = JSONValue{typ: Array, value: []JSONValue{str}}
	assert.Equal(t, []any{"bar"}, val.Visit(DefaultVisitor))

	// String "bar"
	assert.Equal(t, "bar", str.Visit(DefaultVisitor))

	// Number 0
	val = JSONValue{typ: Number, value: big.NewFloat(0)}
	assert.Equal(t, big.NewFloat(0), val.Visit(DefaultVisitor))

	// Boolean true
	val = JSONValue{typ: Boolean, value: true}
	assert.Equal(t, true, val.Visit(DefaultVisitor))

	// Null
	val = NullValue
	assert.Nil(t, val.Visit(DefaultVisitor))
}

func TestConversionVisitor(t *testing.T) {
	// Object {}
	var (
		visitor = ConversionVisitor(nil, nil, nil)
		val     = JSONValue{typ: Object, value: map[string]JSONValue{}}
	)
	assert.Equal(t, map[string]any{}, val.Visit(visitor))

	// Object {"foo": "bar"}
	str := JSONValue{typ: String, value: "bar"}
	val = JSONValue{typ: Object, value: map[string]JSONValue{"foo": str}}
	assert.Equal(t, map[string]any{"foo": "bar"}, val.Visit(visitor))

	// Array []
	val = JSONValue{typ: Array, value: []JSONValue{}}
	assert.Equal(t, []any{}, val.Visit(visitor))

	// Array ["bar"]
	val = JSONValue{typ: Array, value: []JSONValue{str}}
	assert.Equal(t, []any{"bar"}, val.Visit(visitor))

	// String "bar"
	assert.Equal(t, "bar", str.Visit(visitor))

	// Number 0
	val = JSONValue{typ: Number, value: big.NewFloat(0)}
	assert.Equal(t, big.NewFloat(0), val.Visit(visitor))

	// Boolean true
	val = JSONValue{typ: Boolean, value: true}
	assert.Equal(t, true, val.Visit(visitor))

	// Null
	assert.Nil(t, NullValue.Visit(visitor))
}
