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
)

// BinaryOperator is an operator with two arguments
type BinaryOperator Operator

const (
	// Four basic ops
	Add BinaryOperator = iota
	Sub
	Mul
	Div

	// String
	Concat

	// Logical
	And
	Or

	// Relational
	LessThan
	LessThanEquals
	Equals
	GreaterEquals
	Greater

	// Bitwise
	BitAnd
	BitOr
	BitXor
)

// Expr is an expression
// If the Operator is a UnaryOperator, then Val2 is empty
type Expr struct {
	Op   Operator
	Val1 Val
	Val2 union.Maybe[Val]
}

// StmtKind describes the type of statement
type StmtKind uint

const (
	Local      StmtKind = iota // Local is declaration of a local var
	Assignment                 // Assign a value to a local var
)

// Stmt is a statement
type Stmt struct {
	Kind  StmtKind             // The kind of statement
	Type  union.Maybe[TypeDef] // The TypeDef for a Local
	Value union.Maybe[Expr]    // The Value to assign
}
