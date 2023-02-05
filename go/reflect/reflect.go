package reflect

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/go/funcs"
	"math/big"
	goreflect "reflect"
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

// KindElem describes the Kind and Elem methods common to both Value and Type objects
type KindElem[T any] interface {
	Kind() goreflect.Kind
	Elem() T
}

func IsPrimitive[T KindElem[T]](val T) bool {
	_, hasIt := kindToType[val.Kind()]
	return hasIt
}

// ToBaseType converts a reflect.Value that may be a primitive subtype (eg type byte uint8) to the underlying type (eg uint8).
func ToBaseType(val *goreflect.Value) {
	// Check if val is a primitive subtype
	k := val.Kind()
	pt := kindToType[k]
	if (pt != nil) && (k.String() != val.Type().String()) {
		// If so, then convert the value to the base type so we can pass it to the conversion function
		*val = val.Convert(pt)
		// Check if val is a pointer to a primitive subtype
	} else if k == goreflect.Ptr {
		k = val.Elem().Kind()
		pt = kindToType[k]
		if (pt != nil) && (k.String() != val.Elem().Type().String()) {
			*val = val.Convert(goreflect.PtrTo(pt))
		}
	}
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

// ResolveValueType resolves a value to the real type of value it contains.
// The only case where the result is different from the argument is when the argument is typed as interface{}.
// For example, if the interface{} value is actually an int, then the result will be typed as int.
// This generally only happens in corner cases like iterating the elements of a slice typed as []interface{} - even though
// the elements may be strings, ints, etc, each element will be typed as []interface{}.
func ResolveValueType(val goreflect.Value) goreflect.Value {
	// Check special case
	if val.IsValid() && (val.Kind() == goreflect.Interface) {
		// Return a new wrapper that is typed according to actual value type
		return goreflect.ValueOf(val.Interface())
	}

	return val
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
// ValueMaxOnePtrIsType(reflect.ValueOf(p2)) == nil
func ValueMaxOnePtrType(val goreflect.Value) goreflect.Type {
	adjustedVal := funcs.Ternary(val.Kind() == goreflect.Pointer, val.Elem(), val)
	return funcs.Ternary(adjustedVal.Kind() == goreflect.Pointer, nil, adjustedVal.Type())
}

// FieldsByName collects the fields of a struct into a map.
// Returns the zero value if the type provided does not represent a struct, or a struct that does not have any fields.
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

// IsBigPtr returns true if the given type is a *big.Int, *big.Float, or *big.Rat, and false otherwise
func IsBigPtr(typ goreflect.Type) bool {
	_, isBig := bigPtrTypes[typ]
	return isBig
}
