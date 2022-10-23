package util

// SPDX-License-Identifier: Apache-2.0

// KeyValue is a struct to hold a single key/pair for a map[K]V entry
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

// KVOf constructs a KeyValue so user does not have to enter generic type, Go can infer them
func KVOf[K comparable, V any](key K, val V) KeyValue[K, V] {
	return KeyValue[K, V]{key, val}
}
