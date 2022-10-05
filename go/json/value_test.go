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

	assert.Equal(t, "true", TrueValue.AsString())

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
	assert.True(t, TrueValue.AsBoolean())

	funcs.TryTo(
		func() {
			NullValue.AsBoolean()
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

	assert.False(t, TrueValue.IsNull())
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
	assert.Equal(t, TrueValue, TrueValue.Visit(noconv))

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

	// Object {"foo": {"bar": "baz"}}
	baz := JSONValue{typ: String, value: "baz"}
	val = JSONValue{typ: Object, value: map[string]JSONValue{
		"foo": {typ: Object, value: map[string]JSONValue{"bar": baz}}},
	}
	assert.Equal(t, map[string]any{"foo": map[string]any{"bar": "baz"}}, val.Visit(DefaultVisitor))

	// Array []
	val = JSONValue{typ: Array, value: []JSONValue{}}
	assert.Equal(t, []any{}, val.Visit(DefaultVisitor))

	// Array ["bar"]
	val = JSONValue{typ: Array, value: []JSONValue{str}}
	assert.Equal(t, []any{"bar"}, val.Visit(DefaultVisitor))

	// Array [ Array ["bar"] ]
	val = JSONValue{typ: Array, value: []JSONValue{{typ: Array, value: []JSONValue{str}}}}
	assert.Equal(t, []any{[]any{"bar"}}, val.Visit(DefaultVisitor))

	// String "bar"
	assert.Equal(t, "bar", str.Visit(DefaultVisitor))

	// Number 0
	val = JSONValue{typ: Number, value: big.NewFloat(0)}
	assert.Equal(t, big.NewFloat(0), val.Visit(DefaultVisitor))

	// Boolean true
	assert.Equal(t, true, TrueValue.Visit(DefaultVisitor))

	// Null
	val = NullValue
	assert.Nil(t, val.Visit(DefaultVisitor))
}

func TestConversionVisitor(t *testing.T) {
	// Object {}
	var (
		defaultVisitor = ConversionVisitor(nil, nil, nil)
		customVisitor  = ConversionVisitor(
			func(str string) any { return str + str },
			func(num *big.Float) any { res := big.NewFloat(0); res.Add(num, num); return res },
			func(b bool) any { return !b },
		)
		intVisitor   = ConversionVisitor(nil, NumberToInt64Conversion, nil)
		floatVisitor = ConversionVisitor(nil, NumberToFloat64Conversion, nil)
		val          = JSONValue{typ: Object, value: map[string]JSONValue{}}
	)
	assert.Equal(t, map[string]any{}, val.Visit(defaultVisitor))
	assert.Equal(t, map[string]any{}, val.Visit(customVisitor))

	// Object {"foo": "bar"}
	str := JSONValue{typ: String, value: "bar"}
	val = JSONValue{typ: Object, value: map[string]JSONValue{"foo": str}}
	assert.Equal(t, map[string]any{"foo": "bar"}, val.Visit(defaultVisitor))
	assert.Equal(t, map[string]any{"foo": "barbar"}, val.Visit(customVisitor))
	assert.Equal(t, map[string]any{"foo": "bar"}, val.Visit(intVisitor))
	assert.Equal(t, map[string]any{"foo": "bar"}, val.Visit(floatVisitor))

	// Object {"foo": {"bar": "baz"}}
	baz := JSONValue{typ: String, value: "baz"}
	val = JSONValue{typ: Object, value: map[string]JSONValue{
		"foo": {typ: Object, value: map[string]JSONValue{"bar": baz}}},
	}
	assert.Equal(t, map[string]any{"foo": map[string]any{"bar": "baz"}}, val.Visit(defaultVisitor))
	assert.Equal(t, map[string]any{"foo": map[string]any{"bar": "bazbaz"}}, val.Visit(customVisitor))
	assert.Equal(t, map[string]any{"foo": map[string]any{"bar": "baz"}}, val.Visit(intVisitor))
	assert.Equal(t, map[string]any{"foo": map[string]any{"bar": "baz"}}, val.Visit(floatVisitor))

	// Array []
	val = JSONValue{typ: Array, value: []JSONValue{}}
	assert.Equal(t, []any{}, val.Visit(defaultVisitor))
	assert.Equal(t, []any{}, val.Visit(customVisitor))
	assert.Equal(t, []any{}, val.Visit(intVisitor))
	assert.Equal(t, []any{}, val.Visit(floatVisitor))

	// Array ["bar"]
	val = JSONValue{typ: Array, value: []JSONValue{str}}
	assert.Equal(t, []any{"bar"}, val.Visit(defaultVisitor))
	assert.Equal(t, []any{"barbar"}, val.Visit(customVisitor))
	assert.Equal(t, []any{"bar"}, val.Visit(intVisitor))
	assert.Equal(t, []any{"bar"}, val.Visit(floatVisitor))

	// Array [ Array ["bar"] ]
	val = JSONValue{typ: Array, value: []JSONValue{{typ: Array, value: []JSONValue{str}}}}
	assert.Equal(t, []any{[]any{"bar"}}, val.Visit(defaultVisitor))
	assert.Equal(t, []any{[]any{"barbar"}}, val.Visit(customVisitor))
	assert.Equal(t, []any{[]any{"bar"}}, val.Visit(intVisitor))
	assert.Equal(t, []any{[]any{"bar"}}, val.Visit(floatVisitor))

	// String "bar"
	assert.Equal(t, "bar", str.Visit(defaultVisitor))
	assert.Equal(t, "barbar", str.Visit(customVisitor))
	assert.Equal(t, "bar", str.Visit(intVisitor))
	assert.Equal(t, "bar", str.Visit(floatVisitor))

	// Number 5
	val = JSONValue{typ: Number, value: big.NewFloat(5)}
	assert.Equal(t, big.NewFloat(5), val.Visit(defaultVisitor))
	assert.Equal(t, big.NewFloat(10), val.Visit(customVisitor))
	assert.Equal(t, int64(5), val.Visit(intVisitor))
	assert.Equal(t, float64(5), val.Visit(floatVisitor))

	// Boolean true
	assert.Equal(t, true, TrueValue.Visit(defaultVisitor))
	assert.Equal(t, false, TrueValue.Visit(customVisitor))
	assert.Equal(t, true, TrueValue.Visit(intVisitor))
	assert.Equal(t, true, TrueValue.Visit(floatVisitor))

	// Null
	assert.Nil(t, NullValue.Visit(defaultVisitor))
	assert.Nil(t, NullValue.Visit(intVisitor))
	assert.Nil(t, NullValue.Visit(floatVisitor))
}

func TestToMap(t *testing.T) {
	// Object {}
	var (
		defaultVisitor = ConversionVisitor(nil, nil, nil)
		customVisitor  = ConversionVisitor(
			func(str string) any { return str + str },
			func(num *big.Float) any { res := big.NewFloat(0); res.Add(num, num); return res },
			func(b bool) any { return !b },
		)
		intVisitor   = ConversionVisitor(nil, NumberToInt64Conversion, nil)
		floatVisitor = ConversionVisitor(nil, NumberToFloat64Conversion, nil)
		val          = JSONValue{typ: Object, value: map[string]JSONValue{}}
	)
	assert.Equal(t, map[string]any{}, val.ToMap())
	assert.Equal(t, map[string]any{}, val.ToMap(defaultVisitor))
	assert.Equal(t, map[string]any{}, val.ToMap(customVisitor))

	// Object {"foo": "bar"}
	str := JSONValue{typ: String, value: "bar"}
	val = JSONValue{typ: Object, value: map[string]JSONValue{"foo": str}}
	assert.Equal(t, map[string]any{"foo": "bar"}, val.ToMap())
	assert.Equal(t, map[string]any{"foo": "bar"}, val.ToMap(defaultVisitor))
	assert.Equal(t, map[string]any{"foo": "barbar"}, val.ToMap(customVisitor))
	assert.Equal(t, map[string]any{"foo": "bar"}, val.ToMap(intVisitor))
	assert.Equal(t, map[string]any{"foo": "bar"}, val.ToMap(floatVisitor))

	// Object {"foo": 5}
	num := JSONValue{typ: Number, value: big.NewFloat(5)}
	val = JSONValue{typ: Object, value: map[string]JSONValue{"foo": num}}
	assert.Equal(t, map[string]any{"foo": big.NewFloat(5)}, val.ToMap())
	assert.Equal(t, map[string]any{"foo": big.NewFloat(5)}, val.ToMap(defaultVisitor))
	assert.Equal(t, map[string]any{"foo": big.NewFloat(10)}, val.ToMap(customVisitor))
	assert.Equal(t, map[string]any{"foo": int64(5)}, val.ToMap(intVisitor))
	assert.Equal(t, map[string]any{"foo": 5.0}, val.ToMap(floatVisitor))

	// Object {"foo": true}
	val = JSONValue{typ: Object, value: map[string]JSONValue{"foo": TrueValue}}
	assert.Equal(t, map[string]any{"foo": true}, val.ToMap())
	assert.Equal(t, map[string]any{"foo": true}, val.ToMap(defaultVisitor))
	assert.Equal(t, map[string]any{"foo": false}, val.ToMap(customVisitor))
	assert.Equal(t, map[string]any{"foo": true}, val.ToMap(intVisitor))
	assert.Equal(t, map[string]any{"foo": true}, val.ToMap(floatVisitor))

	// Object {"foo": nil}
	val = JSONValue{typ: Object, value: map[string]JSONValue{"foo": NullValue}}
	assert.Equal(t, map[string]any{"foo": nil}, val.ToMap())
	assert.Equal(t, map[string]any{"foo": nil}, val.ToMap(defaultVisitor))
	assert.Equal(t, map[string]any{"foo": nil}, val.ToMap(customVisitor))
	assert.Equal(t, map[string]any{"foo": nil}, val.ToMap(intVisitor))
	assert.Equal(t, map[string]any{"foo": nil}, val.ToMap(floatVisitor))

	// Panic if value is not an object
	funcs.TryTo(
		func() {
			TrueValue.ToMap()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotObject, e)
		},
	)
}

func TestToSlice(t *testing.T) {
	// Array []
	var (
		defaultVisitor = ConversionVisitor(nil, nil, nil)
		customVisitor  = ConversionVisitor(
			func(str string) any { return str + str },
			func(num *big.Float) any { res := big.NewFloat(0); res.Add(num, num); return res },
			func(b bool) any { return !b },
		)
		intVisitor   = ConversionVisitor(nil, NumberToInt64Conversion, nil)
		floatVisitor = ConversionVisitor(nil, NumberToFloat64Conversion, nil)
		val          = JSONValue{typ: Array, value: []JSONValue{}}
	)
	assert.Equal(t, []any{}, val.ToSlice())
	assert.Equal(t, []any{}, val.ToSlice(defaultVisitor))
	assert.Equal(t, []any{}, val.ToSlice(customVisitor))

	// Array ["bar"]
	str := JSONValue{typ: String, value: "bar"}
	val = JSONValue{typ: Array, value: []JSONValue{str}}
	assert.Equal(t, []any{"bar"}, val.ToSlice())
	assert.Equal(t, []any{"bar"}, val.ToSlice(defaultVisitor))
	assert.Equal(t, []any{"barbar"}, val.ToSlice(customVisitor))
	assert.Equal(t, []any{"bar"}, val.ToSlice(intVisitor))
	assert.Equal(t, []any{"bar"}, val.ToSlice(floatVisitor))

	// Array [5]
	num := JSONValue{typ: Number, value: big.NewFloat(5)}
	val = JSONValue{typ: Array, value: []JSONValue{num}}
	assert.Equal(t, []any{big.NewFloat(5)}, val.ToSlice())
	assert.Equal(t, []any{big.NewFloat(5)}, val.ToSlice(defaultVisitor))
	assert.Equal(t, []any{big.NewFloat(10)}, val.ToSlice(customVisitor))
	assert.Equal(t, []any{int64(5)}, val.ToSlice(intVisitor))
	assert.Equal(t, []any{5.0}, val.ToSlice(floatVisitor))

	// Array [true]
	val = JSONValue{typ: Array, value: []JSONValue{TrueValue}}
	assert.Equal(t, []any{true}, val.ToSlice())
	assert.Equal(t, []any{true}, val.ToSlice(defaultVisitor))
	assert.Equal(t, []any{false}, val.ToSlice(customVisitor))
	assert.Equal(t, []any{true}, val.ToSlice(intVisitor))
	assert.Equal(t, []any{true}, val.ToSlice(floatVisitor))

	// Array [nil]
	val = JSONValue{typ: Array, value: []JSONValue{NullValue}}
	assert.Equal(t, []any{nil}, val.ToSlice())
	assert.Equal(t, []any{nil}, val.ToSlice(defaultVisitor))
	assert.Equal(t, []any{nil}, val.ToSlice(customVisitor))
	assert.Equal(t, []any{nil}, val.ToSlice(intVisitor))
	assert.Equal(t, []any{nil}, val.ToSlice(floatVisitor))

	// Panic if value is not an array
	funcs.TryTo(
		func() {
			TrueValue.ToSlice()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotArray, e)
		},
	)
}

func TestToInt(t *testing.T) {
	val := JSONValue{typ: Number, value: big.NewFloat(5)}
	assert.Equal(t, int64(5), val.ToInt())

	funcs.TryTo(
		func() {
			NullValue.ToInt()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotNumber, e)
		},
	)
}

func TestToFloat(t *testing.T) {
	val := JSONValue{typ: Number, value: big.NewFloat(5)}
	assert.Equal(t, 5.0, val.ToFloat())

	funcs.TryTo(
		func() {
			NullValue.ToFloat()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotNumber, e)
		},
	)
}
