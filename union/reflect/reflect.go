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

// GetMaybeValue gets the value of a Maybe
// returns Invalid Value if the Maybe is empty
// returns Valid Value   if the Maybe is present
// panics if val is not a Maybe
func GetMaybeValue(val goreflect.Value) goreflect.Value {
	if val.MethodByName("Present").Call(nil)[0].Bool() {
		return val.MethodByName("Get").Call(nil)[0]
	}

	return goreflect.Value{}
}

// SetMaybeValue copies the value of val into dst
// Dst must be a Maybe[T], and val must be a T, otherwise a panic will occur
func SetMaybeValue(dst, val goreflect.Value) {
	dst.MethodByName("Set").Call([]goreflect.Value{val})
}

// SetMaybeValueEmpty sets a Maybe to empty
// Dst must be a Maybe[T], otherwise a panic will occur
func SetMaybeValueEmpty(dst goreflect.Value) {
	dst.MethodByName("SetEmpty").Call(nil)
}
