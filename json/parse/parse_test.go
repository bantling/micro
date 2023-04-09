package parse

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bantling/micro/iter"
	"github.com/bantling/micro/json"
	"github.com/bantling/micro/stream"
	"github.com/bantling/micro/union"
	"github.com/bantling/micro/util"
	"github.com/stretchr/testify/assert"
)

func TestParseValue_(t *testing.T) {
	assert.Equal(t, union.OfResult(json.FromMap(map[string]any{"foo": "bar"})), union.OfResultError(parseValue(lexer(iter.OfStringAsRunes(`{"foo": "bar"}`)))))
	assert.Equal(t, union.OfResult(json.FromSlice([]any{"bar"})), union.OfResultError(parseValue(lexer(iter.OfStringAsRunes(`["bar"]`)))))
	assert.Equal(t, union.OfResult(json.FromString("bar")), union.OfResultError(parseValue(lexer(iter.OfStringAsRunes(`"bar"`)))))
	assert.Equal(t, union.OfResult(json.FromNumberString("1.25")), union.OfResultError(parseValue(lexer(iter.OfStringAsRunes("1.25")))))
	assert.Equal(t, union.OfResult(json.TrueValue), union.OfResultError(parseValue(lexer(iter.OfStringAsRunes("true")))))
	assert.Equal(t, union.OfResult(json.NullValue), union.OfResultError(parseValue(lexer(iter.OfStringAsRunes("null")))))

	// Array that returns a problem
	anErr := fmt.Errorf("An err")
	assert.Equal(t, union.OfError[json.Value](anErr), union.OfResultError(parseValue(lexer(iter.SetError(iter.OfStringAsRunes(`[`), anErr)))))

	// parseValue returns (invalid Value, nil) for tokens that cannot be a value - up to caller to return better error
	assert.Equal(t, union.OfResult(json.Value{}), union.OfResultError(parseValue(lexer(iter.OfStringAsRunes("}")))))
}

func TestParseObject_(t *testing.T) {
	assert.Equal(t, union.OfResult(json.FromMap(map[string]any{})), union.OfResultError(parseObject(lexer(iter.OfStringAsRunes(`{}`)))))
	assert.Equal(t, union.OfResult(json.FromMap(map[string]any{"foo": "bar"})), union.OfResultError(parseObject(lexer(iter.OfStringAsRunes(`{"foo": "bar"}`)))))

	anErr := fmt.Errorf("An err")

	// Case 1
	assert.Equal(t, union.OfError[json.Value](errObjectRequiresKeyOrBrace), union.OfResultError(parseObject(lexer(iter.OfStringAsRunes(`{`)))))
	// Case 2
	assert.Equal(t, union.OfError[json.Value](anErr), union.OfResultError(parseObject(lexer(iter.SetError(iter.OfStringAsRunes(`{`), anErr)))))
	// Case 3
	assert.Equal(t, union.OfError[json.Value](errObjectRequiresKeyOrBrace), union.OfResultError(parseObject(lexer(iter.OfStringAsRunes(`{{`)))))
	// Case 4
	assert.Equal(t, union.OfError[json.Value](fmt.Errorf("A JSON object cannot have duplicate key \"foo\"")), union.OfResultError(parseObject(lexer(iter.OfStringAsRunes(`{"foo": "bar","foo":"baz"`)))))
	// Case 5
	assert.Equal(t, union.OfError[json.Value](fmt.Errorf("The JSON object key \"foo\" just be followed by a colon")), union.OfResultError(parseObject(lexer(iter.OfStringAsRunes(`{"foo"`)))))
	// Case 6
	assert.Equal(t, union.OfError[json.Value](anErr), union.OfResultError(parseObject(lexer(iter.SetError(iter.OfStringAsRunes(`{"foo"`), anErr)))))
	// Case 7
	assert.Equal(t, union.OfError[json.Value](fmt.Errorf("The JSON object key \"foo\" just be followed by a colon")), union.OfResultError(parseObject(lexer(iter.OfStringAsRunes(`{"foo"{`)))))
	// Case 8
	assert.Equal(t, union.OfError[json.Value](fmt.Errorf("The JSON object key \"foo\" must be have a value that is an object, arrray, string, number, boolean, or null")), union.OfResultError(parseObject(lexer(iter.OfStringAsRunes(`{"foo":`)))))
	// Case 9
	assert.Equal(t, union.OfError[json.Value](anErr), union.OfResultError(parseObject(lexer(iter.SetError(iter.OfStringAsRunes(`{"foo":`), anErr)))))
	// Case 10
	assert.Equal(t, union.OfError[json.Value](anErr), union.OfResultError(parseObject(lexer(iter.SetError(iter.OfStringAsRunes(`{"foo":{`), anErr)))))
	// Case 11
	assert.Equal(t, union.OfError[json.Value](fmt.Errorf("The JSON key/value pair \"foo\" must be followed by a colon or closing brace")), union.OfResultError(parseObject(lexer(iter.OfStringAsRunes(`{"foo":1`)))))
	// Case 12 - need space after key value so that error occurs after successfully returning number
	assert.Equal(t, union.OfError[json.Value](anErr), union.OfResultError(parseObject(lexer(iter.SetError(iter.OfStringAsRunes(`{"foo":1 `), anErr)))))
	// Case 13
	assert.Equal(t, union.OfError[json.Value](fmt.Errorf("The JSON key/value pair \"foo\" must be followed by a colon or closing brace")), union.OfResultError(parseObject(lexer(iter.OfStringAsRunes(`{"foo":1{`)))))
}

func TestParseArray_(t *testing.T) {
	assert.Equal(t, union.OfResult([]json.Value{json.FromString("bar")}), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`["bar"]`))))))
	assert.Equal(t, union.OfResult([]json.Value{json.FromString("foo"), json.FromString("bar")}), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`["foo", "bar"]`))))))

	anErr := fmt.Errorf("An err")

	// Case 1
	assert.Equal(t, union.OfError[[]json.Value](errArrayRequiresValueOrBracket), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`[`))))))
	// Case 2
	assert.Equal(t, union.OfError[[]json.Value](anErr), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.SetError(iter.OfStringAsRunes(`[`), anErr))))))
	// Case 3
	assert.Equal(t, union.OfResult([]json.Value{}), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`[]`))))))
	// Case 4
	assert.Equal(t, union.OfError[[]json.Value](anErr), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.SetError(iter.OfStringAsRunes(`[[`), anErr))))))
	// Case 5
	assert.Equal(t, union.OfError[[]json.Value](errArrayRequiresCommaOrBracket), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`[1`))))))
	// Case 6 - Need a space after value so that error occurs after successfully returning number
	assert.Equal(t, union.OfError[[]json.Value](anErr), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.SetError(iter.OfStringAsRunes(`[1 `), anErr))))))
	// Case 7
	assert.Equal(t, union.OfError[[]json.Value](errArrayRequiresValue), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`[1,`))))))
	// Case 8
	assert.Equal(t, union.OfError[[]json.Value](anErr), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.SetError(iter.OfStringAsRunes(`[1,`), anErr))))))
	// Case 9
	assert.Equal(t, union.OfError[[]json.Value](anErr), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.SetError(iter.OfStringAsRunes(`[1,[`), anErr))))))
	// Case 10 - ordinary success case
	assert.Equal(t, union.OfResult([]json.Value{json.FromNumberType(1)}), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`[1]`))))))
	// Case 11
	assert.Equal(t, union.OfError[[]json.Value](errArrayRequiresCommaOrBracket), iter.Maybe(stream.ReduceToSlice(parseArray(lexer(iter.OfStringAsRunes(`[1}`))))))
}

func TestIterate_(t *testing.T) {
	assert.Equal(t, union.OfResult([]json.Value{json.FromMap(map[string]any{"foo": "bar"})}), iter.Maybe(stream.ReduceToSlice(Iterate(strings.NewReader(`{"foo": "bar"}`)))))
	assert.Equal(t, union.OfResult([]json.Value{json.FromMap(map[string]any{"foo": "bar"})}), iter.Maybe(stream.ReduceToSlice(Iterate(strings.NewReader(`[{"foo": "bar"}]`)))))

	vals, err := stream.ReduceToSlice(Iterate(strings.NewReader(`[{"foo": "bar", "baz": 1}, ["fooey", 2], true, null]`))).Next()
	assert.Nil(t, err)
	assert.Equal(t, 4, len(vals))
	assert.Equal(t, json.FromMap(map[string]any{"foo": "bar", "baz": 1}), vals[0])
	assert.Equal(t, json.FromSlice([]any{"fooey", 2}), vals[1])
	assert.Equal(t, json.TrueValue, vals[2])
	assert.Equal(t, json.NullValue, vals[3])

	anErr := fmt.Errorf("An err")

	// Case 1
	assert.Equal(t, union.OfError[json.Value](errEmptyDocument), iter.Maybe(Iterate(strings.NewReader(``))))
	// Case 2
	assert.Equal(t, union.OfError[json.Value](anErr), iter.Maybe(Iterate(util.NewErrorReader([]byte(``), anErr))))
	// Case 3
	assert.Equal(t, union.OfError[json.Value](errObjectRequiresKeyOrBrace), iter.Maybe(Iterate(strings.NewReader(`{`))))
	// Case 4
	assert.Equal(t, union.OfError[json.Value](errObjectOrArrayRequired), iter.Maybe(Iterate(strings.NewReader(`:`))))
}

func TestParse_(t *testing.T) {
	assert.Equal(t, union.OfResult(json.FromMap(map[string]any{"foo": "bar"})), union.OfResultError(Parse(strings.NewReader(`{"foo": "bar"}`))))
	assert.Equal(t, union.OfResult(json.FromSlice([]any{map[string]any{"foo": "bar"}})), union.OfResultError(Parse(strings.NewReader(`[{"foo": "bar"}]`))))

	anErr := fmt.Errorf("An err")

	// Case 1
	assert.Equal(t, union.OfError[json.Value](errEmptyDocument), union.OfResultError(Parse(strings.NewReader(``))))
	// Case 2
	assert.Equal(t, union.OfError[json.Value](anErr), union.OfResultError(Parse(util.NewErrorReader([]byte(``), anErr))))
	// Case 3
	assert.Equal(t, union.OfError[json.Value](anErr), union.OfResultError(Parse(util.NewErrorReader([]byte(`{`), anErr))))
	// Case 4
	assert.Equal(t, union.OfError[json.Value](errObjectOrArrayRequired), union.OfResultError(Parse(strings.NewReader(`:`))))
}
