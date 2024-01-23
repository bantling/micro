package golang

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/code"
)

// Map general primitive types to Go types
var (
	typeMap = map[string]string{
		// Boolean
		code.Bool: "bool",

		// Unsigned ints
		code.Uint8:  "uint8",
		code.Uint16: "uint16",
		code.Uint32: "uint32",
		code.Uint64: "uint64",

		// Signed ints
		code.Int8:  "int8",
		code.Int16: "int16",
		code.Int32: "int32",
		code.Int64: "int64",

		// String, UUID, JSON
		code.String: "string",
		code.Uuid:   "uuid.UUID",  // provided by Google library github.com/google/uuid
		code.Json:   "json.Value", // provided by this library in encoding/json

		// Date, DateTime, and Interval
		code.Date:                 "time.Time",
		code.DateTimeSecs:      "time.Time",
		code.DateTimeMillis: "time.Time",
		code.DurationDays:         "time.Duration",
		code.DurationSecs:      "time.Duration",
		code.DurationMillis: "time.Duration",
	}
)

// TypeDefString produces a type declaration
func TypeDefString(def code.TypeDef, top bool) string {
  switch def.Type {
  case code.Object:
    return funcs.Ternary(top, "struct ", "") + TypeDeclaration(def.ObjectType, false) + " {"
  case code.Map:
    return "map[" + TypeDefString(def.KeyType, false) + "]" + TypeDefString(def.ValueType, false)
  case code.Set:
    return "map[" + TypeDefString(def.ValueType, false) + "]bool"
  case code.Array:
    return iter.Maybe(
      stream.Reduce(
        func(a, b string) {return a + b},
      )(
        stream.Map(func(bound int) string { return "[" + strconv.Itoa(bound) + "]" }) (iter.Of(def.Bounds...))
      )
    ).Get() + TypeDefString(def.ValueType, false)
  case code.List:
    return "[]" + TypeDefString(def.ValueType, false)
  }
}
