package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"io"
)

// ==== Constants

var (
	errNewIterNeedsIterator = fmt.Errorf("NewIter requires a non-nil iterating function")
	errValueExpected        = fmt.Errorf("Value has to be called after Next")
	errDone                 = fmt.Errorf("Next cannot be called again after returning false without calling Unread")
	errNextExpected         = fmt.Errorf("Next has to be called before Value")
)

// ==== Types

// Iter iterates elements of type T
// Next must be called before Value
// Unread can be called any number of times to build up an unbounded buffer of values to iterate
// Taken together, call sequence is (Unread*, Next, Unread*, Value)*
//
// When Next is called, the buffer provided by Unread is consulted first, and if it is empty, then the iterating
// function is consulted.
// Note that Unread values are iterated in reverse order: Unreading (1, 2, 3) iterates (3, 2, 1).
//
// Panics if:
// - Next is called a second time before Value
// - Next returns false and is called a second time without an Unread call
type Iter[T any] struct {
	buffer     []T
	iterFn     func() (T, bool)
	nextCalled bool
	value      T
}

// ==== Construct

// NewIter constructs an Iter[T] from an iterating function that returns (T, bool).
// The function must return (nextItem, true) for every item available to iterate, then return (invalid, false) on the
// next call after the last item, where invalid is any value of type T.
// Once the function returns a false bool value, it will never be called again.
// Panics if iterFn is nil.
func NewIter[T any](iterFn func() (T, bool)) *Iter[T] {
	if iterFn == nil {
		panic(errNewIterNeedsIterator)
	}

	return &Iter[T]{iterFn: iterFn}
}

// Of constructs an Iter[T] that iterates the items passed.
func Of[T any](items ...T) *Iter[T] {
	return NewIter[T](SliceIterGen[T](items))
}

// OfEmpty constructs an Iter[T] that iterates no values.
func OfEmpty[T any]() *Iter[T] {
	return NewIter[T](NoValueIterGen[T]())
}

// OfOne constructs an Iter[T] that iterates a single value.
func OfOne[T any](item T) *Iter[T] {
	return NewIter[T](SingleValueIterGen[T](item))
}

// Of constructs an Iter[KeyValue[K, V]] that iterates the items passed.
func OfMap[K comparable, V any](items map[K]V) *Iter[KeyValue[K, V]] {
	return NewIter[KeyValue[K, V]](MapIterGen[K, V](items))
}

// OfReader constructs an Iter[byte] that iterates the bytes of a Reader.
// See ReaderIterGen for details.
func OfReader(src io.Reader) *Iter[byte] {
	return NewIter[byte](ReaderIterGen(src))
}

// OfReaderAsRunes constructs an Iter[rune] that iterates the UTF-8 runes of a Reader.
// See ReaderAsRunesIterGen for details.
func OfReaderAsRunes(src io.Reader) *Iter[rune] {
	return NewIter[rune](ReaderAsRunesIterGen(src))
}

// OfReaderAsLines constructs an Iter[string] that iterates the UTF-8 lines of a Reader.
// See ReaderAsLinesIterGen for details.
func OfReaderAsLines(src io.Reader) *Iter[string] {
	return NewIter[string](ReaderAsLinesIterGen(src))
}

// Concat
func Concat[T any](iters ...*Iter[T]) *Iter[T] {
	var (
		i    int
		iter *Iter[T]
	)

	return NewIter[T](func() (T, bool) {
		for {
			if i == len(iters) {
				var zv T
				return zv, false
			}

			if iter == nil {
				iter = iters[i]
			}

			if (iter != nil) && iter.Next() {
				return iter.Value(), true
			}

			iter = nil
			i++
		}
	})
}

// ==== Methods

// Next returns true if there is another item to be read by Value.
// Panics if:
// - Called more than once before calling Value
// - Called again after returning false without calling Unread first
func (it *Iter[T]) Next() bool {
	// Die if Next called twice before Value
	if it.nextCalled {
		panic(errValueExpected)
	}

	it.nextCalled = true

	// Check buffer before consulting iterating function in case items have been unread
	if l := len(it.buffer); l > 0 {
		it.value = it.buffer[l-1]
		it.buffer = it.buffer[:l-1]
		return true
	}

	// If the iterating func is nil, we already returned false on the previous call
	if it.iterFn == nil {
		panic(errDone)
	}

	// Try to get next item from iterating function
	if value, haveIt := it.iterFn(); haveIt {
		// If we have it, keep the value for call to Value() and return true
		it.value = value
		return true
	}

	// First call with no more items, mark as iterated
	it.iterFn = nil
	return false
}

// Value returns value found by last call to Next.
// Panics if called before Next
func (it *Iter[T]) Value() T {
	// Die if Value called twice before Next
	if !it.nextCalled {
		panic(errNextExpected)
	}

	it.nextCalled = false

	// Return value
	return it.value
}

// Unread can be called at any time to add a value to a buffer that has to be exhausted before any further calls are
// made to the iterating function
func (it *Iter[T]) Unread(value T) {
	it.buffer = append(it.buffer, value)
}
