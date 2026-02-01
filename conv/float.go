package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"
	goreflect "reflect"

	"github.com/bantling/micro/constraint"
	"github.com/bantling/micro/funcs"
)

// ==== ToFloat

// IntToFloat converts any kind of signed or unssigned integer into any kind of float.
// Returns an error if the int value cannot be exactly represented without rounding.
func IntToFloat[I constraint.Integer, F constraint.Float](ival I, oval *F) error {
	// Convert int to float type, which may round if int has more bits than float type mantissa
	inter := F(ival)

	// If converting the float back to the int type is not the same value, rounding occurred
	if ival != I(inter) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%d", ival), goreflect.TypeOf(inter).Name())
	}

	*oval = inter
	return nil
}

// MustIntToFloat is a Must version of IntToFloat
func MustIntToFloat[I constraint.Integer, F constraint.Float](ival I, oval *F) {
	funcs.Must(IntToFloat(ival, oval))
}

// FloatToFloat converts a float32 or float64 to a float32 or float64
// Returns an error if the float64 is outside the range of a float32
func FloatToFloat[I constraint.Float, O constraint.Float](ival I, oval *O) error {
	ival64 := float64(ival)
	if math.IsInf(ival64, 0) || math.IsNaN(ival64) {
		*oval = O(ival)
		return nil
	}

	if _, isa := any(oval).(*float32); isa && (((ival64 != 0.0) && (math.Abs(ival64) < math.SmallestNonzeroFloat32)) || (math.Abs(ival64) > math.MaxFloat32)) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%g", ival64), "float32")
	}

	*oval = O(ival)
	return nil
}

// MustFloatToFloat is a Must version of FloatToFloat
func MustFloatToFloat[I constraint.Float, O constraint.Float](ival I, oval *O) {
	funcs.Must(FloatToFloat(ival, oval))
}

// ==== ToFloat32

// BigIntToFloat32 converts a *big.Int to a float32
// Returns an error if the *big.Int cannot be represented as a float32
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

// MustBigIntToFloat32 is a Must version of BigIntToFloat32
func MustBigIntToFloat32(ival *big.Int, oval *float32) {
	funcs.Must(BigIntToFloat32(ival, oval))
}

// BigFloatToFloat32 converts a *big.Float to a float32
// Returns an error if the *big.Float cannot be represented as a float32
func BigFloatToFloat32(ival *big.Float, oval *float32) error {
	var acc big.Accuracy
	if *oval, acc = ival.Float32(); acc != big.Exact {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%.f", ival), "float32")
	}

	return nil
}

// MustBigFloatToFloat32 is a Must version of BigFloatToFloat32
func MustBigFloatToFloat32(ival *big.Float, oval *float32) {
	funcs.Must(BigFloatToFloat32(ival, oval))
}

// BigRatToFloat32 converts a *big.Rat to a float32
// Returns an error if the *big.Rat cannot be represented as a float32
func BigRatToFloat32(ival *big.Rat, oval *float32) error {
	var exact bool
	if *oval, exact = ival.Float32(); !exact {
		return fmt.Errorf(errMsg, ival, ival.String(), "float32")
	}

	return nil
}

// MustBigRatToFloat32 is a Must version of BigRatToFloat32
func MustBigRatToFloat32(ival *big.Rat, oval *float32) {
	funcs.Must(BigRatToFloat32(ival, oval))
}

// StringToFloat32 converts a string to a float32
// Returns an error if the string cannot be represented as a float32
func StringToFloat32(ival string, oval *float32) error {
	var inter *big.Float

	if ival == "NaN" {
		*oval = float32(math.NaN())
		return nil
	}

	if (StringToBigFloat(ival, &inter) != nil) || (BigFloatToFloat32(inter, oval) != nil) {
		return fmt.Errorf(errMsg, ival, ival, "float32")
	}

	return nil
}

// MustStringToFloat32 is a Must version of StringToFloat32
func MustStringToFloat32(ival string, oval *float32) {
	funcs.Must(StringToFloat32(ival, oval))
}

// ==== ToFloat64

// BigIntToFloat64 converts a *big.Int to a float64
// Returns an error if the *big.Int cannot be represented as a float64
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

// MustBigIntToFloat64 is a Must version of BigIntToFloat64
func MustBigIntToFloat64(ival *big.Int, oval *float64) {
	funcs.Must(BigIntToFloat64(ival, oval))
}

// BigFloatToFloat64 converts a *big.Float to a float64
// Returns an error if the *big.Float cannot be represented as a float64
func BigFloatToFloat64(ival *big.Float, oval *float64) error {
	var acc big.Accuracy
	if *oval, acc = ival.Float64(); acc != big.Exact {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%.f", ival), "float64")
	}

	return nil
}

// MustBigFloatToFloat64 is a Must version of BigFloatToFloat64
func MustBigFloatToFloat64(ival *big.Float, oval *float64) {
	funcs.Must(BigFloatToFloat64(ival, oval))
}

// BigRatToFloat64 converts a *big.Rat to a float64
// Returns an error if the *big.Rat cannot be represented as a float64
func BigRatToFloat64(ival *big.Rat, oval *float64) error {
	var exact bool
	if *oval, exact = ival.Float64(); !exact {
		return fmt.Errorf(errMsg, ival, ival.String(), "float64")
	}

	return nil
}

// MustBigRatToFloat64 is a Must version of BigRatToFloat64
func MustBigRatToFloat64(ival *big.Rat, oval *float64) {
	funcs.Must(BigRatToFloat64(ival, oval))
}

// StringToFloat64 converts a string to a float64
// Returns an error if the string cannot be represented as a float64
func StringToFloat64(ival string, oval *float64) error {
	var inter *big.Float

	if ival == "NaN" {
		*oval = math.NaN()
		return nil
	}

	if (StringToBigFloat(ival, &inter) != nil) || (BigFloatToFloat64(inter, oval) != nil) {
		return fmt.Errorf(errMsg, ival, ival, "float64")
	}

	return nil
}

// MustStringToFloat64 is a Must version of StringToFloat64
func MustStringToFloat64(ival string, oval *float64) {
	funcs.Must(StringToFloat64(ival, oval))
}
