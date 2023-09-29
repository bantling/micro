package json

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/union"
	"github.com/stretchr/testify/assert"
)

func assertObject(t *testing.T, e map[string]Value, a union.Result[Value]) {
	assert.Equal(t, union.OfResult(Value{typ: Object, val: union.Of4T[map[string]Value, []Value, string, bool](e)}), a)
}

func assertArray(t *testing.T, e []Value, a union.Result[Value]) {
	assert.Equal(t, union.OfResult(Value{typ: Array, val: union.Of4U[map[string]Value, []Value, string, bool](e)}), a)
}

func assertString(t *testing.T, e string, a union.Result[Value]) {
	assert.Equal(t, union.OfResult(Value{typ: String, val: union.Of4V[map[string]Value, []Value, string, bool](e)}), a)
}

func assertNumber(t *testing.T, e string, a union.Result[Value]) {
	assert.Equal(t, union.OfResult(Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool](e)}), a)
}

func assertBoolean(t *testing.T, e bool, a union.Result[Value]) {
	assert.Equal(t, union.OfResult(Value{typ: Boolean, val: union.Of4W[map[string]Value, []Value, string, bool](e)}), a)
}

func assertNull(t *testing.T, a union.Result[Value]) {
	assert.Equal(t, union.OfResult(NullValue), a)
}

func assertError(t *testing.T, e error, a union.Result[Value]) {
	assert.Equal(t, union.OfError[Value](e), a)
}

func TestString_(t *testing.T) {
	assert.Equal(t, "Object", fmt.Sprintf("%s", Object))
	assert.Equal(t, "Array", fmt.Sprintf("%s", Array))
	assert.Equal(t, "String", fmt.Sprintf("%s", String))
	assert.Equal(t, "Number", fmt.Sprintf("%s", Number))
	assert.Equal(t, "Boolean", fmt.Sprintf("%s", Boolean))
	assert.Equal(t, "Null", fmt.Sprintf("%s", Null))
}

func TestToValue_(t *testing.T) {
	// Object
	assertObject(t, map[string]Value{"foo": StringToValue("bar")}, union.OfResultError(ToValue(map[string]any{"foo": "bar"})))
	assertObject(t, map[string]Value{"foo": StringToValue("bar")}, union.OfResult(MustToValue(map[string]any{"foo": "bar"})))
	assertObject(t, map[string]Value{"foo": StringToValue("bar")}, union.OfResult(MustMapToValue(map[string]any{"foo": "bar"})))

	// Array
	assertArray(t, []Value{StringToValue("bar")}, union.OfResultError(ToValue([]any{"bar"})))
	assertArray(t, []Value{StringToValue("bar")}, union.OfResult(MustSliceToValue([]any{"bar"})))

	// String
	assertString(t, "bar", union.OfResultError(ToValue("bar")))

	// Number - int
	assertNumber(t, "1", union.OfResultError(ToValue(int(1))))
	assertNumber(t, "1", union.OfResult(MustNumberToValue(int(1))))

	assertNumber(t, "2", union.OfResultError(ToValue(int8(2))))
	assertNumber(t, "3", union.OfResultError(ToValue(int16(3))))
	assertNumber(t, "4", union.OfResultError(ToValue(int32(4))))
	assertNumber(t, "5", union.OfResultError(ToValue(int64(5))))

	// Number - uint
	assertNumber(t, "1", union.OfResultError(ToValue(uint(1))))
	assertNumber(t, "2", union.OfResultError(ToValue(uint8(2))))
	assertNumber(t, "3", union.OfResultError(ToValue(uint16(3))))
	assertNumber(t, "4", union.OfResultError(ToValue(uint32(4))))
	assertNumber(t, "5", union.OfResultError(ToValue(uint64(5))))

	// Number - float
	assertNumber(t, "1.25", union.OfResultError(ToValue(float32(1.25))))
	assertNumber(t, "2.5", union.OfResultError(ToValue(float64(2.5))))

	// Number - *big.Int, *big.Float, *big.Rat
	var (
		bi *big.Int
		bf *big.Float
		br *big.Rat
	)
	conv.IntToBigInt(3, &bi)
	assertNumber(t, "3", union.OfResultError(ToValue(bi)))
	conv.FloatToBigFloat(3.75, &bf)
	assertNumber(t, "3.75", union.OfResultError(ToValue(bf)))
	conv.FloatToBigRat(4.25, &br)
	assertNumber(t, "4.25", union.OfResultError(ToValue(br)))

	// Number - NumberString
	// fromValue accepts type any, so explicit conversion required
	assertNumber(t, "5.75", union.OfResultError(ToValue(NumberString("5.75"))))

	// Boolean - true
	assertBoolean(t, true, union.OfResultError(ToValue(true)))

	// Boolean - false
	assertBoolean(t, false, union.OfResultError(ToValue(false)))

	// Null
	assertNull(t, union.OfResultError(ToValue((*big.Int)(nil))))

	// Map failure
	var c = make(chan bool)
	assert.Equal(t, union.OfError[Value](fmt.Errorf("chan bool cannot be converted to string")), union.OfResultError(ToValue(map[string]any{"foo": c})))

	// Slice failure
	assert.Equal(t, union.OfError[Value](fmt.Errorf("chan bool cannot be converted to string")), union.OfResultError(ToValue([]any{c})))
}

func TestMapToValue_(t *testing.T) {
	assertObject(t, map[string]Value{"foo": StringToValue("bar")}, union.OfResultError(MapToValue(map[string]any{"foo": "bar"})))
	assertObject(t, map[string]Value{"foo": StringToValue("bar")}, union.OfResultError(MapToValue(map[string]Value{"foo": StringToValue("bar")})))
}

func TestSliceToValue_(t *testing.T) {
	assertArray(t, []Value{StringToValue("bar")}, union.OfResultError(SliceToValue([]any{"bar"})))
	assertArray(t, []Value{StringToValue("bar")}, union.OfResultError(SliceToValue([]Value{StringToValue("bar")})))
}

func TestStringToValue_(t *testing.T) {
	assertString(t, "bar", union.OfResult(StringToValue("bar")))
}

func TestNumberToValue_(t *testing.T) {
	assertNumber(t, "123", union.OfResultError(NumberToValue(123)))
	assertNumber(t, "1", union.OfResultError(NumberToValue(int(1))))
	assertNumber(t, "2", union.OfResultError(NumberToValue(int8(2))))
	assertNumber(t, "3", union.OfResultError(NumberToValue(int16(3))))
	assertNumber(t, "4", union.OfResultError(NumberToValue(int32(4))))
	assertNumber(t, "5", union.OfResultError(NumberToValue(int64(5))))

	assertNumber(t, "1", union.OfResultError(NumberToValue(uint(1))))
	assertNumber(t, "2", union.OfResultError(NumberToValue(uint8(2))))
	assertNumber(t, "3", union.OfResultError(NumberToValue(uint16(3))))
	assertNumber(t, "4", union.OfResultError(NumberToValue(uint32(4))))
	assertNumber(t, "5", union.OfResultError(NumberToValue(uint64(5))))

	assertNumber(t, "1.25", union.OfResultError(NumberToValue(float32(1.25))))
	assertNumber(t, "2.5", union.OfResultError(NumberToValue(float64(2.5))))

	assertNumber(t, "1", union.OfResultError(NumberToValue(byte(1))))
	assertNumber(t, "2", union.OfResultError(NumberToValue(rune(2))))

	var (
		bi *big.Int
		bf *big.Float
		br *big.Rat
	)
	conv.IntToBigInt(3, &bi)
	assert.Equal(t, big.NewInt(3), bi)
	assertNumber(t, "3", union.OfResultError(NumberToValue(bi)))
	conv.FloatToBigFloat(3.75, &bf)
	assertNumber(t, "3.75", union.OfResultError(NumberToValue(bf)))
	conv.FloatToBigRat(4.25, &br)
	assertNumber(t, "4.25", union.OfResultError(NumberToValue(br)))

	assertNumber(t, "5.75", union.OfResultError(NumberToValue(NumberString("5.75"))))

	// NumberString cannot be empty or some other non-number value
	assertError(t, errNotNumber, union.OfResultError(NumberToValue(NumberString(""))))
	assertError(t, errNotNumber, union.OfResultError(NumberToValue(NumberString("foo"))))
	assertError(t, errNotNumber, union.OfResultError(NumberToValue(NumberString("-"))))
}

func TestBoolToValue_(t *testing.T) {
	assertBoolean(t, true, union.OfResult(BoolToValue(true)))
	assertBoolean(t, false, union.OfResult(BoolToValue(false)))
}

func TestType_(t *testing.T) {
	val, _ := MapToValue(map[string]any{})
	assert.Equal(t, Object, val.Type())

	val, _ = SliceToValue([]any{})
	assert.Equal(t, Array, val.Type())

	val = StringToValue("foo")
	assert.Equal(t, String, val.Type())

	val, _ = NumberToValue(123)
	assert.Equal(t, Number, val.Type())

	val = BoolToValue(true)
	assert.Equal(t, Boolean, val.Type())

	val = NullValue
	assert.Equal(t, Null, val.Type())
}

func TestAsMap_(t *testing.T) {
	mp := map[string]any{"foo": "bar"}
	val, _ := MapToValue(mp)
	assert.Equal(t, map[string]Value{"foo": StringToValue("bar")}, val.AsMap())

	funcs.TryTo(
		func() {
			NullValue.AsMap()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, errNotObject, e)
		},
	)
}

func TestAsSlice_(t *testing.T) {
	slc := []any{"foo"}
	val, _ := SliceToValue(slc)
	assert.Equal(t, []Value{StringToValue("foo")}, val.AsSlice())

	funcs.TryTo(
		func() {
			NullValue.AsSlice()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, errNotArray, e)
		},
	)
}

func TestAsString_(t *testing.T) {
	val := StringToValue("foo")
	assert.Equal(t, "foo", val.AsString())

	val, _ = NumberToValue(1)
	assert.Equal(t, "1", val.AsString())

	assert.Equal(t, "true", TrueValue.AsString())

	funcs.TryTo(
		func() {
			NullValue.AsString()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, errNotStringable, e)
		},
	)
}

func TestAsBigRat_(t *testing.T) {
	val, _ := NumberToValue(1)
	assert.Equal(t, NumberString("1"), val.AsNumber())

	funcs.TryTo(
		func() {
			NullValue.AsNumber()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, errNotNumber, e)
		},
	)
}

func TestAsBool_(t *testing.T) {
	assert.True(t, TrueValue.AsBool())

	funcs.TryTo(
		func() {
			NullValue.AsBool()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, errNotBoolean, e)
		},
	)
}

func TestIsNull_(t *testing.T) {
	assert.True(t, NullValue.IsNull())
	assert.False(t, TrueValue.IsNull())
}

func TestToAny_(t *testing.T) {
	m := map[string]any{"foo": "bar"}

	assert.Equal(t, m, funcs.MustValue(MapToValue(m)).ToAny())

	s := []any{"foo", "bar"}
	assert.Equal(t, s, funcs.MustValue(SliceToValue(s)).ToAny())

	assert.Equal(t, "str", StringToValue("str").ToAny())
	assert.Equal(t, NumberString("1"), funcs.MustValue(NumberToValue(NumberString("1"))).ToAny())
	assert.Equal(t, true, BoolToValue(true).ToAny())
	assert.Nil(t, NullValue.ToAny())
}

func TestToMap_(t *testing.T) {
	m := map[string]any{
		"map": map[string]any{"foo": "bar"},
		"slc": []any{"foo", "bar"},
		"str": "foo",
		"num": NumberString("1"),
		"bln": true,
		"nil": nil,
	}
	assert.Equal(t, m, funcs.MustValue(MapToValue(m)).ToMap())

	// Panic if value is not an object
	funcs.TryTo(
		func() {
			TrueValue.ToMap()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, errNotObject, e)
		},
	)
}

func TestToSlice_(t *testing.T) {
	s := []any{
		map[string]any{"foo": "bar"},
		[]any{"foo", "bar"},
		"str",
		NumberString("1"),
		true,
		nil,
	}
	assert.Equal(t, s, funcs.MustValue(SliceToValue(s)).ToSlice())

	// Panic if value is not an array
	funcs.TryTo(
		func() {
			TrueValue.ToSlice()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, errNotArray, e)
		},
	)
}

func TestIsDocument_(t *testing.T) {
	assert.True(t, funcs.MustValue(ToValue(map[string]any{})).IsDocument())
	assert.True(t, funcs.MustValue(ToValue([]any{})).IsDocument())
	assert.False(t, funcs.MustValue(ToValue("")).IsDocument())
	assert.False(t, funcs.MustValue(ToValue(0)).IsDocument())
	assert.False(t, funcs.MustValue(ToValue(true)).IsDocument())
	assert.False(t, NullValue.IsDocument())
}
