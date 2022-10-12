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

	return BigRatToBigFloat(val).String()
}

// ==== int/uint to int/uint, float to int, float64 to float32

// NumBits provides the number of bits of any integer or float type
func NumBits[T constraint.Signed](val T) int {
	return int(reflect.ValueOf(val).Type().Size() * 8)
}

// IntToInt converts any signed integer type into any signed integer type
// Panics if the source value cannot be represented by the target type
func IntToInt[S constraint.SignedInteger, T constraint.SignedInteger](sval S, tval *T) {
	var (
		srcSize = NumBits(sval)
		tgtSize = NumBits(*tval)
	)

	if (srcSize > tgtSize) && ((sval < S(minIntValue[tgtSize])) || (sval > S(maxIntValue[tgtSize]))) {
		panic(fmt.Errorf(errMsg, sval, fmt.Sprintf("%d", sval), fmt.Sprintf("%T", *tval)))
	}

	*tval = T(sval)
}

// IntToUint converts any signed integer type into any unsigned integer type
// Panics if the signed int cannot be represented by the unsigned type
func IntToUint[I constraint.SignedInteger, U constraint.UnsignedInteger](ival I, uval *U) {
	var (
		intSize  = NumBits(ival)
		uintSize = NumBits(*uval)
	)

	if (ival < 0) || ((intSize > uintSize) && (ival > I(maxUintValue[uintSize]))) {
		panic(fmt.Errorf(errMsg, ival, fmt.Sprintf("%d", ival), fmt.Sprintf("%T", *uval)))
	}

	*uval = U(ival)
}

// UintToInt converts any unsigned integer type into any signed integer type
// Panics if the unsigned int cannot be represented by the signed type
func UintToInt[U constraint.UnsignedInteger, I constraint.SignedInteger](uval U, ival *I) {
	var (
		uintSize = NumBits(uval)
		intSize  = NumBits(*ival)
	)

	if (uintSize >= intSize) && (uval > U(maxIntValue[intSize])) {
		panic(fmt.Errorf(errMsg, uval, fmt.Sprintf("%d", uval), fmt.Sprintf("%T", *ival)))
	}

	*ival = I(uval)
}

// UintToUint converts any unsigned integer type into any unsigned integer type
// Panics if the source value cannot be represented by the target type
func UintToUint[S constraint.UnsignedInteger, T constraint.UnsignedInteger](sval S, tval *T) {
	var (
		srcSize = NumBits(sval)
		tgtSize = NumBits(*tval)
	)

	if (srcSize > tgtSize) && (sval > S(maxUintValue[tgtSize])) {
		panic(fmt.Errorf(errMsg, sval, fmt.Sprintf("%d", sval), fmt.Sprintf("%T", *tval)))
	}

	*tval = T(sval)
}

// FloatToInt converts and float type to any signed int type
// Panics if the float value cannot be represented by the int type
func FloatToInt[F constraint.Float, I constraint.SignedInteger](fval F, ival *I) {
	funcs.TryTo(
		func() { IntToInt(BigRatToInt64(FloatToBigRat(fval)), ival) },
		func(_ any) {
			panic(fmt.Errorf(errMsg, fval, fmt.Sprintf("%f", fval), reflect.TypeOf(*ival).Name()))
		},
	)
}

// FloatToUint converts and float type to any unsigned int type
// Panics if the float value cannot be represented by the unsigned int type
func FloatToUint[F constraint.Float, I constraint.UnsignedInteger](fval F, ival *I) {
	funcs.TryTo(
		func() { UintToUint(BigRatToUint64(FloatToBigRat(fval)), ival) },
		func(_ any) {
			panic(fmt.Errorf(errMsg, fval, fmt.Sprintf("%f", fval), reflect.TypeOf(*ival).Name()))
		},
	)
}

// Float64ToFloat32 converts a float64 to a float32
// Panics if the float64 is outside the range of a float32
func Float64ToFloat32(val float64) float32 {
	if (!math.IsInf(val, 0)) && ((val < math.SmallestNonzeroFloat32) || (val > math.MaxFloat32)) {
		panic(fmt.Errorf(errMsg, val, fmt.Sprintf("%f", val), "float32"))
	}

	return float32(val)
}

// ==== ToInt64

// BigIntToInt converts a *big.Int to a signed integer
// Panics if the *big.Int cannot be represented as an int64
func BigIntToInt64(val *big.Int) int64 {
	if !val.IsInt64() {
		panic(fmt.Errorf(errMsg, val, val.String(), "int64"))
	}

	return val.Int64()
}

// BigFloatToInt64 converts a *big.Float to an int64
// Panics if the *big.Float cannot be represented as an int64
func BigFloatToInt64(val *big.Float) int64 {
	var res int64
	funcs.TryTo(
		func() { res = BigIntToInt64(BigFloatToBigInt(val)) },
		func(_ any) {
			panic(fmt.Errorf(errMsg, val, val.String(), "int64"))
		},
	)

	return res
}

// BigRatToInt64 converts a *big.Rat to an int64
// Panics if the *big.Rat cannot be represented as an int64
func BigRatToInt64(val *big.Rat) int64 {
	if (!val.IsInt()) || (!val.Num().IsInt64()) {
		panic(fmt.Errorf(errMsg, val, val.String(), "int64"))
	}

	return val.Num().Int64()
}

// StringToInt64 converts a string to an int64
// Panics if the string cannot be represented as an int64
func StringToInt64(val string) int64 {
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(fmt.Errorf(errMsg, val, val, "int64"))
	}

	return i
}

// ==== ToUint64

// BigIntToUint64 converts a *big.Int to a uint64
// Panics if the *big.Int cannot be represented as a uint64
func BigIntToUint64(val *big.Int) uint64 {
	if !val.IsUint64() {
		panic(fmt.Errorf(errMsg, val, val.String(), "uint64"))
	}

	return val.Uint64()
}

// BigFloatToUint64 converts a *big.Float to a uint64
// Panics if the *big.Float cannot be represented as a uint64
func BigFloatToUint64(val *big.Float) uint64 {
	var res uint64
	funcs.TryTo(
		func() { res = BigIntToUint64(BigFloatToBigInt(val)) },
		func(_ any) {
			panic(fmt.Errorf(errMsg, val, val.String(), "uint64"))
		},
	)

	return res
}

// BigRatToUint64 converts a *big.Rat to a uint64
// Panics if the *big.Rat cannot be represented as a uint64
func BigRatToUint64(val *big.Rat) uint64 {
	if (!val.IsInt()) || (!val.Num().IsUint64()) {
		panic(fmt.Errorf(errMsg, val, val.String(), "uint64"))
	}

	return val.Num().Uint64()
}

// StringToUint64 converts a string to a uint64
// Panics if the string cannot be represented as a uint64
func StringToUint64(val string) uint64 {
	i, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(fmt.Errorf(errMsg, val, val, "uint64"))
	}

	return i
}

// ==== ToFloat32

// BigIntToFloat32 converts a *big.Int to a float32
// Panics if the *big.Int cannot be represented as a float32
func BigIntToFloat32(val *big.Int) float32 {
	f32, acc := BigIntToBigFloat(val).Float32()
	if acc != big.Exact {
		panic(fmt.Errorf(errMsg, val, val.String(), "float32"))
	}

	return f32
}

// BigFloatToFloat32 converts a *big.Float to a float32
// Panics if the *big.Float cannot be represented as a float32
func BigFloatToFloat32(val *big.Float) float32 {
	f32, acc := val.Float32()
	if acc != big.Exact {
		panic(fmt.Errorf(errMsg, val, fmt.Sprintf("%.f", val), "float32"))
	}

	return f32
}

// BigRatToFloat32 converts a *big.Rat to a float32
// Panics if the *big.Rat cannot be represented as a float32
func BigRatToFloat32(val *big.Rat) float32 {
	f32, exact := val.Float32()
	if !exact {
		panic(fmt.Errorf(errMsg, val, val.String(), "float32"))
	}

	return f32
}

// StringToFloat32 converts a string to a float32
// Panics if the string cannot be represented as a float32
func StringToFloat32(val string) float32 {
	var res float32
	funcs.TryTo(
		func() { res = BigFloatToFloat32(StringToBigFloat(val)) },
		func(e any) { panic(fmt.Errorf(errMsg, val, val, "float32")) },
	)

	return res
}

// ==== ToFloat64

// BigIntToFloat64 converts a *big.Int to a float64
// Panics if the *big.Int cannot be represented as a float64
func BigIntToFloat64(val *big.Int) float64 {
	f64, acc := BigIntToBigFloat(val).Float64()

	if acc != big.Exact {
		panic(fmt.Errorf(errMsg, val, val.String(), "float64"))
	}

	return f64
}

// BigFloatToFloat64 converts a *big.Float to a float64
// Panics if the *big.Float cannot be represented as a float64
func BigFloatToFloat64(val *big.Float) float64 {
	f64, acc := val.Float64()
	if acc != big.Exact {
		panic(fmt.Errorf(errMsg, val, fmt.Sprintf("%.f", val), "float64"))
	}

	return f64
}

// BigRatToFloat64 converts a *big.Rat to a float64
// Panics if the *big.Rat cannot be represented as a float64
func BigRatToFloat64(val *big.Rat) float64 {
	f64, exact := val.Float64()
	if !exact {
		panic(fmt.Errorf(errMsg, val, val.String(), "float64"))
	}

	return f64
}

// StringToFloat64 converts a string to a float64
// Panics if the string cannot be represented as a float64
func StringToFloat64(val string) float64 {
	var res float64
	funcs.TryTo(
		func() { res = BigFloatToFloat64(StringToBigFloat(val)) },
		func(e any) { panic(fmt.Errorf(errMsg, val, val, "float64")) },
	)

	return res
}

// ==== ToBigInt

// IntToBigInt converts any signed int type into a *big.Int
func IntToBigInt[T constraint.SignedInteger](val T) *big.Int {
	i := big.NewInt(0)
	return i.SetInt64(int64(val))
}

// UintToBigInt converts any unsigned int type into a *big.Int
func UintToBigInt[T constraint.UnsignedInteger](val T) *big.Int {
	i := big.NewInt(0)
	return i.SetUint64(uint64(val))
}

// FloatToBigInt converts any float type to a *big.Int
// Panics if the float has fractional digits
func FloatToBigInt[T constraint.Float](val T) *big.Int {
	var r big.Rat
	r.SetFloat64(float64(val))
	if !r.IsInt() {
		panic(fmt.Errorf(errMsg, val, fmt.Sprintf("%f", val), "*big.Int"))
	}

	return r.Num()
}

// BigFloatToBigInt converts a *big.Float to a *big.Int.
// Panics if the *big.Float has any fractional digits.
func BigFloatToBigInt(val *big.Float) *big.Int {
	r, acc := val.Rat(nil)
	if (!r.IsInt()) || (acc != big.Exact) {
		panic(fmt.Errorf(errMsg, val, val.String(), "*big.Int"))
	}

	return r.Num()
}

// BigRatToBigInt converts a *big.Rat to a *big.Int
// Panics if the *big.Rat is not an int
func BigRatToBigInt(val *big.Rat) *big.Int {
	if !val.IsInt() {
		panic(fmt.Errorf(errMsg, val, val.String(), "*big.Int"))
	}

	return val.Num()
}

// String toBigInt converts a string to a *big.Int.
// Panics if the string is not an integer.
func StringToBigInt(val string) *big.Int {
	jv := big.NewInt(0)
	if _, ok := jv.SetString(val, 10); !ok {
		panic(fmt.Errorf(errMsg, val, val, "*big.Int"))
	}

	return jv
}

// ==== ToBigFloat

// IntToBigFloat converts any signed int type into a *big.Float
func IntToBigFloat[T constraint.SignedInteger](val T) *big.Float {
	f := big.NewFloat(0)
	return f.SetInt64(int64(val))
}

// UintToBigFloat converts any unsigned int type into a *big.Float
func UintToBigFloat[T constraint.UnsignedInteger](val T) *big.Float {
	f := big.NewFloat(0)
	return f.SetUint64(uint64(val))
}

// FloatToBigFloat converts any float type into a *big.Float
func FloatToBigFloat[T constraint.Float](val T) *big.Float {
	return big.NewFloat(float64(val))
}

// BigIntToBigFloat converts a *big.Int into a *big.Float
func BigIntToBigFloat(val *big.Int) *big.Float {
	return StringToBigFloat(val.String())
}

// BigRatToBigFloat converts a *big.Rat to a *big.Float
func BigRatToBigFloat(val *big.Rat) *big.Float {
	// Use numerator to calculate the precision, shd be accurate since denominator is basically the exponent
	prec := int(math.Ceil(math.Max(float64(53), float64(len(val.Num().String()))*log2Of10)))
	res, _, _ := big.ParseFloat(val.FloatString(prec), 10, uint(prec), big.ToNearestEven)

	// Set accuracy to exact
	res.SetMode(res.Mode())

	return res
}

// StringToBigFloat converts a string to a *big.Float
// Panics if the string is not a valid float string
func StringToBigFloat(val string) *big.Float {
	// A *big.Float is imprecise, but you can set the precision
	// The crude measure we use is the largest of 53 (number of bits in a float64) and ceiling(string length * Log2(10))
	// If every char was a significant digit, the ceiling calculation would be the minimum number of bits required
	numBits := uint(math.Max(53, math.Ceil(float64(len(val))*log2Of10)))
	f, _, err := big.ParseFloat(val, 10, numBits, big.ToNearestEven)
	if err != nil {
		panic(fmt.Errorf(errMsg, val, val, "*big.Float"))
	}

	return f
}

// ==== ToBigRat

// IntToBigRat converts any signed int type into a *big.Rat
func IntToBigRat[T constraint.SignedInteger](val T) *big.Rat {
	r := big.NewRat(1, 1)
	return r.SetInt64(int64(val))
}

// UintToBigRat converts any unsigned int type into a *big.Rat
func UintToBigRat[T constraint.UnsignedInteger](val T) *big.Rat {
	r := big.NewRat(1, 1)
	return r.SetUint64(uint64(val))
}

// FloatToBigRat converts any float type into a *big.Rat
func FloatToBigRat[T constraint.Float](val T) *big.Rat {
	r := big.NewRat(1, 1)
	if r.SetFloat64(float64(val)) == nil {
		panic(fmt.Errorf(errMsg, val, val, "*big.Rat"))
	}

	return r
}

// BigIntToBigRat converts a *big.Int into a *big.Rat
func BigIntToBigRat(val *big.Int) *big.Rat {
	r := big.NewRat(1, 1)
	r.SetFrac(val, big.NewInt(1))

	return r
}

// BigFloatToBigRat converts a *big.Float into a *big.Rat
func BigFloatToBigRat(val *big.Float) *big.Rat {
	r := big.NewRat(1, 1)
	r.SetString(BigFloatToString(val))

	return r
}

// StringToBigRat converts a string into a *big.Rat
func StringToBigRat(val string) *big.Rat {
	r := big.NewRat(1, 1)
	if _, ok := r.SetString(val); !ok {
		panic(fmt.Errorf(errMsg, val, val, "*big.Rat"))
	}

	return r
}