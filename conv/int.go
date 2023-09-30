package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"
	goreflect "reflect"
	"strconv"

	"github.com/bantling/micro/constraint"
	"github.com/bantling/micro/funcs"
)

var (
	minIntValue = map[int]int{
		8:  math.MinInt8,
		16: math.MinInt16,
		32: math.MinInt32,
		64: math.MinInt64,
	}

	maxIntValue = map[int]int{
		8:  math.MaxInt8,
		16: math.MaxInt16,
		32: math.MaxInt32,
		64: math.MaxInt64,
	}

	maxUintValue = map[int]uint{
		8:  math.MaxUint8,
		16: math.MaxUint16,
		32: math.MaxUint32,
		64: math.MaxUint64,
	}
)

// NumBits provides the number of bits of any integer or float type
func NumBits[T constraint.Signed | constraint.UnsignedInteger](val T) int {
	return int(goreflect.ValueOf(val).Type().Size() * 8)
}

// ==== ToInt

// IntToInt converts any signed integer type into any signed integer type
// Returns an error if the source value cannot be represented by the target type
func IntToInt[S constraint.SignedInteger, T constraint.SignedInteger](ival S, oval *T) error {
	var (
		srcSize = NumBits(ival)
		tgtSize = NumBits(*oval)
	)

	if (srcSize > tgtSize) && ((ival < S(minIntValue[tgtSize])) || (ival > S(maxIntValue[tgtSize]))) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%d", ival), fmt.Sprintf("%T", *oval))
	}

	*oval = T(ival)
	return nil
}

// MustIntToInt is a Must version of IntToInt
func MustIntToInt[S constraint.SignedInteger, T constraint.SignedInteger](ival S, oval *T) {
	funcs.Must(IntToInt(ival, oval))
}

// UintToInt converts any unsigned integer type into any signed integer type
// Returns an error if the unsigned int cannot be represented by the signed type
func UintToInt[U constraint.UnsignedInteger, I constraint.SignedInteger](ival U, oval *I) error {
	var (
		uintSize = NumBits(ival)
		intSize  = NumBits(*oval)
	)

	if (uintSize >= intSize) && (ival > U(maxIntValue[intSize])) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%d", ival), fmt.Sprintf("%T", *oval))
	}

	*oval = I(ival)
	return nil
}

// MustUintToInt is a Must version of UintToInt
func MustUintToInt[U constraint.UnsignedInteger, I constraint.SignedInteger](ival U, oval *I) {
	funcs.Must(UintToInt(ival, oval))
}

// FloatToInt converts and float type to any signed int type
// Returns an error if the float value cannot be represented by the int type
func FloatToInt[F constraint.Float, I constraint.SignedInteger](ival F, oval *I) error {
	var (
		inter1 *big.Rat
		inter2 int64
	)

	if math.IsNaN(float64(ival)) || (FloatToBigRat(ival, &inter1) != nil) || (BigRatToInt64(inter1, &inter2) != nil) || (IntToInt(inter2, oval) != nil) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%g", ival), goreflect.TypeOf(*oval).Name())
	}

	return nil
}

// MustFloatToInt is a Must version of FloatToInt
func MustFloatToInt[F constraint.Float, I constraint.SignedInteger](ival F, oval *I) {
	funcs.Must(FloatToInt(ival, oval))
}

// ==== ToUint

// IntToUint converts any signed integer type into any unsigned integer type
// Returns an error if the signed int cannot be represented by the unsigned type
func IntToUint[I constraint.SignedInteger, U constraint.UnsignedInteger](ival I, oval *U) error {
	var (
		intSize  = NumBits(ival)
		uintSize = NumBits(*oval)
	)

	if (ival < 0) || ((intSize > uintSize) && (ival > I(maxUintValue[uintSize]))) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%d", ival), fmt.Sprintf("%T", *oval))
	}

	*oval = U(ival)
	return nil
}

// MustIntToUint is a Must version of IntToUint
func MustIntToUint[I constraint.SignedInteger, U constraint.UnsignedInteger](ival I, oval *U) {
	funcs.Must(IntToUint(ival, oval))
}

// UintToUint converts any unsigned integer type into any unsigned integer type
// Returns an error if the source value cannot be represented by the target type
func UintToUint[S constraint.UnsignedInteger, T constraint.UnsignedInteger](ival S, oval *T) error {
	var (
		srcSize = NumBits(ival)
		tgtSize = NumBits(*oval)
	)

	if (srcSize > tgtSize) && (ival > S(maxUintValue[tgtSize])) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%d", ival), fmt.Sprintf("%T", *oval))
	}

	*oval = T(ival)
	return nil
}

// MustUintToInt is a Must version of UintToInt
func MustUintToUint[S constraint.UnsignedInteger, T constraint.UnsignedInteger](ival S, oval *T) {
	funcs.Must(UintToUint(ival, oval))
}

// FloatToUint converts and float type to any unsigned int type
// Returns an error if the float value cannot be represented by the unsigned int type
func FloatToUint[F constraint.Float, I constraint.UnsignedInteger](ival F, oval *I) error {
	var (
		inter1 *big.Rat
		inter2 uint64
	)

	if math.IsNaN(float64(ival)) || (FloatToBigRat(ival, &inter1) != nil) || (BigRatToUint64(inter1, &inter2) != nil) || (UintToUint(inter2, oval) != nil) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%g", ival), goreflect.TypeOf(*oval).Name())
	}

	return nil
}

// MustFloatToUint is a Must version of FloatToUint
func MustFloatToUint[F constraint.Float, I constraint.UnsignedInteger](ival F, oval *I) {
	funcs.Must(FloatToUint(ival, oval))
}

// ==== ToInt64

// BigIntToInt64 converts a *big.Int to a signed integer
// Returns an error if the *big.Int cannot be represented as an int64
func BigIntToInt64(ival *big.Int, oval *int64) error {
	if !ival.IsInt64() {
		return fmt.Errorf(errMsg, ival, ival.String(), "int64")
	}

	*oval = ival.Int64()
	return nil
}

// MustBigIntToInt64 is a Must version of BigIntToInt64
func MustBigIntToInt64(ival *big.Int, oval *int64) {
	funcs.Must(BigIntToInt64(ival, oval))
}

// BigFloatToInt64 converts a *big.Float to an int64
// Returns an error if the *big.Float cannot be represented as an int64
func BigFloatToInt64(ival *big.Float, oval *int64) error {
	inter := big.NewInt(0)
	if (BigFloatToBigInt(ival, &inter) != nil) || (BigIntToInt64(inter, oval) != nil) {
		return fmt.Errorf(errMsg, ival, ival.String(), "int64")
	}

	return nil
}

// MustBigFloatToInt64 is a Must version of BigFloatToInt64
func MustBigFloatToInt64(ival *big.Float, oval *int64) {
	funcs.Must(BigFloatToInt64(ival, oval))
}

// BigRatToInt64 converts a *big.Rat to an int64
// Returns an error if the *big.Rat cannot be represented as an int64
func BigRatToInt64(ival *big.Rat, oval *int64) error {
	if (!ival.IsInt()) || (!ival.Num().IsInt64()) {
		return fmt.Errorf(errMsg, ival, ival.String(), "int64")
	}

	*oval = ival.Num().Int64()
	return nil
}

// MustBigRatToInt64 is a Must version of BigRatToInt64
func MustBigRatToInt64(ival *big.Rat, oval *int64) {
	funcs.Must(BigRatToInt64(ival, oval))
}

// StringToInt64 converts a string to an int64
// Returns an error if the string cannot be represented as an int64
func StringToInt64(ival string, oval *int64) error {
	var err error
	*oval, err = strconv.ParseInt(ival, 10, 64)
	if err != nil {
		return fmt.Errorf(errMsg, ival, ival, "int64")
	}

	return nil
}

// MustStringToInt64 is a Must version of StringToInt64
func MustStringToInt64(ival string, oval *int64) {
	funcs.Must(StringToInt64(ival, oval))
}

// ==== ToUint64

// BigIntToUint64 converts a *big.Int to a uint64
// Returns an error if the *big.Int cannot be represented as a uint64
func BigIntToUint64(ival *big.Int, oval *uint64) error {
	if !ival.IsUint64() {
		return fmt.Errorf(errMsg, ival, ival.String(), "uint64")
	}

	*oval = ival.Uint64()
	return nil
}

// MustBigIntToUint64 is a Must version of BigIntToUint64
func MustBigIntToUint64(ival *big.Int, oval *uint64) {
	funcs.Must(BigIntToUint64(ival, oval))
}

// BigFloatToUint64 converts a *big.Float to a uint64
// Returns an error if the *big.Float cannot be represented as a uint64
func BigFloatToUint64(ival *big.Float, oval *uint64) error {
	var inter *big.Int
	if (BigFloatToBigInt(ival, &inter) != nil) || (BigIntToUint64(inter, oval) != nil) {
		return fmt.Errorf(errMsg, ival, ival.String(), "uint64")
	}

	return nil
}

// MustBigFloatToUint64 is a Must version of BigFloatToUint64
func MustBigFloatToUint64(ival *big.Float, oval *uint64) {
	funcs.Must(BigFloatToUint64(ival, oval))
}

// BigRatToUint64 converts a *big.Rat to a uint64
// Returns an error if the *big.Rat cannot be represented as a uint64
func BigRatToUint64(ival *big.Rat, oval *uint64) error {
	if (!ival.IsInt()) || (!ival.Num().IsUint64()) {
		return fmt.Errorf(errMsg, ival, ival.String(), "uint64")
	}

	*oval = ival.Num().Uint64()
	return nil
}

// MustBigRatToUint64 is a Must version of BigRatToUint64
func MustBigRatToUint64(ival *big.Rat, oval *uint64) {
	funcs.Must(BigRatToUint64(ival, oval))
}

// StringToUint64 converts a string to a uint64
// Returns an error if the string cannot be represented as a uint64
func StringToUint64(ival string, oval *uint64) error {
	var err error
	if *oval, err = strconv.ParseUint(ival, 10, 64); err != nil {
		return fmt.Errorf(errMsg, ival, ival, "uint64")
	}

	return nil
}

// MustStringToUint64 is a Must version of StringToUint64
func MustStringToUint64(ival string, oval *uint64) {
	funcs.Must(StringToUint64(ival, oval))
}
