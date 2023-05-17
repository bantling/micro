package bcd

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/tuple"
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

// Negate returns opposite sign (Negative -> Positive, Zero -> Zero, Positive -> Negative)
func (s Sign) Negate() Sign {
	switch s {
	case Negative:
		return Positive
	case Positive:
		return Negative
	}

	return Zero
}

// ==== State Type

// State is an enum of Number states
type State uint

const (
	Normal    State = iota // Normal state, number contains useful digits
	Overflow               // Overflow state, digits are same as before operation that would have overflowed
	Underflow              // Underflow state, digits are same as before operation that would have underflowed
)

// String is the Stringer interface for State
func (s State) String() string {
	switch s {
	case Normal:
		return "Normal"
	case Overflow:
		return "Overflow"
	}

	return "Underflow"
}

// ==== Number type

var (
	// numberRegex is a regex for a number
	numberRegex = regexp.MustCompile("^([-+])?([0-9]+)([.][0-9]+)?$")

	// splitMasks is a slice of digit and decimal masks to split a number into integer and fractional parts
	splitMasks = []tuple.Two[uint64, uint64]{
		tuple.Of2[uint64, uint64](0xFF_FF_FF_FF_FF_FF_FF_FF, 0x00_00_00_00_00_00_00_00), //  0 decimals
		tuple.Of2[uint64, uint64](0xFF_FF_FF_FF_FF_FF_FF_F0, 0x00_00_00_00_00_00_00_0F), //  1 decimal
		tuple.Of2[uint64, uint64](0xFF_FF_FF_FF_FF_FF_FF_00, 0x00_00_00_00_00_00_00_FF), //  2 decimals
		tuple.Of2[uint64, uint64](0xFF_FF_FF_FF_FF_FF_F0_00, 0x00_00_00_00_00_00_0F_FF), //  3 decimals
		tuple.Of2[uint64, uint64](0xFF_FF_FF_FF_FF_FF_00_00, 0x00_00_00_00_00_00_FF_FF), //  4 ecimals
		tuple.Of2[uint64, uint64](0xFF_FF_FF_FF_FF_F0_00_00, 0x00_00_00_00_00_0F_FF_FF), //  5 decimals
		tuple.Of2[uint64, uint64](0xFF_FF_FF_FF_FF_00_00_00, 0x00_00_00_00_00_FF_FF_FF), //  6 decimals
		tuple.Of2[uint64, uint64](0xFF_FF_FF_FF_F0_00_00_00, 0x00_00_00_00_0F_FF_FF_FF), //  7 decimals
		tuple.Of2[uint64, uint64](0xFF_FF_FF_FF_00_00_00_00, 0x00_00_00_00_FF_FF_FF_FF), //  8 decimals
		tuple.Of2[uint64, uint64](0xFF_FF_FF_F0_00_00_00_00, 0x00_00_00_0F_FF_FF_FF_FF), //  9 decimals
		tuple.Of2[uint64, uint64](0xFF_FF_FF_00_00_00_00_00, 0x00_00_00_FF_FF_FF_FF_FF), // 10 decimals
		tuple.Of2[uint64, uint64](0xFF_FF_F0_00_00_00_00_00, 0x00_00_0F_FF_FF_FF_FF_FF), // 11 decimals
		tuple.Of2[uint64, uint64](0xFF_FF_00_00_00_00_00_00, 0x00_00_FF_FF_FF_FF_FF_FF), // 12 decimals
		tuple.Of2[uint64, uint64](0xFF_F0_00_00_00_00_00_00, 0x00_0F_FF_FF_FF_FF_FF_FF), // 13 decimals
		tuple.Of2[uint64, uint64](0xFF_00_00_00_00_00_00_00, 0x00_FF_FF_FF_FF_FF_FF_FF), // 14 decimals
		tuple.Of2[uint64, uint64](0xF0_00_00_00_00_00_00_00, 0x0F_FF_FF_FF_FF_FF_FF_FF), // 15 decimals
		tuple.Of2[uint64, uint64](0x00_00_00_00_00_00_00_00, 0xFF_FF_FF_FF_FF_FF_FF_FF), // 16 decimals
	}
)

const (
	// numberDecimalsErrMsg is an invalid number of decimals
	numberDecimalsErrMsg = "Invalid number of decimals %d: the valid range is [0 .. 16]"

	// numberDigitsErrMsg is an invalid number input
	numberDigitsErrMsg = `Invalid Number "0x%X": the value must contain only decimal digits for each hex group`

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

	// numberSubOverflowMsg is an attempt to subtract a negative from a positive that requires a 17th digit
	numberSubOverflowMsg = "Overflow subtracting %s from %s"

	// numberSubUnderflowMsg is an attempt to subtract a positive from a negative that requires a 17th digit
	numberSubUnderflowMsg = "Underflow subtracting %s from %s"

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
	state    State
	stateMsg string
}

// == Constructors

// ofHexInternal constructs a Number where we know that no error can occur
func ofHexInternal(psign Sign, digits uint64, decimals uint) Number {
	res := Number{psign, digits, decimals, Normal, ""}

	// Adjust sign
	switch {
	case digits == 0:
		res.sign = Zero
		res.decimals = 0
	case (digits > 0) && (psign == Zero):
		res.sign = Positive
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

// MustHex is a Must version of OfHex
func MustHex(psign Sign, digits uint64, decimals uint) Number {
	return funcs.MustValue(OfHex(psign, digits, decimals))
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
	return ofHexInternal(sign, digits, decimals), nil
}

// MustString is a Must version of OfString
func MustString(str string) Number {
	return funcs.MustValue(OfString(str))
}

// == Operations

// String is Stringer interface
func (s Number) String() string {
	// If the state is not normal, return the state message
	if s.state != Normal {
		return s.state.String()
	}

	var (
		str   strings.Builder
		mask  uint64 = highestDigitMask
		shift        = highestDigitShift
		digit rune
		i     uint
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

// AdjustedToPositive adjusts the sign of a Number by ensuring that the sign is Positive when the digits are zero.
// Returns the adjusted sign, the Number is unmodified.
func (s Number) AdjustedToPositive() Sign {
	return funcs.Ternary(s.digits == 0, Positive, s.sign)
}

// Negate returns the same digits with the negated sign
// If the state is Overflow or Underflow, the value is returned as is
func (s Number) Negate() Number {
	if s.state != Normal {
		return s
	}

	return Number{sign: s.sign.Negate(), digits: s.digits, decimals: s.decimals}
}

// IsNormal returns true if the number is in the Normal state
func (s Number) IsNormal() bool {
	return s.state == Normal
}

// State returns the state of the number
func (s Number) State() State {
	return s.state
}

// StateMsg returns the state message of the number, which is the empty string for the Normal state
func (s Number) StateMsg() string {
	return s.stateMsg
}

// ConvertDecimals converts the number to have the specified number of decimals.
// New decimals < old decimals: a rounding is performed, such that 0-4 are rounded down, 5-9 are rounded up
// New decimals = old decimals: no operation
// New decimals > old decimals: trailing zeros are added by shifting left 4 bits per extra decimal digit required, causing
// the most significant decimal digit to be lost in each shift. The lost digits must be 0.
//
// No over/under flow can occur when rounding, since at least one leading zero is introduced.
// Over/under flow can occur when adding trailing zeros, if a leading non-zero digit is shifted off.
// In this case it is an underflow is the number is negative, else it is an overflow.
//
// An error occurs if:
// - New number of decimals > 16
// - New number of decimals > old number of decimals, and non-zero digits would be lost (over/underflow discussed above)
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

// alignDecimals aligns the decimal point of two numbers so they are the same:
// - If they are already the same, return them as is
// - Try ConvertDecimals on number with fewer decimals to extend to larger. If an error occurs, round more decimals to fewer.
func alignDecimals(a, b Number) (Number, Number) {
	ad, bd := a.decimals, b.decimals

	switch {
	case ad < bd:
		// Try extending a, adding trailing zeroes
		if a.ConvertDecimals(bd) != nil {
			// Overflowed, shorten b
			b.ConvertDecimals(ad)
		}

	default: // ad > bd
		// Try extending b, adding trailing zeroes
		if b.ConvertDecimals(ad) != nil {
			// Overflowed, shorten a
			a.ConvertDecimals(bd)
		}
	}

	return a, b
}

// Cmp compares this number against another number, returning:
//
// +1 = s > n
//
//	0 = s == n
//
// -1 = s < n
func (s Number) Cmp(n Number) int {
	// The two numbers may have different numbers of decimal places.
	// - Find the position of most significant digit
	// - Adjust the power of the digit based on the number of decimals
	// - If the powers are different, return +1 if larger power is first, -1 if second
	// - If the powers are the same, compare digits from right to left until:
	//   - A difference occurs, returning a +1 if larger digit is first, -1 if second
	//   - All digits are the same, returning 0
	// - Optimisations:
	//   - Same sign, digits, and decimals mean same value
	//   - Different signs mean greater sign is greater number
	//   - Same sign and decimals and different digits means greater digits is greater
	//   - Same sign and different decimals:
	//     - Both digits are 0 means equal
	//     - One digits are 0 and other is not means non-0 digits are greater
	//     - Otherwise, compare digits to determine less than, greater than, or equal

	switch {
	case s == n:
		// Signs, digits, and decimals are all equal
		return 0
	case s.sign < n.sign:
		// This sign less than other sign
		return -1
	case s.sign > n.sign:
		// This sign greater than other sign
		return +1
	case s.decimals == n.decimals:
		// Signs and decimals are equal, digits must differ
		switch {
		case s.digits < n.digits:
			// This value less than other value (unless signs are negative)
			return funcs.Ternary(s.sign == Positive, -1, +1)
		default:
			// This value greater than other value
			return funcs.Ternary(s.sign == Positive, +1, -1)
		}
	}

	// To reach this point, the two numbers:
	// - Are different
	// - Have the same sign
	// - Have different number of decimals
	// - Neither has digits = 0
	// - May have same digits, but cannot be equal, since different number of decimals
	// - May have different digits, but are equal, with different number of trailing zeroes
	// Strategy:
	// - Get the integer parts as a pair of uint64 values
	// - Shift the integers to the right such that the rightmost digit is in the ones column
	// - Compare integers with < and > operators, continue if equal
	// - Get the fractional parts as a pair of uint64 values
	// - Shift the shorter number of decimals left by the difference, so that leftmost digits are in same column
	// - Compare integers with < and > operators, may be equal
	//
	// Example of equality:
	// 0000 0000 0000 01.20
	// 0000 0000 0000 001.2

	// Integer parts, shifting right by number of digits * 4 bits per digit to align rightmost digits with ones column
	smasks, nmasks := splitMasks[s.decimals], splitMasks[n.decimals]
	si, ni := (s.digits&smasks.T)>>(s.decimals*4), (n.digits&nmasks.T)>>(n.decimals*4)

	switch {
	case si < ni:
		return funcs.Ternary(s.sign == Positive, -1, +1)
	case si > ni:
		return funcs.Ternary(s.sign == Negative, +1, -1)
	}

	// Integers must be equal
	// Decimal parts, shifting shorter length left by difference to align leftmost digits
	sf, nf := s.digits&smasks.U, n.digits&nmasks.U
	if s.decimals > n.decimals {
		nf <<= (s.decimals - n.decimals) * 4
	} else {
		sf <<= (n.decimals - s.decimals) * 4
	}

	switch {
	case sf < nf:
		return funcs.Ternary(s.sign == Positive, -1, +1)
	case sf > nf:
		return funcs.Ternary(s.sign == Positive, +1, -1)
	}

	// Must have same logical value
	return 0
}

// Add this number to another number, returning a new number with the same number of decimals as the number with the most decimals.
//
// The state of the result is Overflow if:
// - s is already Overflow, returning s as is
// - s is Normal and o is Overflow, returning s with state = Overflow and same digits
// - the addition overflows, returning s with state = Overflow, and same digits
// - in last two cases, stateMsg = "Overflow adding o to s"
//
// The state of and the result can become Underflow if:
// - s is already Underflow, returning s as is
// - s is Normal and o is Underflow, returning s with state = Underflow and same digits
// - the addition underflows, returning s with state = Underflow, and same digits
// - in last two cases, stateMsg = "Underflow adding o to s"
func (s Number) Add(o Number) Number {
	// If s is not Normal, return as is
	if s.state != Normal {
		return s
	}

	// If o is not Normal, return s with o State and a message
	switch o.state {
	case Overflow:
		return Number{s.sign, s.digits, s.decimals, o.state, fmt.Sprintf(numberAddOverflowMsg, o, s)}
	case Underflow:
		return Number{s.sign, s.digits, s.decimals, o.state, fmt.Sprintf(numberAddUnderflowMsg, o, s)}
	}

	// Align the decimals of the two numbers
	a, b := alignDecimals(s, o)

	//  9 +  5 = add 9 + 5 =  14
	//  5 +  9 = add 5 + 9 =  14
	// -9 + -5 = add 9 + 5 = -14
	// -5 + -9 = add 5 + 9 = -14
	//
	//  9 + -5 = sub 9 - 5 =  4
	//  5 + -9 = sub 5 - 9 = -4
	// -9 +  5 = sub 9 - 5 = -4
	// -5 +  9 = sub 9 - 5 =  4
	// If adjusted signs differ, it is actually subtraction
	asgn, bsgn := a.AdjustedToPositive(), b.AdjustedToPositive()
	if asgn != bsgn {
		// Call sub with the negative number altered to positive
		switch {
		case asgn == Positive:
			return a.Sub(b.Negate())
		default:
			// If this is negative, we also have to negate result, which cannot over/under flow
			return a.Negate().Sub(b).Negate()
		}
	}

	// Split the two numbers into integer and fractional parts

	// Add the digits one column at a time, from right to left.
	// If the result of a column >= 10, subtract 10 for that column, and have a carry of 1 for next column.
	var (
		carry     uint64
		mask      = lowestDigitMask
		maskShift = 0
		digit     uint64
		sum       uint64 = s.digits
	)

	for i := 0; i < 16; i++ {
		// Add next column and any carry from previous column
		digit = ((sum & mask) >> maskShift) + ((b.digits & mask) >> maskShift) + carry

		// If column >= 10, we need to subtract 10 and carry to next column
		if carry = funcs.Ternary[uint64](digit >= 10, 1, 0); carry == 1 {
			digit -= 10
		}

		// Set next digit of sum
		sum = (sum & (allBitsMask ^ mask)) | (digit << maskShift)

		// Next mask and shift value
		mask <<= 4
		maskShift += 4
	}

	// If we have a final carry, that is an overflow (adding positives) or underflow (adding negatives)
	if carry == 1 {
		var (
			state    State
			stateMsg string
		)

		if asgn == Positive {
			state = Overflow
			stateMsg = numberAddOverflowMsg
		} else {
			state = Underflow
			stateMsg = numberAddUnderflowMsg
		}

		return Number{s.sign, s.digits, s.decimals, state, fmt.Sprintf(stateMsg, o, s)}
	}

	// Addition was successful
	return ofHexInternal(a.sign, sum, a.decimals)
}

// Subtract another number from this number, returning a new number with the same number of decimals as the number with the most decimals.
//
// The state of the result is Overflow if:
// - s is already Overflow, returning s as is
// - s is Normal and o is Overflow, returning s with state = Overflow and same digits
// - in last case, stateMsg = "Overflow subtracting o from s"
//
// The state of and the result can become Underflow if:
// - s is already Underflow, returning s as is
// - s is Normal and o is Underflow, returning s with state = Underflow and same digits
// - the subtraction underflows, returning s with state = Underflow, and same digits
// - in last two cases, stateMsg = "Underflow subtracting o from s"
func (s Number) Sub(o Number) Number {
	// If s is not Normal, return as is
	if s.state != Normal {
		return s
	}

	// If o is not Normal, return s with o State and a message
	switch o.state {
	case Overflow:
		return Number{s.sign, s.digits, s.decimals, o.state, fmt.Sprintf(numberSubOverflowMsg, o, s)}
	case Underflow:
		return Number{s.sign, s.digits, s.decimals, o.state, fmt.Sprintf(numberSubUnderflowMsg, o, s)}
	}

	// Align the decimals of the two numbers
	a, b := alignDecimals(s, o)

	//  9 -  5 = sub 9 - 5 =  4
	//  5 -  9 = sub 9 - 5 = -4
	// -9 - -5 = sub 9 - 5 = -4
	// -5 - -9 = sub 9 - 5 =  4
	//
	//  9 - -5 = add 9 + 5 =  14
	//  5 - -9 = add 5 + 9 =  14
	// -9 -  5 = add 9 + 5 = -14
	// -5 -  9 = add 5 + 9 = -14
	//
	// If adjusted signs differ, it is actually addition
	ssgn, osgn := a.AdjustedToPositive(), b.AdjustedToPositive()
	if ssgn != osgn {
		// Call add with the negative number altered to positive
		var r Number
		switch {
		case ssgn == Positive:
			r = a.Add(b.Negate())
		default:
			// If this is negative, we also have to negate result, which may underflow
			r = a.Negate().Add(b).Negate()
		}

		if r.state != Normal {
			var stateMsg string

			if ssgn == Positive {
				r.state = Overflow
				stateMsg = numberSubOverflowMsg
			} else {
				r.sign = Negative
				r.state = Underflow
				stateMsg = numberSubUnderflowMsg
			}

			r.stateMsg = fmt.Sprintf(stateMsg, b, a)
		}

		return r
	}

	// Borrowing requires the smaller magnitude to be subtracted from the larger magnitude.
	// The resulting sign is the same as this sign, unless we have to flip, in which case it is the opposite of this sign.
	var (
		top, bot = a.digits, b.digits
		rsgn     = ssgn
	)
	if top < bot {
		top, bot = b.digits, a.digits
		rsgn = rsgn.Negate()
	}

	// Subtract the digits one column at a time, from right to left.
	// If a column has top digit < bottom digit, start borrowing by adding 10 to top digit.
	// On next column, if it top digit <= bottom digit, add 9 to top digit.
	// Continue until a column has top digit > bottom digit, subtract one from top and stop borrowing.
	// Example:
	//  201
	// -199
	//  002
	//
	// 1 - 9        -> 1 + 10 - 9 -> 2 start borrow
	// 0 - 9 borrow -> 0 +  9 - 9 -> 0 continue borrow
	// 2 - 1 borrow -> 2 -  1 - 1 -> 0 stop borrow
	var (
		borrow    bool // true if borrowing continues to next column
		mask      = lowestDigitMask
		maskShift = 0
		sub       uint64
		topDigit  uint64
		botDigit  uint64
		subDigit  uint64
	)

	for i := 0; i < 16; i++ {
		topDigit = (top & mask) >> maskShift
		botDigit = (bot & mask) >> maskShift

		switch {
		case !borrow:
			switch {
			case topDigit < botDigit:
				// Borrow 10
				subDigit = topDigit + 10 - botDigit
				borrow = true
			default:
				subDigit = topDigit - botDigit
			}

		case borrow:
			switch {
			case topDigit <= botDigit:
				// Borrow 9
				subDigit = topDigit + 9 - botDigit
			default:
				// Subtract additional 1 and stop borrowing
				subDigit = topDigit - 1 - botDigit
				borrow = false
			}
		}

		sub = (sub & (allBitsMask ^ mask)) | (subDigit << maskShift)

		mask <<= 4
		maskShift += 4
	}

	return ofHexInternal(rsgn, sub, a.decimals)
}

// Mul multiplies this number by another number, returning a new number.
// If this number has N decimals and the other number has M decimals, then multiplication produces N+M decimals.
// If this number has N integer digits and the other number has M integer digits:
// - If both numbers < 1, then no integer digits are produced.
// - If one number >= 1, and the other < 1, then fewer integer digits may be produced, possibly none.
// - If both numbers >= 1, then the integer produced is between max(N, M) and N+M digits, inclusive.
//
// If there are not enough decimal digits available to store the resulting number of decimal digits, then the decimal
// digits are rounded to what is available. Decimal digits may also be rounded further to make room for integer digits.
//
// If there are not enough integer digits available to store the resulting number of integer digits, then an overflow
// (positive number too large) or underflow (negative number too low) occurs.
func (s Number) Mul(o Number) Number {

	return Number{}
}
