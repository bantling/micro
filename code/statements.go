package code

import (
	"github.com/bantling/micro/union"
)

// SPDX-License-Identifier: Apache-2.0

// Operator describes all types of operators
type Operator uint

// UnaryOperator is an operator with one argument
type UnaryOperator Operator

const (
	// Increment/decrement
	Inc UnaryOperator = iota
	Dec

	// Negation
	Neg

	// Logical not
	Not

	// Bitwise not
	BitNot

	// internal constant for one past last unary
	afterUnary
)

// BinaryOperator is an operator with two arguments
type BinaryOperator Operator

const (
	// Four basic ops
	Add BinaryOperator = BinaryOperator(afterUnary)
	Sub
	Mul
	Div

	// String
	Concat

	// Bitwise
	BitAnd
	BitOr
	BitXor

	// internal constant for one past last binary
	afterBinary
)

type BooleanOperator Operator

const (
	// Logical
	And BooleanOperator = BooleanOperator(afterBinary)
	Or

	// Relational
	LessThan
	LessThanEquals
	Equals
	GreaterEquals
	Greater
)

// IsUnary is true if the Operator is a UnaryOperator
func IsUnary(op Operator) bool {
	return uint(op) < uint(afterUnary)
}

// IsBinary is true if the Operator is a BinaryOperator
func IsBinary(op Operator) bool {
	return (uint(op) >= uint(afterUnary)) && (uint(op) < uint(afterBinary))
}

// IsBoolean is true if the Operator is a BooleanOperator
func IsBoolean(op Operator) bool {
	return uint(op) >= uint(afterBinary)
}

// ExprDef is an expression
// If the Operator is a UnaryOperator, then Val2 is empty
type ExprDef struct {
	Op   Operator
	Val1 Val
	Val2 union.Maybe[Val]
}

// StmtKind describes the type of statement
type StmtKind uint

const (
	Local      StmtKind = iota // Local is declaration of a local var
	Assignment                 // Assign a value to a local var
	Case                       // Conditional
)

// StmtDef is a statement
type StmtDef struct {
	Kind StmtKind             // The kind of statement
	Type union.Maybe[TypeDef] // The TypeDef for a Local
	Expr union.Maybe[ExprDef] // The Value to assign a local
}
