package code

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
	BitShiftLeft
	BitShiftRight
  BitShiftRightArithmetic

	// internal constant for one past last binary
	afterBinary
)

// BooleanOperator is a boolean operator with one to three arguments
type BooleanOperator Operator

const (
	// Logical
	And BooleanOperator = BooleanOperator(afterBinary)
	Or
  Not
  Ternary

	// Relational
	Lesser
	LesserEquals
	Equals
	GreaterEquals
	Greater
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
  op  UnaryOperator,
  val *Val,
) Expr {
  return Expr{
    Op: Operator(op),
    Val1: funcs.MustNonNilValue(val),
  }
}

// OfBinaryExpr constructs a binary Expr
func OfBinaryExpr(
  op   BinaryOperator,
  val1 *Val,
  val2 *Val,
) Expr {
  return Expr{
    Op: Operator(op),
    Val1: funcs.MustNonNilValue(val1),
    Val2: union.Present(val2),
  }
}

// OfBooleanExpr constructs a binary boolean Expr
func OfBooleanExpr(
  op   BooleanOperator,
  val1 *Val,
  val2 *Val,
  val3 ...*Val,
) Expr {
  var ternVal union.Maybe[*Val]
  if op == Ternary {
    ternVal = union.Present(val3[0])
  }
  
  return Expr{
    Op:   Operator(op),
    Val1: funcs.MustNonNilValue(val1),
    Val2: union.Present(val2),
    Val3: ternVal,
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
