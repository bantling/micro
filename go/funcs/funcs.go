// Package funcs is useful Go functions
package funcs

import (
	"fmt"
	"reflect"

	"github.com/bantling/micro/go/constraint"
)

const (
	notNilableMsg = "Type %s is not a nillable type"
)

// And converts any number of filter funcs (func(T) bool) into the conjunction of all the funcs.
// Short-circuit logic will return false on the first function that returns false.
// If no filters are provided, the result is a function that always returns true.
func And[T any](filters ...func(T) bool) func(T) bool {
	return func(t T) bool {
		result := true

		for _, nextFilter := range filters {
			if result = nextFilter(t); !result {
				break
			}
		}

		return result
	}
}

// Or converts any number of filter funcs (func(T) bool) into the disjunction of all the funcs.
// Short-circuit logic will return true on the first function that returns true.
// If no filters are provided, the result is a function that always returns true.
func Or[T any](filters ...func(T) bool) func(T) bool {
	return func(t T) bool {
		result := true

		for _, nextFilter := range filters {
			if result = nextFilter(t); result {
				break
			}
		}

		return result
	}
}

// Not (filter func) adapts a filter func(any) bool to the negation of the func.
func Not[T any](filter func(T) bool) func(T) bool {
	return func(t T) bool {
		return !filter(t)
	}
}

// EqualTo returns a func(T) bool that returns true if it accepts a value that equals the given value
func EqualTo[T comparable](val T) func(T) bool {
	return func(t T) bool {
		return t == val
	}
}

// LessThan returns a func(T) bool that returns true if it accepts a value that is less than the given value
func LessThan[T constraint.Ordered](val T) func(T) bool {
	return func(t T) bool {
		return t < val
	}
}

// LessThanEquals returns a func(T) bool that returns true if it accepts a value that is less than or equal to the given value
func LessThanEqual[T constraint.Ordered](val T) func(T) bool {
	return func(t T) bool {
		return t <= val
	}
}

// IsNegative returns a func(T) bool that returns true if it accepts a negative value.
func IsNegative[T constraint.Signed]() func(T) bool {
	return func(t T) bool {
		return t < 0
	}
}

// IsNonNegative returns a func(T) bool that returns true if it accepts a non-negative value.
func IsNonNegative[T constraint.Signed]() func(T) bool {
	return func(t T) bool {
		return t >= 0
	}
}

// IsPositive returns a func(T) bool that returns true if it accepts a positive value.
func IsPositive[T constraint.Signed]() func(T) bool {
	return func(t T) bool {
		return t > 0
	}
}

// Nillable returns true if the given reflect.Type represents a chan, func, map, pointer, and slice
func Nillable(typ reflect.Type) bool {
	nillable := true

	switch typ.Kind() {
	case reflect.Chan:
	case reflect.Func:
	case reflect.Map:
	case reflect.Pointer:
	case reflect.Slice:
	default:
		nillable = false
	}

	return nillable
}

// MustBeNillable panics if Nillable(typ) returns false
func MustBeNillable(typ reflect.Type) {
	if !Nillable(typ) {
		panic(fmt.Errorf(notNilableMsg, typ.Name()))
	}
}

// IsNil returns true if the value given is nil.
// A type constraint cannot be used to describe nillable types at compile time, so reflection is used.
// Nillable types are .
func IsNil[T any]() func(T) bool {
	var n T
	MustBeNillable(reflect.TypeOf(n))

	return func(t T) bool {
		return reflect.ValueOf(t).IsNil()
	}
}

// IsNonNil returns a func(T) bool that returns true if it accepts a non-nil value.
func IsNonNil[T any]() func(T) bool {
	var n T
	MustBeNillable(reflect.TypeOf(n))

	return func(t T) bool {
		return !reflect.ValueOf(t).IsNil()
	}
}

// Supplier returns a func() T that returns the given value
func Supplier[T any](value T) func() T {
	return func() T {
		return value
	}
}

// Must panics if the error is non-nil, else returns
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// MustValue panics if the error is non-nil, else returns the value of type T
func MustValue[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}
