package strukt

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	goreflect "reflect"

	"github.com/bantling/micro/funcs"
	unionreflect "github.com/bantling/micro/union/reflect"
)

const (
  errNotAStructPtrMsg   = "The pointer type %s does not point to a struct"
  errNilStructPtrMsg    = "The struct pointer of type %s cannot be nil"
  errMissingStructFieldMsg = "The struct of type %s does not contain a field named %s"
  errNotASubMapMsg         = "The map value of type %T is not the expected map type map[string]any for a struct"
)

// MapStructIgnoreMode indicates whether or not a map key that does not correspond to a struct field is ignored
type MapStructIgnoreMode bool

const (
  UnmappedFieldError  MapStructIgnoreMode = false
  UnmappedFieldIgnore MapStructIgnoreMode = true
)

// MapToStruct recurses through a map[string]any, populating the provided struct instance.
// conv.To is used to convert a map key value into a struct field value.
// The provided struct instance cannot be nil.
// The map key is converted from snake_case to CamelCase, so that only exported fields can be populated.
//
// If the map key is a sub map, then the related struct field must be a struct, maybe[struct], or pointer(s) to struct.
// If the map key value is an array/sliceof type T, then the related struct field must be an array/slice of type T.
// If any map key value cannot be converted, an error is returned.
// If the submap is empty, a struct value is set to the zero value, a struct * is set to nil.
// If the submap is not empty, a struct value is populated and if needed, pointers are allocated or it is wrapped in a Maybe.
//
// If the ignore flag is set to UnmappedFieldIgnore, then the map may contain values that do not exist in the struct, but
// it is still an error if the field exists in the struct and the map value cannot be converted.
// The default ignore flag value is UnmappedFieldError.
//
// The following errors can occur:
// - type T is not a struct type
// - strukt is a nil pointer
func MapToStruct[T any](mp map[string]any, strukt *T, ignore ...MapStructIgnoreMode) error {
  return mapToStruct(mp, strukt, funcs.SliceIndex(ignore, 0))
}

func mapToStruct(mp map[string]any, strukt any, ignore MapStructIgnoreMode) error {
  var (
    struktVal = goreflect.ValueOf(strukt)
    struktTyp = struktVal.Type().Elem()
  )

  // strukt * must point to a struct
  if struktTyp.Kind() != goreflect.Struct {
    return fmt.Errorf(errNotAStructPtrMsg, struktVal.Type())
  }

  // strukt * cannot be nil
  if struktVal.IsNil() {
    return fmt.Errorf(errNilStructPtrMsg, struktVal.Type())
  }

  // Iterate the fields of the struct
  for k, v := range mp {
    fieldName := funcs.SnakeToCamelCase(k)
    if fieldVal := struktVal.FieldByName(fieldName); !fieldVal.IsValid() {
      if !ignore {
        // The field does not exist in the struct, and we are not ignoring unmapped fields
        return fmt.Errorf(errMissingStructFieldMsg, struktTyp, fieldName)
      }
      // We are ignoring unmapped fields, do nothing
    } else if mapVal := goreflect.ValueOf(v); mapVal.Kind() == goreflect.Map {
      // The field may be a substruct, which can be any of Struct, *Struct, or Maybe[Struct]
      if maybeTyp := unionreflect.GetMaybeType(fieldVal.Type());
         ((maybeTyp != nil) && (maybeTyp.Kind() == goreflect.Struct)) ||
         (fieldVal.Kind() == goreflect.Struct) ||
         ((fieldVal.Kind() == goreflect.Pointer) && (fieldVal.Type().Elem().Kind() == goreflect.Struct)) {
        // Is the value of the map key a submap?
        submap, isa := v.(map[string]any)
        if !isa {
          return fmt.Errorf(errNotASubMapMsg, v)
        }

        switch {
        case maybeTyp != nil:
          // Recurse the submap into the substruct inside the Maybe
          if err := mapToStruct(submap, unionreflect.GetMaybeValue(mapVal).Addr().Interface(), ignore); err != nil {
            return err
          }

        case fieldVal.Kind() == goreflect.Struct:
          // Recurse the submap into the substruct value
          if err := mapToStruct(submap, fieldVal.Addr().Interface(), ignore); err != nil {
            return err
          }

        default:
        // Recurse the submap into the *substruct value
        if fieldVal.IsNil() {
          return fmt.Errorf(errNilStructPtrMsg, fieldVal.Type())
        }

        if err := mapToStruct(submap, fieldVal.Interface(), ignore); err != nil {
          return err
        }
      }
    }
  }
}
  return nil
}
