package conv

// SPDX-License-Identifier: Apache-2.0

import (
  "fmt"
  goreflect "reflect"

  unionreflect "github.com/bantling/micro/union/reflect"
)

var (
  errMaybeNotConvertibleMsg = "A %s cannot be converted to %s"
  errMaybeNotSettableMsg = "A %s cannot be set to a(n) %s"
)

// MaybeWrapperInfo provides the semantics of the Maybe wrapper:
// - Any type is accepted
// - No conversions can occur, only returns exactly what is stored
// - Empty values can be stored
type MaybeWrapperInfo int

func (mwi MaybeWrapperInfo) PackagePath() string {
  return "github.com/bantling/micro/union"
}

func (mwi MaybeWrapperInfo) TypeNamePrefix() string {
  return "Maybe"
}

func (mwi MaybeWrapperInfo) AcceptsType(instance goreflect.Value, typ goreflect.Type) bool {
  return unionreflect.GetMaybeType(instance.Type()) == typ
}

func (mwi MaybeWrapperInfo) CanBeEmpty(instance goreflect.Value) bool {
  return true
}

func (mwi MaybeWrapperInfo) ConvertibleTo(instance goreflect.Value, typ goreflect.Type) bool {
  // If the instance is not a Maybe type, then reflect.GetMaybeType returns nil. If type given is also nil, return false.
  return (typ != nil) && (unionreflect.GetMaybeType(instance.Type()) == typ)
}

func (mwi MaybeWrapperInfo) Get(instance goreflect.Value, typ goreflect.Type) (goreflect.Value, bool, error) {
  if mwi.ConvertibleTo(instance, typ) {
    mval := unionreflect.GetMaybeValue(instance)
    return mval, mval.IsValid(), nil
  }

  return goreflect.Value{}, false, fmt.Errorf(errMaybeNotConvertibleMsg, instance.Type(), typ)
}

func (mwi MaybeWrapperInfo) Set(instance, val goreflect.Value, present bool) error {
  if !present {
    unionreflect.SetMaybeValueEmpty(instance)
    return nil
  }

  if mwi.ConvertibleTo(instance, val.Type()) {
    unionreflect.SetMaybeValue(instance, val)
    return nil
  }

  return fmt.Errorf(errMaybeNotSettableMsg, instance.Type(), val.Type())
}
