package math

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"
	"math/cmplx"
	"reflect"
	"strings"

	"github.com/bantling/micro/constraint"
	"github.com/bantling/micro/conv"
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

// CmpOrdered compares two ordered types and returns -1, 0, 1, depending on whether val1 is <, =, or > val2.
func CmpOrdered[T constraint.Ordered](val1, val2 T) (res int) {
	if val1 < val2 {
		res = -1
	} else if val1 > val2 {
		res = 1
	}

	return
}

// CmpComplex compares two complex types and returns -1, 0, 1, depending on whether val1 is <, =, or > val2.
func CmpComplex[T constraint.Complex](val1, val2 T) (res int) {
	var (
		abs1 = cmplx.Abs(complex128(val1))
		abs2 = cmplx.Abs(complex128(val2))
	)
	if abs1 < abs2 {
		res = -1
	} else if abs1 > abs2 {
		res = 1
	}

	return
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
	conv.IntegerToBigInt(mp, &mpBI)
	conv.IntegerToBigInt(*ma, &maBI)
	maBI.Mul(mpBI, maBI)

	if err := conv.To(maBI, ma); err != nil {
		if maBI.Sign() >= 0 {
			return OverflowErr
		}

		return UnderflowErr
	}

	return nil
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

// MinOrdered returns the minimum value of two ordered types
func MinOrdered[T constraint.Ordered](val1, val2 T) T {
	if val1 > val2 {
		return val2
	}

	return val1
}

// MinComplex returns the minimum value of two complex types
func MinComplex[T constraint.Complex](val1, val2 T) T {
	if cmplx.Abs(complex128(val1)) > cmplx.Abs(complex128(val2)) {
		return val2
	}

	return val1
}

// MinCmp returns the minimum value of two comparable types
func MinCmp[T constraint.Cmp[T]](val1, val2 T) T {
	if val1.Cmp(val2) > 0 {
		return val2
	}

	return val1
}

// MaxOrdered returns the maximum value of two ordered types
func MaxOrdered[T constraint.Ordered](val1, val2 T) T {
	if val1 < val2 {
		return val2
	}

	return val1
}

// MaxComplex returns the maximum value of two complex types
func MaxComplex[T constraint.Complex](val1, val2 T) T {
	if cmplx.Abs(complex128(val1)) < cmplx.Abs(complex128(val2)) {
		return val2
	}

	return val1
}

// MaxCmp returns the maximum value of two comparable types
func MaxCmp[T constraint.Cmp[T]](val1, val2 T) T {
	if val1.Cmp(val2) < 0 {
		return val2
	}

	return val1
}
