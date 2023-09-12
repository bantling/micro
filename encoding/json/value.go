package json

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
  "math/big"
  goreflect "reflect"
  "regexp"

	"github.com/bantling/micro/constraint"
	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/union"
)

// NumberString regex
var numberStringRegex    = regexp.MustCompile("-?[0-9]+([.][0-9]+)?(e[0-9]+)?")

// Error constants
var (
	errInvalidGoValueMsg = "A value of type %T is not a valid type to convert to a Value. Acceptable types are map[string]any, []any, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *big.Int, *big.Float, *big.Rat, NumberString, bool, and nil"
	errNotObject         = fmt.Errorf("The Value is not an object")
	errNotArray          = fmt.Errorf("The Value is not an array")
	errNotString         = fmt.Errorf("The Value is not a string")
	errNotNumber         = fmt.Errorf("The Value is not a number")
	errNotBoolean        = fmt.Errorf("The Value is not a boolean")
	errNotStringable     = fmt.Errorf("The Value is not a string, number, or boolean")

)

// Type is an enum of json types
type Type uint

// The value types
const (
  Invalid Type = iota
	Object
	Array
	String
	Number
	Boolean
	Null
)

// The value types as strings
var (
	valueTypeToString = map[Type]string{
		Object:  "Object",
		Array:   "Array",
		String:  "String",
		Number:  "Number",
		Boolean: "Boolean",
		Null:    "Null",
	}
)

// ToString is the Stringer interface for fmt
func (typ Type) String() string {
	return valueTypeToString[typ]
}

// NumberString is a type that allows a plain string to be considered a JSON Number.
// Allows differentiation between an actual string value, and a string that is really a number value.
type NumberString string

// Value represents any kind of JSON value - object, array, string, number, boolean, null
type Value struct {
	typ Type
  val union.Four[map[string]Value, []Value, string, bool]
}

// Constant values for a invalid, true, false, and null
var (
  invalidValue = Value{}
	TrueValue    = Value{typ: Boolean, val: union.Of4W[map[string]Value, []Value, string, bool](true)}
	FalseValue   = Value{typ: Boolean, val: union.Of4W[map[string]Value, []Value, string, bool](false)}
	NullValue    = Value{typ: Null, val: union.Four[map[string]Value, []Value, string, bool]{}}
)

// toValue is common code for ToValue, MapToValue, and SliceToValue
func toValue(val any) (res Value, err error) {
  // Object
  if mp, isa := val.(map[string]any); isa {
    // recursive call
    if res, err = MapToValue(mp); err != nil {
      return
    }

  // Array
  } else if arr, isa := val.([]any); isa {
    // recursive call
    if res, err = SliceToValue(arr); err != nil {
      return
    }

  // Boolean
  } else if bln, isa := val.(bool); isa {
    res = BoolToValue(bln)

  // Null : cannot test an any value of some unknown type as == nil, have to type assert or use reflection
  } else if funcs.IsNilValue(val) {
    // nil only occurs if a *big is nil
    res = NullValue

  } else {
    // string and NumberString
    typ := goreflect.TypeOf(val)

    if typ.Kind() == goreflect.String {
      if typ.AssignableTo(goreflect.TypeOf(NumberString(""))) {
        res, err = numberToValue(val.(NumberString))
      } else {
        res = StringToValue(val.(string))
      }

      return

    // Must be Number
    } else {
      res, err = numberToValue(val)
    }
  }

  return
}

// ToValue converts value types any into a Value
func ToValue[T map[string]any | map[string]Value | []any | []Value | string | NumberString | constraint.Numeric | bool](val T) (res Value, err error) {
  return toValue(any(val))
}

// MustToValue is a must version of ToValue
func MustToValue[T map[string]any | map[string]Value | []any | []Value | string | NumberString | constraint.Numeric | bool](val T) Value {
  return funcs.MustValue(ToValue(val))
}

// MapToValue converts a map[string]any to a Value of type Object, where the map values must be:
// - map[string]any for a sub Object
// - []any for a sub Array
// - string for a String
// - NumberString or any constraint.Numeric type for a Number
// - bool for a Boolean
// - nil for a Null
//
// Any other key value results in (Invalid Value, error)
func MapToValue[T map[string]any | map[string]Value](val T) (mval Value, err error) {
  var (
    av = any(val)
    mp map[string]Value
    jval Value
  )

  if v, isa := av.(map[string]Value); isa {
    mp = v
  } else {
    mp = map[string]Value{}

    for k, v := range av.(map[string]any) {
      if jval, err = toValue(v); err != nil {
        return
      }

      mp[k] = jval
    }
  }

  mval = Value{typ: Object, val: union.Of4T[map[string]Value, []Value, string, bool](mp)}

  return
}

// MustMapToValue is a must version of MapToValue
func MustMapToValue[T map[string]any | map[string]Value](val T) Value {
  return funcs.MustValue(MapToValue(val))
}

// SliceToValue converts a []any to a Value of type Array, where the slice values must be the same as map key values
// (see MapToValue).
//
// Any other slice value results in the (Invalid Value, error)
func SliceToValue[T []any | []Value](val T) (sval Value, err error) {
  var (
    av = any(val)
    slc []Value
    jval Value
  )

  if v, isa := av.([]Value); isa {
    slc = v
  } else {
    slc = make([]Value, len(val))

    for i, v := range av.([]any) {
      if jval, err = toValue(v); err != nil {
        return
      }

      slc[i] = jval
    }
  }

  sval = Value{typ: Array, val: union.Of4U[map[string]Value, []Value, string, bool](slc)}

  return
}

// MustSliceToValue is a must version of SliceToValue
func MustSliceToValue[T []any | []Value](val T) Value {
  return funcs.MustValue(SliceToValue(val))
}

// StringToValue converts a string to a Value
func StringToValue(val string) Value {
  return Value{typ: String, val: union.Of4V[map[string]Value, []Value, string, bool](val)}
}

// numberToValue is common code for ToValue and NumberToValue
func numberToValue(val any) (nval Value, err error) {
	var str string

	if br, isa := val.(*big.Rat); isa {
	  // *big.Rat must be converted to a normalized string.
		str = conv.BigRatToNormalizedString(br)
	} else {
    // All other numeric types that remain are convertible to string by conv.To
    if err = conv.To(val, &str); err != nil {
      return
    }

    // A NumberString could be the empty string, or "foo" or some other value that is not a number
    if ns, isa := val.(NumberString); isa && (!numberStringRegex.MatchString(string(ns))) {
      err = errNotNumber
      return
    }
  }

  nval = Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool](str)}
  return
}

// NumberToValue converts any numeric type into a Value
func NumberToValue[T NumberString | constraint.Numeric](val T) (Value, error) {
  return numberToValue(val)
}

// MustNumberToValue is a must version of NumberToValue
func MustNumberToValue[T NumberString | constraint.Numeric](val T) Value {
  return funcs.MustValue(NumberToValue(val))
}

// BoolToValue converts a bool into a Value
func BoolToValue(val bool) Value {
  return funcs.Ternary(val, TrueValue, FalseValue)
}

// Type returns the type of value this Value contains
func (jv Value) Type() Type {
	return jv.typ
}

// AsMap returns a map representation of a Value.
// Panics if the Value is not an object.
func (jv Value) AsMap() map[string]Value {
	if jv.typ != Object {
		panic(errNotObject)
	}

	return jv.val.T()
}

// AsSlice returns a slice representation of a Value.
// Panics if the Value is not an array.
func (jv Value) AsSlice() []Value {
	if jv.typ != Array {
		panic(errNotArray)
	}

	return jv.val.U()
}

// AsString returns a string representation of a Value.
// Panics if the Value is not a string, number, or boolean.
func (jv Value) AsString() string {
	switch jv.typ {
	case String:
		fallthrough
	case Number:
		return jv.val.V()
	case Boolean:
		return fmt.Sprintf("%t", jv.val.W())
	}

	panic(errNotStringable)
}

// AsNumber returns a NumberString representation of a Value.
// Panics if the Value is not a number.
func (jv Value) AsNumber() NumberString {
	if jv.typ != Number {
		panic(errNotNumber)
	}

	return NumberString(jv.val.V())
}

// AsBoolean returns a bool representation of a Value.
// Panics if the Value is not a boolean.
func (jv Value) AsBool() bool {
	if jv.typ != Boolean {
		panic(errNotBoolean)
	}

	return jv.val.W()
}

// IsNull returns true if the Value is a null, else false
func (jv Value) IsNull() bool {
	return jv.typ == Null
}

// ToAny converts the Value to the approriate go type:
// Object  = map[string]any
// Array   = []any
// String  = string
// Number  = NumberString
// Boolean = bool
// Null    = nil
func (jv Value) ToAny() any {
	switch jv.typ {
	case Object:
		return jv.ToMap()
	case Array:
		return jv.ToSlice()
	case String:
		return jv.val.V()
	case Number:
		return NumberString(jv.val.V())
  case Boolean:
    return jv.val.W()
	default:
		return nil
	}
}

// ToMap returns a map[string]any representation of a Value.
// Panics if the Value is not an object.
func (jv Value) ToMap() map[string]any {
	if jv.typ != Object {
		panic(errNotObject)
	}

	m := map[string]any{}

	for k, v := range jv.val.T() {
		m[k] = v.ToAny()
	}

	return m
}

// ToSlice returns a []any representation of a Value.
// Panics if the Value is not an Array.
// If the optional visitor func is not provided, then DefaultVisitor is used.
func (jv Value) ToSlice(visitor ...func(Value) any) []any {
	if jv.typ != Array {
		panic(errNotArray)
	}

	s := []any{}

	for _, v := range jv.val.U() {
		s = append(s, v.ToAny())
	}

	return s
}

// IsDocument returns true if a Value is a document (Object or Array)
func (jv Value) IsDocument() bool {
	return (jv.typ == Object) || (jv.typ == Array)
}
