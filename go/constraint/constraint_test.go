package constraint

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSignedInt(t *testing.T) {
	var (
		i   int
		i8  int8
		i16 int16
		i32 int32
		i64 int64

		ui   uint
		ui8  uint8
		ui16 uint16
		ui32 uint32
		ui64 uint64

		f32 float32
		f64 float64

		bi *big.Int
		bf *big.Float
		br *big.Rat
	)

	assert.True(t, IsSignedInt(i))
	assert.True(t, IsSignedInt(i8))
	assert.True(t, IsSignedInt(i16))
	assert.True(t, IsSignedInt(i32))
	assert.True(t, IsSignedInt(i64))

	assert.False(t, IsSignedInt(ui))
	assert.False(t, IsSignedInt(ui8))
	assert.False(t, IsSignedInt(ui16))
	assert.False(t, IsSignedInt(ui32))
	assert.False(t, IsSignedInt(ui64))

	assert.False(t, IsSignedInt(f32))
	assert.False(t, IsSignedInt(f64))

	assert.False(t, IsSignedInt(bi))
	assert.False(t, IsSignedInt(bf))
	assert.False(t, IsSignedInt(br))
}

func TestIsUnsignedInt(t *testing.T) {
	var (
		i   int
		i8  int8
		i16 int16
		i32 int32
		i64 int64

		ui   uint
		ui8  uint8
		ui16 uint16
		ui32 uint32
		ui64 uint64

		f32 float32
		f64 float64

		bi *big.Int
		bf *big.Float
		br *big.Rat
	)

	assert.False(t, IsUnsignedInt(i))
	assert.False(t, IsUnsignedInt(i8))
	assert.False(t, IsUnsignedInt(i16))
	assert.False(t, IsUnsignedInt(i32))
	assert.False(t, IsUnsignedInt(i64))

	assert.True(t, IsUnsignedInt(ui))
	assert.True(t, IsUnsignedInt(ui8))
	assert.True(t, IsUnsignedInt(ui16))
	assert.True(t, IsUnsignedInt(ui32))
	assert.True(t, IsUnsignedInt(ui64))

	assert.False(t, IsUnsignedInt(f32))
	assert.False(t, IsUnsignedInt(f64))

	assert.False(t, IsUnsignedInt(bi))
	assert.False(t, IsUnsignedInt(bf))
	assert.False(t, IsUnsignedInt(br))
}

func TestIsFloat(t *testing.T) {
	var (
		i   int
		i8  int8
		i16 int16
		i32 int32
		i64 int64

		ui   uint
		ui8  uint8
		ui16 uint16
		ui32 uint32
		ui64 uint64

		f32 float32
		f64 float64

		bi *big.Int
		bf *big.Float
		br *big.Rat
	)

	assert.False(t, IsFloat(i))
	assert.False(t, IsFloat(i8))
	assert.False(t, IsFloat(i16))
	assert.False(t, IsFloat(i32))
	assert.False(t, IsFloat(i64))

	assert.False(t, IsFloat(ui))
	assert.False(t, IsFloat(ui8))
	assert.False(t, IsFloat(ui16))
	assert.False(t, IsFloat(ui32))
	assert.False(t, IsFloat(ui64))

	assert.True(t, IsFloat(f32))
	assert.True(t, IsFloat(f64))

	assert.False(t, IsFloat(bi))
	assert.False(t, IsFloat(bf))
	assert.False(t, IsFloat(br))
}

func TestIsBig(t *testing.T) {
	var (
		i   int
		i8  int8
		i16 int16
		i32 int32
		i64 int64

		ui   uint
		ui8  uint8
		ui16 uint16
		ui32 uint32
		ui64 uint64

		f32 float32
		f64 float64

		bi *big.Int
		bf *big.Float
		br *big.Rat
	)

	assert.False(t, IsBig(i))
	assert.False(t, IsBig(i8))
	assert.False(t, IsBig(i16))
	assert.False(t, IsBig(i32))
	assert.False(t, IsBig(i64))

	assert.False(t, IsBig(ui))
	assert.False(t, IsBig(ui8))
	assert.False(t, IsBig(ui16))
	assert.False(t, IsBig(ui32))
	assert.False(t, IsBig(ui64))

	assert.False(t, IsBig(f32))
	assert.False(t, IsBig(f64))

	assert.True(t, IsBig(bi))
	assert.True(t, IsBig(bf))
	assert.True(t, IsBig(br))
}
