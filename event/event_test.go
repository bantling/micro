package event

// SPDX-License-Identifier: Apache-2.0

import (
	// "fmt"
	"testing"

	"github.com/bantling/micro/conv"
	"github.com/stretchr/testify/assert"
)

// TestRegistry
func TestRegistry(t *testing.T) {
	var (
		r   Registry[int, string]
		str string
		d                             = Data[int, string]{0, 0, &str}
		f   ReceiverFunc[int, string] = ReceiverFunc[int, string](func(d Data[int, string]) {
			conv.To(d.Input, d.Result)
		})
	)
	r.Register(0, f)

	d.Input = 5
	assert.True(t, r.Send(d))
	assert.Equal(t, "5", *d.Result)
}
