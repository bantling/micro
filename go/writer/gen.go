package writer

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"io"
	"runtime"

	"github.com/bantling/micro/go/funcs"
	"github.com/bantling/micro/go/util"
)

// Constants
var (
	errIOByteWriterMsg   = "The byte 0x%x could not be written"
	errIORuneWriterMsg   = "The rune \\U%08x could not be written"
	errIOStringWriterMsg = "Only the first %d bytes were written of a string of %d runes that is encoded as %d UTF-8 bytes"

	osEOLSequence = funcs.Ternary(runtime.GOOS == "windows", "\r\n", "\n")
)

// SliceWritergen generates a writing function for a slice of type T.
// A pointer to the slice is required so that if the slice has to be reallocated, the address of the caller slice can be changed.
func SliceWriterGen[T any](slc *[]T) func(T) error {
	return func(val T) error {
		*slc = append(*slc, val)
		return nil
	}
}

// MapWriterGen generates a writing function for a map[K]V.
// A pointer to the map is not required, as the map internally reallocates as needed, so that the caller map address never changes.
func MapWriterGen[K comparable, V any](m map[K]V) func(util.Tuple2[K, V]) error {
	return func(val util.Tuple2[K, V]) error {
		m[val.T] = val.U
		return nil
	}
}

// IOWriterGen generates a writing function for writing bytes to an io.Writer.
// Returns an error if the byte is not written:
// - If the Writer returns an error, it is returned as is
// - If the Writer wrote 0 bytes with no error, an custom error containing the hex value of the byte is returned
func IOWriterGen(dst io.Writer) func(byte) error {
	var p = []byte{0}

	return func(val byte) error {
		p[0] = val
		if n, err := dst.Write(p); err != nil {
			return err
		} else if n == 0 {
			return fmt.Errorf(errIOByteWriterMsg, val)
		}

		return nil
	}
}

// IOWriterAsRunesGen generates a writing function for writing runes to an io.Writer.
// Returns an error if the rune is not written:
// - If the Writer returns an error, it is returned as is
// - If the Writer wrote fewer bytes than the UTF-8 encoding requires, a custom error containing the hex value of the rune is returned
func IOWriterAsRunesGen(dst io.Writer) func(rune) error {
	return func(val rune) error {
		// No actual encoding function provided in Go unicode/utf8 package as none is needed.
		// Simply convert rune > string > []byte
		p := []byte(string(val))
		if n, err := dst.Write(p); err != nil {
			return err
		} else if n < len(p) {
			return fmt.Errorf(errIORuneWriterMsg, val)
		}

		return nil
	}
}

// IOWriterAsStringsGen generates a writing function for writing strings to an io.Writer.
// Returns an error if any rune is not written:
// - If the Writer returns an error, it is returned as is
// - If the Writer wrote fewer bytes than the UTF-8 encoding requires, a custom error containing the number of bytes written is returned
func IOWriterAsStringsGen(dst io.Writer) func(string) error {
	return func(val string) error {
		// No actual encoding function provided in Go unicode/utf8 package as none is needed.
		// Simply convert string > []byte
		p := []byte(val)
		if n, err := dst.Write(p); err != nil {
			return err
		} else if n < len(p) {
			return fmt.Errorf(errIOStringWriterMsg, n, len(val), len(p))
		}

		return nil
	}
}

// IOWriterAsLinesGen generates a writing function for writing lines to an io.Writer.
// The idea is that each string is one line, but does not end in any line ending sequence.
// The OS line ending sequence is used.
// Returns an error if any rune is not written:
// - If the Writer returns an error, it is returned as is
// - If the Writer wrote fewer bytes than the UTF-8 encoding requires, a custom error containing the number of bytes written is returned
func IOWriterAsLinesGen(dst io.Writer) func(string) error {
	return func(val string) error {
		// No actual encoding function provided in Go unicode/utf8 package as none is needed.
		// Simply add OS line ending sequence and convert string > []byte
		var (
			valEOL = val + osEOLSequence
			p      = []byte(valEOL)
		)
		if n, err := dst.Write(p); err != nil {
			return err
		} else if n < len(p) {
			return fmt.Errorf(errIOStringWriterMsg, n, len(valEOL), len(p))
		}

		return nil
	}
}
