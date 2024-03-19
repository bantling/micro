package golang

//
// // SPDX-License-Identifier: Apache-2.0
//
// import (
// 	"github.com/bantling/micro/code"
// 	"github.com/bantling/micro/funcs"
// )
//
// // Map general primitive types to Go types
// var (
// 	typeMap = map[string]string{
// 		// Boolean
// 		code.Bool: "bool",
//
// 		// Unsigned ints
// 		code.Uint8:  "uint8",
// 		code.Uint16: "uint16",
// 		code.Uint32: "uint32",
// 		code.Uint64: "uint64",
//
// 		// Signed ints
// 		code.Int8:  "int8",
// 		code.Int16: "int16",
// 		code.Int32: "int32",
// 		code.Int64: "int64",
//
// 		// String, UUID, JSON
// 		code.String: "string",
// 		code.Uuid:   "uuid.UUID",  // provided by Google library github.com/google/uuid
// 		code.Json:   "json.Value", // provided by this library in encoding/json
//
// 		// Date, DateTime, and Interval
// 		code.Date:                 "time.Time",
// 		code.DateTimeSecs:      "time.Time",
// 		code.DateTimeMillis: "time.Time",
// 		code.DurationDays:         "time.Duration",
// 		code.DurationSecs:      "time.Duration",
// 		code.DurationMillis: "time.Duration",
// 	}
// )
//
// // TypeDefString produces a type declaration
// // The optional position is only used in recursive calls, the caller must not provide a value for it
// func TypeDefString(def code.TypeDef, pos ...int) string {
//   var (
//     typStr string
//     posVal = funcs.SliceIndex(pos, 0)
//   )
//
//   switch def.Type {
//   case code.Array:
//     // Arrays can be map values, but not keys: position does not affect declaration
//     for _, dim := range def.Bounds {
//       typStr = "[" + strconv.Itoa(dim) + "]"
//     }
//     typStr += TypeDefString(def.ValueType, 1)
//
//   case code.Enum:
//     // Enums can only be top level
//     typStr = "type " + def.Name + " int\n\nconst (\n  " + def.Names[0] + " " + def.Name + " = iota\n"
//     for _, konst := range def.Names[1:] {
//       typStr += "  " + konst + "\n"
//     }
//     typStr += ")\n\n"
//
//   case code.List:
//     // Lists can be map values, but not keys: position does not affect declaration
//     typStr = strings.Repeat("[]", def.ListDimension) + TypeDefString(def.ValueType, 1)
//
//   case code.Map:
//     // Maps can be map values, but not keys: position does not affect declaration
//     typStr = "map[" + TypeDefString(def.KeyType, 1) + "]" + TypeDefString(def.ValueType, 1)
//
//   case code.Maybe:
//     // Maybes can be map values, but not keys: position does not affect declaration
//     typStr = "Maybe[" + TypeDefString(def.ValueType, 1) + "]"
//
//   case code.Object:
//     // Objects can be map values, but not keys: position DOES affect declaration
//     typStr = funcs.Ternary(pos == 0, "struct ", "") + def.Name
//     if len(def.Names) > 0 {
//       // Add generic bounds
//       typStr += "[" + def.Names[0]
//       for _, bound := range def.Names[1:] {
//         typStr += ", " + bound
//       }
//       typStr += "]"
//     }
//   case code.Set:
//     return "map[" + TypeDefString(def.ValueType, 1) + "]bool"
//   }
//
//   return typStr
// }
