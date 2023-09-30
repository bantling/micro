package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"

	"github.com/bantling/micro/constraint"
	"github.com/bantling/micro/funcs"
)

// IntToBigRat converts any signed int type into a *big.Rat
func IntToBigRat[T constraint.SignedInteger](ival T, oval **big.Rat) {
	*oval = big.NewRat(1, 1).SetInt64(int64(ival))
}

// UintToBigRat converts any unsigned int type into a *big.Rat
func UintToBigRat[T constraint.UnsignedInteger](ival T, oval **big.Rat) {
	*oval = big.NewRat(1, 1).SetUint64(uint64(ival))
}

// FloatToBigRat converts any float type into a *big.Rat
func FloatToBigRat[T constraint.Float](ival T, oval **big.Rat) error {
	if math.IsInf(float64(ival), 0) || math.IsNaN(float64(ival)) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%g", ival), "*big.Rat")
	}

	var inter *big.Float
	FloatToBigFloat(ival, &inter)
	BigFloatToBigRat(inter, oval)
	return nil
}

// MustFloatToBigRat is a Must version of FloatToBigRat
func MustFloatToBigRat[T constraint.Float](ival T, oval **big.Rat) {
	funcs.Must(FloatToBigRat(ival, oval))
}

// BigIntToBigRat converts a *big.Int into a *big.Rat
func BigIntToBigRat(ival *big.Int, oval **big.Rat) {
	*oval = big.NewRat(1, 1).SetFrac(ival, big.NewInt(1))
}

// BigFloatToBigRat converts a *big.Float into a *big.Rat
func BigFloatToBigRat(ival *big.Float, oval **big.Rat) error {
	if ival.IsInf() {
		return fmt.Errorf(errMsg, ival, ival.String(), "*big.Rat")
	}

	*oval, _ = big.NewRat(1, 1).SetString(BigFloatToString(ival))
	return nil
}

// MustBigFloatToBigRat is a Must version of FloatToBigRat
func MustBigFloatToBigRat(ival *big.Float, oval **big.Rat) {
	funcs.Must(BigFloatToBigRat(ival, oval))
}

// BigRatToBigRat makes a copy of a *big.Rat such that ival and *oval are different pointers
func BigRatToBigRat(ival *big.Rat, oval **big.Rat) {
	*oval = big.NewRat(0, 1)
	(*oval).Set(ival)
}

// StringToBigRat converts a string into a *big.Rat
func StringToBigRat(ival string, oval **big.Rat) error {
	var ok bool
	if *oval, ok = big.NewRat(1, 1).SetString(ival); !ok {
		return fmt.Errorf(errMsg, ival, ival, "*big.Rat")
	}

	return nil
}

// MustStringToBigRat is a Must version of StringToBigRat
func MustStringToBigRat(ival string, oval **big.Rat) {
	funcs.Must(StringToBigRat(ival, oval))
}

// FloatStringToBigRat converts a float string to a *big.Rat.
// Unlike StringToBigRat, it will not accept a ratio string like 5/4.
func FloatStringToBigRat(ival string, oval **big.Rat) error {
	// ensure the string is a float string, and not a ratio
	var err error
	if (ival == "+Inf") || (ival == "-Inf") || (ival == "NaN") {
		return fmt.Errorf("The float string value of %s cannot be converted to *big.Rat", ival)
	}

	if _, _, err = big.NewFloat(0).Parse(ival, 10); err != nil {
		return fmt.Errorf("The float string value of %s cannot be converted to *big.Rat", ival)
	}

	// If it is a float string, cannot fail to be parsed by StringToBigRat
	StringToBigRat(ival, oval)
	return nil
}

// MustFloatStringToBigRat is a Must version of FloatStringToBigRat
func MustFloatStringToBigRat(ival string, oval **big.Rat) {
	funcs.Must(FloatStringToBigRat(ival, oval))
}
