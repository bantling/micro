package json

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"

	"github.com/bantling/micro/conv"
)

// init registers conversions between json.Value and other types
func init() {
	// Object
	conv.MustRegisterConversion[map[string]any, Value](map[string]any{}, nil, func(src map[string]any, tgt *Value) (err error) { *tgt, err = MapToValue(src); return })
	conv.MustRegisterConversion[map[string]Value, Value](map[string]Value{}, nil, func(src map[string]Value, tgt *Value) (err error) { *tgt, err = MapToValue(src); return })

	// Array
	conv.MustRegisterConversion[[]any, Value]([]any{}, nil, func(src []any, tgt *Value) (err error) { *tgt, err = SliceToValue(src); return })
	conv.MustRegisterConversion[[]Value, Value]([]Value{}, nil, func(src []Value, tgt *Value) (err error) { *tgt, err = SliceToValue(src); return })

	// String
	conv.MustRegisterConversion[string, Value]("", nil, func(src string, tgt *Value) (err error) { *tgt = StringToValue(src); return })

	// NumberString
	conv.MustRegisterConversion[NumberString, Value](NumberString(""), nil, func(src NumberString, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })

	// Signed ints
	conv.MustRegisterConversion[int, Value](0, nil, func(src int, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion[int8, Value](0, nil, func(src int8, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion[int16, Value](0, nil, func(src int16, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion[int32, Value](0, nil, func(src int32, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion[int64, Value](0, nil, func(src int64, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })

	// Unsigned ints
	conv.MustRegisterConversion[uint, Value](0, nil, func(src uint, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion[uint8, Value](0, nil, func(src uint8, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion[uint16, Value](0, nil, func(src uint16, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion[uint32, Value](0, nil, func(src uint32, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion[uint64, Value](0, nil, func(src uint64, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })

	// Floats
	conv.MustRegisterConversion[float32, Value](0, nil, func(src float32, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion[float64, Value](0, nil, func(src float64, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })

	// Bigs
	conv.MustRegisterConversion[*big.Float, Value](nil, nil, func(src *big.Float, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion[*big.Int, Value](nil, nil, func(src *big.Int, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion[*big.Rat, Value](nil, nil, func(src *big.Rat, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
}
