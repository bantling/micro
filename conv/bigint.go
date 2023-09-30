package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"

	"github.com/bantling/micro/constraint"
	"github.com/bantling/micro/funcs"
)

// IntToBigInt converts any signed int type into a *big.Int
func IntToBigInt[T constraint.SignedInteger](ival T, oval **big.Int) {
	*oval = big.NewInt(int64(ival))
}

// UintToBigInt converts any unsigned int type into a *big.Int
func UintToBigInt[T constraint.UnsignedInteger](ival T, oval **big.Int) {
	*oval = big.NewInt(0)
	(*oval).SetUint64(uint64(ival))
}

// FloatToBigInt converts any float type to a *big.Int
// Returns an error if the float has fractional digits
func FloatToBigInt[T constraint.Float](ival T, oval **big.Int) error {
	if math.IsInf(float64(ival), 0) || math.IsNaN(float64(ival)) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%g", ival), "*big.Int")
	}

	var inter big.Rat
	inter.SetFloat64(float64(ival))
	if !inter.IsInt() {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%g", ival), "*big.Int")
	}

	*oval = inter.Num()
	return nil
}

// MustFloatToBigInt is a Must version of FloatToBigInt
func MustFloatToBigInt[T constraint.Float](ival T, oval **big.Int) {
	funcs.Must(FloatToBigInt(ival, oval))
}

// BigIntToBigInt makes a copy of a *big.Int such that ival and *oval are different pointers
func BigIntToBigInt(ival *big.Int, oval **big.Int) {
	*oval = big.NewInt(0)
	(*oval).Set(ival)
}

// BigFloatToBigInt converts a *big.Float to a *big.Int.
// Returns an error if the *big.Float has any fractional digits.
func BigFloatToBigInt(ival *big.Float, oval **big.Int) error {
	inter, acc := ival.Rat(nil)
	if (inter == nil) || (!inter.IsInt()) || (acc != big.Exact) {
		return fmt.Errorf(errMsg, ival, ival.String(), "*big.Int")
	}

	*oval = inter.Num()
	return nil
}

// MustBigFloatToBigInt is a Must version of BigFloatToBigInt
func MustBigFloatToBigInt(ival *big.Float, oval **big.Int) {
	funcs.Must(BigFloatToBigInt(ival, oval))
}

// BigRatToBigInt converts a *big.Rat to a *big.Int
// Returns an error if the *big.Rat is not an int
func BigRatToBigInt(ival *big.Rat, oval **big.Int) error {
	if !ival.IsInt() {
		return fmt.Errorf(errMsg, ival, ival.String(), "*big.Int")
	}

	*oval = ival.Num()
	return nil
}

// MustBigRatToBigInt is a Must version of BigFloatToBigInt
func MustBigRatToBigInt(ival *big.Rat, oval **big.Int) {
	funcs.Must(BigRatToBigInt(ival, oval))
}

// StringtoBigInt converts a string to a *big.Int.
// Returns an error if the string is not an integer.
func StringToBigInt(ival string, oval **big.Int) error {
	*oval = big.NewInt(0)
	if _, ok := (*oval).SetString(ival, 10); !ok {
		return fmt.Errorf(errMsg, ival, ival, "*big.Int")
	}

	return nil
}

// MustStringToBigInt is a Must version of StringToBigInt
func MustStringToBigInt(ival string, oval **big.Int) {
	funcs.Must(StringToBigInt(ival, oval))
}
