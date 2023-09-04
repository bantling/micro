package json

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/funcs"
)

// Error constants
var (
	errInvalidGoValueMsg = "A value of type %T is not a valid type to convert to a Value. Acceptable types are map[string]any, []any, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *big.Int, *big.Float, *big.Rat, NumberString, bool, and nil"
	errNotObject         = fmt.Errorf("The Value is not an object")
	errNotArray          = fmt.Errorf("The Value is not an array")
	errNotString         = fmt.Errorf("The Value is not a string")
	errNotNumber         = fmt.Errorf("The Value is not a number")
	errNotBoolean        = fmt.Errorf("The Value is not a boolean")
	errNotStringable     = fmt.Errorf("The Value is not a string, number, or boolean")
	errNotAStructMsg     = "The value of type %T is not a struct"
	errNilPtrMsg         = "The value of type %T has multiple pointers, where one of the leading pointers is nil"
)

// ValueType is an enum of value types
type ValueType uint

// The value types
const (
	Object ValueType = iota
	Array
	String
	Number
	Boolean
	Null
)

// The value types as strings
var (
	valueTypeToString = map[ValueType]string{
		Object:  "Object",
		Array:   "Array",
		String:  "String",
		Number:  "Number",
		Boolean: "Boolean",
		Null:    "Null",
	}
)

// ToString is the Stringer interface for fmt
func (typ ValueType) String() string {
	return valueTypeToString[typ]
}

// NumberString is a special type that allows a plain string to be considered a JSON Number
type NumberString string

// Value represents any kind of JSON value - object, array, string, number, boolean, null
type Value struct {
	typ   ValueType
	value any
}

// Constant values for a true, false, and Null
var (
	TrueValue    = Value{typ: Boolean, value: true}
	FalseValue   = Value{typ: Boolean, value: false}
	NullValue    = Value{typ: Null, value: nil}
	InvalidValue = Value{}
)

// FromValue converts a Go value into a Value, where the Go value must be as follows:
//
// Object: map[string]any
// Array: []any
// String: string
// Number: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *big.Int, *big.Float,
//
//	*big.Rat, or NumberString
//
// Boolean: bool
// Null: nil
//
// To make recursive algorithms whose base case returns calls to functions like FromMap and FromSlice easier and more
// efficient to implement, the value can also already be a Value, which is used as is.
//
// Panics if any other kind of value is provided
func FromValue(v any) Value {
	var jval Value

	if jv, isa := v.(map[string]any); isa {
		jval = FromMap(jv)
	} else if jv, isa := v.([]any); isa {
		jval = FromSlice(jv)
	} else if jv, isa := v.(string); isa {
		jval = FromString(jv)
	} else if jv, isa := v.(bool); isa {
		jval = FromBool(jv)
	} else if v == nil {
		jval = NullValue
	} else if jv, isa := v.(Value); isa {
		jval = jv
	} else if jval = FromNumber(v); (jval == Value{}) {
		panic(fmt.Errorf(errInvalidGoValueMsg, v))
	}

	return jval
}

// FromMap converts a map[string]any into a Value.
// The types of the map keys must be acceptable to FromValue.
func FromMap(m map[string]any) Value {
	jv := map[string]Value{}

	for k, v := range m {
		jv[k] = FromValue(v)
	}

	return Value{typ: Object, value: jv}
}

// FromMapOfValue converts a map[string]Value into a Value
func FromMapOfValue(m map[string]Value) Value {
	return Value{typ: Object, value: m}
}

// FromSlice converts a []any into a Value.
// The types of the slice elements must be acceptable to FromValue.
func FromSlice(a []any) Value {
	js := make([]Value, len(a))

	for i, v := range a {
		js[i] = FromValue(v)
	}

	return Value{typ: Array, value: js}
}

// FromSliceOfValue converts a []Value into a Value
func FromSliceOfValue(a []Value) Value {
	return Value{typ: Array, value: a}
}

// FromDocument converts a map[string]any or []any into a Value.
// See FromMap and FromSlice.
func FromDocument[T map[string]any | []any](doc T) Value {
	if m, isa := any(doc).(map[string]any); isa {
		return FromMap(m)
	}

	return FromSlice(any(doc).([]any))
}

// FromDocumentOfValue converts a map[string]Value or []Value into a Value.
// See FromMapOfValue and FromSliceOfValue.
func FromDocumentOfValue[T map[string]Value | []Value](doc T) Value {
	if m, isa := any(doc).(map[string]Value); isa {
		return FromMapOfValue(m)
	}

	return FromSliceOfValue(any(doc).([]Value))
}

// FromString converts a string into a Value
func FromString(s string) Value {
	return Value{typ: String, value: s}
}

// FromNumeric converts any constraint.Numeric type to a Value
// If the conversion fails, an Invalid Value is returned
func FromNumber(n any) Value {
	// The value can be any value conv.To accepts.
	// *big.Rat must be converted to a normalized string.
	var s NumberString

	if br, isa := n.(*big.Rat); isa {
		s = NumberString(conv.BigRatToNormalizedString(br))
	} else if str, isa := n.(string); isa && (strings.TrimSpace(str) == "") {
		return Value{}
	} else {
		if err := conv.To(n, &s); err != nil {
			return Value{}
		}
	}

	return Value{typ: Number, value: s}
}

// FromBool converts a bool into a Value
func FromBool(b bool) Value {
	return funcs.Ternary(b, TrueValue, FalseValue)
}

// Type returns the type of value this Value contains
func (jv Value) Type() ValueType {
	return jv.typ
}

// AsMap returns a map representation of a Value.
// Panics if the Value is not an object.
func (jv Value) AsMap() map[string]Value {
	if jv.typ != Object {
		panic(errNotObject)
	}

	return jv.value.(map[string]Value)
}

// AsSlice returns a slice representation of a Value.
// Panics if the Value is not an array.
func (jv Value) AsSlice() []Value {
	if jv.typ != Array {
		panic(errNotArray)
	}

	return jv.value.([]Value)
}

// AsString returns a string representation of a Value.
// Panics if the Value is not a string, number, or boolean.
func (jv Value) AsString() string {
	switch jv.typ {
	case String:
		return jv.value.(string)
	case Number:
		return string(jv.value.(NumberString))
	case Boolean:
		return fmt.Sprintf("%t", jv.value.(bool))
	}

	panic(errNotStringable)
}

// AsBigRat returns a NumberString representation of a Value.
// Panics if the Value is not a number.
func (jv Value) AsNumberString() NumberString {
	if jv.typ != Number {
		panic(errNotNumber)
	}

	return jv.value.(NumberString)
}

// AsBoolean returns a bool representation of a Value.
// Panics if the Value is not a boolean.
func (jv Value) AsBool() bool {
	if jv.typ != Boolean {
		panic(errNotBoolean)
	}

	return jv.value.(bool)
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
// NUll    = nil
func (jv Value) ToAny() any {
	switch jv.typ {
	case Object:
		return jv.ToMap()
	case Array:
		return jv.ToSlice()
	default:
		return jv.value
	}
}

// ToMap returns a map[string]any representation of a Value.
// Panics if the Value is not an object.
func (jv Value) ToMap() map[string]any {
	if jv.typ != Object {
		panic(errNotObject)
	}

	m := map[string]any{}

	for k, v := range jv.value.(map[string]Value) {
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

	for _, v := range jv.value.([]Value) {
		s = append(s, v.ToAny())
	}

	return s
}

// IsDocument returns true if a Value is a document (Object or Array)
func (jv Value) IsDocument() bool {
	return (jv.typ == Object) || (jv.typ == Array)
}
