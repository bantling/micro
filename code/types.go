package code

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/union"
)

// Type is the type of any kind of value
type Type uint

// Scalar is a Type that has only one value
type ScalarType Type

const (
	Bool ScalarType = iota

	// Unsigned ints
	Uint
	Uint8
	Uint16
	Uint32
	Uint64

	// Signed ints
	Int
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
	afterScalar    // internal constant for one past last scalar
)

// Aggregate is a type that has multiple values
type AggregateType Type

const (
	JSON AggregateType = iota + AggregateType(afterScalar)
	Array
	Enum
	List
	Map
	Maybe
	Object
	Set
)

// Types is a constraint on Type and all known subtypes of it
type Types interface {
	Type | ScalarType | AggregateType
}

// IsScalar is true if the Type is a Scalar Type
func IsScalar(typ Type) bool {
	return uint(typ) < uint(afterScalar)
}

// IsAggregate is true if the Type is an Aggregate Type
func IsAggregate(typ Type) bool {
	return uint(typ) >= uint(afterScalar)
}

// AccessLevel indicates the level of access to a source member
// Depending on the target language, not all access levels may be supported
// The default level is Private
type AccessLevel int

const (
	Private    AccessLevel = iota // Private access
	Package                       // Same package
	PackageSub                    // Same package or sub objects
	Public                        // All code
)

// TypeDef is a type definition
type TypeDef struct {
	Access      AccessLevel // Level of access, if applicable. Empty means default level.
	Typ         Type        // Type
	ArrayBounds []uint      // Bounds of an array. A dimension can be -1 for unspecified dimension, slice can
	// be zero length for one unspecified dimension.
	// For a list, one element >= 1 that indicates list, or list of list, etc.
	Name      string                // Name of an enum or object
	Names     []string              // Names of enum constants or object generics - zero-based indexes are the enum values, strings are the names
	KeyType   union.Maybe[*TypeDef] // Map key type
	ValueType union.Maybe[*TypeDef] // Array element type, Enum base type, list element type, map value type, maybe type, object base type, or set element type
}

// OfScalarType constructs a scalar TypeDef
func OfScalarType(
	typ ScalarType,
	access ...AccessLevel,
) *TypeDef {
	return &TypeDef{
		Access: funcs.SliceIndex(access, 0, Private),
		Typ:    Type(typ),
	}
}

// OfJSONType constructs a JSON TypeDef
func OfJSONType(access ...AccessLevel) *TypeDef {
	return &TypeDef{
		Access: funcs.SliceIndex(access, 0, Private),
		Typ:    Type(JSON),
	}
}

// OfArrayType constructs an array TypeDef
func OfArrayType(
	elementTyp *TypeDef,
	bounds []uint,
	access ...AccessLevel,
) *TypeDef {
	return &TypeDef{
		Access:      funcs.SliceIndex(access, 0, Private),
		Typ:         Type(Array),
		ArrayBounds: bounds,
		ValueType:   union.Present(elementTyp),
	}
}

// OfEnumType constructs an enum TypeDef
func OfEnumType(
	name string,
	baseType *TypeDef,
	constants []string,
	access ...AccessLevel,
) *TypeDef {
	return &TypeDef{
		Access:    funcs.SliceIndex(access, 0, Private),
		Typ:       Type(Enum),
		Name:      name,
		Names:     funcs.MustNonEmptySlice(constants),
		ValueType: union.Present(baseType),
	}
}

// OfListType constructs a list TypeDef
func OfListType(
	elementType *TypeDef,
	access ...AccessLevel,
) *TypeDef {
	return &TypeDef{
		Access:    funcs.SliceIndex(access, 0, Private),
		Typ:       Type(List),
		ValueType: union.Present(elementType),
	}
}

// OfMapType constructs a map TypeDef
func OfMapType(
	keyType *TypeDef,
	valueType *TypeDef,
	access ...AccessLevel,
) *TypeDef {
	return &TypeDef{
		Access:    funcs.SliceIndex(access, 0, Private),
		Typ:       Type(Map),
		KeyType:   union.Present(keyType),
		ValueType: union.Present(valueType),
	}
}

// OfMaybeType constructs a Maybe TypeDef
func OfMaybeType(
	elementType *TypeDef,
	access ...AccessLevel,
) *TypeDef {
	return &TypeDef{
		Access:    funcs.SliceIndex(access, 0, Private),
		Typ:       Type(Maybe),
		ValueType: union.Present(elementType),
	}
}

// OfObjectType constructs an Object TypeDef
func OfObjectType(
	baseType *TypeDef,
	name string,
	generics []string,
	access ...AccessLevel,
) *TypeDef {
	return &TypeDef{
		Access:    funcs.SliceIndex(access, 0, Private),
		Typ:       Type(Object),
		Name:      name,
		Names:     generics,
		ValueType: union.Present(baseType),
	}
}

// OfSetType constructs a set TypeDef
func OfSetType(
	elementType *TypeDef,
	access ...AccessLevel,
) *TypeDef {
	return &TypeDef{
		Access:    funcs.SliceIndex(access, 0, Private),
		Typ:       Type(Set),
		ValueType: union.Present(elementType),
	}
}

// ValKind is the kind of a value
type ValKind int

const (
	LitVal   ValKind = iota // A literal value
	VarVal                  // A variable value
	VarConst                // A constant variable value
)

// Val represents a value of some type
// It is a (non-)constant variable or literal
type Val struct {
	Access  union.Maybe[AccessLevel] // Level of access
	Kind    ValKind                  // The kind of value
	Typ     *TypeDef                 // The TypeDef
	Value   string                   // The literal value or variable name
}

// OfLitVal constructs a literal value
func OfLitVal(
	typeDef *TypeDef,
	val string,
) Val {
	return Val{
		Kind : LitVal,
		Typ  : funcs.MustNonNilValue(typeDef),
		Value: val,
	}
}

// OfVarVal constructs a variable value, which may be constant
func OfVarVal(
	konst bool,
	typeDef *TypeDef,
	val string,
	access ...AccessLevel,
) Val {
	return Val{
		Access:  union.First(access),
		Kind:    funcs.Ternary(konst, VarConst, VarVal),
		Typ:    funcs.MustNonNilValue(typeDef),
		Value:   val,
	}
}

// FuncDef is a function definition
type FuncDef struct {
	Access  AccessLevel    // The level of access
	Params  map[string]Val // Parameters of function
	Locals  map[string]Val // Local constants and vars
	Results []TypeDef      // Results of function
	//Code                            // Code of function
}

// ObjectDef is an object, which can have fields and functions that operate on them
type ObjectDef struct {
	TypeDef TypeDef            // Name of the object type, and any generics
	Fields  map[string]TypeDef // Object fields
	Funcs   map[string]FuncDef // Object methods
}

// SrcDef is a source file
type SrcDef struct {
	Globals map[string]TypeDef   // Global constants and vars
	Objects map[string]ObjectDef // Objects
	Funcs   map[string]FuncDef   // Top level functions that are not methods
	Main    FuncDef              // Main function
}

// PkgDef is a directory of source files
type PkgDef struct {
	Path    string   // Relative path to dir containing sources
	Sources []SrcDef // Source files

	// Init is an optional initialization function for the package.
	// The set of all package init functions execute in some arbitrary order at runtime.
	// Depending on the target language, they may all execute before main starts, or they may execute some time later.
	// such as when files that need them are loaded.
	Init union.Maybe[FuncDef]
}
