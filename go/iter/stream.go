package iter

// SPDX-License-Identifier: Apache-2.0

// ==== Functions that provide the foundation for all other functions

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

// SetReduce is similar to Reduce, except that:
// - Multiple reductions can be made, such that the resulting Iter[T] returns multiple elements.
// - The reducer is a function generator
//   - The function is generated at the beginning of each new reduction, to allow each reduction to start with some
//     initial state. A such, no identity is needed, that can be part of the initial state.
//   - The generated function is a func(Iter[T]) (T, bool) instead of func(T, T) T to give the reducer the ability to
//     determine on its own how to deal with cases like no elements, one element, multiple elements, and
//     incomplete input.
//
// A simple example is reducing a set of integers into the sum of pairs of integers, so that the
// input [1,2,3,4,5] is reduced to the output [3,7,5].
func SetReduce[T, U any](
	reducer func() func(*Iter[T]) (U, bool),
) func(*Iter[T]) *Iter[U] {
	return func(it *Iter[T]) *Iter[U] {
		return NewIter(func() (U, bool) {
			return reducer()(it)
		})
	}
}
