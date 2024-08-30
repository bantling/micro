package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	goio "io"
	"strings"

	"github.com/bantling/micro/tuple"
	"github.com/bantling/micro/union"
)

// ==== Constants

var (
	errOfIterNeedsIterator = fmt.Errorf("OfIter requires a non-nil iterating function")
	errValueExpected       = fmt.Errorf("Value has to be called after Next")
	errNextExpected        = fmt.Errorf("Next has to be called before Value")
	errNoMoreValues        = fmt.Errorf("Value cannot be called after Next returns false")
	EOI                    = fmt.Errorf("End of Iteration")
)

// ==== Types

// Iter defines iteration for type T, which works as follows:
// - Next returns (next value, nil) if there is another value, or (zero value, EOI) if there is not.
//   - It is possible for Next to return (zero value, some other error).
//   - Next should never return (non zero value, non nil error).
//   - Once Next returns (zero value, EOI or problem error), Next will continue to return (zero value, EOI or problem error).
//
// - Unread places the given value at the end of a buffer of values
//   - Next consults the buffer before calling the underlying iterating function
//   - Next returns values in reverse order of Unreads (eg Unread(1); Unread(2) results in Next returning 2, then 1)
type Iter[T any] interface {
	Next() (T, error)
	Unread(T)
}

// IterImpl is the common implementation of Iter[T], based on an underlying iterating function.
type IterImpl[T any] struct {
	iterFn  func() (T, error)
	buffer  []T
	lastErr error
}

// ==== Construct

// OfIter constructs an Iter[T] from an iterating function that returns (T, error).
// The function must return (nextItem, nil) for every item available to iterate, then return (invalid, EOI) on the
// next call after the last item, where invalid is any value of type T.
// If some actual error occurs attempting to read the next value, then the function must return (invalid, non-nil non-EOI error).
// Once the function returns a non-nil error, it will never be called again.
// Panics if iterFn is nil.
//
// See IterImpl.
func OfIter[T any](iterFn func() (T, error)) Iter[T] {
	if iterFn == nil {
		panic(errOfIterNeedsIterator)
	}

	return &IterImpl[T]{iterFn: iterFn}
}

// Of constructs an Iter[T] that iterates the items passed.
// The intention is hard-coded values are passed.
//
// See SliceIterGen.
func Of[T any](items ...T) Iter[T] {
	return OfIter[T](SliceIterGen[T](items))
}

// OfEmpty constructs an Iter[T] that iterates no values.
//
// See NoValueIterGen.
func OfEmpty[T any]() Iter[T] {
	return OfIter[T](NoValueIterGen[T]())
}

// OfOne constructs an Iter[T] that iterates a single value.
//
// See SingleValueIterGen.
func OfOne[T any](item T) Iter[T] {
	return OfIter[T](SingleValueIterGen[T](item))
}

// OfSlice constructs an Iter[T] that iterates the slice values passed.
// The intention is the slice may be large, and passing the slice by reference is better than using varargs like Of(...).
//
// See SliceIterGen
func OfSlice[T any](items []T) Iter[T] {
	return OfIter[T](SliceIterGen[T](items))
}

// Of constructs an Iter[tuple.Two[K, V]] that iterates the items passed.
//
// See MapIterGen.
func OfMap[K comparable, V any](items map[K]V) Iter[tuple.Two[K, V]] {
	return OfIter[tuple.Two[K, V]](MapIterGen[K, V](items))
}

// OfReader constructs an Iter[byte] that iterates the bytes of a Reader.
//
// See ReaderIterGen.
func OfReader(src goio.Reader) Iter[byte] {
	return OfIter[byte](ReaderIterGen(src))
}

// OfReaderAsRunes constructs an Iter[rune] that iterates the UTF-8 runes of a Reader.
//
// See ReaderAsRunesIterGen.
func OfReaderAsRunes(src goio.Reader) Iter[rune] {
	return OfIter(ReaderAsRunesIterGen(src))
}

// OfStringAsRunes constructs an Iter[rune] that iterates runes of a string.
//
// See SliceIterGen.
func OfStringAsRunes(src string) Iter[rune] {
	return OfIter(SliceIterGen([]rune(src)))
}

// OfReaderAsLines constructs an Iter[string] that iterates the UTF-8 lines of a Reader.
//
// See ReaderAsLinesIterGen.
func OfReaderAsLines(src goio.Reader) Iter[string] {
	return OfIter(ReaderAsLinesIterGen(src))
}

// OfStringAsLines constructs an Iter[rune] that iterates lines of a string.
//
// See ReaderAsLinesIterGen.
func OfStringAsLines(src string) Iter[string] {
	return OfIter(ReaderAsLinesIterGen(strings.NewReader(src)))
}

// OfCSV constructs an Iter[[]string] that iterates lines of a csv document inan io.Reader.
func OfCSV(src goio.Reader) Iter[[]string] {
	return OfIter(CSVIterGen(src))
}

// Concatenate any number of Iter[T] into a single Iter[T] that iterates all the elements of each Iter[T], until the
// last element of the last iterator has been returned.
func Concat[T any](iters ...Iter[T]) Iter[T] {
	return OfIter(ConcatIterGen(iters))
}

// ==== IterImpl Methods

// Next returns (value, nil) if there is another item to be read by Value.
// When Next returns (zero value, EOI), further calls return (zero value, EOI).
// Next does not ever return both a value and an error.
func (it *IterImpl[T]) Next() (T, error) {
	// Check buffer for values placed by Unread
	if len(it.buffer) > 0 {
		val := it.buffer[len(it.buffer)-1]
		it.buffer = it.buffer[0 : len(it.buffer)-1]
		return val, nil
	}

	// Check if we still have values to acquire via iterating function
	if it.iterFn != nil {
		// Try to get next value
		val, err := it.iterFn()

		if err == nil {
			// Got a value, may still be more left
			return val, nil
		}

		// The value is invalid, error could be EOI or a problem.
		// Don't try to call iterating function again, previous value was the last - function can't change its mind.
		it.iterFn = nil
		it.lastErr = err

		// Return (zero value, EOI or problem)
		var zv T
		return zv, err
	}

	// Called again after already returning (zero value, EOI or problem).
	// Continue to repeat (zero value, EOI or problem).
	var zv T
	return zv, it.lastErr
}

// Unread adds the given value to an internal buffer, to be returned by Next in reverse order.
func (it *IterImpl[T]) Unread(val T) {
	it.buffer = append(it.buffer, val)
}

// ==== Operations on an Iter

// Maybe converts the result of Next into a Result[T] to represent the result as a single type.
func Maybe[T any](it Iter[T]) union.Result[T] {
	return union.OfResultError(it.Next())
}

// SetError sets a particular error to occur instead of the first non-nil error the given iterator returns.
func SetError[T any](it Iter[T], err error) Iter[T] {
	return OfIter[T](func() (T, error) {
		if v, e := it.Next(); e == nil {
			return v, e
		} else {
			var zv T
			return zv, err
		}
	})
}
