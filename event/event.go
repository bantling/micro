package event

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/encoding/json"
	"github.com/bantling/micro/funcs"
)

var (
	// DefaultRegistry is a JSON registry that is easy to use
	DefaultRegistry = Registry[json.Value]{}
)

const (
	// The ALL constant is for the Remove method
	ALL = true
)

// None is an empty struct, for cases where receivers do not need any data to process
type None struct{}

// ==== Registry

// Receiver is a single method interface for event receivers that process the data for an event.

// Receiving and returning the same type allows D to be a value type, or to contain value types, no need for pointers.
// If the receiver can have an error, then D needs to contain an error, and the semantics of error handling are defined
// by the receivers.
//
// Im particular, if:
// - D contains an error
// - there are multiple receivers
// - a receiver before the last receiver has an error
// Then ideally the remaining receivers should do nothing, so that the original error is reported, to make debugging easier.
//
// Notes:
// - Since the Process method accepts and returns the same type, a receiver can be created by composing functions.
// - The None type is provided for cases where the receivers do not need any data.
type Receiver[D any] interface {
	Process(D) D
}

// ReceiverFunc is an adapter to allow the user of ordinary functions as Receivers
type ReceiverFunc[D any] func(D) D

// Process is the Receiver implementation
func (r ReceiverFunc[D]) Process(d D) D {
	return r(d)
}

// Registry is an event registry of event receivers. Senders only register if they are also a receiver.
// There can be any number of receivers for a particular operation, executed in the order they are registered.
// T is the type of event, a registry can handle any number of types. Ideally, T should be an enum type.
//
// The zero value is ready to use.
type Registry[D any] struct {
	receivers []Receiver[D]
}

// Register a receiver.
// Receivers can be registered at any time, typically they are registered during initialization.
// The same receiver can be registered multiple times.
// It can be particularly useful to add a debugging receiver after each normal receiver, to debug the data processing after
// each step.
func (r *Registry[D]) Register(rcvr Receiver[D]) {
	// ensure receivers is not nil
	if r.receivers == nil {
		r.receivers = []Receiver[D]{}
	}

	// Add receiver to end of existing list of receivers for the given operation
	r.receivers = append(r.receivers, rcvr)
}

// Remove a receiver for a specific type of event.
// Removes only the first occurrence, unless the optional all flag is true.
func (r *Registry[D]) Remove(rcvr Receiver[D], all ...bool) {
	funcs.SliceRemoveUncomparable(&(r.receivers), rcvr, all...)
}

// Send an event to any receivers of the specified operation.
// If there are no receivers, then the call is a no operation.
//
// D may be cumulative, where each receiver can add to it, or D can just be the result of the final operation.
// It is up to the user to decide on the semantics of D.
// If there are no receivers, the D value passed will be returned.
//
// See comments on Receiver interface regarding the None type and error handling.
func (r *Registry[D]) Send(d D) D {
	if r.receivers != nil {
		for _, rcvr := range r.receivers {
			d = rcvr.Process(d)
		}
	}

	return d
}

// ==== DefaultRegistry

// Register registers a Receiver for the DefaultRegistry
func Register(rcvr Receiver[json.Value]) {
	DefaultRegistry.Register(rcvr)
}

// Remove a receiver from the DefaultRegistry.
// Removes only the first occurrence, unless the optional all flag is true.
func Remove(rcvr Receiver[json.Value], all ...bool) {
	DefaultRegistry.Remove(rcvr, all...)
}

// Send a JSONData to DefaultRegistry receivers
func Send(data json.Value) json.Value {
	return DefaultRegistry.Send(data)
}
