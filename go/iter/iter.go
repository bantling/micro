// Package iter provides iterators
// SPDX-License-Identifier: Apache-2.0
package iter

import (
	"io"
	"reflect"
)

const (
	READER_BUF_SIZE int = 4 * 1024
)

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

// KeyValue is a struct to hold a single key/pair for a map[K]V entry
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

// MapIterGen generates an iterating function for a map[K]V
// First len(m) calls to iterating function return (KeyValue[K, V]{m key, m value}, true)
// All remaining calls return (KeyValue[K, V] zero value, false)
func MapIterGen[K comparable, V any](m map[K]V) func() (KeyValue[K, V], bool) {
	// Unlike a slice, we don't know the set of indexes ahead of time
	// Use reflection.Value.MapIter to iterate the keys via a stateful object that tracks the progress of key iteration internally
	// We could use a go routine that writes a key/value pair to a channel, but that would cause a memory leak if map is not fully iterated

	var (
		mi   = reflect.ValueOf(m).MapRange()
		zkv  KeyValue[K, V]
		done bool
	)

	return func() (KeyValue[K, V], bool) {
		if done {
			return zkv, false
		}
		done = !mi.Next()
		if done {
			return zkv, false
		}
		return KeyValue[K, V]{mi.Key().Interface().(K), mi.Value().Interface().(V)}, true
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

// ReaderIterGen generates an iterating function that iterates all the bytes of an io.Reader
func ReaderIterGen(src io.Reader) func() (byte, bool) {
	var (
		buf     = make([]byte, READER_BUF_SIZE)
		bufSize int
		pos     int
		done    bool
	)

	return func() (byte, bool) {
		if done {
			return 0, false
		}

		// Read next buffer of data
		if pos == bufSize {
			if bufSize, err := src.Read(buf); err != nil {
				// Die on non-EOF error
				if err != io.EOF {
					panic(err)
				}

				// May get (bufSize > 0, EOF), where next call to src.Read will return (bufSize = 0, EOF)
				if bufSize == 0 {
					done = true
					return 0, false
				}
			}

			pos = 0
		}

		// Must have at least one available byte in buffer
		val := buf[pos]
		pos++
		return val, true
	}
}
