package json

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math/big"
	goreflect "reflect"
	"strings"

	"github.com/bantling/micro/go/funcs"
	"github.com/bantling/micro/go/reflect"
)

var (
	ErrNotValidArrayElementTypeMsg = "A value of type %s is not a valid JSON array element"
)

// FromStruct converts a struct to a Value, as follows:
// - The Value type will be Object
// - Field names are map keys, with auto lowercasing of first letter (eg FirstName field name -> firstName JSON key)
// - In the special case where every letter of the field name is uppercase, the whole name is lowercased (eg BI -> bi)
// - a struct, map[string]Value, or map[string]any field is a recursive Object Value
// - a []Value, []any, []string, []NumberType, []bool field is a recursive Array Value
// - a string field is a String Value
// - a NumberType field is a Number Value (*big.Rat is normalized to a float string)
// - a bool field is a Boolean Value
// - a pointer field may point to any of the above types, and if it is null then the Null Value is used
// - any other kind of field is ignored, and assumed to be for an unrelated purpose
//
// # If a pointer field has multiple pointers, the Null Value is used if any pointer is nil
//
// Returns an error if:
// - The value passed is not zero or more pointers to a struct
// - The value has a nil pointer
func FromStruct(strukt any) (Value, error) {
	var (
		struktVal   = reflect.DerefValue(goreflect.ValueOf(strukt))
		handleSlice func(reflectFld goreflect.Value) ([]any, error)
	)

	handleSlice = func(reflectGiven goreflect.Value) ([]any, error) {
		// If the field value is a nil pointer to slice or a nil slice, return nil
		if !reflectGiven.IsValid() || reflectGiven.IsNil() {
			return nil, nil
		}

		// Each slice element may contain bool | string | constraint.Numeric | Value | struct | slice
		slc := []any{}
		for i, l := 0, reflectGiven.Len(); i < l; i++ {
			// In case the element value is one or more pointers, resolve it to no more than one pointer
			reflectElemMaxOnePtr := reflect.DerefValueMaxOnePtr(reflectGiven.Index(i))

			// Deref the max one pointer to get the actual value, if available
      // Resolve the type in case the slice was typed as []interface{}, causing all elements to be typed as interface{}
			reflectElem := reflect.ResolveValueType(reflect.DerefValue(reflectElemMaxOnePtr))

      // In case the elem was typed as []interface{}, reconstruct the max one ptr from it
      if reflectElemMaxOnePtr.Kind() == goreflect.Pointer {
        reflectElemMaxOnePtr = reflectElem.Addr()
      } else {
        reflectElemMaxOnePtr = reflectElem
      }

			// Check if the field is a big pointer
			reflectElemMaxOnePtrTyp := reflectElemMaxOnePtr.Type()

			isBigPtr := (reflectElemMaxOnePtrTyp == goreflect.TypeOf((*big.Int)(nil))) ||
				(reflectElemMaxOnePtrTyp == goreflect.TypeOf((*big.Float)(nil))) ||
				(reflectElemMaxOnePtrTyp == goreflect.TypeOf((*big.Rat)(nil)))

			reflectElemTyp := reflect.DerefType(reflectElemMaxOnePtrTyp)
			reflectElemKind := reflectElemTyp.Kind()

			// See if the element is any type we can work with
			switch {
			case (reflectElemKind == goreflect.Struct) && (!isBigPtr):
				// If the element is a nil pointer to struct, map the field name to json null
				if !reflectElem.IsValid() {
					slc = append(slc, nil)
					continue
				}

				// Make a recursive call and add the Value as is
				if subJSONVal, subErr := FromStruct(reflectElemMaxOnePtr.Interface()); subErr == nil {
					slc = append(slc, subJSONVal)
				} else {
					return nil, subErr
				}

			case reflectElemKind == goreflect.Slice:
				// Make a recursive call and add the Value as is
				if subSlice, subErr := handleSlice(reflectElem); subErr == nil {
					slc = append(slc, subSlice)
				} else {
					return nil, subErr
				}

			case reflectElemKind == goreflect.String:
				// If the element value is a nil pointer to string, add a json null
				if !reflectElem.IsValid() {
					slc = append(slc, nil)
					continue
				}

				slc = append(slc, reflectElem.String())

			case reflectElemKind == goreflect.Bool:
				// If the element value is a nil pointer to bool, add a json null
				if !reflectElem.IsValid() {
					slc = append(slc, nil)
					continue
				}

				slc = append(slc, reflectElem.Bool())

			default:
				// Is the element a NumberType?
				isNumberType := ((reflectElemKind >= goreflect.Int) && (reflectElemKind <= goreflect.Float64) && (reflectElemKind != goreflect.Uintptr)) ||
					isBigPtr || reflectElemTyp == goreflect.TypeOf(NumberString(""))

				// If it isn't a NumberType or pointer to it, then it is an unrelated type - skip it, we cannot convert to JSON
				if !isNumberType {
					continue
				}

				// If the element value is a nil pointer, add a json null
				if !reflectElem.IsValid() {
					slc = append(slc, nil)
					continue
				}

				// Must be convertible to json.NumberType
				slc = append(slc, funcs.Ternary(isBigPtr, reflectElemMaxOnePtr.Interface(), reflectElem.Interface()))
			}
		}

		return slc, nil
	}

	if !struktVal.IsValid() {
		// One or more pointers where at least one of them is nil
		return InvalidValue, fmt.Errorf(errNilPtrMsg, strukt)
	}

	// strukt must deref to a struct
	if struktVal.Kind() != goreflect.Struct {
		return InvalidValue, fmt.Errorf(errNotAStructMsg, strukt)
	}

	// non-nil map to populate and convert to a json Value
	jsonMap := map[string]any{}

	// Iterate the fields
	for fieldName := range reflect.FieldsByName(struktVal.Type()) {
		// Translate field name to a json key name
		var jsonKey string
		if allUpper := strings.ToUpper(fieldName); allUpper == fieldName {
			// Special case of all uppercase field name, lowercase whole thing for json key
			jsonKey = strings.ToLower(fieldName)
		} else {
			// Usual case of mixed case field name, lowercase first letter for json key (eg FirstName -> firstName)
			fieldNameRunes := []rune(fieldName)
			jsonKey = strings.ToLower(string(fieldNameRunes[0])) + string(fieldNameRunes[1:])
		}

		// In case the field value is one or more pointers, resolve it to no more than one pointer
		reflectFldMaxOnePtr := reflect.DerefValueMaxOnePtr(struktVal.FieldByName(fieldName))

		// Deref the max one pointer to get the actual value, if available
		reflectFld := reflect.DerefValue(reflectFldMaxOnePtr)

		// Check if the field is a big pointer
		reflectFldMaxOnePtrTyp := reflectFldMaxOnePtr.Type()
		isBigPtr := (reflectFldMaxOnePtrTyp == goreflect.TypeOf((*big.Int)(nil))) ||
			(reflectFldMaxOnePtrTyp == goreflect.TypeOf((*big.Float)(nil))) ||
			(reflectFldMaxOnePtrTyp == goreflect.TypeOf((*big.Rat)(nil)))

		reflectFldTyp := reflect.DerefType(reflectFldMaxOnePtrTyp)
		reflectFldKind := reflectFldTyp.Kind()

		// See if the field is any type we can work with
		switch {
		case (reflectFldKind == goreflect.Struct) && (!isBigPtr):
			// If the field value is nil pointer to struct, map the field name to json null
			if !reflectFld.IsValid() {
				jsonMap[jsonKey] = nil
				continue
			}

			// Make a recursive call and add the Value as is
			if subJSONVal, subErr := FromStruct(reflectFldMaxOnePtr.Interface()); subErr == nil {
				jsonMap[jsonKey] = subJSONVal
			} else {
				return InvalidValue, subErr
			}

		case reflectFldKind == goreflect.Slice:
			if slc, subErr := handleSlice(reflectFld); subErr == nil {
				jsonMap[jsonKey] = slc
			} else {
				return InvalidValue, subErr
			}

		case reflectFldKind == goreflect.String:
			// If the field value is a nil pointer to string, map the field name to json null
			if !reflectFld.IsValid() {
				jsonMap[jsonKey] = nil
				continue
			}

			jsonMap[jsonKey] = reflectFld.String()

		case reflectFldKind == goreflect.Bool:
			// If the field value is a nil pointer to bool, map the field name to json null
			if !reflectFld.IsValid() {
				jsonMap[jsonKey] = NullValue
				continue
			}

			jsonMap[jsonKey] = reflectFld.Bool()
		default:
			// Is the field a NumberType?
			isNumberType := ((reflectFldKind >= goreflect.Int) && (reflectFldKind <= goreflect.Float64) && (reflectFldKind != goreflect.Uintptr)) ||
				isBigPtr || reflectFldTyp == goreflect.TypeOf(NumberString(""))

			// If it isn't a NumberType or pointer to it, then it is an unrelated type - skip it, we cannot convert to JSON
			if !isNumberType {
				continue
			}

			// If the field value is a nil pointer, map the field name to json null
			if !reflectFld.IsValid() {
				jsonMap[jsonKey] = nil
				continue
			}

			// Must be convertible to json.NumberType
			jsonMap[jsonKey] = funcs.Ternary(isBigPtr, reflectFldMaxOnePtr.Interface(), reflectFld.Interface())
		}
	}

	// Return successful Value conversion
	return FromMap(jsonMap), nil
}

// ToStruct copies the structure of a json.Value Object to a go reflect.Value wrapper of zero or more pointers to a struct.
// In the case of one or more pointers, if any pointers are nil, they are allocated so the converted json.Value can be stored.
//
// The struct does not have to contain fields for all the Value keys, any Value key that does not have a corresponding
// field is ignored. Field name first letters are downcased automatically, so that FirstName becomes firstName Value key.
//
// The struct may contain additional fields that are unrelated to the Object keys, they are ignored.
// If the Value contains sub objects, they can be stored in a sub struct.
//
// Any struct in the hierarchy can have a field named JSON that is a compatible type for an Object or Array, and the
// JSON field will get a copy of the entire Object or Array structure. For an Object, the same struct can also have other
// fields that are named after Object keys, and those fields will also be set to the Object key values, allowing the struct
// to have both a copy of the whole structure, and individual fields.
//
// The possible compatible struct field types for a given Object key are:
// Value                    - copy Value as is
// map[string]Value         - copy underlying Object map
// map[string]any           - convert underlying Object map using ToMap
// []Value                  - copy underlying Array slice
// []any                    - convert underlying Array slice using ToSlice
// struct                   - convert Object to struct, recursively
// slice                    - convert Array to slice (elements must be Value, any, string, constraint.Numeric types, or bool)
// string                   - convert String to string
// constraint.Numeric types - convert Number to target type
// bool                     - convert Boolean to bool
//
// An error is returned if:
// - The value passed is not zero or more pointers to a struct
// func ToStruct(src Value, tgt goreflect.Value) error {
//   return nil
// }
