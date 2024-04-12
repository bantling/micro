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
// The map key is converted from snake_case to CamelCase, as only exported fields can be populated.
//
// If the map key is a sub map, then the related struct field must be a (*)struct or Maybe[(*)struct].
// If the map key value is an array/slice of type T, then the related struct field must be an array/slice of type T.
// If any map key value cannot be converted, an error is returned, and the target struct will be partially populated.
// If the submap is empty, a (Maybe)struct value is set to the zero value, while a (Maybe)*struct is set to nil.
// If the submap is not empty, a struct value is populated and if needed, pointers are allocated and/oror it is wrapped in a Maybe.
//
// If the ignore flag is set to UnmappedFieldIgnore, then the map may contain keys that do not exist in the struct, but
// it is still an error if the field exists in the struct and the map value cannot be converted.
// The default ignore flag value is UnmappedFieldError, which requires every map key to exist in the struct.
//
// The following errors can occur:
// - type T is not a struct type
// - strukt is a nil pointer
func MapToStruct[T any](mp map[string]any, strukt *T, ignore ...MapStructIgnoreMode) error {
	return mapToStruct(mp, goreflect.ValueOf(strukt), funcs.SliceIndex(ignore, 0))
}

// mapToStruct is the internal non-generic function that recurses the map and the struct.
// This function is necessary because you can't recurse a generic function when the value of T for a substruct is not
// known at compile time.
func mapToStruct(mp map[string]any, struktPtr goreflect.Value, ignore MapStructIgnoreMode) error {
	struktTyp := struktPtr.Type().Elem()

	// struktPtr must be a pointer to a struct
	if struktTyp.Kind() != goreflect.Struct {
		return fmt.Errorf(errNotAStructPtrMsg, struktPtr.Type())
	}

	// struktPtr cannot be nil - the top call stores the result in the struct the pointer references
	if struktPtr.IsNil() {
		return fmt.Errorf(errNilStructPtrMsg, struktPtr.Type())
	}

  // Deref the strukt so we can access fields of it
  struktVal := struktPtr.Elem()

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
			// The field must be a substruct, which can be any of Struct, *Struct, Maybe[Struct], or Maybe[*Struct]
      if fieldVal.Kind() == goreflect.Struct {
        // Case Struct
        // Recurse the subMap into the substruct value address, overwriting the fields
        if err := mapToStruct(subMap, fieldVal.Addr().Interface(), ignore); err != nil {
          return err
        }
      } else if (fieldVal.Kind() == goreflect.Pointer) && (fieldVal.Type().Elem().Kind() == goreflect.Struct) {
        // Case *Struct
        // If the field is nil, allocate a new struct and set the field to a copy of the pointer
        if fieldVal.IsNil() {
          fieldVal.Addr().Set(goreflect.New(fieldVal.Type().Elem()))
        }

        // Recurse the subMap into the substruct value pointer
        if err := mapToStruct(subMap, fieldVal.Interface(), ignore); err != nil {
          return err
        }
      } else if maybeTyp := unionreflect.GetMaybeType(fieldVal.Type()); (maybeTyp != nil) && (maybeTyp.Kind() == goreflect.Struct) {
        // Case Maybe[Struct]
        // Get a copy the Struct value - it may be initialized to a desired state with defaults
        maybeVal := unionreflect.GetMaybeValue(fieldVal)

        // Recurse into a pointer to the copied value
        if err := mapToStruct(subMap, maybeVal.Addr().Interface(), ignore); err != nil {
          return err
        }

        // Copy the modified Struct value into the Maybe
        unionreflect.SetMaybeValue(fieldVal.Addr(), maybeVal)
      } else if (maybeTyp != nil) && (maybeTyp.Kind() == goreflect.Pointer) && (maybeTyp.Elem().Kind() == goreflect.Struct) {
        // Case Maybe[*Struct]
        // If the Maybe value is nil, allocate a new struct and set the value to a copy of the pointer
        maybeVal := unionreflect.GetMaybeValue(fieldVal)
        if (! maybeVal.IsValid()) {
          unionreflect.SetMaybeValue(maybeVal, goreflect.New(maybeTyp))
        }

        // Recurse the subMap into the substruct pointed to inside the Maybe
        if err := mapToStruct(subMap, maybeVal.Elem().Interface(), ignore); err != nil {
          return err
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
        return nil //fmt.Errorf(err)
      }

      // The slice elements may be map[string]any
      if sliceOfMaps, isa := v.([]map[string]any); isa {
        for _, v := range sliceOfMaps {
          if err := mapToStruct(v, struktVal, ignore); err != nil {
            return fmt.Errorf(errNotASliceMsg, struktTyp, fieldName)
          }
        }
      } else {

      }
    }
	}
	return nil
}
