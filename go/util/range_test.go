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
			assert.Equal(t, fmt.Errorf(errMinMaxMsg, "3", "1"), e)
		},
	)

	min, minOpen := r.GetMin()
	assert.Equal(t, 1, min)
	assert.Equal(t, Open, minOpen)

	max, maxOpen := r.GetMax()
	assert.Equal(t, 3, max)
	assert.Equal(t, Closed, maxOpen)

	assert.Equal(t, 2, r.GetValue())
	assert.Nil(t, r.SetValue(3))
	assert.Equal(t, 3, r.GetValue())

	// Error setting to 1, as open min is 1, so val must be > 1
	assert.Equal(t, fmt.Errorf(errOutsideRangeMsg, 1, "1", ">", 1, "<=", 3), r.SetValue(1))

	// Error setting to 4, as closed max is 3, so val must be <= 3
	assert.Equal(t, fmt.Errorf(errOutsideRangeMsg, 4, "4", ">", 1, "<=", 3), r.SetValue(4))
}
