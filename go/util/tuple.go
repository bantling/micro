package util

// SPDX-License-Identifier: Apache-2.0

// ==== Types

// Tuple2 is a tuple of two values for cases where it isn't worthwhile to declare your own struct.
// Particularly useful for functions that need to return (T, error), they can instead return a single value Tuple2[T, error]
type Tuple2[T, U any] struct {
	T T
	U U
}

// Tuple3 is analogous to Tuple2, but with three values, where the last value might be an error
type Tuple3[T, U, V any] struct {
	T T
	U U
	V V
}

// Tuple4 is analogous to Tuple2, but with four values, where the last value might be an error
type Tuple4[T, U, V, W any] struct {
	T T
	U U
	V V
	W W
}

// ==== Constructors

func Of2[T, U any](t T, u U) Tuple2[T, U] {
	return Tuple2[T, U]{t, u}
}

func Of2Same[T any](t T, u T) Tuple2[T, T] {
	return Tuple2[T, T]{t, u}
}

func Of2Error[T any](t T, err error) Tuple2[T, error] {
	return Tuple2[T, error]{t, err}
}

func Of3[T, U, V any](t T, u U, v V) Tuple3[T, U, V] {
	return Tuple3[T, U, V]{t, u, v}
}

func Of3Same[T any](t T, u T, v T) Tuple3[T, T, T] {
	return Tuple3[T, T, T]{t, u, v}
}

func Of3Error[T, U any](t T, u U, err error) Tuple3[T, U, error] {
	return Tuple3[T, U, error]{t, u, err}
}

func Of4[T, U, V, W any](t T, u U, v V, w W) Tuple4[T, U, V, W] {
	return Tuple4[T, U, V, W]{t, u, v, w}
}

func Of4Same[T any](t T, u T, v T, w T) Tuple4[T, T, T, T] {
	return Tuple4[T, T, T, T]{t, u, v, w}
}

func Of4Error[T, U, V any](t T, u U, v V, err error) Tuple4[T, U, V, error] {
	return Tuple4[T, U, V, error]{t, u, v, err}
}

// ==== Methods

func (t Tuple2[T, U]) Values() (T, U) {
	return t.T, t.U
}

func (t Tuple3[T, U, V]) Values() (T, U, V) {
	return t.T, t.U, t.V
}

func (t Tuple4[T, U, V, W]) Values() (T, U, V, W) {
	return t.T, t.U, t.V, t.W
}
