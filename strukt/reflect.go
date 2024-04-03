package strukt

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	goreflect "reflect"

	"github.com/bantling/micro/funcs"
	unionreflect "github.com/bantling/micro/union/reflect"
)

const (
  errInvalidMapTypeMsg     = "The map key %s has an invalid map value of type %T"
	errMissingStructFieldMsg = "The struct of type %s does not contain a field named %s"
	errNilStructPtrMsg       = "The struct pointer of type %s cannot be nil"
  errNotASliceMsg          = "The struct field %s.%s is not a slice"
	errNotAStructPtrMsg      = "The pointer type %s does not point to a struct"
	errNotASubMapMsg         = "The struct field %s.%s is a sub struct but the map value of type %T is not the expected type map[string]any"
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
    var (
	    fieldName = funcs.SnakeToCamelCase(k)
      mapVal    = goreflect.ValueOf(v)
    )
		if fieldVal := struktVal.FieldByName(fieldName); !fieldVal.IsValid() {
			if !ignore {
				// The field does not exist in the struct, and we are not ignoring unmapped fields
				return fmt.Errorf(errMissingStructFieldMsg, struktTyp, fieldName)
			}
			// We are ignoring unmapped fields, do nothing
		} else if subMap, isa := v.(map[string]any); isa {
			// The field may be a substruct, which can be any of Struct, *Struct, or Maybe[Struct]
			if maybeTyp := unionreflect.GetMaybeType(fieldVal.Type()); ((maybeTyp != nil) && (maybeTyp.Kind() == goreflect.Struct)) ||
				(fieldVal.Kind() == goreflect.Struct) ||
				((fieldVal.Kind() == goreflect.Pointer) && (fieldVal.Type().Elem().Kind() == goreflect.Struct)) {
				switch {
				case maybeTyp != nil:
					// Recurse the subMap into the substruct inside the Maybe
					if err := mapToStruct(subMap, unionreflect.GetMaybeValue(fieldVal).Interface(), ignore); err != nil {
						return err
					}

				case fieldVal.Kind() == goreflect.Struct:
					// Recurse the subMap into the substruct value
					if err := mapToStruct(subMap, fieldVal.Interface(), ignore); err != nil {
						return err
					}

				default:
					if fieldVal.IsNil() {
            // Allocate a new struct and set pointer to it
            fieldVal.Set(goreflect.New(fieldVal.Type().Elem()))
					}

            // Recurse the subMap into the *substruct value
					if err := mapToStruct(subMap, fieldVal.Elem().Interface(), ignore); err != nil {
						return err
					}
				}
      } else if fieldVal.Type() == goreflect.TypeOf((map[string]any)(nil)) {
        // Just copy map directly into struct
        fieldVal.Set(goreflect.ValueOf(subMap))
			} else {
        // We can't populate a struct field any other type with a map[string]any as input
        return fmt.Errorf(errInvalidMapTypeMsg, k, v)
      }
		} else if mapVal.Kind() == goreflect.Slice {
      // The field value must be a slice
      if fieldVal.Kind() != goreflect.Slice {
        return fmt.Errorf(err)
      }

      // The slice elements may be map[string]any
      if sliceOfMaps, isa := v.([]map[string]any); isa {
        for i, v := range sliceOfMaps {
          if err := mapToStruct(v, strukt, ignore); err != nil {
            return fmt.Errorf(errNotASliceMsg, struktTyp, fieldName)
          }
        }
      } else {

      }
    }
	}
	return nil
}