package one28

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"
	"testing"

	"github.com/bantling/micro/conv"
	"github.com/stretchr/testify/assert"
)

func TestAdd_(t *testing.T) {
	carry, upper, lower := Add(1, 20, 2, 40)
	assert.Equal(t, uint64(0), carry)
	assert.Equal(t, uint64(3), upper)
	assert.Equal(t, uint64(60), lower)

	carry, upper, lower = Add(0xFF_00_00_00_00_00_00_00, 0, 0x01_00_00_00_00_00_00_00, 0)
	assert.Equal(t, uint64(1), carry)
	assert.Equal(t, uint64(0), upper)
	assert.Equal(t, uint64(0), lower)
}

func TestTwosComplement_(t *testing.T) {
  upper, lower := TwosComplement(0, 0xFF_FF_FF_FF_FF_FF_FF_FF)
  assert.Equal(t, uint64(0xFF_FF_FF_FF_FF_FF_FF_FF), upper)
  assert.Equal(t, uint64(1), lower)

  upper, lower = TwosComplement(0x80_00_00_00_00_00_00_00, 0)
  assert.Equal(t, uint64(0x80_00_00_00_00_00_00_00), upper)
  assert.Equal(t, uint64(0), lower)
}

func TestSub_(t *testing.T) {
	borrow, upper, lower := Sub(2, 40, 1, 20)
	assert.Equal(t, uint64(0), borrow)
	assert.Equal(t, uint64(1), upper)
	assert.Equal(t, uint64(20), lower)

	borrow, upper, lower = Sub(0xFF_00_00_00_00_00_00_00, 0, 0x01_00_00_00_00_00_00_00, 0)
	assert.Equal(t, uint64(0), borrow)
	assert.Equal(t, uint64(0xFE_00_00_00_00_00_00_00), upper)
	assert.Equal(t, uint64(0), lower)
}

func TestLsh_(t *testing.T) {
	carry, upper, lower := Lsh(1, 20)
	assert.Zero(t, carry)
	assert.Equal(t, uint64(2), upper)
	assert.Equal(t, uint64(40), lower)

	carry, upper, lower = Lsh(20, 1)
	assert.Zero(t, carry)
	assert.Equal(t, uint64(40), upper)
	assert.Equal(t, uint64(2), lower)

	carry, upper, lower = Lsh(0x87_00_00_00_00_00_00_12, 0x87_00_00_00_00_00_00_12)
	assert.Equal(t, uint64(1), carry)
	assert.Equal(t, uint64(0x0E_00_00_00_00_00_00_25), upper)
	assert.Equal(t, uint64(0x0E_00_00_00_00_00_00_24), lower)
}

func TestRsh_(t *testing.T) {
	carry, upper, lower := Rsh(1, 20)
	assert.Zero(t, carry)
	assert.Equal(t, uint64(0), upper)
	assert.Equal(t, uint64(0x80_00_00_00_00_00_00_0A), lower)

	carry, upper, lower = Rsh(20, 1)
	assert.Equal(t, uint64(1), carry)
	assert.Equal(t, uint64(10), upper)
	assert.Equal(t, uint64(0), lower)

	carry, upper, lower = Rsh(0x87_00_00_00_00_00_00_21, 0x87_00_00_00_00_00_00_21)
	assert.Equal(t, uint64(1), carry)
	assert.Equal(t, uint64(0x43_80_00_00_00_00_00_10), upper)
	assert.Equal(t, uint64(0xC3_80_00_00_00_00_00_10), lower)
}

func TestMul_(t *testing.T) {
	//==== 10 * 20
	var a, b uint64 = 10, 20
	c, d := Mul(a, b)
	assert.Equal(t, uint64(0), c)
	assert.Equal(t, uint64(200), d)

	//==== 0x10_00_00_0 * 0x20_00_00_00
	a, b = 0x10_00_00_00, 0x20_00_00_00
	c, d = Mul(a, b)
	assert.Equal(t, uint64(0), c)
	assert.Equal(t, a*b, d)

	//==== 0x10_20_30_40 * 0x50_60_70_80
	a, b = 0x10_20_30_40, 0x50_60_70_80
	c, d = Mul(a, b)
	assert.Equal(t, uint64(0), c)
	assert.Equal(t, a*b, d)

	//==== 0x10_20_30_40_50_60_70_80 * 0x90_A0_B0_C0_D0_E0_F0_00
	a, b = 0x10_20_30_40_50_60_70_80, 0x90_A0_B0_C0_D0_E0_F0_00

	// Calculate expected result using big.Int
	var abi, bbi, er *big.Int
	conv.To(a, &abi)
	conv.To(b, &bbi)
	conv.To(0, &er)
	er.Mul(abi, bbi)

	// Calculate actual result
	c, d = Mul(a, b)

	// Combine c and d into a 128-bit result for comparison
	var cdbi, dbi *big.Int
	conv.To(c, &cdbi)
	cdbi.Lsh(cdbi, 64)
	conv.To(d, &dbi)
	cdbi.Or(cdbi, dbi)

	// Assert we got the same result as big.Int
	assert.Zero(t, er.Cmp(cdbi))

	//==== 0xFF_FF_FF_FF_FF_FF_FF_FF * 0xFF_FF_FF_FF_FF_FF_FF_FF
	a, b = 0xFF_FF_FF_FF_FF_FF_FF_FF, 0xFF_FF_FF_FF_FF_FF_FF_FF

	// Calculate expected result using big.Int
	conv.To(a, &abi)
	conv.To(b, &bbi)
	conv.To(0, &er)
	er.Mul(abi, bbi)

	// Calculate actual result
	c, d = Mul(a, b)

	// Compare result to what we know it should be, based on pattern of multiplying highest bit patterns together
	assert.Equal(t, uint64(0xFF_FF_FF_FF_FF_FF_FF_FE), c)
	assert.Equal(t, uint64(0x00_00_00_00_00_00_00_01), d)

	// Combine c and d into a 128-bit result for comparison
	conv.To(c, &cdbi)
	cdbi.Lsh(cdbi, 64)
	conv.To(d, &dbi)
	cdbi.Or(cdbi, dbi)

	// Assert we got the same result as big.Int
	assert.Zero(t, er.Cmp(cdbi))
}
