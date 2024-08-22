package iter

// SPDX-License-Identifier: Apache-2.0

import (
  //goio "io"
  "testing"
  
  "github.com/bantling/micro/union"
  "github.com/stretchr/testify/assert"
)

func TestToReader_(t *testing.T) {
  // One zero value
  it := Of(byte(0))
  r := ToReader(it)
  assert.Equal(t, union.OfResult(byte(0)), union.OfResultError(r.Next()))
  assert.Equal(t, union.OfError[byte](EOI), union.OfResultError(r.Next()))
  assert.Equal(t, union.OfError[byte](EOI), union.OfResultError(r.Next()))
}
