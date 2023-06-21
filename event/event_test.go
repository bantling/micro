package event

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRegistry
func TestRegistry(t *testing.T) {
	var (
		r Registry[string]
		f Receiver[string] = ReceiverFunc[string](func(d string) string {
			return d + d
		})
	)

	r.Register(f)
	assert.Equal(t, "55", r.Send("5"))

	r.Remove(f)
	assert.Equal(t, "5", r.Send("5"))
}
