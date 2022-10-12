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
	ErrInvalidGoValueMsg       = "A value of type %T is not a valid type to convert to a JSONValue. Acceptable types are map[string]any, []any, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *big.Int, *big.Float, *big.Rat, NumberString, bool, and nil"
	ErrInvalidGoNumberValueMsg = "A value of type %T is not a valid type to convert to a JSON Number. Acceptable types are int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *big.Int, *big.Float, *big.Rat, and NumberString"
	ErrNotObject               = fmt.Errorf("The JSONValue is not an object")
	ErrNotArray                = fmt.Errorf("The JSONValue is not an array")
	ErrNotString               = fmt.Errorf("The JSONValue is not a string")
	ErrNotNumber               = fmt.Errorf("The JSONValue is not a number")
	ErrNotBoolean              = fmt.Errorf("The JSONValue is not a boolean")
	ErrNotStringable           = fmt.Errorf("The JSONValue is not a string, number, or boolean")
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

// JSONValue represents any kind of JSON value - object, array, string, number, boolean, null
type JSONValue struct {
	typ   ValueType
	value any
}

// Constant values for a true, false, and Null
var (
	TrueValue  = JSONValue{typ: Boolean, value: true}
	FalseValue = JSONValue{typ: Boolean, value: false}
	NullValue  = JSONValue{typ: Null, value: nil}
)

// Constant for default value visitor
var DefaultConversionVisitor func(JSONValue) any = ConversionVisitor(
	funcs.Passthrough[string],
	funcs.Passthrough[*big.Rat],
	funcs.Passthrough[bool],
)

// fromNumberInternal converts any kind of number into a JSONValue
// returns zero value if the given value is not any recognized numeric type
func fromNumberInternal(n any) JSONValue {
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

	return JSONValue{}
}

// FromValue converts a Go value into a JSONValue, where the Go value must be as follows:
//
// Object: map[string]any
// Array: []any
// String: string
// Number: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *big.Int, *big.Float, *big.Rat, or NumberString
// Boolean: bool
// Null: nil
//
// Panics if any other kind of Go value is provided
func FromValue(v any) JSONValue {
	var jval JSONValue

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
	} else if jval = fromNumberInternal(v); jval.value == nil {
		panic(fmt.Errorf(ErrInvalidGoValueMsg, v))
	}

	return jval
}

// FromMap converts a map[string]any into a JSONValue
// The types of the map keys must be acceptable to FromValue
func FromMap(m map[string]any) JSONValue {
	jv := map[string]JSONValue{}

	for k, v := range m {
		jv[k] = FromValue(v)
	}

	return JSONValue{typ: Object, value: jv}
}

// FromSlice converts a []any into a JSONValue
// The types of the slice elements must be acceptable to FromValue
func FromSlice(a []any) JSONValue {
	jv := make([]JSONValue, len(a))

	for i, v := range a {
		jv[i] = FromValue(v)
	}

	return JSONValue{typ: Array, value: jv}
}

// FromString converts a string into a JSONValue
func FromString(s string) JSONValue {
	return JSONValue{typ: String, value: s}
}

// FromSignedInt converts any kind of signed int into a JSONValue
func FromSignedInt[T constraint.SignedInteger](n T) JSONValue {
	return JSONValue{typ: Number, value: conv.IntToBigRat(n)}
}

// FromUnsignedInt converts any kind of unsigned int into a JSONValue
func FromUnsignedInt[T constraint.UnsignedInteger](n T) JSONValue {
	return JSONValue{typ: Number, value: conv.UintToBigRat(n)}
}

// FromFloat converts any kind of float into a JSONValue
func FromFloat[T constraint.Float](n T) JSONValue {
	return JSONValue{typ: Number, value: conv.FloatToBigRat(n)}
}

// FromBigInt converts a *big.Int into a JSONValue
func FromBigInt(n *big.Int) JSONValue {
	return JSONValue{typ: Number, value: conv.BigIntToBigRat(n)}
}

// FromBigFloat converts a *big.Float into a JSONValue
func FromBigFloat(n *big.Float) JSONValue {
	return JSONValue{typ: Number, value: conv.BigFloatToBigRat(n)}
}

// FromBigRat converts a *big.Rat into a JSONValue
func FromBigRat(n *big.Rat) JSONValue {
	return JSONValue{typ: Number, value: n}
}

// FromNumberString converts a NumberString into a JSONValue
func FromNumberString(n NumberString) JSONValue {
	// Convert to *big.Float first, to ensure only a floating point string is acceptable.
	// Then convert to *big.Rat, as that is the internal value for numbers.
	return JSONValue{typ: Number, value: conv.BigFloatToBigRat(conv.StringToBigFloat(string(n)))}
}

// FromNumber converts an int, int8, int16, int32, int64, uint, uin8, uint16, uint32, uint64, float32, float64, *big.Int,
// *big.Float, *big.Rat, or NumberString into a JSONValue
func FromNumber(n any) JSONValue {
	jv := fromNumberInternal(n)
	if jv.value == nil {
		panic(fmt.Errorf(ErrInvalidGoNumberValueMsg, n))
	}
	return jv
}

// FromBool converts a bool into a JSONValue
func FromBool(b bool) JSONValue {
	return funcs.Ternary(b, TrueValue, FalseValue)
}

// Type returns the type of value this JSONValue contains
func (jv JSONValue) Type() ValueType {
	return jv.typ
}

// AsMap returns a map representation of a JSONValue.
// Panics if the JSONValue is not an object.
func (jv JSONValue) AsMap() map[string]JSONValue {
	if jv.typ != Object {
		panic(ErrNotObject)
	}

	return jv.value.(map[string]JSONValue)
}

// AsSlice returns a slice representation of a JSONValue.
// Panics if the JSONValue is not an array.
func (jv JSONValue) AsSlice() []JSONValue {
	if jv.typ != Array {
		panic(ErrNotArray)
	}

	return jv.value.([]JSONValue)
}

// AsString returns a string representation of a JSONValue.
// Panics if the JSONValue is not a string, number, or boolean.
// In the case of number, if it is an int, then it is formatted as an int, otherwise it is formatted as a float.
func (jv JSONValue) AsString() string {
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

// AsBigRat returns a *big.Rat representation of a JSONValue.
// Panics if the JSONValue is not a number.
func (jv JSONValue) AsBigRat() *big.Rat {
	if jv.typ != Number {
		panic(ErrNotNumber)
	}

	return jv.value.(*big.Rat)
}

// AsBoolean returns a bool representation of a JSONValue.
// Panics if the JSONValue is not a boolean.
func (jv JSONValue) AsBoolean() bool {
	if jv.typ != Boolean {
		panic(ErrNotBoolean)
	}

	return jv.value.(bool)
}

// IsNull returns true if the JSONValue is a null, else false
func (jv JSONValue) IsNull() bool {
	return jv.typ == Null
}

// Visit implements a very simple visitor pattern, where the provided visitor function is applied to each value in an
// object, each element in an array, and each string, number, boolean and null primitive value.
// It is up to the visitor func to recursively call the JSONValue.Visit method on the values of each object key or array
// element. Empty objects and arrays are returned as non-nil empty maps and slices.
//
// The purpose is to allow conversion of json values to arbitrary go values, eg convert all numbers to ints.
func (jv JSONValue) Visit(visitor func(JSONValue) any) any {
	switch jv.typ {
	case Object:
		var (
			obj = jv.value.(map[string]JSONValue)
			res = map[string]any{}
		)

		for k, v := range obj {
			res[k] = visitor(v)
		}

		return res

	case Array:
		var (
			obj = jv.value.([]JSONValue)
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
// Object and Array: return a recursive call to JSONValue.Visit(generated visitor function)
// String: stringConv(string)
// Number: numberConv(*big.Rat)
// Boolean: boolConv(bool)
// Null: nil
//
// The conversion funcs cannot be nil, but can be funcs.Passthrough[string], funcs.Passthrough[*big.Rat], and
// funcs.Passthrough[bool].
func ConversionVisitor[S, N, B any](
	stringConv func(string) S,
	numberConv func(*big.Rat) N,
	boolConv func(bool) B,
) func(JSONValue) any {
	var conversionVisitorFn func(JSONValue) any

	conversionVisitorFn = func(jv JSONValue) any {
		switch jv.typ {
		case Object:
			return jv.Visit(conversionVisitorFn)
		case Array:
			return jv.Visit(conversionVisitorFn)
		case String:
			return stringConv(jv.value.(string))
		case Number:
			return numberConv(jv.value.(*big.Rat))
		case Boolean:
			return boolConv(jv.value.(bool))
		}

		// Must be Null
		return nil
	}

	return conversionVisitorFn
}

// DefaultVisitorFunc is the default visitor function that converts JSON values as follows:
// Object and Array: return a recursive call to JSONValue.Visit(DefaultVisitor)
// String: string
// Number: *big.Rat
// Boolean: bool
// Null: nil
func DefaultVisitorFunc(jv JSONValue) any {
	if (jv.typ == Object) || (jv.typ == Array) {
		return jv.Visit(DefaultConversionVisitor)
	}

	// Must be String, Number, Boolean, or Null, which is already string, *big.Rat, bool, or nil
	return jv.value
}

// ToMap returns a map[string]any representation of a JSONValue.
// Panics if the JSONValue is not an object.
// If the optional visitor func is not provided, then DefaultVisitor is used.
func (jv JSONValue) ToMap(visitor ...func(JSONValue) any) map[string]any {
	if jv.typ != Object {
		panic(ErrNotObject)
	}

	return jv.Visit(funcs.SliceIndex(visitor, 0, DefaultVisitorFunc)).(map[string]any)
}

// ToSlice returns a []any representation of a JSONValue.
// Panics if the JSONValue is not an Array.
// If the optional visitor func is not provided, then DefaultVisitor is used.
func (jv JSONValue) ToSlice(visitor ...func(JSONValue) any) []any {
	if jv.typ != Array {
		panic(ErrNotArray)
	}

	return jv.Visit(funcs.SliceIndex(visitor, 0, DefaultVisitorFunc)).([]any)
}
