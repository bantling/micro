package code

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/union"
)

// SPDX-License-Identifier: Apache-2.0

// Operator describes all types of operators
type Operator uint

// UnaryOperator is a non-boolean operator with one argument
type UnaryOperator Operator

const (
	// Increment/decrement
	Inc UnaryOperator = iota
	Dec

	// Negation
	Neg

	// Bitwise not
	BitNot

	// internal constant for one past last unary
	afterUnary
)

// BinaryOperator is a non-boolean operator with two arguments
type BinaryOperator Operator

const (
	// Four basic ops
	Add BinaryOperator = iota + BinaryOperator(afterUnary)
	Sub
	Mul
	Div

	// String
	Concat

	// Bitwise
	BitAnd
	BitOr
	BitXor
	BitShiftLeft
	BitShiftRight
	BitShiftRightArithmetic

	// internal constant for one past last binary
	afterBinary
)

// BooleanOperator is a boolean operator with one to three arguments
type BooleanOperator Operator

const (
	// Logical Unary
	Not BooleanOperator = iota + BooleanOperator(afterBinary)
	// Logical Binary
	And
	Or

	// Relational Binary
	Lesser
	LesserEquals
	Equals
	GreaterEquals
	Greater

	// Logical Ternary
	Ternary
)

// IsUnary is true if the Operator is a UnaryOperator
func (op Operator) IsUnary() bool {
	return uint(op) < uint(afterUnary)
}

// IsBinary is true if the Operator is a BinaryOperator
func (op Operator) IsBinary() bool {
	return (uint(op) >= uint(afterUnary)) && (uint(op) < uint(afterBinary))
}

// IsBoolean is true if the Operator is a BooleanOperator
func (op Operator) IsBoolean() bool {
	return uint(op) >= uint(afterBinary)
}

// Expr is an expression
// If the Operator is a UnaryOperator, then Val2 is empty
// If the Operator is not Ternary, then Val3 is empty
type Expr struct {
	Op   Operator
	Val1 *Val
	Val2 union.Maybe[*Val]
	Val3 union.Maybe[*Val]
}

// OfUnaryExpr constructs a unary Expr
func OfUnaryExpr(
	op UnaryOperator,
	val *Val,
) Expr {
	return Expr{
		Op:   Operator(op),
		Val1: funcs.MustNonNilValue(val),
	}
}

// OfBinaryExpr constructs a binary Expr
func OfBinaryExpr(
	op BinaryOperator,
	val1 *Val,
	val2 *Val,
) Expr {
	return Expr{
		Op:   Operator(op),
		Val1: funcs.MustNonNilValue(val1),
		Val2: union.Present(val2),
	}
}

// OfBooleanExpr constructs a binary boolean Expr
// If the operator is Not, then val23 is ignored
// If the operator is Ternary, then val23 must have two values
// All other operators are binary, so val23 must have one value
func OfBooleanExpr(
	op BooleanOperator,
	val1 *Val,
	val23 ...*Val,
) Expr {
	var val2, val3 union.Maybe[*Val]
	switch op {
	case Not:
		// Only need one value
	case Ternary:
		// Need three values
		val2 = union.Present(val23[0])
		val3 = union.Present(val23[1])
	default:
		// Rest are binary operators, need two values
		val2 = union.Present(val23[0])
	}

	return Expr{
		Op:   Operator(op),
		Val1: funcs.MustNonNilValue(val1),
		Val2: val2,
		Val3: val3,
	}
}

// StmtKind describes the type of statement
type StmtKind uint

const (
	Constant   StmtKind = iota // Constant is a local constant
	Local                      // Local is a local var
	Assignment                 // Assign a value to a local var
	Case                       // Conditional, flexible like SQL or go switch
)

// StmtDef is a statement
type StmtDef struct {
	Kind StmtKind              // The kind of statement
	Type union.Maybe[*TypeDef] // The TypeDef for a Constant or Local
	Expr union.Maybe[*Expr]    // The Value to assign a Constant or Local
}
