package math

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strings"

	"github.com/bantling/micro/go/constraint"
)

// Constants
var (
	absErrMsg    = "Absolute value error for %d: there is no corresponding positive value in type %T"
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

	// map strings of type names to func(any, any, any) error that perform an div calculation on a value of the type
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
				return DivByZeroErr
			}

			qbr.Quo(debr, dvbr)
			return nil
		},
	}
)

// Signed division, for toDiv map above
func divSInt(de, dv int64, q *int64) {
  *q = de / dv

  r, halfdv := de%dv, dv/2
  if r < 0 {
    r = -r
  }
  if halfdv < 0 {
    halfdv = -halfdv
  }

  if (((dv & 1) == 0) && (r >= halfdv)) || (((dv & 1) == 1) && (r > halfdv)) {
    if *q >= 0 {
      *q++
    } else {
      *q--
    }
  }
}

// Unsigned division, for toDiv map above
divUInt = func(de, dv uint64, q *uint64) {
  *q = de / dv

  r, halfdv := de%dv, dv/2

  if (((dv & 1) == 0) && (r >= halfdv)) || (((dv & 1) == 1) && (r > halfdv)) {
    *q++
  }
}

// Abs calculates the absolute value of any numeric type.
// Integer types have a range of -N ... +(N-1), which means that if you try to calculate abs(-N), the result is -N.
// The reason for this is that there is no corresponding +N, that would require more bits.
// In this special case, an error is returned, otherwise nil is returned.
func Abs[T constraint.Numeric](val *T) error {
	typval := reflect.TypeOf(val).Elem().String()

	// No need to calculate absolute value of an unsigned int
	if strings.HasPrefix(typval, "uint") {
		return nil
	}

	return toAbs[typval](val)
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
