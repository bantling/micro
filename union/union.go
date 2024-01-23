package union

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	goreflect "reflect"

	"github.com/bantling/micro/funcs"
)

// Errors

var (
	errWhichMsg   = "Member %s is not available"
	errEmptyMaybe = fmt.Errorf("Empty Maybe cannot return a value")
  errPresentMaybe = fmt.Errorf("Present Maybe cannot be overwritten with SetOrError")
)

// ==== Types

// Which describes which member to return
type Which uint

const (
	T Which = iota
	U
	V
	W
)

var (
	whichStr = map[Which]string{
		T: "T",
		U: "U",
		V: "V",
		W: "W",
	}
)

// Stringer for Which
func (w Which) String() string {
	return whichStr[w]
}

// Two is a union of two values
type Two[T, U any] struct {
	t     T
	u     U
	which Which
}

// Three is a union of three values`
type Three[T, U, V any] struct {
	t     T
	u     U
	v     V
	which Which
}

// Four is a union of four values
type Four[T, U, V, W any] struct {
	t     T
	u     U
	v     V
	w     W
	which Which
}

// Maybe is a wrapper around a single value that may or may not exist, as an alternative to a pointer or an Iter
type Maybe[T any] struct {
	v       T
	present bool
}

// Result is a variation of Two where the second value is predefined as an error
type Result[R any] struct {
	r R
	e error
}

// ==== Constructors

// Of2T constructs a Two that holds a T
func Of2T[TT, UU any](t TT) Two[TT, UU] {
	return Two[TT, UU]{t: t, which: T}
}

// Of2U constructs a Two that holds a U
func Of2U[TT, UU any](u UU) Two[TT, UU] {
	return Two[TT, UU]{u: u, which: U}
}

// Of3T constructs a Three that holds a T
func Of3T[TT, UU, VV any](t TT) Three[TT, UU, VV] {
	return Three[TT, UU, VV]{t: t, which: T}
}

// Of3U constructs a Three that holds a U
func Of3U[TT, UU, VV any](u UU) Three[TT, UU, VV] {
	return Three[TT, UU, VV]{u: u, which: U}
}

// Of3V constructs a Three that holds a V
func Of3V[TT, UU, VV any](v VV) Three[TT, UU, VV] {
	return Three[TT, UU, VV]{v: v, which: V}
}

// Of4T constructs a Four that holds a T
func Of4T[TT, UU, VV, WW any](t TT) Four[TT, UU, VV, WW] {
	return Four[TT, UU, VV, WW]{t: t, which: T}
}

// Of4U constructs a Four that holds a U
func Of4U[TT, UU, VV, WW any](u UU) Four[TT, UU, VV, WW] {
	return Four[TT, UU, VV, WW]{u: u, which: U}
}

// Of4V constructs a Four that holds a V
func Of4V[TT, UU, VV, WW any](v VV) Four[TT, UU, VV, WW] {
	return Four[TT, UU, VV, WW]{v: v, which: V}
}

// Of4W constructs a Four that holds a W
func Of4W[TT, UU, VV, WW any](w WW) Four[TT, UU, VV, WW] {
	return Four[TT, UU, VV, WW]{w: w, which: W}
}

// Of constructs a Maybe that holds a T.
// If t is nil, the Maybe is Empty, otherwise it is Present.
func Of[TT any](t TT) Maybe[TT] {
	return Maybe[TT]{v: t, present: !funcs.IsNilValue(t)}
}

// Empty constructs an empty Maybe
func Empty[TT any]() Maybe[TT] {
	return Maybe[TT]{}
}

// OfResult constructs a Result that holds an R
func OfResult[R any](r R) Result[R] {
	return Result[R]{r: r}
}

// OfError constructs a Result that holds an error
// Panics if the error is nil
func OfError[R any](err error) Result[R] {
	if err == nil {
		panic(fmt.Errorf("A Result cannot be set to a nil error"))
	}

	return Result[R]{e: err}
}

// OfResultError constructs a Result from a value of type R and an error
// Panics if the error is non-nil and R is not the zero value
func OfResultError[R any](r R, err error) Result[R] {
	if err != nil {
		// R is not comparable
		if !goreflect.ValueOf(r).IsZero() {
			panic(fmt.Errorf("A Result cannot have both a non-zero R value and a non-nil error"))
		}
	}

	return Result[R]{r: r, e: err}
}

// ==== Helpers

// check if the desired member is selected
func check[R any](asked, have Which, res R) R {
	if asked != have {
		panic(fmt.Errorf(errWhichMsg, asked))
	}

	return res
}

// ==== Two

// Which of Two
func (s Two[TT, UU]) Which() Which {
	return s.which
}

// T of Two
// Panics if which != T
func (s Two[TT, UU]) T() TT {
	return check[TT](T, s.which, s.t)
}

// SetT of Two
func (s *Two[TT, UU]) SetT(t TT) {
	s.t = t
	s.which = T
}

// U of Two
// Panics if which != U
func (s Two[TT, UU]) U() UU {
	return check[UU](U, s.which, s.u)
}

// SetU of Two
func (s *Two[TT, UU]) SetU(u UU) {
	s.u = u
	s.which = U
}

// ==== Three

// Which of Three
func (s Three[TT, UU, VV]) Which() Which {
	return s.which
}

// T of Three
// Panics if which != T
func (s Three[TT, UU, VV]) T() TT {
	return check[TT](T, s.which, s.t)
}

// SetT of Three
func (s *Three[TT, UU, VV]) SetT(t TT) {
	s.t = t
	s.which = T
}

// U of Three
// Panics if which != U
func (s Three[TT, UU, VV]) U() UU {
	return check[UU](U, s.which, s.u)
}

// SetU of Three
func (s *Three[TT, UU, VV]) SetU(u UU) {
	s.u = u
	s.which = U
}

// V of Three
// Panics if which != V
func (s Three[TT, UU, VV]) V() VV {
	return check[VV](V, s.which, s.v)
}

// SetV of Three
func (s *Three[TT, UU, VV]) SetV(v VV) {
	s.v = v
	s.which = V
}

// ==== Four

// Which of Four
func (s Four[TT, UU, VV, WW]) Which() Which {
	return s.which
}

// T of Four
// Panics if which != T
func (s Four[TT, UU, VV, WW]) T() TT {
	return check[TT](T, s.which, s.t)
}

// SetT of Four
func (s *Four[TT, UU, VV, WW]) SetT(t TT) {
	s.t = t
	s.which = T
}

// U of Four
// Panics if which != U
func (s Four[TT, UU, VV, WW]) U() UU {
	return check[UU](U, s.which, s.u)
}

// SetU of Four
func (s *Four[TT, UU, VV, WW]) SetU(u UU) {
	s.u = u
	s.which = U
}

// V of Four
// Panics if which != V
func (s Four[TT, UU, VV, WW]) V() VV {
	return check[VV](V, s.which, s.v)
}

// SetV of Four
func (s *Four[TT, UU, VV, WW]) SetV(v VV) {
	s.v = v
	s.which = V
}

// W of Four
// Panics if which != W
func (s Four[TT, UU, VV, WW]) W() WW {
	return check[WW](W, s.which, s.w)
}

// SetW of Four
func (s *Four[TT, UU, VV, WW]) SetW(w WW) {
	s.w = w
	s.which = W
}

// ==== Maybe

// Present returns true if Maybe contains a value
func (m Maybe[T]) Present() bool {
	return m.present
}

// Empty returns true if Maybe does not contain a value
func (m Maybe[T]) Empty() bool {
	return !m.present
}

// Get returns the value of the Maybe, which panics if it is not present
func (m Maybe[T]) Get() T {
	if !m.present {
		panic(errEmptyMaybe)
	}

	return m.v
}

// OrElse returns the value of the Maybe if present, or the else value provided if it is empty
func (m Maybe[T]) OrElse(elseVal T) T {
	return funcs.Ternary(m.present, m.v, elseVal)
}

// OrError returns (value of the Maybe, nil) if present, or (zero value of T, error provided) if not present
func (m Maybe[T]) OrError(e error) (res T, err error) {
	if m.present {
		res = m.v
	} else {
		err = e
	}

	return
}

// Set overwrites the current value with newVal, and sets m as present unless the newVal is a nil pointer
func (m *Maybe[T]) Set(newVal T) {
	// Store new value, which may be nil if T is a pointer type
	m.v = newVal

	// Set present to true unless T is a pointer type and newVal is nil
	m.present = !funcs.IsNilValue(newVal)
}

// SetEmpty overwrites the current value with the zero value, and sets m as empty
func (m *Maybe[T]) SetEmpty() {
	var zv T
	m.v = zv
	m.present = false
}

// SetOrError sets the current value with val if empty, else panics if a present val has already been set
func (m *Maybe[T]) SetOrError(val T) {
  if m.present {
    panic(errPresentMaybe)
  }

  m.v = val
  m.present = true
}

// ==== Result

// HasResult is true if Result contains an R
func (r Result[R]) HasResult() bool {
	return r.e == nil
}

// HasError is true if Result contains an error
func (r Result[R]) HasError() bool {
	return r.e != nil
}

// Get returns an R if there is a result, or zero value if there is an error
func (r Result[R]) Get() R {
	return r.r
}

// Error returns a nil error if there is a result, or a non-nil error if there is an error
func (r Result[R]) Error() error {
	return r.e
}
