package json

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"

	"github.com/bantling/micro/conv"
)

// init registers conversions between json.Value and other types
func init() {
	// Object
	conv.RegisterConversion[map[string]any, Value](map[string]any{}, nil, func(src map[string]any, tgt *Value) (err error) { *tgt, err = MapToValue(src); return })
	conv.RegisterConversion[map[string]Value, Value](map[string]Value{}, nil, func(src map[string]Value, tgt *Value) (err error) { *tgt, err = MapToValue(src); return })

	// Array
	conv.RegisterConversion[[]any, Value]([]any{}, nil, func(src []any, tgt *Value) (err error) { *tgt, err = SliceToValue(src); return })
	conv.RegisterConversion[[]Value, Value]([]Value{}, nil, func(src []Value, tgt *Value) (err error) { *tgt, err = SliceToValue(src); return })

	// String
	conv.RegisterConversion[string, Value]("", nil, func(src string, tgt *Value) (err error) { *tgt = StringToValue(src); return })

	// NumberString
	conv.RegisterConversion[NumberString, Value](NumberString(""), nil, func(src NumberString, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })

	// Signed ints
	conv.RegisterConversion[int, Value](0, nil, func(src int, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.RegisterConversion[int8, Value](0, nil, func(src int8, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.RegisterConversion[int16, Value](0, nil, func(src int16, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.RegisterConversion[int32, Value](0, nil, func(src int32, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.RegisterConversion[int64, Value](0, nil, func(src int64, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })

	// Unsigned ints
	conv.RegisterConversion[uint, Value](0, nil, func(src uint, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.RegisterConversion[uint8, Value](0, nil, func(src uint8, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.RegisterConversion[uint16, Value](0, nil, func(src uint16, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.RegisterConversion[uint32, Value](0, nil, func(src uint32, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.RegisterConversion[uint64, Value](0, nil, func(src uint64, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })

	// Floats
	conv.RegisterConversion[float32, Value](0, nil, func(src float32, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.RegisterConversion[float64, Value](0, nil, func(src float64, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })

	// Bigs
	conv.RegisterConversion[*big.Float, Value](nil, nil, func(src *big.Float, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.RegisterConversion[*big.Int, Value](nil, nil, func(src *big.Int, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.RegisterConversion[*big.Rat, Value](nil, nil, func(src *big.Rat, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
}
