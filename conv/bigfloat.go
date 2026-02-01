package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"

	"github.com/bantling/micro/constraint"
	"github.com/bantling/micro/funcs"
)

// IntToBigFloat converts any signed int type into a *big.Float
func IntToBigFloat[T constraint.SignedInteger](ival T, oval **big.Float) {
	prec := uint(math.Ceil(float64(len(IntToString(ival))) * log2Of10))
	*oval = big.NewFloat(0)
	(*oval).SetPrec(prec)
	(*oval).SetInt64(int64(ival))
}

// UintToBigFloat converts any unsigned int type into a *big.Float
func UintToBigFloat[T constraint.UnsignedInteger](ival T, oval **big.Float) {
	*oval = big.NewFloat(0).SetUint64(uint64(ival))
}

// FloatToBigFloat converts any float type into a *big.Float
func FloatToBigFloat[T constraint.Float](ival T, oval **big.Float) error {
	if math.IsNaN(float64(ival)) {
		return fmt.Errorf(errMsg, ival, "NaN", "*big.Float")
	}

	*oval = big.NewFloat(float64(ival))
	return nil
}

// MustFloatToBigFloat is a Must version of FloatToBigFloat
func MustFloatToBigFloat[T constraint.Float](ival T, oval **big.Float) {
	funcs.Must(FloatToBigFloat(ival, oval))
}

// BigIntToBigFloat converts a *big.Int into a *big.Float
func BigIntToBigFloat(ival *big.Int, oval **big.Float) {
	StringToBigFloat(ival.String(), oval)
}

// BigRatToBigFloat converts a *big.Rat to a *big.Float
func BigRatToBigFloat(ival *big.Rat, oval **big.Float) {
	// Use numerator to calculate the precision, shd be accurate since denominator is basically the exponent
	prec := int(math.Ceil(math.Max(float64(53), float64(len(ival.Num().String()))*log2Of10)))
	*oval, _, _ = big.ParseFloat(ival.FloatString(prec), 10, uint(prec), big.ToNearestEven)

	// Set accuracy to exact
	(*oval).SetMode((*oval).Mode())
}

// StringToBigFloat converts a string to a *big.Float
// Returns an error if the string is not a valid float string
func StringToBigFloat(ival string, oval **big.Float) error {
	// A *big.Float is imprecise, but you can set the precision
	// The crude measure we use is the largest of 53 (number of bits in a float64) and ceiling(string length * Log2(10))
	// If every char was a significant digit, the ceiling calculation would be the minimum number of bits required
	var (
		numBits = uint(math.Max(53, math.Ceil(float64(len(ival))*log2Of10)))
		err     error
	)

	if *oval, _, err = big.ParseFloat(ival, 10, numBits, big.ToNearestEven); err != nil {
		return fmt.Errorf(errMsg, ival, ival, "*big.Float")
	}

	return nil
}

// MustStringToBigFloat is a Must version of FloatToBigFloat
func MustStringToBigFloat(ival string, oval **big.Float) {
	funcs.Must(StringToBigFloat(ival, oval))
}
