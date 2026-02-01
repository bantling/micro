package constraint

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"
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
// Embeds comparable so that the Cmp interface can be a map key
type Cmp[T any] interface {
	comparable
	// Returns <0 if this value < argument
	//          0 if this value = argument
	//         >0 if this value > argument
	Cmp(T) int
}

// BigOps is an interface that describes all the common methods of the provided go big types
type BigOps[T any] interface {
	Abs(T) T
	Add(T, T) T
	Cmp(T) int
	Mul(T, T) T
	Neg(T) T
	Quo(T, T) T
	Set(T) T
	SetInt64(int64) T
	SetUint64(uint64) T
	Sign() int
	String() string
	Sub(T, T) T
}
