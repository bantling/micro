package writer

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"io"

	"github.com/bantling/micro/go/util"
)

// ==== constants

var (
	errNewWriterNeedsFn = fmt.Errorf("NewWriter requires a non-nil writing function")
)

// Writer defines a single method for writing values of type T to some destination.
// Panics if a value of type T cannot be written to the destination successfully.
// The error object provided to panic contains whatever object the underlying storage function provides.
type Writer[T any] interface {
	Write(val T)
}

// WriterImpl is the base implementation of Writer[T]
type WriterImpl[T any] struct {
	writerFn func(T) error
}

// ==== Construct

// NewWriter constructs a Writer[T] from a writing function that accepts T and returns error.
// Panics if writerFn is nil.
func NewWriter[T any](writerFn func(T) error) Writer[T] {
	if writerFn == nil {
		panic(errNewWriterNeedsFn)
	}

	return WriterImpl[T]{writerFn: writerFn}
}

// OfSliceWriter returns a Writer[T] that writes to the given slice
func OfSliceWriter[T any](dst *[]T) Writer[T] {
	return NewWriter(SliceWriterGen(dst))
}

// OfMapWriter returns a Writer[KeyValue[K, V]] that writes to the given map
func OfMapWriter[K comparable, V any](dst map[K]V) Writer[util.KeyValue[K, V]] {
	return NewWriter(MapWriterGen(dst))
}

// OfIOWriterAsBytes returns a Writer[byte] that writes bytes to the given io.Writer.
// Panics if writing a byte fails.
func OfIOWriterAsBytes(dst io.Writer) Writer[byte] {
	return NewWriter(IOWriterGen(dst))
}

// OfIOWriterAsRunes returns a Writer[rune] that writes UTF-8 encoded runes to the given io.Writer.
// Panics if writing any byte of a rune fails.
func OfIOWriterAsRunes(dst io.Writer) Writer[rune] {
	return NewWriter(IOWriterAsRunesGen(dst))
}

// OfIOWriterAsStrings returns a Writer[string] that writes UTF-8 encoded strings to the given io.Writer.
// Panics if writing any byte of the string fails.
func OfIOWriterAsStrings(dst io.Writer) Writer[string] {
	return NewWriter(IOWriterAsStringsGen(dst))
}

// OfIOWriterAsLines returns a Writer[string] that writes UTF-8 encoded strings to the given io.Writer.
// Each string is followed by the OS line ending sequence, where Windows uses CRLF and all other systems use LF.
// Panics if writing any byte of the string fails.
func OfIOWriterAsLines(dst io.Writer) Writer[string] {
	return NewWriter(IOWriterAsLinesGen(dst))
}

// ==== WriterImpl method

// Write panics if the value given cannot be written to the destination.
func (w WriterImpl[T]) Write(val T) {
	if err := w.writerFn(val); err != nil {
		panic(err)
	}
}