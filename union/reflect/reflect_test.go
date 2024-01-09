package reflect

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"
	goreflect "reflect"
	"testing"

	// "github.com/bantling/micro/funcs"
	// "github.com/bantling/micro/tuple"
	"github.com/bantling/micro/union"
	"github.com/stretchr/testify/assert"
)

func TestGetMaybeType_(t *testing.T) {
	assert.Equal(t, goreflect.TypeOf(0), GetMaybeType(goreflect.TypeOf(union.Maybe[int]{})))
	assert.Equal(t, goreflect.TypeOf(0), GetMaybeType(goreflect.TypeOf((*union.Maybe[int])(nil))))

	assert.Equal(t, goreflect.TypeOf((*big.Int)(nil)), GetMaybeType(goreflect.TypeOf(union.Maybe[*big.Int]{})))
	assert.Equal(t, goreflect.TypeOf((*big.Int)(nil)), GetMaybeType(goreflect.TypeOf((*union.Maybe[*big.Int])(nil))))

	assert.Nil(t, GetMaybeType(goreflect.TypeOf(0)))
}

func TestGetMaybeValue_(t *testing.T) {
	val := GetMaybeValue(goreflect.ValueOf(union.Of(1)))
	assert.Equal(t, 1, GetMaybeValue(goreflect.ValueOf(union.Of(1))).Interface())

	val = GetMaybeValue(goreflect.ValueOf(union.Empty[int]()))
	assert.False(t, val.IsValid())
}

func TestSetMaybeValue_(t *testing.T) {
	var m union.Maybe[int]
	SetMaybeValue(goreflect.ValueOf(&m), goreflect.ValueOf(1))
	assert.True(t, m.Present())
	assert.Equal(t, 1, m.Get())
}

func TestSetMaybeValueEmpty_(t *testing.T) {
	m := union.Of(1)
	SetMaybeValueEmpty(goreflect.ValueOf(&m))
	assert.False(t, m.Present())
}
