package util

// SPDX-License-Identifier: Apache-2.0
import (
	"fmt"
	"testing"

	"github.com/bantling/micro/go/funcs"
	"github.com/stretchr/testify/assert"
)

func TestOpenRange(t *testing.T) {
	r := OfRange(1, Open, 3, Closed, 2)

	// Die with nonsensical min/max values in constructor
	funcs.TryTo(
		func() {
			OfRange(3, Open, 1, Closed, 2)
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(ErrMinMaxMsg, "3", "1"), e)
		},
	)

	min, minOpen := r.GetMin()
	assert.Equal(t, 1, min)
	assert.Equal(t, Open, minOpen)

	max, maxOpen := r.GetMax()
	assert.Equal(t, 3, max)
	assert.Equal(t, Closed, maxOpen)

	assert.Equal(t, 2, r.GetValue())
	r.SetValue(3)
	assert.Equal(t, 3, r.GetValue())

	// Die setting to 1, as open min is 1, so val must be > 1
	funcs.TryTo(
		func() {
			r.SetValue(1)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(ErrOutsideRangeMsg, 1, "1", ">", 1, "<=", 3), e)
		},
	)

	// Die setting to 4, as closed max is 3, so val must be <= 3
	funcs.TryTo(
		func() {
			r.SetValue(4)
			assert.Fail(t, "Must die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf(ErrOutsideRangeMsg, 4, "4", ">", 1, "<=", 3), e)
		},
	)
}
