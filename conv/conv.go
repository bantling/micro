package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"
	goreflect "reflect"
	"strconv"

	"github.com/bantling/micro/constraint"
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/reflect"
)

var (
	errMsg = "The %T value of %s cannot be converted to %s"

	log2Of10 = math.Log2(10)

	minIntValue = map[int]int{
		8:  math.MinInt8,
		16: math.MinInt16,
		32: math.MinInt32,
		64: math.MinInt64,
	}

	maxIntValue = map[int]int{
		8:  math.MaxInt8,
		16: math.MaxInt16,
		32: math.MaxInt32,
		64: math.MaxInt64,
	}

	maxUintValue = map[int]uint{
		8:  math.MaxUint8,
		16: math.MaxUint16,
		32: math.MaxUint32,
		64: math.MaxUint64,
	}

	// map strings of from/to conversion pairs to func(any, any) error that perform the specified conversion
	// no map entries are provided for from/to pairs where from and to are the same type.
	convertFromTo = map[string]func(any, any) error{
		// ==== To int
		"int8int": func(t any, u any) error {
			*(u.(*int)) = int(t.(int8))
			return nil
		},
		"int16int": func(t any, u any) error {
			*(u.(*int)) = int(t.(int16))
			return nil
		},
		"int32int": func(t any, u any) error {
			*(u.(*int)) = int(t.(int32))
			return nil
		},
		"int64int": func(t any, u any) error {
			return IntToInt(t.(int64), u.(*int))
		},
		"uintint": func(t any, u any) error {
			return UintToInt(t.(uint), u.(*int))
		},
		"uint8int": func(t any, u any) error {
			*(u.(*int)) = int(t.(uint8))
			return nil
		},
		"uint16int": func(t any, u any) error {
			*(u.(*int)) = int(t.(uint16))
			return nil
		},
		"uint32int": func(t any, u any) error {
			return UintToInt(t.(uint32), u.(*int))
		},
		"uint64int": func(t any, u any) error {
			return UintToInt(t.(uint64), u.(*int))
		},
		"float32int": func(t any, u any) error {
			return FloatToInt(t.(float32), u.(*int))
		},
		"float64int": func(t any, u any) error {
			return FloatToInt(t.(float64), u.(*int))
		},
		"*big.Intint": func(t any, u any) error {
			var inter int64
			if err := BigIntToInt64(t.(*big.Int), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int))
		},
		"*big.Floatint": func(t any, u any) error {
			var inter int64
			if err := BigFloatToInt64(t.(*big.Float), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int))
		},
		"*big.Ratint": func(t any, u any) error {
			var inter int64
			if err := BigRatToInt64(t.(*big.Rat), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int))
		},
		"stringint": func(t any, u any) error {
			var inter int64
			if err := StringToInt64(t.(string), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int))
		},

		// ==== To int8
		"intint8": func(t any, u any) error {
			return IntToInt(t.(int), u.(*int8))
		},
		"int16int8": func(t any, u any) error {
			return IntToInt(t.(int16), u.(*int8))
		},
		"int32int8": func(t any, u any) error {
			return IntToInt(t.(int32), u.(*int8))
		},
		"int64int8": func(t any, u any) error {
			return IntToInt(t.(int64), u.(*int8))
		},
		"uintint8": func(t any, u any) error {
			return UintToInt(t.(uint), u.(*int8))
		},
		"uint8int8": func(t any, u any) error {
			return UintToInt(t.(uint8), u.(*int8))
		},
		"uint16int8": func(t any, u any) error {
			return UintToInt(t.(uint16), u.(*int8))
		},
		"uint32int8": func(t any, u any) error {
			return UintToInt(t.(uint32), u.(*int8))
		},
		"uint64int8": func(t any, u any) error {
			return UintToInt(t.(uint64), u.(*int8))
		},
		"float32int8": func(t any, u any) error {
			return FloatToInt(t.(float32), u.(*int8))
		},
		"float64int8": func(t any, u any) error {
			return FloatToInt(t.(float64), u.(*int8))
		},
		"*big.Intint8": func(t any, u any) error {
			var inter int64
			if err := BigIntToInt64(t.(*big.Int), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int8))
		},
		"*big.Floatint8": func(t any, u any) error {
			var inter int64
			if err := BigFloatToInt64(t.(*big.Float), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int8))
		},
		"*big.Ratint8": func(t any, u any) error {
			var inter int64
			if err := BigRatToInt64(t.(*big.Rat), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int8))
		},
		"stringint8": func(t any, u any) error {
			var inter int64
			if err := StringToInt64(t.(string), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int8))
		},

		// ==== To int16
		"intint16": func(t any, u any) error {
			return IntToInt(t.(int), u.(*int16))
		},
		"int8int16": func(t any, u any) error {
			*(u.(*int16)) = int16(t.(int8))
			return nil
		},
		"int32int16": func(t any, u any) error {
			return IntToInt(t.(int32), u.(*int16))
		},
		"int64int16": func(t any, u any) error {
			return IntToInt(t.(int64), u.(*int16))
		},
		"uintint16": func(t any, u any) error {
			return UintToInt(t.(uint), u.(*int16))
		},
		"uint8int16": func(t any, u any) error {
			*(u.(*int16)) = int16(t.(uint8))
			return nil
		},
		"uint16int16": func(t any, u any) error {
			return UintToInt(t.(uint16), u.(*int16))
		},
		"uint32int16": func(t any, u any) error {
			return UintToInt(t.(uint32), u.(*int16))
		},
		"uint64int16": func(t any, u any) error {
			return UintToInt(t.(uint64), u.(*int16))
		},
		"float32int16": func(t any, u any) error {
			return FloatToInt(t.(float32), u.(*int16))
		},
		"float64int16": func(t any, u any) error {
			return FloatToInt(t.(float64), u.(*int16))
		},
		"*big.Intint16": func(t any, u any) error {
			var inter int64
			if err := BigIntToInt64(t.(*big.Int), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int16))
		},
		"*big.Floatint16": func(t any, u any) error {
			var inter int64
			if err := BigFloatToInt64(t.(*big.Float), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int16))
		},
		"*big.Ratint16": func(t any, u any) error {
			var inter int64
			if err := BigRatToInt64(t.(*big.Rat), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int16))
		},
		"stringint16": func(t any, u any) error {
			var inter int64
			if err := StringToInt64(t.(string), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int16))
		},

		// ==== To int32
		"intint32": func(t any, u any) error {
			return IntToInt(t.(int), u.(*int32))
		},
		"int8int32": func(t any, u any) error {
			*(u.(*int32)) = int32(t.(int8))
			return nil
		},
		"int16int32": func(t any, u any) error {
			*(u.(*int32)) = int32(t.(int16))
			return nil
		},
		"int64int32": func(t any, u any) error {
			return IntToInt(t.(int64), u.(*int32))
		},
		"uintint32": func(t any, u any) error {
			return UintToInt(t.(uint), u.(*int32))
		},
		"uint8int32": func(t any, u any) error {
			*(u.(*int32)) = int32(t.(uint8))
			return nil
		},
		"uint16int32": func(t any, u any) error {
			*(u.(*int32)) = int32(t.(uint16))
			return nil
		},
		"uint32int32": func(t any, u any) error {
			return UintToInt(t.(uint32), u.(*int32))
		},
		"uint64int32": func(t any, u any) error {
			return UintToInt(t.(uint64), u.(*int32))
		},
		"float32int32": func(t any, u any) error {
			return FloatToInt(t.(float32), u.(*int32))
		},
		"float64int32": func(t any, u any) error {
			return FloatToInt(t.(float64), u.(*int32))
		},
		"*big.Intint32": func(t any, u any) error {
			var inter int64
			if err := BigIntToInt64(t.(*big.Int), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int32))
		},
		"*big.Floatint32": func(t any, u any) error {
			var inter int64
			if err := BigFloatToInt64(t.(*big.Float), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int32))
		},
		"*big.Ratint32": func(t any, u any) error {
			var inter int64
			if err := BigRatToInt64(t.(*big.Rat), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int32))
		},
		"stringint32": func(t any, u any) error {
			var inter int64
			if err := StringToInt64(t.(string), &inter); err != nil {
				return err
			}
			return IntToInt(inter, u.(*int32))
		},

		// ==== To int64
		"intint64": func(t any, u any) error {
			*(u.(*int64)) = int64(t.(int))
			return nil
		},
		"int8int64": func(t any, u any) error {
			*(u.(*int64)) = int64(t.(int8))
			return nil
		},
		"int16int64": func(t any, u any) error {
			*(u.(*int64)) = int64(t.(int16))
			return nil
		},
		"int32int64": func(t any, u any) error {
			*(u.(*int64)) = int64(t.(int32))
			return nil
		},
		"uintint64": func(t any, u any) error {
			return UintToInt(t.(uint), u.(*int64))
		},
		"uint8int64": func(t any, u any) error {
			*(u.(*int64)) = int64(t.(uint8))
			return nil
		},
		"uint16int64": func(t any, u any) error {
			*(u.(*int64)) = int64(t.(uint16))
			return nil
		},
		"uint32int64": func(t any, u any) error {
			*(u.(*int64)) = int64(t.(uint32))
			return nil
		},
		"uint64int64": func(t any, u any) error {
			return UintToInt(t.(uint64), u.(*int64))
		},
		"float32int64": func(t any, u any) error {
			return FloatToInt(t.(float32), u.(*int64))
		},
		"float64int64": func(t any, u any) error {
			return FloatToInt(t.(float64), u.(*int64))
		},
		"*big.Intint64": func(t any, u any) error {
			return BigIntToInt64(t.(*big.Int), u.(*int64))
		},
		"*big.Floatint64": func(t any, u any) error {
			return BigFloatToInt64(t.(*big.Float), u.(*int64))
		},
		"*big.Ratint64": func(t any, u any) error {
			return BigRatToInt64(t.(*big.Rat), u.(*int64))
		},
		"stringint64": func(t any, u any) error {
			return StringToInt64(t.(string), u.(*int64))
		},

		// ==== To uint
		"intuint": func(t any, u any) error {
			return IntToUint(t.(int), u.(*uint))
		},
		"int8uint": func(t any, u any) error {
			return IntToUint(t.(int8), u.(*uint))
		},
		"int16uint": func(t any, u any) error {
			return IntToUint(t.(int16), u.(*uint))
		},
		"int32uint": func(t any, u any) error {
			return IntToUint(t.(int32), u.(*uint))
		},
		"int64uint": func(t any, u any) error {
			return IntToUint(t.(int64), u.(*uint))
		},
		"uint8uint": func(t any, u any) error {
			*(u.(*uint)) = uint(t.(uint8))
			return nil
		},
		"uint16uint": func(t any, u any) error {
			*(u.(*uint)) = uint(t.(uint16))
			return nil
		},
		"uint32uint": func(t any, u any) error {
			*(u.(*uint)) = uint(t.(uint32))
			return nil
		},
		"uint64uint": func(t any, u any) error {
			return UintToUint(t.(uint64), u.(*uint))
		},
		"float32uint": func(t any, u any) error {
			return FloatToUint(t.(float32), u.(*uint))
		},
		"float64uint": func(t any, u any) error {
			return FloatToUint(t.(float64), u.(*uint))
		},
		"*big.Intuint": func(t any, u any) error {
			var inter uint64
			if err := BigIntToUint64(t.(*big.Int), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint))
		},
		"*big.Floatuint": func(t any, u any) error {
			var inter uint64
			if err := BigFloatToUint64(t.(*big.Float), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint))
		},
		"*big.Ratuint": func(t any, u any) error {
			var inter uint64
			if err := BigRatToUint64(t.(*big.Rat), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint))
		},
		"stringuint": func(t any, u any) error {
			var inter uint64
			if err := StringToUint64(t.(string), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint))
		},

		// ==== To uint8
		"intuint8": func(t any, u any) error {
			return IntToUint(t.(int), u.(*uint8))
		},
		"int8uint8": func(t any, u any) error {
			return IntToUint(t.(int8), u.(*uint8))
		},
		"int16uint8": func(t any, u any) error {
			return IntToUint(t.(int16), u.(*uint8))
		},
		"int32uint8": func(t any, u any) error {
			return IntToUint(t.(int32), u.(*uint8))
		},
		"int64uint8": func(t any, u any) error {
			return IntToUint(t.(int64), u.(*uint8))
		},
		"uintuint8": func(t any, u any) error {
			return UintToUint(t.(uint), u.(*uint8))
		},
		"uint16uint8": func(t any, u any) error {
			return UintToUint(t.(uint16), u.(*uint8))
		},
		"uint32uint8": func(t any, u any) error {
			return UintToUint(t.(uint32), u.(*uint8))
		},
		"uint64uint8": func(t any, u any) error {
			return UintToUint(t.(uint64), u.(*uint8))
		},
		"float32uint8": func(t any, u any) error {
			return FloatToUint(t.(float32), u.(*uint8))
		},
		"float64uint8": func(t any, u any) error {
			return FloatToUint(t.(float64), u.(*uint8))
		},
		"*big.Intuint8": func(t any, u any) error {
			var inter uint64
			if err := BigIntToUint64(t.(*big.Int), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint8))
		},
		"*big.Floatuint8": func(t any, u any) error {
			var inter uint64
			if err := BigFloatToUint64(t.(*big.Float), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint8))
		},
		"*big.Ratuint8": func(t any, u any) error {
			var inter uint64
			if err := BigRatToUint64(t.(*big.Rat), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint8))
		},
		"stringuint8": func(t any, u any) error {
			var inter uint64
			if err := StringToUint64(t.(string), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint8))
		},

		// ==== To uint16
		"intuint16": func(t any, u any) error {
			return IntToUint(t.(int), u.(*uint16))
		},
		"int8uint16": func(t any, u any) error {
			return IntToUint(t.(int8), u.(*uint16))
		},
		"int16uint16": func(t any, u any) error {
			return IntToUint(t.(int16), u.(*uint16))
		},
		"int32uint16": func(t any, u any) error {
			return IntToUint(t.(int32), u.(*uint16))
		},
		"int64uint16": func(t any, u any) error {
			return IntToUint(t.(int64), u.(*uint16))
		},
		"uintuint16": func(t any, u any) error {
			return UintToUint(t.(uint), u.(*uint16))
		},
		"uint8uint16": func(t any, u any) error {
			*(u.(*uint16)) = uint16(t.(uint8))
			return nil
		},
		"uint32uint16": func(t any, u any) error {
			return UintToUint(t.(uint32), u.(*uint16))
		},
		"uint64uint16": func(t any, u any) error {
			return UintToUint(t.(uint64), u.(*uint16))
		},
		"float32uint16": func(t any, u any) error {
			return FloatToUint(t.(float32), u.(*uint16))
		},
		"float64uint16": func(t any, u any) error {
			return FloatToUint(t.(float64), u.(*uint16))
		},
		"*big.Intuint16": func(t any, u any) error {
			var inter uint64
			if err := BigIntToUint64(t.(*big.Int), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint16))
		},
		"*big.Floatuint16": func(t any, u any) error {
			var inter uint64
			if err := BigFloatToUint64(t.(*big.Float), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint16))
		},
		"*big.Ratuint16": func(t any, u any) error {
			var inter uint64
			if err := BigRatToUint64(t.(*big.Rat), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint16))
		},
		"stringuint16": func(t any, u any) error {
			var inter uint64
			if err := StringToUint64(t.(string), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint16))
		},

		// ==== To uint32
		"intuint32": func(t any, u any) error {
			return IntToUint(t.(int), u.(*uint32))
		},
		"int8uint32": func(t any, u any) error {
			return IntToUint(t.(int8), u.(*uint32))
		},
		"int16uint32": func(t any, u any) error {
			return IntToUint(t.(int16), u.(*uint32))
		},
		"int32uint32": func(t any, u any) error {
			return IntToUint(t.(int32), u.(*uint32))
		},
		"int64uint32": func(t any, u any) error {
			return IntToUint(t.(int64), u.(*uint32))
		},
		"uintuint32": func(t any, u any) error {
			return UintToUint(t.(uint), u.(*uint32))
		},
		"uint8uint32": func(t any, u any) error {
			*(u.(*uint32)) = uint32(t.(uint8))
			return nil
		},
		"uint16uint32": func(t any, u any) error {
			*(u.(*uint32)) = uint32(t.(uint16))
			return nil
		},
		"uint64uint32": func(t any, u any) error {
			return UintToUint(t.(uint64), u.(*uint32))
		},
		"float32uint32": func(t any, u any) error {
			return FloatToUint(t.(float32), u.(*uint32))
		},
		"float64uint32": func(t any, u any) error {
			return FloatToUint(t.(float64), u.(*uint32))
		},
		"*big.Intuint32": func(t any, u any) error {
			var inter uint64
			if err := BigIntToUint64(t.(*big.Int), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint32))
		},
		"*big.Floatuint32": func(t any, u any) error {
			var inter uint64
			if err := BigFloatToUint64(t.(*big.Float), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint32))
		},
		"*big.Ratuint32": func(t any, u any) error {
			var inter uint64
			if err := BigRatToUint64(t.(*big.Rat), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint32))
		},
		"stringuint32": func(t any, u any) error {
			var inter uint64
			if err := StringToUint64(t.(string), &inter); err != nil {
				return err
			}
			return UintToUint(inter, u.(*uint32))
		},

		// ==== To uint64
		"intuint64": func(t any, u any) error {
			return IntToUint(t.(int), u.(*uint64))
		},
		"int8uint64": func(t any, u any) error {
			return IntToUint(t.(int8), u.(*uint64))
		},
		"int16uint64": func(t any, u any) error {
			return IntToUint(t.(int16), u.(*uint64))
		},
		"int32uint64": func(t any, u any) error {
			return IntToUint(t.(int32), u.(*uint64))
		},
		"int64uint64": func(t any, u any) error {
			return IntToUint(t.(int64), u.(*uint64))
		},
		"uintuint64": func(t any, u any) error {
			*(u.(*uint64)) = uint64(t.(uint))
			return nil
		},
		"uint8uint64": func(t any, u any) error {
			*(u.(*uint64)) = uint64(t.(uint8))
			return nil
		},
		"uint16uint64": func(t any, u any) error {
			*(u.(*uint64)) = uint64(t.(uint16))
			return nil
		},
		"uint32uint64": func(t any, u any) error {
			*(u.(*uint64)) = uint64(t.(uint32))
			return nil
		},
		"float32uint64": func(t any, u any) error {
			return FloatToUint(t.(float32), u.(*uint64))
		},
		"float64uint64": func(t any, u any) error {
			return FloatToUint(t.(float64), u.(*uint64))
		},
		"*big.Intuint64": func(t any, u any) error {
			return BigIntToUint64(t.(*big.Int), u.(*uint64))
		},
		"*big.Floatuint64": func(t any, u any) error {
			return BigFloatToUint64(t.(*big.Float), u.(*uint64))
		},
		"*big.Ratuint64": func(t any, u any) error {
			return BigRatToUint64(t.(*big.Rat), u.(*uint64))
		},
		"stringuint64": func(t any, u any) error {
			return StringToUint64(t.(string), u.(*uint64))
		},

		// ==== To float32
		"intfloat32": func(t any, u any) error {
			return IntToFloat(t.(int), u.(*float32))
		},
		"int8float32": func(t any, u any) error {
			*(u.(*float32)) = float32(t.(int8))
			return nil
		},
		"int16float32": func(t any, u any) error {
			*(u.(*float32)) = float32(t.(int16))
			return nil
		},
		"int32float32": func(t any, u any) error {
			return IntToFloat(t.(int32), u.(*float32))
		},
		"int64float32": func(t any, u any) error {
			return IntToFloat(t.(int64), u.(*float32))
		},
		"uintfloat32": func(t any, u any) error {
			return IntToFloat(t.(uint), u.(*float32))
		},
		"uint8float32": func(t any, u any) error {
			*(u.(*float32)) = float32(t.(uint8))
			return nil
		},
		"uint16float32": func(t any, u any) error {
			*(u.(*float32)) = float32(t.(uint16))
			return nil
		},
		"uint32float32": func(t any, u any) error {
			return IntToFloat(t.(uint32), u.(*float32))
		},
		"uint64float32": func(t any, u any) error {
			return IntToFloat(t.(uint64), u.(*float32))
		},
		"float64float32": func(t any, u any) error {
			return FloatToFloat(t.(float64), u.(*float32))
		},
		"*big.Intfloat32": func(t any, u any) error {
			var inter float64
			if err := BigIntToFloat64(t.(*big.Int), &inter); err != nil {
				return err
			}
			return FloatToFloat(inter, u.(*float32))
		},
		"*big.Floatfloat32": func(t any, u any) error {
			var inter float64
			if err := BigFloatToFloat64(t.(*big.Float), &inter); err != nil {
				return err
			}
			return FloatToFloat(inter, u.(*float32))
		},
		"*big.Ratfloat32": func(t any, u any) error {
			var inter float64
			if err := BigRatToFloat64(t.(*big.Rat), &inter); err != nil {
				return err
			}
			return FloatToFloat(inter, u.(*float32))
		},
		"stringfloat32": func(t any, u any) error {
			return StringToFloat32(t.(string), u.(*float32))
		},

		// ==== To float64
		"intfloat64": func(t any, u any) error {
			return IntToFloat(t.(int), u.(*float64))
		},
		"int8float64": func(t any, u any) error {
			*(u.(*float64)) = float64(t.(int8))
			return nil
		},
		"int16float64": func(t any, u any) error {
			*(u.(*float64)) = float64(t.(int16))
			return nil
		},
		"int32float64": func(t any, u any) error {
			*(u.(*float64)) = float64(t.(int32))
			return nil
		},
		"int64float64": func(t any, u any) error {
			return IntToFloat(t.(int64), u.(*float64))
		},
		"uintfloat64": func(t any, u any) error {
			return IntToFloat(t.(uint), u.(*float64))
		},
		"uint8float64": func(t any, u any) error {
			*(u.(*float64)) = float64(t.(uint8))
			return nil
		},
		"uint16float64": func(t any, u any) error {
			*(u.(*float64)) = float64(t.(uint16))
			return nil
		},
		"uint32float64": func(t any, u any) error {
			*(u.(*float64)) = float64(t.(uint32))
			return nil
		},
		"uint64float64": func(t any, u any) error {
			return IntToFloat(t.(uint64), u.(*float64))
		},
		"float32float64": func(t any, u any) error {
			*(u.(*float64)) = float64(t.(float32))
			return nil
		},
		"*big.Intfloat64": func(t any, u any) error {
			return BigIntToFloat64(t.(*big.Int), u.(*float64))
		},
		"*big.Floatfloat64": func(t any, u any) error {
			return BigFloatToFloat64(t.(*big.Float), u.(*float64))
		},
		"*big.Ratfloat64": func(t any, u any) error {
			return BigRatToFloat64(t.(*big.Rat), u.(*float64))
		},
		"stringfloat64": func(t any, u any) error {
			return StringToFloat64(t.(string), u.(*float64))
		},

		// ==== To *big.Int
		"int*big.Int": func(t any, u any) error {
			IntToBigInt(t.(int), u.(**big.Int))
			return nil
		},
		"int8*big.Int": func(t any, u any) error {
			IntToBigInt(t.(int8), u.(**big.Int))
			return nil
		},
		"int16*big.Int": func(t any, u any) error {
			IntToBigInt(t.(int16), u.(**big.Int))
			return nil
		},
		"int32*big.Int": func(t any, u any) error {
			IntToBigInt(t.(int32), u.(**big.Int))
			return nil
		},
		"int64*big.Int": func(t any, u any) error {
			IntToBigInt(t.(int64), u.(**big.Int))
			return nil
		},
		"uint*big.Int": func(t any, u any) error {
			UintToBigInt(t.(uint), u.(**big.Int))
			return nil
		},
		"uint8*big.Int": func(t any, u any) error {
			UintToBigInt(t.(uint8), u.(**big.Int))
			return nil
		},
		"uint16*big.Int": func(t any, u any) error {
			UintToBigInt(t.(uint16), u.(**big.Int))
			return nil
		},
		"uint32*big.Int": func(t any, u any) error {
			UintToBigInt(t.(uint32), u.(**big.Int))
			return nil
		},
		"uint64*big.Int": func(t any, u any) error {
			UintToBigInt(t.(uint64), u.(**big.Int))
			return nil
		},
		"float32*big.Int": func(t any, u any) error {
			return FloatToBigInt(t.(float32), u.(**big.Int))
		},
		"float64*big.Int": func(t any, u any) error {
			return FloatToBigInt(t.(float64), u.(**big.Int))
		},
		"*big.Int*big.Int": func(t any, u any) error {
			BigIntToBigInt(t.(*big.Int), u.(**big.Int))
			return nil
		},
		"*big.Float*big.Int": func(t any, u any) error {
			return BigFloatToBigInt(t.(*big.Float), u.(**big.Int))
		},
		"*big.Rat*big.Int": func(t any, u any) error {
			return BigRatToBigInt(t.(*big.Rat), u.(**big.Int))
		},
		"string*big.Int": func(t any, u any) error {
			return StringToBigInt(t.(string), u.(**big.Int))
		},

		// ==== To *big.Float
		"int*big.Float": func(t any, u any) error {
			IntToBigFloat(t.(int), u.(**big.Float))
			return nil
		},
		"int8*big.Float": func(t any, u any) error {
			IntToBigFloat(t.(int8), u.(**big.Float))
			return nil
		},
		"int16*big.Float": func(t any, u any) error {
			IntToBigFloat(t.(int16), u.(**big.Float))
			return nil
		},
		"int32*big.Float": func(t any, u any) error {
			IntToBigFloat(t.(int32), u.(**big.Float))
			return nil
		},
		"int64*big.Float": func(t any, u any) error {
			IntToBigFloat(t.(int64), u.(**big.Float))
			return nil
		},
		"uint*big.Float": func(t any, u any) error {
			UintToBigFloat(t.(uint), u.(**big.Float))
			return nil
		},
		"uint8*big.Float": func(t any, u any) error {
			UintToBigFloat(t.(uint8), u.(**big.Float))
			return nil
		},
		"uint16*big.Float": func(t any, u any) error {
			UintToBigFloat(t.(uint16), u.(**big.Float))
			return nil
		},
		"uint32*big.Float": func(t any, u any) error {
			UintToBigFloat(t.(uint32), u.(**big.Float))
			return nil
		},
		"uint64*big.Float": func(t any, u any) error {
			UintToBigFloat(t.(uint64), u.(**big.Float))
			return nil
		},
		"float32*big.Float": func(t any, u any) error {
			return FloatToBigFloat(t.(float32), u.(**big.Float))
		},
		"float64*big.Float": func(t any, u any) error {
			return FloatToBigFloat(t.(float64), u.(**big.Float))
		},
		"*big.Int*big.Float": func(t any, u any) error {
			BigIntToBigFloat(t.(*big.Int), u.(**big.Float))
			return nil
		},
		"*big.Float*big.Float": func(t any, u any) error {
			BigFloatToBigFloat(t.(*big.Float), u.(**big.Float))
			return nil
		},
		"*big.Rat*big.Float": func(t any, u any) error {
			BigRatToBigFloat(t.(*big.Rat), u.(**big.Float))
			return nil
		},
		"string*big.Float": func(t any, u any) error {
			return StringToBigFloat(t.(string), u.(**big.Float))
		},

		// ==== To *big.Rat
		"int*big.Rat": func(t any, u any) error {
			IntToBigRat(t.(int), u.(**big.Rat))
			return nil
		},
		"int8*big.Rat": func(t any, u any) error {
			IntToBigRat(t.(int8), u.(**big.Rat))
			return nil
		},
		"int16*big.Rat": func(t any, u any) error {
			IntToBigRat(t.(int16), u.(**big.Rat))
			return nil
		},
		"int32*big.Rat": func(t any, u any) error {
			IntToBigRat(t.(int32), u.(**big.Rat))
			return nil
		},
		"int64*big.Rat": func(t any, u any) error {
			IntToBigRat(t.(int64), u.(**big.Rat))
			return nil
		},
		"uint*big.Rat": func(t any, u any) error {
			UintToBigRat(t.(uint), u.(**big.Rat))
			return nil
		},
		"uint8*big.Rat": func(t any, u any) error {
			UintToBigRat(t.(uint8), u.(**big.Rat))
			return nil
		},
		"uint16*big.Rat": func(t any, u any) error {
			UintToBigRat(t.(uint16), u.(**big.Rat))
			return nil
		},
		"uint32*big.Rat": func(t any, u any) error {
			UintToBigRat(t.(uint32), u.(**big.Rat))
			return nil
		},
		"uint64*big.Rat": func(t any, u any) error {
			UintToBigRat(t.(uint64), u.(**big.Rat))
			return nil
		},
		"float32*big.Rat": func(t any, u any) error {
			return FloatToBigRat(t.(float32), u.(**big.Rat))
		},
		"float64*big.Rat": func(t any, u any) error {
			return FloatToBigRat(t.(float64), u.(**big.Rat))
		},
		"*big.Int*big.Rat": func(t any, u any) error {
			BigIntToBigRat(t.(*big.Int), u.(**big.Rat))
			return nil
		},
		"*big.Float*big.Rat": func(t any, u any) error {
			return BigFloatToBigRat(t.(*big.Float), u.(**big.Rat))
		},
		"*big.Rat*big.Rat": func(t any, u any) error {
			BigRatToBigRat(t.(*big.Rat), u.(**big.Rat))
			return nil
		},
		"string*big.Rat": func(t any, u any) error {
			return StringToBigRat(t.(string), u.(**big.Rat))
		},

		// ==== To string
		"intstring": func(t any, u any) error {
			*(u.(*string)) = IntToString(t.(int))
			return nil
		},
		"int8string": func(t any, u any) error {
			*(u.(*string)) = IntToString(t.(int8))
			return nil
		},
		"int16string": func(t any, u any) error {
			*(u.(*string)) = IntToString(t.(int16))
			return nil
		},
		"int32string": func(t any, u any) error {
			*(u.(*string)) = IntToString(t.(int32))
			return nil
		},
		"int64string": func(t any, u any) error {
			*(u.(*string)) = IntToString(t.(int64))
			return nil
		},
		"uintstring": func(t any, u any) error {
			*(u.(*string)) = UintToString(t.(uint))
			return nil
		},
		"uint8string": func(t any, u any) error {
			*(u.(*string)) = UintToString(t.(uint8))
			return nil
		},
		"uint16string": func(t any, u any) error {
			*(u.(*string)) = UintToString(t.(uint16))
			return nil
		},
		"uint32string": func(t any, u any) error {
			*(u.(*string)) = UintToString(t.(uint32))
			return nil
		},
		"uint64string": func(t any, u any) error {
			*(u.(*string)) = UintToString(t.(uint64))
			return nil
		},
		"float32string": func(t any, u any) error {
			*(u.(*string)) = FloatToString(t.(float32))
			return nil
		},
		"float64string": func(t any, u any) error {
			*(u.(*string)) = FloatToString(t.(float64))
			return nil
		},
		"*big.Intstring": func(t any, u any) error {
			*(u.(*string)) = BigIntToString(t.(*big.Int))
			return nil
		},
		"*big.Floatstring": func(t any, u any) error {
			*(u.(*string)) = BigFloatToString(t.(*big.Float))
			return nil
		},
		"*big.Ratstring": func(t any, u any) error {
			*(u.(*string)) = BigRatToString(t.(*big.Rat))
			return nil
		},
	}
)

// ToString

// IntToString converts any signed int type into a string
func IntToString[T constraint.SignedInteger](val T) string {
	return strconv.FormatInt(int64(val), 10)
}

// UintToString converts any unsigned int type into a string
func UintToString[T constraint.UnsignedInteger](val T) string {
	return strconv.FormatUint(uint64(val), 10)
}

// FloatToString converts any float type into a string
func FloatToString[T constraint.Float](val T) string {
	_, is32 := any(val).(float32)
	return strconv.FormatFloat(float64(val), 'f', -1, funcs.Ternary(is32, 32, 64))
}

// BigIntToString converts a *big.Int to a string
func BigIntToString(val *big.Int) string {
	return val.String()
}

// BigFloatToString converts a *big.Float to a string
func BigFloatToString(val *big.Float) string {
	return val.String()
}

// BigRatToString converts a *big.Rat to a string.
// The string will be a ratio like 5/4, if it is int it will be a ratio like 5/1.
func BigRatToString(val *big.Rat) string {
	return val.String()
}

// BigRatToNormalizedString converts a *big.Rat to a string.
// The string will be formatted like an integer if the ratio is an int, else formatted like a float if it is not an int.
func BigRatToNormalizedString(val *big.Rat) string {
	if val.IsInt() {
		return val.Num().String()
	}

	var inter *big.Float
	BigRatToBigFloat(val, &inter)
	return inter.String()
}

// Converts any signed or unsigned int type, any float type, *big.Int, *big.Float, or *big.Rat to a string.
// The *big.Rat string will be normalized (see BigRatToNormalizedString).
func ToString[T constraint.Numeric](val T) string {
	if v, isa := any(val).(int); isa {
		return IntToString(v)
	} else if v, isa := any(val).(int8); isa {
		return IntToString(v)
	} else if v, isa := any(val).(int16); isa {
		return IntToString(v)
	} else if v, isa := any(val).(int32); isa {
		return IntToString(v)
	} else if v, isa := any(val).(int64); isa {
		return IntToString(v)
	} else if v, isa := any(val).(uint); isa {
		return UintToString(v)
	} else if v, isa := any(val).(uint8); isa {
		return UintToString(v)
	} else if v, isa := any(val).(uint16); isa {
		return UintToString(v)
	} else if v, isa := any(val).(uint32); isa {
		return UintToString(v)
	} else if v, isa := any(val).(uint64); isa {
		return UintToString(v)
	} else if v, isa := any(val).(float32); isa {
		return FloatToString(v)
	} else if v, isa := any(val).(float64); isa {
		return FloatToString(v)
	} else if v, isa := any(val).(*big.Int); isa {
		return BigIntToString(v)
	} else if v, isa := any(val).(*big.Float); isa {
		return BigFloatToString(v)
	}

	// Must be *big.Rat
	return BigRatToNormalizedString(any(val).(*big.Rat))
}

// ==== int/uint to int/uint, float to int, float64 to float32

// NumBits provides the number of bits of any integer or float type
func NumBits[T constraint.Signed | constraint.UnsignedInteger](val T) int {
	return int(goreflect.ValueOf(val).Type().Size() * 8)
}

// IntToInt converts any signed integer type into any signed integer type
// Returns an error if the source value cannot be represented by the target type
func IntToInt[S constraint.SignedInteger, T constraint.SignedInteger](ival S, oval *T) error {
	var (
		srcSize = NumBits(ival)
		tgtSize = NumBits(*oval)
	)

	if (srcSize > tgtSize) && ((ival < S(minIntValue[tgtSize])) || (ival > S(maxIntValue[tgtSize]))) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%d", ival), fmt.Sprintf("%T", *oval))
	}

	*oval = T(ival)
	return nil
}

// MustIntToInt is a Must version of IntToInt
func MustIntToInt[S constraint.SignedInteger, T constraint.SignedInteger](ival S, oval *T) {
	funcs.Must(IntToInt(ival, oval))
}

// IntToUint converts any signed integer type into any unsigned integer type
// Returns an error if the signed int cannot be represented by the unsigned type
func IntToUint[I constraint.SignedInteger, U constraint.UnsignedInteger](ival I, oval *U) error {
	var (
		intSize  = NumBits(ival)
		uintSize = NumBits(*oval)
	)

	if (ival < 0) || ((intSize > uintSize) && (ival > I(maxUintValue[uintSize]))) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%d", ival), fmt.Sprintf("%T", *oval))
	}

	*oval = U(ival)
	return nil
}

// MustIntToUint is a Must version of IntToUint
func MustIntToUint[I constraint.SignedInteger, U constraint.UnsignedInteger](ival I, oval *U) {
	funcs.Must(IntToUint(ival, oval))
}

// UintToInt converts any unsigned integer type into any signed integer type
// Returns an error if the unsigned int cannot be represented by the signed type
func UintToInt[U constraint.UnsignedInteger, I constraint.SignedInteger](ival U, oval *I) error {
	var (
		uintSize = NumBits(ival)
		intSize  = NumBits(*oval)
	)

	if (uintSize >= intSize) && (ival > U(maxIntValue[intSize])) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%d", ival), fmt.Sprintf("%T", *oval))
	}

	*oval = I(ival)
	return nil
}

// MustUintToInt is a Must version of UintToInt
func MustUintToInt[U constraint.UnsignedInteger, I constraint.SignedInteger](ival U, oval *I) {
	funcs.Must(UintToInt(ival, oval))
}

// UintToUint converts any unsigned integer type into any unsigned integer type
// Returns an error if the source value cannot be represented by the target type
func UintToUint[S constraint.UnsignedInteger, T constraint.UnsignedInteger](ival S, oval *T) error {
	var (
		srcSize = NumBits(ival)
		tgtSize = NumBits(*oval)
	)

	if (srcSize > tgtSize) && (ival > S(maxUintValue[tgtSize])) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%d", ival), fmt.Sprintf("%T", *oval))
	}

	*oval = T(ival)
	return nil
}

// MustUintToInt is a Must version of UintToInt
func MustUintToUint[S constraint.UnsignedInteger, T constraint.UnsignedInteger](ival S, oval *T) {
	funcs.Must(UintToUint(ival, oval))
}

// IntToFloat converts any kind of signed or unssigned integer into any kind of float.
// Returns an error if the int value cannot be exactly represented without rounding.
func IntToFloat[I constraint.Integer, F constraint.Float](ival I, oval *F) error {
	// Convert int to float type, which may round if int has more bits than float type mantissa
	inter := F(ival)

	// If converting the float back to the int type is not the same value, rounding occurred
	if ival != I(inter) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%d", ival), goreflect.TypeOf(inter).Name())
	}

	*oval = inter
	return nil
}

// MustIntToFloat is a Must version of IntToFloat
func MustIntToFloat[I constraint.Integer, F constraint.Float](ival I, oval *F) {
	funcs.Must(IntToFloat(ival, oval))
}

// FloatToInt converts and float type to any signed int type
// Returns an error if the float value cannot be represented by the int type
func FloatToInt[F constraint.Float, I constraint.SignedInteger](ival F, oval *I) error {
	var (
		inter1 *big.Rat
		inter2 int64
	)

	if math.IsNaN(float64(ival)) || (FloatToBigRat(ival, &inter1) != nil) || (BigRatToInt64(inter1, &inter2) != nil) || (IntToInt(inter2, oval) != nil) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%g", ival), goreflect.TypeOf(*oval).Name())
	}

	return nil
}

// MustFloatToInt is a Must version of FloatToInt
func MustFloatToInt[F constraint.Float, I constraint.SignedInteger](ival F, oval *I) {
	funcs.Must(FloatToInt(ival, oval))
}

// FloatToUint converts and float type to any unsigned int type
// Returns an error if the float value cannot be represented by the unsigned int type
func FloatToUint[F constraint.Float, I constraint.UnsignedInteger](ival F, oval *I) error {
	var (
		inter1 *big.Rat
		inter2 uint64
	)

	if math.IsNaN(float64(ival)) || (FloatToBigRat(ival, &inter1) != nil) || (BigRatToUint64(inter1, &inter2) != nil) || (UintToUint(inter2, oval) != nil) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%g", ival), goreflect.TypeOf(*oval).Name())
	}

	return nil
}

// MustFloatToUint is a Must version of FloatToUint
func MustFloatToUint[F constraint.Float, I constraint.UnsignedInteger](ival F, oval *I) {
	funcs.Must(FloatToUint(ival, oval))
}

// FloatToFloat converts a float32 or float64 to a float32 or float64
// Returns an error if the float64 is outside the range of a float32
func FloatToFloat[I constraint.Float, O constraint.Float](ival I, oval *O) error {
	ival64 := float64(ival)
	if math.IsInf(ival64, 0) || math.IsNaN(ival64) {
		*oval = O(ival)
		return nil
	}

	if _, isa := any(oval).(*float32); isa && (((ival64 != 0.0) && (math.Abs(ival64) < math.SmallestNonzeroFloat32)) || (math.Abs(ival64) > math.MaxFloat32)) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%g", ival64), "float32")
	}

	*oval = O(ival)
	return nil
}

// MustFloatToFloat is a Must version of FloatToFloat
func MustFloatToFloat[I constraint.Float, O constraint.Float](ival I, oval *O) {
	funcs.Must(FloatToFloat(ival, oval))
}

// ==== ToInt64

// BigIntToInt converts a *big.Int to a signed integer
// Returns an error if the *big.Int cannot be represented as an int64
func BigIntToInt64(ival *big.Int, oval *int64) error {
	if !ival.IsInt64() {
		return fmt.Errorf(errMsg, ival, ival.String(), "int64")
	}

	*oval = ival.Int64()
	return nil
}

// MustBigIntToInt64 is a Must version of BigIntToInt64
func MustBigIntToInt64(ival *big.Int, oval *int64) {
	funcs.Must(BigIntToInt64(ival, oval))
}

// BigFloatToInt64 converts a *big.Float to an int64
// Returns an error if the *big.Float cannot be represented as an int64
func BigFloatToInt64(ival *big.Float, oval *int64) error {
	inter := big.NewInt(0)
	if (BigFloatToBigInt(ival, &inter) != nil) || (BigIntToInt64(inter, oval) != nil) {
		return fmt.Errorf(errMsg, ival, ival.String(), "int64")
	}

	return nil
}

// MustBigFloatToInt64 is a Must version of BigFloatToInt64
func MustBigFloatToInt64(ival *big.Float, oval *int64) {
	funcs.Must(BigFloatToInt64(ival, oval))
}

// BigRatToInt64 converts a *big.Rat to an int64
// Returns an error if the *big.Rat cannot be represented as an int64
func BigRatToInt64(ival *big.Rat, oval *int64) error {
	if (!ival.IsInt()) || (!ival.Num().IsInt64()) {
		return fmt.Errorf(errMsg, ival, ival.String(), "int64")
	}

	*oval = ival.Num().Int64()
	return nil
}

// MustBigRatToInt64 is a Must version of BigRatToInt64
func MustBigRatToInt64(ival *big.Rat, oval *int64) {
	funcs.Must(BigRatToInt64(ival, oval))
}

// StringToInt64 converts a string to an int64
// Returns an error if the string cannot be represented as an int64
func StringToInt64(ival string, oval *int64) error {
	var err error
	*oval, err = strconv.ParseInt(ival, 10, 64)
	if err != nil {
		return fmt.Errorf(errMsg, ival, ival, "int64")
	}

	return nil
}

// MustStringToInt64 is a Must version of StringToInt64
func MustStringToInt64(ival string, oval *int64) {
	funcs.Must(StringToInt64(ival, oval))
}

// ==== ToUint64

// BigIntToUint64 converts a *big.Int to a uint64
// Returns an error if the *big.Int cannot be represented as a uint64
func BigIntToUint64(ival *big.Int, oval *uint64) error {
	if !ival.IsUint64() {
		return fmt.Errorf(errMsg, ival, ival.String(), "uint64")
	}

	*oval = ival.Uint64()
	return nil
}

// MustBigIntToUint64 is a Must version of BigIntToUint64
func MustBigIntToUint64(ival *big.Int, oval *uint64) {
	funcs.Must(BigIntToUint64(ival, oval))
}

// BigFloatToUint64 converts a *big.Float to a uint64
// Returns an error if the *big.Float cannot be represented as a uint64
func BigFloatToUint64(ival *big.Float, oval *uint64) error {
	var inter *big.Int
	if (BigFloatToBigInt(ival, &inter) != nil) || (BigIntToUint64(inter, oval) != nil) {
		return fmt.Errorf(errMsg, ival, ival.String(), "uint64")
	}

	return nil
}

// MustBigFloatToUint64 is a Must version of BigFloatToUint64
func MustBigFloatToUint64(ival *big.Float, oval *uint64) {
	funcs.Must(BigFloatToUint64(ival, oval))
}

// BigRatToUint64 converts a *big.Rat to a uint64
// Returns an error if the *big.Rat cannot be represented as a uint64
func BigRatToUint64(ival *big.Rat, oval *uint64) error {
	if (!ival.IsInt()) || (!ival.Num().IsUint64()) {
		return fmt.Errorf(errMsg, ival, ival.String(), "uint64")
	}

	*oval = ival.Num().Uint64()
	return nil
}

// MustBigRatToUint64 is a Must version of BigRatToUint64
func MustBigRatToUint64(ival *big.Rat, oval *uint64) {
	funcs.Must(BigRatToUint64(ival, oval))
}

// StringToUint64 converts a string to a uint64
// Returns an error if the string cannot be represented as a uint64
func StringToUint64(ival string, oval *uint64) error {
	var err error
	if *oval, err = strconv.ParseUint(ival, 10, 64); err != nil {
		return fmt.Errorf(errMsg, ival, ival, "uint64")
	}

	return nil
}

// MustStringToUint64 is a Must version of StringToUint64
func MustStringToUint64(ival string, oval *uint64) {
	funcs.Must(StringToUint64(ival, oval))
}

// ==== ToFloat32

// BigIntToFloat32 converts a *big.Int to a float32
// Returns an error if the *big.Int cannot be represented as a float32
func BigIntToFloat32(ival *big.Int, oval *float32) error {
	var (
		inter *big.Float
		acc   big.Accuracy
	)
	BigIntToBigFloat(ival, &inter)
	if *oval, acc = inter.Float32(); acc != big.Exact {
		return fmt.Errorf(errMsg, ival, ival.String(), "float32")
	}

	return nil
}

// MustBigIntToFloat32 is a Must version of BigIntToFloat32
func MustBigIntToFloat32(ival *big.Int, oval *float32) {
	funcs.Must(BigIntToFloat32(ival, oval))
}

// BigFloatToFloat32 converts a *big.Float to a float32
// Returns an error if the *big.Float cannot be represented as a float32
func BigFloatToFloat32(ival *big.Float, oval *float32) error {
	var acc big.Accuracy
	if *oval, acc = ival.Float32(); acc != big.Exact {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%.f", ival), "float32")
	}

	return nil
}

// MustBigFloatToFloat32 is a Must version of BigFloatToFloat32
func MustBigFloatToFloat32(ival *big.Float, oval *float32) {
	funcs.Must(BigFloatToFloat32(ival, oval))
}

// BigRatToFloat32 converts a *big.Rat to a float32
// Returns an error if the *big.Rat cannot be represented as a float32
func BigRatToFloat32(ival *big.Rat, oval *float32) error {
	var exact bool
	if *oval, exact = ival.Float32(); !exact {
		return fmt.Errorf(errMsg, ival, ival.String(), "float32")
	}

	return nil
}

// MustBigRatToFloat32 is a Must version of BigRatToFloat32
func MustBigRatToFloat32(ival *big.Rat, oval *float32) {
	funcs.Must(BigRatToFloat32(ival, oval))
}

// StringToFloat32 converts a string to a float32
// Returns an error if the string cannot be represented as a float32
func StringToFloat32(ival string, oval *float32) error {
	var inter *big.Float

	if ival == "NaN" {
		*oval = float32(math.NaN())
		return nil
	}

	if (StringToBigFloat(ival, &inter) != nil) || (BigFloatToFloat32(inter, oval) != nil) {
		return fmt.Errorf(errMsg, ival, ival, "float32")
	}

	return nil
}

// MustStringToFloat32 is a Must version of StringToFloat32
func MustStringToFloat32(ival string, oval *float32) {
	funcs.Must(StringToFloat32(ival, oval))
}

// ==== ToFloat64

// BigIntToFloat64 converts a *big.Int to a float64
// Returns an error if the *big.Int cannot be represented as a float64
func BigIntToFloat64(ival *big.Int, oval *float64) error {
	var (
		inter *big.Float
		acc   big.Accuracy
	)

	BigIntToBigFloat(ival, &inter)
	if *oval, acc = inter.Float64(); acc != big.Exact {
		return fmt.Errorf(errMsg, ival, ival.String(), "float64")
	}

	return nil
}

// MustBigIntToFloat64 is a Must version of BigIntToFloat64
func MustBigIntToFloat64(ival *big.Int, oval *float64) {
	funcs.Must(BigIntToFloat64(ival, oval))
}

// BigFloatToFloat64 converts a *big.Float to a float64
// Returns an error if the *big.Float cannot be represented as a float64
func BigFloatToFloat64(ival *big.Float, oval *float64) error {
	var acc big.Accuracy
	if *oval, acc = ival.Float64(); acc != big.Exact {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%.f", ival), "float64")
	}

	return nil
}

// MustBigFloatToFloat64 is a Must version of BigFloatToFloat64
func MustBigFloatToFloat64(ival *big.Float, oval *float64) {
	funcs.Must(BigFloatToFloat64(ival, oval))
}

// BigRatToFloat64 converts a *big.Rat to a float64
// Returns an error if the *big.Rat cannot be represented as a float64
func BigRatToFloat64(ival *big.Rat, oval *float64) error {
	var exact bool
	if *oval, exact = ival.Float64(); !exact {
		return fmt.Errorf(errMsg, ival, ival.String(), "float64")
	}

	return nil
}

// MustBigRatToFloat64 is a Must version of BigRatToFloat64
func MustBigRatToFloat64(ival *big.Rat, oval *float64) {
	funcs.Must(BigRatToFloat64(ival, oval))
}

// StringToFloat64 converts a string to a float64
// Returns an error if the string cannot be represented as a float64
func StringToFloat64(ival string, oval *float64) error {
	var inter *big.Float

	if ival == "NaN" {
		*oval = math.NaN()
		return nil
	}

	if (StringToBigFloat(ival, &inter) != nil) || (BigFloatToFloat64(inter, oval) != nil) {
		return fmt.Errorf(errMsg, ival, ival, "float64")
	}

	return nil
}

// MustStringToFloat64 is a Must version of StringToFloat64
func MustStringToFloat64(ival string, oval *float64) {
	funcs.Must(StringToFloat64(ival, oval))
}

// ==== ToBigInt

// IntToBigInt converts any signed int type into a *big.Int
func IntToBigInt[T constraint.SignedInteger](ival T, oval **big.Int) {
	*oval = big.NewInt(int64(ival))
}

// UintToBigInt converts any unsigned int type into a *big.Int
func UintToBigInt[T constraint.UnsignedInteger](ival T, oval **big.Int) {
	*oval = big.NewInt(0)
	(*oval).SetUint64(uint64(ival))
}

// FloatToBigInt converts any float type to a *big.Int
// Returns an error if the float has fractional digits
func FloatToBigInt[T constraint.Float](ival T, oval **big.Int) error {
	if math.IsInf(float64(ival), 0) || math.IsNaN(float64(ival)) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%g", ival), "*big.Int")
	}

	var inter big.Rat
	inter.SetFloat64(float64(ival))
	if !inter.IsInt() {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%g", ival), "*big.Int")
	}

	*oval = inter.Num()
	return nil
}

// MustFloatToBigInt is a Must version of FloatToBigInt
func MustFloatToBigInt[T constraint.Float](ival T, oval **big.Int) {
	funcs.Must(FloatToBigInt(ival, oval))
}

// BigIntToBigInt makes a copy of a *big.Int such that ival and *oval are different pointers
func BigIntToBigInt(ival *big.Int, oval **big.Int) {
	*oval = big.NewInt(0)
	(*oval).Set(ival)
}

// BigFloatToBigInt converts a *big.Float to a *big.Int.
// Returns an error if the *big.Float has any fractional digits.
func BigFloatToBigInt(ival *big.Float, oval **big.Int) error {
	inter, acc := ival.Rat(nil)
	if (inter == nil) || (!inter.IsInt()) || (acc != big.Exact) {
		return fmt.Errorf(errMsg, ival, ival.String(), "*big.Int")
	}

	*oval = inter.Num()
	return nil
}

// MustBigFloatToBigInt is a Must version of BigFloatToBigInt
func MustBigFloatToBigInt(ival *big.Float, oval **big.Int) {
	funcs.Must(BigFloatToBigInt(ival, oval))
}

// BigRatToBigInt converts a *big.Rat to a *big.Int
// Returns an error if the *big.Rat is not an int
func BigRatToBigInt(ival *big.Rat, oval **big.Int) error {
	if !ival.IsInt() {
		return fmt.Errorf(errMsg, ival, ival.String(), "*big.Int")
	}

	*oval = ival.Num()
	return nil
}

// MustBigRatToBigInt is a Must version of BigFloatToBigInt
func MustBigRatToBigInt(ival *big.Rat, oval **big.Int) {
	funcs.Must(BigRatToBigInt(ival, oval))
}

// StringtoBigInt converts a string to a *big.Int.
// Returns an error if the string is not an integer.
func StringToBigInt(ival string, oval **big.Int) error {
	*oval = big.NewInt(0)
	if _, ok := (*oval).SetString(ival, 10); !ok {
		return fmt.Errorf(errMsg, ival, ival, "*big.Int")
	}

	return nil
}

// MustStringToBigInt is a Must version of StringToBigInt
func MustStringToBigInt(ival string, oval **big.Int) {
	funcs.Must(StringToBigInt(ival, oval))
}

// ==== ToBigFloat

// IntToBigFloat converts any signed int type into a *big.Float
func IntToBigFloat[T constraint.SignedInteger](ival T, oval **big.Float) {
	prec := uint(math.Ceil(float64(len(IntToString(ival))) * log2Of10))
	*oval = big.NewFloat(0)
	(*oval).SetPrec(prec)
	(*oval).SetInt64(int64(ival))
}

// UintToBigFloat converts any unsigned int type into a *big.Float
func UintToBigFloat[T constraint.UnsignedInteger](ival T, oval **big.Float) {
	*oval = big.NewFloat(0).SetUint64(uint64(ival))
}

// FloatToBigFloat converts any float type into a *big.Float
func FloatToBigFloat[T constraint.Float](ival T, oval **big.Float) error {
	if math.IsNaN(float64(ival)) {
		return fmt.Errorf(errMsg, ival, "NaN", "*big.Float")
	}

	*oval = big.NewFloat(float64(ival))
	return nil
}

// MustFloatToBigFloat is a Must version of FloatToBigFloat
func MustFloatToBigFloat[T constraint.Float](ival T, oval **big.Float) {
	funcs.Must(FloatToBigFloat(ival, oval))
}

// BigIntToBigFloat converts a *big.Int into a *big.Float
func BigIntToBigFloat(ival *big.Int, oval **big.Float) {
	StringToBigFloat(ival.String(), oval)
}

// BigFloatToBigFloat makes a copy of a *big.Float such that ival and *oval are different pointers
func BigFloatToBigFloat(ival *big.Float, oval **big.Float) {
	*oval = big.NewFloat(0)
	(*oval).SetMode(ival.Mode())
	(*oval).SetPrec(ival.Prec())
	(*oval).Set(ival)
}

// BigRatToBigFloat converts a *big.Rat to a *big.Float
func BigRatToBigFloat(ival *big.Rat, oval **big.Float) {
	// Use numerator to calculate the precision, shd be accurate since denominator is basically the exponent
	prec := int(math.Ceil(math.Max(float64(53), float64(len(ival.Num().String()))*log2Of10)))
	*oval, _, _ = big.ParseFloat(ival.FloatString(prec), 10, uint(prec), big.ToNearestEven)

	// Set accuracy to exact
	(*oval).SetMode((*oval).Mode())
}

// StringToBigFloat converts a string to a *big.Float
// Returns an error if the string is not a valid float string
func StringToBigFloat(ival string, oval **big.Float) error {
	// A *big.Float is imprecise, but you can set the precision
	// The crude measure we use is the largest of 53 (number of bits in a float64) and ceiling(string length * Log2(10))
	// If every char was a significant digit, the ceiling calculation would be the minimum number of bits required
	var (
		numBits = uint(math.Max(53, math.Ceil(float64(len(ival))*log2Of10)))
		err     error
	)

	if *oval, _, err = big.ParseFloat(ival, 10, numBits, big.ToNearestEven); err != nil {
		return fmt.Errorf(errMsg, ival, ival, "*big.Float")
	}

	return nil
}

// MustStringToBigFloat is a Must version of FloatToBigFloat
func MustStringToBigFloat(ival string, oval **big.Float) {
	funcs.Must(StringToBigFloat(ival, oval))
}

// ==== ToBigRat

// IntToBigRat converts any signed int type into a *big.Rat
func IntToBigRat[T constraint.SignedInteger](ival T, oval **big.Rat) {
	*oval = big.NewRat(1, 1).SetInt64(int64(ival))
}

// UintToBigRat converts any unsigned int type into a *big.Rat
func UintToBigRat[T constraint.UnsignedInteger](ival T, oval **big.Rat) {
	*oval = big.NewRat(1, 1).SetUint64(uint64(ival))
}

// FloatToBigRat converts any float type into a *big.Rat
func FloatToBigRat[T constraint.Float](ival T, oval **big.Rat) error {
	if math.IsInf(float64(ival), 0) || math.IsNaN(float64(ival)) {
		return fmt.Errorf(errMsg, ival, fmt.Sprintf("%g", ival), "*big.Rat")
	}

	var inter *big.Float
	FloatToBigFloat(ival, &inter)
	BigFloatToBigRat(inter, oval)
	return nil
}

// MustFloatToBigRat is a Must version of FloatToBigRat
func MustFloatToBigRat[T constraint.Float](ival T, oval **big.Rat) {
	funcs.Must(FloatToBigRat(ival, oval))
}

// BigIntToBigRat converts a *big.Int into a *big.Rat
func BigIntToBigRat(ival *big.Int, oval **big.Rat) {
	*oval = big.NewRat(1, 1).SetFrac(ival, big.NewInt(1))
}

// BigFloatToBigRat converts a *big.Float into a *big.Rat
func BigFloatToBigRat(ival *big.Float, oval **big.Rat) error {
	if ival.IsInf() {
		return fmt.Errorf(errMsg, ival, ival.String(), "*big.Rat")
	}

	*oval, _ = big.NewRat(1, 1).SetString(BigFloatToString(ival))
	return nil
}

// MustBigFloatToBigRat is a Must version of FloatToBigRat
func MustBigFloatToBigRat(ival *big.Float, oval **big.Rat) {
	funcs.Must(BigFloatToBigRat(ival, oval))
}

// BigRatToBigRat makes a copy of a *big.Rat such that ival and *oval are different pointers
func BigRatToBigRat(ival *big.Rat, oval **big.Rat) {
	*oval = big.NewRat(0, 1)
	(*oval).Set(ival)
}

// StringToBigRat converts a string into a *big.Rat
func StringToBigRat(ival string, oval **big.Rat) error {
	var ok bool
	if *oval, ok = big.NewRat(1, 1).SetString(ival); !ok {
		return fmt.Errorf(errMsg, ival, ival, "*big.Rat")
	}

	return nil
}

// MustStringToBigRat is a Must version of StringToBigRat
func MustStringToBigRat(ival string, oval **big.Rat) {
	funcs.Must(StringToBigRat(ival, oval))
}

// FloatStringToBigRat converts a float string to a *big.Rat.
// Unlike StringToBigRat, it will not accept a ratio string like 5/4.
func FloatStringToBigRat(ival string, oval **big.Rat) error {
	// ensure the string is a float string, and not a ratio
	var err error
	if (ival == "+Inf") || (ival == "-Inf") || (ival == "NaN") {
		return fmt.Errorf("The float string value of %s cannot be converted to *big.Rat", ival)
	}

	if _, _, err = big.NewFloat(0).Parse(ival, 10); err != nil {
		return fmt.Errorf("The float string value of %s cannot be converted to *big.Rat", ival)
	}

	// If it is a float string, cannot fail to be parsed by StringToBigRat
	StringToBigRat(ival, oval)
	return nil
}

// MustFloatStringToBigRat is a Must version of FloatStringToBigRat
func MustFloatStringToBigRat(ival string, oval **big.Rat) {
	funcs.Must(FloatStringToBigRat(ival, oval))
}

// To converts any numeric or string into any other such type.
// The actual conversion is performed by other funcs.
func To[S constraint.Numeric | ~string, T constraint.Numeric | ~string](src S, tgt *T) error {
	var (
		valsrc = goreflect.ValueOf(src)
		valtgt = goreflect.ValueOf(tgt)
	)

	// Convert source and target to base types
	reflect.ToBaseType(&valsrc)
	reflect.ToBaseType(&valtgt)

	// No conversion function exists if src and *tgt are the same type, unless they are *big types
	copy := valsrc.Type() == valtgt.Elem().Type()
	if copy {
		if _, isa := any(src).(*big.Int); isa {
			copy = false
		}

		if _, isa := any(src).(*big.Float); isa {
			copy = false
		}

		if _, isa := any(src).(*big.Rat); isa {
			copy = false
		}
	}

	if copy {
		valtgt.Elem().Set(valsrc)
		return nil
	}

	// Types differ, lookup conversion using base types and execute it, returning result
	return convertFromTo[valsrc.Type().String()+valtgt.Type().Elem().String()](valsrc.Interface(), valtgt.Interface())
}

// MustTo is a Must version of To
func MustTo[S constraint.Numeric | ~string, T constraint.Numeric | ~string](src S, tgt *T) {
	funcs.Must(To(src, tgt))
}

// ToBigOps is the BigOps version of To
func ToBigOps[S constraint.Numeric | ~string, T constraint.BigOps[T]](src S, tgt *T) error {
	var (
		valsrc = goreflect.ValueOf(src)
		valtgt = goreflect.ValueOf(tgt)
	)

	// Convert source to base type
	reflect.ToBaseType(&valsrc)

	// No conversion function exists if src and *tgt are the same type
	if valsrc.Type() == valtgt.Elem().Type() {
		valtgt.Elem().Set(valsrc)
		return nil
	}

	// Types differ, lookup conversion using base types and execute it, returning result
	return convertFromTo[valsrc.Type().String()+valtgt.Type().Elem().String()](valsrc.Interface(), valtgt.Interface())
}

// MustToBigOps is a Must version of ToBigOps
func MustToBigOps[S constraint.Numeric | ~string, T constraint.BigOps[T]](src S, tgt *T) {
	funcs.Must(ToBigOps(src, tgt))
}
