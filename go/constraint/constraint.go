// Package constraint defines useful constraints that are not in the go builtin package
package constraint

// Complex is copied from golang.org/x/exp/constraints#Complex
type Complex interface {
	~complex64 | ~complex128
}

// Float is copied from golang.org/x/exp/constraints#Float
type Float interface {
	~float32 | ~float64
}

// Integer is equivalent to from golang.org/x/exp/constraints#Integer
type Integer interface {
	SignedInteger | UnsignedInteger
}

// Ordered is copied from golang.org/x/exp/constraints#Ordered
type Ordered interface {
	Integer | Float | ~string
}

// Signed differs from golang.org/x/exp/constraints#Signed - it includes Float
type Signed interface {
	Integer | Float
}

// SignedInteger is copied from golang.org/x/exp/constraints#Signed
type SignedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// UnsignedInteger is copied from golang.org/x/exp/constraints#Unsigned
type UnsignedInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
