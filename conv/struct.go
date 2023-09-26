package conv

// SPDX-License-Identifier: Apache-2.0

// MapToStruct populates a struct from a map[string]any.
// The struct may contain sub structs as values or pointers.
// The struct may be recursive (eg Customer{child *Customer}).
// The conv.LookupConversion function is used to locate a suitable conversion, if one exists.
//
// This func is not generic on struct type, to support cases where the struct type is not known ahead of time, such as:
// - The map is recursive, and the struct has child structs that cannot be known ahead of time.
// - A generalized algorithm that can work with any struct (eg, JSON -> struct, struct -> database row, etc)
func MapToStruct(val map[string]any, str any) {

}
