package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"reflect"
)

// ==== Functions that provide the foundation for all other functions

// MustValue combines Next and Value together in a single call.
// If there is another value, then the next value is returned, else a panic occurs.
// First is not a method so it can be used as the last funciton in a composition of Iter functions
func First[T any](it *Iter[T]) T {
	it.Next()
	return it.Value()
}

// Map constructs an new Iter[U] from an Iter[T] and a func that transforms a T to a U.
func Map[T, U any](mapper func(T) U) func(*Iter[T]) *Iter[U] {
	return func(it *Iter[T]) *Iter[U] {
		return NewIter(func() (U, bool) {
			if it.Next() {
				return mapper(it.Value()), true
			}

			var zv U
			return zv, false
		})
	}
}

// Filter constructs a new Iter[T] from an Iter[T] and a func that returns true if a T passes the filter.
func Filter[T any](filter func(T) bool) func(*Iter[T]) *Iter[T] {
	return func(it *Iter[T]) *Iter[T] {
		return NewIter(func() (T, bool) {
			for it.Next() {
				if val := it.Value(); filter(val) {
					return val, true
				}
			}

			var zv T
			return zv, false
		})
	}
}

// Reduce reduces all elements in the input set to a single element of the same type, depending on two factors:
// - the number of elements in the input set (0, 1, or multiple)
// - whether or not the optional identity is provided
//
// If the optional identity is NOT provided:
// - 0 elements: empty
// - 1 elements: the element
// - multiple elements: reducer(reducer(reducer(first, second), third), ...)
//
//  If the optional identity IS provided:
// - 0 elements: identity
// - 1 elements: reducer(identity, the element)
// - multiple elements: reducer(reducer(reducer(identity, first), second), ...)
func Reduce[T any](
	reducer func(T, T) T,
	identity ...T,
) func(*Iter[T]) *Iter[T] {
	return func(it *Iter[T]) *Iter[T] {
		var done bool

		return NewIter(func() (T, bool) {
			var zv T

			if done {
				return zv, false
			}

			done = true

			if len(identity) == 0 {
				// No identity
				if !it.Next() {
					// 0 elements = Empty set
					return zv, false
				} else {
					// At least one element, start with first element
					result := it.Value()

					for it.Next() {
						// If there are more elements, make cumulative reducer calls
						result = reducer(result, it.Value())
					}

					return result, true
				}
			}

			// There is an identity
			identityVal := identity[0]
			if !it.Next() {
				// 0 elements = identity
				return identityVal, true
			} else {
				// At least one element, call reducer with identity and first element
				result := reducer(identityVal, it.Value())

				for it.Next() {
					// If there are more elements, make cumulative reducer calls
					result = reducer(result, it.Value())
				}

				return result, true
			}
		})
	}
}

// ReduceTo is similar to Reduce, except that:
// - The result does not have to be the same type
func ReduceTo[T, U any](
	reducer func(U, T) U,
	identity ...U,
) func(*Iter[T]) *Iter[U] {
	return func(it *Iter[T]) *Iter[U] {
		var done bool

		return NewIter(func() (U, bool) {
			var zv U

			if done {
				return zv, false
			}

			done = true

			if len(identity) == 0 {
				// No identity
				if !it.Next() {
					// 0 elements = Empty set
					return zv, true
				} else {
					// At least one element, call reducer with zero value and first element
					result := reducer(zv, it.Value())

					for it.Next() {
						// If there are more elements, make cumulative reducer calls, combining old and new results
						result = reducer(result, it.Value())
					}

					return result, true
				}
			}

			// There is an identity
			identityVal := identity[0]

			if !it.Next() {
				// 0 elements = identity
				return identityVal, true
			} else {
				// At least one element, call reducer with identity and first element
				result := reducer(identityVal, it.Value())

				for it.Next() {
					// If there are more elements, make cumulative reducer calls, combining old and new results
					result = reducer(result, it.Value())
				}

				return result, true
			}
		})
	}
}

// ReduceToSlice reduces an Iter[T] into a Iter[[]T] that contains a single element if type []T.
// Eg, an Iter[int] of 1,2,3,4,5 becomes an Iter[[]int] of [1,2,3,4,5].
func ReduceToSlice[T any](it *Iter[T]) *Iter[[]T] {
	var done bool

	return NewIter(func() ([]T, bool) {
		if done {
			return nil, false
		}

		var slc []T
		for it.Next() {
			slc = append(slc, it.Value())
		}

		done = true

		return slc, true
	})
}

// ExpandSlices is the opposite of ReduceToSlice: an Iter[[]int] of [1,2,3,4,5] becomes an Iter[int] of 1,2,3,4,5.
// If the source Iter contains multiple slices, they are combined together into one set of data (skipping nils),
// so that an Iter[[]int] of [1,2,3], nil, [], [4,5] also becomes an Iter[int] of 1,2,3,4,5.
func ExpandSlices[T any](it *Iter[[]T]) *Iter[T] {
	var (
		slc []T
		idx int
	)

	return NewIter(func() (T, bool) {
		if (slc == nil) || (idx == len(slc)) {
			// Search for next non-nil non-empty slice
			// Nilify slc var in case we just finished iterating last element of last slice, which is non-nil
			slc = nil

			for it.Next() {
				if slc = it.Value(); (slc != nil) && (len(slc) > 0) {
					idx = 0
					break
				}
			}

			// Stop if no more non-nil non-empty slices available
			if slc == nil {
				var zv T
				return zv, false
			}
		}

		val := slc[idx]
		idx++

		return val, true
	})
}

// ReduceToMap reduces an Iter[KeyValue[K, V]]] into a Iter[map[K, V]] that contains a single element if type map[K, V].
// Eg, an Iter[KeyValue[int]string] of {1: "1"}, {2: "2"} becomes an Iter[map[int]string] of {1: "1", 2: "2"}.
// If multiple KeyValue objects in the Iter have the same key, the last such object in iteration order determines the
// value for the key in the resulting map.
func ReduceToMap[K comparable, V any](it *Iter[KeyValue[K, V]]) *Iter[map[K]V] {
	var done bool

	return NewIter(func() (map[K]V, bool) {
		if done {
			return nil, false
		}

		m := map[K]V{}
		for it.Next() {
			kv := it.Value()
			m[kv.Key] = kv.Value
		}

		done = true

		return m, true
	})
}

// ExpandMaps is the opposite of ReduceToMap: an Iter[map[int]string] of {1: "1", 2: "2", 3: "3"} becomes an
// Iter[KeyValue[int, string]] of {1: "1"}, {2: "2"}, {3: "3"].
// If the source Iter contains multiple maps, they are combined together into one set of data (skipping nils),
// so that an Iter[map[int]string] of {1: "1", 2: "2"}, nil, {}, {3: "3"} also becomes
// an Iter[KeyValue[int, string]] of {1: "1"}, {2: "2"}, {3: "3"}.
func ExpandMaps[K comparable, V any](it *Iter[map[K]V]) *Iter[KeyValue[K, V]] {
	var (
		m  map[K]V
		mr *reflect.MapIter
	)

	return NewIter(func() (KeyValue[K, V], bool) {
		if (m == nil) || (!mr.Next()) {
			// Search for next non-nil non-empty map
			// Nilify m var in case we just finished iterating last element of last map, which is non-nil
			m = nil

			for it.Next() {
				if m = it.Value(); m != nil {
					if mr = reflect.ValueOf(m).MapRange(); mr.Next() {
						break
					}
				}
			}

			// Stop if no more non-nil non-empty maps available
			if m == nil {
				var zv KeyValue[K, V]
				return zv, false
			}
		}

		val := KeyValue[K, V]{mr.Key().Interface().(K), mr.Value().Interface().(V)}

		return val, true
	})
}

// Transform allows for an arbitrary transform of an Iter[T] to an Iter[U], where:
// - type T may or may not be the same as type U
// - a single U value may require iterating multiple T values (reduction)
// - a single T value may result in multiple U values (expansion)
//
// A simple example is reducing a set of integers into the sum of pairs of integers, so that an
// iter of 1,2,3,4,5 is reduced to an iter of 3,7,5.
//
// All iteration logic is handled by the transformer function, it can iterate as many elements as necessary.
// A transformer should only be used for cases where using other provided functions like Map, Filter, and Reduce(To)
// either won't work, or result in a process nobody understands.
func Transform[T, U any](
	transformer func(*Iter[T]) (U, bool),
) func(*Iter[T]) *Iter[U] {
	return func(it *Iter[T]) *Iter[U] {
		return NewIter(func() (U, bool) {
			return transformer(it)
		})
	}
}
