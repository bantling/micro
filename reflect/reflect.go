package reflect

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math/big"
	goreflect "reflect"
	"strings"

	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/union"
)

const (
	errTypeAssertMsg      = "%s%s is %s, not %s"
	errNoSuchFieldNameMsg = "type %s does not have a field named %s"
	errNotAStructMsg      = "Cannot recurse type %s, it is not a struct"
	errFieldHandlingMsg   = "Type %s should never be encountered when recursing a type"
)

var (
	// Map kinds to base types to convert to
	kindToType = map[goreflect.Kind]goreflect.Type{
		goreflect.Bool:       goreflect.TypeOf(false),
		goreflect.Int:        goreflect.TypeOf(int(0)),
		goreflect.Int8:       goreflect.TypeOf(int8(0)),
		goreflect.Int16:      goreflect.TypeOf(int16(0)),
		goreflect.Int32:      goreflect.TypeOf(int32(0)),
		goreflect.Int64:      goreflect.TypeOf(int64(0)),
		goreflect.Uint:       goreflect.TypeOf(uint(0)),
		goreflect.Uint8:      goreflect.TypeOf(uint8(0)),
		goreflect.Uint16:     goreflect.TypeOf(uint16(0)),
		goreflect.Uint32:     goreflect.TypeOf(uint32(0)),
		goreflect.Uint64:     goreflect.TypeOf(uint64(0)),
		goreflect.Float32:    goreflect.TypeOf(float32(0)),
		goreflect.Float64:    goreflect.TypeOf(float64(0)),
		goreflect.Complex64:  goreflect.TypeOf(complex64(0)),
		goreflect.Complex128: goreflect.TypeOf(complex128(0)),
		goreflect.String:     goreflect.TypeOf(""),
	}

	// Map big type pointers to true for testing if a pointer is a big type
	bigPtrTypes = map[goreflect.Type]bool{
		goreflect.TypeOf((*big.Int)(nil)):   true,
		goreflect.TypeOf((*big.Float)(nil)): true,
		goreflect.TypeOf((*big.Rat)(nil)):   true,
	}
)

// RecurseMode indicates the progress of recursing a struct
type RecurseMode uint

const (
	Start RecurseMode = iota // Start recursing into a struct
	Field                    // A field of current struct
	End                      // Occurs after last field of current struct
)

var (
	recurseModeString = map[RecurseMode]string{
		Start: "Start",
		Field: "Field",
		End:   "End",
	}
)

// String is Stringer interface for RecurseMode
func (rm RecurseMode) String() string {
	return recurseModeString[rm]
}

// KindElem describes the Kind and Elem methods common to both Value and Type objects
type KindElem[T any] interface {
	Kind() goreflect.Kind
	Elem() T
}

// DerefType returns the element type of zero or more pointers to a type.
func DerefType(typ goreflect.Type) goreflect.Type {
	for typ.Kind() == goreflect.Pointer {
		typ = typ.Elem()
	}

	return typ
}

// DerefTypeMaxOnePtr returns zero or one pointers to a type.
// If the type is more than one pointer, it is derefd to one pointer, otherwse it is returned as is.
func DerefTypeMaxOnePtr(typ goreflect.Type) goreflect.Type {
	res := typ
	for res.Kind() == goreflect.Pointer {
		d := res.Elem()
		if d.Kind() != goreflect.Pointer {
			break
		}

		res = d
	}

	return res
}

// DerefValue returns the element of zero or more pointers to a value.
// If any pointer is nil, an invalid Value is returned.
func DerefValue(val goreflect.Value) goreflect.Value {
	for val.Kind() == goreflect.Pointer {
		if val.IsNil() {
			var zv goreflect.Value
			return zv
		}

		val = val.Elem()
	}

	return val
}

// DerefValueMaxOnePtr returns zero or one pointers to a value.
// If the value is more than one pointer, it is derefd to one pointer, otherwise it is returned as is.
// If any pointer except the last one is nil, an invalid Value is returned.
//
// There are 3 cases of results:
// - a valid Value for a non-pointer
// - a valid Value for a nil pointer to a non-pointer
// - an invalid Value for a nil pointer to a pointer
//
// Examples:
// DerefValueMaxOnePtr(reflect.ValueOf(0)) is a valid Value
//
// var p *int
// DerefValueMaxOnePtr(p) is a valid Value
//
// var p *int
// var p2 = &p
// DerefValueMaxOnePtr(p2) is a valid Value, since the outer pointer p2 is non-nil (doesn't matter p is nil)
//
// var p2 **int
// DerefValueMaxOnePtr(p) is an invalid Value, since the outer pointer is nil
func DerefValueMaxOnePtr(val goreflect.Value) goreflect.Value {
	res := val

	for res.Kind() == goreflect.Pointer {
		if res.IsNil() && (res.Type().Elem().Kind() == goreflect.Pointer) {
			// Nil pointer to a pointer type is invalid
			var zv goreflect.Value
			return zv
		}

		d := res.Elem()
		if d.Kind() != goreflect.Pointer {
			break
		}

		res = d
	}

	return res
}

// FieldsByName collects the fields of a struct into a map.
// Returns the zero value if the type provided does not represent a struct, or a struct that does not have any fields.
// If a given struct field is a struct, then another call would have to made on that struct.
// If a given struct field is a *struct, then it is possible it is a recursive struct (eg Customer{child *Customer}).
func FieldsByName(typ goreflect.Type) map[string]goreflect.StructField {
	var fields map[string]goreflect.StructField

	if (typ.Kind() == goreflect.Struct) && (typ.NumField() > 0) {
		fields = map[string]goreflect.StructField{}

		for i := 0; i < typ.NumField(); i++ {
			fld := typ.Field(i)
			fields[fld.Name] = fld
		}
	}

	return fields
}

// SetPointerValue copies the value of val into dst after dereffing all pointers
// Dst must be at least one pointer that derefs to some type T, and val must be convertible to T, otherwise a panic will occur
func SetPointerValue(dst, val goreflect.Value) {
	deref := DerefValue(dst)
	deref.Set(val.Convert(deref.Type()))
}

// IsBigPtr returns true if the given type is a *big.Int, *big.Float, or *big.Rat, and false otherwise
func IsBigPtr(typ goreflect.Type) bool {
	_, isBig := bigPtrTypes[typ]
	return isBig
}

// IsNumeric returns true if the given type satisfies constraint.Numeric
func IsNumeric(typ goreflect.Type) bool {
	knd := typ.Kind()
	return ((knd >= goreflect.Int) && (knd <= goreflect.Float64) && (knd != goreflect.Uintptr)) || IsBigPtr(typ)
}

// IsNillable returns true if ke.Kind() is nillable, which means it is Chan, Func, Interface, Map, Pointer, or Slice.
//
// If ke is nil, it means ke is reflect.TypeOf(a nil value of some interface type).
// If ke.Kind() is Invalid, it means ke is reflect.ValueOf(a nil value of any type).
//
// If ke is nil or Invalid, the result is true.
func IsNillable[T KindElem[T]](ke T) bool {
	if any(ke) == nil {
		return true
	}

	knd := ke.Kind()
	return (knd == goreflect.Invalid) || ((knd >= goreflect.Chan) && (knd <= goreflect.Slice))
}

// IsNil returns true if the value is invalid or a nillable type whose value is nil
func IsNil(val goreflect.Value) bool {
	return (!val.IsValid()) || (IsNillable(val.Type()) && val.IsNil())
}

func IsPrimitive[T KindElem[T]](val T) bool {
	_, hasIt := kindToType[val.Kind()]
	return hasIt
}

// NumPointers returns the number of pointers a type represents
func NumPointers(val goreflect.Type) (res int) {
	for val.Kind() == goreflect.Pointer {
		val = val.Elem()
		res++
	}

	return
}

// ResolveValueType resolves a value to the real type of value it contains.
// The only case where the result is different from the argument is when the argument is typed as interface{}.
// For example, if the interface{} value is actually an int, then the result will be typed as int.
// This generally only happens in corner cases like iterating the elements of a slice typed as []interface{} - even though
// the elements may be strings, ints, etc, each element will be typed as interface{}.
func ResolveValueType(val goreflect.Value) goreflect.Value {
	// Check special case
	if val.IsValid() && (val.Kind() == goreflect.Interface) {
		// Return a new wrapper that is typed according to actual value type
		return goreflect.ValueOf(val.Interface())
	}

	return val
}

// TypeAssert asserts that the value given has the same type as the type given.
// If not, the error returned contains a message that is similiar to the one Go provides when the type assertion syntax fails.
// Unlike Go's type assertion syntax, this function can be called with any kind of value and any type.
// If desired, an optional message can be provided that is placed at the beginning of the error message, followed by a colon and space.
func TypeAssert(val goreflect.Value, typ goreflect.Type, msg ...string) error {
	if rt := ResolveValueType(val).Type(); rt != typ {
		prefix := funcs.SliceIndex(msg, 0)
		return fmt.Errorf(errTypeAssertMsg, funcs.Ternary(prefix == "", prefix, prefix+": "), val.Type(), rt, typ)
	}

	return nil
}

// MustTypeAssert is a must version of TypeAssert
func MustTypeAssert(val goreflect.Value, typ goreflect.Type, msg ...string) {
	funcs.Must(TypeAssert(val, typ, msg...))
}

// TypeToBaseType converts a reflect.Type that may be a primitive subtype (eg type foo uint8) to the underlying type (eg uint8).
// If the given type is not a primitive subtype, nil is returned.
func TypeToBaseType(typ goreflect.Type) goreflect.Type {
	// Check if typ is a primitive subtype
	k := typ.Kind()
	if pt, isa := kindToType[k]; isa && (k.String() != typ.String()) {
		// If so, then return the base type
		return pt
	}

	return nil
}

// ValueToBaseType converts a reflect.Value that may be a primitive subtype (eg type byte uint8) to the underlying type (eg uint8).
// If the value is a pointer to a primitive subtype, the value is converted to a pointer to the underlying type.
func ValueToBaseType(val goreflect.Value) goreflect.Value {
	// Check if val is a primitive subtype
	k := val.Kind()
	pt := kindToType[k]
	if (pt != nil) && (k.String() != val.Type().String()) {
		// If so, then convert the value to the base type
		return val.Convert(pt)
	} else if k == goreflect.Ptr {
		k = val.Elem().Kind()
		pt = kindToType[k]
		if (pt != nil) && (k.String() != val.Elem().Type().String()) {
			return val.Convert(goreflect.PtrTo(pt))
		}
	}

	return val
}

// ValueMaxOnePtrType returns the underlying type of zero or one pointers to a value.
// If the value given has multiple pointers, the value is not a valid parameter value, and the result is nil.
//
// Examples:
// ValueMaxOnePtrType(reflect.ValueOf(0)) == reflect.TypeOf(0)
//
// var p *int
// ValueMaxOnePtrType(reflect.ValueOf(p)) == reflect.TypeOf(0)
//
// var p2 **int
// ValueMaxOnePtrType(reflect.ValueOf(p2)) == nil
func ValueMaxOnePtrType(val goreflect.Value) goreflect.Type {
	if !val.IsValid() {
		return nil
	}

	typ := val.Type()
	if typ.Kind() == goreflect.Pointer {
		typ = typ.Elem()
		if typ.Kind() == goreflect.Pointer {
			// ** is not a valid parameter
			return nil
		}
	}

	return typ
}

// StructFieldByName is a more convenient version of reflect.Type.FieldByName - it only returns a reflect.StructField,.
// If the field does not exist, it panics.
// Mostly useful in unit testing
func GetFieldByName(typ goreflect.Type, name string) goreflect.StructField {
	if sf, hasIt := typ.FieldByName(name); hasIt {
		return sf
	}

	panic(fmt.Errorf(errNoSuchFieldNameMsg, typ.String(), name))
}

// FieldHandler is a function to handle a single field of a struct.
// The path is the location of the field within the struct.
// The path can never be empty, it will always be at least one element.
type FieldHandler func(
	mode RecurseMode,
	path []union.Two[string, int],
	fld goreflect.StructField,
	val goreflect.Value,
) error

// RescurseFields recurses the fields of a struct, calling the given FieldHandler for each field.
// If a field is a struct, *struct, or Maybe[struct or *struct] it will be recursed if non-nil, except for go built-in types.
// If a field is an array, slice or map, it is up to the caller to handle them.
// Nil fields and empty Maybe fields are ignored.
//
// Detecting built in types is done by checking if the package path of the dereferenced type has at least one dot in it.
// Any value type that is not a built-in type should be pulled from the network, and be of the following form:
// <domain name>/<account>/<project>(/pkg)*
// Where the domain name will have at least one dot in it; built in types never have dots, just pkg(/pkg)*
//
// Example:
//
//			type Address Struct {
//			  Line string
//			  City string
//			}
//
//			type Customer struct {
//			  Name string
//			  Address
//	     SecondaryAddress Maybe[Address]
//			  Updated *time.Time
//			  Codes []string
//			}
//
// The following paths would be provided for a Customer, with modes shown in the order shown below.
// If the SecondaryAddress is empty, the lines shown below for SecondaryAddress would not be present.
// Similarly, if the Updated field is nil, the line shown below for Updated would not be present.
// The lines shown below for Codes assumes it contains the two values "CodeA", "CodeB" in that order.
// If Codes is nil or empty, the lines shown below for Codes would not be present.
//
// Start, []
// Field, [Name]
// Start, [Address]
// Field, [Address, Line]
// Field, [Address, City]
// End,   [Address]
// Start, [SecondaryAddress]
// Field, [SecondaryAddress, Line]
// Field, [SecondaryAddress, City]
// End,   [SecondaryAddress]
// Field, [Updated]
// Start, [Codes]
// Field, [Codes,0]
// Field, [Codes,1]
// End,   [Codes]
// End,   []
//
// Notes:
// - The provided context is passed to every invocation of the handler, to allow persistent user data across all calls
// - When mode is Start or End, the StructTag and Value handler parameters are zero values
// - If Customer.Address was declared as a pointer, and/or codes was an array, the same paths would still be returned.
// - Customer.Codes is not recursed, as no arrays, slices, or maps are recursed
// - Uintptr, Chan, and UnsafePointer fields are ignored
// - If the caller wants to modify the struct, then the caller must pass a pointer to the struct
//
// If the FieldHandler returns an error, then recursion stops, and that error is returned.
// Otherwise, all applicable fields are recursed, and nil is returned.
func RecurseFields(strukt goreflect.Value, handler FieldHandler) (err error) {
	// Error if the derefd type is not a struct
	ds := DerefValue(strukt)
	if ds.Kind() != goreflect.Struct {
		return fmt.Errorf(errNotAStructMsg, ds.Type())
	}

	var (
		recurse func(goreflect.Value)
		noField goreflect.StructField
		noValue goreflect.Value
		path    = []union.Two[string, int]{}
	)

	// Signal start of top struct
	if err := handler(Start, path, noField, noValue); err != nil {
		panic(err)
	}

	recurse = func(val goreflect.Value) {
		// Iterate all fields, if any
		for i, nf := 0, val.NumField(); i < nf; i++ {
			tf, df := val.Type().Field(i), DerefValue(val.Field(i))
			if knd := df.Kind(); (knd != goreflect.Uintptr) &&
				(knd != goreflect.Chan) &&
				(knd != goreflect.UnsafePointer) {
				// Add field name to end of path
				path = append(path, union.Of2T[string, int](tf.Name))

				switch {
				case (df.Kind() == goreflect.Struct) &&
					(
					// Struct is not from a go builtin package (has a dot in the package name, eg github.com)
					(strings.IndexRune(df.Type().PkgPath(), '.') >= 0) &&
						// Struct Field does not have tag of recurse:"-"
						tf.Tag.Get("recurse") != "-"):

					// Signal start of recursing sub struct
					if err := handler(Start, append([]union.Two[string, int]{}, path...), tf, df); err != nil {
						// Unwind recursion on first handler error
						panic(err)
					}

					// Recurse sub struct
					recurse(df)

					// Signal end of recursing sub struct
					if err := handler(End, append([]union.Two[string, int]{}, path...), tf, df); err != nil {
						// Unwind recursion on first handler error
						panic(err)
					}

				case ((df.Kind() == goreflect.Array) || (df.Kind() == goreflect.Slice)) && (tf.Tag.Get("recurse") != "-"):
					// Signal start of recursing sub slice
					if err := handler(Start, append([]union.Two[string, int]{}, path...), tf, df); err != nil {
						// Unwind recursion on first handler error
						panic(err)
					}

					// Recurse each array/slice index
					for i := 0; i < df.Len(); i++ {
						recurse(df.Index(i))
					}

					// Signal end of recursing indexes
					if err := handler(End, append([]union.Two[string, int]{}, path...), tf, df); err != nil {
						// Unwind recursion on first handler error
						panic(err)
					}

				default:
					// Signal field of current struct
					if err := handler(Field, append([]union.Two[string, int]{}, path...), tf, df); err != nil {
						// Unwind recursion on first handler error
						panic(err)
					}
				}

				// Remove field name from end of path
				path = path[:len(path)-1]
			}
		}
	}

	funcs.TryTo(
		func() { recurse(ds) },
		func(e any) {
			err = e.(error)
		},
	)

	// Signal end of top struct
	if err := handler(End, []union.Two[string, int]{}, noField, noValue); err != nil {
		panic(err)
	}

	return
}

// FieldsToMap converts the fields of a struct to a recursive map[string]any.
// It is an error if strukt is not a struct instance.
// The map field names are converted from snake case to camel case.
func FieldsToMap(strukt goreflect.Value) (map[string]any, error) {
	var (
		resultMap  = map[string]any{} // result map to return
		err        error
		currentMap = resultMap          // current map to work with during recursion
		mapList    = []map[string]any{} // list of maps to add/remove during recursion
	)

	err = RecurseFields(
		strukt,
		func(
			mode RecurseMode,
			path []union.Two[string, int],
			fld goreflect.StructField,
			val goreflect.Value,
		) error {
			switch mode {
			case Start:
				// Create a new submap to recurse into and append it to the list, so we can return to it later
				currentMap = map[string]any{}
				mapList = append(mapList, currentMap)

			case Field:
				var (
					lastPathPart = funcs.SliceIndex(path, -1)
					derefdVal    = DerefValue(val)
					knd          = derefdVal.Kind()
				)

				switch {
				case knd == goreflect.Invalid:
					// One or more pointers, where some pointer is nil: put a correctly typed nil value in the map
					currentMap[lastPathPart.T()] = goreflect.Zero(fld.Type).Interface()

				case ((knd >= goreflect.Bool) && (knd <= goreflect.Complex128)) || (knd == goreflect.String):
					// It's a scalar type, store the value in the map
					currentMap[lastPathPart.T()] = val.Interface()

				default:
					// Should never reach this code
					panic(fmt.Errorf(errFieldHandlingMsg, fld.Type))
				}

			default: // Must be End
				// Remove map at end of list, and set current map to new end of list (if it is not empty)
				mapList = mapList[0 : len(mapList)-1]
				currentMap = funcs.SliceIndex(mapList, -1)
			}

			return nil
		},
	)

	return resultMap, err
}
