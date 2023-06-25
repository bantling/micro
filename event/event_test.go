package event

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/bantling/micro/encoding/json"
	"github.com/stretchr/testify/assert"
)

// Test Registry
func TestRegistry(t *testing.T) {
	var (
		r Registry[string]
		f Receiver[string] = ReceiverFunc[string](func(d string) string {
			return d + d
		})
	)

	r.Register(f) // 5 -> 55
	r.Register(f) // 55 -> 5555
	r.Register(f) // 5555 -> 55555555
	assert.Equal(t, "55555555", r.Send("5"))

	r.Remove(f)
	assert.Equal(t, "5555", r.Send("5"))

	r.Remove(f, ALL)
	assert.Equal(t, "5", r.Send("5"))
}

// Test DefaultRegistry
func TestDefaultRegistry(t *testing.T) {
	var (
		f Receiver[json.Value] = ReceiverFunc[json.Value](func(d json.Value) json.Value {
			return json.FromString(d.AsString() + d.AsString())
		})
	)

	Register(f) // 5 -> 55
	Register(f) // 55 -> 5555
	Register(f) // 5555 -> 55555555
	assert.Equal(t, "55555555", Send(json.FromString("5")).AsString())

	Remove(f)
	assert.Equal(t, "5555", Send(json.FromString("5")).AsString())

	Remove(f, ALL)
	assert.Equal(t, "5", Send(json.FromString("5")).AsString())
}
