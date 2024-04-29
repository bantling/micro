package reflect

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	goreflect "reflect"
	"strings"

	"github.com/bantling/micro/reflect"
)

var (
	unionPkgPath = "github.com/bantling/micro/union"
)

var (
	errGetMaybeValueEmptyMsg         = "Cannot get the Maybe value of an empty %s"
	errSetMaybeValueUnsafeMsg        = "Cannot set the Maybe value of type %s"
	errSetMaybeValueUnaddressableMsg = "Cannot set the Maybe value of type %s as it is not a pointer and not addressable"
)

// GetMaybeType gets the generic type of the value wrapped in a union.Maybe (which may have pointers to it).
// If the type is not a union.Maybe, then it returns a nil Type.
func GetMaybeType(typ goreflect.Type) goreflect.Type {
	dtyp := reflect.DerefType(typ)
	if (dtyp.PkgPath() == unionPkgPath) && strings.HasPrefix(dtyp.Name(), "Maybe[") {
		if get, hasIt := dtyp.MethodByName("Get"); hasIt {
			return get.Type.Out(0)
		}
	}

	return nil
}

// safeMaybeAccess returns true if the given reflect.Value can safely be accessed as a Maybe value, which is true if:
// - The reflect.Value is valid
// - GetMaybeType(val.Type()) returns a non-nil reflect.Type
// - The reflect.Value is a non-nil pointer or a value
func safeMaybeAccess(val goreflect.Value) bool {
	return val.IsValid() &&
		(GetMaybeType(val.Type()) != nil) &&
		((val.Kind() != goreflect.Pointer) || !val.IsNil())
}

// MaybeValueIsPresent indicates if the given reflect.Value wraps a present Maybe. which is true if:
// - The reflect.Value is valid
// - GetMaybeType(val.Type()) returns a non-nil reflect.Type
// - derefing the reflect.Value to at most one pointer and calling the Present method returns true
func MaybeValueIsPresent(val goreflect.Value) bool {
	return safeMaybeAccess(val) && reflect.DerefValueMaxOnePtr(val).MethodByName("Present").Call(nil)[0].Bool()
}

// GetMaybeValue gets the value of a Maybe, returning reflect.Value of present value.
// A invalid reflect.Value is returned if
// - The reflect.Value is invalid
// - The reflect.Value is not a Maybe
// - The reflect.Value is an empty Maybe
// When an error ocurs, (invalid reflect.Value, error) is returned
func GetMaybeValue(val goreflect.Value) goreflect.Value {
	if !MaybeValueIsPresent(val) {
		return goreflect.Value{}
	}

	return reflect.DerefValueMaxOnePtr(val).MethodByName("Get").Call(nil)[0]
}

// SetMaybeValue copies the value of val into dst.
// Dst must be zero or more pointers to an addressable Maybe[T], and val must be a T, otherwise an error will occur.
func SetMaybeValue(dst, val goreflect.Value) error {
	if !safeMaybeAccess(dst) {
		return fmt.Errorf(errSetMaybeValueUnsafeMsg, reflect.TypeOf(dst))
	}

	if (dst.Kind() != goreflect.Pointer) && (!dst.CanAddr()) {
		return fmt.Errorf(errSetMaybeValueUnaddressableMsg, reflect.TypeOf(dst))
	}

	reflect.DerefValue(dst).Addr().MethodByName("Set").Call([]goreflect.Value{val})
	return nil
}

// SetMaybeValueEmpty sets a Maybe to empty. If the Maybe is already Empty, it is effectively a non operation.
// Dst must be a Maybe, otherwise an error will occur
func SetMaybeValueEmpty(dst goreflect.Value) error {
	if !safeMaybeAccess(dst) {
		return fmt.Errorf(errSetMaybeValueUnsafeMsg, reflect.TypeOf(dst))
	}

	if (dst.Kind() != goreflect.Pointer) && (!dst.CanAddr()) {
		return fmt.Errorf(errSetMaybeValueUnaddressableMsg, reflect.TypeOf(dst))
	}

	reflect.DerefValue(dst).Addr().MethodByName("SetEmpty").Call(nil)
	return nil
}
