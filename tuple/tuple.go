package tuple

// SPDX-License-Identifier: Apache-2.0

// ==== Types

// Two is a tuple of two values for cases where it isn't worthwhile to declare your own struct.
// Particularly useful for functions that need to return (T, error), they can instead return a single value Two[T, error]
type Two[T, U any] struct {
	T T
	U U
}

// Three is analogous to Two, but with three values, where the last value might be an error
type Three[T, U, V any] struct {
	T T
	U U
	V V
}

// Four is analogous to Two, but with four values, where the last value might be an error
type Four[T, U, V, W any] struct {
	T T
	U U
	V V
	W W
}

// ==== Constructors

func Of2[T, U any](t T, u U) Two[T, U] {
	return Two[T, U]{t, u}
}

func Of2Same[T any](t T, u T) Two[T, T] {
	return Two[T, T]{t, u}
}

func Of3[T, U, V any](t T, u U, v V) Three[T, U, V] {
	return Three[T, U, V]{t, u, v}
}

func Of3Same[T any](t T, u T, v T) Three[T, T, T] {
	return Three[T, T, T]{t, u, v}
}

func Of4[T, U, V, W any](t T, u U, v V, w W) Four[T, U, V, W] {
	return Four[T, U, V, W]{t, u, v, w}
}

func Of4Same[T any](t T, u T, v T, w T) Four[T, T, T, T] {
	return Four[T, T, T, T]{t, u, v, w}
}

// ==== Methods

func (t Two[T, U]) Values() (T, U) {
	return t.T, t.U
}

func (t Three[T, U, V]) Values() (T, U, V) {
	return t.T, t.U, t.V
}

func (t Four[T, U, V, W]) Values() (T, U, V, W) {
	return t.T, t.U, t.V, t.W
}
