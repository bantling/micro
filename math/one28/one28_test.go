package one28

// SPDX-License-Identifier: Apache-2.0

import (
	gomath "math"
	"math/big"
	"testing"

	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/math"
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
	upper, lower := Sub(2, 40, 1, 20)
	assert.Equal(t, uint64(1), upper)
	assert.Equal(t, uint64(20), lower)

	upper, lower = Sub(0xFF_00_00_00_00_00_00_00, 0, 0x01_00_00_00_00_00_00_00, 0)
	assert.Equal(t, uint64(0xFE_00_00_00_00_00_00_00), upper)
	assert.Equal(t, uint64(0), lower)
}

func TestLsh_(t *testing.T) {
	carry, upper, lower := Lsh(1, 20)
	assert.Zero(t, carry)
	assert.Equal(t, uint64(2), upper)
	assert.Equal(t, uint64(40), lower)

	carry, upper, lower = Lsh(1, 20, 4)
	assert.Zero(t, carry)
	assert.Equal(t, uint64(16), upper)
	assert.Equal(t, uint64(320), lower)

	carry, upper, lower = Lsh(20, 1)
	assert.Zero(t, carry)
	assert.Equal(t, uint64(40), upper)
	assert.Equal(t, uint64(2), lower)

	carry, upper, lower = Lsh(20, 1, 3)
	assert.Zero(t, carry)
	assert.Equal(t, uint64(160), upper)
	assert.Equal(t, uint64(8), lower)

	//                           1000 0111     0001 0010    1000 0111      0001 0002
	carry, upper, lower = Lsh(0x87_00_00_00_00_00_00_12, 0x87_00_00_00_00_00_00_12)
	assert.Equal(t, uint64(1), carry)
	assert.Equal(t, uint64(0x0E_00_00_00_00_00_00_25), upper)
	assert.Equal(t, uint64(0x0E_00_00_00_00_00_00_24), lower)

	//                          1110 0111      0001 0010    1000 0111      0001 0002
	carry, upper, lower = Lsh(0xE7_00_00_00_00_00_00_12, 0x87_00_00_00_00_00_00_12, 2)
	assert.Equal(t, uint64(3), carry)
	assert.Equal(t, uint64(0x9C_00_00_00_00_00_00_4A), upper)
	assert.Equal(t, uint64(0x1C_00_00_00_00_00_00_48), lower)
}

func TestRsh_(t *testing.T) {
	upper, lower, carry := Rsh(1, 20)
	assert.Equal(t, uint64(0), upper)
	assert.Equal(t, uint64(0x80_00_00_00_00_00_00_0A), lower)
	assert.Zero(t, carry)

	upper, lower, carry = Rsh(1, 20, 4)
	assert.Equal(t, uint64(0), upper)
	assert.Equal(t, uint64(0x10_00_00_00_00_00_00_01), lower)
	assert.Equal(t, uint64(0x40_00_00_00_00_00_00_00), carry)

	upper, lower, carry = Rsh(20, 1)
	assert.Equal(t, uint64(10), upper)
	assert.Equal(t, uint64(0), lower)
	assert.Equal(t, uint64(0x80_00_00_00_00_00_00_00), carry)

	upper, lower, carry = Rsh(20, 1, 3)
	assert.Equal(t, uint64(2), upper)
	assert.Equal(t, uint64(0x80_00_00_00_00_00_00_00), lower)
	assert.Equal(t, uint64(0x20_00_00_00_00_00_00_00), carry)

	upper, lower, carry = Rsh(0x87_00_00_00_00_00_00_21, 0x87_00_00_00_00_00_00_21)
	assert.Equal(t, uint64(0x43_80_00_00_00_00_00_10), upper)
	assert.Equal(t, uint64(0xC3_80_00_00_00_00_00_10), lower)
	assert.Equal(t, uint64(0x80_00_00_00_00_00_00_00), carry)

	upper, lower, carry = Rsh(0x87_00_00_00_00_00_00_21, 0x87_00_00_00_00_00_00_21, 2)
	assert.Equal(t, uint64(0x21_C0_00_00_00_00_00_08), upper)
	assert.Equal(t, uint64(0x61_C0_00_00_00_00_00_08), lower)
	assert.Equal(t, uint64(0x40_00_00_00_00_00_00_00), carry)
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

	// Compare actual result to what we know it should be, based on pattern of multiplying highest bit patterns together
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

func TestQuoRem_(t *testing.T) {
	// Die if division by zero
	funcs.TryTo(
		func() { QuoRem(1, 2, 0) },
		func(e any) { assert.Equal(t, math.DivByZeroErr, e) },
	)

	// Shortcut that uses built in operators for case of upper quotient = 0
	uq, lq, rm := QuoRem(0, 100, 11)
	assert.Equal(t, uint64(0), uq)
	assert.Equal(t, uint64(9), lq)
	assert.Equal(t, uint64(1), rm)

	// Long case of upper quotient > 0, where remainder = 0
	uq, lq, rm = QuoRem(100, 0, 2) // 100 * 2^64 / 2 = 100 * 2^32 rmdr 0
	assert.Equal(t, uint64(50), uq)
	assert.Equal(t, uint64(0), lq)
	assert.Equal(t, uint64(0), rm)

	// Long case of upper quotient > 0, where remainder = 1
	uq, lq, rm = QuoRem(100, 3, 2) // 100 * 2^64 + 3 / 2 = 50 * 2^64 + 1 rmdr 1
	assert.Equal(t, uint64(50), uq)
	assert.Equal(t, uint64(1), lq)
	assert.Equal(t, uint64(1), rm)

	// Long case of upper quotient > 0, stupidly dividing by 1
	uq, lq, rm = QuoRem(100, 3, 1) // 100 * 2^64 + 3
	assert.Equal(t, uint64(100), uq)
	assert.Equal(t, uint64(3), lq)
	assert.Equal(t, uint64(0), rm)

	//// Long case of a 32 digit number
	var numbi *big.Int
	conv.To("12345678901234567890123456789012", &numbi)

	// Split long num into upper and lower 64 bits
	var (
		unumbi, lnumbi, lmask *big.Int
		unum, lnum            uint64
	)
	// unumbi = upper 64
	conv.To(numbi, &unumbi)
	unumbi.Rsh(unumbi, 64)

	// lnumbi = lower 64
	conv.To(numbi, &lnumbi)
	conv.To(uint64(gomath.MaxUint64), &lmask)
	lnumbi.And(lnumbi, lmask)

	// extract upper and lower 64 into uint64s
	conv.To(unumbi, &unum)
	conv.To(lnumbi, &lnum)

	// Divide long number by 10 using our function
	uq, lq, rm = QuoRem(unum, lnum, 10)

	// Calculate expected result using bigInt calcs
	var (
		tenbi, uqbi, lqbi, rmbi *big.Int
		ueq, leq, er            uint64
	)

	// Divide original numbi by 10
	conv.To(10, &tenbi)
	conv.To(0, &uqbi)
	conv.To(0, &rmbi)
	uqbi.QuoRem(numbi, tenbi, rmbi)

	// uqbi = 128 bit result, copy it to lqbi
	conv.To(uqbi, &lqbi)

	// uqbi = upper 64
	uqbi.Rsh(uqbi, 64)

	// lqbi = lower 64
	lqbi.And(lqbi, lmask)

	// extract upper and lower 64 and remainder into uint64s
	conv.To(uqbi, &ueq)
	conv.To(lqbi, &leq)
	conv.To(rmbi, &er)

	// Check expected result is correct - combine upper and lower 64, multiply by 10, add remainder, and compare
	var cbi *big.Int
	conv.To(uqbi, &cbi)
	cbi.Lsh(cbi, 64).Or(cbi, lqbi).Mul(cbi, tenbi).Add(cbi, rmbi)
	assert.Equal(t, numbi, cbi)

	// Check our result
	assert.Equal(t, ueq, uq)
	assert.Equal(t, leq, lq)
	assert.Equal(t, er, rm)
}
