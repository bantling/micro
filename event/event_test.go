package event

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/bantling/micro/funcs"
	"github.com/stretchr/testify/assert"
)

type addReceiver struct {
	s string
}

func (nr addReceiver) Process(s string) string {
	return s + nr.s
}

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

	// Remove one copy of f
	r.Remove(0, f)
	assert.Equal(t, "5555", r.Send("5"))

	// Remove all copies of f - no receivers, but still have an entry for id 0
	r.Remove(0, f, funcs.ALL)
	assert.Equal(t, "5", r.Send("5"))

	// Use id 2
	r.Register(2, f) // 5 -> 55
	assert.Equal(t, "55", r.Send("5"))

	// Use id 1
	r.Register(1, f) // 5 -> 5555
	assert.Equal(t, "5555", r.Send("5"))

	// Remove id 1, leaving only id 2
	r.RemoveId(1)
	assert.Equal(t, "55", r.Send("5"))

	rcvr1 := addReceiver{" blahdy"}
	r.Register(3, rcvr1)
	assert.Equal(t, "55 blahdy", r.Send("5"))

	rcvr2 := addReceiver{" dude"}
	r.Register(3, rcvr2)
	assert.Equal(t, "55 blahdy dude", r.Send("5"))

	r.Remove(3, rcvr2)
	assert.Equal(t, "55 blahdy", r.Send("5"))

	r.Remove(3, rcvr1)
	assert.Equal(t, "55", r.Send("5"))
}
