package parse

// SPDX-License-Identifier: Apache-2.0

import (
	"io"

	"github.com/bantling/micro/go/iter"
	"github.com/bantling/micro/go/json"
)

// Parse parses a JSON document (object or array) into json.Value objects.
// The result is an *iter.Iter[json.Value], where an object is an iter of one json.Value,
// while an array is an iter of zero or more json.Value.
//
// The caller has to invoke an iter method like NextValue that pulls a json.Value from the iter, in order to get the parser
// to parse the next json.Value.
func Parse(src io.Reader) *iter.Iter[json.Value] {
	return nil
}
