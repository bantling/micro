package json

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
	"fmt"
	"math/big"
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

func TestFromNumberInternal_(t *testing.T) {
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

	var (
		bi *big.Int
		bf *big.Float
		br *big.Rat
	)
	conv.IntToBigInt(3, &bi)
	assert.Equal(t, big.NewInt(3), bi)
	assertNumber(t, NumberString("3"), fromNumberInternal(bi))
	conv.FloatToBigFloat(3.75, &bf)
	assertNumber(t, NumberString("3.75"), fromNumberInternal(bf))
	conv.FloatToBigRat(4.25, &br)
	assertNumber(t, NumberString("4.25"), fromNumberInternal(br))

	// fromNumberInternal accepts type any, so explicit conversion required
	assertNumber(t, NumberString("5.75"), fromNumberInternal(NumberString("5.75")))

	// Any other type results in invalid zero value
	assert.Equal(t, Value{}, fromNumberInternal(""))
}

func TestFromValue_(t *testing.T) {
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
	var (
		bi *big.Int
		bf *big.Float
		br *big.Rat
	)
	conv.IntToBigInt(3, &bi)
	assertNumber(t, NumberString("3"), FromValue(bi))
	conv.FloatToBigFloat(3.75, &bf)
	assertNumber(t, NumberString("3.75"), FromValue(bf))
	conv.FloatToBigRat(4.25, &br)
	assertNumber(t, NumberString("4.25"), FromValue(br))

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
		func(e any) {
			assert.Equal(t, fmt.Errorf("A value of type complex128 is not a valid type to convert to a Value. Acceptable types are map[string]any, []any, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *big.Int, *big.Float, *big.Rat, NumberString, bool, and nil"), e)
		},
	)
}

func TestFromMap_(t *testing.T) {
	assertObject(t, map[string]Value{"foo": {typ: String, value: "bar"}}, FromMap(map[string]any{"foo": "bar"}))
}

func TestFromMapOfValue_(t *testing.T) {
	assertObject(t, map[string]Value{"foo": {typ: String, value: "bar"}}, FromMapOfValue(map[string]Value{"foo": {typ: String, value: "bar"}}))
}

func TestFromSlice_(t *testing.T) {
	assertArray(t, []Value{{typ: String, value: "bar"}}, FromSlice([]any{"bar"}))
}

func TestFromSliceOfValue_(t *testing.T) {
	assertArray(t, []Value{{typ: String, value: "bar"}}, FromSliceOfValue([]Value{{typ: String, value: "bar"}}))
}

func TestFromDocument_(t *testing.T) {
	assertObject(t, map[string]Value{"foo": {typ: String, value: "bar"}}, FromDocument(map[string]any{"foo": "bar"}))
	assertArray(t, []Value{{typ: String, value: "bar"}}, FromDocument([]any{"bar"}))
}

func TestFromDocumentOfValue_(t *testing.T) {
	assertObject(t, map[string]Value{"foo": {typ: String, value: "bar"}}, FromDocumentOfValue(map[string]Value{"foo": {typ: String, value: "bar"}}))
	assertArray(t, []Value{{typ: String, value: "bar"}}, FromDocumentOfValue([]Value{{typ: String, value: "bar"}}))
}

func TestFromString_(t *testing.T) {
	assertString(t, "bar", FromString("bar"))
}

func TestFromSignedInt_(t *testing.T) {
	assertNumber(t, NumberString("1"), FromSignedInt(int(1)))
	assertNumber(t, NumberString("2"), FromSignedInt(int8(2)))
	assertNumber(t, NumberString("3"), FromSignedInt(int16(3)))
	assertNumber(t, NumberString("4"), FromSignedInt(int32(4)))
	assertNumber(t, NumberString("5"), FromSignedInt(int64(5)))
}

func TestFromUnsignedInt_(t *testing.T) {
	assertNumber(t, NumberString("1"), FromUnsignedInt(uint(1)))
	assertNumber(t, NumberString("2"), FromUnsignedInt(uint8(2)))
	assertNumber(t, NumberString("3"), FromUnsignedInt(uint16(3)))
	assertNumber(t, NumberString("4"), FromUnsignedInt(uint32(4)))
	assertNumber(t, NumberString("5"), FromUnsignedInt(uint64(5)))
}

func TestFromFloat_(t *testing.T) {
	assertNumber(t, NumberString("1.25"), FromFloat(float32(1.25)))
	assertNumber(t, NumberString("2.5"), FromFloat(float64(2.5)))
}

func TestFromBigInt_(t *testing.T) {
	var bi *big.Int
	conv.IntToBigInt(3, &bi)
	assertNumber(t, NumberString("3"), FromBigInt(bi))
}

func TestFromBigFloat_(t *testing.T) {
	var bf *big.Float
	conv.FloatToBigFloat(3.75, &bf)
	assertNumber(t, NumberString("3.75"), FromBigFloat(bf))
}

func TestFromBigRat_(t *testing.T) {
	var br *big.Rat
	conv.FloatToBigRat(4.25, &br)
	assertNumber(t, NumberString("4.25"), FromBigRat(br))
}

func TestFromNumberString_(t *testing.T) {
	// fromNumberString accepts type NumberString, so implicit conversion allowed
	assertNumber(t, NumberString("5.75"), FromNumberString("5.75"))
}

func TestFromNumber_(t *testing.T) {
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
	var (
		bi *big.Int
		bf *big.Float
		br *big.Rat
	)
	conv.IntToBigInt(3, &bi)
	assertNumber(t, NumberString("3"), FromNumber(bi))
	conv.FloatToBigFloat(3.75, &bf)
	assertNumber(t, NumberString("3.75"), FromNumber(bf))
	conv.FloatToBigRat(4.25, &br)
	assertNumber(t, NumberString("4.25"), FromNumber(br))

	// Number - NumberString
	// FromNumber accepts type any, so explicit conversion required
	assertNumber(t, NumberString("5.75"), FromNumber(NumberString("5.75")))
}

func TestFromBool_(t *testing.T) {
	assertBoolean(t, true, FromBool(true))
	assertBoolean(t, false, FromBool(false))
}

func TestType_(t *testing.T) {
	val := FromMap(map[string]any{})
	assert.Equal(t, Object, val.Type())
}

func TestAsMap_(t *testing.T) {
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

func TestAsSlice_(t *testing.T) {
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

func TestAsString_(t *testing.T) {
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

func TestAsBigRat_(t *testing.T) {
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
	assert.Equal(t, m, FromMap(m).ToAny())

	s := []any{"foo", "bar"}
	assert.Equal(t, s, FromSlice(s).ToAny())

	assert.Equal(t, "str", FromString("str").ToAny())
	assert.Equal(t, NumberString("1"), FromNumberString("1").ToAny())
	assert.Equal(t, true, FromBool(true).ToAny())
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

func TestToSlice_(t *testing.T) {
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

func TestIsDocument_(t *testing.T) {
	assert.True(t, FromMap(map[string]any{}).IsDocument())
	assert.True(t, FromSlice([]any{}).IsDocument())
	assert.False(t, FromString("").IsDocument())
	assert.False(t, FromNumber(0).IsDocument())
	assert.False(t, FromBool(true).IsDocument())
	assert.False(t, NullValue.IsDocument())
}

func TestWrite_(t *testing.T) {
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

func TestWriteObject_(t *testing.T) {
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

func TestWriteArray_(t *testing.T) {
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
