package code

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/union"
)

// Type is the type of any kind of value
type Type int

// Scalar is a Type that has only one value
type ScalarType Type

const (
	Bool ScalarType = iota

	// Unsigned ints
	Uint8
	Uint16
	Uint32
	Uint64

	// Signed ints
	Int8
	Int16
	Int32
	Int64

	// Floating point
	Float
	Double

	// String, UUID, JSON
	String
	UUID

	// Date, DateTime, and Duration
	Date           // days since 1970
	DateTimeSecs   // seconds since 1970
	DateTimeMillis // milliseconds since 1970
	DurationDays   // days elapsed
	DurationSecs   // seconds elapsed
	DurationMillis // elapsed milliseconds
)

// Aggregate is a type that has multiple values
type AggregateType Type

const (
	JSON AggregateType = iota
	Array
	Enum
	List
	Map
	Maybe
	Object
	Set
)

// AccessLevel indicates the level of access to a source member
// Depending on the target language, not all access levels may be supported
type AccessLevel int

const (
	Private    AccessLevel = iota // Private access
	Package                       // Same package
	PackageSub                    // Same package or sub objects
	Public                        // All code
)

// TypeDef is a type definition
type TypeDef struct {
	typ           Type                     // Type
	arrayBounds   []int                    // Bounds of an array, cannot be empty
	name          string                   // Name of an enum or object
	names         []string                 // Names of enum constants or object generics - zero-based indexes are the enum values, strings are the names
	listDimension uint                     // Dimension of the list (eg, list of string, list of list of string, etc), cannot be 0
	keyType       *TypeDef                 // Map key type
	valueType     *TypeDef                 // Array element type, list element type, map value type, maybe type, or set type
	access        union.Maybe[AccessLevel] // Level of access, if applicable
}

// OfScalarType constructs a scalar TypeDef
func OfScalarType(
	typ ScalarType,
	access ...AccessLevel,
) TypeDef {
	return TypeDef{
		typ:    Type(typ),
		access: union.First(access...),
	}
}

// ConstDef is a constant definition
type ConstDef struct {
	TypeName string      // TypeDef.Name
	Name     string      // Constant name
	Value    string      // Value
	Access   AccessLevel // Level of access
}

// VarDef is a variable definition
type VarDef struct {
	TypeName string      // TypeDef.Name
	Name     string      // Var name
	Access   AccessLevel // Level of access
}

// LitDef is a literal definition
type LitDef struct {
	TypeName string // TypeDef.Name
	Value    string // Literal value
}

// ValKind is the kind of a value
type ValKind int

const (
	ConstVal = iota // A constant value
	VarVal          // A variable value
	LitVal          // A literal value
)

// Val represents a value of some type
// It is a constant, variable or literal
type Val struct {
	Kind  ValKind // The kind of value
	Name  string  // The ConstDef.Name or VarDef.Name, if applicable
	Value string  // The literal value, if applicable
}

// FuncDef is a function definition
type FuncDef struct {
	Params      map[string]VarDef   // Parameters of function
	LocalConsts map[string]ConstDef // Local constants
	LocalVars   map[string]VarDef   // Local vars
	Results     []TypeDef           // Results of function
	//Code                            // Code of function
	Access AccessLevel // The level of access
}

// ObjectDef is an object, which can have fields and functions that operate on them
type ObjectDef struct {
	TypeDef TypeDef            // Name of the object type, and any generics
	Fields  map[string]TypeDef // Object fields
	Funcs   map[string]FuncDef // Object methods
}

// Src is a source file
type Src struct {
	GlobalConsts map[string]TypeDef   // Global constants
	GlobalVars   map[string]TypeDef   // Global vars
	Objects      map[string]ObjectDef // Objects
	Funcs        map[string]FuncDef   // Top level functions that are not methods
	Main         FuncDef              // Main function
}

// Pkg is a directory of source files
type Pkg struct {
	Path    string // Relative path to dir containing sources
	Sources []Src  // Source files

	// Init is an optional initialization function for the package.
	// The set of all package init functions execute in some arbitrary order at runtime.
	// Depending on the target language, they may all execute before main starts, or they may execute some time later.
	// such as when files that need them are loaded.
	Init union.Maybe[FuncDef]
}
