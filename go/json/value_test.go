package json

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"
	"testing"

	"github.com/bantling/micro/go/funcs"
	"github.com/bantling/micro/go/util"
	"github.com/stretchr/testify/assert"
)

func assertObject(t *testing.T, e map[string]JSONValue, a JSONValue) {
	assert.Equal(t, JSONValue{typ: Object, value: e}, a)
}

func assertArray(t *testing.T, e []JSONValue, a JSONValue) {
	assert.Equal(t, JSONValue{typ: Array, value: e}, a)
}

func assertString(t *testing.T, e string, a JSONValue) {
	assert.Equal(t, JSONValue{typ: String, value: e}, a)
}

func assertNumber(t *testing.T, e *big.Float, a JSONValue) {
	assert.Equal(t, JSONValue{typ: Number, value: e}, a)
}

func assertBoolean(t *testing.T, e bool, a JSONValue) {
	assert.Equal(t, JSONValue{typ: Boolean, value: e}, a)
}

func assertNull(t *testing.T, a JSONValue) {
	assert.Equal(t, NullValue, a)
}

func TestFromNumberInternal(t *testing.T) {
	assertNumber(t, util.IntToBigFloat(1), fromNumberInternal(int(1)))
	assertNumber(t, util.IntToBigFloat(2), fromNumberInternal(int8(2)))
	assertNumber(t, util.IntToBigFloat(3), fromNumberInternal(int16(3)))
	assertNumber(t, util.IntToBigFloat(4), fromNumberInternal(int32(4)))
	assertNumber(t, util.IntToBigFloat(5), fromNumberInternal(int64(5)))

	assertNumber(t, util.UintToBigFloat(uint64(1)), fromNumberInternal(uint(1)))
	assertNumber(t, util.UintToBigFloat(uint64(2)), fromNumberInternal(uint8(2)))
	assertNumber(t, util.UintToBigFloat(uint64(3)), fromNumberInternal(uint16(3)))
	assertNumber(t, util.UintToBigFloat(uint64(4)), fromNumberInternal(uint32(4)))
	assertNumber(t, util.UintToBigFloat(uint64(5)), fromNumberInternal(uint64(5)))

	assertNumber(t, util.FloatToBigFloat(1.25), fromNumberInternal(float32(1.25)))
	assertNumber(t, util.FloatToBigFloat(2.5), fromNumberInternal(float64(2.5)))

	assertNumber(t, big.NewFloat(3.0), fromNumberInternal(util.IntToBigInt(3)))
	assertNumber(t, big.NewFloat(3.75), fromNumberInternal(util.FloatToBigFloat(3.75)))

	// fromNumberInternal accepts type any, so explicit conversion required
	assertNumber(t, util.FloatToBigFloat(4.25), fromNumberInternal(NumberString("4.25")))

	// Any other type results in invalid zero value
	assert.Equal(t, JSONValue{}, fromNumberInternal(""))
}

func TestFromValue(t *testing.T) {
	// Object
	assertObject(t, map[string]JSONValue{"foo": {typ: String, value: "bar"}}, FromValue(map[string]any{"foo": "bar"}))

	// Array
	assertArray(t, []JSONValue{{typ: String, value: "bar"}}, FromValue([]any{"bar"}))

	// String
	assertString(t, "bar", FromValue("bar"))

	// Number - int
	assertNumber(t, util.IntToBigFloat(1), FromValue(int(1)))
	assertNumber(t, util.IntToBigFloat(2), FromValue(int8(2)))
	assertNumber(t, util.IntToBigFloat(3), FromValue(int16(3)))
	assertNumber(t, util.IntToBigFloat(4), FromValue(int32(4)))
	assertNumber(t, util.IntToBigFloat(5), FromValue(int64(5)))

	// Number - uint
	assertNumber(t, util.IntToBigFloat(1), FromValue(uint(1)))
	assertNumber(t, util.IntToBigFloat(2), FromValue(uint8(2)))
	assertNumber(t, util.IntToBigFloat(3), FromValue(uint16(3)))
	assertNumber(t, util.IntToBigFloat(4), FromValue(uint32(4)))
	assertNumber(t, util.IntToBigFloat(5), FromValue(uint64(5)))

	// Number - float
	assertNumber(t, util.FloatToBigFloat(1.25), FromValue(float32(1.25)))
	assertNumber(t, util.FloatToBigFloat(2.5), FromValue(float64(2.5)))

	// Number - *big.Int, *big.Float
	assertNumber(t, big.NewFloat(3.0), FromValue(util.IntToBigInt(3)))
	assertNumber(t, big.NewFloat(3.75), FromValue(util.FloatToBigFloat(3.75)))

	// Number - NumberString
	// fromValue accepts type any, so explicit conversion required
	assertNumber(t, big.NewFloat(4.25), FromValue(NumberString("4.25")))

	// Boolean - true
	assertBoolean(t, true, FromValue(true))

	// Boolean - false
	assertBoolean(t, false, FromValue(false))

	// Null
	assertNull(t, FromValue(nil))
}

func TestFromMap(t *testing.T) {
	assertObject(t, map[string]JSONValue{"foo": {typ: String, value: "bar"}}, FromMap(map[string]any{"foo": "bar"}))
}

func TestFromSlice(t *testing.T) {
	assertArray(t, []JSONValue{{typ: String, value: "bar"}}, FromSlice([]any{"bar"}))
}

func TestFromString(t *testing.T) {
	assertString(t, "bar", FromString("bar"))
}

func TestFromSignedInt(t *testing.T) {
	assertNumber(t, util.IntToBigFloat(1), FromSignedInt(int(1)))
	assertNumber(t, util.IntToBigFloat(2), FromSignedInt(int8(2)))
	assertNumber(t, util.IntToBigFloat(3), FromSignedInt(int16(3)))
	assertNumber(t, util.IntToBigFloat(4), FromSignedInt(int32(4)))
	assertNumber(t, util.IntToBigFloat(5), FromSignedInt(int64(5)))
}

func TestFromUnsignedInt(t *testing.T) {
	assertNumber(t, util.IntToBigFloat(1), FromUnsignedInt(uint(1)))
	assertNumber(t, util.IntToBigFloat(2), FromUnsignedInt(uint8(2)))
	assertNumber(t, util.IntToBigFloat(3), FromUnsignedInt(uint16(3)))
	assertNumber(t, util.IntToBigFloat(4), FromUnsignedInt(uint32(4)))
	assertNumber(t, util.IntToBigFloat(5), FromUnsignedInt(uint64(5)))
}

func TestFromFloat(t *testing.T) {
	assertNumber(t, util.FloatToBigFloat(1.25), FromFloat(float32(1.25)))
	assertNumber(t, util.FloatToBigFloat(2.5), FromFloat(float64(2.5)))
}

func TestFromBigInt(t *testing.T) {
	assertNumber(t, util.IntToBigFloat(3), FromBigInt(util.IntToBigInt(3)))
}

func TestFromBigFloat(t *testing.T) {
	assertNumber(t, util.FloatToBigFloat(3.75), FromBigFloat(util.FloatToBigFloat(3.75)))
}

func TestFromNumberString(t *testing.T) {
	// fromNumberString accepts type NumberString, so implicit conversion allowed
	assertNumber(t, util.FloatToBigFloat(4.25), FromNumberString("4.25"))
}

func TestFromNumber(t *testing.T) {
	// Number - int
	assertNumber(t, util.IntToBigFloat(1), FromNumber(int(1)))
	assertNumber(t, util.IntToBigFloat(2), FromNumber(int8(2)))
	assertNumber(t, util.IntToBigFloat(3), FromNumber(int16(3)))
	assertNumber(t, util.IntToBigFloat(4), FromNumber(int32(4)))
	assertNumber(t, util.IntToBigFloat(5), FromNumber(int64(5)))

	// Number - uint
	assertNumber(t, util.IntToBigFloat(1), FromNumber(uint(1)))
	assertNumber(t, util.IntToBigFloat(2), FromNumber(uint8(2)))
	assertNumber(t, util.IntToBigFloat(3), FromNumber(uint16(3)))
	assertNumber(t, util.IntToBigFloat(4), FromNumber(uint32(4)))
	assertNumber(t, util.IntToBigFloat(5), FromNumber(uint64(5)))

	// Number - float
	assertNumber(t, util.FloatToBigFloat(1.25), FromNumber(float32(1.25)))
	assertNumber(t, util.FloatToBigFloat(2.5), FromNumber(float64(2.5)))

	// Number - *big.Int, *big.Float
	assertNumber(t, big.NewFloat(3.0), FromNumber(util.IntToBigInt(3)))
	assertNumber(t, big.NewFloat(3.75), FromNumber(util.FloatToBigFloat(3.75)))

	// Number - NumberString
	// FromNumber accepts type any, so explicit conversion required
	assertNumber(t, big.NewFloat(4.25), FromNumber(NumberString("4.25")))
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

func TestAsBigInt(t *testing.T) {
	val := FromNumber(1)
	assert.Equal(t, big.NewInt(1), val.AsBigInt())

	funcs.TryTo(
		func() {
			NullValue.AsBigInt()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, ErrNotNumber, e)
		},
	)
}

func TestAsBigFloat(t *testing.T) {
	val := FromNumber(1)
	assert.Equal(t, val.value, val.AsBigFloat())

	funcs.TryTo(
		func() {
			NullValue.AsBigFloat()
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
	assert.True(t, NullValue.IsNull())
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
