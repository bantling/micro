package generate

// SPDX-License-Identifier: Apache-2.0

// FieldType enum represents the type of a field
type FieldType uint

// FieldType constants
const (
  Int8_t FieldType = iota
  Int16_t
  Int32_t
  Int64_t
  Uint8_t
  Uint16_t
  Uint32_t
  Uint64_t
  Float32_t
  Float64_t
  Date_t
  Datetime_t
  Time_t
  Interval_t
  String_t
)

// Field represents a single field of a data type
type Field struct {
  Comment string
  Name string
  Type FieldType // Field type
  Array bool // true if it is an array (slice) of Type
  Ref bool // true for a reference (pointer), when used with arrays it is applied to Type, not the array (eg []*Type)
}

// DataType represents a data type that needs to be transferred over the wire (eg HTTPS, SQL)
type DataType struct {
  Comment string
  Name string
  Union bool // True if this type is a union, which means only one field actually available to be used
  Fields []Field
}
