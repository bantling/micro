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
	conv.MustRegisterConversion(func(src Value, tgt *map[string]any) (err error) { *tgt = src.ToMap(); return })

	conv.MustRegisterConversion(func(src map[string]Value, tgt *Value) (err error) { *tgt, err = MapToValue(src); return })
	conv.MustRegisterConversion(func(src Value, tgt *map[string]Value) (err error) { *tgt = src.AsMap(); return })

	// Array
	conv.MustRegisterConversion(func(src []any, tgt *Value) (err error) { *tgt, err = SliceToValue(src); return })
	conv.MustRegisterConversion(func(src Value, tgt *[]any) (err error) { *tgt = src.ToSlice(); return })

	conv.MustRegisterConversion(func(src []Value, tgt *Value) (err error) { *tgt, err = SliceToValue(src); return })
	conv.MustRegisterConversion(func(src Value, tgt *[]Value) (err error) { *tgt = src.AsSlice(); return })

	// String
	conv.MustRegisterConversion(func(src string, tgt *Value) (err error) { *tgt = StringToValue(src); return })
	conv.MustRegisterConversion(func(src Value, tgt *string) (err error) { *tgt = src.AsString(); return })

	// NumberString
	conv.MustRegisterConversion(func(src NumberString, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src Value, tgt *NumberString) (err error) { *tgt = src.AsNumber(); return })

	// Signed ints
	conv.MustRegisterConversion(func(src int, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src int8, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src int16, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src int32, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src int64, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })

	conv.MustRegisterConversion(func(src Value, tgt *int) error { ns := src.AsNumber(); return conv.To(ns, tgt) })
	conv.MustRegisterConversion(func(src Value, tgt *int8) error { ns := src.AsNumber(); return conv.To(ns, tgt) })
	conv.MustRegisterConversion(func(src Value, tgt *int16) error { ns := src.AsNumber(); return conv.To(ns, tgt) })
	conv.MustRegisterConversion(func(src Value, tgt *int32) error { ns := src.AsNumber(); return conv.To(ns, tgt) })
	conv.MustRegisterConversion(func(src Value, tgt *int64) error { ns := src.AsNumber(); return conv.To(ns, tgt) })

	// Unsigned ints
	conv.MustRegisterConversion(func(src uint, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src uint8, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src uint16, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src uint32, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src uint64, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })

	conv.MustRegisterConversion(func(src Value, tgt *uint) error { ns := src.AsNumber(); return conv.To(ns, tgt) })
	conv.MustRegisterConversion(func(src Value, tgt *uint8) error { ns := src.AsNumber(); return conv.To(ns, tgt) })
	conv.MustRegisterConversion(func(src Value, tgt *uint16) error { ns := src.AsNumber(); return conv.To(ns, tgt) })
	conv.MustRegisterConversion(func(src Value, tgt *uint32) error { ns := src.AsNumber(); return conv.To(ns, tgt) })
	conv.MustRegisterConversion(func(src Value, tgt *uint64) error { ns := src.AsNumber(); return conv.To(ns, tgt) })

	// Floats
	conv.MustRegisterConversion(func(src float32, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src float64, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })

	conv.MustRegisterConversion(func(src Value, tgt *float32) (err error) { ns := src.AsNumber(); return conv.To(ns, tgt) })
	conv.MustRegisterConversion(func(src Value, tgt *float64) (err error) { ns := src.AsNumber(); return conv.To(ns, tgt) })

	// Bigs
	conv.MustRegisterConversion(func(src *big.Float, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src Value, tgt **big.Float) (err error) {
		ns := src.AsNumber()
		return conv.StringToBigFloat(string(ns), tgt)
	})

	conv.MustRegisterConversion(func(src *big.Int, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src Value, tgt **big.Int) (err error) {
		ns := src.AsNumber()
		return conv.StringToBigInt(string(ns), tgt)
	})

	conv.MustRegisterConversion(func(src *big.Rat, tgt *Value) (err error) { *tgt, err = NumberToValue(src); return })
	conv.MustRegisterConversion(func(src Value, tgt **big.Rat) (err error) {
		ns := src.AsNumber()
		return conv.StringToBigRat(string(ns), tgt)
	})
}
