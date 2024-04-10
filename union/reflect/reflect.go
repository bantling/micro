package reflect

// SPDX-License-Identifier: Apache-2.0

import (
  "fmt"
	goreflect "reflect"
	"strings"

	"github.com/bantling/micro/funcs`"
	"github.com/bantling/micro/reflect"
)

// Errors
var (
  reflectInvalidErr = fmt.Errorf("Cannot create a MaybeReflect from an invalid reflect.Value")
)

// Error messages
const (
  reflectNonMaybeErrMsg = "Cannot create a MaybeReflect from a value of type %s"
)

// Constants
const (
	unionPkgPath = "github.com/bantling/micro/union"
)

// MaybeReflect wraps a Maybe and offers a similar api to Maybe via reflection.
// The internal goreflect.Value wrapper will be a *Maybe[T] when possible, and a Maybe[T] if it is not addressable.
// If the Value is not addressable, then the Maybe[T] can only be read, it cannot be written.
type MaybeReflect struct {
  typ goreflect.Type
  present goreflect.Value
  empty goreflect.Value
  get goreflect.Value
  orElse goreflect.Value
  orError goreflect.Value
  set goreflect.Value
  setEmpty goreflect.Value
  setOrError goreflect.Value
  string goreflect.Value
}

// Reflect returns (MaybeReflect, nil) if maybe wraps a Maybe value, or (zero value, error) if it does not.
func Reflect(maybe goreflect.Value) (MaybeReflect, error) {
  // Cannot reflect an invalid Value
  if !maybe.IsValid() {
    return MaybeReflect{}, reflectInvalidErr
  }

  // Cannot reflect a Value that does not wrap a Maybe
  dtyp := reflect.DerefType(maybe.Type())
	if !((dtyp.PkgPath() == unionPkgPath) && strings.HasPrefix(dtyp.Name(), "Maybe[")) {
    return MaybeReflect{}, fmt.Errorf(reflectNonMaybeErrMsg, maybe.Type())
	}

  // The Value wraps a Maybe.
  // If it has multiple pointers, deref it to one pointer.
  dmaybe := DerefValueMaxOnePtr(maybe)

  // If it is a value and addressable, get the address of it
  if (dmaybe.Kind() == goreflect.Struct) && dmaybe.CanAddr() {
    dmaybe = dmaybe.Addr()
  }

  // Create a result
  return MaybeReflect{
    typ: dtyp,
    present: dmaybe.MethodByName("Present"),
    empty: dmaybe.MethodByName("Empty"),
    get: dmaybe.MethodByName("Get"),
    orElse: dmaybe.MethodByName("OrElse"),
    orError: dmaybe.MethodByName("OrError"),
    set: dmaybe.MethodByName("Set"),
    setEmpty: dmaybe.MethodByName("SetEmpty"),
    setOrError: dmaybe.MethodByName("SetOrError"),
    string: dmaybe.MethodByName("String"),
  }
}

// Type returns the generic type of the MaybeReflect.
// Eg, if the MaybeReflect is a *Maybe[int] or Maybe[int], the result is equal to goreflect.TypeOf(0).
func (mr *MaybeReflect) Type() goreflect.Type {
  return mr.typ
}

// Present returns true if the MaybeReflect is present
func (mr *MaybeReflect) Present() bool {
  return mr.present.Call(nil)[0].Bool()
}

// Empty returns true if the MaybeReflect is empty
func (mr *MaybeReflect) Empty() bool {
  return mr.empty.Call(nil)[0].Bool()
}

// Get returns MaybeReflect value, panicking if it is not present
func (mr *MaybeReflect) Get() goreflect.Value {
  return mr.get.Call(nil)[0]
}

// OrElse returns the MaybeReflect value if present, else it returns elseVal.
// No check is made that elseVal is an appropriate type of value, panics if it is not
func (mr *MaybeReflect) OrElse(elseVal goreflect.Value) goreflect.Value {
  return mr.orElse.Call([]goreflect.Value{elseVal}])[0]
}

// OrEerror returns (value, nil) if present, else (zero value, error provided) if not present.
func (mr *MaybeReflect) OrError(err error) (goreflect.Value, error) {
  resAndErr := mr.orError.Call([]goreflect.Value{goreflect.ValueOf(err)})

  if resAndErr[1].IsNil() {
    return resAndErr[0], nil
  }

  return resAndErr[0], resAndErr[1].Interface().(error)
}

// Set sets the MaybeReflect to a new value, which is present as long as the new value is not nil.
// Panics if newVal is not the correct type.
func (mr *MaybeReflect) Set(newVal goreflect.Value) {
  mr.set.Call([]goreflect.Value{newVal})
}

// SetEmpty sets the MaybeReflect to an empty value
func (mr *MaybeReflect) SetEmpty() {
  mr.setEmpty.Call(nil)
}

// SetOrError sets the MaybeReflect to a new value, unless it already has a present value, in which case an error occurs
func (mr *MaybeReflect) SetOrError(newVal goreflect.Value) error {
  mr.setOrError.Call([]goreflect.Value{newVal}})[0].Interface().(error)
}

// String is the Stringer interface
func (mr *MaybeReflect) String() string {
  return mr.string.Call(nil)[0].String()
}
