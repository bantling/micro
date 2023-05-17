package math

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strings"

	"github.com/bantling/micro/constraint"
	"github.com/bantling/micro/conv"
	// "github.com/bantling/micro/funcs"
)

// Constants
var (
	absErrMsg    = "Absolute value error for %d: there is no corresponding positive value in type %T"
	OverflowErr  = fmt.Errorf("Overflow error")
	UnderflowErr = fmt.Errorf("Underflow error")
	DivByZeroErr = fmt.Errorf("Division by zero error")

	// map strings of type names to func(any) error that perform an abs calculation on a value of the type.
	// no map entries are provided for uints.
	toAbs = map[string]func(any) error{
		"int": func(t any) error {
			i := t.(*int)
			if *i < 0 {
				if *i = -*i; *i < 0 {
					return fmt.Errorf(absErrMsg, *i, *i)
				}
			}

			return nil
		},

		"int8": func(t any) error {
			i := t.(*int8)
			if *i < 0 {
				if *i = -*i; *i < 0 {
					return fmt.Errorf(absErrMsg, *i, *i)
				}
			}

			return nil
		},

		"int16": func(t any) error {
			i := t.(*int16)
			if *i < 0 {
				if *i = -*i; *i < 0 {
					return fmt.Errorf(absErrMsg, *i, *i)
				}
			}

			return nil
		},

		"int32": func(t any) error {
			i := t.(*int32)
			if *i < 0 {
				if *i = -*i; *i < 0 {
					return fmt.Errorf(absErrMsg, *i, *i)
				}
			}

			return nil
		},

		"int64": func(t any) error {
			i := t.(*int64)
			if *i < 0 {
				if *i = -*i; *i < 0 {
					return fmt.Errorf(absErrMsg, *i, *i)
				}
			}

			return nil
		},

		"float32": func(t any) error {
			f := t.(*float32)
			*f = float32(math.Abs(float64(*f)))

			return nil
		},

		"float64": func(t any) error {
			f := t.(*float64)
			*f = math.Abs(*f)

			return nil
		},

		"*big.Int": func(t any) error {
			i := t.(**big.Int)
			(*i).Abs(*i)

			return nil
		},

		"*big.Float": func(t any) error {
			f := t.(**big.Float)
			(*f).Abs(*f)

			return nil
		},

		"*big.Rat": func(t any) error {
			r := t.(**big.Rat)
			(*r).Abs(*r)

			return nil
		},
	}

	// map strings of type names to func(any, any, any) error that perform a div calculation on a value of the type
	toDiv = map[string]func(any, any, any) error{
		"int": func(de, dv, q any) error {
			dvi := any(dv).(int)
			if dvi == 0 {
				return DivByZeroErr
			}

			qi := any(q).(*int)
			qi64 := int64(*qi)
			divSInt(int64(de.(int)), int64(dvi), &qi64)
			*qi = int(qi64)

			return nil
		},

		"int8": func(de, dv, q any) error {
			dvi := any(dv).(int8)
			if dvi == 0 {
				return DivByZeroErr
			}

			qi := any(q).(*int8)
			qi64 := int64(*qi)
			divSInt(int64(de.(int8)), int64(dvi), &qi64)
			*qi = int8(qi64)

			return nil
		},

		"int16": func(de, dv, q any) error {
			dvi := any(dv).(int16)
			if dvi == 0 {
				return DivByZeroErr
			}

			qi := any(q).(*int16)
			qi64 := int64(*qi)
			divSInt(int64(de.(int16)), int64(dvi), &qi64)
			*qi = int16(qi64)

			return nil
		},

		"int32": func(de, dv, q any) error {
			dvi := any(dv).(int32)
			if dvi == 0 {
				return DivByZeroErr
			}

			qi := any(q).(*int32)
			qi64 := int64(*qi)
			divSInt(int64(de.(int32)), int64(dvi), &qi64)
			*qi = int32(qi64)

			return nil
		},

		"int64": func(de, dv, q any) error {
			dvi := any(dv).(int64)
			if dvi == 0 {
				return DivByZeroErr
			}

			qi := any(q).(*int64)
			divSInt(de.(int64), dvi, qi)

			return nil
		},

		"uint": func(de, dv, q any) error {
			dvi := dv.(uint)
			if dvi == 0 {
				return DivByZeroErr
			}

			qi := q.(*uint)
			qi64 := uint64(*qi)
			divUInt(uint64(de.(uint)), uint64(dvi), &qi64)
			*qi = uint(qi64)

			return nil
		},

		"uint8": func(de, dv, q any) error {
			dvi := dv.(uint8)
			if dvi == 0 {
				return DivByZeroErr
			}

			qi := q.(*uint8)
			qi64 := uint64(*qi)
			divUInt(uint64(de.(uint8)), uint64(dvi), &qi64)
			*qi = uint8(qi64)

			return nil
		},

		"uint16": func(de, dv, q any) error {
			dvi := dv.(uint16)
			if dvi == 0 {
				return DivByZeroErr
			}

			qi := q.(*uint16)
			qi64 := uint64(*qi)
			divUInt(uint64(de.(uint16)), uint64(dvi), &qi64)
			*qi = uint16(qi64)

			return nil
		},

		"uint32": func(de, dv, q any) error {
			dvi := dv.(uint32)
			if dvi == 0 {
				return DivByZeroErr
			}

			qi := q.(*uint32)
			qi64 := uint64(*qi)
			divUInt(uint64(de.(uint32)), uint64(dvi), &qi64)
			*qi = uint32(qi64)

			return nil
		},

		"uint64": func(de, dv, q any) error {
			dvi := dv.(uint64)
			if dvi == 0 {
				return DivByZeroErr
			}

			qi := q.(*uint64)
			divUInt(de.(uint64), dvi, qi)

			return nil
		},

		"float32": func(de, dv, q any) error {
			// Let FP handle corner cases like +/- infinity, NaN, -0, and division by zero.
			*(q.(*float32)) = de.(float32) / dv.(float32)
			return nil
		},

		"float64": func(de, dv, q any) error {
			// Let FP handle corner cases like +/- infinity, NaN, -0, and division by zero.
			*(q.(*float64)) = de.(float64) / dv.(float64)
			return nil
		},

		"*big.Int": func(de, dv, q any) error {
			// If q (**big.Int) points to nil, then allocate a quotient
			debi, dvbi, qbip := de.(*big.Int), dv.(*big.Int), q.(**big.Int)
			if *qbip == nil {
				*qbip = big.NewInt(0)
			}
			qbi := *qbip

			// *big.Int panics if you divide by zero
			if dvbi.Sign() == 0 {
				return DivByZeroErr
			}

			r := big.NewInt(0)
			qbi.QuoRem(debi, dvbi, r)
			r.Abs(r)

			halfdv := big.NewInt(0).SetBytes(dvbi.Bytes())
			halfdv.Abs(halfdv)
			halfdv.Rsh(halfdv, 1)

			dvbiEven := big.NewInt(0).SetBytes(dvbi.Bytes())
			dvbiEven.And(dvbiEven, big.NewInt(1))
			if ((dvbiEven.Sign() == 0) && (r.Cmp(halfdv) >= 0)) || ((dvbiEven.Sign() == 1) && (r.Cmp(halfdv) > 0)) {
				if qbi.Sign() >= 0 {
					qbi.Add(qbi, big.NewInt(1))
				} else {
					qbi.Sub(qbi, big.NewInt(1))
				}
			}

			return nil
		},

		"*big.Float": func(de, dv, q any) error {
			// If q (**big.Float) points to nil, then allocate a quotient
			debf, dvbf, qbfp := de.(*big.Float), dv.(*big.Float), q.(**big.Float)
			if *qbfp == nil {
				*qbfp = big.NewFloat(0)
			}
			qbf := *qbfp

			// *big.Float cannot store a NaN, it panics if you try 0/0 or +-Inf/+-Inf
			if ((debf.Sign() == 0) && (dvbf.Sign() == 0)) || (debf.IsInf() && dvbf.IsInf()) {
				return big.ErrNaN{}
			}

			qbf.Quo(debf, dvbf)
			return nil
		},

		"*big.Rat": func(de, dv, q any) error {
			// If q (**big.Rat) points to nil, then allocate a quotient
			debr, dvbr, qbrp := de.(*big.Rat), dv.(*big.Rat), q.(**big.Rat)
			if *qbrp == nil {
				*qbrp = big.NewRat(0, 1)
			}
			qbr := *qbrp

			if dvbr.Sign() == 0 {
				// *big.Rat panics if you divide by zero
				return DivByZeroErr
			}

			qbr.Quo(debr, dvbr)
			return nil
		},
	}
)

// Signed integer division, for toDiv map above
func divSInt(de, dv int64, q *int64) {
	*q = de / dv

	// Calc abs of remainder and of half divisor
	r, halfdv := de%dv, dv/2
	if r < 0 {
		r = -r
	}
	if halfdv < 0 {
		halfdv = -halfdv
	}

	// If divisor is odd and r >= half divisor, or divisor is even and r > half divisor, adjust quotient by one to round
	if (((dv & 1) == 0) && (r >= halfdv)) || (((dv & 1) == 1) && (r > halfdv)) {
		if *q >= 0 {
			*q++
		} else {
			*q--
		}
	}
}

// Unsigned division, for toDiv map above
func divUInt(de, dv uint64, q *uint64) {
	*q = de / dv

	// Calc remainder and hald divisor
	r, halfdv := de%dv, dv/2

	// If divisor is odd and r >= half divisor, or divisor is even and r > half divisor, adjust quotient by one to round
	if (((dv & 1) == 0) && (r >= halfdv)) || (((dv & 1) == 1) && (r > halfdv)) {
		*q++
	}
}

// Abs calculates the absolute value of any numeric type.
// Integer types have a range of -N ... +(N-1), which means that if you try to calculate abs(-N), the result is -N.
// The reason for this is that there is no corresponding +N, that would require more bits.
// In this special case, an error is returned, otherwise nil is returned.
// Note that while the constraint allows unsigned ints for completeness, no operation is performed.
func Abs[T constraint.Numeric](val *T) error {
	typval := reflect.TypeOf(val).Elem().String()

	// No need to calculate absolute value of an unsigned int
	if strings.HasPrefix(typval, "uint") {
		return nil
	}

	return toAbs[typval](val)
}

// AddInt adds two signed integers overwriting the second value, and returns an error if over/underflow occurs.
// Over/underflow occurrs if two numbers of the same sign are added, and the sign of result has changed.
// EG, two positive ints are added to create a result too large to be represented in the same number of bits,
//
//	or two negative ints are added to create a result too small to be represented in the same number of bits.
func AddInt[T constraint.SignedInteger](ae1 T, ae2 *T) error {
	sign1, sign2 := ae1 >= 0, *ae2 >= 0
	*ae2 += ae1
	newSign := *ae2 >= 0

	if (sign1 == sign2) && (sign2 != newSign) {
		// Same sign added, sign of result differs
		if sign2 {
			// Positive wrapped around to negative
			return OverflowErr
		}

		// Negative wrapped around to positive
		return UnderflowErr
	}

	return nil
}

// AddUint adds two unsigned integers overwriting the second value, and returns an error if overflow occurs.
// Overflow occurs if two numbers are added, and the magnitude of result is smaller.
func AddUint[T constraint.UnsignedInteger](ae1 T, ae2 *T) error {
	ae2Orig := *ae2
	*ae2 += ae1

	if *ae2 < ae2Orig {
		return OverflowErr
	}

	return nil
}

// SubInt subtracts two integers as sed = me - sed, overwriting the second value, and returns an error if over/underflow occurs.
// See AddInt.
func SubInt[T constraint.SignedInteger](me T, sed *T) error {
	// Instead of actually subtracting, just add the additive inverse
	*sed = -*sed

	return AddInt(me, sed)
}

// SubUint subtracts two integers as sed = me - sed, overwriting the second value, and returns an error if underflow occurs.
// Underflow occurs if sed > me before the subtraction is performed.
func SubUint[T constraint.UnsignedInteger](me T, sed *T) error {
	*sed = me - *sed
	if *sed > me {
		return UnderflowErr
	}

	return nil
}

// Mul multiples two integers overwriting the second value, and returns an error if over/underflow occurs.
// Over/underflow occurs if the magnitude of the result requires more bits than the type provides.
// Unsigned types can only overflow.
func Mul[T constraint.Integer](mp T, ma *T) error {
	var mpBI, maBI *big.Int
	conv.To(mp, &mpBI)
	conv.To(*ma, &maBI)
	maBI.Mul(mpBI, maBI)

	if err := conv.To(maBI, ma); err != nil {
		if maBI.Sign() >= 0 {
			return OverflowErr
		}

		return UnderflowErr
	}

	return nil
}

// MulU64 multiplies two uint64 values into a pair of uint64 values that represent a 128-bit result.
func MulU64(mp, ma uint64) (upper, lower uint64) {
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
	// The folowing explaanation shows how to multiply two maximum value 16-bit numbers:
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

// Div calculates the quotient of a division operation of a pair of integers or a pair of floating point types.
// In the case of integers, the remainder is handled as follows:
// - If the divisor is even, and abs(remainder) >= abs(divisor) / 2, then increase magnitude of quotient by one.
// - If the divisor is odd,  and abs(remainder) >  abs(divisor) / 2, then increase magnitude of quotient by one.
//
// Increasing the magnitude by one means add one of the quotient is positive, subtract one if negative.
// The purpose of adjusting the magnitude is to get the same result as rounding the floating point calculation.
// Floats are not used since 32 and 64 bit integer values can have a magnitude too large to be expressed accurately
// in a float.
//
// Integer Examples:
// 18 / 4 = 4 remainder 2. Since divisor 4 is even and remainder 2 is     >= (4 / 2 = 2), increment quotient to 5 (round 4.5 up).
// 17 / 5 = 3 remainder 2. Since divisor 5 is odd  and remainder 2 is not >  (5 / 2 = 2), leave quotient as is (round 3.4 down).
func Div[T constraint.Numeric](dividend T, divisor T, quotient *T) error {
	// Cast args to any for functions to accept
	var typ, ade, adv, aq = reflect.TypeOf(dividend).String(), any(dividend), any(divisor), any(quotient)
	return toDiv[typ](ade, adv, aq)
}

// DivBigOps is the BigOps version of Div
func DivBigOps[T constraint.BigOps[T]](dividend T, divisor T, quotient *T) error {
	// Cast args to any for functions to accept
	var typ, ade, adv, aq = reflect.TypeOf(dividend).String(), any(dividend), any(divisor), any(quotient)
	return toDiv[typ](ade, adv, aq)
}
