package union

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/bantling/micro/go/funcs"
	"github.com/stretchr/testify/assert"
)

var (
	anErr = fmt.Errorf("An error")
)

// ==== Constructors

func Test2_(t *testing.T) {
  var e error

  // T
  {
  	u2T := Of2T[string, int]("a")
  	assert.Equal(t, T, u2T.Which())
  	assert.Equal(t, "a", u2T.T())
    u2T.SetT("b")
  	assert.Equal(t, T, u2T.Which())
  	assert.Equal(t, "b", u2T.T())

    funcs.TryTo(
      func() {
        u2T.U()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member U is not available"), e)
  }

  // U
  {
  	u2U := Of2U[string, int](1)
  	assert.Equal(t, U, u2U.Which())
  	assert.Equal(t, 1, u2U.U())
    u2U.SetU(2)
  	assert.Equal(t, U, u2U.Which())
  	assert.Equal(t, 2, u2U.U())

    funcs.TryTo(
      func() {
        u2U.T()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member T is not available"), e)
  }
}

func Test3_(t *testing.T) {
  var e error

  // T
  {
  	u3T := Of3T[string, int, string]("a")
  	assert.Equal(t, T, u3T.Which())
  	assert.Equal(t, "a", u3T.T())
    u3T.SetT("b")
  	assert.Equal(t, T, u3T.Which())
  	assert.Equal(t, "b", u3T.T())

    funcs.TryTo(
      func() {
        u3T.U()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member U is not available"), e)

    funcs.TryTo(
      func() {
        u3T.V()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member V is not available"), e)
  }

  // U
  {
  	u3U := Of3U[string, int, string](1)
  	assert.Equal(t, U, u3U.Which())
  	assert.Equal(t, 1, u3U.U())
    u3U.SetU(2)
  	assert.Equal(t, U, u3U.Which())
  	assert.Equal(t, 2, u3U.U())

    funcs.TryTo(
      func() {
        u3U.T()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member T is not available"), e)

    funcs.TryTo(
      func() {
        u3U.V()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member V is not available"), e)
  }

  // V
  {
  	u3V := Of3V[string, int, string]("a")
  	assert.Equal(t, V, u3V.Which())
  	assert.Equal(t, "a", u3V.V())
    u3V.SetV("b")
  	assert.Equal(t, V, u3V.Which())
  	assert.Equal(t, "b", u3V.V())

    funcs.TryTo(
      func() {
        u3V.T()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member T is not available"), e)

    funcs.TryTo(
      func() {
        u3V.U()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member U is not available"), e)
  }
}

func Test4_(t *testing.T) {
  var e error

  // T
  {
  	u4T := Of4T[string, int, string, int]("a")
  	assert.Equal(t, T, u4T.Which())
  	assert.Equal(t, "a", u4T.T())
    u4T.SetT("b")
  	assert.Equal(t, T, u4T.Which())
  	assert.Equal(t, "b", u4T.T())

    funcs.TryTo(
      func() {
        u4T.U()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member U is not available"), e)

    funcs.TryTo(
      func() {
        u4T.V()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member V is not available"), e)

    funcs.TryTo(
      func() {
        u4T.W()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member W is not available"), e)
  }

  // U
  {
  	u4U := Of4U[string, int, string, int](1)
  	assert.Equal(t, U, u4U.Which())
  	assert.Equal(t, 1, u4U.U())
    u4U.SetU(2)
  	assert.Equal(t, U, u4U.Which())
  	assert.Equal(t, 2, u4U.U())

    funcs.TryTo(
      func() {
        u4U.T()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member T is not available"), e)

    funcs.TryTo(
      func() {
        u4U.V()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member V is not available"), e)

    funcs.TryTo(
      func() {
        u4U.W()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member W is not available"), e)
  }

  // V
  {
  	u4V := Of4V[string, int, string, int]("a")
  	assert.Equal(t, V, u4V.Which())
  	assert.Equal(t, "a", u4V.V())
    u4V.SetV("b")
  	assert.Equal(t, V, u4V.Which())
  	assert.Equal(t, "b", u4V.V())

    funcs.TryTo(
      func() {
        u4V.T()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member T is not available"), e)

    funcs.TryTo(
      func() {
        u4V.U()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member U is not available"), e)

    funcs.TryTo(
      func() {
        u4V.W()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member W is not available"), e)
  }

  // W
  {
  	u4W := Of4W[string, int, string, int](1)
  	assert.Equal(t, W, u4W.Which())
  	assert.Equal(t, 1, u4W.W())
    u4W.SetW(2)
  	assert.Equal(t, W, u4W.Which())
  	assert.Equal(t, 2, u4W.W())

    funcs.TryTo(
      func() {
        u4W.T()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member T is not available"), e)

    funcs.TryTo(
      func() {
        u4W.U()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member U is not available"), e)

    funcs.TryTo(
      func() {
        u4W.V()
      },
      func(r any) {
        e = r.(error)
      },
    )
    assert.Equal(t, fmt.Errorf("Member V is not available"), e)
  }
}
