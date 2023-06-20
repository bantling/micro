package union

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/bantling/micro/funcs"
	"github.com/stretchr/testify/assert"
)

var (
	anErr = fmt.Errorf("An error")
)

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
				assert.Fail(t, "Must die")
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
				assert.Fail(t, "Must die")
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
				assert.Fail(t, "Must die")
			},
			func(r any) {
				e = r.(error)
			},
		)
		assert.Equal(t, fmt.Errorf("Member U is not available"), e)

		funcs.TryTo(
			func() {
				u3T.V()
				assert.Fail(t, "Must die")
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
				assert.Fail(t, "Must die")
			},
			func(r any) {
				e = r.(error)
			},
		)
		assert.Equal(t, fmt.Errorf("Member T is not available"), e)

		funcs.TryTo(
			func() {
				u3U.V()
				assert.Fail(t, "Must die")
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
				assert.Fail(t, "Must die")
			},
			func(r any) {
				e = r.(error)
			},
		)
		assert.Equal(t, fmt.Errorf("Member T is not available"), e)

		funcs.TryTo(
			func() {
				u3V.U()
				assert.Fail(t, "Must die")
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
				assert.Fail(t, "Must die")
			},
			func(r any) {
				e = r.(error)
			},
		)
		assert.Equal(t, fmt.Errorf("Member U is not available"), e)

		funcs.TryTo(
			func() {
				u4T.V()
				assert.Fail(t, "Must die")
			},
			func(r any) {
				e = r.(error)
			},
		)
		assert.Equal(t, fmt.Errorf("Member V is not available"), e)

		funcs.TryTo(
			func() {
				u4T.W()
				assert.Fail(t, "Must die")
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
				assert.Fail(t, "Must die")
			},
			func(r any) {
				e = r.(error)
			},
		)
		assert.Equal(t, fmt.Errorf("Member T is not available"), e)

		funcs.TryTo(
			func() {
				u4U.V()
				assert.Fail(t, "Must die")
			},
			func(r any) {
				e = r.(error)
			},
		)
		assert.Equal(t, fmt.Errorf("Member V is not available"), e)

		funcs.TryTo(
			func() {
				u4U.W()
				assert.Fail(t, "Must die")
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
				assert.Fail(t, "Must die")
			},
			func(r any) {
				e = r.(error)
			},
		)
		assert.Equal(t, fmt.Errorf("Member T is not available"), e)

		funcs.TryTo(
			func() {
				u4V.U()
				assert.Fail(t, "Must die")
			},
			func(r any) {
				e = r.(error)
			},
		)
		assert.Equal(t, fmt.Errorf("Member U is not available"), e)

		funcs.TryTo(
			func() {
				u4V.W()
				assert.Fail(t, "Must die")
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
				assert.Fail(t, "Must die")
			},
			func(r any) {
				e = r.(error)
			},
		)
		assert.Equal(t, fmt.Errorf("Member T is not available"), e)

		funcs.TryTo(
			func() {
				u4W.U()
				assert.Fail(t, "Must die")
			},
			func(r any) {
				e = r.(error)
			},
		)
		assert.Equal(t, fmt.Errorf("Member U is not available"), e)

		funcs.TryTo(
			func() {
				u4W.V()
				assert.Fail(t, "Must die")
			},
			func(r any) {
				e = r.(error)
			},
		)
		assert.Equal(t, fmt.Errorf("Member V is not available"), e)
	}
}

func TestResult_(t *testing.T) {
	// Result
	{
		res := OfResult("a")
		assert.True(t, res.HasResult())
		assert.False(t, res.HasError())
		assert.Equal(t, "a", res.Get())
		assert.Zero(t, res.Error())
	}

	// Error
	{
		e := fmt.Errorf("An Error")
		res := OfError[string](e)
		assert.False(t, res.HasResult())
		assert.True(t, res.HasError())
		assert.Zero(t, res.Get())
		assert.Equal(t, e, res.Error())

		funcs.TryTo(
			func() {
				OfError[string](nil)
				assert.Fail(t, "Must die")
			},
			func(r any) {
				e = r.(error)
			},
		)
		assert.Equal(t, fmt.Errorf("A Result cannot be set to a nil error"), e)
	}

	// ResultError
	{
		e := fmt.Errorf("An Error")
		res := OfResultError(1, nil)
		assert.True(t, res.HasResult())
		assert.False(t, res.HasError())
		assert.Equal(t, 1, res.Get())
		assert.Zero(t, res.Error())

		res = OfResultError(0, e)
		assert.False(t, res.HasResult())
		assert.True(t, res.HasError())
		assert.Zero(t, res.Get())
		assert.Equal(t, e, res.Error())

		funcs.TryTo(
			func() {
				OfResultError(1, e)
				assert.Fail(t, "Must die")
			},
			func(r any) {
				e = r.(error)
			},
		)
		assert.Equal(t, fmt.Errorf("A Result cannot have both a non-zero R value and a non-nil error"), e)
	}
}
