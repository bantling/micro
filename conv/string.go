package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"
	"strconv"

	"github.com/bantling/micro/constraint"
	"github.com/bantling/micro/funcs"
)

// IntToString converts any signed int type into a string
func IntToString[T constraint.SignedInteger](val T) string {
	return strconv.FormatInt(int64(val), 10)
}

// UintToString converts any unsigned int type into a string
func UintToString[T constraint.UnsignedInteger](val T) string {
	return strconv.FormatUint(uint64(val), 10)
}

// FloatToString converts any float type into a string
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
