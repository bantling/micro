package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/bits"
	"reflect"
	"sync"

	"github.com/bantling/micro/go/constraint"
	"github.com/bantling/micro/go/funcs"
)

var (
	errSkipLimitValueCannotBeNegative = fmt.Errorf("The Skip or Limit value cannot be negative")
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

// ==== Functions that provide the foundation for all other functions

// First combines Next and Value together in a single call.
// If there is another value, then the next value is returned, else a panic occurs.
// First is not a method so it can be used in a composition.
func First[T any](it *Iter[T]) T {
	it.Next()
	return it.Value()
}

// Map constructs a new Iter[U] from an Iter[T] and a func that transforms a T to a U.
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
// - If no identity value is given, the zero value is used
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

// ReduceToBool is similar to ReduceTo, except that it uses boolean short circuit logic to stop iterating early if
// possible. If stopIfTrue is true, then early termination occurs on first call to reducer that returns true, else it
// occurs on first call to reducer that returns false.
func ReduceToBool[T any](
	predicate func(T) bool,
	identity bool,
	stopIfTrue bool,
) func(*Iter[T]) *Iter[bool] {
	return func(it *Iter[T]) *Iter[bool] {
		var done bool

		return NewIter(func() (bool, bool) {
			if done {
				return false, false
			}

			done = true

			if !it.Next() {
				// 0 elements = identity
				return identity, true
			} else {
				// At least one element, call reducer with identity and first element
				result := predicate(it.Value())
				if result == stopIfTrue {
					// Stop early if result matches stopping condition
					return result, true
				}

				for it.Next() {
					// If there are more elements, make cumulative reducer calls, combining old and new results
					result = predicate(it.Value())
					if result == stopIfTrue {
						// Stop early if result matches stopping condition
						return result, true
					}
				}

				return result, true
			}
		})
	}
}

// ReduceToSlice reduces an Iter[T] into a Iter[[]T] that contains a single element if type []T.
// Eg, an Iter[int] of 1,2,3,4,5 becomes an Iter[[]int] of [1,2,3,4,5].
// An empty Iter is reduced to a zero length slice.
func ReduceToSlice[T any](it *Iter[T]) *Iter[[]T] {
	var done bool

	return NewIter(func() ([]T, bool) {
		if done {
			return nil, false
		}

		done = true

		slc := []T{}
		for it.Next() {
			slc = append(slc, it.Value())
		}

		return slc, true
	})
}

// ReduceIntoSlice is the same as ReduceToSlice, except that:
// - It accepts a target slice to append results to
// - It generates a transform
//
// The generated transform panics if the target slice length is not at least as many elements as the source iter
func ReduceIntoSlice[T any](slc []T) func(*Iter[T]) *Iter[[]T] {
	return func(it *Iter[T]) *Iter[[]T] {
		var (
			done bool
			i    int
		)

		return NewIter(func() ([]T, bool) {
			if done {
				return nil, false
			}

			done = true

			for it.Next() {
				slc[i] = it.Value()
				i++
			}

			return slc, true
		})
	}
}

// ExpandSlices is the opposite of ReduceToSlice: an Iter[[]int] of [1,2,3,4,5] becomes an Iter[int] of 1,2,3,4,5.
// If the source Iter contains multiple slices, they are combined together into one set of data (skipping nil and empty
// slices), so that an Iter[[]int] of [1,2,3], nil, [], [4,5] also becomes an Iter[int] of 1,2,3,4,5.
// An empty Iter or an Iter with nil/empty slices is expanded to an empty Iter.
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
// An empty Iter is reduced to an empty map.
func ReduceToMap[K comparable, V any](it *Iter[KeyValue[K, V]]) *Iter[map[K]V] {
	var done bool

	return NewIter(func() (map[K]V, bool) {
		if done {
			return nil, false
		}

		done = true

		m := map[K]V{}
		for it.Next() {
			kv := it.Value()
			m[kv.Key] = kv.Value
		}

		return m, true
	})
}

// ExpandMaps is the opposite of ReduceToMap: an Iter[map[int]string] of {1: "1", 2: "2", 3: "3"} becomes an
// Iter[KeyValue[int, string]] of {1: "1"}, {2: "2"}, {3: "3"].
// If the source Iter contains multiple maps, they are combined together into one set of data (skipping nils),
// so that an Iter[map[int]string] of {1: "1", 2: "2"}, nil, {}, {3: "3"} also becomes
// an Iter[KeyValue[int, string]] of {1: "1"}, {2: "2"}, {3: "3"}.
// An empty Iter or an Iter with nil/empty maps is expanded to an empty Iter.
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

// Skip skips the first n elements, then iteration continues from there.
// If there are n or fewer elements in total, then the resulting iter is empty.
// Panics if n < 0.
//
// Note that for the set 1,2,3,4,5 the composition os Skip(1),Limit(3) will first skip 1 then limit to 2,3,4;
// whereas the composition of Limit(3),Skip(1) will first limit to 1,2,3 then skip 1 returning 2,3.
func Skip[T any](n int) func(*Iter[T]) *Iter[T] {
	if n < 0 {
		panic(errSkipLimitValueCannotBeNegative)
	}

	return func(it *Iter[T]) *Iter[T] {
		skip := n

		return NewIter(func() (T, bool) {
			for ; (skip > 0) && it.Next(); skip-- {
				it.Value()
			}

			if it.Next() {
				return it.Value(), true
			}

			var zv T
			return zv, false
		})
	}
}

// Limit returns the first n elements, then iteration stops and all further elements are ignored.
// If there fewer than n elements in total, then all n elements are returned.
// Panics if n < 0.
//
// Note that for the set 1,2,3,4,5 the composition of Skip(1),Limit(3) will first skip 1 then limit to 2,3,4;
// whereas the composition of Limit(3),Skip(1) will first limit to 1,2,3 then skip 1 returning 2,3.
func Limit[T any](n int) func(*Iter[T]) *Iter[T] {
	if n < 0 {
		panic(errSkipLimitValueCannotBeNegative)
	}

	return func(it *Iter[T]) *Iter[T] {
		limit := n

		return NewIter(func() (T, bool) {
			if (limit > 0) && it.Next() {
				limit--
				return it.Value(), true
			}

			var zv T
			return zv, false
		})
	}
}

// Peek executes a func for every item being iterated, which is a side effect.
func Peek[T any](fn func(T)) func(*Iter[T]) *Iter[T] {
	return func(it *Iter[T]) *Iter[T] {
		return NewIter(func() (T, bool) {
			if it.Next() {
				val := it.Value()
				fn(val)
				return val, true
			}

			var zv T
			return zv, false
		})
	}
}

// Generator receives a generator (a func of no args that returns a func of Iter[T] -> Iter[U], and detects if the
// *Iter[T] has changed address. If so, it internally generates a new function by invoking the generator.
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
func Generator[T, U any](gen func() func(*Iter[T]) *Iter[U]) func(*Iter[T]) *Iter[U] {
	var (
		currentIter *Iter[T]
		fn          func(*Iter[T]) *Iter[U]
	)

	return func(it *Iter[T]) *Iter[U] {
		if currentIter != it {
			fn = gen()
		}

		return fn(it)
	}
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
func Transform[T, U any](transformer func(*Iter[T]) (U, bool)) func(*Iter[T]) *Iter[U] {
	return func(it *Iter[T]) *Iter[U] {
		return NewIter(func() (U, bool) {
			return transformer(it)
		})
	}
}

// ==== Functions based on foundational functions

// AllMatch reduces Iter[T] to an Iter[bool] with a single value that is true if the Iter[T] is empty or all elements
// pass the filter. Boolean short circuit logic stops on first case where filter returns false.
// Calls ReduceToBool(filter, true, false).
func AllMatch[T any](filter func(T) bool) func(*Iter[T]) *Iter[bool] {
	return ReduceToBool(filter, true, false)
}

// AnyMatch reduces Iter[T] to an Iter[bool] with a single value that is true if the Iter[T] is non-empty and any
// element passes the filter. Boolean short circuit logic stops on first case where filter returns true.
// Calls ReduceToBool(filter, true, false).
func AnyMatch[T any](filter func(T) bool) func(*Iter[T]) *Iter[bool] {
	return ReduceToBool(filter, false, true)
}

// NoneMatch reduces Iter[T] to an Iter[bool] with a single value that is true if the Iter[T] is empty or no elements
// pass the filter. Boolean short circuit logic stops on first case where filter returns true.
// Calls ReduceToBool(!filter, true, false).
func NoneMatch[T any](filter func(T) bool) func(*Iter[T]) *Iter[bool] {
	return ReduceToBool(func(t T) bool { return !filter(t) }, true, false)
}

// Count reduces Iter[T] to an Iter[int] with a single value that is the number of elements in the Iter[T].
func Count[T any]() func(*Iter[T]) *Iter[int] {
	return ReduceTo[T, int](func(c int, _ T) int { return c + 1 })
}

// Distinct reduces Iter[T] to an Iter[T] with distinct values.
// Distinct is a stateful transform that has to track unique values across iterator Next and Value calls.
//
// Distinct uses Generator internally to ensure what whenever a new Iter is encountered, a new install state of an empty
// set of values is generated. This allows a composition to be stored in a variable and reused across data
// sets correctly.
//
// If you want Distinct to have one state across multiple Iters, use Concat to create a single Iter that traverses them.
func Distinct[T comparable]() func(*Iter[T]) *Iter[T] {
	return Generator(func() func(*Iter[T]) *Iter[T] {
		vals := map[T]bool{}

		return Filter[T](func(val T) bool {
			haveIt := vals[val]
			if !haveIt {
				vals[val] = true
			}

			return !haveIt
		})
	})
}

// Duplicate reduces Iter[T] to an Iter[T] with duplicate values.
// Like Distinct, a given duplicate will only appear once.
// The order of the elements is the order in which the second occurence of each value appears.
// Eg, for the input 1,2,2,1 the result is 2,1 since the second value of 2 appears before the second value of 1.
// See Distinct for an explanation of statefulness and the usage of Generator.
func Duplicate[T comparable]() func(*Iter[T]) *Iter[T] {
	return Generator(func() func(*Iter[T]) *Iter[T] {
		vals := map[T]int{}

		return Filter[T](func(val T) bool {
			count := vals[val]
			if count < 3 {
				count++
			}
			vals[val] = count

			return count == 2
		})
	})
}

// Reverse reverse all the elements.
// The input iter must have a finite size.
func Reverse[T any]() func(*Iter[T]) *Iter[T] {
	return func(it *Iter[T]) *Iter[T] {
		slc := First(ReduceToSlice(it))
		funcs.Reverse(slc)

		return Of(slc...)
	}
}

// SortOrdered sorts an Ordered type that is implicitly sortable.
// The input iter must have a finite size.
func SortOrdered[T constraint.Ordered]() func(*Iter[T]) *Iter[T] {
	return func(it *Iter[T]) *Iter[T] {
		slc := First(ReduceToSlice(it))
		funcs.SortOrdered(slc)

		return Of(slc...)
	}
}

// SortComplex sorts a Complex type using funcs.SortComplex.
// The input iter must have a finite size.
func SortComplex[T constraint.Complex]() func(*Iter[T]) *Iter[T] {
	return func(it *Iter[T]) *Iter[T] {
		slc := First(ReduceToSlice(it))
		funcs.SortComplex(slc)

		return Of(slc...)
	}
}

// SortCmp sorts a Cmp type using funcs.SortCmp.
// The input iter must have a finite size.
func SortCmp[T constraint.Cmp[T]]() func(*Iter[T]) *Iter[T] {
	return func(it *Iter[T]) *Iter[T] {
		slc := First(ReduceToSlice(it))
		funcs.SortCmp(slc)

		return Of(slc...)
	}
}

// SortBy sorts any type using funcs.SortBy and the given comparator.
// The input iter must have a finite size.
func SortBy[T any](less func(T, T) bool) func(*Iter[T]) *Iter[T] {
	return func(it *Iter[T]) *Iter[T] {
		slc := First(ReduceToSlice(it))
		funcs.SortBy(slc, less)

		return Of(slc...)
	}
}

// ==== Parallel

// Parallel collects all the items of the source iter into a []T and divvies them up into buckets, then uses a set of
// threads, one per bucket, to process the items using the given set of transforms. If the optional PInfo is provided,
// The number N is interpreted in one of two ways, depending on the PUnit:
//
// Threads: N is the number of threads, the bucket size for each thread is number of items / N, with remainder r
//          distributed across first r threads. If number of items <= N, a single thread is used.
// Items:   N is the bucket size, the number of threads is number of items / N, with remainder r handled by an
//          additional thread. If number of items <= N, a single thread is used.
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
func Parallel[T, U any](transforms func(*Iter[T]) *Iter[U], info ...PInfo) func(*Iter[T]) *Iter[U] {
	return func(source *Iter[T]) *Iter[U] {
		var (
			// Collect all items from source iterator into a slice
			input = First(ReduceToSlice(source))
			// Get number of items fromm source
			numItems = uint(len(input))
		)

		switch numItems {
		case 0:
			// If the source is empty, nothing to do, just return an empty Iter
			return OfEmpty[U]()
		case 1:
			// If the source has 1 element, don't bother with a separate thread, just return the result
			return transforms(OfOne(input[0]))
		}

		// Determine the slice ranges of each thread
		var sliceRanges [][]uint

		if len(info) == 0 {
			// Use square root when no PInfo given
			bucketSize := uint(math.Round(math.Sqrt(float64(numItems))))
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
				numThreads := uint(math.Min(float64(numItems), math.Max(1, float64(inf.N))))
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
				bucketSize := uint(math.Min(float64(numItems), math.Max(1, float64(inf.N))))
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

		// Determine a slice to use for the output
		var (
			zt     T
			zu     U
			output []U
		)
		if reflect.TypeOf(zt) == reflect.TypeOf(zu) {
			// T and U aare the same type, modify input slice in place
			output = any(input).([]U)
		} else {
			// Create a separate output slice
			output = make([]U, numItems)
		}

		// Create a WaitGroup that can wait for all threads to complete
		var wg sync.WaitGroup

		// The function to execute in each thread, accepting source and target subslices
		threadFn := func(in []T, out []U) {
			// Decrement number of threads remaining once transforms are complete
			defer wg.Done()

			// Perform transforms and copy to output
			First(ReduceIntoSlice(out)(transforms(Of(in...))))
		}

		// Create the threads, passing a subslice of input to each thread for processing
		for _, sliceRange := range sliceRanges {
			wg.Add(1)
			go threadFn(input[sliceRange[0]:sliceRange[1]], output[sliceRange[0]:sliceRange[1]])
		}

		// Wait for threads to complete
		wg.Wait()

		return Of(output...)
	}
}
