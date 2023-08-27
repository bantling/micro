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
	"github.com/bantling/micro/funcs"
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

const (
	// decimalMaxScale is the maximum decimal scale, which is also the maximum precision
	decimalMaxScale = 18

	// decimalDefaultScale is the default decimal scale, which is2, since most uses will probably be for money
	decimalDefaultScale = 2

  // decimalMaxValue is the maximum decimal value
	//                 123 456 789 012 345 678
  decimalMaxValue int64 = +999_999_999_999_999_999

  // decimalMinValue is the minimum decimal value
	//                 123 456 789 012 345 678
  decimalMinValue int64 = -999_999_999_999_999_999

	// decimalCheck18SignificantDigits is the smallest value of 18 significant digits
	//                                       123 456 789 012 345 678
	decimalCheck18SignificantDigits int64 = 100_000_000_000_000_000

  // decimalRoundMaxValue is the maximum decimal value that can be rounded up without requiring a 19th digit
	//                        123 456 789 012 345 678
  decimalRoundMaxValue int64 = +999_999_999_999_999_994

  // decimalRoundMinValue is the minimum decimal value that can be rounded down without requiring a 19th digit
	//                             123 456 789 012 345 678
  decimalRoundMinValue int64 = -999_999_999_999_999_994

  // errScaleTooLargeMsg is the error message for a decimal scale value that is too large
	errScaleTooLargeMsg = "The Decimal scale %d is too large: the value must be <= 18"

  // errValueTooLargeMsg is the error message for a decimal value that is too large
	errValueTooLargeMsg = "The Decimal value %d is too large: the value must be <= 999_999_999_999_999_999"

  // errValueTooSmallMsg is the error message for a decimal value that is too small
	errValueTooSmallMsg = "The Decimal value %d is too small: the value must be >= -999_999_999_999_999_999"

  // errValueTooLargeToRoundMsg is the error message for aligning decimals by rounding up a number too large to round
  errValueTooLargeToRoundMsg = "The decimal value %s is too large to round up"

  // errValueTooSmallToRoundMsg is the error message for aligning decimals by rounding down a number too small to round
  errValueTooSmallToRoundMsg = "The decimal value %s is too small to round down"
)

// Decimal is like SQL Decimal(precision, scale):
// - precision is always 18, the maximum number of decimal digits a signed 64 bit value can store
// - scale is number of digits after decimal place, must be <= 18 (default 2 as most popular use is money)
//
// The zero value is ready to use
type Decimal struct {
  value int64
	scale  uint
}

// OfDecimal creates a Decimal with the given sign, digits, and optional scale (default 0)
func OfDecimal(value int64, scale ...uint) (d Decimal, err error) {
	scaleVal := funcs.SliceIndex(scale, 0, decimalDefaultScale)
	if scaleVal > decimalMaxScale {
		err = fmt.Errorf(errScaleTooLargeMsg, scaleVal)
		return
	}

  if value > decimalMaxValue {
    err = fmt.Errorf(errValueTooLargeMsg, value)
    return
  }

  if value < decimalMinValue {
    err = fmt.Errorf(errValueTooSmallMsg, value)
    return
  }

  d.value = value
	d.scale = scaleVal
	return
}

// MustDecimal is a must version of OfDecimal
func MustDecimal(value int64, scale ...uint) Decimal {
  return funcs.MustValue(OfDecimal(value, scale...))
}

// Sign returns the sign of the number:
// -1 if value < 0
//  0 if value = 0
// +1 if value > 0
func (d Decimal) Sign() (sgn int) {
	switch {
  case d.value > 0:
    sgn = 1
  case d.value < 0:
    sgn = -1
  }

  return
}

// String is the Stringer interface
func (d Decimal) String() (str string) {
	// Convert the abs value of the int to a string to start
	conv.To(funcs.Ternary(d.value < 0, -d.value, d.value), &str)

  // Get number of significant digits (length of string)
	numSig := uint(len(str))

	switch {
	// No digits after the decimal point, just an integer
	case d.scale == 0:
		break

		// The number of significant digits is <= the number of decimals. Add leading "0." + (scale - digits) zeros.
	case numSig <= d.scale:
		str = "0." + strings.Repeat("0", int(d.scale - numSig)) + str

	// At least one digit before and after decimal point, insert decimal at appropriate position
	default:
		numDigitsBeforeDecimal := numSig - d.scale
		str = str[:numDigitsBeforeDecimal] + "." + str[numDigitsBeforeDecimal:]
	}

  // Add a leading minus if negative
  if d.value < 0 {
    str = "-" + str
  }

	return
}

// AdjustDecimalScale adjusts the scale of d1 and d2:
// - If both numbers have the same scale, no adjustment is made
// - Otherwise, the number with the smaller scale is usually adjusted to the same scale as the other number
// - Increasing the scale can cause some most significant digits to be lost, in which case the other number is rounded
//   down to match the scale
//
// Examples:
//
// 1.5 and 1.25 -> 1.50 and 1.25
// 1.5 and 18 digits with no decimals -> 18 digits cannot increase scale, so round 1.5 to 2
// 99_999_999_999_999_999.5 and 1 -> the 18 digits round to a 19 digit value, an error occurs
func AdjustDecimalScale(d1, d2 *Decimal) error {
  if d1.scale == d2.scale {
    return nil
  }

  // Swap if necessary so that d1 has larger scale
  if d1.scale < d2.scale {
    t := d1
    d1 = d2
    d2 = t
  }

  // Convert d1 and d2 to strings of digits only, to see how many significant digits they possess
  var str1, str2 string
  conv.To(funcs.Ternary(d1.value >= 0, d1.value, -d1.value), &str1)
  conv.To(funcs.Ternary(d2.value >= 0, d2.value, -d2.value), &str2)

  var (
    len1 = len(str1)
    len2 = len(str2)
    d2Capacity = decimalMaxScale - len2
    scaleDiff = int(d1.scale - d2.scale)
  )

  // Does d2 have enough remaining capacity for the required trailing zeroes to increase the scale?
  if d2Capacity >= scaleDiff {
    // Easy solution - multiply d2 by 10 ^ scaleDiff, and set scale to match d1
    for i := 0; i < scaleDiff; i++ {
      d2.value *= 10
    }

    d2.scale = d1.scale
  } else {
    // Harder solution - round d1 away from 0 to the same scale as d2, and set scale to match d2

    // First check if d2 value can actually be rounded
    if d2.value > decimalRoundMaxValue {
      return fmt.Errorf(errValueTooLargeToRoundMsg, d2.String())
    }

    if d2.value < decimalRoundMinValue {
      return fmt.Errorf(errValueTooSmallToRoundMsg, d2.String())
    }

    // Round by manipulating digits directly in string as a []byte
    dig1 := []byte(str1)

    // If last digit >= 5, then enter decimal rounding loop:
    // - applies only to digits we throw away
    // - as long as prior digits are >= 4, adding 1 makes it >= 5, so round = true
    // - if a digit is < 4 , adding 1 makes it < 5, so round = false and stop
    // - stop if scaleDiff digits have been rounded and round is still true
    // - if resulting round is true, then continue to integer rounding, else stop
    var round bool = dig1[len1-1] >= 5
    for i := len1 - 2; round && (i >= len1 - scaleDiff); i-- {
      round = dig1[i] >= 4
    }
    // throw away decimal digits
    dig1 = dig1[:len1 - scaleDiff]
    len1 = len(dig1)

    // Integer rounding only occurs if final decimal round is true, and affects digits of d1 we're keeping
    // - as long as integer digit = 9, set to 0
    // - if a digit < 9 is encountered, increment and stop
    // - if all digits are 9, add additional 1 digit on left
    var dig byte
    for i := len1 - 1; round && (i >= 0); i-- {
      dig = dig1[i]
      round = dig == 9
      dig1[i] = funcs.Ternary(round, 0, dig + 1)
    }

    // If final round is true, all integer digits are 9, add a leading 1
    //  + strings.Repeat("0", int(scaleDiff))
    str1 = funcs.Ternary(round, "1" + string(dig1), string(dig1))

    // Set d1 scale
    d1.scale = d2.scale

    // Set d1 value (preserving sign)
    neg := d1.value < 0
    conv.To(str1, &d1.value)
    if neg {
      d1.value = -d1.value
    }
  }

  return nil
}

// Add adds
// func (d Decimal) Add(x, y Decimal) Decimal {
//
// }
