package json

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math/big"

	"github.com/bantling/micro/go/constraint"
	"github.com/bantling/micro/go/conv"
	"github.com/bantling/micro/go/funcs"
)

// Error constants
var (
	ErrInvalidGoValueMsg = "A value of type %T is not a valid type to convert to a Value. Acceptable types are map[string]any, []any, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *big.Int, *big.Float, *big.Rat, NumberString, bool, and nil"
	ErrNotObject         = fmt.Errorf("The Value is not an object")
	ErrNotArray          = fmt.Errorf("The Value is not an array")
	ErrNotString         = fmt.Errorf("The Value is not a string")
	ErrNotNumber         = fmt.Errorf("The Value is not a number")
	ErrNotBoolean        = fmt.Errorf("The Value is not a boolean")
	ErrNotStringable     = fmt.Errorf("The Value is not a string, number, or boolean")
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

// NumberString is a special type that allows a plain string to be considered a JSON Number
type NumberString string

// NumberType is a constraint of all possible number types
// int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *big.Int, *big.Float, *big.Rat, NumberString
type NumberType interface {
	constraint.Signed | constraint.UnsignedInteger | *big.Int | *big.Float | *big.Rat | NumberString
}

// Value represents any kind of JSON value - object, array, string, number, boolean, null
type Value struct {
	typ   ValueType
	value any
}

// Constant values for a true, false, and Null
var (
	TrueValue  = Value{typ: Boolean, value: true}
	FalseValue = Value{typ: Boolean, value: false}
	NullValue  = Value{typ: Null, value: nil}
)

// Constant for default value visitor
var DefaultConversionVisitor func(Value) any = ConversionVisitor(
	funcs.Passthrough[string],
	funcs.Passthrough[any],
	funcs.Passthrough[bool],
)

// fromNumberInternal converts any kind of number into a Value
// returns zero value if the given value is not any recognized numeric type
func fromNumberInternal(n any) Value {

	if v, isa := n.(int); isa {
		return FromSignedInt(v)
	} else if v, isa := n.(int8); isa {
		return FromSignedInt(v)
	} else if v, isa := n.(int16); isa {
		return FromSignedInt(v)
	} else if v, isa := n.(int32); isa {
		return FromSignedInt(v)
	} else if v, isa := n.(int64); isa {
		return FromSignedInt(v)
	} else if v, isa := n.(uint); isa {
		return FromUnsignedInt(v)
	} else if v, isa := n.(uint8); isa {
		return FromUnsignedInt(v)
	} else if v, isa := n.(uint16); isa {
		return FromUnsignedInt(v)
	} else if v, isa := n.(uint32); isa {
		return FromUnsignedInt(v)
	} else if v, isa := n.(uint64); isa {
		return FromUnsignedInt(v)
	} else if v, isa := n.(float32); isa {
		return FromFloat(v)
	} else if v, isa := n.(float64); isa {
		return FromFloat(v)
	} else if v, isa := n.(*big.Int); isa {
		return FromBigInt(v)
	} else if v, isa := n.(*big.Float); isa {
		return FromBigFloat(v)
	} else if v, isa := n.(*big.Rat); isa {
		return FromBigRat(v)
	} else if v, isa := n.(NumberString); isa {
		return FromNumberString(v)
	}

	return Value{}
}

// FromValue converts a Go value into a Value, where the Go value must be as follows:
//
// Object: map[string]any
// Array: []any
// String: string
// Number: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *big.Int, *big.Float, *big.Rat, or NumberString
// Boolean: bool
// Null: nil
//
// Panics if any other kind of Go value is provided
// If the optional custom number conversion is provided, it is only called if the value is NOT map[string]any, []any,
// string, bool, or nil.
func FromValue(v any, cnf ...func(any) Value) Value {
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
	} else if len(cnf) > 0 {
		return cnf[0](v)
	} else if jval = fromNumberInternal(v); jval.value == nil {
		panic(fmt.Errorf(ErrInvalidGoValueMsg, v))
	}

	return jval
}

// FromMap converts a map[string]any into a Value.
// The types of the map keys must be acceptable to FromValue.
// An optional custom number conversion can be provided, which is passed to FromValue.
func FromMap(m map[string]any, cnf ...func(any) Value) Value {
	jv := map[string]Value{}

	for k, v := range m {
		jv[k] = FromValue(v, cnf...)
	}

	return Value{typ: Object, value: jv}
}

// FromSlice converts a []any into a Value.
// The types of the slice elements must be acceptable to FromValue.
// An optional custom number conversion can be provided, which is passed to FromValue.
func FromSlice(a []any, cnf ...func(any) Value) Value {
	jv := make([]Value, len(a))

	for i, v := range a {
		jv[i] = FromValue(v, cnf...)
	}

	return Value{typ: Array, value: jv}
}

// FromString converts a string into a Value
func FromString(s string) Value {
	return Value{typ: String, value: s}
}

// FromSignedInt converts any kind of signed int into a Value
func FromSignedInt[T constraint.SignedInteger](n T) Value {
	return Value{typ: Number, value: conv.IntToBigRat(n)}
}

// FromUnsignedInt converts any kind of unsigned int into a Value
func FromUnsignedInt[T constraint.UnsignedInteger](n T) Value {
	return Value{typ: Number, value: conv.UintToBigRat(n)}
}

// FromFloat converts any kind of float into a Value
func FromFloat[T constraint.Float](n T) Value {
	return Value{typ: Number, value: conv.FloatToBigRat(n)}
}

// FromBigInt converts a *big.Int into a Value
func FromBigInt(n *big.Int) Value {
	return Value{typ: Number, value: conv.BigIntToBigRat(n)}
}

// FromBigFloat converts a *big.Float into a Value
func FromBigFloat(n *big.Float) Value {
	return Value{typ: Number, value: conv.BigFloatToBigRat(n)}
}

// FromBigRat converts a *big.Rat into a Value
func FromBigRat(n *big.Rat) Value {
	return Value{typ: Number, value: n}
}

// FromNumberString converts a NumberString into a Value
func FromNumberString(n NumberString) Value {
	// Convert to *big.Float first, to ensure only a floating point string is acceptable.
	// Then convert to *big.Rat, as that is the internal value for numbers.
	return Value{typ: Number, value: conv.StringToBigRat(string(n))}
}

// FromNumber converts an int, int8, int16, int32, int64, uint, uin8, uint16, uint32, uint64, float32, float64, *big.Int,
// *big.Float, *big.Rat, or NumberString into a Value
func FromNumber[N NumberType](n N) Value {
	return fromNumberInternal(n)
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
		panic(ErrNotObject)
	}

	return jv.value.(map[string]Value)
}

// AsSlice returns a slice representation of a Value.
// Panics if the Value is not an array.
func (jv Value) AsSlice() []Value {
	if jv.typ != Array {
		panic(ErrNotArray)
	}

	return jv.value.([]Value)
}

// AsString returns a string representation of a Value.
// Panics if the Value is not a string, number, or boolean.
// In the case of number, if it is an int, then it is formatted as an int, otherwise it is formatted as a float.
func (jv Value) AsString() string {
	switch jv.typ {
	case String:
		return jv.value.(string)
	case Number:
		return conv.BigRatToNormalizedString(jv.value.(*big.Rat))
	case Boolean:
		return fmt.Sprintf("%t", jv.value.(bool))
	}

	panic(ErrNotStringable)
}

// AsBigRat returns a *big.Rat representation of a Value.
// Panics if the Value is not a number.
func (jv Value) AsBigRat() *big.Rat {
	if jv.typ != Number {
		panic(ErrNotNumber)
	}

	return jv.value.(*big.Rat)
}

// AsBoolean returns a bool representation of a Value.
// Panics if the Value is not a boolean.
func (jv Value) AsBool() bool {
	if jv.typ != Boolean {
		panic(ErrNotBoolean)
	}

	return jv.value.(bool)
}

// IsNull returns true if the Value is a null, else false
func (jv Value) IsNull() bool {
	return jv.typ == Null
}

// Visit implements a very simple visitor pattern, where the provided visitor function is applied to each value in an
// object, each element in an array, and each string, number, boolean and null primitive value.
// It is up to the visitor func to recursively call the Value.Visit method on the values of each object key or array
// element. Empty objects and arrays are returned as non-nil empty maps and slices.
//
// The purpose is to allow conversion of json values to arbitrary go values, eg convert all numbers to ints.
func (jv Value) Visit(visitor func(Value) any) any {
	switch jv.typ {
	case Object:
		var (
			obj = jv.value.(map[string]Value)
			res = map[string]any{}
		)

		for k, v := range obj {
			res[k] = visitor(v)
		}

		return res

	case Array:
		var (
			obj = jv.value.([]Value)
			res = make([]any, len(obj))
		)

		for i, v := range obj {
			res[i] = visitor(v)
		}

		return res
	}

	// Must be a primitive value
	return visitor(jv)
}

// ConversionVisitor generates a visitor function - given conversion functions for String, Number,and Boolean values - that
// converts as follows:
// Object and Array: return a recursive call to Value.Visit(generated visitor function)
// String: stringConv(string)
// Number: numberConv(any)
// Boolean: boolConv(bool)
// Null: nil
//
// The conversion funcs cannot be nil, but can be funcs.Passthrough[string], funcs.Passthrough[any], and
// funcs.Passthrough[bool].
//
// The numberConv accepts any because it is possible to create a Value using a custom numeric conversion function rather
// than the default conversion to *big.Rat. (See FromValue, FromMap, and FromSlice functions).
//
// Unless you provide a custom numeric conversion function when creating your JSON values, you can assume the numberConv
// will receive a *big.Rat.
func ConversionVisitor[S, N, B any](
	stringConv func(string) S,
	numberConv func(any) N,
	boolConv func(bool) B,
) func(Value) any {
	var conversionVisitorFn func(Value) any

	conversionVisitorFn = func(jv Value) any {
		switch jv.typ {
		case Object:
			return jv.Visit(conversionVisitorFn)
		case Array:
			return jv.Visit(conversionVisitorFn)
		case String:
			return stringConv(jv.value.(string))
		case Number:
			return numberConv(jv.value)
		case Boolean:
			return boolConv(jv.value.(bool))
		}

		// Must be Null
		return nil
	}

	return conversionVisitorFn
}

// DefaultVisitorFunc is the default visitor function that converts JSON values as follows:
// Object and Array: return a recursive call to Value.Visit(DefaultVisitor)
// String: string
// Number: *big.Rat
// Boolean: bool
// Null: nil
func DefaultVisitorFunc(jv Value) any {
	if (jv.typ == Object) || (jv.typ == Array) {
		return jv.Visit(DefaultConversionVisitor)
	}

	// Must be String, Number, Boolean, or Null, which is already string, *big.Rat, bool, or nil
	return jv.value
}

// ToMap returns a map[string]any representation of a Value.
// Panics if the Value is not an object.
// If the optional visitor func is not provided, then DefaultVisitor is used.
func (jv Value) ToMap(visitor ...func(Value) any) map[string]any {
	if jv.typ != Object {
		panic(ErrNotObject)
	}

	return jv.Visit(funcs.SliceIndex(visitor, 0, DefaultVisitorFunc)).(map[string]any)
}

// ToSlice returns a []any representation of a Value.
// Panics if the Value is not an Array.
// If the optional visitor func is not provided, then DefaultVisitor is used.
func (jv Value) ToSlice(visitor ...func(Value) any) []any {
	if jv.typ != Array {
		panic(ErrNotArray)
	}

	return jv.Visit(funcs.SliceIndex(visitor, 0, DefaultVisitorFunc)).([]any)
}
