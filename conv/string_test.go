package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntToString_(t *testing.T) {
	assert.Equal(t, IntToString(int8(1)), "1")
	assert.Equal(t, IntToString(int16(2)), "2")
	assert.Equal(t, IntToString(int32(3)), "3")
	assert.Equal(t, IntToString(int64(4)), "4")
	assert.Equal(t, IntToString(int(5)), "5")
}

func TestUintToString_(t *testing.T) {
	assert.Equal(t, UintToString(uint8(1)), "1")
	assert.Equal(t, UintToString(uint16(2)), "2")
	assert.Equal(t, UintToString(uint32(3)), "3")
	assert.Equal(t, UintToString(uint64(4)), "4")
	assert.Equal(t, UintToString(uint(5)), "5")
}

func TestFloatToString_(t *testing.T) {
	assert.Equal(t, "1.25", FloatToString(float32(1.25)))
	assert.Equal(t, "1.25", FloatToString(float64(1.25)))
	assert.Equal(t, "-Inf", FloatToString(float32(math.Inf(-1))))
	assert.Equal(t, "+Inf", FloatToString(math.Inf(1)))
	assert.Equal(t, "NaN", FloatToString(math.NaN()))
	assert.Equal(t, "-0", FloatToString(-1/math.Inf(1)))
}

func TestBigIntToString_(t *testing.T) {
	assert.Equal(t, "1234", BigIntToString(big.NewInt(1234)))
}

func TestBigFloatToString_(t *testing.T) {
	assert.Equal(t, "1234.5678", BigFloatToString(big.NewFloat(1234.5678)))
	assert.Equal(t, "-Inf", BigFloatToString(big.NewFloat(math.Inf(-1))))
	assert.Equal(t, "+Inf", BigFloatToString(big.NewFloat(math.Inf(1))))
	assert.Equal(t, "-0", BigFloatToString(big.NewFloat(-1/math.Inf(1))))
}

func TestBigRatToString_(t *testing.T) {
	assert.Equal(t, "5/4", BigRatToString(big.NewRat(125, 100)))
}

func TestBigRatToNormalizedString_(t *testing.T) {
	assert.Equal(t, "1234", BigRatToNormalizedString(big.NewRat(1234, 1)))
	assert.Equal(t, "1.25", BigRatToNormalizedString(big.NewRat(125, 100)))
}
