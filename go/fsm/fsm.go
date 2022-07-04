package fsm

import (
	"fmt"
)

const (
	errCannotFindStateForInput = "Cannot find next state for input %v at row %v"
)

// SPDX-License-Identifier: Apache-2.0

// Machine is a Finite State Machine that given an input, it transitions to a new State S',
// where S' may be the same state S. A machine produces a result of type R.
type Machine[I comparable, R any, S any] interface {
	// Transition to a new state. A single call to transition may traverse multiple state transitions before arriving at
	// a stopping transition.
	Transition(I) S

	// The result of executing the machine on a given input.
	Result() R
}

// TableState describes a map entry in a TableMachine.
// I is the input type, R is the result type.
type TableState[I comparable, R any] interface {
	// NextRow returns the index of the next row to transition to, which is negative to indicate a stop state.
	NextRow() int

	// Result returns the new result from the input and old result. The new result may be the same as the old result.
	Result(I) R

	// Error generates an error for the current result R. If there is no error, the result is nil.
	Error(R) error
}

// TableMachine is a table driven Machine that contains an internal slice of map of inputs to states.
// Each slice row is a map that represents the current transition, where the initial transition is row 0.
//
// A call to Transition(I) looks up the given input in the current row map to determine a new row.
// A map does not have to contain all possible inputs, it can have a predefined catch-all entry that is used
// when the input does not exist as a map key.
//
// Transitions continue to occur until a TableState is reached where NextRow() returns a negative value.
type TableMachine[I comparable, R any, S TableState[I, R]] struct {
	// The table
	table []map[I]S

	// The current table row
	row int

	// The current result
	result R

	// True if a catch-all key is used
	haveCatchAll bool

	// Map key for catch-all
	catchAll I
}

// Construct a new TableMachine
func NewTableMachine[I comparable, R any, S TableState[I, R]](
	table []map[I]S,
	initialResult R,
	catchAll ...I,
) *TableMachine[I, R, S] {
	var (
		catchAllVal  I
		haveCatchAll = len(catchAll) > 0
	)
	if haveCatchAll {
		catchAllVal = catchAll[0]
	}

	return &TableMachine[I, R, S]{
		table:        table,
		row:          0,
		result:       initialResult,
		haveCatchAll: haveCatchAll,
		catchAll:     catchAllVal,
	}
}

// Transition to new states until a stopping state occurs as follows:
// - lookup a new state in the map of the current row with the input
// - if the state does not exist and there is a provided catch-all, lookup the catch-all state
// - if no state can be found, an error occurs
// - if the state next row is negative, then it is a stopping state
func (t *TableMachine[I, R, S]) Transition(input I) R {
	for {
		// Get next state using given input and current row
		state, haveIt := t.table[t.row][input]
		if (!haveIt) && t.haveCatchAll {
			// No entry for input, check catch-all if available
			if state, haveIt = t.table[t.row][t.catchAll]; !haveIt {
				// No catch-all either, die
				panic(fmt.Errorf(errCannotFindStateForInput, input, t.row))
			}
		}

		// Go to next row if this a non-final state
		if row := state.NextRow(); row >= 0 {
			t.row = row
		}

		// Calculate the result for a final state
		res := state.Result(input)

		// Die if this an error state
		if err := state.Error(res); err != nil {
			panic(err)
		}

		return res
	}
}

// Reset the TableMachine to start at row 0
func (t *TableMachine[I, R, S]) Reset() {
	t.row = 0
}
