package parse

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bantling/micro/go/funcs"
	"github.com/bantling/micro/go/iter"
	"github.com/bantling/micro/go/json"
	"github.com/stretchr/testify/assert"
)

func mkIter(str string) iter.Iter[token] {
	return lexer(iter.OfStringAsRunes(str))
}

func TestParseValue(t *testing.T) {
	assert.Equal(t, json.FromMap(map[string]any{"foo": "bar"}), parseValue(mkIter(`{"foo": "bar"}`)))
	assert.Equal(t, json.FromSlice([]any{"bar"}), parseValue(mkIter(`["bar"]`)))
	assert.Equal(t, json.FromString("bar"), parseValue(mkIter(`"bar"`)))
	assert.Equal(t, json.FromNumberString("1.25"), parseValue(mkIter("1.25")))
	assert.Equal(t, json.TrueValue, parseValue(mkIter("true")))
	assert.Equal(t, json.NullValue, parseValue(mkIter("null")))

	// parseValue returns invalid Value for tokens that cannot be a value
	assert.Equal(t, json.Value{}, parseValue(mkIter("}")))
}

func TestParseObject(t *testing.T) {
	assert.Equal(t, json.FromMap(map[string]any{}), parseObject(mkIter(`{}`)))
	assert.Equal(t, json.FromMap(map[string]any{"foo": "bar"}), parseObject(mkIter(`{"foo": "bar"}`)))

	funcs.TryTo(
		func() {
			parseObject(mkIter(`{`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, errObjectRequiresKeyOrBrace, e) },
	)

	funcs.TryTo(
		func() {
			parseObject(mkIter(`{{`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, errObjectRequiresKeyOrBrace, e) },
	)

	funcs.TryTo(
		func() {
			parseObject(mkIter(`{"foo": "bar", "foo": "bar"`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errObjectDuplicateKeyMsg, "foo"), e) },
	)

	funcs.TryTo(
		func() {
			parseObject(mkIter(`{"foo"{`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errObjectKeyRequiresColonMsg, "foo"), e) },
	)

	funcs.TryTo(
		func() {
			parseObject(mkIter(`{"foo"{`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errObjectKeyRequiresColonMsg, "foo"), e) },
	)

	funcs.TryTo(
		func() {
			parseObject(mkIter(`{"foo":}`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errObjectKeyRequiresValueMsg, "foo"), e) },
	)

	funcs.TryTo(
		func() {
			parseObject(mkIter(`{"foo":1`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errObjectKeyValueRequiresCommaOrBraceMsg, "foo"), e) },
	)

	funcs.TryTo(
		func() {
			parseObject(mkIter(`{"foo":1{`))
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, fmt.Errorf(errObjectKeyValueRequiresCommaOrBraceMsg, "foo"), e) },
	)
}

func TestParseArray(t *testing.T) {
	assert.Equal(t, json.FromSlice([]any{}), json.FromSliceOfValue(iter.ReduceToSlice(parseArray(mkIter(`[]`))).Must()))
	assert.Equal(t, json.FromSlice([]any{"bar"}), json.FromSliceOfValue(iter.ReduceToSlice(parseArray(mkIter(`["bar"]`))).Must()))
	assert.Equal(t, json.FromSlice([]any{"foo", "bar"}), json.FromSliceOfValue(iter.ReduceToSlice(parseArray(mkIter(`["foo", "bar"]`))).Must()))

	funcs.TryTo(
		func() {
			parseArray(mkIter(`[`)).Must()
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, errArrayRequiresValueOrBracket, e) },
	)

	funcs.TryTo(
		func() {
			iter.ReduceToSlice(parseArray(mkIter(`[}`))).Must()
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, errArrayRequiresValueOrBracket, e) },
	)

	funcs.TryTo(
		func() {
			iter.ReduceToSlice(parseArray(mkIter(`["bar"`))).Must()
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, errArrayRequiresCommaOrBracket, e) },
	)

	funcs.TryTo(
		func() {
			iter.ReduceToSlice(parseArray(mkIter(`["bar",`))).Must()
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, errArrayRequiresValue, e) },
	)

	funcs.TryTo(
		func() {
			iter.ReduceToSlice(parseArray(mkIter(`["bar",}`))).Must()
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, errArrayRequiresValue, e) },
	)

	funcs.TryTo(
		func() {
			iter.ReduceToSlice(parseArray(mkIter(`["bar"}`))).Must()
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, errArrayRequiresCommaOrBracket, e) },
	)
}

func TestIterate(t *testing.T) {
	assert.Equal(
		t,
		json.FromMap(map[string]any{"foo": "bar"}),
		iter.ReduceToSlice(Iterate(strings.NewReader(`{"foo": "bar"}`))).Must()[0],
	)

	assert.Equal(
		t,
		json.FromMap(map[string]any{"foo": "bar"}),
		iter.ReduceToSlice(Iterate(strings.NewReader(`[{"foo": "bar"}]`))).Must()[0],
	)

	vals := iter.ReduceToSlice(Iterate(strings.NewReader(`[{"foo": "bar", "baz": 1}, ["fooey", 2], true, null]`))).Must()
	assert.Equal(t, 4, len(vals))
	assert.Equal(t, json.FromMap(map[string]any{"foo": "bar", "baz": 1}), vals[0])
	assert.Equal(t, json.FromSlice([]any{"fooey", 2}), vals[1])
	assert.Equal(t, json.TrueValue, vals[2])
	assert.Equal(t, json.NullValue, vals[3])

	funcs.TryTo(
		func() {
			iter.ReduceToSlice(Iterate(strings.NewReader(``))).Must()
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, errEmptyDocument, e) },
	)

	funcs.TryTo(
		func() {
			iter.ReduceToSlice(Iterate(strings.NewReader(`,`))).Must()
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, errObjectOrArrayRequired, e) },
	)
}

func TestParse(t *testing.T) {
	val, err := Parse(strings.NewReader(`{"foo": "bar"}`))
	assert.Nil(t, err)
	assert.Equal(t, json.FromMap(map[string]any{"foo": "bar"}), val)

	val, err = Parse(strings.NewReader(`[{"foo": "bar"}]`))
	assert.Nil(t, err)
	assert.Equal(t, json.FromSlice([]any{map[string]any{"foo": "bar"}}), val)

	val, err = Parse(strings.NewReader(``))
	assert.Equal(t, json.Value{}, val)
	assert.Equal(t, errEmptyDocument, err)

	val, err = Parse(strings.NewReader(`,`))
	assert.Equal(t, json.Value{}, val)
	assert.Equal(t, errObjectOrArrayRequired, err)
}
