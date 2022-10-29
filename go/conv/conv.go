package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"

	"github.com/bantling/micro/go/constraint"
	"github.com/bantling/micro/go/funcs"
)

var (
	errMsg      = "The %T value of %s cannot be converted to a %s"
	log2Of10    = math.Log2(10)
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

// ToString

// IntToStrinng converts any signed int type into a string
func IntToString[T constraint.SignedInteger](val T) string {
	return strconv.FormatInt(int64(val), 10)
}

// UintToStrinng converts any unsigned int type into a string
func UintToString[T constraint.UnsignedInteger](val T) string {
	return strconv.FormatUint(uint64(val), 10)
}

// FloatToStrinng converts any float type into a string
func FloatToString[T constraint.Float](val T) string {
	_, is32 := any(val).(float32)
	return strconv.FormatFloat(float64(val), 'f', -1, funcs.Ternary(is32, 32, 64))
}

// BigIntToString converts a *big.Int to a string
func BigIntToString(val *big.Int) string {
	return val.String()
}

// BigFloatToString converts a *big.Float to a string
func BigFloatToString(val *big.Float) string {
	return val.String()
}

// BigRatToString converts a *big.Rat to a string.
// The string will be a ratio like 5/4, if it is int it will be a ratio like 5/1.
func BigRatToString(val *big.Rat) string {
	return val.String()
}

// BigRatToNormalizedString converts a *big.Rat to a string.
// The string will be formatted like an integer if the ratio is an int, else formatted like a float if it is not an int.
func BigRatToNormalizedString(val *big.Rat) string {
	if val.IsInt() {
		return val.Num().String()
	}

	var inter *big.Float
	BigRatToBigFloat(val, &inter)
	return inter.String()
}

// Converts any signed or unsigned int type, any float type, *big.Int, *big.Float, or *big.Rat to a string.
// The *big.Rat string will be normalized (see BigRatToNormalizedString).
func ToString[T constraint.Integer | constraint.Float | *big.Int | *big.Float | *big.Rat](val T) string {
	if v, isa := any(val).(int); isa {
		return IntToString(v)
	} else if v, isa := any(val).(int8); isa {
		return IntToString(v)
	} else if v, isa := any(val).(int16); isa {
		return IntToString(v)
	} else if v, isa := any(val).(int32); isa {
		return IntToString(v)
	} else if v, isa := any(val).(int64); isa {
		return IntToString(v)
	} else if v, isa := any(val).(uint); isa {
		return UintToString(v)
	} else if v, isa := any(val).(uint8); isa {
		return UintToString(v)
	} else if v, isa := any(val).(uint16); isa {
		return UintToString(v)
	} else if v, isa := any(val).(uint32); isa {
		return UintToString(v)
	} else if v, isa := any(val).(uint64); isa {
		return UintToString(v)
	} else if v, isa := any(val).(float32); isa {
		return FloatToString(v)
	} else if v, isa := any(val).(float64); isa {
		return FloatToString(v)
	} else if v, isa := any(val).(*big.Int); isa {
		return BigIntToString(v)
	} else if v, isa := any(val).(*big.Float); isa {
		return BigFloatToString(v)
	}

	// Must be *big.Rat
	return BigRatToNormalizedString(any(val).(*big.Rat))
}

// ==== int/uint to int/uint, float to int, float64 to float32

// NumBits provides the number of bits of any integer or float type
func NumBits[T constraint.Signed | constraint.UnsignedInteger](val T) int {
	return int(reflect.ValueOf(val).Type().Size() * 8)
}

// IntToInt converts any signed integer type into any signed integer type
// Panics if the source value cannot be represented by the target type
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

// IntToUint converts any signed integer type into any unsigned integer type
// Panics if the signed int cannot be represented by the unsigned type
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

// UintToInt converts any unsigned integer type into any signed integer type
// Panics if the unsigned int cannot be represented by the signed type
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

// UintToUint converts any unsigned integer type into any unsigned integer type
// Panics if the source value cannot be represented by the target type
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

// IntToFloat converts any kind of signed or unssigned integer into any kind of float.
// Panics if the int value cannot be exactly represented without rounding.
func IntToFloat[I constraint.Integer, F constraint.Float](ival I, oval *F) error {
	// Convert int to float type, which may round if int has more bits than float type mantissa
	inter := F(ival)

	// If converting the float back to the int type is not the same value, rounding occurred
	if ival != I(inter) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%d", ival), reflect.TypeOf(inter).Name())
	}

	*oval = inter
	return nil
}

// FloatToInt converts and float type to any signed int type
// Panics if the float value cannot be represented by the int type
func FloatToInt[F constraint.Float, I constraint.SignedInteger](ival F, oval *I) error {
	var (
		inter1 *big.Rat
		inter2 int64
	)

	if (FloatToBigRat(ival, &inter1) != nil) || (BigRatToInt64(inter1, &inter2) != nil) || (IntToInt(inter2, oval) != nil) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%f", ival), reflect.TypeOf(*oval).Name())
	}

	return nil
}

// FloatToUint converts and float type to any unsigned int type
// Panics if the float value cannot be represented by the unsigned int type
func FloatToUint[F constraint.Float, I constraint.UnsignedInteger](ival F, oval *I) error {
	var (
		inter1 *big.Rat
		inter2 uint64
	)

	if (FloatToBigRat(ival, &inter1) != nil) || (BigRatToUint64(inter1, &inter2) != nil) || (UintToUint(inter2, oval) != nil) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%f", ival), reflect.TypeOf(*oval).Name())
	}

	return nil
}

// Float64ToFloat32 converts a float64 to a float32
// Panics if the float64 is outside the range of a float32
func Float64ToFloat32(ival float64, oval *float32) error {
	if (!math.IsInf(ival, 0)) && (ival != 0) && ((ival < math.SmallestNonzeroFloat32) || (ival > math.MaxFloat32)) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%f", ival), "float32")
	}

	*oval = float32(ival)
	return nil
}

// ==== ToInt64

// BigIntToInt converts a *big.Int to a signed integer
// Panics if the *big.Int cannot be represented as an int64
func BigIntToInt64(ival *big.Int, oval *int64) error {
	if !ival.IsInt64() {
		return fmt.Errorf(errMsg, ival, ival.String(), "int64")
	}

	*oval = ival.Int64()
	return nil
}

// BigFloatToInt64 converts a *big.Float to an int64
// Panics if the *big.Float cannot be represented as an int64
func BigFloatToInt64(ival *big.Float, oval *int64) error {
	inter := big.NewInt(0)
	if (BigFloatToBigInt(ival, &inter) != nil) || (BigIntToInt64(inter, oval) != nil) {
		return fmt.Errorf(errMsg, ival, ival.String(), "int64")
	}

	return nil
}

// BigRatToInt64 converts a *big.Rat to an int64
// Panics if the *big.Rat cannot be represented as an int64
func BigRatToInt64(ival *big.Rat, oval *int64) error {
	if (!ival.IsInt()) || (!ival.Num().IsInt64()) {
		return fmt.Errorf(errMsg, ival, ival.String(), "int64")
	}

	*oval = ival.Num().Int64()
	return nil
}

// StringToInt64 converts a string to an int64
// Panics if the string cannot be represented as an int64
func StringToInt64(ival string, oval *int64) error {
	var err error
	*oval, err = strconv.ParseInt(ival, 10, 64)
	if err != nil {
		return fmt.Errorf(errMsg, ival, ival, "int64")
	}

	return nil
}

// ==== ToUint64

// BigIntToUint64 converts a *big.Int to a uint64
// Panics if the *big.Int cannot be represented as a uint64
func BigIntToUint64(ival *big.Int, oval *uint64) error {
	if !ival.IsUint64() {
		return fmt.Errorf(errMsg, ival, ival.String(), "uint64")
	}

	*oval = ival.Uint64()
	return nil
}

// BigFloatToUint64 converts a *big.Float to a uint64
// Panics if the *big.Float cannot be represented as a uint64
func BigFloatToUint64(ival *big.Float, oval *uint64) error {
	var inter *big.Int
	if (BigFloatToBigInt(ival, &inter) != nil) || (BigIntToUint64(inter, oval) != nil) {
		return fmt.Errorf(errMsg, ival, ival.String(), "uint64")
	}

	return nil
}

// BigRatToUint64 converts a *big.Rat to a uint64
// Panics if the *big.Rat cannot be represented as a uint64
func BigRatToUint64(ival *big.Rat, oval *uint64) error {
	if (!ival.IsInt()) || (!ival.Num().IsUint64()) {
		return fmt.Errorf(errMsg, ival, ival.String(), "uint64")
	}

	*oval = ival.Num().Uint64()
	return nil
}

// StringToUint64 converts a string to a uint64
// Panics if the string cannot be represented as a uint64
func StringToUint64(ival string, oval *uint64) error {
	var err error
	if *oval, err = strconv.ParseUint(ival, 10, 64); err != nil {
		return fmt.Errorf(errMsg, ival, ival, "uint64")
	}

	return nil
}

// ==== ToFloat32

// BigIntToFloat32 converts a *big.Int to a float32
// Panics if the *big.Int cannot be represented as a float32
func BigIntToFloat32(ival *big.Int, oval *float32) error {
	var (
		inter *big.Float
		acc   big.Accuracy
	)
	BigIntToBigFloat(ival, &inter)
	if *oval, acc = inter.Float32(); acc != big.Exact {
		return fmt.Errorf(errMsg, ival, ival.String(), "float32")
	}

	return nil
}

// BigFloatToFloat32 converts a *big.Float to a float32
// Panics if the *big.Float cannot be represented as a float32
func BigFloatToFloat32(ival *big.Float, oval *float32) error {
	var acc big.Accuracy
	if *oval, acc = ival.Float32(); acc != big.Exact {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%.f", ival), "float32")
	}

	return nil
}

// BigRatToFloat32 converts a *big.Rat to a float32
// Panics if the *big.Rat cannot be represented as a float32
func BigRatToFloat32(ival *big.Rat, oval *float32) error {
	var exact bool
	if *oval, exact = ival.Float32(); !exact {
		return fmt.Errorf(errMsg, ival, ival.String(), "float32")
	}

	return nil
}

// StringToFloat32 converts a string to a float32
// Panics if the string cannot be represented as a float32
func StringToFloat32(ival string, oval *float32) error {
	var inter *big.Float
	if (StringToBigFloat(ival, &inter) != nil) || (BigFloatToFloat32(inter, oval) != nil) {
		return fmt.Errorf(errMsg, ival, ival, "float32")
	}

	return nil
}

// ==== ToFloat64

// BigIntToFloat64 converts a *big.Int to a float64
// Panics if the *big.Int cannot be represented as a float64
func BigIntToFloat64(ival *big.Int, oval *float64) error {
	var (
		inter *big.Float
		acc   big.Accuracy
	)

	BigIntToBigFloat(ival, &inter)
	if *oval, acc = inter.Float64(); acc != big.Exact {
		return fmt.Errorf(errMsg, ival, ival.String(), "float64")
	}

	return nil
}

// BigFloatToFloat64 converts a *big.Float to a float64
// Panics if the *big.Float cannot be represented as a float64
func BigFloatToFloat64(ival *big.Float, oval *float64) error {
	var acc big.Accuracy
	if *oval, acc = ival.Float64(); acc != big.Exact {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%.f", ival), "float64")
	}

	return nil
}

// BigRatToFloat64 converts a *big.Rat to a float64
// Panics if the *big.Rat cannot be represented as a float64
func BigRatToFloat64(ival *big.Rat, oval *float64) error {
	var exact bool
	if *oval, exact = ival.Float64(); !exact {
		return fmt.Errorf(errMsg, ival, ival.String(), "float64")
	}

	return nil
}

// StringToFloat64 converts a string to a float64
// Panics if the string cannot be represented as a float64
func StringToFloat64(ival string, oval *float64) error {
	var inter *big.Float
	if (StringToBigFloat(ival, &inter) != nil) || (BigFloatToFloat64(inter, oval) != nil) {
		return fmt.Errorf(errMsg, ival, ival, "float64")
	}

	return nil
}

// ==== ToBigInt

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
// Panics if the float has fractional digits
func FloatToBigInt[T constraint.Float](ival T, oval **big.Int) error {
	var inter big.Rat
	inter.SetFloat64(float64(ival))
	if !inter.IsInt() {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%f", ival), "*big.Int")
	}

	*oval = inter.Num()
	return nil
}

// BigFloatToBigInt converts a *big.Float to a *big.Int.
// Panics if the *big.Float has any fractional digits.
func BigFloatToBigInt(ival *big.Float, oval **big.Int) error {
	inter, acc := ival.Rat(nil)
	if (inter == nil) || (!inter.IsInt()) || (acc != big.Exact) {
		return fmt.Errorf(errMsg, ival, ival.String(), "*big.Int")
	}

	*oval = inter.Num()
	return nil
}

// BigRatToBigInt converts a *big.Rat to a *big.Int
// Panics if the *big.Rat is not an int
func BigRatToBigInt(ival *big.Rat, oval **big.Int) error {
	if !ival.IsInt() {
		return fmt.Errorf(errMsg, ival, ival.String(), "*big.Int")
	}

	*oval = ival.Num()
	return nil
}

// String toBigInt converts a string to a *big.Int.
// Panics if the string is not an integer.
func StringToBigInt(ival string, oval **big.Int) error {
	*oval = big.NewInt(0)
	if _, ok := (*oval).SetString(ival, 10); !ok {
		return fmt.Errorf(errMsg, ival, ival, "*big.Int")
	}

	return nil
}

// ==== ToBigFloat

// IntToBigFloat converts any signed int type into a *big.Float
func IntToBigFloat[T constraint.SignedInteger](ival T, oval **big.Float) {
	*oval = big.NewFloat(0)
	(*oval).SetInt64(int64(ival))
}

// UintToBigFloat converts any unsigned int type into a *big.Float
func UintToBigFloat[T constraint.UnsignedInteger](ival T, oval **big.Float) {
	*oval = big.NewFloat(0).SetUint64(uint64(ival))
}

// FloatToBigFloat converts any float type into a *big.Float
func FloatToBigFloat[T constraint.Float](ival T, oval **big.Float) {
	*oval = big.NewFloat(float64(ival))
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
// Panics if the string is not a valid float string
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

// ==== ToBigRat

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
	if math.IsInf(float64(ival), 0) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%f", ival), "*big.Rat")
	}

	var inter *big.Float
	FloatToBigFloat(ival, &inter)
	BigFloatToBigRat(inter, oval)
	return nil
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

// StringToBigRat converts a string into a *big.Rat
func StringToBigRat(ival string, oval **big.Rat) error {
	var ok bool
	if *oval, ok = big.NewRat(1, 1).SetString(ival); !ok {
		return fmt.Errorf(errMsg, ival, ival, "*big.Rat")
	}

	return nil
}

// FloatStringToBigRat converts a float string to a *big.Rat.
// Unlike StringToBigRat, it will not accept a ratio string like 5/4.
func FloatStringToBigRat(ival string, oval **big.Rat) error {
	// ensure the string is a float string, and not a ratio
	var err error

	if _, _, err = big.NewFloat(0).Parse(ival, 10); err != nil {
		return fmt.Errorf(errMsg, ival, ival, "*big.Rat")
	}

	if err = StringToBigRat(ival, oval); err != nil {
		return err
	}
	return nil
}
