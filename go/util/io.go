package util

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/go/funcs"
)

// ErrorReader is an io.Reader that returns a non-eof error after some input bytes have been read.
// This is mostly useful for testing, where the input bytes can get the code being tested to reach a desired state first,
// then verify that code correctly handles non-eof errors.
//
// When the error is returned, the number of bytes is always 0.
// If there is at least one byte remaining to read and a call to read asks for more bytes than remain, the result of the
// call is (num remaining bytes, nil), and the next call to read returns (0, error).
// If the input is an empty slice, the first call to read returns (0, error).
//
// Once all input bytes have been read and the error has been returned, any further calls to read will continue to
// return (0, error).
type ErrorReader struct {
	input []byte
	pos   int
	err   error
}

// NewErrorReader constructs an ErrorReader from a set of input bytes and an error
func NewErrorReader(input []byte, err error) *ErrorReader {
	return &ErrorReader{
		input: input,
		pos:   0,
		err:   err,
	}
}

// Read is the io.Reader method, that eventually returns with the provided error
func (r *ErrorReader) Read(p []byte) (int, error) {
	// If we have no bytes left, return (0, error)
	if r.pos >= len(r.input) {
		return 0, r.err
	}

	// Return the number of bytes asked for, or what we have remaining, whichever is less
	numBytes := funcs.MinOrdered(len(r.input), len(p))

	copy(p[:numBytes], r.input[r.pos:r.pos+numBytes])
	r.pos += numBytes
	return numBytes, nil
}

// ErrorWriter is the Writer analog to ErrorReader, it returns a non-eof error after some output bytes have been written.
// A preallocated slice of bytes tracks the bytes written, so the caller can compare.
type ErrorWriter struct {
	count  int
	output []byte
	err    error
}

// NewErrorWriter constructs an ErrorWriter from a count of output bytes to allow and an error
func NewErrorWriter(count int, err error) *ErrorWriter {
	return &ErrorWriter{count, make([]byte, 0, funcs.MaxOrdered(0, count)), err}
}

// Write is the io.Writer method, that eventually returns with the provided error
func (w *ErrorWriter) Write(p []byte) (int, error) {
	// If we have no bytes left to allow writing, return (0, error)
	if len(w.output) == w.count {
		return 0, w.err
	}

	// If count < 0, return (0, nil) to simulate no bytes read and no error
	if w.count < 0 {
		return 0, nil
	}

	// Return the number of bytes copied, or what we have remaining, whichever is less
	numBytes := funcs.MinOrdered(len(p), w.count-len(w.output))

	w.output = append(w.output, p[:numBytes]...)
	return numBytes, nil
}

// Output provides the output written to the write so far, for comparison
func (w *ErrorWriter) Output() []byte {
	return w.output
}
