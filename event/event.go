package event

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/json"
)

var (
	// DefaultRegistry is a JSON registry that is easy to use
	DefaultRegistry = JSONRegistry{}
)

// ==== Registry

// Data is the data to send or receive
// - The Op field indicates the operation being performed, such as write to database
// - I is the type of input the event receivers get
// - R is the type of result the event receivers return
type Data[I, R any] struct {
	Op     int
	Input  I
	Result *R
}

// Receiver is a single method interface for event receivers that process the input for an event
type Receiver[I, R any] interface {
	Process(Data[I, R])
}

// ReceiverFunc is an adapter to allow the user of ordinary functions as Receivers
type ReceiverFunc[I, R any] func(Data[I, R])

// Process calls Process(d)
func (r ReceiverFunc[I, R]) Process(d Data[I, R]) {
	r(d)
}

// Registry is an event registry of event receivers. Senders only register if they are also a receiver.
// There can be any number of receivers for a particular operation, executed in the order they are registered.
// As a rule, receivers should not modify the input I, only the result R.
// The result R may be cumulative, where each receiver can add to it, or R can just be the result of the final operation.
// R may not even be neccessary. It is up to the user to decide on the semantics of R.
//
// The sender of the event can get the value of the result by examining the Result field of the Data.
// This requires the Result field to be a pointer.
//
// The zero value is ready to use.
//
// Note: all operations in a single registry have to receive the same input type I and return the same result type R
type Registry[I, R any] struct {
	receivers map[int][]Receiver[I, R]
}

// Register a receiver for a specific type of event
func (r *Registry[I, R]) Register(op int, rcvr Receiver[I, R]) {
	// ensure receivers is not nil
	if r.receivers == nil {
		r.receivers = map[int][]Receiver[I, R]{}
	}

	// Add receiver to end of existing list of receivers for the given operation
	r.receivers[op] = append(r.receivers[op], rcvr)
}

// Send an event to any receivers of the specified operation.
// Returns true if there exists at least one receiver of the event, false if there are no receivers.
func (r *Registry[I, R]) Send(data Data[I, R]) bool {
	var received bool

	if r.receivers != nil {
		for _, rcvr := range r.receivers[data.Op] {
			received = true
			rcvr.Process(data)
		}
	}

	return received
}

// ==== SimpleRegistry

// SimpleData is Data with the same type for input and output
type SimpleData[T any] struct {
	Data[T, T]
}

// SimpleReceiver is a Receiver with the same type for input and output
type SimpleReceiver[T any] interface {
	Receiver[T, T]
}

// SimpleRegistry is a registry with the same type for input and output
type SimpleRegistry[T any] struct {
	Registry[T, T]
}

// Register registers a SimpleReceiver for a SimpleRegistry
func (s *SimpleRegistry[T]) Register(op int, rcvr SimpleReceiver[T]) {
	s.Registry.Register(op, rcvr)
}

// Send a SimpleData to SimpleReceivers
func (s *SimpleRegistry[T]) Send(data SimpleData[T]) {
	s.Registry.Send(data.Data)
}

// ==== JSONRegistry

// JSONData is Data with json.Value type for input and output
type JSONData struct {
	Data[json.Value, json.Value]
}

// JSONReceiver is a Receiver with json.Value for input and output
type JSONReceiver interface {
	Receiver[json.Value, json.Value]
}

// JSONRegistry creates a simple registry using a jaon.Value object
type JSONRegistry struct {
	Registry[json.Value, json.Value]
}

// Register registers a JSONReceiver for a JSONRegistry
func (j *JSONRegistry) Register(op int, rcvr JSONReceiver) {
	j.Registry.Register(op, rcvr)
}

// Send a JSONData to JSONReceivers
func (j *JSONRegistry) Send(data JSONData) {
	j.Registry.Send(data.Data)
}

// ==== DefaultRegistry

// Register registers a JSONReceiver for the DefaultRegistry
func Register(op int, rcvr JSONReceiver) {
	DefaultRegistry.Registry.Register(op, rcvr)
}

// Send a JSONData to DefaultRegistry receivers
func Send(data JSONData) {
	DefaultRegistry.Registry.Send(data.Data)
}
