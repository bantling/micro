package json

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/union"
	"github.com/stretchr/testify/assert"
)

func TestInit_(t *testing.T) {
  var v Value
  assert.Nil(t, conv.To(0, &v))
  assert.Equal(t, Value{typ: Number, val: union.Of4V[map[string]Value, []Value, string, bool]("0")}, v)
}
