package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"io"
	"strings"

	"github.com/bantling/micro/go/tuple"
)

// ==== Constants

var (
	errNewIterNeedsIterator = fmt.Errorf("NewIter requires a non-nil iterating function")
	errValueExpected        = fmt.Errorf("Value has to be called after Next")
	errNextExpected         = fmt.Errorf("Next has to be called before Value")
	errNoMoreValues         = fmt.Errorf("Value cannot be called after Next returns false")
	EOI                     = fmt.Errorf("End of Iteration")
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

// NewIter constructs an Iter[T] from an iterating function that returns (T, error).
// The function must return (nextItem, nil) for every item available to iterate, then return (invalid, EOI) on the
// next call after the last item, where invalid is any value of type T.
// If some actual error occurs attempting read the next value, then the function must return (invalid, non-nil non-EOI error).
// Once the function returns a non-nil error, it will never be called again.
// Panics if iterFn is nil.
//
// See IterImpl.
func NewIter[T any](iterFn func() (T, error)) Iter[T] {
	if iterFn == nil {
		panic(errNewIterNeedsIterator)
	}

	return &IterImpl[T]{iterFn: iterFn}
}

// Of constructs an Iter[T] that iterates the items passed.
//
// See SliceIterGen.
func Of[T any](items ...T) Iter[T] {
	return NewIter[T](SliceIterGen[T](items))
}

// OfEmpty constructs an Iter[T] that iterates no values.
//
// See NoValueIterGen.
func OfEmpty[T any]() Iter[T] {
	return NewIter[T](NoValueIterGen[T]())
}

// OfOne constructs an Iter[T] that iterates a single value.
//
// See SingleValueIterGen.
func OfOne[T any](item T) Iter[T] {
	return NewIter[T](SingleValueIterGen[T](item))
}

// Of constructs an Iter[tuple.Two[K, V]] that iterates the items passed.
//
// See MapIterGen.
func OfMap[K comparable, V any](items map[K]V) Iter[tuple.Two[K, V]] {
	return NewIter[tuple.Two[K, V]](MapIterGen[K, V](items))
}

// OfReader constructs an Iter[byte] that iterates the bytes of a Reader.
//
// See ReaderIterGen.
func OfReader(src io.Reader) Iter[byte] {
	return NewIter[byte](ReaderIterGen(src))
}

// OfReaderAsRunes constructs an Iter[rune] that iterates the UTF-8 runes of a Reader.
//
// See ReaderAsRunesIterGen.
func OfReaderAsRunes(src io.Reader) Iter[rune] {
	return NewIter(ReaderAsRunesIterGen(src))
}

// OfStringAsRunes constructs an Iter[rune] that iterates runes of a string.
//
// See SliceIterGen.
func OfStringAsRunes(src string) Iter[rune] {
	return NewIter(SliceIterGen([]rune(src)))
}

// OfReaderAsLines constructs an Iter[string] that iterates the UTF-8 lines of a Reader.
//
// See ReaderAsLinesIterGen.
func OfReaderAsLines(src io.Reader) Iter[string] {
	return NewIter(ReaderAsLinesIterGen(src))
}

// OfStringAsLines constructs an Iter[rune] that iterates lines of a string.
//
// See ReaderAsLinesIterGen.
func OfStringAsLines(src string) Iter[string] {
	return NewIter(ReaderAsLinesIterGen(strings.NewReader(src)))
}

// Concatenate any number of Iter[T] into a single Iter[T] that iterates all the elements of each Iter[T], until the
// last element of the last iterator has been returned.
func Concat[T any](iters ...Iter[T]) Iter[T] {
	return NewIter(ConcatIterGen(iters))
}

// ==== IterImpl Methods

// Next returns (true, nil) if there is another item to be read by Value.
// When Next returns (zero value, EOI), further calls return (zero value, EOI).
func (it *IterImpl[T]) Next() (T, error) {
	// Check buffer for values placed by Unread
	if len(it.buffer) > 0 {
		val := it.buffer[len(it.buffer)-1]
		it.buffer = it.buffer[0 : len(it.buffer)-1]
		return val, nil
	}

	// Check if we still may have values to acquire via iterating function
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

// Unread adds the given value to an internal buffer, to be returned by Next in reverse order
func (it *IterImpl[T]) Unread(val T) {
	it.buffer = append(it.buffer, val)
}

// ==== Operations on an Iter

// Maybe converts the result of Next into a Tuple2[T, error] to represent the result as a single type.
func Maybe[T any](it Iter[T]) tuple.Two[T, error] {
	return tuple.Of2Error(it.Next())
}

// SetError sets a particular error to occur instead of the first non-nil error the given iterator returns.
func SetError[T any](it Iter[T], err error) Iter[T] {
	return NewIter[T](func() (T, error) {
		if v, e := it.Next(); e == nil {
			return v, e
		} else {
			var zv T
			return zv, err
		}
	})
}
