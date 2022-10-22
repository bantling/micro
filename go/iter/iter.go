package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"io"
	"strings"
)

// ==== Constants

var (
	errNewIterNeedsIterator = fmt.Errorf("NewIter requires a non-nil iterating function")
	errValueExpected        = fmt.Errorf("Value has to be called after Next")
	errNextExpected         = fmt.Errorf("Next has to be called before Value")
	errNoMoreValues         = fmt.Errorf("Value cannot be called after Next returns false")
)

// ==== Types

// Iter defines iteration operations for type T, which work as follows:
// - Next and Value must be called in interleaving pairs
// - Unread can be called any number of times to build up an unbounded buffer of values to iterate
// - An underlying iteration function that returns (T, bool) provides the values to iterate
// - NextValue combines Next and Value in a single call, returning value and bool flag to tell if value is valid
// - Must also combines Next and Value in one call, but only returns the value, and panics if a value does not exist to return
//
// Taken together, in the call sequence Unread*, (Next, Value)*, the following happens:
// - Unread calls build a buffer
// - Next/Value pairs return the values in the buffer provided by Unread until it is empty, then the function is consulted.
// - Note that Unread values are iterated in reverse order: Unreading (1, 2, 3) iterates (3, 2, 1).
// - Unread calls can be made between Next and Value calls.
//
// Panics if:
// - Next is called a second time before Value
// - Next returns false and is called a second time without an Unread call
type Iter[T any] interface {
	Next() bool
	Value() T
	Unread(value T)
	NextValue() (T, bool)
	Must() T
}

// IterImpl is the base implementation of Iter[T]
type IterImpl[T any] struct {
	buffer     []T
	iterFn     func() (T, bool)
	nextCalled bool
	haveValue  bool
	value      T
}

// IOByteIterImpl is a byte override of Iter[rune] that alters the base implementation in a single respect:
// Unreading a zero value is ignored, as that means unreading eof, which causes Next/Value to return an actual zero.
// By ignoring unreads of 0, Next will return false, and the 0 does not get returned as a value.
type IOByteIterImpl struct {
	IterImpl[byte]
}

// IORuneIterImpl is a rune override of Iter[rune] that alters the base implementation in a single respect:
// Unreading a zero value is ignored, as that means unreading eof, which causes Next/Value to return an actual zero.
// By ignoring unreads of 0, Next will return false, and the 0 does not get returned as a value.
type IORuneIterImpl struct {
	IterImpl[rune]
}

// ==== Construct

// NewIter constructs an Iter[T] from an iterating function that returns (T, bool).
// The function must return (nextItem, true) for every item available to iterate, then return (invalid, false) on the
// next call after the last item, where invalid is any value of type T.
// Once the function returns a false bool value, it will never be called again.
// Panics if iterFn is nil.
func NewIter[T any](iterFn func() (T, bool)) Iter[T] {
	if iterFn == nil {
		panic(errNewIterNeedsIterator)
	}

	return &IterImpl[T]{iterFn: iterFn}
}

// Of constructs an Iter[T] that iterates the items passed.
func Of[T any](items ...T) Iter[T] {
	return NewIter[T](SliceIterGen[T](items))
}

// OfEmpty constructs an Iter[T] that iterates no values.
func OfEmpty[T any]() Iter[T] {
	return NewIter[T](NoValueIterGen[T]())
}

// OfOne constructs an Iter[T] that iterates a single value.
func OfOne[T any](item T) Iter[T] {
	return NewIter[T](SingleValueIterGen[T](item))
}

// Of constructs an Iter[KeyValue[K, V]] that iterates the items passed.
func OfMap[K comparable, V any](items map[K]V) Iter[KeyValue[K, V]] {
	return NewIter[KeyValue[K, V]](MapIterGen[K, V](items))
}

// OfReader constructs an Iter[byte] that iterates the bytes of a Reader.
// See ReaderIterGen and IOByteIterImpl for details.
func OfReader(src io.Reader) Iter[byte] {
	return &IOByteIterImpl{IterImpl: IterImpl[byte]{iterFn: ReaderIterGen(src)}}
}

// OfReaderAsRunes constructs an Iter[rune] that iterates the UTF-8 runes of a Reader.
//
// See ReaderAsRunesIterGen for details.
func OfReaderAsRunes(src io.Reader) Iter[rune] {
	return &IORuneIterImpl{IterImpl: IterImpl[rune]{iterFn: ReaderAsRunesIterGen(src)}}
}

// OfStringAsRunes constructs an Iter[rune] that iterates runes of a string.
func OfStringAsRunes(src string) Iter[rune] {
	// return Of([]rune(src)...)
	return &IORuneIterImpl{IterImpl: IterImpl[rune]{iterFn: SliceIterGen([]rune(src))}}
}

// OfReaderAsLines constructs an Iter[string] that iterates the UTF-8 lines of a Reader.
// See ReaderAsLinesIterGen for details.
func OfReaderAsLines(src io.Reader) Iter[string] {
	return NewIter(ReaderAsLinesIterGen(src))
}

// OfStringAsLines constructs an Iter[rune] that iterates lines of a string.
func OfStringAsLines(src string) Iter[string] {
	return NewIter(ReaderAsLinesIterGen(strings.NewReader(src)))
}

// Concatenate any number of Iter[T] into a single Iter[T] that iterates all the elements of each Iter[T], until the
// last element of the last iterator has been returned.
func Concat[T any](iters ...Iter[T]) Iter[T] {
	var (
		i    int
		iter Iter[T]
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

// ==== IterImpl Methods

// Next returns true if there is another item to be read by Value.
// When Next returns false, Next can be called any number of times, it just continues to return false.
// When Next returns true, a panic occurs if Next is called again before calling Value.
func (it *IterImpl[T]) Next() bool {
	// Assume no value until we know differently
	it.haveValue = false

	// Die if Next called twice before Value, unless prior Next call exhausted iter
	if it.nextCalled && (it.iterFn != nil) {
		panic(errValueExpected)
	}

	it.nextCalled = true

	// Check buffer before consulting iterating function in case items have been unread
	if l := len(it.buffer); l > 0 {
		// Read items from buffer in order they were unread - eg unread(x), unread(y) returns x first, then y
		it.haveValue = true
		it.value = it.buffer[0]
		it.buffer = it.buffer[1:]
		return true
	}

	// If the iterating func is nil, we must have exhausted the func, Unread was called, and the buffer also exhausted
	if it.iterFn == nil {
		return false
	}

	// Try to get next item from iterating function
	if value, haveIt := it.iterFn(); haveIt {
		// If we have it, keep the value for call to Value() and return true
		it.haveValue = true
		it.value = value
		return true
	}

	// First call with no more items, mark as iterated
	it.iterFn = nil
	return false
}

// Value returns value found by last call to Next.
// Panics if called before Next
func (it *IterImpl[T]) Value() T {
	// Die if Value called twice before Next
	if !it.nextCalled {
		panic(errNextExpected)
	}

	if !it.haveValue {
		// Die if next indicated no more values exist
		panic(errNoMoreValues)
	}

	it.nextCalled = false
	it.haveValue = false

	// Return value
	return it.value
}

// Unread can be called any time to add a value to a buffer that has to be exhausted before any further calls are made
// to the iterating function.
func (it *IterImpl[T]) Unread(value T) {
	it.buffer = append(it.buffer, value)
	if it.iterFn == nil {
		it.nextCalled = false
	}
}

// NextValue combines Next and Value together in a single call.
// If there is another value, then (next value, true) is returned, else (zero value, false) is returned.
// NextValue may be called after Next has already returned false without a panic.
func (it *IterImpl[T]) NextValue() (T, bool) {
	if (it.iterFn != nil) && it.Next() {
		return it.Value(), true
	}

	var zv T
	return zv, false
}

// Must combines Next and Value together in a single call.
// If there is another value, then the next value is returned, else a panic occurs.
func (it *IterImpl[T]) Must() T {
	it.Next()
	return it.Value()
}

// ==== IOByteIterImpl Methods

// Don't unread a zero byte - this should only occur if caller mistakenly unreads the byte returned by NextValue when
// the bool is false.
func (it *IOByteIterImpl) Unread(value byte) {
	if value > 0 {
		it.IterImpl.Unread(value)
	}
}

// ==== IORuneIterImpl Methods

// Don't unread a zero rune - this should only occur if caller mistakenly unreads the rune returned by NextValue when
// the bool is false.
func (it *IORuneIterImpl) Unread(value rune) {
	if value > 0 {
		it.IterImpl.Unread(value)
	}
}
