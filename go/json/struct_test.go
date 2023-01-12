package json

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"
	"testing"

	"github.com/bantling/micro/go/util"
	"github.com/stretchr/testify/assert"
)

func TestFromStruct_(t *testing.T) {
	strukt := struct {
    Str   string
    Bln   bool
    I     int
    I8    int8
    I16   int16
    I32   int32
    I64   int64
    UI    uint
    UI8   uint8
    UI16  uint16
    UI32  uint32
    UI64  uint64
    F32   float32
    F64   float64
    BI    *big.Int
    BF    *big.Float
    BR    *big.Rat
    Inner struct {
      Foo string
      Bar int
    }
  } {
		Str:  "foo",
		Bln:  true,
		I:    1,
		I8:   2,
		I16:  3,
		I32:  4,
		I64:  5,
		UI:   6,
		UI8:  7,
		UI16: 8,
		UI32: 9,
		UI64: 10,
		F32:  11.25,
		F64:  12.5,
		BI:   big.NewInt(13),
		BF:   big.NewFloat(14.5),
		BR:   big.NewRat(15, 16),
		Inner: struct {
      Foo string
      Bar int
    } {
			Foo: "foo",
			Bar: 1,
		},
  }

	m := map[string]any{
		"str":  "foo",
		"bln":  true,
		"i":    1,
		"i8":   2,
		"i16":  3,
		"i32":  4,
		"i64":  5,
		"ui":   6,
		"ui8":  7,
		"ui16": 8,
		"ui32": 9,
		"ui64": 10,
		"f32":  11.25,
		"f64":  12.5,
		"bi":   13,
		"bf":   14.5,
		"br":   big.NewRat(15, 16),
		"inner": map[string]any{
			"foo": "foo",
			"bar": 1,
		},
	}

	assert.Equal(t, util.Of2Error(FromMap(m), nil), util.Of2Error(FromStruct(strukt)))
}
