package reflect

// SPDX-License-Identifier: Apache-2.0

import (
	goreflect "reflect"
	"strings"

	"github.com/bantling/micro/reflect"
)

var (
	unionPkgPath = "github.com/bantling/micro/union"
)

// GetMaybeType gets the generic type of the value wrapped in a union.Maybe (which may have pointers to it).
// If the type is not a union.Maybe, then it returns a nil Type.
func GetMaybeType(typ goreflect.Type) goreflect.Type {
	dtyp := reflect.DerefType(typ)
	if (dtyp.PkgPath() == unionPkgPath) && strings.HasPrefix(dtyp.Name(), "Maybe") {
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
// - safeMaybeAccess(val) is true
// - drefing the reflect.Value to at most pointer and calling the Present method returns true
func MaybeValueIsPresent(val goreflect.Value) bool {
  return safeMaybeAccess(val) && reflect.DerefValueMaxOnePtr(val).MethodByName("Present").Call(nil)[0].Bool()
}

// GetMaybeValue gets the value of a Maybe
// returns Invalid Value if the Maybe is empty
// returns Valid Value   if the Maybe is present
//
// See MaybeValueIsPresent
func GetMaybeValue(val goreflect.Value) goreflect.Value {
	if MaybeValueIsPresent(val) {
		return reflect.DerefValueMaxOnePtr(val).MethodByName("Get").Call(nil)[0]
	}

	return goreflect.Value{}
}

// SetMaybeValue copies the value of val into dst
// Dst must be one or more pointers to a Maybe[T], and val must be a T, otherwise a panic will occur
func SetMaybeValue(dst, val goreflect.Value) error {
  if !safeMaybeAccess(dst) {
    return fmt.Errorf(errSetMaybeValueUnsafeMsg, dst.Type())
  }

  if (dst.Kind() != goreflect.Pointer) || (!dst.CanAddr()) {
    return fmt.Errorf(errSetMaybeValueUnaddressableMsg, dst.Type())
  }
  
	reflect.DerefValueMaxOnePtr(dst).MethodByName("Set").Call([]goreflect.Value{val})
}

// SetMaybeValueEmpty sets a Maybe to empty
// Dst must be a Maybe, otherwise a panic will occur
func SetMaybeValueEmpty(dst goreflect.Value) {
	reflect.DerefValueMaxOnePtr(dst).MethodByName("SetEmpty").Call(nil)
}
