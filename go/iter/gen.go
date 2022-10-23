package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/bantling/micro/go/util"
)

// ==== Constants

// Error constants
var (
	InvalidUTF8EncodingError = fmt.Errorf("Invalid UTF 8 encoding")
)

// Internal constants
var (
	zeroUTF8Buffer = []byte{0, 0, 0, 0}
)

// ==== Iterating function generators

// SliceIterGen generates an iterating function for a slice of type T
// First len(slc) calls to iterating function return (slc element, true)
// All remaining calls return (T zero value, false)
func SliceIterGen[T any](slc []T) func() (T, bool) {
	// Simple, just track index on each call
	var idx int

	return func() (value T, haveIt bool) {
		if haveIt = idx < len(slc); haveIt {
			value = slc[idx]
			idx++
			return
		}

		// Once idx = len(slc), all further calls will land here
		return
	}
}

// MapIterGen generates an iterating function for a map[K]V
// First len(m) calls to iterating function return (KeyValue[K, V]{m key, m value}, true)
// All remaining calls return (KeyValue[K, V] zero value, false)
func MapIterGen[K comparable, V any](m map[K]V) func() (util.KeyValue[K, V], bool) {
	// Unlike a slice, we don't know the set of indexes ahead of time
	// Use reflection.Value.MapIter to iterate the keys via a stateful object that tracks the progress of key iteration internally
	// We could use a go routine that writes a key/value pair to a channel, but that would cause a memory leak if map is not fully iterated

	var (
		mi   = reflect.ValueOf(m).MapRange()
		zkv  util.KeyValue[K, V]
		done bool
	)

	return func() (util.KeyValue[K, V], bool) {
		if done {
			return zkv, false
		}
		done = !mi.Next()
		if done {
			return zkv, false
		}
		return util.KVOf(mi.Key().Interface().(K), mi.Value().Interface().(V)), true
	}
}

// NoValueIterGen generates an iterating function that has no values
func NoValueIterGen[T any]() func() (T, bool) {
	var zv T

	return func() (T, bool) {
		return zv, false
	}
}

// SingleValueIterGen generates an iterating function that has one value
func SingleValueIterGen[T any](value T) func() (T, bool) {
	var (
		zv   T
		done bool
	)

	return func() (T, bool) {
		if done {
			return zv, false
		}

		done = true
		return value, true
	}
}

// InfiniteIterGen generates an iterative function based on an iterative function and zero or more initial values.
// The initial values are handled as follows:
//   - zero initial values: the zero value of T is used as the seed value
//   - one initial values: the value given is used as the seed value
//   - multiple initial values: the first n-1 values are returned from the first n-1 calls to the generated function,
//     and the last value is the seed value
//
// The seed value is used as the argument to the first call of the given function.
// The generated values are the first n-1 initialValues followed by the inifinite series
// f(seed), f(f(seed)), f(f(f(seed))), ...
func InfiniteIterGen[T any](iterative func(T) T, initialValues ...T) func() (T, bool) {
	var (
		lastIndex     = len(initialValues) - 1
		literalValues []T
		accumulator   T // start with zero value in case no seed provided
	)

	// Do we have any initial values?
	if lastIndex >= 0 {
		// literal values to return are all but the last initial value
		literalValues = initialValues[:lastIndex]
		// accumulator is last initial value, which is the seed for first call to iterative func
		accumulator = initialValues[lastIndex]
	}

	return func() (result T, haveIt bool) {
		// Inifinite series always have a value to return
		haveIt = true

		// Do we still have literal values left to return?
		if l := len(literalValues); l > 0 {
			// Return literal values in order provided
			result = literalValues[0]

			// Are there more literal values after this one?
			if l > 1 {
				// We have more literals, shorten slice to all but value we're returning
				literalValues = literalValues[1:]
			} else {
				// No more literals, nullify slice so the memory can be freed
				literalValues = nil
			}

			return
		}

		// No literal values left, execute iterative func with accumulator (could be seed value) to get next accumulator
		accumulator = iterative(accumulator)

		// Return next accumulator
		result = accumulator
		return
	}
}

// FibonnaciIterGen generates an iterating function that iterates the Fibonacci series 1, 1, 2, 3, 5, 8, 13, ...
func FibonnaciIterGen() func() (int, bool) {
	// The value returned two calls ago (initially zero)
	var prev2 int

	return InfiniteIterGen(
		// The function actually returns 1, 2, 3, 5, 8, 13, ... - it is missing the leading 1 value
		// This is the easiest way to do the math correctly without futzing around with special initial cases
		func(prev1 int) int {
			// prev 1 is the value returned from the last call
			next := prev2 + prev1
			prev2 = prev1

			return next
		},
		1, // This initial value provides the missing leading 1 value, it is returned without calling the above func
		1, // The seed value for the first call to above func
	)
}

// ReaderIterGen generates an iterating function that iterates all the bytes of an io.Reader
func ReaderIterGen(src io.Reader) func() (byte, bool) {
	var (
		done = src == nil
		buf  = make([]byte, 1)
	)

	return func() (byte, bool) {
		if done {
			return 0, false
		}

		if _, err := src.Read(buf); err != nil {
			if err != io.EOF {
				panic(err)
			}

			return 0, false
		}

		return buf[0], true
	}
}

// ReaderAsRunesIterGen generates an iterating function that iterates all the UTF-8 runes of an io.Reader
func ReaderAsRunesIterGen(src io.Reader) func() (rune, bool) {
	// UTF-8 requires at most 4 bytes for a code point
	var (
		done   = src == nil
		buf    = make([]byte, 4)
		bufPos int
	)

	return func() (rune, bool) {
		if done {
			return 0, false
		}

		// Read next up to 4 bytes from reader into subslice of buffer, after any remaining bytes from last read
		_, err := src.Read(buf[bufPos:])
		if (err != nil) && (err != io.EOF) {
			panic(err)
		}

		// No more to read if first buf pos is 0 and EOF
		if done = (buf[0] == 0) && (err == io.EOF); done {
			return 0, false
		}

		// Decode up to 4 bytes for next code point
		r, rl := utf8.DecodeRune(buf)
		if r == utf8.RuneError {
			panic(InvalidUTF8EncodingError)
		}

		// Shift any remaining unused bytes back to the beginning of the buffer
		copy(buf, buf[rl:])

		// Next time read up to as many bytes as were shifted from source, overwriting remaining bytes
		bufPos = 4 - rl

		// Clear out the unused bytes at the end, in case we don't have enough bytes left to fill them
		copy(buf[bufPos:], zeroUTF8Buffer)

		return r, true
	}
}

// StringAsRunesIterGen generates an iterating function that iterates the runes of a string
func StringAsRunesIterGen(src string) func() (rune, bool) {
	return SliceIterGen([]rune(src))
}

// readLines is common functionality for ReaderAsLinesIterGen and StringAsLinesIterGen
func readLines(it func() (rune, bool)) func() (string, bool) {
	var (
		str    strings.Builder
		lastCR bool
	)

	return func() (string, bool) {
		str.Reset()

		for {
			codePoint, haveIt := it()

			if !haveIt {
				if str.Len() > 0 {
					return str.String(), true
				}

				return "", false
			}

			if codePoint == '\r' {
				lastCR = true
				return str.String(), true
			}

			if codePoint == '\n' {
				if lastCR {
					lastCR = false
					continue
				}

				return str.String(), true
			}

			str.WriteRune(codePoint)
		}
	}
}

// ReaderAsLinesIterGen generates an iterating function that iterates all the UTF-8 lines of an io.Reader
func ReaderAsLinesIterGen(src io.Reader) func() (string, bool) {
	// Use ReaderAsRunesIterGen to read individual runes until a line is read
	return readLines(ReaderAsRunesIterGen(src))
}

// StringAsLinesIterGen generates an iterating function that iterates all the UTF-8 lines of a String
func StringAsLinesIterGen(src string) func() (string, bool) {
	// Use StringAsRunesIterGen to read individual runes until a line is read
	return readLines(StringAsRunesIterGen(src))
}
