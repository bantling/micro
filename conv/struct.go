package conv

// import (
//   "fmt"
//   goreflect "reflect"
//
//   "github.com/bantling/micro/funcs"
//   unionreflect "github.com/bantling/micro/union/reflect"
// )
//
// // SPDX-License-Identifier: Apache-2.0
//
// var (
//   notAStructMsg = "Type %q is not a struct"
//   // current map key, top object name, path in object, field name
//   noStructFieldForKeyMsg = "Input has key %q, but %q path %q does not have a matching field named %q"
//   // current map key, top object name, path in object, field name
//   noSubMapForStructFieldMsg = "Input does not have a submap %q, but %q path %q is a field of type struct named %q"
// )
//
// const (
//   LAX = true
// )
//
//
//
// MapToStruct populates a struct from a map[string]any.
// The struct may contain sub structs as values or pointers.
// The struct may be recursive (eg Customer{child *Customer}).
// A map key is translated to a struct field by snake case to camel case conversion (eg first_name -> FirstName).
// The conv.LookupConversion function is used to locate a suitable conversion, if one exists.
//
// This func is not generic on struct type, to support cases where the struct type is not known ahead of time, such as:
// - The map is recursive, and the struct has child structs that cannot be known ahead of time.
// - A generalized algorithm that can work with any struct (eg, JSON -> struct, struct -> database row, etc).
//
// Returns an error if:
// - type S is not a struct
// - *any map key cannot be mapped to a struct field
// - *any map key cannot be converted to a struct field
//
// If the optional lax parameter is true (best to use provided constant LAX), then the starred errors above are ignored.
// If there are any fields of the struct without a matching map key, their values are unmodified.
// func ToStruct[S map[string]any | []any, T any](src S, dst *T, lax ...bool) error {
//   // Lax mode
//   isLax := funcs.SliceIndex(lax, 0)
//
//   // Die if dst is not a struct
//   topObj := goreflect.ValueOf(dst).Elem()
//   if topObj.Kind() != goreflect.Struct {
//     return fmt.Errorf(notAStructMsg, topObj.Type())
//   }
//
//   // Recurse through fields of src map
//   var (
//     k string
//     v any
//     path []string
//     fieldName string
//     field Value
//     err error
//     maybeType goreflect.Type
//     isStructField bool
//     convFn func(any, any) error
//
//     makeStructFieldErr = func(msg string) error {
//       // current map key, top object name, path in object, field name
//       return fmt.Errorf(msg, k, topObj, strings.Join(path, ","), fieldName)
//     }
//
//     recurseMap func(m map[string]any, obj goreflect.Value)
//     recurseSlice func([]any, obj goreflect.Value)
//   )
//
//   recurse = func(o map[string]any, a []any, obj goreflect.Value) {
//     for k, v = range m {
//       // Add key as next path part
//       path = append(path, k)
//
//       // Convert the key name from snake_case to CamelCase, to match Go exported struct field name conventions
//       // Check if such a field exists
//       fieldName = funcs.SnakeToCamelCase(k)
//       if field = obj.FieldByName(fieldName); !field.IsValid() {
//         err = makeStructFieldErr(noStructFieldForKeyMsg)
//         goto done
//       }
//
//       // Is the field a struct or Maybe[struct]?
//       maybeType = unionreflect.GetMaybeType(field.Type())
//       isStructField = (field.Kind() == goreflect.Struct) || ((maybeType != nil) && (mayType.Kind() == goreflect.Struct))
//       if isStruct {
//         // Is the map key value a map[string]any?
//         if submap, isa := v.(map[string]any); isa {
//           if maybeType == nil {
//             // If the field is a struct, recurse submap into it
//             recurse(submap, field)
//           } else {
//             // If the field is a Maybe[struct], then create a struct instance to recurse the submap into
//             maybeObj := goreflect.New(maybeType).Elem()
//             recurse(submap, maybeObj)
//
//             // If no errors occurred, copy maybeObj into Maybe field
//             unionreflect.SetMaybeValue(field, maybeObj)
//           }
//         } else if !isLax {
//           // - The field is a struct
//           // - The map key value is not a submap
//           // - Lax mode is false
//           err = makeStructFieldErr(noSubMapForStructFieldMsg)
//           goto done
//         }
//       }
//
//       // Is the field a slice?
//       if field.Kind() == goreflect.Slice {
//         // Checking the field for compatible slice element types is too much work, too many possibilities
//         // Fail on elements if they can't be converted
//         // Is the map key value a []any?
//         if slc, isa := v.([]any); isa {
//           // Recurse for each element of the slice, and try converting
//
//         }
//       }
//
//       // Try to convert value into field using LookupConversion
//       if convFn, err = LookupConversion(reflect.TypeOf(v), field.Type()); err != nil {
//         goto done
//       }
//
//       // Remove the last path part
//       path = path[0:len(path)-1]
//     }
//   }
//
//   // Start recursion at provided map and object
//   recurse(src, topObj)
//
//   done:
//   return err
// }
