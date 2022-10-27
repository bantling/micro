package json

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/bantling/micro/go/conv"
	"github.com/bantling/micro/go/funcs"
	"github.com/bantling/micro/go/util"
	"github.com/bantling/micro/go/writer"
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

func assertNumber(t *testing.T, e NumberString, a Value) {
	assert.Equal(t, Value{typ: Number, value: e}, a)
}

func assertBoolean(t *testing.T, e bool, a Value) {
	assert.Equal(t, Value{typ: Boolean, value: e}, a)
}

func assertNull(t *testing.T, a Value) {
	assert.Equal(t, NullValue, a)
}

func TestString(t *testing.T) {
	assert.Equal(t, "Object", fmt.Sprintf("%s", Object))
	assert.Equal(t, "Array", fmt.Sprintf("%s", Array))
	assert.Equal(t, "String", fmt.Sprintf("%s", String))
	assert.Equal(t, "Number", fmt.Sprintf("%s", Number))
	assert.Equal(t, "Boolean", fmt.Sprintf("%s", Boolean))
	assert.Equal(t, "Null", fmt.Sprintf("%s", Null))
}

func TestFromNumberInternal(t *testing.T) {
	assertNumber(t, NumberString("1"), fromNumberInternal(int(1)))
	assertNumber(t, NumberString("2"), fromNumberInternal(int8(2)))
	assertNumber(t, NumberString("3"), fromNumberInternal(int16(3)))
	assertNumber(t, NumberString("4"), fromNumberInternal(int32(4)))
	assertNumber(t, NumberString("5"), fromNumberInternal(int64(5)))

	assertNumber(t, NumberString("1"), fromNumberInternal(uint(1)))
	assertNumber(t, NumberString("2"), fromNumberInternal(uint8(2)))
	assertNumber(t, NumberString("3"), fromNumberInternal(uint16(3)))
	assertNumber(t, NumberString("4"), fromNumberInternal(uint32(4)))
	assertNumber(t, NumberString("5"), fromNumberInternal(uint64(5)))

	assertNumber(t, NumberString("1.25"), fromNumberInternal(float32(1.25)))
	assertNumber(t, NumberString("2.5"), fromNumberInternal(float64(2.5)))

	assertNumber(t, NumberString("3"), fromNumberInternal(conv.IntToBigInt(3)))
	assertNumber(t, NumberString("3.75"), fromNumberInternal(conv.FloatToBigFloat(3.75)))
	assertNumber(t, NumberString("4.25"), fromNumberInternal(conv.FloatToBigRat(4.25)))

	// fromNumberInternal accepts type any, so explicit conversion required
	assertNumber(t, NumberString("5.75"), fromNumberInternal(NumberString("5.75")))

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
	assertNumber(t, NumberString("1"), FromValue(int(1)))
	assertNumber(t, NumberString("2"), FromValue(int8(2)))
	assertNumber(t, NumberString("3"), FromValue(int16(3)))
	assertNumber(t, NumberString("4"), FromValue(int32(4)))
	assertNumber(t, NumberString("5"), FromValue(int64(5)))

	// Number - uint
	assertNumber(t, NumberString("1"), FromValue(uint(1)))
	assertNumber(t, NumberString("2"), FromValue(uint8(2)))
	assertNumber(t, NumberString("3"), FromValue(uint16(3)))
	assertNumber(t, NumberString("4"), FromValue(uint32(4)))
	assertNumber(t, NumberString("5"), FromValue(uint64(5)))

	// Number - float
	assertNumber(t, NumberString("1.25"), FromValue(float32(1.25)))
	assertNumber(t, NumberString("2.5"), FromValue(float64(2.5)))

	// Number - *big.Int, *big.Float, *big.Rat
	assertNumber(t, NumberString("3"), FromValue(conv.IntToBigInt(3)))
	assertNumber(t, NumberString("3.75"), FromValue(conv.FloatToBigFloat(3.75)))
	assertNumber(t, NumberString("4.25"), FromValue(conv.FloatToBigRat(4.25)))

	// Number - NumberString
	// fromValue accepts type any, so explicit conversion required
	assertNumber(t, NumberString("5.75"), FromValue(NumberString("5.75")))

	// Boolean - true
	assertBoolean(t, true, FromValue(true))

	// Boolean - false
	assertBoolean(t, false, FromValue(false))

	// Null
	assertNull(t, FromValue(nil))

	// Error
	funcs.TryTo(
		func() { FromValue((1 + 2i)) },
		func(e any) { assert.Equal(t, fmt.Errorf(errInvalidGoValueMsg, (1+2i)), e) },
	)
}

func TestFromMap(t *testing.T) {
	assertObject(t, map[string]Value{"foo": {typ: String, value: "bar"}}, FromMap(map[string]any{"foo": "bar"}))
}

func TestFromMapOfValue(t *testing.T) {
	assertObject(t, map[string]Value{"foo": {typ: String, value: "bar"}}, FromMapOfValue(map[string]Value{"foo": {typ: String, value: "bar"}}))
}

func TestFromSlice(t *testing.T) {
	assertArray(t, []Value{{typ: String, value: "bar"}}, FromSlice([]any{"bar"}))
}

func TestFromSliceOfValue(t *testing.T) {
	assertArray(t, []Value{{typ: String, value: "bar"}}, FromSliceOfValue([]Value{{typ: String, value: "bar"}}))
}

func TestFromDocument(t *testing.T) {
	assertObject(t, map[string]Value{"foo": {typ: String, value: "bar"}}, FromDocument(map[string]any{"foo": "bar"}))
	assertArray(t, []Value{{typ: String, value: "bar"}}, FromDocument([]any{"bar"}))
}

func TestFromDocumentOfValue(t *testing.T) {
	assertObject(t, map[string]Value{"foo": {typ: String, value: "bar"}}, FromDocumentOfValue(map[string]Value{"foo": {typ: String, value: "bar"}}))
	assertArray(t, []Value{{typ: String, value: "bar"}}, FromDocumentOfValue([]Value{{typ: String, value: "bar"}}))
}

func TestFromString(t *testing.T) {
	assertString(t, "bar", FromString("bar"))
}

func TestFromSignedInt(t *testing.T) {
	assertNumber(t, NumberString("1"), FromSignedInt(int(1)))
	assertNumber(t, NumberString("2"), FromSignedInt(int8(2)))
	assertNumber(t, NumberString("3"), FromSignedInt(int16(3)))
	assertNumber(t, NumberString("4"), FromSignedInt(int32(4)))
	assertNumber(t, NumberString("5"), FromSignedInt(int64(5)))
}

func TestFromUnsignedInt(t *testing.T) {
	assertNumber(t, NumberString("1"), FromUnsignedInt(uint(1)))
	assertNumber(t, NumberString("2"), FromUnsignedInt(uint8(2)))
	assertNumber(t, NumberString("3"), FromUnsignedInt(uint16(3)))
	assertNumber(t, NumberString("4"), FromUnsignedInt(uint32(4)))
	assertNumber(t, NumberString("5"), FromUnsignedInt(uint64(5)))
}

func TestFromFloat(t *testing.T) {
	assertNumber(t, NumberString("1.25"), FromFloat(float32(1.25)))
	assertNumber(t, NumberString("2.5"), FromFloat(float64(2.5)))
}

func TestFromBigInt(t *testing.T) {
	assertNumber(t, NumberString("3"), FromBigInt(conv.IntToBigInt(3)))
}

func TestFromBigFloat(t *testing.T) {
	assertNumber(t, NumberString("3.75"), FromBigFloat(conv.FloatToBigFloat(3.75)))
}

func TestFromBigRat(t *testing.T) {
	assertNumber(t, NumberString("4.25"), FromBigFloat(conv.FloatToBigFloat(4.25)))
}

func TestFromNumberString(t *testing.T) {
	// fromNumberString accepts type NumberString, so implicit conversion allowed
	assertNumber(t, NumberString("5.75"), FromNumberString("5.75"))
}

func TestFromNumber(t *testing.T) {
	// Number - int
	assertNumber(t, NumberString("1"), FromNumber(int(1)))
	assertNumber(t, NumberString("2"), FromNumber(int8(2)))
	assertNumber(t, NumberString("3"), FromNumber(int16(3)))
	assertNumber(t, NumberString("4"), FromNumber(int32(4)))
	assertNumber(t, NumberString("5"), FromNumber(int64(5)))

	// Number - uint
	assertNumber(t, NumberString("1"), FromNumber(uint(1)))
	assertNumber(t, NumberString("2"), FromNumber(uint8(2)))
	assertNumber(t, NumberString("3"), FromNumber(uint16(3)))
	assertNumber(t, NumberString("4"), FromNumber(uint32(4)))
	assertNumber(t, NumberString("5"), FromNumber(uint64(5)))

	// Number - float
	assertNumber(t, NumberString("1.25"), FromNumber(float32(1.25)))
	assertNumber(t, NumberString("2.5"), FromNumber(float64(2.5)))

	// Number - *big.Int, *big.Float, *big.Rat
	assertNumber(t, NumberString("3"), FromNumber(conv.IntToBigInt(3)))
	assertNumber(t, NumberString("3.75"), FromNumber(conv.FloatToBigFloat(3.75)))
	assertNumber(t, NumberString("4.25"), FromNumber(conv.FloatToBigRat(4.25)))

	// Number - NumberString
	// FromNumber accepts type any, so explicit conversion required
	assertNumber(t, NumberString("5.75"), FromNumber(NumberString("5.75")))
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
			assert.Equal(t, errNotObject, e)
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
			assert.Equal(t, errNotArray, e)
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
			assert.Equal(t, errNotStringable, e)
		},
	)
}

func TestAsBigRat(t *testing.T) {
	val := FromNumber(1)
	assert.Equal(t, NumberString("1"), val.AsNumberString())

	funcs.TryTo(
		func() {
			NullValue.AsNumberString()
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, errNotNumber, e)
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
			assert.Equal(t, errNotBoolean, e)
		},
	)
}

func TestIsNull(t *testing.T) {
	assert.True(t, NullValue.IsNull())
	assert.False(t, TrueValue.IsNull())
}

func TestToAny(t *testing.T) {
	m := map[string]any{"foo": "bar"}
	assert.Equal(t, m, FromMap(m).ToAny())

	s := []any{"foo", "bar"}
	assert.Equal(t, s, FromSlice(s).ToAny())

	assert.Equal(t, "str", FromString("str").ToAny())
	assert.Equal(t, NumberString("1"), FromNumberString("1").ToAny())
	assert.Equal(t, true, FromBool(true).ToAny())
	assert.Nil(t, NullValue.ToAny())
}

func TestToMap(t *testing.T) {
	m := map[string]any{
		"map": map[string]any{"foo": "bar"},
		"slc": []any{"foo", "bar"},
		"str": "foo",
		"num": NumberString("1"),
		"bln": true,
		"nil": nil,
	}
	assert.Equal(t, m, FromMap(m).ToMap())

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

func TestToSlice(t *testing.T) {
	s := []any{
		map[string]any{"foo": "bar"},
		[]any{"foo", "bar"},
		"str",
		NumberString("1"),
		true,
		nil,
	}
	assert.Equal(t, s, FromSlice(s).ToSlice())

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

func TestIsDocument(t *testing.T) {
	assert.True(t, FromMap(map[string]any{}).IsDocument())
	assert.True(t, FromSlice([]any{}).IsDocument())
	assert.False(t, FromString("").IsDocument())
	assert.False(t, FromNumber(0).IsDocument())
	assert.False(t, FromBool(true).IsDocument())
	assert.False(t, NullValue.IsDocument())
}

func TestWrite(t *testing.T) {
	var str strings.Builder

	assert.Nil(t, FromMap(map[string]any{}).Write(writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, "{}", str.String())

	str.Reset()
	assert.Nil(t, FromSlice([]any{}).Write(writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, "[]", str.String())

	str.Reset()
	assert.Nil(t, FromString("foo").Write(writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, `"foo"`, str.String())

	str.Reset()
	assert.Nil(t, FromNumberString("1").Write(writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, "1", str.String())

	str.Reset()
	assert.Nil(t, TrueValue.Write(writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, "true", str.String())

	str.Reset()
	assert.Nil(t, NullValue.Write(writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, "null", str.String())
}

func TestWriteObject(t *testing.T) {
	var (
		str strings.Builder
		m   = map[string]any{
			"obj": map[string]any{"foo": "bar"},
			"arr": []any{"foo"},
			"str": "foo",
			"num": NumberString("1"),
			"bln": false,
			"nul": nil,
		}
	)

	assert.Nil(t, FromMap(m).Write(writer.OfIOWriterAsRunes(&str)))

	// Can't rely on map ordering in string result, and can't use our parser since it imports this package making a cycle.
	// So use go built in JSON parser to parse string into a map struct. It parses number as a float64 when using a map.
	var mc map[string]any
	json.Unmarshal([]byte(str.String()), &mc)
	mc["num"] = NumberString(conv.FloatToString(mc["num"].(float64)))

	assert.Equal(t, mc, m)

	// Test errors
	err := fmt.Errorf("An error")

	// Fail to write opening {
	w := util.NewErrorWriter(0, err)
	assert.Equal(t, err, FromMap(m).Write(writer.OfIOWriterAsRunes(w)))

	// Fail to write first key
	w = util.NewErrorWriter(1, err)
	assert.Equal(t, err, FromMap(m).Write(writer.OfIOWriterAsRunes(w)))

	// Fail to write first value
	w = util.NewErrorWriter(7, err)
	assert.Equal(t, err, FromMap(map[string]any{"foo": "bar"}).Write(writer.OfIOWriterAsRunes(w)))
}

func TestWriteArray(t *testing.T) {
	var (
		str strings.Builder
		s   = []any{
			map[string]any{"foo": "bar"},
			[]any{"foo"},
			"foo",
			NumberString("1"),
			false,
			nil,
		}
	)

	assert.Nil(t, FromSlice(s).Write(writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, `[{"foo":"bar"},["foo"],"foo",1,false,null]`, str.String())

	// Test errors
	err := fmt.Errorf("An error")

	// Fail to write opening [
	w := util.NewErrorWriter(0, err)
	assert.Equal(t, err, FromSlice(s).Write(writer.OfIOWriterAsRunes(w)))

	// Fail to write first comma
	w = util.NewErrorWriter(14, err)
	assert.Equal(t, err, FromSlice(s).Write(writer.OfIOWriterAsRunes(w)))

	// Fail to write second value
	w = util.NewErrorWriter(15, err)
	assert.Equal(t, err, FromSlice(s).Write(writer.OfIOWriterAsRunes(w)))
}
