package bcd

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"regexp"
	"strings"

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

	// numberAddDecimalsErr is an attempt to increase the number of decimals where leading significant digits would be lost
	numberAddDecimalsMsg = "Cannot convert %s to %d decimal(s), as significant leading digits would be lost"

  // numberAddOverflowMsg is an attempt to add two non-negative numbers that requires a 17th digit
  numberAddOverflowMsg = "Overflow adding %s to %s"

  // numberAddUnderflowMsg is an attempt to add two negative numbers that requires a 17th digit
  numberAddUnderflowMsg = "Underflow adding %s to %s"

	// highestDigitMask is the bit mask to read the highest digit
	highestDigitMask uint64 = 0xF0_00_00_00_00_00_00_00

  // highestDigitShift is the amount of shifting required to read the decimal value
  highestDigitShift = 60

	// lowestDigitMask is the bit mask to read the lowest digit
	lowestDigitMask uint64 = 0xF

	// allBitsMask is the bit mask for all bits (useful for subtracting out bits you don't want to exclude digits)
	allBitsMask uint64 = 0xFF_FF_FF_FF_FF_FF_FF_FF // all 64 bits set
)

// Number is 16 decimal digits of precision, with two digits stored in each of the 8 bytes of a uint64.
// The decimal place can be located before the first digit(0), between any two digits, or after the last digit (16).
type Number struct {
	sign     Sign
	digits   uint64
	decimals uint
}

// == Constructors

// ofHexInternal constructs a Number where we know that no error can occur
func ofHexInternal(psign Sign, digits uint64, decimals uint) Number {
  res := Number{psign, digits, decimals}

  // Adjust sign
  switch {
  case digits == 0:
    res.sign = Zero
  case digits > 0:
    res.sign = funcs.Ternary(psign == Zero, Positive, psign)
  }

  return res
}

// OfHex constructs a Number from a sign, uint64 bcd, and number of digits that come after the decimal place.
// The most convenient and readable way to specify the digits is to use hex of the form 0x1_234 for the digits 1234.
// If the sign is non-Zero and the digits passed are 0, the sign is adjusted to Zero.
// If the sign is Zero and the digits passed are non-0, the sign is adjusted to Positive.
//
// Returns an error if:
// - the number of decimals > 16
// - the provided uint64 has any digits outside the decimal range of 0 - 9.
func OfHex(psign Sign, digits uint64, decimals uint) (Number, error) {
	var (
		zv   Number
		sign = psign
	)

	if decimals > 16 {
		return zv, fmt.Errorf(numberDecimalsErrMsg, decimals)
	}

	// Test if each digit is decimal
	for check := digits; check != 0; check = check >> 4 {
		if (check & lowestDigitMask) > 9 {
			return zv, fmt.Errorf(numberDigitsErrMsg, digits)
		}
	}

	// Adjust the sign to zero in the result
  return ofHexInternal(sign, digits, decimals), nil
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

	// Get the sign
	sign := OfSign(signStr)

	// Populate the digits from left to right
	var digits uint64
	for _, d := range numFrac {
		digits = (digits << 4) | uint64(d-'0')
	}

	// The number of decimal digits is the length of the fractional string - 1 for the dot (0 if there is no fractional string)
	decimals := uint(funcs.Ternary(fracStr == "", 0, len(fracStr)-1))

	// Return the number representation, with the sign adjusted to zero
	return OfHex(sign, digits, decimals)
}

// == Operations

// String is Stringer interface
func (s Number) String() string {
	var (
		str   strings.Builder
		mask  uint64 = highestDigitMask
		shift        = highestDigitShift
		digit rune
    i uint
	)

  // If the number is negative, start with leading minus sign
  if s.sign == Negative {
    str.WriteRune('-')
  }

	// Search for most significant non-zero digit that comes before the decimal (if any)
	for i = 16 - s.decimals; i > 0; i-- {
    // Get digit value
		digit = rune((s.digits & mask) >> shift)

    // Prepare to get next digit value
		mask >>= 4
		shift -= 4

    // Stop if digit is significant
		if digit > 0 {
			break
		}
	}

  // Is there a significant digit before the decimal point?
  if i == 0 {
    // No, so start with a 0
    str.WriteRune('0')
  } else {
    // Yes, print it and any remaining digits before the decimal, regardless of value
    str.WriteRune(digit + '0')
    for i--; i > 0; i-- {
      // Prior loop always altered mask and shift before terminating
      digit = rune((s.digits & mask) >> shift)
      str.WriteRune(digit + '0')

      mask >>= 4
      shift -= 4
    }
  }

  // Do we have any decimals?
  if s.decimals > 0 {
    // Yes, print a dot, then remaining digits
    str.WriteRune('.')

    for i = s.decimals; i > 0; i-- {
      // Prior loop always altered mask and shift before terminating
      digit = rune((s.digits & mask) >> shift)
      str.WriteRune(digit + '0')

      mask >>= 4
      shift -= 4
    }
  }

	return str.String()
}

// AdjustToZero adjusts the sign of a Number by ensuring that the sign is Zero when the digits are zero.
// No result is returned, the sign of the Number is modified.
func (s *Number) AdjustToZero() {
  s.sign = funcs.Ternary(s.digits == 0, Zero, s.sign)
}

// AdjustedToPositive adjusts the sign of a Number by ensuring that the sign is Positive when the digits are zero.
// Returns the adjusted sign, the Number is unmodified.
func (s Number) AdjustedToPositive() Sign {
  return funcs.Ternary(s.digits == 0, Positive, s.sign)
}

// Negated returns the opposite sign after adjusting to positive.
// Returns the negated sign, the Number is unmodified.
func (s Number) Negated() Sign {
  return funcs.Ternary(s.AdjustedToPositive() == Positive, Negative, Positive)
}

// checkDecimals checks that the two numbers passed have the same number of decimals, returning an error if not
func checkDecimals(a, b Number) error {
	if a.decimals != b.decimals {
		return fmt.Errorf(numberDecimalsDifferMsg, a.decimals, b.decimals)
	}

  return nil
}

// Negate returns the negation of this number
func (s Number) Negate() Number {
  return Number{sign: s.Negated(), digits: s.digits, decimals: s.decimals}
}

// ConvertDecimals converts the number to have the specified number of decimals.
// New decimals < old decimals: a rounding is performed, such that 0-4 are rounded down, 5-9 are rounded up
// New decimals = old decimals: no operation
// New decimals > old decimals: trailing zeros are added by shifting left 4 bits per extra decimal digit required, causing
//
//	the most significant decimal digit to be lost in each shift. The digits lost must be 0.
//
// No over/under flow can occur when rounding, since at least one leading zero is introduced.
// Over/under flow can occur when adding trailing zeros, if a leading non-zero digit is shifted off.
// In this case it is an underflow is the number is negative, else it is an overflow.
//
// An error occurs if:
// - New number of decimals > 16
// - New number of decimals > old number of decimals, and non-zero digits would be lost (over/under flow discussed above)
//
// # When an error occurs, this object is not modified
//
// Examples:
//   - 1.285 is shortened to 2 decimals. Sequence is 1.285 (shift right and round) => 01.29.
//   - 1.295 is shortened to 2 decimals. Sequence is 1.295 (shift right and round) => 01.20 (round in place) => 01.30.
//   - 1.2995 is shortened to 3 decimals. Sequence is 1.2995 (shift right and round) => 01.290 (round in place) => 01.200 (round in place) => 01.300.
//   - 999_999_999_999_999.9 is shortened to 0 decimals. Sequence is 9...9.9 (shift right and round) => 09...90.
//     The rounding ripples across all the 9s, ending with 1_000_000_000_000_000.
//
// - 01.285 is expanded to 4 decimals. Sequence is 01.285 (shift left and add zero) => 1.2850.
func (s *Number) ConvertDecimals(decimals uint) error {
	// There can't be more than 16 digits
	if decimals > 16 {
		return fmt.Errorf(numberDecimalsErrMsg, decimals)
	}

	switch {
	// Shortening number of digits
	case decimals < s.decimals:
		{
			var (
				digit     uint64                   // a single digit to round, typed as uint64 because it can be shifted to any digit position
				mask      uint64 = lowestDigitMask // initial mask to grab digit is last four bits
				maskShift        = 0               // how many bits to shift grabbed digit to line it up on far right, so the value is 0 - 9
				roundNext bool                     // true if this digit rounded from 9 to 0, so that next digit has to be rounded
			)

			// Round 0-4 down, 5-9 up, from right to left, rounding and removing all digits we no longer want
			for i := s.decimals; i > decimals; i-- {
				// Get digit
				digit = s.digits & mask

				// If last digit rounded to a value >= 5, then this digit must be rounded, regardless of value.
				if roundNext {
					digit += 1
				}

				// Round next digit if this digit >= 5
				roundNext = digit >= 5

				// Remove this digit by shifting all digits right once place (4 bits), introducing a 0 in leading position
				s.digits >>= 4
			}

			// If the last removed digit was >= 5:
			// - Always have to round lowest digit regardless of value
			// - Continue rounding remaining digits until a rounded digit is < 9
			// - Modify digits instead of removing them
			for roundNext {
				// Get digit, and shift it to the far right place to examine the value
				digit = (s.digits & mask) >> maskShift

				// Round this digit, 9 becomes 0
				if digit++; digit == 10 {
					digit = 0
				}

				// Replace the digit with the rounded value, whether or not any further rounding occurs
				s.digits = (s.digits & (allBitsMask ^ mask)) | (digit << maskShift)

				// Continue rounding if this digit wrapped around to 0
				if roundNext = digit == 0; roundNext {
					// Prepare mask and shift for next digit
					mask <<= 4
					maskShift += 4
				}
			}
		}

		// Expanding number of digits
	case decimals > s.decimals:
		{
			var (
				mask   uint64 = highestDigitMask // Highest digit only
				digits        = s.digits         // Copy original value in case we have an error, so original remains unmodified
			)

			for i := s.decimals; i < decimals; i++ {
				if (digits & mask) > 0 {
					return fmt.Errorf(numberAddDecimalsMsg, s.String(), decimals)
				}

				digits = digits << 4
			}

			// No error, copy result to receiver
			s.digits = digits
		}
	}

	// Copy new decimals value if no error
	s.decimals = decimals
	return nil
}

// Cmp compares this number against another number, returning:
// +1 = s > n
//  0 = s == n
// -1 = s < n
//
// Returns an error if this number has a different number of decimals than the provided number
func (s Number) Cmp(n Number) (int, error) {
  // Must have same number of decimals
  if err := checkDecimals(s, n); err != nil {
    return 0, err
  }

  // If two numbers are opposite sign, the positive one is always larger (consider zero positive)
  ssgn, nsgn := s.AdjustedToPositive(), n.AdjustedToPositive()
  if ssgn != nsgn {
    return funcs.Ternary(ssgn == Positive, +1, -1), nil
  }

  // If two numbers are the same sign, compare from left to right, stopping at first digit that differs
  var (
    mask = highestDigitMask
    maskShift = 60
    sdigit, ndigit uint64
    cmp int
  )

  for i := 0; i < 16; i++ {
    sdigit, ndigit = s.digits & mask, n.digits & mask

    // The larger digit is the larger magnitude
    switch {
    case sdigit > ndigit:
      cmp = +1

    case sdigit < ndigit:
      cmp = -1
    }

    if cmp != 0 {
      // Consider the sign: if they are positive, sdigit > ndigit means s > n; otherwise s < n
      return funcs.Ternary(ssgn == Positive, cmp, -cmp), nil
    }

    mask >>= 4
    maskShift -= 4
  }

  // Must have all the same digits
  return 0, nil
}

// Add this number to another number, returning a new number with the same number of decimals.
//
// Returns an error if:
// - this number has a different number of decimals than the provided number
// - the addition overflows  (adding two positives is too large)
// - the addition underflows (adding two negatives is too low)
func (s Number) Add(o Number) (Number, error) {
	var zv Number

  // Must have same number of decimals
  if err := checkDecimals(s, o); err != nil {
    return zv, err
  }

  // If adjusted signs differ, it is actually subtraction
  sgn := s.AdjustedToPositive()
  if sgn != o.AdjustedToPositive() {
    return zv, nil
  }

  // Add the digits one column at a time, from right to left.
  // If the result of a column >= 10, subtract 10 for that column, and have a carry of 1 for next column.
  var (
    carry uint64
    mask = lowestDigitMask
    maskShift = 0
    digit uint64
    zero uint64
    one uint64 = 1
    sum uint64 = s.digits
  )

  for i := 0; i < 16; i++ {
    digit = ((sum & mask) >> maskShift) + ((o.digits & mask) >> maskShift) + carry
    if carry = funcs.Ternary(digit >= 10, one, zero); carry == 1 {
      digit -= 10
    }

    sum = (sum & (allBitsMask ^ mask)) | (digit << maskShift)

    mask <<= 4
    maskShift += 4
  }

  // If we have a final carry, that is an overflow (adding non-negatives) or underflow (adding negatives)
  if carry == one {
    return zv, fmt.Errorf(funcs.Ternary(sgn == Positive, numberAddOverflowMsg, numberAddUnderflowMsg), s.String(), o.String())
  }

  // Addition was successful
	return ofHexInternal(s.sign, sum, s.decimals), nil
}

// Subtract another number from this number, returning a new number with the same number of decimals.
//
// Returns an error if:
// - this number has a different number of decimals than the provided number
// - the subtraction overflows (subtracting a negative from a positive is too large)
// - the subtraction underflows (subtracting a positive from a negative is too low)
// func (s Number) Subtract(o Number) (Number, error) {
// 	var zv Number
//
//  // Must have same number of decimals
//  if err := checkDecimals(s, n); err != nil {
//    return err
//  }
//
//   //  9 -  5 = sub 9 - 5 =  4
//   //  5 -  9 = sub 9 - 5 = -4
//   //  9 - -5 = add 9 + 5 =  14
//   //  5 - -9 = add 5 + 9 =  14
//   // -9 -  5 = add 9 + 5 = -14
//   // -5 -  9 = add 5 + 9 =  14
//   // -9 - -5 = sub 9 - 5 = -4
//   // -5 - -9 = sub 9 - 5 =  4
//   sgn := s.sign.AdjustToPositive(s.digits)
//   if sgn != o.sign.AdjustToPositive(o.digits) {
//     return zv, nil
//   }
//
//   // Add the digits one column at a time, from right to left.
//   // If the result of a column >= 10, subtract 10 for that column, and have a carry of 1 for next column.
//   var (
//     carry uint64
//     mask = lowestDigitMask
//     maskShift = 0
//     digit uint64
//     zero uint64
//     one uint64 = 1
//     sum uint64 = s.digits
//   )
//
//   for i := 0; i < 16; i++ {
//     digit = ((sum & mask) >> maskShift) + ((o.digits & mask) >> maskShift) + carry
//     if carry = funcs.Ternary(digit >= 10, one, zero); carry == 1 {
//       digit -= 10
//     }
//
//     sum = (sum & (allBitsMask ^ mask)) | (digit << maskShift)
//
//     mask <<= 4
//     maskShift += 4
//   }
//
//   // If we have a final carry, that is an overflow (adding non-negatives) or underflow (adding negatives)
//   if carry == one {
//     return zv, fmt.Errorf(funcs.Ternary(sgn == Positive, numberAddOverflowMsg, numberAddUnderflowMsg), s.String(), o.String())
//   }
//
//   // Addition was successful
// 	return Number{sign: s.sign.AdjustToZero(sum), digits: sum, decimals: s.decimals}, nil
// }
