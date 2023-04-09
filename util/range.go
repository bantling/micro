package util

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/bantling/micro/constraint"
	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/funcs"
)

// RangeMode indicates if a minimum or maximum value in a range is open or closed
type RangeMode bool

const (
	Open   RangeMode = false
	Closed RangeMode = true
)

// Error messages
var (
	errMinMaxMsg       = "The (min, max) values of (%s, %s) are not allowed, min must be < max and max must be > min"
	errOutsideRangeMsg = "The %T value %s is not valid, as the value must be %s %s and %s %s"
)

// Range represents a range of values.
// The range may be open, half open, or closed.
type Range[T constraint.IntegerAndFloat] struct {
	min     T
	minMode RangeMode
	max     T
	maxMode RangeMode
	val     T
}

// OfRange constructs a range from minimum value and mode, maximum value and mode, and initial value.
// An initial value is required since min and max could both be open and the type could a float, so there is no sensible
// default initial value.
//
// Panics if the initial value is not in the specified range.
func OfRange[T constraint.IntegerAndFloat](
	min T,
	minMode RangeMode,
	max T,
	maxMode RangeMode,
	initial T,
) *Range[T] {
	if (min >= max) || (max <= min) {
		panic(fmt.Errorf(errMinMaxMsg, conv.ToString(min), conv.ToString(max)))
	}

	return &Range[T]{min, minMode, max, maxMode, initial}
}

// GetMin returns the minimum value and mode
func (r Range[T]) GetMin() (T, RangeMode) {
	return r.min, r.minMode
}

// GetMax returns the maximum value and mode
func (r Range[T]) GetMax() (T, RangeMode) {
	return r.max, r.maxMode
}

// GetValue returns the current value
func (r Range[T]) GetValue() T {
	return r.val
}

// SetValue sets the value.
//
//	Returns an error if the value cannot be set because it is outside the defined bounds, otherwise returns nil
func (r *Range[T]) SetValue(val T) error {
	if (((r.minMode == Open) && (val > r.min)) || ((r.minMode == Closed) && (val >= r.min))) &&
		(((r.maxMode == Open) && (val < r.max)) || ((r.maxMode == Closed) && (val <= r.max))) {
		r.val = val
		return nil
	}

	return fmt.Errorf(
		errOutsideRangeMsg,
		val,
		conv.ToString(val),
		funcs.Ternary(r.minMode == Open, ">", ">="), conv.ToString(r.min),
		funcs.Ternary(r.maxMode == Open, "<", "<="), conv.ToString(r.max),
	)
}
