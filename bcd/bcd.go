package bcd

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"regexp"

	"github.com/bantling/micro/funcs"
)

// SPDX-License-Identifier: Apache-2.0

var (
	// fixedRegex is a regex for a fixed number
	fixedRegex = regexp.MustCompile("^([-+])?([0-9]+)([.][0-9]+)?$")
)

const (
	// fixedErrMsg is the error message for an invalid fixed number input string
	fixedErrMsg = "Invalid fixed decimal string %q: the value must be an optional sign, at least 1 digit, an optional dot and at least one digit, with no more than 16 digits in total"
)

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

// OfSign constructs a Sign from a string, which may be "-" or "+". Any other string is the same as "+".
// The result will be either Negative or Positive.
func OfSign(s string) Sign {
	return funcs.Ternary(s == negStr, Negative, Positive)
}

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

// Fixed is 16 decimal digits of precision.
// The decimal place can be located before the first digit(0), between any two digits, or after the last digit (16).
type Fixed struct {
	sign     Sign
	digits   uint64
	decimals int
}

// Of constructs a Fixed from a string described by the regex ^([-+])?([0-9]+)([.][0-9]+)?$,
// where the number of digits must <= 16.
//
// Returns (zero value, error) if string does not match above regex.
func Of(str string) (Fixed, error) {
	var zv Fixed

	// Attempt to find the parts, returning an error if the string does not match
	parts := fixedRegex.FindStringSubmatch(str)
	if parts == nil {
		return zv, fmt.Errorf(fixedErrMsg, str)
	}

	// Grab all the parts
	signStr, numStr, fracStr := parts[1], parts[2], parts[3]

	// Combine all digits - if fracStr is empty use fracStr[0:], otherwise fracStr[1:] to skip dot
	numFrac := []rune(numStr + fracStr[funcs.Ternary(len(fracStr) == 0, 0, 1):])

	// There can't be more than 16 digits
	if len(numFrac) > 16 {
		return zv, fmt.Errorf(fixedErrMsg, str)
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
	decimals := funcs.Ternary(fracStr == "", 0, len(fracStr) - 1)

	// Return the fixed representation
	return Fixed{sign, digits, decimals}, nil
}
