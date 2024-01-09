package code

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/union"
)

// Type is an enum of all basic types
type Type uint

const (
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
	Uuid //   = "uuid.UUID"  // provided by Google library github.com/google/uuid
	Json //   = "json.Value" // provided by this library in encoding/json

	// Date, DateTime, and Interval
	Date                 // = "time.Time" // provided by standard library, resolution is days since 2970
	DateTimeSeconds      // = "time.Time" // provided by standard library, resolution is seconds since 1970
	DateTimeMilliseconds // = "time.Time" // provided by standard library, resolution is milliseconds since 1970
	IntervalDays         // = "time.Duration" // provided by standard library, resolution is  days
	IntervalSeconds      // = "time.Duration" // provided by standard library, resolution is seconds
	IntervalMilliseconds //  = "time.Duration" // provided by standard library, resolution is milliseconds

	// Enum
	Enum

	// Object, Map, Set, List, Maybe
	Object
	Map
	Set
	List
	Maybe
)

// VarDef is a variable type definition
type VarDef struct {
	Type      Type     // The type of variable
	KeyType   Type     // The key type, if the variable is a map
	ValueType Type     // The value type, if the variable is an object, map, set, list, or maybe
	Names     []string // The names of enum constants - zero-based indexes are the enum values, strings are the names
	Maybe     bool     // True if the var is a Maybe, which can wrap around anything except map, set, list, and maybe
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

// Package is a package of code, usually not a complete program (usually some other files will be hand written)
type Package struct {
	Dirs map[string]Dir
}
