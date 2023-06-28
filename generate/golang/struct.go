package golang

// SPDX-License-Identifier: Apache-2.0

import (
  "strings"
  "text/template"

  "github.com/bantling/micro/funcs"
  "github.com/bantling/micro/generate"
)

const (
  structTemplate = "struct.tmpl"
)

// FieldType strings
var (
  fieldTypeString = map[generate.FieldType]string{
    generate.Int8_t: "int8",
    generate.Int16_t: "int16",
    generate.Int32_t: "int32",
    generate.Int64_t: "int64",
    generate.Uint8_t: "uint8",
    generate.Uint16_t: "uint16",
    generate.Uint32_t: "uint32",
    generate.Uint64_t: "uint64",
    generate.Float32_t: "float32",
    generate.Float64_t: "float64",
    generate.Date_t: "time.Time",
    generate.Datetime_t: "time.Time",
    generate.Time_t: "time.Time",
    generate.Interval_t: "time.Duration",
    generate.String_t: "string",
  }
)

// String is the Stringer interface for scalar FieldType values
func FieldTypeString(ft generate.FieldType) string {
  return fieldTypeString[ft]
}

// Field represents a single field of a data type
type Field struct {
  Comment string
  Name string
  Type generate.FieldType
  Ref bool // true for a reference (pointer), not applicable to arrays or maps
  Array bool // true if it is an array (slice) of Type
  ValueType generate.FieldType // for maps only, the value type (ignored if Type is not a map), limited to int8_t .. string_t
  Union []generate.FieldType // for unions only, the types in the union
}

// DataType represents a data type that needs to be transferred over the wire (eg HTTPS, SQL)
type DataType struct {
  Comment string
  Name string
  Fields []Field
}

// String is the Stringer implementation
func (dt DataType) String() string {
  var str strings.Builder

  tmpl, err := template.New(
    structTemplate,
  ).Funcs(
    template.FuncMap{"FieldTypeString": FieldTypeString},
  ).ParseFiles(structTemplate)
  funcs.Must(err)
  
  funcs.Must(tmpl.Execute(&str, dt))
  return str.String()
}
