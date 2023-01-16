package json

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"
	"testing"

	"github.com/bantling/micro/go/util"
	"github.com/stretchr/testify/assert"
)

func TestFromStruct0Ptr_(t *testing.T) {
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
	}{
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
		}{
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

func TestFromStruct1Ptr_(t *testing.T) {
	type InnerStrukt struct {
		Foo *string
		Bar *int
	}

	var (
		foo          = "foo"
		bar          = 1
		bln          = true
		i            = 1
		i8   int8    = 2
		i16  int16   = 3
		i32  int32   = 4
		i64  int64   = 5
		ui   uint    = 6
		ui8  uint8   = 7
		ui16 uint16  = 8
		ui32 uint32  = 9
		ui64 uint64  = 10
		f32  float32 = 11.25
		f64  float64 = 12.5
		bi           = big.NewInt(13)
		bf           = big.NewFloat(14.5)
		br           = big.NewRat(15, 16)
	)

	strukt := struct {
		Str   *string
		Bln   *bool
		I     *int
		I8    *int8
		I16   *int16
		I32   *int32
		I64   *int64
		UI    *uint
		UI8   *uint8
		UI16  *uint16
		UI32  *uint32
		UI64  *uint64
		F32   *float32
		F64   *float64
		BI    **big.Int
		BF    **big.Float
		BR    **big.Rat
		Inner *InnerStrukt
	}{
		Str:  &foo,
		Bln:  &bln,
		I:    &i,
		I8:   &i8,
		I16:  &i16,
		I32:  &i32,
		I64:  &i64,
		UI:   &ui,
		UI8:  &ui8,
		UI16: &ui16,
		UI32: &ui32,
		UI64: &ui64,
		F32:  &f32,
		F64:  &f64,
		BI:   &bi,
		BF:   &bf,
		BR:   &br,
		Inner: &InnerStrukt{
			Foo: &foo,
			Bar: &bar,
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

func TestFromStruct2Ptr_(t *testing.T) {
	type InnerStrukt struct {
		Foo **string
		Bar **int
	}

	var (
		foo          = "foo"
		bar          = 1
		bln          = true
		i            = 1
		i8   int8    = 2
		i16  int16   = 3
		i32  int32   = 4
		i64  int64   = 5
		ui   uint    = 6
		ui8  uint8   = 7
		ui16 uint16  = 8
		ui32 uint32  = 9
		ui64 uint64  = 10
		f32  float32 = 11.25
		f64  float64 = 12.5
		bi           = big.NewInt(13)
		bf           = big.NewFloat(14.5)
		br           = big.NewRat(15, 16)

		foop  = &foo
		barp  = &bar
		blnp  = &bln
		ip    = &i
		i8p   = &i8
		i16p  = &i16
		i32p  = &i32
		i64p  = &i64
		uip   = &ui
		ui8p  = &ui8
		ui16p = &ui16
		ui32p = &ui32
		ui64p = &ui64
		f32p  = &f32
		f64p  = &f64
		bip   = &bi
		bfp   = &bf
		brp   = &br

		is = InnerStrukt{
			Foo: &foop,
			Bar: &barp,
		}
		isp = &is
	)

	strukt := struct {
		Str   **string
		Bln   **bool
		I     **int
		I8    **int8
		I16   **int16
		I32   **int32
		I64   **int64
		UI    **uint
		UI8   **uint8
		UI16  **uint16
		UI32  **uint32
		UI64  **uint64
		F32   **float32
		F64   **float64
		BI    ***big.Int
		BF    ***big.Float
		BR    ***big.Rat
		Inner **InnerStrukt
	}{
		Str:   &foop,
		Bln:   &blnp,
		I:     &ip,
		I8:    &i8p,
		I16:   &i16p,
		I32:   &i32p,
		I64:   &i64p,
		UI:    &uip,
		UI8:   &ui8p,
		UI16:  &ui16p,
		UI32:  &ui32p,
		UI64:  &ui64p,
		F32:   &f32p,
		F64:   &f64p,
		BI:    &bip,
		BF:    &bfp,
		BR:    &brp,
		Inner: &isp,
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
