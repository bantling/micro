package bcd

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"regexp"

	"github.com/bantling/micro/funcs"
)

// ==== Sign type

// Sign is -1 for values < 0, 0 for 0, and 1 for values > 0
type Sign int8

const (
	Negative Sign = iota - 1
	Zero
	Positive
)

var (
	negStr     string = "-"
	zeroPosStr string = ""
)

// == Constructors

// OfSign constructs a Sign from a string, which may be "-" or "+". Any other string is the same as "+".
// The result will be either Negative or Positive.
func OfSign(s string) Sign {
	return funcs.Ternary(s == negStr, Negative, Positive)
}

// == Operations

// String is the Stringer interface for Sign
func (s Sign) String() string {
	return funcs.Ternary(s == Negative, negStr, zeroPosStr)
}

// Adjust adjusts the sign of a uint64 by ensuring that the sign is Zero when the uint64 is zero
func (s *Sign) Adjust(n uint64) {
	if n == 0 {
		*s = Zero
	}
}

// ==== Number type

var (
	// numberRegex is a regex for a number
	numberRegex = regexp.MustCompile("^([-+])?([0-9]+)([.][0-9]+)?$")
)

const (
  // numberDecimalsErrMsg is an invalid number of decimals
  numberDecimalsErrMsg = "Invalid number of decimals %d: the valid range is [0 .. 16]"

  // numberDigitsErrMsg is an invalid number input
  numberDigitsErrMsg = `Invalid Number "%x": the value must contain only decimal digits for each hex group`

	// numberStringErrMsg is an invalid number string input
	numberStringErrMsg = "Invalid Number string %q: the value must be an optional sign, at least 1 digit, an optional dot and at least one digit, with no more than 16 digits in total"

  // numberDecimalsDifferMsg is an invalid pair of numbers to operate on, because the number of decimals differs
  numberDecimalsDifferMsg = "Invalid Number pair: the number of decimals do not match (%d and %d)"
)

// Number is 16 decimal digits of precision, with two digits stored in each of the 8 bytes of a uint64.
// The decimal place can be located before the first digit(0), between any two digits, or after the last digit (16).
type Number struct {
	sign     Sign
	digits   uint64
	decimals uint
}

// == Constructors

// OfHex constructs a Number from a sign, uint64 bcd, and number of digits that come after the decimal place.
// The most convenient and readable way to specify the digits is to use hex of the form 0x1_234 for the digits 1234.
// If the sign is non-Zero and the digits passed are 0, the sign is adjusted to Zero.
//
// Returns an error if the provided uint64 has any digits outside the decimal range of 0 - 9.
func OfHex(psign Sign, digits uint64, decimalsOpt ...uint) (Number, error) {
    var (
      zv Number
      sign = psign
      decimals uint = 6
    )

    if len(decimalsOpt) > 0 {
      if decimals = decimalsOpt[0]; decimals > 16 {
        return zv, fmt.Errorf(numberDecimalsErrMsg, decimals)
      }
    }

    // Test if each digit is decimal
    for check := digits; check != 0; check = check >> 4 {
      if (check & 0xF) > 9 {
        return zv, fmt.Errorf(numberDigitsErrMsg, digits)
      }
    }

    // Adjust the sign
    sign.Adjust(digits)

    return Number{sign: sign, digits: digits, decimals: decimals}, nil
}

// OfString constructs a Number from a string described by the regex ^([-+])?([0-9]+)([.][0-9]+)?$,
// where the number of digits must <= 16.
//
// Returns an error if the string does not match above regex.
func OfString(str string) (Number, error) {
	var zv Number

	// Attempt to find the parts, returning an error if the string does not match
	parts := numberRegex.FindStringSubmatch(str)
	if parts == nil {
		return zv, fmt.Errorf(numberStringErrMsg, str)
	}

	// Grab all the parts
	signStr, numStr, fracStr := parts[1], parts[2], parts[3]

	// Combine all digits - if fracStr is empty use fracStr[0:], otherwise fracStr[1:] to skip dot
	numFrac := []rune(numStr + fracStr[funcs.Ternary(len(fracStr) == 0, 0, 1):])

	// There can't be more than 16 digits
	if len(numFrac) > 16 {
		return zv, fmt.Errorf(numberStringErrMsg, str)
	}

	// Get the sign
	sign := OfSign(signStr)

	// Populate the digits from left to right
	var digits uint64
	for _, d := range numFrac {
		digits = (digits << 4) | uint64(d-'0')
	}

	// Adjust the sign
	sign.Adjust(digits)

	// The number of decimal digits is the length of the fractional string - 1 for the dot (0 if there is no fractional string)
	decimals := uint(funcs.Ternary(fracStr == "", 0, len(fracStr) - 1))

	// Return the number representation
	return Number{sign, digits, decimals}, nil
}

// == Operations

// ConvertDecimals converts the number to have the specified number of decimals.
// New decimals < old decimals: a rounding is performed, such that 0-4 are rounded down, 5-9 are rounded up
// New decimals = old decimals: no operation
// New decimals > old decimals: trailing zeros are added by shifting left 4 bits per extra decimal digit required, causing
//   the most significant decimal digit to be lost in each shift. The digits lost must be 0.
//
// No over/under flow can occur when rounding, since at least one leading zero is introduced.
// Over/under flow can occur when adding trailing zeros, if a leading non-zero digit is shifted off.
// In this case it is an underflow is the number is negative, else it is an overflow.
//
// An error occurs if:
// - New number of decimals > 16
// - New number of decimals > old number of decimals, and non-zero digits would be lost (over/under flow discussed above)
//
// When an error occurs, this object is not modified
//
// Examples:
// - 1.285 is shortened to 2 decimals. Sequence is 1.285 (shift right and round) => 01.29.
// - 1.295 is shortened to 2 decimals. Sequence is 1.295 (shift right and round) => 01.20 (round in place) => 01.30.
// - 1.2995 is shortened to 3 decimals. Sequence is 1.2995 (shift right and round) => 01.290 (round in place) => 01.200 (round in place) => 01.300.
// - 999_999_999_999_999.9 is shortened to 0 decimals. Sequence is 9...9.9 (shift right and round) => 09...90.
//   The rounding ripples across all the 9s, ending with 1_000_000_000_000_000.
//
// - 01.285 is expanded to 4 decimals. Sequence is 01.285 (shift left and add zero) => 1.2850.
func (s *Number) ConvertDecimals(decimals uint) error {
  // There can't be more than 16 digits
  if len(decimals) > 16 {
    return fmt.Errorf(numberDecimalsErrMsg, decimals)
  }

  // Diff in number of decimals is -n if shortening by n digits, +n if expanding by n digits, 0 if no change occurs
  diff := decimals - s.decimals

  switch {
  case (decimals < s.decimals): {
    var (
      digit uint8 // a single digit to round
      mask uint64 = 0xF // initial mask to grab digit is last four bits
      maskShift = 0 // how many bits to shift grabbed digit to line it up on far right, so the value is 0 - 9
      allBits uint64 = 0xFF_FF_FF_FF_FF_FF_FF_FF // all 64 bits set
      bool roundNext // true if this digit rounded from 9 to 0, so that next digit has to be rounded
    )

    // Round 0-4 down, 5-9 up, from right to left, rounding and removing all digits we no longer want
    for i := s.decimals; i > decimals; i--  {
      // Get digit
      digit = uint8(s.digits & mask)

      // If last digit rounded to a value >= 5, then this digit must be rounded, regardless of value.
      if roundNext {
        digit += 1
      }

      // Round next digit if this digit >= 5
      roundNext = digit >= 5

      // Remove this digit by shifting all digits right once place (4 bits), introducing a 0 in leading position
      s.decimals >>= 4
    }

    // If the last removed digit was >= 5:
    // - Always have to round lowest digit regardless of value
    // - Continue rounding remaining digits until a rounded digit is < 9
    // - Modify digits instead of removing them
    for roundNext {
      // Get digit, and shift it to the far right place to examine the value
      digit = uint8((s.digits & mask) >> maskShift)

      // Round this digit, 9 becomes 0
      if digit++; digit == 10 {
        digit = 0
      }

      // Replace the digit with the rounded value, whether or not any further rounding occurs
      s.digit = (s.digits & (allBits ^ mask)) | (digit << maskShift)

      // Continue rounding if this digit wrapped around to 0
      if roundNext = digit == 0; roundNext {
        // Prepare mask and shift for next digit
        mask <<= 4
        maskShift += 4
      }
    }
  }

  case (decimals > s.decimals):
  }
}

// Add this number to another number, returning a new number
//
// Returns an error if:
// - this number has a different number of decimals than the provided number
// - the addition overflows  (adding two positives is too large)
// - the addition underflows (adding two negatives is too low)
func (s Number) Add(o Number) (Number, error) {
  var zv Number

  if (a.decimals != b.decimals) {
    return zv, fmt.Errorf(numberDecimalsDifferMsg, a.decimals, b.decimals)
  }

  return zv, nil
}
