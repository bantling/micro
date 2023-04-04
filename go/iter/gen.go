package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/bantling/micro/go/funcs"
	"github.com/bantling/micro/go/tuple"
)

// ==== Constants

// Error constants
var (
	InvalidUTF8EncodingError = fmt.Errorf("Invalid UTF 8 encoding")
)

// ==== Iterating function generators

// SliceIterGen generates an iterating function for a slice of type T
// First len(slc) calls to iterating function return (slc element, nil)
// All remaining calls return (T zero value, EOI)
func SliceIterGen[T any](slc []T) func() (T, error) {
	// Simple, just track index on each call
	var idx int

	return func() (T, error) {
		if idx < len(slc) {
			value := slc[idx]
			idx++
			return value, nil
		}

		// Once idx = len(slc), all further calls will land here
		var zv T
		return zv, EOI
	}
}

// MapIterGen generates an iterating function for a map[K]V
// First len(m) calls to iterating function return (tuple.Two[K, V]{m key, m value}, nil)
// All remaining calls return (tuple.Two[K, V] zero value, EOI)
func MapIterGen[K comparable, V any](m map[K]V) func() (tuple.Two[K, V], error) {
	// Unlike a slice, we don't know the set of indexes ahead of time
	// Use reflection.Value.MapIter to iterate the keys via a stateful object that tracks the progress of key iteration internally
	// We could use a go routine that writes a key/value pair to a channel, but that would cause a memory leak if map is not fully iterated

	var (
		mi   = reflect.ValueOf(m).MapRange()
		zv   tuple.Two[K, V]
		done bool
	)

	return func() (tuple.Two[K, V], error) {
		if done {
			return zv, EOI
		}

		done = !mi.Next()
		if done {
			return zv, EOI
		}

		return tuple.Of2(mi.Key().Interface().(K), mi.Value().Interface().(V)), nil
	}
}

// NoValueIterGen generates an iterating function that has no values.
// Always returns (zero value, EOI)
func NoValueIterGen[T any]() func() (T, error) {
	var zv T

	return func() (T, error) {
		return zv, EOI
	}
}

// SingleValueIterGen generates an iterating function that has one value
func SingleValueIterGen[T any](value T) func() (T, error) {
	var (
		zv   T
		done bool
	)

	return func() (T, error) {
		if done {
			return zv, EOI
		}

		done = true
		return value, nil
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
func InfiniteIterGen[T any](iterative func(T) T, initialValues ...T) func() (T, error) {
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

	return func() (T, error) {
		// Inifinite series always have a value to return
		var result T

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

			return result, nil
		}

		// No literal values left, execute iterative func with accumulator (could be seed value) to get next accumulator
		accumulator = iterative(accumulator)

		// Return next accumulator
		result = accumulator
		return result, nil
	}
}

// FibonnaciIterGen generates an iterating function that iterates the Fibonacci series 1, 1, 2, 3, 5, 8, 13, ...
func FibonnaciIterGen() func() (int, error) {
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

// ReaderIterGen generates an iterating function that iterates all the bytes of an io.Reader.
// If the reader returns an EOF, it is translated to an EOI, any other error is returned as is.
// If the iter is called again after returning a non-nil error, it returns (0, same error).
func ReaderIterGen(src io.Reader) func() (byte, error) {
	var (
		done = src == nil
		err  = EOI
		buf  = make([]byte, 1)
	)

	return func() (byte, error) {
		if done {
			return 0, err
		}

		if _, err := src.Read(buf); err != nil {
			done = true
			err = funcs.Ternary(err == io.EOF, EOI, err)
			return 0, err
		}

		return buf[0], nil
	}
}

// ReaderAsRunesIterGen generates an iterating function that iterates all the UTF-8 runes of an io.Reader.
// Up to four UTF-8 bytes are read to produce a single rune.
// If the reader returns an EOF, it is translated to an EOI, any other error is returned as is.
// If the iter is called again after returning a non-nil error, it returns (0, same error).
func ReaderAsRunesIterGen(src io.Reader) func() (rune, error) {
	// UTF-8 encodes the bytes as follows:
	//
	// +==============================================================================================+
	// | First code point | Last code point | Byte 1   | Byte 2   | Byte 3   | Byte 4   | Code points |
	// | U+0000           | U+007F          | 0xxxxxxx |          |          |          | 128         |
	// | U+0080           | U+07FF          | 110xxxxx | 10xxxxxx |          |          | 1920        |
	// | U+0800           | U+FFFF          | 1110xxxx | 10xxxxxx | 10xxxxxx |          | 61440       |
	// | U+10000          | U+10FFFF        | 11110xxx | 10xxxxxx | 10xxxxxx | 10xxxxxx | 1048576     |
	// +==============================================================================================+
	//
	// Minimum and maximum values for each byte of encoding
	// Two byte encodings:
	// - range from 0000 0000 1000 0000 thru 0000 0111 1111 1111
	// - encoded as 1100 0aaa bbbb cccc thru 1100 0aaa bbbb cccc
	// +==============================================================================================+
	// | Byte 1                | Byte 2                | Byte 3                | Byte 4               |
	// | 0 0000000 - 0 1111111 |                       |                       |                      |
	// | 110 00000 - 110 11111 | 10 000000 - 10 111111 |                       |                      |
	// | 1110 0000 - 1110 1111 | 10 000000 - 10 111111 | 10 000000 - 10 111111 |                      |
	// | 11110 000 - 11110 100 | 10 000000 - 10 001111 | 10 000000 - 10 111111 | 10 00000 - 10 111111 |
	// +==============================================================================================+
	//
	// See https://en.wikipedia.org/wiki/UTF-8 for further details

	var (
		done  = src == nil
		err   = EOI
		n, eb int
		buf   = make([]byte, 3)
		b     byte
		r     rune
	)

	return func() (rune, error) {
		if done {
			return 0, err
		}

		// Read up to next 4 bytes from reader, depending on bit pattern of first byte
		eb = 0
		if n, err = src.Read(buf[0:1]); (n == 0) || (err != nil) {
			done = true

			// Should be a non-nil error if 0 bytes were returned, but don't assume
			// Translate EOF to EOI
			if ((n == 0) && (err == nil)) || (err == io.EOF) {
				err = EOI
			}

			return 0, err
		}

		// First byte indicates how many more bytes are needed, if any
		b = buf[0]
		switch {
		case (b & 0b1_0000000) == 0b0_0000000:
			// One byte
			return rune(b), nil
		case (b & 0b111_00000) == 0b110_00000:
			// Two bytes = 1 extra byte
			eb = 1
		case (b & 0b1111_0000) == 0b1110_0000:
			// Three bytes = 2 extra bytes
			eb = 2
		case (b & 0b11111_000) == 0b11110_000:
			// Four bytes = 3 extra bytes
			eb = 3
		default:
			// Incorrect leading bit pattern, must be either:
			// - 10_xxxxxx reserved for extra bytes
			// - 11111_xxx that is not used
			done = true
			err = InvalidUTF8EncodingError
			return 0, err
		}

		// Multiple bytes
		if n, err = src.Read(buf[0:eb]); (n != eb) || (err != nil) {
			// Not enough extra bytes exist or some other error
			if (err == nil) || (err == io.EOF) {
				err = InvalidUTF8EncodingError
			}

			return 0, err
		}

		// We have read all required extra bytes
		// Ensure they all start with 10
		for i := 0; i < eb; i++ {
			if buf[i]&0b11_000000 != 0b10_000000 {
				// Invalid extra byte
				err = InvalidUTF8EncodingError
				return 0, err
			}
		}

		// At least two bytes in the encoding
		switch eb {
		case 1: // Two bytes
			r = (rune(b&0b000_11111) << 6) |
				(rune(buf[0] & 0b00_111111))

		case 2: // Three bytes
			r = (rune(b&0b0000_1111) << (6 + 6)) |
				(rune(buf[0]&0b00_111111) << 6) |
				(rune(buf[1] & 0b00_111111))

		default: // Four bytes
			r = (rune(b&0b00000_111) << (6 + 6 + 6)) |
				(rune(buf[0]&0b00_111111) << (6 + 6)) |
				(rune(buf[1]&0b00_111111) << 6) |
				(rune(buf[2] & 0b00_111111))

			// Cannot exceed maximum value
			if r > utf8.MaxRune {
				err = InvalidUTF8EncodingError
				return 0, err
			}
		}

		// Valid multiple bytes
		return r, nil
	}
}

// StringAsRunesIterGen generates an iterating function that iterates the runes of a string.
// See ReaderAsRunesIterGen.
func StringAsRunesIterGen(src string) func() (rune, error) {
	// If the string is invalid UTF8, when converted to a []rune, it will contain a utf8.RuneError value.
	// By using ReaderAsRunesIterGen, that rune will be converted to an InvalidUTF8EncodingError.
	return ReaderAsRunesIterGen(strings.NewReader(src))
}

// readLines is common functionality for ReaderAsLinesIterGen and StringAsLinesIterGen.
// If the reader returns an EOF, it is translated to an EOI, any other error is returned as is.
// If an EOF occurs after some input is buffered, the buffer is returned with a nil error, further calls return 0, EOI.
// If the iter is called again after returning a non-nil error, it returns (0, same error).
func readLines(it func() (rune, error)) func() (string, error) {
	var (
		str    strings.Builder
		lastCR bool
		err    error
	)

	return func() (string, error) {
		if err != nil {
			return "", err
		}

		str.Reset()

		for {
			var codePoint rune
			codePoint, err = it()

			if err != nil {
				if err = funcs.Ternary(err == io.EOF, EOI, err); err != EOI {
					return "", err
				}

				if str.Len() > 0 {
					return str.String(), nil
				}

				return "", err
			}

			if codePoint == '\r' {
				lastCR = true
				return str.String(), nil
			}

			if codePoint == '\n' {
				if lastCR {
					lastCR = false
					continue
				}

				return str.String(), nil
			}

			str.WriteRune(codePoint)
		}
	}
}

// ReaderAsLinesIterGen generates an iterating function that iterates all the UTF-8 lines of an io.Reader
// See readLines.
func ReaderAsLinesIterGen(src io.Reader) func() (string, error) {
	// Use ReaderAsRunesIterGen to read individual runes until a line is read
	return readLines(ReaderAsRunesIterGen(src))
}

// StringAsLinesIterGen generates an iterating function that iterates all the UTF-8 lines of a String
// See readLines.
func StringAsLinesIterGen(src string) func() (string, error) {
	// Use StringAsRunesIterGen to read individual runes until a line is read
	return readLines(StringAsRunesIterGen(src))
}

// ConcatIterGen generates an iterating function that iterates all the values of all the Iters passed.
// If a non-nil non-EOI error is returned from an underlying iter, then (zero value, error) is returned.
// After returning (zero value, non-nil error), all further calls return (zero value, same error).
func ConcatIterGen[T any](src []Iter[T]) func() (T, error) {
	var (
		i    int
		iter Iter[T]
		zv   T
		err  = EOI
	)

	return func() (T, error) {
		for {
			if i == len(src) {
				return zv, err
			}

			if iter == nil {
				iter = src[i]
			}

			if iter != nil {
				var val T
				if val, err = iter.Next(); err == nil {
					return val, nil
				}

				if err != EOI {
					i = len(src)
					return zv, err
				}
			}

			iter = nil
			i++
		}
	}
}
