package stream

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	gomath "math"
	"math/bits"
	"reflect"
	"sync"

	"github.com/bantling/micro/constraint"
	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/iter"
	"github.com/bantling/micro/math"
	"github.com/bantling/micro/tuple"
)

// PUnit indicates how to interpret a parallel quantity
type PUnit bool

const (
	Threads PUnit = false // NumThreads indicates the quantity is the number of threads
	Items   PUnit = true  // NumItems indicates the quantity is the number of items each thread processes
)

// PInfo includes the number of items and a unit
type PInfo struct {
	N int
	PUnit
}

// Constants
var (
	absErrMsg = "Absolute value error for %d: there is no corresponding positive value in type %T"
)

// ==== Functions that provide the foundation for all other functions

// Map constructs a new Iter[U] from an Iter[T] and a func that transforms a T to a U.
//
// The resulting iter can return any kind of error from source iter, or EOI.
func Map[T, U any](mapper func(T) U) func(iter.Iter[T]) iter.Iter[U] {
	return func(it iter.Iter[T]) iter.Iter[U] {
		return iter.NewIter(func() (U, error) {
			val, err := it.Next()
			if err == nil {
				return mapper(val), nil
			}

			var zv U
			return zv, err
		})
	}
}

// MapError is similar to Map, except the mapper function returns (U, error), and the first element that returns a non-nil
// error results in iteration being cut short.
//
// The resulting iter can return any kind of error from source iter, or EOI.
func MapError[T, U any](mapper func(T) (U, error)) func(iter.Iter[T]) iter.Iter[U] {
	return func(it iter.Iter[T]) iter.Iter[U] {
		return iter.NewIter(func() (U, error) {
			val, err := it.Next()
			if err == nil {
				var mval U
				if mval, err = mapper(val); err == nil {
					return mval, nil
				}
			}

			var zv U
			return zv, err
		})
	}
}

// Filter constructs a new Iter[T] from an Iter[T] and a func that returns true if a T passes the filter.
//
// The resulting iter can return any kind of error from source iter, or EOI.
func Filter[T any](filter func(T) bool) func(iter.Iter[T]) iter.Iter[T] {
	return func(it iter.Iter[T]) iter.Iter[T] {
		return iter.NewIter(func() (T, error) {
			var (
				val T
				err error
			)
			for {
				val, err = it.Next()
				if err == nil {
					if filter(val) {
						return val, nil
					}
				} else {
					break
				}
			}

			var zv T
			return zv, err
		})
	}
}

// Reduce reduces all elements in the input set to an empty or single element output set, depending on two factors:
// - the number of elements in the input set (0, 1, or multiple)
// - whether or not the optional identity is provided
//
// If the optional identity is NOT provided:
// - 0 elements: empty
// - 1 element: the element
// - multiple elements: reducer(reducer(reducer(first, second), third), ...)
//
//	If the optional identity IS provided:
//
// - 0 elements: identity
// - 1 element: reducer(identity, the element)
// - multiple elements: reducer(reducer(reducer(identity, first), second), ...)
//
// The resulting iter can return any kind of error from source iter, or EOI.
func Reduce[T any](
	reducer func(T, T) T,
	identity ...T,
) func(iter.Iter[T]) iter.Iter[T] {
	return func(it iter.Iter[T]) iter.Iter[T] {
		var done bool

		return iter.NewIter(func() (T, error) {
			var (
				val T
				zv  T
				err error
			)

			if done {
				return zv, iter.EOI
			}

			done = true

			if len(identity) == 0 {
				// No identity
				val, err = it.Next()
				if err != nil {
					// 0 elements = Empty set, return error that may be EOI or a problem
					return zv, err
				} else {
					// At least one element, start with first element
					result := val

					for {
						val, err = it.Next()
						if err == nil {
							// If there are more elements, make cumulative reducer calls
							result = reducer(result, val)
						} else if err != iter.EOI {
							// A problem occurred, toss result
							return zv, err
						} else {
							// An EOI occurred
							break
						}
					}

					// Successfully return result
					return result, nil
				}
			}

			// There is an identity
			identityVal := identity[0]
			val, err = it.Next()
			if err != nil {
				if err == iter.EOI {
					// 0 elements = identity
					return identityVal, nil
				}
				// A problem
				return zv, err
			} else {
				// At least one element, call reducer with identity and first element
				result := reducer(identityVal, val)

				for {
					val, err = it.Next()
					if err == nil {
						// If there are more elements, make cumulative reducer calls
						result = reducer(result, val)
					} else if err != iter.EOI {
						// A problem occurred, toss result
						return zv, err
					} else {
						// An EOI occurred
						break
					}
				}

				// Successfully return result
				return result, nil
			}
		})
	}
}

// ReduceTo is similar to Reduce, except that:
// - The result does not have to be the same type
// - If no identity value is given, the zero value is used
//
// The resulting iter can return any kind of error from source iter, or EOI.
func ReduceTo[T, U any](
	reducer func(U, T) U,
	identity ...U,
) func(iter.Iter[T]) iter.Iter[U] {
	return func(it iter.Iter[T]) iter.Iter[U] {
		var done bool

		return iter.NewIter(func() (U, error) {
			var (
				val T
				zv  U
				err error
			)

			if done {
				return zv, iter.EOI
			}

			done = true

			if len(identity) == 0 {
				// No identity
				val, err = it.Next()
				if err != nil {
					// 0 elements = Empty set, return error that may be EOI or a problem
					return zv, err
				} else {
					// At least one element, call reducer with zero value and first element
					result := reducer(zv, val)

					for {
						val, err = it.Next()
						if err == nil {
							// If there are more elements, make cumulative reducer calls, combining old and new results
							result = reducer(result, val)
						} else if err != iter.EOI {
							// A problem occurred, toss result
							return zv, err
						} else {
							// An EOI occurred
							break
						}
					}

					// Successfully return result
					return result, nil
				}
			}

			// There is an identity
			identityVal := identity[0]
			val, err = it.Next()
			if err != nil {
				if err == iter.EOI {
					// 0 elements = identity
					return identityVal, nil
				}
				// A problem
				return zv, err
			} else {
				// At least one element, call reducer with identity and first element
				result := reducer(identityVal, val)

				for {
					val, err = it.Next()
					if err == nil {
						// If there are more elements, make cumulative reducer calls
						result = reducer(result, val)
					} else if err != iter.EOI {
						// A problem occurred, toss result
						return zv, err
					} else {
						// An EOI occurred
						break
					}
				}

				// Successfully return result
				return result, nil
			}
		})
	}
}

// ReduceToBool is similar to ReduceTo, except that it uses boolean short circuit logic to stop iterating early if
// possible. If stopVal is true, then early termination occurs on the first call to reducer that returns true, else it
// occurs on the first call to reducer that returns false.
func ReduceToBool[T any](
	predicate func(T) bool,
	identity bool,
	stopVal bool,
) func(iter.Iter[T]) iter.Iter[bool] {
	return func(it iter.Iter[T]) iter.Iter[bool] {
		var done bool

		return iter.NewIter(func() (bool, error) {
			if done {
				return false, iter.EOI
			}

			done = true

			var (
				val T
				err error
			)

			val, err = it.Next()
			if err != nil {
				if err == iter.EOI {
					// 0 elements = identity
					return identity, nil
				}
				// A problem occurred
				return false, err
			} else {
				// At least one element, call reducer with identity and first element
				result := predicate(val)
				if result == stopVal {
					// Stop early if result matches stopping condition
					return result, nil
				}

				for {
					val, err = it.Next()
					if err != nil {
						if err == iter.EOI {
							// Successfully return result
							break
						}
						// A problem
						return false, err
					}

					// If there are more elements, call predicate
					result = predicate(val)
					if result == stopVal {
						// Stop early if result matches stopping condition
						return result, nil
					}
				}

				// Successfully return result
				return result, nil
			}
		})
	}
}

// ReduceToSlice reduces an Iter[T] into a Iter[[]T] that contains a single element of type []T.
// Eg, an Iter[int] of 1,2,3,4,5 becomes an Iter[[]int] of [1,2,3,4,5].
// An empty Iter is reduced to a zero length slice.
func ReduceToSlice[T any](it iter.Iter[T]) iter.Iter[[]T] {
	var done bool

	return iter.NewIter(func() ([]T, error) {
		if done {
			return nil, iter.EOI
		}

		done = true

		var (
			val T
			err error
			slc = []T{}
		)

		for {
			val, err = it.Next()
			if err != nil {
				if err == iter.EOI {
					// Successfully iterated all values
					break
				}
				// A problem
				var zv []T
				return zv, err
			}

			// Append element
			slc = append(slc, val)
		}

		// Successfully return result
		return slc, nil
	})
}

// ReduceIntoSlice is the same as ReduceToSlice, except that:
// - It accepts a target slice to append results to
// - It generates a transform
//
// The generated transform panics if the target slice length is not at least as many elements as the source iter.
// If the underlying iter returns a non-nil non-EOI error, the provided slice will have zero values.
func ReduceIntoSlice[T any](slc []T) func(iter.Iter[T]) iter.Iter[[]T] {
	return func(it iter.Iter[T]) iter.Iter[[]T] {
		var (
			done bool
			i    int
		)

		return iter.NewIter(func() ([]T, error) {
			if done {
				return nil, iter.EOI
			}

			done = true

			var (
				val T
				err error
			)

			for {
				val, err = it.Next()
				if err != nil {
					if err == iter.EOI {
						// Successfully iterated all values
						break
					}
					// A problem
					var zv T
					for i := range slc {
						slc[i] = zv
					}
					return slc, err
				}

				// Set next slice index
				slc[i] = val
				i++
			}

			// Successfully return result
			return slc, nil
		})
	}
}

// ExpandSlices is the opposite of ReduceToSlice: an Iter[[]int] of [1,2,3,4,5] becomes an Iter[int] of 1,2,3,4,5.
// If the source Iter contains multiple slices, they are combined together into one set of data (skipping nil and empty
// slices), so that an Iter[[]int] of [1,2,3], nil, [], [4,5] also becomes an Iter[int] of 1,2,3,4,5.
// An empty Iter or an Iter with nil/empty slices is expanded to an empty Iter.
func ExpandSlices[T any](it iter.Iter[[]T]) iter.Iter[T] {
	var (
		slc []T
		idx int
	)

	return iter.NewIter(func() (T, error) {
		if (slc == nil) || (idx == len(slc)) {
			// Search for next non-nil non-empty slice
			// Nilify slc var in case we just finished iterating last element of last slice, which is non-nil
			slc = nil

			var (
				zv  T
				err error
			)

			for {
				slc, err = it.Next()
				if err != nil {
					if err == iter.EOI {
						// Successfully iterated all values - slc shd be nil, but make sure
						slc = nil
						break
					}
					// A problem
					return zv, err
				}

				if (slc != nil) && (len(slc) > 0) {
					// Found a non-empty slice
					idx = 0
					break
				}
			}

			// Stop if no more non-nil non-empty slices available
			if slc == nil {
				return zv, iter.EOI
			}
		}

		// Successfully acquired an index of a slice to return
		val := slc[idx]
		idx++

		return val, nil
	})
}

// ReduceToMap reduces an Iter[tuple.Two[K, V]]] into a Iter[map[K, V]] that contains a single element if type map[K, V].
// Eg, an Iter[tuple.Two[int]string] of {1: "1"}, {2: "2"} becomes an Iter[map[int]string] of {1: "1", 2: "2"}.
// If multiple tuple.Two objects in the Iter have the same key, the last such object in iteration order determines the
// value for the key in the resulting map.
// An empty Iter is reduced to an empty map.
func ReduceToMap[K comparable, V any](it iter.Iter[tuple.Two[K, V]]) iter.Iter[map[K]V] {
	var done bool

	return iter.NewIter(func() (map[K]V, error) {
		if done {
			return nil, iter.EOI
		}

		done = true

		var (
			m   = map[K]V{}
			kv  tuple.Two[K, V]
			err error
		)

		for {
			kv, err = it.Next()
			if err != nil {
				if err == iter.EOI {
					// Successfully iterated all values
					break
				}
				// A problem
				var zv map[K]V
				return zv, err
			}

			// Successfully acquired a key value pair to put into the map
			m[kv.T] = kv.U
		}

		return m, nil
	})
}

// ExpandMaps is the opposite of ReduceToMap: an Iter[map[int]string] of {1: "1", 2: "2", 3: "3"} becomes an
// Iter[tuple.Two[int, string]] of {1: "1"}, {2: "2"}, {3: "3"].
// If the source Iter contains multiple maps, they are combined together into one set of data (skipping nils),
// so that an Iter[map[int]string] of {1: "1", 2: "2"}, nil, {}, {3: "3"} also becomes
// an Iter[tuple.Two[int, string]] of {1: "1"}, {2: "2"}, {3: "3"}.
// An empty Iter or an Iter with nil/empty maps is expanded to an empty Iter.
func ExpandMaps[K comparable, V any](it iter.Iter[map[K]V]) iter.Iter[tuple.Two[K, V]] {
	var (
		m  map[K]V
		mr *reflect.MapIter
	)

	return iter.NewIter(func() (tuple.Two[K, V], error) {
		var (
			zv  tuple.Two[K, V]
			err error
		)

		if (m == nil) || (!mr.Next()) {
			// Search for next non-nil non-empty map
			// Nilify m var in case last call finished iterating last element of last map
			m = nil

			for {
				m, err = it.Next()
				if err != nil {
					if err == iter.EOI {
						// Unable to find next result, nilify m
						m = nil
						break
					}
					// A problem
					return zv, err
				}

				if m != nil {
					// Found non-nil map, see if it is also non-empty
					if mr = reflect.ValueOf(m).MapRange(); mr.Next() {
						break
					}
				}
			}

			// Stop if no more non-nil non-empty maps available
			if m == nil {
				return zv, iter.EOI
			}
		}

		val := tuple.Of2(mr.Key().Interface().(K), mr.Value().Interface().(V))

		return val, nil
	})
}

// Skip skips the first n elements, then iteration continues from there.
// If there are n or fewer elements in total, then the resulting iter is empty.
//
// Note that for the set 1,2,3,4,5 the composition os Skip(1),Limit(3) will first skip 1 then limit to 2,3,4;
// whereas the composition of Limit(3),Skip(1) will first limit to 1,2,3 then skip 1 returning 2,3.
func Skip[T any](n uint) func(iter.Iter[T]) iter.Iter[T] {
	return func(it iter.Iter[T]) iter.Iter[T] {
		skip := n

		return iter.NewIter(func() (T, error) {
			var (
				val T
				zv  T
				err error
			)

			// Skip first n values only once
			for skip > 0 {
				val, err = it.Next()
				if err != nil {
					// Reached end or problem
					skip = 0
					return zv, err
				}

				// Read a value to skip
				skip--
			}

			val, err = it.Next()
			if err != nil {
				// Reached end or problem
				return zv, err
			}

			// Successfully read a value to return
			return val, nil
		})
	}
}

// Limit returns the first n elements, then iteration stops and all further elements are ignored.
// If there fewer than n elements in total, then all n elements are returned.
//
// Note that for the set 1,2,3,4,5 the composition of Skip(1),Limit(3) will first skip 1 then limit to 2,3,4;
// whereas the composition of Limit(3),Skip(1) will first limit to 1,2,3 then skip 1 returning 2,3.
func Limit[T any](n uint) func(iter.Iter[T]) iter.Iter[T] {
	return func(it iter.Iter[T]) iter.Iter[T] {
		limit := n

		return iter.NewIter(func() (T, error) {
			var (
				val T
				zv  T
				err error
			)

			if limit > 0 {
				// Try to get next value that is within the limit
				val, err = it.Next()
				if err != nil {
					// EOI or problem
					return zv, err
				}

				// Successfully read value, decrement limit for next time
				limit--
				return val, nil
			}

			// Limit = 0, do not read any more values
			return zv, iter.EOI
		})
	}
}

// Peek executes a func for every item being iterated, which is a side effect.
func Peek[T any](fn func(T)) func(iter.Iter[T]) iter.Iter[T] {
	return func(it iter.Iter[T]) iter.Iter[T] {
		return iter.NewIter(func() (T, error) {
			// Read next value
			val, err := it.Next()

			if err != nil {
				// EOI or problem
				var zv T
				return zv, err
			}

			// Successfully found a value
			fn(val)
			return val, nil
		})
	}
}

// Generator receives a generator (a func of no args that returns a func of Iter[T] -> Iter[U], and detects if the
// Iter[T] has changed address. If so, it internally generates a new function by invoking the generator.
//
// This allows stateful transforms of Iter[T] -> Iter[U] that track state across calls to begin with a new initial state
// for each data set the transform is applied to.
//
// Generator is not thread safe, so be careful about storing a composition containing a Generator in a global variable:
// 1. Declare the global variable as a function of no args that generates the composition when executed
// 2. Declare composition in a local variable so each thread makes its own copy
// 3. Store composition in a Context that is visible across methods in the thread
//
// See Distinct for an example of a stateful function that uses Generator internally.
func Generator[T, U any](gen func() func(iter.Iter[T]) iter.Iter[U]) func(iter.Iter[T]) iter.Iter[U] {
	var (
		currentIter iter.Iter[T]
		fn          func(iter.Iter[T]) iter.Iter[U]
	)

	return func(it iter.Iter[T]) iter.Iter[U] {
		if currentIter != it {
			fn = gen()
		}

		return fn(it)
	}
}

// ==== Functions based on foundational functions

// AllMatch reduces Iter[T] to an Iter[bool] with a single value that is true if the Iter[T] is empty or all elements
// pass the filter. Boolean short circuit logic stops on first case where filter returns false.
// Calls ReduceToBool(filter, true, false).
func AllMatch[T any](filter func(T) bool) func(iter.Iter[T]) iter.Iter[bool] {
	return ReduceToBool(filter, true, false)
}

// AnyMatch reduces Iter[T] to an Iter[bool] with a single value that is true if the Iter[T] is non-empty and any
// element passes the filter. Boolean short circuit logic stops on first case where filter returns true.
// Calls ReduceToBool(filter, true, false).
func AnyMatch[T any](filter func(T) bool) func(iter.Iter[T]) iter.Iter[bool] {
	return ReduceToBool(filter, false, true)
}

// NoneMatch reduces Iter[T] to an Iter[bool] with a single value that is true if the Iter[T] is empty or no elements
// pass the filter. Boolean short circuit logic stops on first case where filter returns true.
// Calls ReduceToBool(!filter, true, false).
func NoneMatch[T any](filter func(T) bool) func(iter.Iter[T]) iter.Iter[bool] {
	return ReduceToBool(func(t T) bool { return !filter(t) }, true, false)
}

// Count reduces Iter[T] to an Iter[int] with a single value that is the number of elements in the Iter[T].
func Count[T any](it iter.Iter[T]) iter.Iter[int] {
	return ReduceTo[T, int](func(c int, _ T) int { return c + 1 }, 0)(it)
}

// Distinct reduces Iter[T] to an Iter[T] with distinct values.
// Distinct is a stateful transform that has to track unique values across iterator Next and Value calls.
//
// Distinct uses Generator internally to ensure what whenever a new Iter is encountered, a new state of an empty set of
// values is generated. This allows a composition to be stored in a variable and reused across data sets correctly.
//
// If you want Distinct to have one state across multiple Iters, use Concat to create a single Iter that traverses them.
func Distinct[T comparable](it iter.Iter[T]) iter.Iter[T] {
	return Generator(func() func(iter.Iter[T]) iter.Iter[T] {
		vals := map[T]bool{}

		return Filter[T](func(val T) bool {
			haveIt := vals[val]
			if !haveIt {
				vals[val] = true
			}

			return !haveIt
		})
	})(it)
}

// Duplicate reduces Iter[T] to an Iter[T] with duplicate values.
// Like Distinct, a given duplicate will only appear once.
// The order of the elements is the order in which the second occurence of each value appears.
// Eg, for the input 1,2,2,1 the result is 2,1 since the second value of 2 appears before the second value of 1.
// See Distinct for an explanation of statefulness and the usage of Generator.
func Duplicate[T comparable](it iter.Iter[T]) iter.Iter[T] {
	return Generator(func() func(iter.Iter[T]) iter.Iter[T] {
		vals := map[T]int{}

		return Filter[T](func(val T) bool {
			count := vals[val]
			if count < 3 {
				count++
			}
			vals[val] = count

			return count == 2
		})
	})(it)
}

// Reverse reverses all the elements.
// The input iter must have a finite size.
func Reverse[T any](it iter.Iter[T]) iter.Iter[T] {
	// Get values into a slice
	slc, err := ReduceToSlice(it).Next()

	if err != nil {
		// Unable to get any values to reverse
		return iter.SetError(iter.OfEmpty[T](), err)
	}

	// Reverse elements
	funcs.SliceReverse(slc)

	// Return iterator of reversed elements
	return iter.Of(slc...)
}

// SortOrdered sorts an Ordered type that is implicitly sortable using funcs.SliceSortOrdered.
// The input iter must have a finite size.
func SortOrdered[T constraint.Ordered](it iter.Iter[T]) iter.Iter[T] {
	// Get values into a slice
	slc, err := ReduceToSlice(it).Next()

	if err != nil {
		// Unable to get any values to sort
		return iter.SetError(iter.OfEmpty[T](), err)
	}

	// Sort elements
	funcs.SliceSortOrdered(slc)

	// Successfully return sorted iter
	return iter.Of(slc...)
}

// SortComplex sorts a Complex type using funcs.SliceSortComplex.
// The input iter must have a finite size.
func SortComplex[T constraint.Complex](it iter.Iter[T]) iter.Iter[T] {
	// Get values into a slice
	slc, err := ReduceToSlice(it).Next()

	if err != nil {
		// Unable to get any values to sort
		return iter.SetError(iter.OfEmpty[T](), err)
	}

	// Sort elements
	funcs.SliceSortComplex(slc)

	// Successfully return sorted iter
	return iter.Of(slc...)
}

// SortCmp sorts a Cmp type using funcs.SliceSortCmp.
// The input iter must have a finite size.
func SortCmp[T constraint.Cmp[T]](it iter.Iter[T]) iter.Iter[T] {
	// Get values into a slice
	slc, err := ReduceToSlice(it).Next()

	if err != nil {
		// Unable to get any values to sort
		return iter.SetError(iter.OfEmpty[T](), err)
	}

	// Sort elements
	funcs.SliceSortCmp(slc)

	// Successfully return sorted iter
	return iter.Of(slc...)
}

// SortBy sorts any type using funcs.SliceSortBy and the given comparator.
// The input iter must have a finite size.
func SortBy[T any](less func(T, T) bool) func(iter.Iter[T]) iter.Iter[T] {
	return func(it iter.Iter[T]) iter.Iter[T] {
		// Get values into a slice
		slc, err := ReduceToSlice(it).Next()

		if err != nil {
			// Unable to get any values to sort
			return iter.SetError(iter.OfEmpty[T](), err)
		}

		// Sort elements
		funcs.SliceSortBy(slc, less)

		// Successfully return sorted iter
		return iter.Of(slc...)
	}
}

// ==== Math

// Abs converts all elements into their absolute values.
// See MapError for error handling in cases where there is no corresponding value for a negative integer.
func Abs[T constraint.SignedInteger](it iter.Iter[T]) iter.Iter[T] {
	return MapError(func(v T) (T, error) {
		if v < 0 {
			if v = -v; v < 0 {
				var zv T
				return zv, fmt.Errorf(absErrMsg, v, v)
			}
		}

		return v, nil
	})(it)
}

// AbsBigOps is the *big.Int, *big.Float, *big.Rat specialization of Abs
func AbsBigOps[T constraint.BigOps[T]](it iter.Iter[T]) iter.Iter[T] {
	return Map(func(v T) T {
		return v.Abs(v)
	})(it)
}

// AvgInt reduces all signed integer elements in the input set to their average. If the input set is empty, the result is empty.
// The average is rounded.
// See math.Add, math.Div.
func AvgInt[T constraint.SignedInteger](it iter.Iter[T]) iter.Iter[T] {
	return iter.NewIter(
		func() (T, error) {
			var (
				sum   T
				count int64
				val   T
				avg   T
				zv    T
				err   error
			)

			for {
				if val, err = it.Next(); err == nil {
					// Sum all values and count them, checking for over/underflow
					if err = math.AddInt(val, &sum); err != nil {
						return zv, err
					}
					count++
				} else if err == iter.EOI {
					if count == 0 {
						// Empty result
						return zv, err
					}

					// Non-empty result - divisor = zero handled above, so Div will never fail
					var de, q int64
					conv.To(sum, &de)
					math.Div(de, count, &q)
					conv.To(q, &avg)
					break
				} else {
					return zv, err
				}
			}

			return avg, nil
		},
	)
}

// AvgUint reduces all unsigned integer elements in the input set to their average. If the input set is empty, the result is empty.
// The average is rounded.
// See math.Add, math.Div.
func AvgUint[T constraint.UnsignedInteger](it iter.Iter[T]) iter.Iter[T] {
	return iter.NewIter(
		func() (T, error) {
			var (
				sum   T
				count uint64
				val   T
				avg   T
				zv    T
				err   error
			)

			for {
				if val, err = it.Next(); err == nil {
					// Sum all values and count them, checking for overflow
					if err = math.AddUint(val, &sum); err != nil {
						return zv, err
					}
					count++
				} else if err == iter.EOI {
					if count == 0 {
						// Empty result
						return zv, err
					}

					// Non-empty result - divisor = zero handled above, so Div will never fail
					var de, q uint64
					conv.To(sum, &de)
					math.Div(de, count, &q)
					conv.To(q, &avg)
					break
				} else {
					return zv, err
				}
			}

			return avg, nil
		},
	)
}

// AvgBigOps reduces all *big.Int, *big.Float, or *big.Rat elements in the input set to their average.
// If the input set is empty, the result is empty. The average is rounded only for *big.Int.
// See math.DivBigOps
func AvgBigOps[T constraint.BigOps[T]](it iter.Iter[T]) iter.Iter[T] {
	return iter.NewIter(
		func() (T, error) {
			var (
				sum   T
				count T
				one   T
				val   T
				avg   T
				zv    T
				err   error
			)
			conv.ToBigOps(0, &sum)
			conv.ToBigOps(0, &count)
			conv.ToBigOps(1, &one)

			for {
				if val, err = it.Next(); err == nil {
					// Sum all values and count them
					sum.Add(sum, val)
					count.Add(count, one)
				} else if err == iter.EOI {
					if count.Sign() == 0 {
						// Empty result
						return zv, err
					}

					// Non-empty result - divisor = zero handled above, so Div will never fail
					math.DivBigOps(sum, count, &avg)
					break
				} else {
					return zv, err
				}
			}

			return avg, nil
		},
	)
}

// Max calculates the maximum value of any primitive numeric type or string. If the input set is empty, the result is empty.
func Max[T constraint.Ordered](it iter.Iter[T]) iter.Iter[T] {
	return Reduce(
		func(a, b T) T {
			return funcs.Ternary(a > b, a, b)
		},
	)(it)
}

// MaxCmp calculates the maximum value of any type implementing the Cmp interface, which includes all the big types.
// If the input set is empty, the result is empty.
func MaxCmp[T constraint.Cmp[T]](it iter.Iter[T]) iter.Iter[T] {
	return Reduce(
		func(a, b T) T {
			return funcs.Ternary(a.Cmp(b) > 0, a, b)
		},
	)(it)
}

// Min calculates the minimum value of any primitive numeric type or string. If the input set is empty, the result is empty.
func Min[T constraint.Ordered](it iter.Iter[T]) iter.Iter[T] {
	return Reduce(
		func(a, b T) T {
			return funcs.Ternary(a < b, a, b)
		},
	)(it)
}

// MaxCmp calculates the maximum value of any type implementing the Cmp interface, which includes all the big types.
// If the input set is empty, the result is empty.
func MinCmp[T constraint.Cmp[T]](it iter.Iter[T]) iter.Iter[T] {
	return Reduce(
		func(a, b T) T {
			return funcs.Ternary(a.Cmp(b) < 0, a, b)
		},
	)(it)
}

// Sum reduces all elements in the input set to their sum. If the input set is empty, the result is empty.
func Sum[T constraint.IntegerAndFloat](it iter.Iter[T]) iter.Iter[T] {
	return Reduce(
		func(a, b T) T {
			return a + b
		},
	)(it)
}

// SumBigOps is the *big.Int, *big.Float, *big.Rat specialization of Sum
func SumBigOps[T constraint.BigOps[T]](it iter.Iter[T]) iter.Iter[T] {
	var sum T
	conv.ToBigOps(0, &sum)

	return Reduce(
		func(a, b T) T {
			return sum.Add(a, b)
		},
	)(it)
}

// ==== Parallel

// generateRanges does the work of generating slice ranges from the optional PInfo.
// numItems must be at least 2, or the results may not be correct.
func generateRanges(numItems uint, info []PInfo) [][]uint {
	// Determine the slice ranges of each thread
	var sliceRanges [][]uint

	if len(info) == 0 {
		// Use square root when no PInfo given
		bucketSize := uint(gomath.Round(gomath.Sqrt(float64(numItems))))
		numThreads, remainder := bits.Div(0, numItems, bucketSize)

		// Algorithm has int sqrt number of threads + additional thread if remainder > 0
		for i, start, end := uint(0), uint(0), uint(0); i < numThreads; i++ {
			end = start + bucketSize
			sliceRanges = append(sliceRanges, []uint{start, end})

			start = end
		}

		if remainder > uint(0) {
			sliceRanges = append(sliceRanges, []uint{numItems - remainder, numItems})
		}
	} else {
		inf := info[0]
		if inf.PUnit == Threads {
			// User specified number of threads, calculate bucket size
			// Number of threads cannot exceed number of items
			numThreads := uint(gomath.Min(float64(numItems), gomath.Max(1, float64(inf.N))))
			bucketSize, remainder := bits.Div(0, numItems, numThreads)

			// Algorithm has number of threads given where first remainder threads have 1 additional item
			for start, end := uint(0), uint(0); start < numItems; start = end {
				end = start + bucketSize
				if remainder > 0 {
					end++
					remainder--
				}
				sliceRanges = append(sliceRanges, []uint{start, end})
			}
		} else {
			// User specified bucket size, calculate number of threads
			// Bucket size cannot exceed number of items
			bucketSize := uint(gomath.Min(float64(numItems), gomath.Max(1, float64(inf.N))))
			numThreads, remainder := bits.Div(0, numItems, bucketSize)

			// Algorithm has bucket size given where remainder is an additional thread
			for i, start, end := uint(0), uint(0), uint(0); i < numThreads; i++ {
				end = start + bucketSize
				sliceRanges = append(sliceRanges, []uint{start, end})

				start = end
			}

			if remainder > uint(0) {
				sliceRanges = append(sliceRanges, []uint{numItems - remainder, numItems})
			}
		}
	}

	return sliceRanges
}

// Parallel collects all the items of the source iter into a []T and divvies them up into buckets, then uses a set of
// threads, one per bucket, to process the items using the given set of transforms. If the optional PInfo is provided,
// The number N is interpreted in one of two ways, depending on the PUnit:
//
// Threads: N is the number of threads, the bucket size for each thread is number of items / N, with remainder r
//
//	distributed across first r threads. If number of items <= N, a single thread is used.
//
// Items:   N is the bucket size, the number of threads is number of items / N, with remainder r handled by an
//
//	additional thread. If number of items <= N, a single thread is used.
//
// If no PInfo is provided, the number of threads is the square root of the number of items, so that each thread has the
// same number of items - except for the last thread, which may have slightly less.
//
// If the input has no items, an empty Iter is returned.
// If the input has one item, the transforms are performed in the same thread.
// If the input has two or more items, the above algorithm is used to perform transforms in two or more threads.
//
// If types T and U are the same, then a single slice is allocated to contain the input and modified in place to produce
// the output. Otherwise, two slices are allocated, one for input and one for output.
func Parallel[T, U any](transforms func(iter.Iter[T]) iter.Iter[U], info ...PInfo) func(iter.Iter[T]) iter.Iter[U] {
	return func(source iter.Iter[T]) iter.Iter[U] {
		// Get values into a slice
		input, err := ReduceToSlice(source).Next()

		if err != nil {
			// Unable to get any values
			return iter.SetError(iter.OfEmpty[U](), err)
		}

		// Get number of items from source
		numItems := uint(len(input))

		switch numItems {
		case 0:
			// If the source is empty, nothing to do, just return an empty Iter
			return iter.OfEmpty[U]()
		case 1:
			// If the source has 1 element, don't bother with a separate thread, just return the result
			return transforms(iter.OfOne(input[0]))
		}

		// At least two items, as required by generateRanges. Get slice ranges.
		sliceRanges := generateRanges(numItems, info)

		// Determine a slice to use for the output
		var (
			zt     T
			zu     U
			output []U
		)
		if reflect.TypeOf(zt) == reflect.TypeOf(zu) {
			// T and U are the same type, modify input slice in place
			output = any(input).([]U)
		} else {
			// Create a separate output slice
			output = make([]U, numItems)
		}

		// Create a WaitGroup that can wait for all threads to complete
		var (
			wg   sync.WaitGroup
			errs = make([]error, len(sliceRanges))
		)

		// The function to execute in each thread, accepting source and target subslices
		threadFn := func(threadNum int, in []T, out []U) {
			// Decrement number of threads remaining once transforms are complete
			defer wg.Done()

			// Perform transforms and copy to output
			_, threadErr := ReduceIntoSlice(out)(transforms(iter.Of(in...))).Next()

			// In case an error occurs, populate the appropriate errs slot
			errs[threadNum] = threadErr
		}

		// Create the threads, passing a subslice of input to each thread for processing
		for threadNum, sliceRange := range sliceRanges {
			wg.Add(1)
			go threadFn(threadNum, input[sliceRange[0]:sliceRange[1]], output[sliceRange[0]:sliceRange[1]])
		}

		// Wait for threads to complete
		wg.Wait()

		// If any errors occur, return first error found (may not be first error in execution, as threads complete in any order)
		err = nil
		for _, e := range errs {
			if e != nil {
				err = e
				break
			}
		}

		return funcs.Ternary(err == nil, iter.Of(output...), iter.SetError(iter.OfEmpty[U](), err))
	}
}
