// Package conv is conversions between various types, without loss of precision.
// Functions panic if a conversion cannot be done precisely.
//
// Provides a registration mechanism for conv.To function, in a way that prevents import cycles, with following example:
// - encoding/json imports conv to register conversions between json.Value and other types (maps, slices, etc) via init
// - encoding/json imports conv to take advantage of conv.To
// - conv never immports encoding/json, so no import cycle
// - any other package can use conv.To to convert between json.Value and other types.
//
// SPDX-License-Identifier: Apache-2.0
package conv
