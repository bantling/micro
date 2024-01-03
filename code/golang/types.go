package golang

// SPDX-License-Identifier: Apache-2.0

import (
  "github.com/bantling/micro/code"
)

// Map general primitive types to Go types
var (
  typeMap = map[code.Type]string {
    // Boolean
    code.Bool: "bool",

    // Unsigned ints
    code.Uint8: "uint8",
    code.Uint16: "uint16",
    code.Uint32: "uint32",
    code.Uint64: "uint64",

    // Signed ints
    code.Int8: "int8",
    code.Int16: "int16",
    code.Int32: "int32",
    code.Int64: "int64",

    // String, UUID, JSON
    code.String: "string",
    code.Uuid: "uuid.UUID",  // provided by Google library github.com/google/uuid
    code.Json: "json.Value", // provided by this library in encoding/json

    // Date, DateTime, and Interval
    code.Date: "time.Time", // provided by standard library, resolution is days since 2970
    code.DateTimeSeconds: "time.Time", // provided by standard library, resolution is seconds since 1970
    code.DateTimeMilliseconds: "time.Time", // provided by standard library, resolution is milliseconds since 1970
    code.IntervalDays: "time.Duration", // provided by standard library, resolution is days
    code.IntervalSeconds: "time.Duration", // provided by standard library, resolution is seconds
    code.IntervalMilliseconds: "time.Duration", // provided by standard library, resolution is milliseconds
  }
)
