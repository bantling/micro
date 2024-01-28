package code

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/union"
)

const (
	// Boolean
	Bool   = "bool"
	Uint8  = "uint8"
	Uint16 = "uint16"
	Uint32 = "uint32"
	Uint64 = "uint64"

	// Signed ints
	Int8  = "int8"
	Int16 = "int16"
	Int32 = "int32"
	Int64 = "int64"

	// String, UUID, JSON
	String = "string"
	Uuid   = "uuid"
	Json   = "json"

	// Date, DateTime, and Interval
	Date           = "date"           // days since 2970
	DateTimeSecs   = "dateTimeSecs"   // seconds since 1970
	DateTimeMillis = "dateTimeMillis" // milliseconds since 1970
	DurationDays   = "durationDays"   // days elapsed
	DurationSecs   = "durationSecs"   // seconds elapsed
	DurationMillis = "durationMillis" // elapsed milliseconds

	// Aray, Enum, List, Map, Maybe, Object, Set
  Array  = "array"
  Enum = "enum"
  List   = "list"
  Map    = "map"
  Maybe  = "maybe"
	Object = "object"
	Set    = "set"

  Main = "main" // name of main function
)

type TypeDef struct {
  Type      string // The type
  ArrayBounds    []int   // The bounds of an array, cannot be empty
  Name string // The name of an enum or object
  Names []string // The names of enum constants or object generics - zero-based indexes are the enum values, strings are the names
  ListDimension int // The dimension of the list (eg, list of string, list of list of string, etc), cannot be 0
  KeyType   TypeDef // The map key type
  ValueType TypeDef // The array element type, list element type, map value type, maybe type, or set type
}

// VarDef is a variable type definition
type VarDef struct {
	Type      TypeDef  // The type of variable
  Name      string   // The name of the variable
}

// FuncDef is a function definition
type FuncDef struct {
	Params      []VarDef          // Parameters of function
	LocalConsts map[string]VarDef // Local constants
	LocalVars   map[string]VarDef // Local vars
	Results     []VarDef          // Results of function
}

// ObjectDef is an object, which can have fields and functions that operate on them
type ObjectDef struct {
	Fields map[string]VarDef
	Funcs  map[string]FuncDef
}

// Src is a source file
type Src struct {
	GlobalConsts map[string]VarDef
	GlobalVars   map[string]VarDef
	Objects      map[string]ObjectDef
	Funcs        map[string]FuncDef
}

// Dir is a directory of source files
type Dir struct {
	Sources []Src

	// Init is an optional initialization function for the directory.
	// The set of all Init functions execute in some arbitrary order at runtime.
	// Depending on the target language, they may all execute before main starts, or they may execute some time later, such
	// as when files that need them are loaded.
	Init union.Maybe[FuncDef]
}

// Package is a package of code, not necessarily a complete program
type Package struct {
	Dirs map[string]Dir
}
