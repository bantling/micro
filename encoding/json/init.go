package json

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"

	"github.com/bantling/micro/conv"
)

// init registers conversions between json.Value and other types
func init() {
	// Object
	conv.MustRegisterConversion(func(src map[string]any, tgt *Value) (err error) { *tgt, err = MapToValue(src); return })
	conv.MustRegisterConversion(func(src map[string]Value, tgt *Value) (err error) { *tgt, err = MapToValue(src); return })

	// Array
	conv.MustRegisterConversion(func(src []any, tgt *Value) (err error) { *tgt, err = SliceToValue(src); return })
	conv.MustRegisterConversion(func(src []Value, tgt *Value) (err error) { *tgt, err = SliceToValue(src); return })

	// String
	conv.MustRegisterConversion(func(src string, tgt *Value) (err error) { *tgt = StringToValue(src); return })

	// NumberString
	conv.MustRegisterConversion(func(src NumberString, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })

	// Signed ints
	conv.MustRegisterConversion(func(src int, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src int8, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src int16, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src int32, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src int64, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })

	// Unsigned ints
	conv.MustRegisterConversion(func(src uint, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src uint8, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src uint16, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src uint32, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src uint64, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })

	// Floats
	conv.MustRegisterConversion(func(src float32, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src float64, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })

	// Bigs
	conv.MustRegisterConversion(func(src *big.Float, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src *big.Int, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src *big.Rat, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
}
