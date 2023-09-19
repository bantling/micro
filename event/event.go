package event

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/constraint"
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/tuple"
)

const (
	// The ALL constant is for the Remove method
	ALL = true
)

// None is an empty struct, for cases where receivers do not need any data to process
type None struct{}

// ==== Registry

// Receiver is a single method interface for event receivers that process the data for an event.

// Receiving and returning the same type allows T to be a value type, or to contain value types, no need for pointers.
// If the receiver can have an error, then T needs to contain an error, and the semantics of error handling are defined
// by the receivers.
//
// Im particular, if:
// - T contains an error
// - there are multiple receivers
// - a receiver before the last receiver has an error
// Then ideally the remaining receivers should do nothing, so that the original error is reported, to make debugging easier.
//
// Notes:
// - Since the Process method accepts and returns the same type, a receiver can be created by composing functions.
// - The None type is provided for cases where the receivers do not need any data.
type Receiver[T any] interface {
	Process(T) T
}

// ReceiverFunc is an adapter to allow the user of ordinary functions as Receivers
type ReceiverFunc[T any] func(T) T

// Process is the Receiver implementation
func (r ReceiverFunc[T]) Process(t T) T {
	return r(t)
}

// Registry is an event registry of event receivers. A receiver may also be a sender.
// There can be any number of receivers for a particular operation, executed in the order they are registered.
// I is the type of id used for operations, it is only used for ordering the calls to receivers.
// Each operation must have a unique id.
//
// The senders and receivers do not provide or receive the id, only the event type T.
// T is the type of event.
//
// The zero value is ready to use.
type Registry[I constraint.Ordered, T any] struct {
	receivers        map[I][]Receiver[T]
	sortedIds        []tuple.Two[I, []Receiver[T]]
	receiversChanged bool
}

// sortIds sorts the ids of the receivers map, but only if ids have been added and/or removed since last call.
// It doesn't matter if receivers were added to an existing id, that has no effect on ordering.
// It is up to other methods to set the idsChanged flag.
func (r *Registry[I, T]) sortIds() []tuple.Two[I, []Receiver[T]] {
	if r.receiversChanged {
		r.receiversChanged = false
		r.sortedIds = funcs.MapSortOrdered(r.receivers)
	}

	return r.sortedIds
}

// Register a receiver.
// Receivers can be registered at any time, typically they are registered during initialization.
// The same receiver can be registered multiple times.
// It can be particularly useful to add a debugging receiver after each normal receiver, to debug the data processing after
// each step.
func (r *Registry[I, T]) Register(id I, rcvr Receiver[T]) {
	// Track id changes so we know when to re-sort
	r.receiversChanged = r.receiversChanged || !funcs.MapTest(r.receivers, id)

	// Add receiver to end of existing list of receivers for the given operation
	funcs.MapSliceAdd(&r.receivers, id, rcvr)
}

// Remove a given id and all receivers associated with it.
func (r *Registry[I, T]) RemoveId(id I) {
	// Remove id from map
	funcs.MapUnset(r.receivers, id)

	// Track receiver changes so we know when to re-sort
	r.receiversChanged = true
}

// Remove a receiver from a specific id.
// Removes only the first occurrence, unless the optional all flag is true.
func (r *Registry[I, T]) Remove(id I, rcvr Receiver[T], all ...bool) {
	funcs.MapSliceRemoveUncomparable(r.receivers, id, rcvr, all...)
	r.receiversChanged = true
}

// Send an event to any receivers of the specified operation.
// If there are no receivers, then the call is a no operation.
//
// T may be cumulative, where each receiver can add to it, or T can just be the result of the final operation.
// It is up to the user to decide on the semantics of T.
// If there are no receivers, the T value passed will be returned.
//
// See comments on Receiver interface regarding the None type and error handling.
func (r *Registry[I, T]) Send(val T) T {
	// []tuple.Two[I, []Receiver[T]]
	if sIds := r.sortIds(); sIds != nil {
		for _, idRcvrs := range sIds {
			for _, rcvr := range idRcvrs.U {
				val = rcvr.Process(val)
			}
		}
	}

	return val
}
