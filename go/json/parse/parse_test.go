package parse

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bantling/micro/go/iter"
	"github.com/bantling/micro/go/json"
	"github.com/bantling/micro/go/stream"
	"github.com/bantling/micro/go/tuple"
	"github.com/bantling/micro/go/util"
	"github.com/stretchr/testify/assert"
)

func TestParseValue_(t *testing.T) {
	assert.Equal(t, tuple.Of2Error(json.FromMap(map[string]any{"foo": "bar"}), nil), tuple.Of2Error(parseValue(lexer(iter.OfStringAsRunes(`{"foo": "bar"}`)))))
	assert.Equal(t, tuple.Of2Error(json.FromSlice([]any{"bar"}), nil), tuple.Of2Error(parseValue(lexer(iter.OfStringAsRunes(`["bar"]`)))))
	assert.Equal(t, tuple.Of2Error(json.FromString("bar"), nil), tuple.Of2Error(parseValue(lexer(iter.OfStringAsRunes(`"bar"`)))))
	assert.Equal(t, tuple.Of2Error(json.FromNumberString("1.25"), nil), tuple.Of2Error(parseValue(lexer(iter.OfStringAsRunes("1.25")))))
	assert.Equal(t, tuple.Of2Error(json.TrueValue, nil), tuple.Of2Error(parseValue(lexer(iter.OfStringAsRunes("true")))))
	assert.Equal(t, tuple.Of2Error(json.NullValue, nil), tuple.Of2Error(parseValue(lexer(iter.OfStringAsRunes("null")))))

	// Array that returns a problem
	anErr := fmt.Errorf("An err")
	assert.Equal(t, tuple.Of2Error(json.Value{}, anErr), tuple.Of2Error(parseValue(lexer(iter.SetError(iter.OfStringAsRunes(`[`), anErr)))))

	// parseValue returns (invalid Value, nil) for tokens that cannot be a value - up to caller to return better error
	assert.Equal(t, tuple.Of2Error(json.Value{}, nil), tuple.Of2Error(parseValue(lexer(iter.OfStringAsRunes("}")))))
}

func TestParseObject_(t *testing.T) {
	assert.Equal(t, tuple.Of2Error(json.FromMap(map[string]any{}), nil), tuple.Of2Error(parseObject(lexer(iter.OfStringAsRunes(`{}`)))))
	assert.Equal(t, tuple.Of2Error(json.FromMap(map[string]any{"foo": "bar"}), nil), tuple.Of2Error(parseObject(lexer(iter.OfStringAsRunes(`{"foo": "bar"}`)))))

	var (
		zv    json.Value
		anErr = fmt.Errorf("An err")
	)

	// Case 1
	assert.Equal(t, tuple.Of2Error(zv, errObjectRequiresKeyOrBrace), tuple.Of2Error(parseObject(lexer(iter.OfStringAsRunes(`{`)))))
	// Case 2
	assert.Equal(t, tuple.Of2Error(zv, anErr), tuple.Of2Error(parseObject(lexer(iter.SetError(iter.OfStringAsRunes(`{`), anErr)))))
	// Case 3
	assert.Equal(t, tuple.Of2Error(zv, errObjectRequiresKeyOrBrace), tuple.Of2Error(parseObject(lexer(iter.OfStringAsRunes(`{{`)))))
	// Case 4
	assert.Equal(t, tuple.Of2Error(zv, fmt.Errorf("A JSON object cannot have duplicate key \"foo\"")), tuple.Of2Error(parseObject(lexer(iter.OfStringAsRunes(`{"foo": "bar","foo":"baz"`)))))
	// Case 5
	assert.Equal(t, tuple.Of2Error(zv, fmt.Errorf("The JSON object key \"foo\" just be followed by a colon")), tuple.Of2Error(parseObject(lexer(iter.OfStringAsRunes(`{"foo"`)))))
	// Case 6
	assert.Equal(t, tuple.Of2Error(zv, anErr), tuple.Of2Error(parseObject(lexer(iter.SetError(iter.OfStringAsRunes(`{"foo"`), anErr)))))
	// Case 7
	assert.Equal(t, tuple.Of2Error(zv, fmt.Errorf("The JSON object key \"foo\" just be followed by a colon")), tuple.Of2Error(parseObject(lexer(iter.OfStringAsRunes(`{"foo"{`)))))
	// Case 8
	assert.Equal(t, tuple.Of2Error(zv, fmt.Errorf("The JSON object key \"foo\" must be have a value that is an object, arrray, string, number, boolean, or null")), tuple.Of2Error(parseObject(lexer(iter.OfStringAsRunes(`{"foo":`)))))
	// Case 9
	assert.Equal(t, tuple.Of2Error(zv, anErr), tuple.Of2Error(parseObject(lexer(iter.SetError(iter.OfStringAsRunes(`{"foo":`), anErr)))))
	// Case 10
	assert.Equal(t, tuple.Of2Error(zv, anErr), tuple.Of2Error(parseObject(lexer(iter.SetError(iter.OfStringAsRunes(`{"foo":{`), anErr)))))
	// Case 11
	assert.Equal(t, tuple.Of2Error(zv, fmt.Errorf("The JSON key/value pair \"foo\" must be followed by a colon or closing brace")), tuple.Of2Error(parseObject(lexer(iter.OfStringAsRunes(`{"foo":1`)))))
	// Case 12 - need space after key value so that error occurs after successfully returning number
	assert.Equal(t, tuple.Of2Error(zv, anErr), tuple.Of2Error(parseObject(lexer(iter.SetError(iter.OfStringAsRunes(`{"foo":1 `), anErr)))))
	// Case 13
	assert.Equal(t, tuple.Of2Error(zv, fmt.Errorf("The JSON key/value pair \"foo\" must be followed by a colon or closing brace")), tuple.Of2Error(parseObject(lexer(iter.OfStringAsRunes(`{"foo":1{`)))))
}

func TestParseArray_(t *testing.T) {
	assert.Equal(t, tuple.Of2Error([]json.Value{json.FromString("bar")}, nil), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`["bar"]`))))))
	assert.Equal(t, tuple.Of2Error([]json.Value{json.FromString("foo"), json.FromString("bar")}, nil), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`["foo", "bar"]`))))))

	var (
		zv    []json.Value
		anErr = fmt.Errorf("An err")
	)

	// Case 1
	assert.Equal(t, tuple.Of2Error(zv, errArrayRequiresValueOrBracket), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`[`))))))
	// Case 2
	assert.Equal(t, tuple.Of2Error(zv, anErr), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.SetError(iter.OfStringAsRunes(`[`), anErr))))))
	// Case 3
	assert.Equal(t, tuple.Of2Error([]json.Value{}, nil), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`[]`))))))
	// Case 4
	assert.Equal(t, tuple.Of2Error(zv, anErr), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.SetError(iter.OfStringAsRunes(`[[`), anErr))))))
	// Case 5
	assert.Equal(t, tuple.Of2Error(zv, errArrayRequiresCommaOrBracket), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`[1`))))))
	// Case 6 - Need a space after value so that error occurs after successfully returning number
	assert.Equal(t, tuple.Of2Error(zv, anErr), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.SetError(iter.OfStringAsRunes(`[1 `), anErr))))))
	// Case 7
	assert.Equal(t, tuple.Of2Error(zv, errArrayRequiresValue), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`[1,`))))))
	// Case 8
	assert.Equal(t, tuple.Of2Error(zv, anErr), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.SetError(iter.OfStringAsRunes(`[1,`), anErr))))))
	// Case 9
	assert.Equal(t, tuple.Of2Error(zv, anErr), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.SetError(iter.OfStringAsRunes(`[1,[`), anErr))))))
	// Case 10 - ordinary success case
	assert.Equal(t, tuple.Of2Error([]json.Value{json.FromNumber(1)}, nil), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`[1]`))))))
	// Case 11
	assert.Equal(t, tuple.Of2Error(zv, errArrayRequiresCommaOrBracket), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`[1}`))))))
}

func TestIterate_(t *testing.T) {
	assert.Equal(t, tuple.Of2Error([]json.Value{json.FromMap(map[string]any{"foo": "bar"})}, nil), iter.Maybe(stream.ReduceToSlice(Iterate(strings.NewReader(`{"foo": "bar"}`)))))
	assert.Equal(t, tuple.Of2Error([]json.Value{json.FromMap(map[string]any{"foo": "bar"})}, nil), iter.Maybe(stream.ReduceToSlice(Iterate(strings.NewReader(`[{"foo": "bar"}]`)))))

	vals, err := stream.ReduceToSlice(Iterate(strings.NewReader(`[{"foo": "bar", "baz": 1}, ["fooey", 2], true, null]`))).Next()
	assert.Nil(t, err)
	assert.Equal(t, 4, len(vals))
	assert.Equal(t, json.FromMap(map[string]any{"foo": "bar", "baz": 1}), vals[0])
	assert.Equal(t, json.FromSlice([]any{"fooey", 2}), vals[1])
	assert.Equal(t, json.TrueValue, vals[2])
	assert.Equal(t, json.NullValue, vals[3])

	var (
		zv    json.Value
		anErr = fmt.Errorf("An err")
	)

	// Case 1
	assert.Equal(t, tuple.Of2Error(zv, errEmptyDocument), iter.Maybe(Iterate(strings.NewReader(``))))
	// Case 2
	assert.Equal(t, tuple.Of2Error(zv, anErr), iter.Maybe(Iterate(util.NewErrorReader([]byte(``), anErr))))
	// Case 3
	assert.Equal(t, tuple.Of2Error(zv, errObjectRequiresKeyOrBrace), iter.Maybe(Iterate(strings.NewReader(`{`))))
	// Case 4
	assert.Equal(t, tuple.Of2Error(zv, errObjectOrArrayRequired), iter.Maybe(Iterate(strings.NewReader(`:`))))
}

func TestParse_(t *testing.T) {
	assert.Equal(t, tuple.Of2Error(json.FromMap(map[string]any{"foo": "bar"}), nil), tuple.Of2Error(Parse(strings.NewReader(`{"foo": "bar"}`))))
	assert.Equal(t, tuple.Of2Error(json.FromSlice([]any{map[string]any{"foo": "bar"}}), nil), tuple.Of2Error(Parse(strings.NewReader(`[{"foo": "bar"}]`))))

	var (
		zv    json.Value
		anErr = fmt.Errorf("An err")
	)

	// Case 1
	assert.Equal(t, tuple.Of2Error(zv, errEmptyDocument), tuple.Of2Error(Parse(strings.NewReader(``))))
	// Case 2
	assert.Equal(t, tuple.Of2Error(zv, anErr), tuple.Of2Error(Parse(util.NewErrorReader([]byte(``), anErr))))
	// Case 3
	assert.Equal(t, tuple.Of2Error(zv, anErr), tuple.Of2Error(Parse(util.NewErrorReader([]byte(`{`), anErr))))
	// Case 4
	assert.Equal(t, tuple.Of2Error(zv, errObjectOrArrayRequired), tuple.Of2Error(Parse(strings.NewReader(`:`))))
}
