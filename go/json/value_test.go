package json

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/bantling/micro/go/conv"
	"github.com/bantling/micro/go/funcs"
	"github.com/stretchr/testify/assert"
)

func assertObject(t *testing.T, e map[string]Value, a Value) {
	assert.Equal(t, Value{typ: Object, value: e}, a)
}

func assertArray(t *testing.T, e []Value, a Value) {
	assert.Equal(t, Value{typ: Array, value: e}, a)
}

func assertString(t *testing.T, e string, a Value) {
	assert.Equal(t, Value{typ: String, value: e}, a)
}

func assertNumber(t *testing.T, e *big.Rat, a Value) {
	assert.Equal(t, Value{typ: Number, value: e}, a)
}

func assertBoolean(t *testing.T, e bool, a Value) {
	assert.Equal(t, Value{typ: Boolean, value: e}, a)
}

func assertNull(t *testing.T, a Value) {
	assert.Equal(t, NullValue, a)
}

func TestFromNumberInternal(t *testing.T) {
	assertNumber(t, conv.IntToBigRat(1), fromNumberInternal(int(1)))
	assertNumber(t, conv.IntToBigRat(2), fromNumberInternal(int8(2)))
	assertNumber(t, conv.IntToBigRat(3), fromNumberInternal(int16(3)))
	assertNumber(t, conv.IntToBigRat(4), fromNumberInternal(int32(4)))
	assertNumber(t, conv.IntToBigRat(5), fromNumberInternal(int64(5)))

	assertNumber(t, conv.UintToBigRat(uint64(1)), fromNumberInternal(uint(1)))
	assertNumber(t, conv.UintToBigRat(uint64(2)), fromNumberInternal(uint8(2)))
	assertNumber(t, conv.UintToBigRat(uint64(3)), fromNumberInternal(uint16(3)))
	assertNumber(t, conv.UintToBigRat(uint64(4)), fromNumberInternal(uint32(4)))
	assertNumber(t, conv.UintToBigRat(uint64(5)), fromNumberInternal(uint64(5)))

	assertNumber(t, conv.FloatToBigRat(1.25), fromNumberInternal(float32(1.25)))
	assertNumber(t, conv.FloatToBigRat(2.5), fromNumberInternal(float64(2.5)))

	assertNumber(t, conv.IntToBigRat(3), fromNumberInternal(conv.IntToBigInt(3)))
	assertNumber(t, conv.FloatToBigRat(3.75), fromNumberInternal(conv.FloatToBigFloat(3.75)))
	assertNumber(t, conv.FloatToBigRat(4.25), fromNumberInternal(conv.FloatToBigRat(4.25)))

	// fromNumberInternal accepts type any, so explicit conversion required
	assertNumber(t, conv.FloatToBigRat(5.75), fromNumberInternal(NumberString("5.75")))

	// Any other type results in invalid zero value
	assert.Equal(t, Value{}, fromNumberInternal(""))
}

func TestFromValue(t *testing.T) {
	// Object
	assertObject(t, map[string]Value{"foo": FromString("bar")}, FromValue(map[string]any{"foo": "bar"}))

	// Array
	assertArray(t, []Value{FromString("bar")}, FromValue([]any{"bar"}))

	// String
	assertString(t, "bar", FromValue("bar"))

	// Number - int
	assertNumber(t, conv.IntToBigRat(1), FromValue(int(1)))
	assertNumber(t, conv.IntToBigRat(2), FromValue(int8(2)))
	assertNumber(t, conv.IntToBigRat(3), FromValue(int16(3)))
	assertNumber(t, conv.IntToBigRat(4), FromValue(int32(4)))
	assertNumber(t, conv.IntToBigRat(5), FromValue(int64(5)))

	// Number - uint
	assertNumber(t, conv.IntToBigRat(1), FromValue(uint(1)))
	assertNumber(t, conv.IntToBigRat(2), FromValue(uint8(2)))
	assertNumber(t, conv.IntToBigRat(3), FromValue(uint16(3)))
	assertNumber(t, conv.IntToBigRat(4), FromValue(uint32(4)))
	assertNumber(t, conv.IntToBigRat(5), FromValue(uint64(5)))

	// Number - float
	assertNumber(t, conv.FloatToBigRat(1.25), FromValue(float32(1.25)))
	assertNumber(t, conv.FloatToBigRat(2.5), FromValue(float64(2.5)))

	// Number - *big.Int, *big.Float, *big.Rat
	assertNumber(t, conv.FloatToBigRat(3.0), FromValue(conv.IntToBigInt(3)))
	assertNumber(t, conv.FloatToBigRat(3.75), FromValue(conv.FloatToBigFloat(3.75)))
	assertNumber(t, conv.FloatToBigRat(4.25), FromValue(conv.FloatToBigRat(4.25)))

	// Number - NumberString
	// fromValue accepts type any, so explicit conversion required
	assertNumber(t, conv.FloatToBigRat(5.75), FromValue(NumberString("5.75")))

	// Number - custom conv to float64
	assert.Equal(t, Value{typ: Number, value: 1.25}, FromValue(1.25, func(v any) Value { return Value{typ: Number, value: v.(float64)} }))

	// Boolean - true
	assertBoolean(t, true, FromValue(true))

	// Boolean - false
	assertBoolean(t, false, FromValue(false))

	// Null
	assertNull(t, FromValue(nil))

	// Error
	funcs.TryTo(
		func() { FromValue((1 + 2i)) },
		func(e any) { assert.Equal(t, fmt.Errorf(ErrInvalidGoValueMsg, (1+2i)), e) },
	)
}

func TestFromMap(t *testing.T) {
	assertObject(t, map[string]Value{"foo": {typ: String, value: "bar"}}, FromMap(map[string]any{"foo": "bar"}))
}

func TestFromSlice(t *testing.T) {
	assertArray(t, []Value{{typ: String, value: "bar"}}, FromSlice([]any{"bar"}))
}

func TestFromString(t *testing.T) {
	assertString(t, "bar", FromString("bar"))
}

func TestFromSignedInt(t *testing.T) {
	assertNumber(t, conv.IntToBigRat(1), FromSignedInt(int(1)))
	assertNumber(t, conv.IntToBigRat(2), FromSignedInt(int8(2)))
	assertNumber(t, conv.IntToBigRat(3), FromSignedInt(int16(3)))
	assertNumber(t, conv.IntToBigRat(4), FromSignedInt(int32(4)))
	assertNumber(t, conv.IntToBigRat(5), FromSignedInt(int64(5)))
}

func TestFromUnsignedInt(t *testing.T) {
	assertNumber(t, conv.IntToBigRat(1), FromUnsignedInt(uint(1)))
	assertNumber(t, conv.IntToBigRat(2), FromUnsignedInt(uint8(2)))
	assertNumber(t, conv.IntToBigRat(3), FromUnsignedInt(uint16(3)))
	assertNumber(t, conv.IntToBigRat(4), FromUnsignedInt(uint32(4)))
	assertNumber(t, conv.IntToBigRat(5), FromUnsignedInt(uint64(5)))
}

func TestFromFloat(t *testing.T) {
	assertNumber(t, conv.FloatToBigRat(1.25), FromFloat(float32(1.25)))
	assertNumber(t, conv.FloatToBigRat(2.5), FromFloat(float64(2.5)))
}

func TestFromBigInt(t *testing.T) {
	assertNumber(t, conv.IntToBigRat(3), FromBigInt(conv.IntToBigInt(3)))
}

func TestFromBigFloat(t *testing.T) {
	assertNumber(t, conv.FloatToBigRat(3.75), FromBigFloat(conv.FloatToBigFloat(3.75)))
}

func TestFromBigRat(t *testing.T) {
	assertNumber(t, conv.FloatToBigRat(4.25), FromBigFloat(conv.FloatToBigFloat(4.25)))
}

func TestFromNumberString(t *testing.T) {
	// fromNumberString accepts type NumberString, so implicit conversion allowed
	assertNumber(t, conv.FloatToBigRat(5.75), FromNumberString("5.75"))
}

func TestFromNumber(t *testing.T) {
	// Number - int
	assertNumber(t, conv.IntToBigRat(1), FromNumber(int(1)))
	assertNumber(t, conv.IntToBigRat(2), FromNumber(int8(2)))
	assertNumber(t, conv.IntToBigRat(3), FromNumber(int16(3)))
	assertNumber(t, conv.IntToBigRat(4), FromNumber(int32(4)))
	assertNumber(t, conv.IntToBigRat(5), FromNumber(int64(5)))

	// Number - uint
	assertNumber(t, conv.IntToBigRat(1), FromNumber(uint(1)))
	assertNumber(t, conv.IntToBigRat(2), FromNumber(uint8(2)))
	assertNumber(t, conv.IntToBigRat(3), FromNumber(uint16(3)))
	assertNumber(t, conv.IntToBigRat(4), FromNumber(uint32(4)))
	assertNumber(t, conv.IntToBigRat(5), FromNumber(uint64(5)))

	// Number - float
	assertNumber(t, conv.FloatToBigRat(1.25), FromNumber(float32(1.25)))
	assertNumber(t, conv.FloatToBigRat(2.5), FromNumber(float64(2.5)))

	// Number - *big.Int, *big.Float, *big.Rat
	assertNumber(t, conv.IntToBigRat(3), FromNumber(conv.IntToBigInt(3)))
	assertNumber(t, conv.FloatToBigRat(3.75), FromNumber(conv.FloatToBigFloat(3.75)))
	assertNumber(t, conv.FloatToBigRat(4.25), FromNumber(conv.FloatToBigRat(4.25)))

	// Number - NumberString
	// FromNumber accepts type any, so explicit conversion required
	assertNumber(t, conv.FloatToBigRat(5.75), FromNumber(NumberString("5.75")))
}

func TestFromBool(t *testing.T) {
	assertBoolean(t, true, FromBool(true))
	assertBoolean(t, false, FromBool(false))
}

func TestType(t *testing.T) {
	val := FromMap(map[string]any{})
	assert.Equal(t, Object, val.Type())
}

func TestAsMap(t *testing.T) {
	val := FromMap(map[string]any{})
	assert.Equal(t, val.value, val.AsMap())

	funcs.TryTo(
		func() {
			NullValue.AsMap()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotObject, e)
		},
	)
}

func TestAsSlice(t *testing.T) {
	val := FromSlice([]any{})
	assert.Equal(t, val.value, val.AsSlice())

	funcs.TryTo(
		func() {
			NullValue.AsSlice()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotArray, e)
		},
	)
}

func TestAsString(t *testing.T) {
	val := FromString("foo")
	assert.Equal(t, val.value, val.AsString())

	val = FromNumber(1)
	assert.Equal(t, "1", val.AsString())

	assert.Equal(t, "true", TrueValue.AsString())

	funcs.TryTo(
		func() {
			NullValue.AsString()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotStringable, e)
		},
	)
}

func TestAsBigRat(t *testing.T) {
	val := FromNumber(1)
	assert.Equal(t, conv.IntToBigRat(1), val.AsBigRat())

	funcs.TryTo(
		func() {
			NullValue.AsBigRat()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotNumber, e)
		},
	)
}

func TestAsBool(t *testing.T) {
	assert.True(t, TrueValue.AsBool())

	funcs.TryTo(
		func() {
			NullValue.AsBool()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotBoolean, e)
		},
	)
}

func TestIsNull(t *testing.T) {
	assert.True(t, NullValue.IsNull())
	assert.False(t, TrueValue.IsNull())
}

func TestVisit(t *testing.T) {
	// Object {}
	var (
		val    = Value{typ: Object, value: map[string]Value{}}
		noconv = func(jv Value) any { return jv }
	)
	assert.Equal(t, map[string]any{}, val.Visit(noconv))

	// Object {"foo": "bar"}
	str := Value{typ: String, value: "bar"}
	val = Value{typ: Object, value: map[string]Value{"foo": str}}
	assert.Equal(t, map[string]any{"foo": str}, val.Visit(noconv))

	// Array []
	val = Value{typ: Array, value: []Value{}}
	assert.Equal(t, []any{}, val.Visit(noconv))

	// Array ["bar"]
	val = Value{typ: Array, value: []Value{str}}
	assert.Equal(t, []any{str}, val.Visit(noconv))

	// String "bar"
	assert.Equal(t, str, str.Visit(noconv))

	// Number 0
	val = Value{typ: Number, value: conv.IntToBigRat(0)}
	assert.Equal(t, val, val.Visit(noconv))

	// Boolean true
	assert.Equal(t, TrueValue, TrueValue.Visit(noconv))

	// Null
	assert.Equal(t, NullValue, NullValue.Visit(noconv))
}

func TestDefaultVisitor(t *testing.T) {
	// Object {}
	val := Value{typ: Object, value: map[string]Value{}}
	assert.Equal(t, map[string]any{}, val.Visit(DefaultVisitorFunc))

	// Object {"foo": "bar"}
	str := Value{typ: String, value: "bar"}
	val = Value{typ: Object, value: map[string]Value{"foo": str}}
	assert.Equal(t, map[string]any{"foo": "bar"}, val.Visit(DefaultVisitorFunc))

	// Object {"foo": {"bar": "baz"}}
	baz := Value{typ: String, value: "baz"}
	val = Value{typ: Object, value: map[string]Value{
		"foo": {typ: Object, value: map[string]Value{"bar": baz}}},
	}
	assert.Equal(t, map[string]any{"foo": map[string]any{"bar": "baz"}}, val.Visit(DefaultVisitorFunc))

	// Array []
	val = Value{typ: Array, value: []Value{}}
	assert.Equal(t, []any{}, val.Visit(DefaultVisitorFunc))

	// Array ["bar"]
	val = Value{typ: Array, value: []Value{str}}
	assert.Equal(t, []any{"bar"}, val.Visit(DefaultVisitorFunc))

	// Array [ Array ["bar"] ]
	val = Value{typ: Array, value: []Value{{typ: Array, value: []Value{str}}}}
	assert.Equal(t, []any{[]any{"bar"}}, val.Visit(DefaultVisitorFunc))

	// String "bar"
	assert.Equal(t, "bar", str.Visit(DefaultVisitorFunc))

	// Number 0
	val = Value{typ: Number, value: conv.IntToBigRat(0)}
	assert.Equal(t, conv.IntToBigRat(0), val.Visit(DefaultVisitorFunc))

	// Boolean true
	assert.Equal(t, true, TrueValue.Visit(DefaultVisitorFunc))

	// Null
	assert.Nil(t, NullValue.Visit(DefaultVisitorFunc))
}

func TestConversionVisitor(t *testing.T) {
	// Object {}
	var (
		defaultVisitor = DefaultConversionVisitor
		customVisitor  = ConversionVisitor(
			func(str string) string { return str + str },
			func(num any) *big.Rat {
				num_ := num.(*big.Rat)
				res := big.NewRat(0, 1)
				res.Add(num_, num_)
				return res
			},
			func(b bool) bool { return !b },
		)
		intVisitor = ConversionVisitor(
			funcs.Passthrough[string],
			func(v any) int64 { return conv.BigRatToInt64(v.(*big.Rat)) },
			funcs.Passthrough[bool],
		)
		floatVisitor = ConversionVisitor(
			funcs.Passthrough[string],
			func(v any) float64 { return conv.BigRatToFloat64(v.(*big.Rat)) },
			funcs.Passthrough[bool],
		)
		val = Value{typ: Object, value: map[string]Value{}}
	)
	assert.Equal(t, map[string]any{}, val.Visit(defaultVisitor))
	assert.Equal(t, map[string]any{}, val.Visit(customVisitor))
	assert.Equal(t, map[string]any{}, val.Visit(intVisitor))
	assert.Equal(t, map[string]any{}, val.Visit(floatVisitor))

	// Object {"foo": "bar"}
	str := Value{typ: String, value: "bar"}
	val = Value{typ: Object, value: map[string]Value{"foo": str}}
	assert.Equal(t, map[string]any{"foo": "bar"}, val.Visit(defaultVisitor))
	assert.Equal(t, map[string]any{"foo": "barbar"}, val.Visit(customVisitor))
	assert.Equal(t, map[string]any{"foo": "bar"}, val.Visit(intVisitor))
	assert.Equal(t, map[string]any{"foo": "bar"}, val.Visit(floatVisitor))

	// Object {"foo": {"bar": "baz"}}
	baz := Value{typ: String, value: "baz"}
	val = Value{typ: Object, value: map[string]Value{
		"foo": {typ: Object, value: map[string]Value{"bar": baz}}},
	}
	assert.Equal(t, map[string]any{"foo": map[string]any{"bar": "baz"}}, val.Visit(defaultVisitor))
	assert.Equal(t, map[string]any{"foo": map[string]any{"bar": "bazbaz"}}, val.Visit(customVisitor))
	assert.Equal(t, map[string]any{"foo": map[string]any{"bar": "baz"}}, val.Visit(intVisitor))
	assert.Equal(t, map[string]any{"foo": map[string]any{"bar": "baz"}}, val.Visit(floatVisitor))

	// Array []
	val = Value{typ: Array, value: []Value{}}
	assert.Equal(t, []any{}, val.Visit(defaultVisitor))
	assert.Equal(t, []any{}, val.Visit(customVisitor))
	assert.Equal(t, []any{}, val.Visit(intVisitor))
	assert.Equal(t, []any{}, val.Visit(floatVisitor))

	// Array ["bar"]
	val = Value{typ: Array, value: []Value{str}}
	assert.Equal(t, []any{"bar"}, val.Visit(defaultVisitor))
	assert.Equal(t, []any{"barbar"}, val.Visit(customVisitor))
	assert.Equal(t, []any{"bar"}, val.Visit(intVisitor))
	assert.Equal(t, []any{"bar"}, val.Visit(floatVisitor))

	// Array [ Array ["bar"] ]
	val = Value{typ: Array, value: []Value{{typ: Array, value: []Value{str}}}}
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
	val = Value{typ: Number, value: big.NewRat(5, 1)}
	assert.Equal(t, big.NewRat(5, 1), val.Visit(defaultVisitor))
	assert.Equal(t, big.NewRat(10, 1), val.Visit(customVisitor))
	assert.Equal(t, int64(5), val.Visit(intVisitor))
	assert.Equal(t, float64(5), val.Visit(floatVisitor))

	// Boolean true
	assert.Equal(t, true, TrueValue.Visit(defaultVisitor))
	assert.Equal(t, false, TrueValue.Visit(customVisitor))
	assert.Equal(t, true, TrueValue.Visit(intVisitor))
	assert.Equal(t, true, TrueValue.Visit(floatVisitor))

	// Null
	assert.Nil(t, NullValue.Visit(defaultVisitor))
	assert.Nil(t, NullValue.Visit(customVisitor))
	assert.Nil(t, NullValue.Visit(intVisitor))
	assert.Nil(t, NullValue.Visit(floatVisitor))
}

func TestToMap(t *testing.T) {
	// Object {}
	var (
		defaultVisitor = DefaultConversionVisitor
		customVisitor  = ConversionVisitor(
			func(str string) string { return str + str },
			func(num any) *big.Rat {
				num_ := num.(*big.Rat)
				res := big.NewRat(0, 1)
				res.Add(num_, num_)
				return res
			},
			func(b bool) bool { return !b },
		)
		intVisitor = ConversionVisitor(
			funcs.Passthrough[string],
			func(num any) int64 { return conv.BigRatToInt64(num.(*big.Rat)) },
			funcs.Passthrough[bool],
		)
		floatVisitor = ConversionVisitor(
			funcs.Passthrough[string],
			func(num any) float64 { return conv.BigRatToFloat64(num.(*big.Rat)) },
			funcs.Passthrough[bool],
		)
		val = Value{typ: Object, value: map[string]Value{}}
	)
	assert.Equal(t, map[string]any{}, val.ToMap())
	assert.Equal(t, map[string]any{}, val.ToMap(defaultVisitor))
	assert.Equal(t, map[string]any{}, val.ToMap(customVisitor))

	// Object {"foo": "bar"}
	str := Value{typ: String, value: "bar"}
	val = Value{typ: Object, value: map[string]Value{"foo": str}}
	assert.Equal(t, map[string]any{"foo": "bar"}, val.ToMap())
	assert.Equal(t, map[string]any{"foo": "bar"}, val.ToMap(defaultVisitor))
	assert.Equal(t, map[string]any{"foo": "barbar"}, val.ToMap(customVisitor))
	assert.Equal(t, map[string]any{"foo": "bar"}, val.ToMap(intVisitor))
	assert.Equal(t, map[string]any{"foo": "bar"}, val.ToMap(floatVisitor))

	// Object {"foo": 5}
	num := Value{typ: Number, value: big.NewRat(5, 1)}
	val = Value{typ: Object, value: map[string]Value{"foo": num}}
	assert.Equal(t, map[string]any{"foo": big.NewRat(5, 1)}, val.ToMap())
	assert.Equal(t, map[string]any{"foo": big.NewRat(5, 1)}, val.ToMap(defaultVisitor))
	assert.Equal(t, map[string]any{"foo": big.NewRat(10, 1)}, val.ToMap(customVisitor))
	assert.Equal(t, map[string]any{"foo": int64(5)}, val.ToMap(intVisitor))
	assert.Equal(t, map[string]any{"foo": 5.0}, val.ToMap(floatVisitor))

	// Object {"foo": true}
	val = Value{typ: Object, value: map[string]Value{"foo": TrueValue}}
	assert.Equal(t, map[string]any{"foo": true}, val.ToMap())
	assert.Equal(t, map[string]any{"foo": true}, val.ToMap(defaultVisitor))
	assert.Equal(t, map[string]any{"foo": false}, val.ToMap(customVisitor))
	assert.Equal(t, map[string]any{"foo": true}, val.ToMap(intVisitor))
	assert.Equal(t, map[string]any{"foo": true}, val.ToMap(floatVisitor))

	// Object {"foo": nil}
	val = Value{typ: Object, value: map[string]Value{"foo": NullValue}}
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
		defaultVisitor = DefaultConversionVisitor
		customVisitor  = ConversionVisitor(
			func(str string) string { return str + str },
			func(num any) *big.Rat {
				num_ := num.(*big.Rat)
				res := big.NewRat(0, 1)
				res.Add(num_, num_)
				return res
			},
			func(b bool) bool { return !b },
		)
		intVisitor = ConversionVisitor(
			funcs.Passthrough[string],
			func(num any) int64 { return conv.BigRatToInt64(num.(*big.Rat)) },
			funcs.Passthrough[bool],
		)
		floatVisitor = ConversionVisitor(
			funcs.Passthrough[string],
			func(num any) float64 { return conv.BigRatToFloat64(num.(*big.Rat)) },
			funcs.Passthrough[bool],
		)
		val = Value{typ: Array, value: []Value{}}
	)
	assert.Equal(t, []any{}, val.ToSlice())
	assert.Equal(t, []any{}, val.ToSlice(defaultVisitor))
	assert.Equal(t, []any{}, val.ToSlice(customVisitor))

	// Array ["bar"]
	str := Value{typ: String, value: "bar"}
	val = Value{typ: Array, value: []Value{str}}
	assert.Equal(t, []any{"bar"}, val.ToSlice())
	assert.Equal(t, []any{"bar"}, val.ToSlice(defaultVisitor))
	assert.Equal(t, []any{"barbar"}, val.ToSlice(customVisitor))
	assert.Equal(t, []any{"bar"}, val.ToSlice(intVisitor))
	assert.Equal(t, []any{"bar"}, val.ToSlice(floatVisitor))

	// Array [5]
	num := Value{typ: Number, value: big.NewRat(5, 1)}
	val = Value{typ: Array, value: []Value{num}}
	assert.Equal(t, []any{big.NewRat(5, 1)}, val.ToSlice())
	assert.Equal(t, []any{big.NewRat(5, 1)}, val.ToSlice(defaultVisitor))
	assert.Equal(t, []any{big.NewRat(10, 1)}, val.ToSlice(customVisitor))
	assert.Equal(t, []any{int64(5)}, val.ToSlice(intVisitor))
	assert.Equal(t, []any{5.0}, val.ToSlice(floatVisitor))

	// Array [true]
	val = Value{typ: Array, value: []Value{TrueValue}}
	assert.Equal(t, []any{true}, val.ToSlice())
	assert.Equal(t, []any{true}, val.ToSlice(defaultVisitor))
	assert.Equal(t, []any{false}, val.ToSlice(customVisitor))
	assert.Equal(t, []any{true}, val.ToSlice(intVisitor))
	assert.Equal(t, []any{true}, val.ToSlice(floatVisitor))

	// Array [nil]
	val = Value{typ: Array, value: []Value{NullValue}}
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
	val := Value{typ: Number, value: big.NewRat(5, 1)}
	assert.Equal(t, int64(5), conv.BigRatToInt64(val.AsBigRat()))

	funcs.TryTo(
		func() {
			NullValue.AsBigRat()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotNumber, e)
		},
	)
}

func TestToFloat(t *testing.T) {
	val := Value{typ: Number, value: big.NewRat(5, 1)}
	assert.Equal(t, 5.0, conv.BigRatToFloat64(val.AsBigRat()))

	funcs.TryTo(
		func() {
			NullValue.AsBigRat()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotNumber, e)
		},
	)
}
