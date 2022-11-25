// Package constraint defines useful constraints that are not in the go builtin package
package constraint

import (
	"math/big"
	"reflect"
	"strings"
)

// SPDX-License-Identifier: Apache-2.0

// SignedInteger is copied from golang.org/x/exp/constraints#Signed
type SignedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// UnsignedInteger is like golang.org/x/exp/constraints#Unsigned, except no uintptr
type UnsignedInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Integer is equivalent to golang.org/x/exp/constraints#Integer
type Integer interface {
	SignedInteger | UnsignedInteger
}

// Float is copied from golang.org/x/exp/constraints#Float
type Float interface {
	~float32 | ~float64
}

// Signed differs from golang.org/x/exp/constraints#Signed - it includes Float
type Signed interface {
	SignedInteger | Float
}

// IntegerAndFloat describes any signed or unsigned integer type and any float type
type IntegerAndFloat interface {
	Integer | Float
}

// Numeric describes any numeric type except complex
type Numeric interface {
	IntegerAndFloat | *big.Int | *big.Float | *big.Rat
}

// Ordered is equivalent to golang.org/x/exp/constraints#Ordered
type Ordered interface {
	Signed | UnsignedInteger | ~string
}

// Complex is copied from golang.org/x/exp/constraints#Complex
type Complex interface {
	~complex64 | ~complex128
}

// Cmp is a companion interface for Ordered
type Cmp[T any] interface {
	// Returns <0 if this value < argument
	//          0 if this value = argument
	//         >0 if this value > argument
	Cmp(T) int
}

// IsSignedInt returns true if the value given is an int, int8, int16, int32, or int64
func IsSignedInt(t any) bool {
	return strings.HasPrefix(reflect.TypeOf(t).Name(), "int")
}

// IsUnsignedInt returns true if the value given is a uint, uint8, uint16, uint32, or uint64
func IsUnsignedInt(t any) bool {
	return strings.HasPrefix(reflect.TypeOf(t).Name(), "uint")
}

// IsFloat returns true if the value given is a float64 or float32
func IsFloat(t any) bool {
	return strings.HasPrefix(reflect.TypeOf(t).Name(), "float")
}

// IsBig returns true if the value given is a *big.Int, *big.Float, or *big.Rat
func IsBig(t any) bool {
	if _, isa := t.(*big.Int); isa {
		return true
	}

	if _, isa := t.(*big.Float); isa {
		return true
	}

	if _, isa := t.(*big.Rat); isa {
		return true
	}

	return false
}
