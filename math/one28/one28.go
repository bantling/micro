package one28

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/math"
)

const (
	highestBitMask uint64 = 0x80_00_00_00_00_00_00_00
	allBitsMask    uint64 = 0xFF_FF_FF_FF_FF_FF_FF_FF
	upper32Mask    uint64 = 0xFF_FF_FF_FF_00_00_00_00
	lower32Mask    uint64 = 0x00_00_00_00_FF_FF_FF_FF
)

// Add adds two 128-bit numbers together, represented as pairs of uint64
func Add(upperAE1, lowerAE1, upperAE2, lowerAE2 uint64) (carry, upper, lower uint64) {
	var (
		// Split all 4 inputs into top and bottom 32 bits, stored as uint64
		ae1ut32, ae1ub32 = (upperAE1 & upper32Mask) >> 32, upperAE1 & lower32Mask
		ae1lt32, ae1lb32 = (lowerAE1 & upper32Mask) >> 32, lowerAE1 & lower32Mask
		ae2ut32, ae2ub32 = (upperAE2 & upper32Mask) >> 32, upperAE2 & lower32Mask
		ae2lt32, ae2lb32 = (lowerAE2 & upper32Mask) >> 32, lowerAE2 & lower32Mask

		// Combine above into 8 32-bit values into 4 32-bit terms ta - td where ta is highest and td is lowest, with carries
		td       = ae1lb32 + ae2lb32
		td_carry = (td & upper32Mask) >> 32

		tc       = td_carry + ae1lt32 + ae2lt32
		tc_carry = ((tc & upper32Mask) >> 32)

		tb       = tc_carry + ae1ub32 + ae2ub32
		tb_carry = ((tb & upper32Mask) >> 32)

		ta = tb_carry + (ae1ut32 + ae2ut32)
	)

	// Combine terms into carry, upper and lower
	carry = (ta & upper32Mask) >> 32
	upper = ((ta & lower32Mask) << 32) | (tb & lower32Mask)
	lower = ((tc & lower32Mask) << 32) | (td & lower32Mask)

	return
}

// TwosComplement calculates the twos complement of a 12-bit number (invert all bits and add one)
func TwosComplement(upper, lower uint64) (upperRes, lowerRes uint64) {
	upperRes = upper ^ allBitsMask
	if lowerRes = (lower ^ allBitsMask) + 1; lowerRes == 0 {
		upperRes++
	}

	return
}

// Sub subtracts a 128-bit subtrahend from a 128-bit minuend, returning a 128-bit result.
// There is no borrow returned, as this is unsigned subtraction. It is up to the caller to ensure minuend >= subtrahend.
func Sub(upperME, lowerME, upperSE, lowerSE uint64) (upper, lower uint64) {
	// Calculate the twos complement and add, it's easier
	upperSE, lowerSE = TwosComplement(upperSE, lowerSE)
	_, upper, lower = Add(upperME, lowerME, upperSE, lowerSE)

	return
}

// Lsh shifts a 128 bit value left n bits (default 1), returning the highest bit as a carry.
func Lsh(upper, lower uint64, n ...uint) (carry, upperRes, lowerRes uint64) {
	upperRes, lowerRes = upper, lower
	for nv := funcs.SliceIndex(n, 0, 1); nv > 0; nv-- {
		// Carry is highest bit that get shifted out
		carry = (carry << 1) | ((upperRes & highestBitMask) >> 63)

		// If highest bit of lower is set, then set lowest bit of higher to shift the bit across to higher value
		upperRes = (upperRes << 1) | ((lowerRes & highestBitMask) >> 63)

		// Lower is always just shifted
		lowerRes <<= 1
	}

	return
}

// Rsh shifts a 128 bit value right 1 bit, returning the lowest bit as a carry.
func Rsh(upper, lower uint64, n ...uint) (upperRes, lowerRes, carry uint64) {
	upperRes, lowerRes = upper, lower
	for nv := funcs.SliceIndex(n, 0, 1); nv > 0; nv-- {
		// Carry is lowest bit that gets shifted out
		carry = (carry >> 1) | ((lowerRes & 1) << 63)

		// If lowest bit of higher is set, then set highest bit of lower to shift the bit across to lower value
		lowerRes = ((upperRes & 1) << 63) | (lowerRes >> 1)

		// Higher is always just shifted
		upperRes >>= 1
	}

	return
}

// Mul multiplies two uint64 values into a pair of uint64 values that represent a 128-bit result.
func Mul(mp, ma uint64) (upper, lower uint64) {
	// There is a simple rule for multiplying two maximum value n-bit integers for some even number n:
	// The result is of the form (F)E(0)1, where the number F and 0 digits is the same: n / 2 - 1.
	// EG:
	// for two  8-bit values, we have  8 / 2 - 1 =  1 F and 0, producing FE01.
	// for two 16-bit values, we have 16 / 4 - 1 =  3 F and 0, producing FFFE_0001.
	// for two 32-bit values, we have 32 / 4 - 1 =  7 F and 0, producing FFFF_FFFE_0000_0001.
	// for two 64-bit values, we have 64 / 4 - 1 = 15 F and 0, producing FFFF_FFFF_FFFF_FFFE_0000_0000_0000_0001.
	//
	// We can perform multiplication of two n-bit values using only n-bit integers, by breaking up the two n-bit values into
	// two pairs of n/2-bit values, which we call (a,b) and (c,d). The results are stored in four n/2-bit slots,
	// which we call e, f, g, and h.
	//
	// We need to break the result down into four multiplications:
	// ab * cd = b*d + b*c + a*d + a*c
	// The results are then placed into the slots.
	//
	// The folowing explanation shows how to multiply two maximum value 16-bit numbers:
	//
	// a = b = c = d = FF, and FF * FF = FE01, so b*d = b*c = a*d = a*c = FF*FF = FE01.
	//
	// The difference between the terms is not their value, but their position:
	// a and c are high bytes, so are multiplied by 2^8, which means shifting left 8 bits.
	// b and d are low  bytes, so are multiplied by 2^0, which means shifting left 0 bits.
	//
	// b*d is  low * low , has a total multiple of 2^0, stored in slots g and h
	// b*c is  low * high, has a total multiple of 2^8, stored in slots f and g
	// a*d is high * low , has a total multiple of 2^8, stored in slots f and g
	// a*c is high * high, has a total multiple of 2^16,stored in slots e and f
	//
	// Slot h = low half of b*d = b*d & 0xFF
	// Slot g = high half of b*d + low half of b*c + low half of a*d = (b*d >> 8) + (b*c & 0xFF) + (a*d & 0xFF)
	// Slot f = carry from g + high half of b*c + high half of a*d + low half of a*c = g carry + (b*c >> 8) + (a*d >> 8) + (a*c & 0xFF)
	// Slot a = carry from f + high half of a*c = f carry + (a*c >> 8)
	//
	// Note how the multiple additions for g and f can produce a carry: adding 3 or 4 integers of n can require up to 2 extra
	// bits. By writing a utility function that adds four 8-bit integers (received as 16-bit integers for convenience) and
	// produces two 16-bit integers (carry and sum), we can use math that multiples two 16 bit integer using only 16-bit math.
	//
	// Adding four FE01 results multiplied by 2^0, 2^8, 2^8, and 2^16:
	//            111 1
	//           0000 FE01 b*d = FE01 * 2^0
	//         + 00FE 0100 b*c = FE01 * 2^8
	//         + 00FE 0100 a*d = FE01 * 2^8
	//         + FE01 0000 a*c = FE01 * 2^16
	//         = FFFE 0001
	//         = eeff gghh
	//
	// The idea can be extended to 64-bit math as follows:
	// - Change all uint16 into uint64
	// - Change all (x & 0xFF) expressions into (x & 0xFF_FF_FF_FF)
	// - Change all (x >> 8) expressions into (x >> 32)
	var (
		// add receives type uint64 for convenience, but they are actually 32 bit values.
		// The result may require 34 bits, and is expressed as a pair of uint64s for convenience.
		add = func(v1, v2, v3, v4 uint64) (carry, result uint64) {
			var res uint64 = v1 + v2 + v3 + v4
			carry = res >> 32
			result = res & 0xFF_FF_FF_FF

			return
		}

		lmp uint64 = (mp & 0xFF_FF_FF_FF)
		hmp uint64 = mp >> 32
		lma uint64 = (ma & 0xFF_FF_FF_FF)
		hma uint64 = ma >> 32

		bd uint64 = lmp * lma
		bc uint64 = lmp * hma
		ad uint64 = hmp * lma
		ac uint64 = hmp * hma

		h          uint64 = bd & 0xFF_FF_FF_FF
		g_carry, g uint64 = add(0, bd>>32, bc&0xFF_FF_FF_FF, ad&0xFF_FF_FF_FF)
		f_carry, f uint64 = add(g_carry, bc>>32, ad>>32, ac&0xFF_FF_FF_FF)
		e          uint64 = f_carry + (ac >> 32)
	)

	lower = (g << 32) | h
	upper = (e << 32) | f

	return
}

// QuoRem divides a pair of uint64 values that represent a 128-bit input into a pair of uint64 128-bit output and
// remainder. The division is performed using bit shifting division, which is similar to bit shifting multiplication.
//
// Example: dividing 121 by 5
//
// m (multiple)  = 5
// f (factor)    = 1
// q (quotient)  = 0
// r (remainder) = 121
//
// Phase 1: Find largest multiple of 5 <= 121, by shifting multiple and factor left 1 bit at a time
// 5   < 121 : m = 10,  f = 2
// 10  < 121 : m = 20,  f = 4
// 20  < 121 : m = 40,  f = 8
// 40  < 121 : m = 80,  f = 16
// 80  < 121 : m = 160, f = 32
// 160 > 121 : m = 80,  f = 16
//
// Phase 2: Subtract multiples, shifting multiple and factor right 1 bit at a time.
// Relevant multiples are <= remainder. Stop when remainder < divisor.
// 80  <= 121 : q = 0  + 16 = 16, r = 121 - 80 = 41, r > 5, m = 40, f = 8
// 40  <=  41 : q = 16 +  8 = 24, r =  41 - 40 =  1, r < 5, stop
//
// Result is 121 / 5 = 24 remainder 1
//
// If upper dividend is 0, just uses division and modulus operators for a fast result.
// Panics if the divisor is 0.
func QuoRem(upperDE, lowerDE, divisor uint64) (upperQ, lowerQ, remainder uint64) {
	// Die if divisor is 0
	if divisor == 0 {
		// Same error Go provides if you execute a,b = 1,0; a/b
		panic(math.DivByZeroErr)
	}

	// Use builtin operators when upper dividend = 0
	if upperDE == 0 {
		lowerQ, remainder = lowerDE/divisor, lowerDE%divisor
		return
	}

	// Phase 1: Find largest multiple of divisor <= dividend.
	// Use bit shifting when upper dividend > 0.
	// Since remainder can be 128-bits for some of the initial subtractions, use dividend parameters for it until we're done.
	// Start with shifting until multiple > dividend (shift while <=).
	// The final shift may require an extra 129th bit that we store in leftCarry.
	var carry, upperM, lowerM, upperF, lowerF uint64 = 0, 0, divisor, 0, 1
	for (carry == 0) && ((upperM < upperDE) || ((upperM == upperDE) && (lowerM <= lowerDE))) {
		carry, upperM, lowerM = Lsh(upperM, lowerM)
		_, upperF, lowerF = Lsh(upperF, lowerF)
	}

	// Stopped at multiple > dividend, bring back one shift, adding carry to the left in case an extra 129th bit was produced
	upperM, lowerM, _ = Rsh(upperM, lowerM)
	upperM |= (carry << 63)
	upperF, lowerF, _ = Rsh(upperF, lowerF)

	// Phase2: Subtract multiples and shift until multiple < divisor.
	// Subtract current multiple from dividend to get new dividend (effectively, new remainder).
	// It is possible that after subtracting this first multiple, we are done.
	upperDE, lowerDE = Sub(upperDE, lowerDE, upperM, lowerM)

	// Add factor to quotient - since quotient is zero, just set it
	upperQ = upperF
	lowerQ = lowerF

	// Continue searching for more multiples to subtract until dividend (current remainder) < divisor
	for (upperDE > 0) || (lowerDE >= divisor) {
		// Find next multiple to subtract from remainder
		for (upperM > upperDE) || (lowerM > lowerDE) {
			upperM, lowerM, _ = Rsh(upperM, lowerM)
			upperF, lowerF, _ = Rsh(upperF, lowerF)
		}

		// Subtract multiple from dividend (current remainder)
		upperDE, lowerDE = Sub(upperDE, lowerDE, upperM, lowerM)

		// Add factor to quotient
		upperQ += upperF
		lowerQ += lowerF
	}

	// Copy final dividend (final remainder) to remainder output.
	// Since the divisor is 64 bits, this final remainder must be 64 bits.
	remainder = lowerDE

	return
}
