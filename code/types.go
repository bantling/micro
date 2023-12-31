package code

// SPDX-License-Identifier: Apache-2.0

import (
  goreflect "reflect"
)

// Type is an enum of all basic types, by language
type Type uint

cont (
  // Boolean
  Bool Type = iota
  Uint8
  Uint16
  Uint32
  Uint64

  // Signed ints
  Int8
  Int16
  Int32
  Int64

  // String, UUID, JSON
  String
  Uuid//   = "uuid.UUID"  // provided by Google library github.com/google/uuid
  Json//   = "json.Value" // provided by this library in encoding/json

  // Date, DateTime, and Interval
  Date// = "time.Time" // provided by standard library, resolution is days since 2970
  DateTimeSeconds// = "time.Time" // provided by standard library, resolution is seconds since 1970
  DateTimeMilliseconds// = "time.Time" // provided by standard library, resolution is milliseconds since 1970
  IntervalDays// = "time.Duration" // provided by standard library, resolution is  days
  IntervalSeconds// = "time.Duration" // provided by standard library, resolution is seconds
  IntervalMilliseconds//  = "time.Duration" // provided by standard library, resolution is milliseconds

  // Enum
  Enum

  // Object, Map, Set, Array, Maybe
  Object,
  Map
  Set
  List
  Maybe
)

// VarDef is a variable type definition
type VarDef struct {
  Type Type  // The type of variable
  KeyType Type // The key type, if the variable is a map
  ValueType Var // The value type, if the variable is an object, map, set, list, or maybe
  Names []string // The names of enum constants - zero-based indexes are the enum values, strings are the names
  Maybe bool // True if the var is a Maybe, which can wrap around anything except map, set, list, and maybe
}

// Var is a variable instance
type Var struct {
  VarDef // The definition of the variable type
  Val string // The value of the variable
}

// Func is a function definition
type Func struct {
  Params []VarDef // Parameters of function
  LocalVars map[string]Var
  Results []VarDef // Results of function
}

// Object is an object, which can have fields and functions that operate on them
type Object struct {
  Fields map[string]Var
  Funcs map[string]Func
}

// Src is a source file
type Src struct {
  GlobalConstants map[string]Var
  GlobalVars map[string]Var
  Objects map[string]Object
  Funcs map[string]Func
}

// Dir is a directory of source files
type Dir struct {
  Sources []Src
}

// Package is a package of code, usually not a complete program (usually some other files will be hand written)
type Program struct {
  Dirs map[string]Dir
}
