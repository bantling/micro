package write

// SPDX-License-Identifier: Apache-2.0

import (
	gojson "encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/json"
	"github.com/bantling/micro/util"
	"github.com/bantling/micro/writer"
	"github.com/stretchr/testify/assert"
)

func TestWrite_(t *testing.T) {
	var str strings.Builder

	assert.Nil(t, Write(json.FromMap(map[string]any{}), writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, "{}", str.String())

	str.Reset()
	assert.Nil(t, Write(json.FromSlice([]any{}), writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, "[]", str.String())

	str.Reset()
	assert.Nil(t, Write(json.FromString("foo"), writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, `"foo"`, str.String())

	str.Reset()
	assert.Nil(t, Write(json.FromNumberString("1"), writer.OfIOWriterAsRunes(&str)))
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
      MustWrite(json.FalseValue, writer.OfIOWriterAsRunes(util.NewErrorWriter(0, err)))
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

	assert.Nil(t, Write(json.FromMap(m), writer.OfIOWriterAsRunes(&str)))

	// Can't rely on map ordering in string result. We could use our parser, but using go built in parser is a better idea.
	// Note that go parses a number as a float64 when using a map.
	var mc map[string]any
	gojson.Unmarshal([]byte(str.String()), &mc)
	mc["num"] = json.NumberString(conv.FloatToString(mc["num"].(float64)))

	assert.Equal(t, mc, m)

	// Test errors
	err := fmt.Errorf("An error")

	// Fail to write opening {
	w := util.NewErrorWriter(0, err)
	assert.Equal(t, err, Write(json.FromMap(m), writer.OfIOWriterAsRunes(w)))

	// Fail to write first key
	w = util.NewErrorWriter(1, err)
	assert.Equal(t, err, Write(json.FromMap(m), writer.OfIOWriterAsRunes(w)))

	// Fail to write first value
	w = util.NewErrorWriter(7, err)
	assert.Equal(t, err, Write(json.FromMap(map[string]any{"foo": "bar"}), writer.OfIOWriterAsRunes(w)))

  funcs.TryTo(
    func() {
      MustWriteObject(json.FromMap(map[string]any{"foo": "bar"}), writer.OfIOWriterAsRunes(util.NewErrorWriter(0, err)))
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

	assert.Nil(t, Write(json.FromSlice(s), writer.OfIOWriterAsRunes(&str)))
	assert.Equal(t, `[{"foo":"bar"},["foo"],"foo",1,false,null]`, str.String())

	// Test errors
	err := fmt.Errorf("An error")

	// Fail to write opening [
	w := util.NewErrorWriter(0, err)
	assert.Equal(t, err, Write(json.FromSlice(s), writer.OfIOWriterAsRunes(w)))

	// Fail to write first comma
	w = util.NewErrorWriter(14, err)
	assert.Equal(t, err, Write(json.FromSlice(s), writer.OfIOWriterAsRunes(w)))

	// Fail to write second value
	w = util.NewErrorWriter(15, err)
	assert.Equal(t, err, Write(json.FromSlice(s), writer.OfIOWriterAsRunes(w)))

  funcs.TryTo(
    func() {
      MustWriteArray(json.FromSlice(s), writer.OfIOWriterAsRunes(util.NewErrorWriter(0, err)))
      assert.Fail(t, "Must die")
    },
    func(e any) { assert.Equal(t, err, e.(error)) },
  )
}
