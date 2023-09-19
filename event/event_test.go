package event

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Registry
func TestRegistry(t *testing.T) {
	var (
		r Registry[int, string]
		f Receiver[string] = ReceiverFunc[string](func(d string) string {
			return d + d
		})
	)

	r.Register(0, f) // 5 -> 55
	r.Register(0, f) // 55 -> 5555
	r.Register(0, f) // 5555 -> 55555555
	assert.Equal(t, "55555555", r.Send("5"))

	r.Remove(0, f)
	assert.Equal(t, "5555", r.Send("5"))

	r.Remove(0, f, ALL)
	assert.Equal(t, "5", r.Send("5"))

	r.Register(1, f) // 5 -> 55
	assert.Equal(t, "55", r.Send("5"))

	r.RemoveId(1)
	assert.Equal(t, "5", r.Send("5"))
}
