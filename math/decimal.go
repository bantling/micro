package math

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/funcs"
)

var (
	decimalRegex = regexp.MustCompile("(-?)([1-9][0-9]*)(?:.([0-9]+))?")
)

const (
	// decimalMaxScale is the maximum decimal scale, which is also the maximum precision
	decimalMaxScale = 18

	// decimalDefaultScale is the default decimal scale, which is 2, since most uses will probably be for money
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

	// errInvalidStringMsg is the error message for an invalid string to construct a decimal from
	errInvalidStringMsg = "The string value %s is not a valid decimal string"

	// errToBigIntMsg is the error message for converting a Decimal whose value is fractional to a *big.Int
	errToBigIntMsg = "The decimal value %s cannot be converted to a *big.Int"

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

	// errDecimalOverflowMsg is the error message for an overflow
	errDecimalOverflowMsg = "The decimal calculation %s %s %s overflowed"

	// errDecimalUnderflowMsg is the error message for an underflow
	errDecimalUnderflowMsg = "The decimal calculation %s %s %s underflowed"

	// errDecimalDivisionByZeroMsg is the error message for dividing by zero
	errDecimalDivisionByZeroMsg = "The decimal calculation %s / 0 is not allowed"

	// errDecimalDivisorTooLargeMsg is the error message for dividing by a divisor that is larger than the dividend
	errDecimalDivisorTooLargeMsg = "The decimal calculation %s / %d is not allowed, the divisor is larger than the dividend"
)

// Decimal is like SQL Decimal(precision, scale):
// - precision is always 18, the maximum number of decimal digits a signed 64 bit value can store
// - scale is number of digits after decimal place, must be <= 18 (default 2 as most popular use is money)
//
// The zero value is ready to use
type Decimal struct {
	value int64
	scale uint
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

// StringToDecimal creates a Decimal from the given string
// The string must contain no more than 18 significant digits (leading zeros are not allowed), and satisfy the following regex:
// -?[1-9][0-9]*(.[0-9]+)?
func StringToDecimal(value string) (d Decimal, err error) {
	parts := decimalRegex.FindStringSubmatch(value)

	// Error if string doesn't match regex
	// Error if total number of digits > 18
	// indexes : 1 = optional leading minus sign, 2 = required digits before decimal, 3 = optional digits after decimal
	if (parts == nil) || ((len(parts[2]) + len(parts[3])) > 18) {
		err = fmt.Errorf(errInvalidStringMsg, value)
		return
	}

	// Set scale to number of digits after decimal, which may be zero
	d.scale = uint(len(parts[3]))

	// Combine digits before and after decimal into a single string, and convert it to the int64 value
	conv.To(parts[2]+parts[3], &d.value)

	// If there is a leading minus sign, negate the value
	if len(parts[1]) > 0 {
		d.value = -d.value
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
		str = "0." + strings.Repeat("0", int(d.scale-numSig)) + str

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

// Precison returns the total number of digits of a decimal, including trailing zeros.
// Effectively, the length of the decimal as a string, without a minus sign or decimal point.
func (d Decimal) Precision() int {
	// Just use the String() method, remove any minus or decimal point, and get the length
	return len(strings.Replace(strings.Replace(d.String(), "-", "", 1), ".", "", 1))
}

// Scale returns the number of digits after the decimal, if any.
func (d Decimal) Scale() uint {
	return d.scale
}

// Sign returns the sign of the number:
// -1 if value < 0
//
//	0 if value = 0
//
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

// AdjustDecimalScale adjusts the scale of d1 and d2:
//   - If both numbers have the same scale, no adjustment is made
//   - Otherwise, the number with the smaller scale is usually adjusted to the same scale as the other number
//   - Increasing the scale can cause some most significant digits to be lost, in which case the other number is rounded
//     down to match the scale
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
		len1       = len(str1)
		len2       = len(str2)
		d2Capacity = decimalMaxScale - len2
		scaleDiff  = int(d1.scale - d2.scale)
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

		// Round by just examining first decimal place, if <= 4 round integer down, else round integer up
		round := dig1[len1-scaleDiff] >= '5'

		// throw away decimal digits
		dig1 = dig1[:len1-scaleDiff]
		len1 = len(dig1)

		// Integer rounding affects digits of d1 we're keeping
		// - as long as integer digit = 9, set to 0
		// - if a digit < 9 is encountered, increment and stop
		// - if all digits are 9, add additional 1 digit on left
		var dig byte
		for i := len1 - 1; round && (i >= 0); i-- {
			dig = dig1[i]
			round = dig == '9'
			dig1[i] = funcs.Ternary(round, '0', dig+1)
		}

		// If final round is true, all integer digits are 9, add a leading 1
		str1 = funcs.Ternary(round, "1"+string(dig1), string(dig1))

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

// AdjustDecimalFormat adjusts the two decimals strings to have the same number of digits before and the decimal,
// and the same number of digits after the decimal. Leading and trailing zeros are added as needed.
// Since minus has a higher ASCII value than plus, a positive number has a leading slash.
//
// The strings returned are almost comparable: "-2" > "-1".
//
// Examples:
// 30, 5 -> " 30", " 05"
// 1.23, -78.295 -> " 01.230", "-78.295"
func AdjustDecimalFormat(d1, d2 Decimal) (string, string) {
	var (
		str1 = d1.String()
		str2 = d2.String()

		minus1 = str1[0] == '-'
		minus2 = str2[0] == '-'

		// Remove optional leading minus from String(), split into before and after decimal
		parts1 = strings.Split(strings.Replace(str1, "-", "", 1), ".")
		parts2 = strings.Split(strings.Replace(str2, "-", "", 1), ".")

		// Integer parts before decimal, which must exist, and lengths
		int1 = parts1[0]
		int2 = parts2[0]

		li1 = len(int1)
		li2 = len(int2)

		// Fractional parts after decimal, which may exist, and lengths
		frac1 = funcs.SliceIndex(parts1, 1)
		frac2 = funcs.SliceIndex(parts2, 1)

		lf1 = len(frac1)
		lf2 = len(frac2)

		// Maximum integer and fractional lengths
		mi = MaxOrdered(li1, li2)
		mf = MaxOrdered(lf1, lf2)

		// Leading zeros for shorter integer
		lz1 = strings.Repeat("0", mi-li1)
		lz2 = strings.Repeat("0", mi-li2)

		// Trailing zeros for shorter fractional
		tz1 = strings.Repeat("0", mf-lf1)
		tz2 = strings.Repeat("0", mf-lf2)

		// Builders for formatted strings, starting with integer parts, which must exist
		bld1, bld2 strings.Builder
	)

	bld1.WriteRune(funcs.Ternary(minus1, '-', '/'))
	bld2.WriteRune(funcs.Ternary(minus2, '-', '/'))

	bld1.WriteString(lz1)
	bld1.WriteString(int1)

	bld2.WriteString(lz2)
	bld2.WriteString(int2)

	if (lf1 > 0) || len(tz1) > 0 {
		bld1.WriteString(".")
		bld1.WriteString(frac1)
		bld1.WriteString(tz1)
	}

	if (lf2 > 0) || len(tz2) > 0 {
		bld2.WriteString(".")
		bld2.WriteString(frac2)
		bld2.WriteString(tz2)
	}

	return bld1.String(), bld2.String()
}

// Cmp compares d against o, and returns -1, 0, or 1 depending on whether d < o, d = o, or d > o, respectively.
func (d Decimal) Cmp(o Decimal) int {
	// Simplest way is to compare adjusted format strings with plain old string comparison
	da, oa := AdjustDecimalFormat(d, o)

	// Have to account for fact that "-2" > "-1"
	compare := CmpOrdered(AdjustDecimalFormat(d, o))

	return funcs.Ternary((da[0] == '-') && (oa[0] == '-'), -compare, compare)
}

// Negate returns the negation of d.
// If 0 is passed, the result is 0.
func (d Decimal) Negate() Decimal {
	return Decimal{value: -d.value, scale: d.scale}
}

// addDecimal is internal function called by Add and Sub
// For Add, o = origO
// For Sub, o = -origO
// origO is only needed for error messages
func addDecimal(d, origO, o Decimal, op string) (Decimal, error) {
	// Adjust scales to be the same
	var (
		r  = d
		oc = o
	)
	if err := AdjustDecimalScale(&r, &oc); err != nil {
		return Decimal{}, err // Stop if an error occurred
	}

	// Add adjusted values
	r.value += oc.value

	// If signs are the same, result may have overflowed or underflowed
	if rs, ocs := r.Sign(), oc.Sign(); rs == ocs {
		if rs == 1 {
			// Positives overflowed if result is > max allowed
			if r.value > decimalMaxValue {
				return Decimal{}, fmt.Errorf(errDecimalOverflowMsg, d, op, origO)
			}
			// Negatives underflow if the result is < min allowed
		} else if r.value < decimalMinValue {
			return Decimal{}, fmt.Errorf(errDecimalUnderflowMsg, d, op, origO)
		}
	}

	return r, nil
}

// Add adds two decimals together by first adjusting them to the same scale, then adding their values
// Returns an error if:
// - Adjusting the scale produces an error
// - Addition overflows or underflows
func (d Decimal) Add(o Decimal) (Decimal, error) {
	return addDecimal(d, o, o, "+")
}

// MustAdd is amjust version of Add
func (d Decimal) MustAdd(o Decimal) Decimal {
	return funcs.MustValue(d.Add(o))
}

// Sub subtracts o from d by first adjusting them to the same scale, then subtracting their values
// Returns an error if:
// - Adjusting the scale produces an error
// - Subtraction overflows or underflows
func (d Decimal) Sub(o Decimal) (Decimal, error) {
	return addDecimal(d, o, o.Negate(), "-")
}

// MustSub is amjust version of Sub
func (d Decimal) MustSub(o Decimal) Decimal {
	return funcs.MustValue(d.Sub(o))
}

// Mul multiplies d by o, then sets the result scale to (d scale) + (o scale)
// Returns an overflow error if the result > 18 9 digits.
// Returns an underflow error if the result < - 18 9 digits.
func (d Decimal) Mul(o Decimal) (Decimal, error) {
	// Start by just multiplying the two 64-bit values together, and adding their scales
	r := d
	r.value *= o.value
	r.scale += o.scale

	// There are two cases of over/under flow:
	// - operation is not reversible: if r != 0 and r <= max value, then r / o != d
	// - abs(value) > 18 9's
	// It is an overflow if the signs are the same, underflow if they differ
	// Note we must do checks in the order shown above:
	// - The resulting value may be storable in a 64 bit int, but roll over/under, so that it has the opposite sign of what it should be
	if (r.value != 0) && (r.value/o.value != d.value) {
		return Decimal{}, fmt.Errorf(funcs.Ternary(d.Sign() == o.Sign(), errDecimalOverflowMsg, errDecimalUnderflowMsg), d, "*", o)
	}
	if r.value > decimalMaxValue {
		return Decimal{}, fmt.Errorf(errDecimalOverflowMsg, d, "*", o)
	}
	if r.value < decimalMinValue {
		return Decimal{}, fmt.Errorf(errDecimalUnderflowMsg, d, "*", o)
	}

	return r, nil
}

// MustMul is a must version of Mul
func (d Decimal) MustMul(o Decimal) Decimal {
	return funcs.MustValue(d.Mul(o))
}

// DivIntQuoRem divides d by unsigned integer o, and returns (quotient, remainder, error).
// The scale of the quotient and remainder are the same as that of d.
// EG, 100.00 / 3 = 33.33 remainder 0.01.
//
// The divisor o cannot be larger than the dividend d.
//
// Returns a division by zero error if o is zero.
// Returns a divisor too large error if the o > d.value.
func (d Decimal) DivIntQuoRem(o uint) (Decimal, Decimal, error) {
	// If o is 0, return division by zero error
	if o == 0 {
		return Decimal{}, Decimal{}, fmt.Errorf(errDecimalDivisionByZeroMsg, d)
	}

	// If o > d, return divisor too large
	// To tell if o > d, we have to convert d to integer part only by dividing d.value by 10 ^ d.scale
	var intPartOfD int64 = funcs.Ternary(d.value >= 0, d.value, -d.value)
	for i := uint(0); i < d.scale; i++ {
		intPartOfD /= 10
	}
	if int64(o) > intPartOfD {
		return Decimal{}, Decimal{}, fmt.Errorf(errDecimalDivisorTooLargeMsg, d, o)
	}

	// Divide d by o, and set to d scale
	// Remainder calculated as d - (q * o), the only way to get the decimal place correct
	var (
		q = Decimal{scale: d.scale, value: d.value / int64(o)}
		r = d.MustSub(Decimal{scale: d.scale, value: q.value * int64(o)})
	)

	return q, r, nil
}

// MustDivIntQuoRem is a must version of DivIntQuoRem
func (d Decimal) MustDivIntQuoRem(o uint) (Decimal, Decimal) {
	return funcs.MustValue2(d.DivIntQuoRem(o))
}

// DintIntAdd is like DivIntQuoRem, except that it returns a slice of values that add up to d.
// EG, 100.00 / 3 = [33.34, 33.33, 33.33].
// This method just calls DivIntQuoRem and spreads the remainder across the first remainder values returned.
func (d Decimal) DivIntAdd(o uint) ([]Decimal, error) {
	// Get the quotient and remainder, returning (nil, error) if an error is returned
	q, r, e := d.DivIntQuoRem(o)
	if e != nil {
		return nil, e
	}

	// Remainder is just a count of how many values need to be increased by 1
	var (
		rc  = r.value
		res = make([]Decimal, o)
	)
	for i := int64(0); i < int64(o); i++ {
		res[i] = Decimal{scale: d.scale, value: q.value + int64(funcs.Ternary(rc > 0, 1, 0))}
		rc--
	}

	return res, nil
}

// MustDivIntAdd is a must version of DivIntAdd
func (d Decimal) MustDivIntAdd(o uint) []Decimal {
	return funcs.MustValue(d.DivIntAdd(o))
}

// Div is division
//
// Examples:
//
// 1. 5000 / 200
// 5000 / 200 = 25
//
// 2. 500.0 / -200
// 5000 / -200 = -25
// Scale 1 - scale 0 = 1 -> Set scale to 1
// Result is -2.5
//
// 3. -500.0 / 2.00
// -5000 / 200 = -25
// Scale 1 - scale 2 = -1 -> Multiply by 10^1
// Result is -250
//
// 4. 5001 / -200
// 5001 / -200 = -25 r 1
// 1 / 200 -> 1000 (1 * 10^3) / 200 = 5 scale 3 = 0.005
// Result is -25 - 0.005 = -25.005
//
// 5. -500.1 / 200
// -5001 / 200 = -25 r 1
// Scale 1 - scale 0 = 1 -> 25 scale 1 = 2.5
// 1 / 200 -> 1000 (1 * 10^3) / 200 = 5 scale (1 + 3) = 0.0005
// Result is 2.5 + 0.0005 = 2.5005
//
// 6. 5.123 / 0.021
// 5123 / 21 = 243 r 20
//   20 / 21 = 200 (20 * 10^1) / 21 = 9 scale 1 r 11     = 0.9         r 11
//   11 / 21 = 110 (11 * 10^1) / 21 = 5 scale 1 + 1 r 5  = 0.05        r 5
//    5 / 21 = 50 (5 * 10^1) / 21   = 2 scale 1 + 2 r 8  = 0.002       r 8
//    8 / 21 = 80 (8 * 10^1) / 21   = 3 scale 1 + 3 r 17 = 0.000_3     r 17
//   17 / 21 = 170 (17 * 10^1) / 21 = 8 scale 1 + 4 r 2  = 0.000_08    r 2
//    2 / 21 = 200 (2 * 10^2) / 21  = 9 scale 2 + 5 r 11 = 0.000_000_9 r 11
// So a repeating decimal sequence of 952380 -> 243.952380952380952
//
// 7. 1.03075 / 0.25
// 103075 / 25 = 4123
// Scale 5 - scale 2 = 3
// Result is 4.123
//
// 8. 1234567890123456.78 / 2.5
// 123456789012345678 / 25 = 4938271560493827 r 3 
// Scale 2 - scale 1 = 1 -> 4938271560493827 scale 1 = 493827156049382.7
// 3 / 25 = 300 (3 * 10^2) / 25 = 12 scale 2 + 1 = 0.012
// Result is 493827156049382.7 + 0.012 = 493827156049382.712
//
// 9. 1234567890123456.78 / 0.25
// 123456789012345678 / 25 = 4938271560493827 r 3
// Scale 2 - scale 2 = scale 0 -> 4938271560493827
// 3 / 25 = 300 (3 * 10^2) / 25 = 12 scale 2 = 0.12
// Result is 4938271560493827 + 0.12 = 4938271560493827.12
//
// 10. 1234567890123456.78 / 0.00025
// 123456789012345678 / 25 = 4938271560493827 r 3
// Scale 2 - scale 5 = -3 -> Multiply by 10^3
// 4938271560493827000 = 19 digits = overflow
//
// 11. 1 / 100_000_000_000_000_000
// 1 / 100_000_000_000_000_000
// = 100_000_000_000_000_000 (1 * 10^17) / 100_000_000_000_000_000
// = 1 scale 17
// = 0.00000000000000001
//
// 12. 1 / 200_000_000_000_000_000
// 1 / 200_000_000_000_000_000
// = 1_000_000_000_000_000_000 (1 * 10^18) / 200_000_000_000_000_000
// = underflow, as 1 * 10^18 is 19 digits
// Note the answer is storable (0.000_000_000_000_000_005) = 5 * 10^-18 = 5 scale 18
//
// 13. 100_000_000_000_000_000 / 0.1
// = 100_000_000_000_000_000 / 1
// = 100_000_000_000_000_000
// Scale 0 - 1 = -1 = Multiply by 10^1
// = 1_000_000_000_000_000_000
// = overflow
//
// Algorithm:
//
// 1. Divide dividend by divisor
//    Scale = dividend scale - divisor scale
//    
//
// 1. While divisor scale > 1 and divisor % 10 = 0 (divisor has trailing zero fractional digits)
//    (eliminate all divisor trailing zero fractional digits)
//    divisor = divisor / 10
//    divisor scale = divisor scale - 1
//
// 2. If divisor scale > 1 (dividend and divisor both have fractional digits)
//    adjustment = min(dividend scale, divisor scale)
//    dividend scale = dividend scale - adjustment
//    divisor scale = divisor scale - adjustment
//
// 3.
//
// 1. while dividend < divisor, multiply dividend by 10, counting as int scale
// 2. int, frac = dividend / divisor, dividend % divisor
// 3. while
// func (d Decimal) Div(o Decimal) (Decimal, error) {
//   // Division by zero is an error
//   if o.value == 0 {
//     return fmt.Errorf(errDecimalDivisionByZeroMsg, d)
//   }
// 
//   // Start with integer division, ignoring scale
//   quo, rem := d.value / o.value, d.value % o.value
// 
//   // Adjust scale of quo if necessary
// //   if
// }
