package math

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math/big"
	"regexp"
	"slices"
	"strings"

	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/funcs"
)

// A function to aid in providing 128-bit constants for powers of ten multipled by 5, 2, and 1.
// These constants are used to efficiently perform subtractions to convert a 128-bit number into digits.
// val is used as base 0 input string for big.Int.SetString, which means it can have underscores in it.
func generate128UpperLower(val string) []uint64 {
	var (
		bi, _  = big.NewInt(0).SetString(val, 0)
		txt    = bi.Text(16)
		ulen   = len(txt) - 16 // lower is always 16 hex chars, upper varies from 1 to 14
		ut, lt = txt[:ulen], txt[ulen:]
	)

    // Grab upper (leading 1 to 14 hex chars)
	bi.SetString(ut, 16)
	upper := bi.Uint64()

    // Grab lower (trailing 16 hex chars)
	bi.SetString(lt, 16)
	lower := bi.Uint64()

	//fmt.Printf("123456789012341234567890123456\n%s\n%s%s\n%x%x\n%x %x\n\n", txt, ut, lt, upper, lower, upper, lower)
	return []uint64{upper, lower}
}

var (
	// optional minus sign, zero or more digits, optional dot and zero or more digits.
	decimalRegex = regexp.MustCompile("(-?)([0-9]*)[.]?([0-9]*)")

	// Powers of 10 constants from 10^0 thru 10^18 (scale can be 0 - 18)
	powersOf10 = []int64{
		1,                         //  0
		10,                        //  1
		100,                       //  2
		1_000,                     //  3
		10_000,                    //  4
		100_000,                   //  5
		1_000_000,                 //  6
		10_000_000,                //  7
		100_000_000,               //  8
		1_000_000_000,             //  9
		10_000_000_000,            // 10
		100_000_000_000,           // 11
		1_000_000_000_000,         // 12
		10_000_000_000_000,        // 13
		100_000_000_000_000,       // 14
		1_000_000_000_000_000,     // 15
		10_000_000_000_000_000,    // 16
		100_000_000_000_000_000,   // 17
		1_000_000_000_000_000_000, // 18
	}

    // upper64PowersOf10 is (5, 2, 1) * 128-bit powers of ten split into upper and lower 64 bit values
    // indexes are [upper length (18 .. 1)][5/2/1][upper/lower]
    // eg, [0][0][0] = 18 digits, 5 * 10^17, upper 64 bits
    //     [1][1][1] = 17 digits, 2 * 10^17, lower 64 bits
    upper64PowersOf10 = [][][]uint64{
        // 18 upper digits
        {
            //                     123456789012345678 123456789012345678
            generate128UpperLower("500000000000000000_000000000000000000"),
        },
        {
            //                     123456789012345678 123456789012345678
            generate128UpperLower("200000000000000000_000000000000000000"),
        },
        {
            //                     123456789012345678 123456789012345678
            generate128UpperLower("100000000000000000_000000000000000000"),
        },

        // 17 upper digits
        {
            //                     12345678901234567 123456789012345678
            generate128UpperLower("50000000000000000_000000000000000000"),
        },
        {
            //                     12345678901234567 123456789012345678
            generate128UpperLower("20000000000000000_000000000000000000"),
        },
        {
            //                     12345678901234567 123456789012345678
            generate128UpperLower("10000000000000000_000000000000000000"),
        },

        // 16 upper digits
        {
            //                     1234567890123456 123456789012345678
            generate128UpperLower("5000000000000000_000000000000000000"),
        },
        {
            //                     1234567890123456 123456789012345678
            generate128UpperLower("2000000000000000_000000000000000000"),
        },
        {
            //                     1234567890123456 123456789012345678
            generate128UpperLower("1000000000000000_000000000000000000"),
        },

        // 15 upper digits
        {
            //                     123456789012345 123456789012345678
            generate128UpperLower("500000000000000_000000000000000000"),
        },
        {
            //                     123456789012345 123456789012345678
            generate128UpperLower("200000000000000_000000000000000000"),
        },
        {
            //                     123456789012345 123456789012345678
            generate128UpperLower("100000000000000_000000000000000000"),
        },

        // 14 upper digits
        {
            //                     12345678901234 123456789012345678
            generate128UpperLower("50000000000000_000000000000000000"),
        },
        {
            //                     12345678901234 123456789012345678
            generate128UpperLower("20000000000000_000000000000000000"),
        },
        {
            //                     12345678901234 123456789012345678
            generate128UpperLower("10000000000000_000000000000000000"),
        },

        // 13 upper digits
        {
            //                     1234567890123 123456789012345678
            generate128UpperLower("5000000000000_000000000000000000"),
        },
        {
            //                     1234567890123 123456789012345678
            generate128UpperLower("2000000000000_000000000000000000"),
        },
        {
            //                     1234567890123 123456789012345678
            generate128UpperLower("1000000000000_000000000000000000"),
        },

        // 12 upper digits
        {
            //                     123456789012 123456789012345678
            generate128UpperLower("500000000000_000000000000000000"),
        },
        {
            //                     123456789012 123456789012345678
            generate128UpperLower("200000000000_000000000000000000"),
        },
        {
            //                     123456789012 123456789012345678
            generate128UpperLower("100000000000_000000000000000000"),
        },

        // 11 upper digits
        {
            //                     12345678901 123456789012345678
            generate128UpperLower("50000000000_000000000000000000"),
        },
        {
            //                     12345678901 123456789012345678
            generate128UpperLower("20000000000_000000000000000000"),
        },
        {
            //                     12345678901 123456789012345678
            generate128UpperLower("10000000000_000000000000000000"),
        },

        // 10 upper digits
        {
            //                     1234567890 123456789012345678
            generate128UpperLower("5000000000_000000000000000000"),
        },
        {
            //                     1234567890 123456789012345678
            generate128UpperLower("2000000000_000000000000000000"),
        },
        {
            //                     1234567890 123456789012345678
            generate128UpperLower("1000000000_000000000000000000"),
        },

        // 9 upper digits
        {
            //                     123456789 123456789012345678
            generate128UpperLower("500000000_000000000000000000"),
        },
        {
            //                     123456789 123456789012345678
            generate128UpperLower("200000000_000000000000000000"),
        },
        {
            //                     123456789 123456789012345678
            generate128UpperLower("100000000_000000000000000000"),
        },

        // 8 upper digits
        {
            //                     12345678 123456789012345678
            generate128UpperLower("50000000_000000000000000000"),
        },
        {
            //                     12345678 123456789012345678
            generate128UpperLower("20000000_000000000000000000"),
        },
        {
            //                     12345678 123456789012345678
            generate128UpperLower("10000000_000000000000000000"),
        },

        // 7 upper digits
        {
            //                     1234567 123456789012345678
            generate128UpperLower("5000000_000000000000000000"),
        },
        {
            //                     1234567 123456789012345678
            generate128UpperLower("2000000_000000000000000000"),
        },
        {
            //                     1234567 123456789012345678
            generate128UpperLower("1000000_000000000000000000"),
        },

        // 6 upper digits
        {
            //                     123456 123456789012345678
            generate128UpperLower("500000_000000000000000000"),
        },
        {
            //                     123456 123456789012345678
            generate128UpperLower("200000_000000000000000000"),
        },
        {
            //                     123456 123456789012345678
            generate128UpperLower("100000_000000000000000000"),
        },

        // 5 upper digits
        {
            //                     12345 123456789012345678
            generate128UpperLower("50000_000000000000000000"),
        },
        {
            //                     12345 123456789012345678
            generate128UpperLower("20000_000000000000000000"),
        },
        {
            //                     12345 123456789012345678
            generate128UpperLower("10000_000000000000000000"),
        },

        // 4 upper digits
        {
            //                     1234 123456789012345678
            generate128UpperLower("5000_000000000000000000"),
        },
        {
            //                     1234 123456789012345678
            generate128UpperLower("2000_000000000000000000"),
        },
        {
            //                     1234 123456789012345678
            generate128UpperLower("1000_000000000000000000"),
        },

        // 3 upper digits
        {
            //                     123 123456789012345678
            generate128UpperLower("500_000000000000000000"),
        },
        {
            //                     123 123456789012345678
            generate128UpperLower("200_000000000000000000"),
        },
        {
            //                     123 123456789012345678
            generate128UpperLower("100_000000000000000000"),
        },

        // 2 upper digits (hex has 2, 1, 1 upper digits)
        {
            //                     12 123456789012345678
            generate128UpperLower("50_000000000000000000"),
        },
        {
            //                     12 123456789012345678
            generate128UpperLower("20_000000000000000000"),
        },
        {
            //                     12 123456789012345678
            generate128UpperLower("10_000000000000000000"),
        },

        // 1 upper digit (hex has 1, 1, 0 upper digits)
        {
            //                     1 123456789012345678
            generate128UpperLower("5_000000000000000000"),
        },
        {
            //                     1 123456789012345678
            generate128UpperLower("2_000000000000000000"),
        },
        nil,
    }
)

const (
	// decimalMaxScale is the maximum decimal scale, which is also the maximum precision
	// range of 64-bit signed int is:
	//   1 234 567 890 123 456 789
	// - 9,223,372,036,854,775,808
	// + 9,223,372,036,854,775,807
	// That's a total of 19 digits, but cannot store 19 9 digits.
	// So we drop back to 18 digits, and we can express values from -18 9s to +18 9s.
	decimalMaxScale = 18

	// decimalMaxValue is the maximum decimal value
	//                       123 456 789 012 345 678
	decimalMaxValue int64 = +999_999_999_999_999_999

	// decimalMinValue is the minimum decimal value
	//                       123 456 789 012 345 678
	decimalMinValue int64 = -999_999_999_999_999_999

	// decimalCheck18SignificantDigits is the smallest value of 18 significant digits
	//                                      123 456 789 012 345 678
	decimalCheck18SignificantDigits int64 = 100_000_000_000_000_000

	// decimalRoundMaxValue is the maximum decimal value that can be rounded up without requiring a 19th digit
	//                            123 456 789 012 345 678
	decimalRoundMaxValue int64 = +999_999_999_999_999_994

	// decimalRoundMinValue is the minimum decimal value that can be rounded down without requiring a 19th digit
	//                            123 456 789 012 345 678
	decimalRoundMinValue int64 = -999_999_999_999_999_994

	// upperHalfMask and lowerHalfMask are bitmasks for the upper and lower 32 bits of a 64 bit value
	upperHalfMask uint64 = 0xFFFF_FFFF_0000_0000
	lowerHalfMask uint64 = 0x0000_0000_FFFF_FFFF

	// Bitmask for lowest bit of a 64 bit unsigned int
	lowestBitMask uint64 = 0x0000_0000_0000_0001

	// Initial power of ten value to use for 128 bit rounding
	round128InitialPowerOfTenUpper64 uint64 = 0x0000_0000_A000_0000

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
	value        int64
	scale        uint
	denormalized bool
}

// OfDecimal creates a Decimal with the given sign, digits, and scale
// For clarity, there is no default scale
func OfDecimal(value int64, scale uint, normalized ...bool) (d Decimal, err error) {
	if scale > decimalMaxScale {
		err = fmt.Errorf(errScaleTooLargeMsg, scale)
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
	d.scale = scale
	d.denormalized = !funcs.SliceIndex(normalized, 0, true)

	// Apply normalization if desired
	d.applyNormalization()
	return
}

// MustDecimal is a must version of OfDecimal
func MustDecimal(value int64, scale uint, normalized ...bool) Decimal {
	return funcs.MustValue(OfDecimal(value, scale, normalized...))
}

// StringToDecimal creates a Decimal from the given string
// The string must contain no more than 18 significant digits, and satisfy the following regex:
// (-?)([0-9]*)(.[0-9]*)?
func StringToDecimal(value string) (d Decimal, err error) {
	parts := decimalRegex.FindStringSubmatch(value)

	// Error if string doesn't match regex
	// Error if total number of digits > 18
	// indexes : 1 = optional leading minus sign, 2 = optional integer digits, 3 = optional decimal digits
	if (parts == nil) || slices.Equal(parts, []string{"", "", "", ""}) || ((len(parts[2]) + len(parts[3])) > 18) {
		err = fmt.Errorf(errInvalidStringMsg, value)
		return
	}

	// Set scale to number of digits after decimal, which may be zero
	d.scale = uint(len(parts[3]))

	// Combine digits before and after decimal into a single string, and convert it to the int64 value
	conv.StringToInt64(parts[2]+parts[3], &d.value)

	// If there is a leading minus sign, negate the value
	if len(parts[1]) > 0 {
		d.value = -d.value
	}

	return
}

// MustStringToDecimal is a must version of StringToDecimal
func MustStringToDecimal(value string) Decimal {
	return funcs.MustValue(StringToDecimal(value))
}

// String is the Stringer interface
func (d Decimal) String() (str string) {
	// Convert the abs value of the int to a string to start
	str = conv.IntToString(funcs.Ternary(d.value < 0, -d.value, d.value))

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

// Normalized returns true if the operations return normalized values
func (d Decimal) Normalized() bool {
	return !d.denormalized
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

	// Swap pointers if necessary so that d1 has larger scale
	if d1.scale < d2.scale {
		t := d1
		d1 = d2
		d2 = t
	}

	// Convert d1 and d2 to strings of digits only, to see how many significant digits they possess
	var str1, str2 string
	str1 = conv.IntToString(funcs.Ternary(d1.value >= 0, d1.value, -d1.value))
	str2 = conv.IntToString(funcs.Ternary(d2.value >= 0, d2.value, -d2.value))

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
		conv.StringToInt64(str1, &d1.value)
		if neg {
			d1.value = -d1.value
		}
	}

	return nil
}

// MustAdjustDecimalScale is a must version of AdjustDecimalScale
func MustAdjustDecimalScale(d1, d2 *Decimal) {
	funcs.Must(AdjustDecimalScale(d1, d2))
}

// AdjustDecimalFormat adjusts the two decimals strings to have the same number of digits before the decimal,
// and the same number of digits after the decimal. Leading and trailing zeros are added as needed.
// A positive number has a leading space.
//
// The strings returned are not directly comparable numerically:
// "-1" > " 1"
// "-2" > "-1"
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

	bld1.WriteRune(funcs.Ternary(minus1, '-', ' '))
	bld2.WriteRune(funcs.Ternary(minus2, '-', ' '))

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

	// Compare the two
	compare := CmpOrdered(da, oa)

	// "-1" > " 2" and "-2" > "-1" and " 2" < "-1"
	// Simple reverse compare in those cases
	ds, os := da[0], oa[0]
	if ((ds == '-') && (os == ' ')) ||
		((ds == '-') && (os == '-')) ||
		((ds == ' ') && (os == '-')) {
		return -compare
	}

	return compare
}

// Negate returns the negation of d.
// If 0 is passed, the result is 0.
func (d Decimal) Negate() Decimal {
	return Decimal{value: -d.value, scale: d.scale, denormalized: d.denormalized}
}

// MagnitudeLessThanOne returns true if the decimal value
// represents a value whose mangitude < 1
func (d Decimal) MagnitudeLessThanOne() bool {
	// If the value is negative, negate it to be positive
	absVal := d.value
	if absVal < 0 {
		absVal = -absVal
	}

	// Use the powersOf10 slice to lookup 10^scale, where scale is the index
	power10 := powersOf10[d.scale]

	// If the absolute value < 10^scale, then all significant digits are the right of the decimal place,
	// which means the Decimal is < 1
	// Examples:
	// if scale = 0 and abs val < 10^0 = 1, then val = 0, which is the only scale 0 value that is < 1.
	// if scale = 1 and abs val < 10^1 = 10, then 0 <= val <= 9, the one digit is right of decimal, value < 1.
	// if scale = 2 and abs val < 10^2 = 100, then 00 <= val <= 99, the two digits are right of decimal, value < 1.
	return absVal < power10
}

// Normalize has two cases:
// value = 0: ensure scale = 0
// scale > 0: eliminate useless trailing zeros
// This can help improve accuracy over successive calculations.
// This method is the only method that ignores the internal normalize field, to allow forcing normalization as desired.
func (d *Decimal) Normalize() {
	if d.value == 0 {
		d.scale = 0
	} else if d.scale > 0 {
		var (
			v   = d.value
			s   = d.scale
			neg = v < 0
		)

		if neg {
			v = -v
		}

		for q, r := v/10, v%10; (s > 0) && (r == 0); q, r = q/10, q%10 {
			v = q
			s--
		}

		if neg {
			v = -v
		}

		d.value = v
		d.scale = s
	}
}

// applyNormalization calls Normalize if the value is 0 or denormalized is true
// Called by other methods that do calculations to maintain the normalization state
func (d *Decimal) applyNormalization() {
	// Skip if not desired
	if (d.value == 0) || (!d.denormalized) {
		d.Normalize()
	}
}

// addDecimal is internal function called by Add and Sub
// For Add, o = origO
// For Sub, o = -origO
// origO is only needed for error messages
// Returns an overflow  error if the result >   18 9 digits
// Returns an underflow error if the result < - 18 9 digits
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

	r.applyNormalization()
	return r, nil
}

// Add adds two decimals together by first adjusting them to the same scale, then adding their values
// Returns an error if:
// - Adjusting the scale produces an error
// - Addition overflows or underflows
//
// See addDecimal
func (d Decimal) Add(o Decimal) (Decimal, error) {
	return addDecimal(d, o, o, "+")
}

// MustAdd is a must version of Add
func (d Decimal) MustAdd(o Decimal) Decimal {
	return funcs.MustValue(d.Add(o))
}

// Sub subtracts o from d by first adjusting them to the same scale, then subtracting their values
// Returns an error if:
// - Adjusting the scale produces an error
// - Subtraction overflows or underflows
//
// See addDecimal
func (d Decimal) Sub(o Decimal) (Decimal, error) {
	return addDecimal(d, o, o.Negate(), "-")
}

// MustSub is a must version of Sub
func (d Decimal) MustSub(o Decimal) Decimal {
	return funcs.MustValue(d.Sub(o))
}

// mul128 uses 128 bits to multiply a pair of 64 bit integers.
// The result cannot overflow, so no error is returned.
func mul128(a, b uint64) (upper, lower uint64) {
	// Split a and b into 32-bit upper and lower halves.
	// Shift the upper halves right 32 bits to align them into the lower 32 bits,
	// so that multiplication can stay within 64 bits.
	var (
		ua = (a & upperHalfMask) >> 32
		la = a & lowerHalfMask

		ub = (b & upperHalfMask) >> 32
		lb = b & lowerHalfMask
	)

	// Perform multiplications of all combinations (none of these can overflow)
	var (
		lalb = la * lb
		laub = la * ub
		ualb = ua * lb
		uaub = ua * ub
	)

	// The above terms line up into four 32-bit sections as follows:
	//   Upper 64 bits     Lower 64 bits
	// Section1 Section2 Section3 Section4
	//                     lalb     lalb
	//            laub     laub
	//            ualb     ualb
	//   uaub     uaub

	// Add lalb and laub bottom 32 bits shifted into upper 32 bits to line up with Section3 (cannot overflow)
	lower = lalb + ((laub & lowerHalfMask) << 32)

	// The above calculation cannot overflow:
	// (la * lb) + (((la * ub) & lowerHalfMask) << 32)
	//
	//   1 * FFFF_FFFF + (((1 * FFFF_FFFF) & lowerHalfMask) << 32)
	// = FFFF_FFFF + ((FFFF_FFFF & lowerHalfMask) << 32)
	// = FFFF_FFFF + (FFFF_FFFF << 32)
	// = FFFF_FFFF + FFFF_FFFF_0000_0000
	// = FFFF_FFFE_0000_0001
	//
	// The problem is as follows:
	// - lb and ub can be at most FFFF_FFFF, as they are 32-bit values
	// - if la = 1, then (((la * ub) & lowerHalfMask) << 32) has max value of FFFF_FFFF_0000_0000
	// - adding FFFF_FFFF to that yields a 64-bit value, no overflow
	// if we increase la, that causes (((la * ub) & lowerHalfMask) << 32) to be a smaller value:
	// - when la > 1, la * ub causes shifting to the left so that there are some zero bits on the right side
	// - when grabbing the bottomm 32 bits, the result is < FFFF_FFFF
	// - when shifting those bits 32 times to the left, the result is < FFFF_FFFF_0000_0000

	// Add ualb bottom 32 bits shifted into upper 32 bits to line up with Section3 (can overflow)
	var temp = lower + ((ualb & lowerHalfMask) << 32)
	if temp < lower {
		// overlow, add 1 to upper 64 bits
		upper++
	}
	lower = temp

	// Add top 32 bits of laub and ualb shifted into lower 32 bits to line up with Section2.
	// Add 64 bits of uaub as is, already aligned with Section1 and Section2.
	// Even though adding 64-bit values can generally overflow, we know the additions result from multiplying two 64 bit
	// values, the result of which cannot exceed 128 bits.
	upper += (laub >> 32) + (ualb >> 32) + uaub

	return
}

// digits36 converts 128 binary bits into 36 decimal digits.
// Since the 128 binary bits derive from multiplying two decimal values, the maximum result comes from multiplying
// 18 9's by 18 9's which equals:
//
// 123456789012345678901234567890123456
// 999999999999999998000000000000000001
//
// For simplicity of accessing the separate digits, each digit is stored in a separate byte
// The dibble-dabble method is used (see https://en.wikipedia.org/wiki/Double_dabble), which works as follows:
//
// - All bytes initialized as zero
// - Shift left 1 bit position, rippling across all bytes
// - Scan all digits, and if any digit is >= 5, add 3 to it (this can result in multiple adds for a single shift)
// - For an input of n bits, then n shifts are required
// - Normally digits are represented in packed BCD form so that each byte has a pair of 4 bit values from 0-9
// - Effectively, it is hex without ever using digits A-F
//
// The algorithm is intended for hardware where accessing each of the 4 bit values and adding 3 can be done in parallel.
// For this implementation, each digit is stored in a separate byte, for ease of access (no bit masking and shifting
// back and forth).
// Example taken for decimal value 65244, a 5 decimal digit 16-bit value:
//     Packed BCD               : Input
//     Dig1 Dig2 Dig3 Dig4 Dig5
// 00. 0000 0000 0000 0000 0000 : 1111 1110 1101 1100 Initial values
//
// 01. 0000 0000 0000 0000 0001 : 1111 1101 1011 1000 Shift
// 02. 0000 0000 0000 0000 0011 : 1111 1011 0111 0000 Shift
// 03. 0000 0000 0000 0000 0111 : 1111 0110 1110 0000 Shift and add Dig5
//     0000 0000 0000 0000 1010
// 04. 0000 0000 0000 0001 0101 : 1110 1101 1100 0000 Shift and add Dig5
//     0000 0000 0000 0001 1000
// 05. 0000 0000 0000 0011 0001 : 1101 1011 1000 0000 Shift
// 06. 0000 0000 0000 0110 0011 : 1011 0111 0000 0000 Shift and add Dig4
//     0000 0000 0000 1001 0011
// 07. 0000 0000 0001 0010 0111 : 0110 1110 0000 0000 Shift and add Dig5
//     0000 0000 0001 0010 1010
// 08. 0000 0000 0010 0101 0100 : 1101 1100 0000 0000 Shift and add Dig4
//     0000 0000 0010 1000 0100
// 09. 0000 0000 0101 0000 1001 : 1011 1000 0000 0000 Shift and add Dig3,Dig5
//     0000 0000 1000 0000 1100
// 10. 0000 0001 0000 0001 1001 : 0111 0000 0000 0000 Shift and add Dig5
//     0000 0001 0000 0001 1100
// 11. 0000 0010 0000 0011 1000 : 1110 0000 0000 0000 Shift and add Dig5
//     0000 0010 0000 0011 1011
// 12. 0000 0100 0000 0111 0111 : 1100 0000 0000 0000 Shift and add Dig4,Dig5
//     0000 0100 0000 1010 1010
// 13. 0000 1000 0001 0101 0101 : 1000 0000 0000 0000 Shift and add Dig2,Dig4,Dig5
//     0000 1011 0001 1000 1000
// 14. 0001 0110 0011 0001 0001 : 0000 0000 0000 0000 Shift and add Dig2
//     0001 1001 0011 0001 0001
// 15. 0011 0010 0110 0010 0010 : 0000 0000 0000 0000 Shift and add Dig3
//     0011 0010 1001 0010 0010
// 16. 0110 0101 0010 0100 0100 : 0000 0000 0000 0000 Shift
// =      6    5    2    4    4
//
// Note the final shift does not perform additions on digits >= 5.
//
// One detail not explained in the article - what if the number requires fewer bits than allowed?
// EG, you have 16 bits for 4 digits, but only a 1 to 3 digit number?
//
// Let's see how to turn 652 into packed BCD when we have 5 digits available:
//     Packed BCD               : Input
//     Dig1 Dig2 Dig3 Dig4 Dig5
// 00. 0000 0000 0000 0000 0000 : 0010 1000 1100 Initial values
//
// 01. 0000 0000 0000 0000 0000 : 0101 0001 1000 Shift
// 02. 0000 0000 0000 0000 0000 : 1010 0011 0000 Shift
// 03. 0000 0000 0000 0000 0001 : 0100 0110 0000 Shift
// 04. 0000 0000 0000 0000 0010 : 1000 1100 0000 Shift
// 05. 0000 0000 0000 0000 0101 : 0001 1000 0000 Shift and add Dig5
//     0000 0000 0000 0000 1000
// 06. 0000 0000 0000 0001 0000 : 0011 0000 0000 Shift
// 07. 0000 0000 0000 0010 0000 : 0110 0000 0000 Shift
// 08. 0000 0000 0000 0100 0000 : 1100 0000 0000 Shift
// 09. 0000 0000 0000 1000 0001 : 1000 0000 0000 Shift and add Dig4
//     0000 0000 0000 1011 0001
// 10. 0000 0000 0001 0110 0011 : 0000 0000 0000 Shift and add Dig4
//     0000 0000 0001 1001 0011
// 11. 0000 0000 0011 0010 0110 : 0000 0000 0000 Shift and add Dig4
//     0000 0000 0011 0010 1001
// 12. 0000 0000 0110 0101 0010 : 0000 0000 0000 Shift
//
// Our initial number 652 is 10 bits in size.
// The next multiple of 4 is 12 bits, so we do 12 shifts.
func digits36(upper, lower uint64) [36]byte {
    // Internally use a uint16 and two uint64 for a total of 36 packed BCD digits
    // This allows for less left shift operations
//     var (
//         hi uint16
//         mid, low uint64
//     )
    return [36]byte{}
}

// Mul calculates d * o using one of two methods:
//
// 1. r = d * o is tried first
// If r = 0, then return 0 scale 0.
// If r / d = o, then if r > max value or < min value, we have a valid result 19 digits in length.
//
//	Round it to 18 with single divide by 10, and add 1 if remainder >= 5.
//	Given max int64 value starts with 92, divide by 10 cannot be max value.
//
// Otherwise, go to method 2 below.
//
// The resulting scale rs is d scale + o scale.
// If rs <= 18, just return r with scale rs.
// Otherwise, round r (rs - 18) times.
// If the result is 0, return 0 scale 0, else return r scale rs.
//
// 2. If r / d != o, we have a case where d * o exceeds bounds of 64 bit integers.
// Split d and o into 32-bit upper/lower pairs, and perform a series of shift and adds
// that generate a 128-bit result.
//
// The 128 bit result is first rounded down to 18 digits, reducing rs by up to 18.
// If rs = 0 and there are more than 18 digitds l
// The result rs must be <= 18, return as is.
func (d Decimal) Mul(o Decimal) (Decimal, error) {
	// Try just multiplying the two 64-bit values together, and see if the result fits in 64 bits
	var (
		dval = d.value
		dpos = dval >= 0

		oval = o.value
		opos = oval >= 0

		rscale = d.scale + o.scale
	)

	if !dpos {
		dval = -dval
	}
	if !opos {
		oval = -oval
	}

	var (
		rval = dval * oval
		rpos = dpos == opos
	)

	// If one or both inputs are 0, nothing more to do
	if (dval != 0) && (oval != 0) {
		// If r / o != d, the result overflowed, and we have to use a different technique
		// Since we already checked rval != 0, we cannot get division by zero
		if rval/oval != dval {
			goto splitHalf
		}

		// The result could be a 19 digit value outside the 18 digit range
		if rval > decimalMaxValue {
			// If the scale > 0 then we can round one time to make it 18 digits
			// Since 19 digit value begins with 92, if we drop a digit and round up,
			// we cannot wind up at 19 digits again
			if rscale == 0 {
				return Decimal{}, fmt.Errorf(funcs.Ternary(d.Sign() == o.Sign(), errDecimalOverflowMsg, errDecimalUnderflowMsg), d, "*", o)
			}

			rmdr := rval % 10
			rval /= 10
			if rmdr >= 5 {
				rval++
			}
			rscale--
		}

		// If scale > 18, we need round until it is 18
		// EG, scale 12 * scale 10 = scale 22
		for rscale > decimalMaxScale {
			rmdr := rval % 10
			rval /= 10
			if rmdr >= 5 {
				rval++
			}
			rscale--
		}
	}
	// Skip the splitHalf algorithm
	goto end

	// Split the two numbers into upper and lower 32 bit halves, and perform a series of multiply and adds to get a
	// 128 bit result. The large result is then rounded down to a 64-bit 18 digit result.
	// If it cannot be rounded down to 64 bits, it is an over/underflow.
splitHalf:

	// If the scale is 0, then we have an integer result that overflows, there is no point in using 128 bit math.
	if rscale == 0 {
		return Decimal{}, fmt.Errorf(funcs.Ternary(dpos == opos, errDecimalOverflowMsg, errDecimalUnderflowMsg), d, "*", o)
	}

	//     upper, lower := mul128(uint64(oval), uint64(dval))

	// Round off digits until one of two results:
	// - A value small enough to fit into 64 bits, which we can return
	// - There are no more decimal places to round, leaving only an integer > 64 bits in size, that is an over/underflow
	//
	// Division by 10 across two 64-bit ints can be performed using a bit shifting algorithm. The idea is as follows:
	// - 121 / 10
	// - Start with 10, 1
	// - Repeatedly shift left (20, 2), (40, 4), ...
	// - Stop when the power of 2 multiple of 10 is the largest multiple <= 121 which is (80, 8)
	// - The multiple 8 is the initial quotient, and 121 - 80 = 41 is the remainder
	// - Now switch to a loop of shift right and subtracts, subtracting only when the multiple of 10 <= remainder
	// - For each subtraction, add multiple of 10 to current quotient
	// - Once the current remainder < 10, we have final result
	// - m, q, r = 80, 8, 41
	// - Shift (80, 8) right = (40, 4)
	// - 40 <= 41, so q, r = (8 + 4, 41 = 40) = (12, 1)
	// - Since r < 10, final result is 121 = 12 * 10 + 1
	//
	// We can optimmize this algorithm by our knowledge that we only need to use this algorithm because we have more
	// 64 bits. So instead of starting at and shifting left, start at the power of 10 in the upper 64 bits where the
	// highest bit is in the middle:
	//
	// (upper m, q) = (0000_0000_1010_0000, 16)
	//
	// To get starting point of bit shifting division:
	// If upper m < upper 64, shift (upper m, q) left until upper m >= upper64. If > upper64, shift (m, q) right once.
	// If upper m > upper 64, shift (upper m, q) right until upper m <= upper64.
	//
	// When shifting right, if the upper m <= 101, then:
	// - The lower m bits have to be shifted right
	// - The highest lower m bit is set to lowest upper m bit
	// - upper m shifted right
	// - eg, when upper m  = 101:
	//   - shift lower m (currently 0) right = 0
	//   - set upper m bit = 1
	//   - shift upper m right = 10
	//   - end result is 0000_0000_0000_0010 1000_0000_000_0000
	//
	// Once the starting point is found
	//     var (
	//         upperQuo := middlePower10
	//         lowerQuo  uint64
	//         upperRmdr uint64
	//         lowerRmdr uint64
	//     )
	//     for (upperQuo > 0) && (rscale > 0) {
	//         temp = upper64 / 10
	//         rmdr = upper64 % 10
	//         if rmdr >= 5 {
	//             temp++
	//         }
	//     }

	// Perform common final operations, regardless of which technique was used to get the result
end:
	if !rpos {
		rval = -rval
	}

	r := Decimal{value: rval, scale: rscale, denormalized: d.denormalized}
	r.applyNormalization()

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
		q = Decimal{scale: d.scale, value: d.value / int64(o), denormalized: d.denormalized}
		r = d.MustSub(Decimal{scale: d.scale, value: q.value * int64(o), denormalized: d.denormalized})
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
		res[i] = Decimal{scale: d.scale, value: q.value + int64(funcs.Ternary(rc > 0, 1, 0)), denormalized: d.denormalized}
		rc--
	}

	return res, nil
}

// MustDivIntAdd is a must version of DivIntAdd
func (d Decimal) MustDivIntAdd(o uint) []Decimal {
	return funcs.MustValue(d.DivIntAdd(o))
}

// Div is the general form of division, suitable for any two Decimal values.
// Integer division is used to generate digits by taking remainders and multiplying them by 10 until they are >= divisor,
// so that the remainder can then be divided further, generating more digits.
//
// Examples:
//
// 1. 5000 / 200
// 5000 / 200 = 25
//
// 2. 500.0 / 200
// 5000 / 200 = 25
// Scale 1 - scale 0 = 1 -> Set scale to 1
// Result is 2.5
//
// 3. 500.0 / 2.00
// 5000 / 200 = 25
// Scale 1 - scale 2 = -1 -> Multiply by 10^1
// Result is 250
//
// 4. 500.1 / 2.00
// 5001 / 200 = 25 r 1
// Scale 1 - scale 2 = -1 -> Multiply by 10^1
// 250 r 10
// 10 / 200 -> 1000 (10 * 10^2) / 200 = 5 scale 2 = 0.05
// Result is 250.05
//
// 5. 5001 / 200
// 5001 / 200 = 25 r 1
// 1 / 200 -> 1000 (1 * 10^3) / 200 = 5 scale 3 = 0.005
// Result is 25 + 0.005 = 25.005
//
// 6. 5001 / -200
// 5001 / -200 = -25 r 1
// 1 / -200 -> 1000 (1 * 10^3) / -200 = -5 scale 3 = -0.005
// Result is -25 + -0.005 = -25.005
//
// 7. -5001 / 200
// -5001 / 200 = -25 r -1
// -1 / 200 -> -1000 (1 * 10^3) / 200 = -5 scale 3 = -0.005
// Result is -25 + -0.005 = -25.005
//
// 8. -5001 / -200
// -5001 / -200 = 25 r -1
// Adjust remainder sign to 1
// 1 / 200 -> 1000 (1 * 10^3) / 200 = 5 scale 3 = 0.005
// Result is 25 + 0.005 = 25.005
//
// 9. -500.1 / 200
// -5001 / 200 = -25 r -1
// Scale 1 - scale 0 = 1 -> -25 scale 1 = -2.5
// -1 / 200 -> -1000 (1 * 10^3) / 200 = -5 scale (1 + 3) = -0.0005
// Result is -2.5 + -0.0005 = -2.5005
//
// 10. 3 / 2
// 3 / 2 = 1 r 1
// 1 / 2 = 10 (1 * 10^1) / 2 = 5 scale 1
// Result is 1.5
//
// 10. 5.123 / 0.021
// 5123 / 21 = 243 r 20
//
//	20 / 21 = 200 (20 * 10^1) / 21 = 9 scale 1 + 0 r 11 = 0.9          r 11
//	11 / 21 = 110 (11 * 10^1) / 21 = 5 scale 1 + 1 r 5  = 0.05         r 5
//	 5 / 21 = 50  (5  * 10^1) / 21 = 2 scale 1 + 2 r 8  = 0.002        r 8
//	 8 / 21 = 80  (8  * 10^1) / 21 = 3 scale 1 + 3 r 17 = 0.000_3      r 17
//	17 / 21 = 170 (17 * 10^1) / 21 = 8 scale 1 + 4 r 2  = 0.000_08     r 2
//	 2 / 21 = 200 (2  * 10^2) / 21 = 9 scale 2 + 5 r 11 = 0.000_000_9  r 11
//
// So a repeating decimal sequence of 952380 -> 243.952380952380952
// After the final 2, the next digit is a 3, which means rounding down
// Final result is still 243.952380952380952
//
// 11. 5 / 9 = 0.555...
// By generating a 19th digit of 5, the result rounds to 0.555_555_555_555_555_556
//
// 12. 1.03075 / 0.25
// 103075 / 25 = 4123
// Scale 5 - scale 2 = 3
// Result is 4.123
//
// 13. 1_234_567_890_123_456.78 / 2.5
// 123_456_789_012_345_678 / 25 = 4_938_271_560_493_827 r 3
// Scale 2 - scale 1 = 1 -> 4_938_271_560_493_827 scale 1 = 493_827_156_049_382.7
// 3 / 25 = 30 (3 * 10^1) / 25 = 1 scale 1 + 1 r 5 = 0.01 r 5
// 5 / 25 = 50 (5 * 10^1) / 25 = 2 scale 1 + 2 r 0 = 0.002
// Result is 493_827_156_049_382.7 + 0.012 = 493_827_156_049_382.712
//
// 14. 1_234_567_890_123_456.78 / 0.25
// 123_456_789_012_345_678 / 25 = 4_938_271_560_493_827 r 3
// Scale 2 - scale 2 = scale 0 -> 4_938_271_560_493_827
// 3 / 25 = 30 (3 * 10^1) / 25 = 1 scale 1 + 0 r 5 = 0.1 r 5
// 5 / 25 = 50 (5 * 10^1) / 25 = 2 scale 1 + 1 r 0 = 0.02
// Result is 4_938_271_560_493_827 + 0.12 = 4_938_271_560_493_827.12
//
// 15. 1_234_567_890_123_456.78 / 0.00025
// 123_456_789_012_345_678 / 25 = 4_938_271_560_493_827 r 3
// Scale 2 - scale 5 = -3 -> Multiply by 10^3
// 4_938_271_560_493_827_000 = 19 digits = overflow
//
// 16. 1 / 100_000_000_000_000_000
// 1 / 100_000_000_000_000_000
// = 100_000_000_000_000_000 (1 * 10^17) / 100_000_000_000_000_000
// = 1 scale 17
// = 0.000_000_000_000_000_01
//
// 17. 1 / 200_000_000_000_000_000
// 1 / 200_000_000_000_000_000
// = 1 * 10^18 / 200_000_000_000_000_000, 1 * 10^18 is too large to store
// = overflow
//
// 18. 100_000_000_000_000_000 / 0.1
// = 100_000_000_000_000_000 / 1
// = 100_000_000_000_000_000
// Scale 0 - 1 = -1 = Multiply by 10^1
// = 1 * 10^18
// = overflow
func (d Decimal) Div(o Decimal) (Decimal, error) {
	// Check if d and o are positive (>= 0)
	// 	fmt.Printf("%s / %s\n", d, o)
	dpos, opos := d.value >= 0, o.value >= 0

	// Make both values positive, for simplicity
	dval, oval := d.value, o.value
	if !dpos {
		dval = -dval
	}
	if !opos {
		oval = -oval
	}

	// Start with plain old division
	q := dval / oval
	r := dval % oval

	// Scale is dividend - divisor, could be negative
	s := int(d.scale - o.scale)

	// If scale is negative, multiply q,r by 10^(-scale), set scale = 0
	// If the multiplications go beyond the limit of digits, we have to error out
	for s < 0 {
		q *= 10
		r *= 10
		s += 1

		if q > decimalMaxValue {
			// Return over/underflow
			return Decimal{}, fmt.Errorf(funcs.Ternary(dpos == opos, errDecimalOverflowMsg, errDecimalUnderflowMsg), d, "/", o)
		}
	}

	// If r != 0, perform successive multiply/divides until r = 0
main_loop:
	for r != 0 {
		// Multiply q,r by 10 until r >= o
		for r < oval {
			nq, nr, ns := q*10, r*10, s+1
			// fmt.Printf("nq = %d, nr = %d, ns = %d\n", nq, nr, ns)

			// If quotient exceeds max value, fall back on previous quotient, as most accurate result we can get
			if nq > decimalMaxValue {
				// fmt.Println("ran out of digits")

				// Get next generated digit, so we know if we should round up the existing quotient
				// Remainder may be < oval, multiply by 10 until it isn't
				for nr < oval {
					nr *= 10
				}
				nr /= oval
				// fmt.Printf("next digit remainder = %d\n", nr)

				// New remainder may be > 10
				// If so, divide by 10 until remainder < 10, to get most significant digit of it
				for nr > 10 {
					nr /= 10
				}
				// fmt.Printf("next digit = %d\n", nr)

				// Is next digit >= 5?
				if nr >= 5 {
					// fmt.Printf("Rounding up")
					q++
				}

				break main_loop
			}

			// If remainder exceeds max value, over/underflow occurs
			if nr > decimalMaxValue {
				// fmt.Println("exceeded max")
				return Decimal{}, fmt.Errorf(funcs.Ternary(dpos == opos, errDecimalOverflowMsg, errDecimalUnderflowMsg), d, "/", o)
			}

			// Copy new quotient, remainder, scale into current
			q = nq
			r = nr
			s = ns
		}

		// Add r / oval to q, set r to r mod oval
		q, r = q+r/oval, r%oval
	}

	// If original signs differed, then result is negative
	if dpos != opos {
		q = -q
	}

	// Return result
	return Decimal{value: q, scale: uint(s), denormalized: d.denormalized}, nil
}

// MustDiv is a must version of Div
func (d Decimal) MustDiv(o Decimal) Decimal {
	return funcs.MustValue(d.Div(o))
}
