package json

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/tuple"
	"github.com/bantling/micro/union"
	"github.com/stretchr/testify/assert"
)

func TestParsePathInt_(t *testing.T) {
	lookups, err := parsePath(".addresses")
	assert.Equal(
		t,
		[]tuple.Two[string, union.Two[string, int]]{
			tuple.Of2(".addresses", union.Of2T[string, int]("addresses")),
		},
		lookups,
	)
	assert.Nil(t, err)

	lookups, err = parsePath("[3]")
	assert.Equal(
		t,
		[]tuple.Two[string, union.Two[string, int]]{
			tuple.Of2("[3]", union.Of2U[string, int](3)),
		},
		lookups,
	)
	assert.Nil(t, err)

	lookups, err = parsePath(".addresses[34].city")
	assert.Equal(
		t,
		[]tuple.Two[string, union.Two[string, int]]{
			tuple.Of2(".addresses", union.Of2T[string, int]("addresses")),
			tuple.Of2(".addresses[34]", union.Of2U[string, int](34)),
			tuple.Of2(".addresses[34].city", union.Of2T[string, int]("city")),
		},
		lookups,
	)
	assert.Nil(t, err)

	// Extra chars error
	lookups, err = parsePath(".addresses]")
	assert.Nil(t, lookups)
	assert.Equal(t, fmt.Errorf("The path .addresses] is not a valid path, it must consist of a series of object keys and indexes, such as .addresses[3].city"), err)

	// Doesn't match error
	lookups, err = parsePath("")
	assert.Nil(t, lookups)
	assert.Equal(t, fmt.Errorf("The path  is not a valid path, it must consist of a series of object keys and indexes, such as .addresses[3].city"), err)

	// Index too large error
	lookups, err = parsePath("[12345678901234567890]")
	assert.Nil(t, lookups)
	assert.Equal(t, fmt.Errorf("The path [12345678901234567890] is not a valid path, it must consist of a series of object keys and indexes, such as .addresses[3].city"), err)
}

func TestLookupsToFunc_(t *testing.T) {
	var (
		obj = MustMapToValue(
			map[string]any{
				"addresses": []any{
					map[string]any{
						"line": "123 Sesame St",
						"city": "New York",
					},
					map[string]any{
						"line": "456 Who Cares Drive",
						"city": "Los Angeles",
					},
				},
			},
		)

		arr = MustSliceToValue(
			[]any{
				map[string]any{
					"line": "123 Sesame St",
					"city": "New York",
				},
				map[string]any{
					"line": "456 Who Cares Drive",
					"city": "Los Angeles",
				},
			},
		)
	)

	// Object
	fn := lookupsToFunc(funcs.MustValue(parsePath(".addresses")))
	assert.Equal(t, union.OfResult(obj.AsMap()["addresses"]), union.OfResultError(fn(obj)))

	fn = lookupsToFunc(funcs.MustValue(parsePath(".addresses[0]")))
	assert.Equal(t, union.OfResult(obj.AsMap()["addresses"].AsSlice()[0]), union.OfResultError(fn(obj)))

	fn = lookupsToFunc(funcs.MustValue(parsePath(".addresses[0].city")))
	assert.Equal(t, union.OfResult(obj.AsMap()["addresses"].AsSlice()[0].AsMap()["city"]), union.OfResultError(fn(obj)))

	// Array
	fn = lookupsToFunc(funcs.MustValue(parsePath("[1]")))
	assert.Equal(t, union.OfResult(arr.AsSlice()[1]), union.OfResultError(fn(arr)))

	fn = lookupsToFunc(funcs.MustValue(parsePath("[1].city")))
	assert.Equal(t, union.OfResult(arr.AsSlice()[1].AsMap()["city"]), union.OfResultError(fn(arr)))

	// Not an Object error
	fn = lookupsToFunc(funcs.MustValue(parsePath(".addresses")))
	assert.Equal(t, union.OfError[Value](fmt.Errorf("The path .addresses cannot be found, as Array is not the correct type, or does not contain the index addresses")), union.OfResultError(fn(arr)))

	// Not such Object key error
	fn = lookupsToFunc(funcs.MustValue(parsePath(".foo")))
	assert.Equal(t, union.OfError[Value](fmt.Errorf("The path .foo cannot be found, as Object is not the correct type, or does not contain the index foo")), union.OfResultError(fn(obj)))

	// Not an Array error
	fn = lookupsToFunc(funcs.MustValue(parsePath("[0]")))
	assert.Equal(t, union.OfError[Value](fmt.Errorf("The path [0] cannot be found, as Object is not the correct type, or does not contain the index 0")), union.OfResultError(fn(obj)))

	// Not such Array index error
	fn = lookupsToFunc(funcs.MustValue(parsePath("[3]")))
	assert.Equal(t, union.OfError[Value](fmt.Errorf("The path [3] cannot be found, as Array is not the correct type, or does not contain the index 3")), union.OfResultError(fn(arr)))
}

func TestParsePath_(t *testing.T) {
	var (
		obj = MustMapToValue(
			map[string]any{
				"addresses": []any{
					map[string]any{
						"line": "123 Sesame St",
						"city": "New York",
					},
					map[string]any{
						"line": "456 Who Cares Drive",
						"city": "Los Angeles",
					},
				},
			},
		)
	)

	fn, err := ParsePath(".addresses[1].city")
	assert.NotNil(t, fn)
	assert.Nil(t, err)
	assert.Equal(t, union.OfResult(obj.AsMap()["addresses"].AsSlice()[1].AsMap()["city"]), union.OfResultError(fn(obj)))

	// Internal path parsing error
	fn, err = ParsePath("addresses")
	assert.Nil(t, fn)
	assert.Equal(t, fmt.Errorf("The path addresses is not a valid path, it must consist of a series of object keys and indexes, such as .addresses[3].city"), err)
}
