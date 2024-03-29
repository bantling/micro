package write

// SPDX-License-Identifier: Apache-2.0

import (
	gojson "encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/encoding/json"
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/io"
	"github.com/bantling/micro/io/writer"
	"github.com/stretchr/testify/assert"
)

func TestWrite_(t *testing.T) {
	var str strings.Builder

	assert.Nil(t, Write(json.MustToValue(map[string]any{}), writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, "{}", str.String())

	str.Reset()
	assert.Nil(t, Write(json.MustToValue([]any{}), writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, "[]", str.String())

	str.Reset()
	assert.Nil(t, Write(json.MustToValue("foo"), writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, `"foo"`, str.String())

	str.Reset()
	assert.Nil(t, Write(json.MustToValue(json.NumberString("1")), writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, "1", str.String())

	str.Reset()
	assert.Nil(t, Write(json.TrueValue, writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, "true", str.String())

	str.Reset()
	assert.Nil(t, Write(json.NullValue, writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, "null", str.String())

	var err = fmt.Errorf("died")
	funcs.TryTo(
		func() {
			MustWrite(json.FalseValue, writer.OfIOWriterAsRunes(io.NewErrorWriter(0, err)))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, err, e.(error)) },
	)
}

func TestWriteObject_(t *testing.T) {
	var (
		str strings.Builder
		m   = map[string]any{
			"obj": map[string]any{"foo": "bar"},
			"arr": []any{"foo"},
			"str": "foo",
			"num": json.NumberString("1"),
			"bln": false,
			"nul": nil,
		}
	)

	assert.Nil(t, Write(json.MustToValue(m), writer.OfIOWriterAsRunes(&str)))

	// Can't rely on map ordering in string result. We could use our parser, but using go built in parser is a better idea.
	// Note that go parses a number as a float64 when using a map.
	var mc map[string]any
	gojson.Unmarshal([]byte(str.String()), &mc)
	mc["num"] = json.NumberString(conv.FloatToString(mc["num"].(float64)))

	assert.Equal(t, mc, m)

	// Test errors
	err := fmt.Errorf("An error")

	// Fail to write opening {
	w := io.NewErrorWriter(0, err)
	assert.Equal(t, err, Write(json.MustToValue(m), writer.OfIOWriterAsRunes(w)))

	// Fail to write first key
	w = io.NewErrorWriter(1, err)
	assert.Equal(t, err, Write(json.MustToValue(m), writer.OfIOWriterAsRunes(w)))

	// Fail to write first value
	w = io.NewErrorWriter(7, err)
	assert.Equal(t, err, Write(json.MustToValue(map[string]any{"foo": "bar"}), writer.OfIOWriterAsRunes(w)))

	funcs.TryTo(
		func() {
			MustWriteObject(json.MustToValue(map[string]any{"foo": "bar"}), writer.OfIOWriterAsRunes(io.NewErrorWriter(0, err)))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, err, e.(error)) },
	)
}

func TestWriteArray_(t *testing.T) {
	var (
		str strings.Builder
		s   = []any{
			map[string]any{"foo": "bar"},
			[]any{"foo"},
			"foo",
			json.NumberString("1"),
			false,
			nil,
		}
	)

	assert.Nil(t, Write(json.MustToValue(s), writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, `[{"foo":"bar"},["foo"],"foo",1,false,null]`, str.String())

	// Test errors
	err := fmt.Errorf("An error")

	// Fail to write opening [
	w := io.NewErrorWriter(0, err)
	assert.Equal(t, err, Write(json.MustToValue(s), writer.OfIOWriterAsRunes(w)))

	// Fail to write first comma
	w = io.NewErrorWriter(14, err)
	assert.Equal(t, err, Write(json.MustToValue(s), writer.OfIOWriterAsRunes(w)))

	// Fail to write second value
	w = io.NewErrorWriter(15, err)
	assert.Equal(t, err, Write(json.MustToValue(s), writer.OfIOWriterAsRunes(w)))

	funcs.TryTo(
		func() {
			MustWriteArray(json.MustToValue(s), writer.OfIOWriterAsRunes(io.NewErrorWriter(0, err)))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, err, e.(error)) },
	)
}
