package json

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math/big"

	"github.com/bantling/micro/go/funcs"
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

// Error constants
var (
	ErrNotObject     = fmt.Errorf("The JSONValue is not an object")
	ErrNotArray      = fmt.Errorf("The JSONValue is not an array")
	ErrNotString     = fmt.Errorf("The JSONValue is not a string")
	ErrNotNumber     = fmt.Errorf("The JSONValue is not a number")
	ErrNotBoolean    = fmt.Errorf("The JSONValue is not a boolean")
	ErrNotStringable = fmt.Errorf("The JSONValue is not a string, number, or boolean")
)

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
func (jv JSONValue) AsString() string {
	switch jv.typ {
	case String:
		return jv.value.(string)
	case Number:
		return jv.value.(*big.Float).String()
	case Boolean:
		return fmt.Sprintf("%t", jv.value.(bool))
	}

	panic(ErrNotStringable)
}

// AsNumber returns a number representation of a JSONValue.
// Panics if the JSONValue is not a number.
func (jv JSONValue) AsNumber() *big.Float {
	if jv.typ != Number {
		panic(ErrNotNumber)
	}

	return jv.value.(*big.Float)
}

// AsBoolean returns a boolean representation of a JSONValue.
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
// The purposse is to allow conversion of json values to arbitrary go values, eg convert all numbers to ints.
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

// DefaultVisitor is the default visitor function that converts JSON values as follows:
// Object and Array: return a recursive call to JSONValue.Visit(DefaultVisitor)
// String: string
// Number: *big.Float
// Boolean: bool
// Null: nil
func DefaultVisitor(jv JSONValue) any {
	if (jv.typ == Object) || (jv.typ == Array) {
		return jv.Visit(DefaultVisitor)
	}

	// Must be String, Number, Boolean, or Null, which is already string, *big.Float, bool, or nil
	return jv.value
}

// ConversionVisitor generates a visitor function, given conversion functions for String, Number,and Boolean values that
// converts as follows:
// Object and Array: return a recursive call to JSONValue.Visit(ConversionVisitor)
// String: stringConv(string)
// Number: numbeerConv(*big.Float)
// Boolean: boolConv(bool)
// Null: nil
//
// If any of the conversion funcs are nil, then no conversion is peformed, resulting in string, *big.Float, or bool
func ConversionVisitor(
	stringConv func(string) any,
	numberConv func(*big.Float) any,
	boolConv func(bool) any,
) func(JSONValue) any {
	var (
		stringConvFn        = funcs.Ternary(stringConv != nil, stringConv, func(val string) any { return val })
		numberConvFn        = funcs.Ternary(numberConv != nil, numberConv, func(val *big.Float) any { return val })
		boolConvFn          = funcs.Ternary(boolConv != nil, boolConv, func(val bool) any { return val })
		conversionVisitorFn func(JSONValue) any
	)

	conversionVisitorFn = func(jv JSONValue) any {
		switch jv.typ {
		case Object:
			return jv.Visit(conversionVisitorFn)
		case Array:
			return jv.Visit(conversionVisitorFn)
		case String:
			return stringConvFn(jv.value.(string))
		case Number:
			return numberConvFn(jv.value.(*big.Float))
		case Boolean:
			return boolConvFn(jv.value.(bool))
		}

		// Must be Null
		return nil
	}

	return conversionVisitorFn
}

// NumberToInt64Conversion is a ConversionVisitor numberConv func that converts Number to int64, typed as any.
func NumberToInt64Conversion(val *big.Float) any {
	i64, _ := val.Int64()
	return i64
}

// NumberToFloat64Conversion is a ConversionVisitor numberConv func that converts Number to float64, typed as any.
func NumberToFloat64Conversion(val *big.Float) any {
	f64, _ := val.Float64()
	return f64
}

// ToMap returns a map[string]any representation of a JSONValue.
// Panics if the JSONValue is not an object.
// If the optional visitor func is not provided, then DefaultVisitor is used.
func (jv JSONValue) ToMap(visitor ...func(JSONValue) any) map[string]any {
	if jv.typ != Object {
		panic(ErrNotObject)
	}

	return jv.Visit(funcs.SliceIndex(visitor, 0, DefaultVisitor)).(map[string]any)
}

// ToSlice returns a []any representation of a JSONValue.
// Panics if the JSONValue is not an Array.
// If the optional visitor func is not provided, then DefaultVisitor is used.
func (jv JSONValue) ToSlice(visitor ...func(JSONValue) any) []any {
	if jv.typ != Array {
		panic(ErrNotArray)
	}

	visitorFunc := DefaultVisitor
	if len(visitor) > 0 {
		visitorFunc = visitor[0]
	}

	return jv.Visit(visitorFunc).([]any)
}

// ToInt returns an int64 representation of a JSONValue.
// Panics if the JSONValue is not a number.
// If the number has a decimal portion, it is rounded.
func (jv JSONValue) ToInt() int64 {
	if jv.typ != Number {
		panic(ErrNotNumber)
	}

	res, _ := jv.value.(*big.Float).Int64()
	return res
}

// ToInt returns an float64 representation of a JSONValue.
// Panics if the JSONValue is not a number.
func (jv JSONValue) ToFloat() float64 {
	if jv.typ != Number {
		panic(ErrNotNumber)
	}

	res, _ := jv.value.(*big.Float).Float64()
	return res
}
