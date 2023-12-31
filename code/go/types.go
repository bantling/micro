package go

// SPDX-License-Identifier: Apache-2.0

import (
  goreflect "reflect"
)

// Map general primitive types to Go types
var (
  typeMap = map[code.Type]string {
    // Boolean
    Bool: "bool",

    // Unsigned ints
    Uint8: "uint8",
    Uint16: "uint16",
    Uint32: "uint32",
    Uint64: "uint64",

    // Signed ints
    Int8: "int8",
    Int16: "int16",
    Int32: "int32",
    Int64: "int64",

    // String, UUID, JSON
    String: "string",
    Uuid: "uuid.UUID",  // provided by Google library github.com/google/uuid
    Json: "json.Value", // provided by this library in encoding/json

    // Date, DateTime, and Interval
    Date: "time.Time", // provided by standard library, resolution is days since 2970
    DateTimeSeconds: "time.Time", // provided by standard library, resolution is seconds since 1970
    DateTimeMilliseconds: "time.Time", // provided by standard library, resolution is milliseconds since 1970
    IntervalDays: "time.Duration", // provided by standard library, resolution is days
    IntervalSeconds: "time.Duration", // provided by standard library, resolution is seconds
    IntervalMilliseconds: "time.Duration", // provided by standard library, resolution is milliseconds
  }
)
