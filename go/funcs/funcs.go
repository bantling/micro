// Package funcs is useful Go functions
// SPDX-License-Identifier: Apache-2.0
package funcs

import (
	"fmt"
	"math/cmplx"
	"reflect"
	"sort"

	"github.com/bantling/micro/go/constraint"
)

const (
	notNilableMsg     = "Type %s is not a nillable type"
	unsortableTypeMsg = "The type %T is not a sortable type - it must implement constraint.Ordered, constraint.Complex, or constraint.Cmp. Provide a custom sorting function."
)

// ==== Element accessors

// SliceIndex returns the first of the following given a slice, index, and optional default value:
// 1. slice[index] if the slice is non-nil and length > index
// 2. default value if provided
// 3. zero value of slice element type
func SliceIndex[T any](slc []T, index uint, defawlt ...T) T {
	// Return index if it exists
	idx := int(index)
	if (slc != nil) && (len(slc) > idx) {
		return slc[idx]
	}

	// Else return default if provided
	if len(defawlt) > 0 {
		return defawlt[0]
	}

	// Else return zero value
	var zv T
	return zv
}

// MapValue returns the first of the following:
// 1. map[key] if the map is non-nil and the key exists in the map
// 2. default if provided
// 3. zero value of map value type
func MapValue[K comparable, V any](mp map[K]V, key K, defawlt ...V) V {
	// Return key value if it exists
	if mp != nil {
		if val, haveIt := mp[key]; haveIt {
			return val
		}
	}

	// Else return default if provided
	if len(defawlt) > 0 {
		return defawlt[0]
	}

	// Else return zero value of map value type
	var zv V
	return zv
}

// ==== Filters

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

// Not (filter func) adapts a filter func (func(T) bool) to the negation of the func.
func Not[T any](filter func(T) bool) func(T) bool {
	return func(t T) bool {
		return !filter(t)
	}
}

// ==== Ternary

// Ternary returns trueVal if expr is true, else it returns falseVal
func Ternary[T any](expr bool, trueVal T, falseVal T) T {
	if expr {
		return trueVal
	}

	return falseVal
}

// TernaryResult returns trueVal() if expr is true, else it returns falseVal()
func TernaryResult[T any](expr bool, trueVal func() T, falseVal func() T) T {
	if expr {
		return trueVal()
	}

	return falseVal()
}

// ==== Comparators

// LessThan returns a filter func (func(T) bool) that returns true if it accepts a value that is less than the given value
func LessThan[T constraint.Ordered](val T) func(T) bool {
	return func(t T) bool {
		return t < val
	}
}

// LessThanEqual returns a filter func (func(T) bool) that returns true if it accepts a value that is less than or equal to the given value
func LessThanEqual[T constraint.Ordered](val T) func(T) bool {
	return func(t T) bool {
		return t <= val
	}
}

// Equal returns a filter func (func(T) bool) that returns true if it accepts a value that equals the given value with ==
func Equal[T comparable](val T) func(T) bool {
	return func(t T) bool {
		return t == val
	}
}

// GreaterThan returns a filter func (func(T) bool) that returns true if it accepts a value that is greater than the given value
func GreaterThan[T constraint.Ordered](val T) func(T) bool {
	return func(t T) bool {
		return t > val
	}
}

// GreaterThanEqual returns a filter func (func(T) bool) that returns true if it accepts a value that is greater than or equal to the given value
func GreaterThanEqual[T constraint.Ordered](val T) func(T) bool {
	return func(t T) bool {
		return t >= val
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

// ==== Sort

// Sort sorts a slice of Ordered
func SortOrdered[T constraint.Ordered](slc []T) {
	sort.Slice(slc, func(i, j int) bool { return slc[i] < slc[j] })
}

// SortComplex sorts a slice of Complex
func SortComplex[T constraint.Complex](slc []T) {
	sort.Slice(slc, func(i, j int) bool { return cmplx.Abs(complex128(slc[i])) < cmplx.Abs(complex128(slc[j])) })
}

// SortCmp sorts a slice of Cmp
func SortCmp[T constraint.Cmp[T]](slc []T) {
	sort.Slice(slc, func(i, j int) bool { return slc[i].Cmp(slc[j]) < 0 })
}

// SortBy sorts a slice of any type with the provided comparator
func SortBy[T any](slc []T, less func(T, T) bool) {
	sort.Slice(slc, func(i, j int) bool { return less(slc[i], slc[j]) })
}

// ==== Nil

// Nillable returns true if the given reflect.Type represents a chan, func, map, pointer, or slice.
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

// IsNil generates a func that returns true if the value given is nil.
// A type constraint cannot be used to describe nillable types at compile time, so reflection is used.
func IsNil[T any]() func(T) bool {
	var n T
	MustBeNillable(reflect.TypeOf(n))

	return func(t T) bool {
		return reflect.ValueOf(t).IsNil()
	}
}

// IsNonNil returns a func(T) bool that returns true if it accepts a non-nil value.
// A type constraint cannot be used to describe nillable types at compile time, so reflection is used.
func IsNonNil[T any]() func(T) bool {
	var n T
	MustBeNillable(reflect.TypeOf(n))

	return func(t T) bool {
		return !reflect.ValueOf(t).IsNil()
	}
}

// ==== Error

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

// ==== Supplier

// SupplierOf returns a func() T that returns the given value
func SupplierOf[T any](value T) func() T {
	return func() T {
		return value
	}
}

// CachingSupplier returns a func() T that caches the result of the given supplier on the first call.
// Any subseqquent calls return the cached value, guaranteeing the provided supplier is invoked at most once.
func CachingSupplier[T any](supplier func() T) func() T {
	var (
		isCached  bool
		cachedVal T
	)

	return func() T {
		if !isCached {
			isCached, cachedVal = true, supplier()
		}

		return cachedVal
	}
}

// IgnoreResult takes a func of no args that returns any type, and converts it to a func of no args and no return value.
// Useful for TryTo function closers.
func IgnoreResult[T any](fn func() T) func() {
	return func() {
		fn()
	}
}

// TryTo executes tryFn, and if a panic occurs, it executes panicFn.
// If any closers are provided, they are deferred in the provided order before the tryFn, to ensure they get closed even if a panic occurs.
// If any closer returns a non-nil error, any remaining closers are still called, as that is go built in behaviour.
//
// This function simplifies the process of "catching" panics over using reverse order code like the following
// (common in unit tests that want to verify the type of object sent to panic):
// func DoSomeStuff() {
//   ...
//   func() {
//     defer zero or more things that have to be closed before we try to recover from any panic
//     defer func() {
//       // Some code that uses recover() to try and deal with a panic
//     }()
//     // Some code that may panic, which is handled by above code
//   }
//   ...
// }
func TryTo(tryFn func(), panicFn func(any), closers ...func()) {
	// Defer code that attempts to recover a value - first func deferred is called last, so this func is called after all provided closers
	defer func() {
		if val := recover(); val != nil {
			panicFn(val)
		}
	}()

	// Defer all closers in provided order, so they get called in reverse order as expected
	for _, closerFn := range closers {
		defer closerFn()
	}

	// Execute code that may panic, which is supposed to panic with a value of type error
	tryFn()
}
