package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"
	goreflect "reflect"

	"github.com/bantling/micro/constraint"
	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/reflect"
)

var (
	errLookupMsg                   = "%v cannot be converted to %v"
	errCopyNilSourceMsg            = "A nil %s cannot be copied to a(n) %s"
	errConvertNilSourceMsg         = "A nil %s cannot be converted to a(n) %s"
	errCopyNilTargetMsg            = "A(n) %s cannot be copied to a nil %s"
	errEmptyMaybeMsg               = "An empty %s cannot be converted to a(n) %s"
	errMsg                         = "The %T value of %s cannot be converted to %s"
	errRegisterMultiplePointersMsg = "The %s type %s has too many pointers"
	errRegisterExistsMsg           = "The conversion from %s to %s has already been registered"

	log2Of10 = math.Log2(10)

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

	badConversionKinds = map[goreflect.Kind]bool{
		goreflect.Uintptr:       true,
		goreflect.Chan:          true,
		goreflect.Func:          true,
		goreflect.UnsafePointer: true,
	}
)

// ==== To and functions that support it

// LookupConversion looks for a conversion from a source type to a target type.
//
// It is an error if either type is more than one pointer, or a uintptr, chan, func, or unsafe pointer.
// If the source and target types are the same, a conversion function that just copies the source to the target is returned.
// The source type may be of the following forms, where T represents any accepted value type:
// - T
// - *T
// - Maybe[T]
// - Maybe[*T]
//
// The target types are the same as the source types.
// For any of the above forms, T may be a primitive sub type (eg, type subint int).
// In such cases, up to four lookups are performed (first match is used):
// - A conversion using the src sub  type and tgt sub  type
// - A conversion using the src base type and tgt sub  type
// - A conversion using the src sub  type and tgt base type
// - A conversion using the src base type and tgt base type
//
// Example lookups, listing possible conversions in search order (shown with the extra * the target type has to have):
// - int to string -> int to string
// - subint to *string -> subint to string, int to string
// - Maybe[int] to string -> int to string
// - int to int -> copy
// - subint to Maybe[*int] -> subint to *int, subint to int, copy
//
// This function returns func, error:
// If a conversion is found (or it is a copy) : returns func, nil
// If a conversion is not found               : returns nil,  nil
// Conversion is not allowed                  : returns nil,  err
func LookupConversion(src, tgt goreflect.Type) (func(any, any) error, error) {
  // Verify src and tgt are not nil
  if (src == nil) || (tgt == nil) {
    return nil, fmt.Errorf(errLookupMsg, src, tgt)
  }

	// Verify the types are not more than one pointer
	if (reflect.NumPointers(src) > 1) || (reflect.NumPointers(tgt) > 1) {
		// A conversion CANNOT be registered for multiple pointers
		return nil, fmt.Errorf(errLookupMsg, src, tgt)
	}

	// Verify the types are not uintptr, chan, func, or unsafe pointer
	for _, check := range []goreflect.Type{src, tgt} {
		if badConversionKinds[check.Kind()] {
			// A conversion CANNOT be registered for these kinds
			return nil, fmt.Errorf(errLookupMsg, src, tgt)
		}
	}

  // Check for conversion from src to tgt as is, most common case
  if conv, haveIt := convertFromTo[src.String()+tgt.String()]; haveIt {
    return conv, nil
  }

  // Search every valid combination of src, tgt = val/subtype, *val/subtype, maybe val/subtype, maybe *val/subtype.
  // Ensure we do not allow invalid combinations, such as *Maybe.
  // If a conversion is found:
  // - create a wrapper func that deals with conversions, *, maybe for both src and tgt
  // - register it so future calls don't have to do same search
  // - return it to caller
  // If no converison found, return nil, error

  var (
    srcBase, srcPtr, srcPtrBase, srcMaybe, srcMaybeBase, srcMaybePtr, srcMaybePtrBase goreflect.Type
    tgtBase, tgtPtr, tgtPtrBase, tgtMaybe, tgtMaybeBase, tgtMaybePtr, tgtMaybePtrBase goreflect.Type
    convFn func(any, any) error
    haveIt bool
  )

  srcBase = reflect.TypeToBaseType(src)

  if srcPtr = funcs.TernaryResult(src.Kind() == goreflect.Pointer, src.Elem, nil); srcPtr != nil {
    srcPtrBase = reflect.TypeToBaseType(srcPtr)

    if reflect.GetMaybeType(srcPtr) != nil {
      // Cannot have a *Maybe, that makes no sense
      return nil, fmt.Errorf(errLookupMsg, src, tgt)
    }
  }

  if srcMaybe = reflect.GetMaybeType(src); srcMaybe != nil {
    srcMaybeBase = reflect.TypeToBaseType(srcMaybe)

    if srcMaybePtr = funcs.TernaryResult(srcMaybe.Kind() == goreflect.Pointer, srcMaybe.Elem, nil); srcMaybePtr != nil {
      srcMaybePtrBase = reflect.TypeToBaseType(srcMaybePtr)
    }
  }

  tgtBase = reflect.TypeToBaseType(tgt)

  if tgtPtr = funcs.TernaryResult(tgt.Kind() == goreflect.Pointer, tgt.Elem, nil); tgtPtr != nil {
    tgtPtrBase = reflect.TypeToBaseType(tgtPtr)

    if reflect.GetMaybeType(tgtPtr) != nil {
      // Cannot have a *Maybe, that makes no sense
      return nil, fmt.Errorf(errLookupMsg, src, tgt)
    }
  }

  if tgtMaybe = reflect.GetMaybeType(tgt); tgtMaybe != nil {
    tgtMaybeBase = reflect.TypeToBaseType(tgtMaybe)

    if tgtMaybePtr = funcs.TernaryResult(tgtMaybe.Kind() == goreflect.Pointer, tgtMaybe.Elem, nil); tgtMaybePtr != nil {
      tgtMaybePtrBase = reflect.TypeToBaseType(tgtMaybePtr)
    }
  }

  for _, srcTyp := range []goreflect.Type{
    src, srcBase, srcPtr, srcPtrBase, srcMaybe, srcMaybeBase, srcMaybePtr, srcMaybePtrBase,
  } {
    for _, tgtTyp := range []goreflect.Type{
      tgt, tgtBase, tgtPtr, tgtPtrBase, tgtMaybe, tgtMaybeBase, tgtMaybePtr, tgtMaybePtrBase,
    } {
      // Cannot lookup conversions for types that don't exist
      if (srcTyp != nil) && (tgtTyp != nil) {
        convFn, haveIt = nil, srcTyp == tgtTyp
        if (!haveIt) {
          convFn, haveIt = convertFromTo[srcTyp.String()+tgtTyp.String()]
        }

        if haveIt {
          // Generate a function to unwrap the src type and read it
          var srcFn func(goreflect.Value) goreflect.Value
          switch srcTyp {
          case src:
            srcFn = func(s goreflect.Value) goreflect.Value { return s }
          case srcBase:
            srcFn = func(s goreflect.Value) goreflect.Value {return s.Convert(srcBase) }
          case srcPtr:
            srcFn = func(s goreflect.Value) goreflect.Value { return s.Elem() }
          case srcPtrBase:
            srcFn = func(s goreflect.Value) goreflect.Value { return s.Elem().Convert(srcPtrBase) }
          case srcMaybe:
            srcFn = func(s goreflect.Value) goreflect.Value { return reflect.GetMaybeValue(s) }
          case srcMaybeBase:
            srcFn = func(s goreflect.Value) goreflect.Value {
              fmt.Printf("a\n")
              if temp := reflect.GetMaybeValue(s); temp.IsValid() {
                fmt.Printf("b\n")
                return temp.Convert(srcMaybeBase)
              } else {
                fmt.Printf("c\n")
                return temp
              }
           }
          case srcMaybePtr:
            srcFn = func(s goreflect.Value) goreflect.Value {
              if temp := reflect.GetMaybeValue(s); temp.IsValid() {
                return temp.Elem()
              } else {
                return temp
              }
            }
          case srcMaybePtrBase:
            srcFn = func(s goreflect.Value) goreflect.Value {
              if temp := reflect.GetMaybeValue(s); temp.IsValid() {
                return temp.Elem().Convert(srcMaybePtrBase)
              } else {
                return temp
              }
            }
          }

          // Generate a function to wrap the tgt type
          var tgtFn func(temp, t goreflect.Value)
          switch tgtTyp {
          case tgt:
            tgtFn = func(temp, t goreflect.Value) { t.Elem().Set(temp.Elem()) }
          case tgtBase:
            tgtFn = func(temp, t goreflect.Value) { t.Elem().Set(temp.Elem().Convert(tgt)) }
          case tgtPtr:
            tgtFn = func(temp, t goreflect.Value) { t.Elem().Set(temp.Elem()) }
          case tgtPtrBase:
            tgtFn = func(temp, t goreflect.Value) { t.Elem().Set(temp.Elem().Convert(tgt.Elem())) }
          case tgtMaybe:
            tgtFn = func(temp, t goreflect.Value) { reflect.SetMaybeValue(t, temp.Elem()) }
          case tgtMaybeBase:
            tgtFn = func(temp, t goreflect.Value) { reflect.SetMaybeValue(t, temp.Elem().Convert(tgtMaybe)) }
          case tgtMaybePtr:
            tgtFn = func(temp, t goreflect.Value) { reflect.SetMaybeValue(t, temp) }
          case tgtMaybePtrBase:
            tgtFn = func(temp, t goreflect.Value) { reflect.SetMaybeValue(t, temp.Convert(goreflect.PtrTo(tgtMaybePtr))) }
          }

          // If convFn is nil and the types are the same, generate a copy function
          if (convFn == nil) && (srcTyp == tgtTyp) {
            convFn = func(s, t any) error {
              goreflect.ValueOf(t).Elem().Set(goreflect.ValueOf(s))
              return nil
            }
          }

          // Return a conversion function that unwraps the source as needed, and wraps the target value as needed
          return func(s, t any) error {
            // Use reflection to do runtime type assertion exactly like a provided conversion function would
            // This ensure two things:
            // - A copy conversion does not inadvertently allow copying any random types that happen to be the same
            // - Registered conversions have errors that indicate if the source or target type is the problem
            srcVal, tgtVal := goreflect.ValueOf(s), goreflect.ValueOf(t)
            reflect.MustTypeAssert(srcVal, src, "source")
            reflect.MustTypeAssert(tgtVal, goreflect.PtrTo(tgt), "target")

            // Unwrap src value, which will be invalid for nil ptr or empty maybe
            srcVal = srcFn(srcVal)
            fmt.Printf("0. %t\n", srcVal.IsValid())

            if !srcVal.IsValid() {
              fmt.Printf("1. %s\n", tgtVal.Type())
              // Tgt must be nillable or maybe
              if reflect.IsNillable(tgt) {
                // Tgt is nillable
                fmt.Printf("2.\n")
                tgtVal.Elem().SetZero()
              } else if tgtMaybe != nil {
                // Tgt is a Maybe
                reflect.SetMaybeValueEmpty(tgtVal)
              } else if srcMaybe == nil {
                // Tgt cannot be nil, src is a nil Ptr
                return fmt.Errorf(errConvertNilSourceMsg, src, tgt)
              } else {
                // Tgt cannot be nil, src is an empty Maybe
                return fmt.Errorf(errEmptyMaybeMsg, src, tgt)
              }

              return nil
            }

            // Create a pointer to the target unwrapped type for the conversion to write to
            temp := goreflect.New(tgtTyp)
            fmt.Printf("2. %s\n", temp.Type())

            // Convert source -> unwrapped target
            err := convFn(funcs.TernaryResult(srcVal.IsValid(), srcVal.Interface, nil), temp.Interface())
            fmt.Printf("3. %s\n", err)

            // Wrap target value only if no error occurred - the target is unmodified if the conversion fails
            if err == nil {
              fmt.Printf("4. %s\n", err)
              tgtFn(temp, tgtVal)
            }

            // Return any error
            return err
          }, nil
        }
      }
    }
  }

  return convFn, nil

  // switch {
  // // Is the src convertible to tgt?
  // if src.CanConvert(tgt) {
  //   if convFn, haveIt = convertFromTo[srcBase.String()+tgt.String()]; haveIt {
  //     return convFn
  //   }
  //   if convFn, haveIt = convertFromTo[src.String()+tgtBase.String()]; haveIt {
  //     return convFn
  //   }
  //   if convFn, haveIt = convertFromTo[srcBase.String()+tgtBase.String()]; haveIt {
  //     return convFn
  //   }
  //
  //   return func(s, t any) error {
  //     // Copy src base -> tgt
  //     goreflect.ValueOf(t).Elem().Set(goreflect.ValueOf(s).Convert(tgt))
  //   }, nil
  // }
  //
  // // Is the source a pointer?
  // if srcIsPtr {
  //   // Check for conversion from *src to tgt
  //   if convFn, haveIt := convertFromTo[src.Elem().String()+tgt.String()]; haveIt {
  //     return func(s, t any) error {
  //       return convFn(goreflect.ValueOf(s).Elem().Interface(), t)
  //     }, nil
  //   }
  //
  //   // Are *src and tgt same types?
  //   if src.Elem() == tgt {
  //     // Copy *src -> tgt
  //     return func(s, t any) error {
  //       goreflect.ValueOf(t).Elem().Set(goreflect.ValueOf(s).Elem())
  //       return nil
  //     }
  //   }
  //
  //   // Is the target a pointer?
  //   if tgtIsPtr {
  //     // Check for conversion from *src to *tgt
  //     if convFn, haveIt = convertFromTo[src.Elem().String()+tgt.Elem().String()]; haveIt {
  //       return func(s, t any) error {
  //         return convFn(goreflect.ValueOf(s).Elem().Interface(), goreflect.ValueOf(t).Elem().Interface())
  //       }, nil
  //     }
  //   }
  //
  //   // Is the target a maybe?
  //   if tgtMaybe != nil {
  //     // Check for conversion from src to Maybe[tgt]
  //     if convFn, haveIt = convertFromTo[src.String()+tgtMaybe.String()]; haveIt {
  //       return func(s, t any) error {
  //         tval := goreflect.New(tgtMaybe)
  //
  //         err := convFn(s, tval.Interface())
  //         if err != nil {
  //           reflect.SetMaybeValue(goreflect.ValueOf(t).Elem(), tval.Elem())
  //         }
  //
  //         return err
  //       }, nil
  //     }
  //   }
  // }

  // // Is the target a pointer?
  // tgtIsPtr := tgt.Kind() == goreflect.Pointer
  // if tgtIsPtr {
  //   // Check for conversion from src to *tgt
  //   if convFn, haveIt := convertFromTo[src.String()+tgt.Elem().String()]; haveIt {
  //     return func(s, t any) error {
  //       return convFn(s, goreflect.ValueOf(t).Elem().Interface())
  //     }, nil
  //   }
  //
  //   // Is the source a Maybe?
  //   if srcMaybe != nil {
  //     // Check for conversion from Maybe[src] to *tgt
  //     if convFn, haveIt := convertFromTo[srcMaybe.String()+tgt.Elem().String()]; haveIt {
  //       return func(s, t any) error {
  //         sval := reflect.GetMaybeValue(goreflect.ValueOf(s))
  //         // If Maybe[src] empty?
  //         if !sval.IsValid() {
  //           // Set target pointer to nil
  //           goreflect.ValueOf(t).Elem().SetZero()
  //           return nil
  //         }
  //
  //         return convFn(sval.Interface(), goreflect.ValueOf(t).Elem().Interface())
  //       }, nil
  //     }
  //   }
  // }
  //
  // // Is the source a Maybe?
  // if srcMaybe {
  //   // Check for conversion from Maybe[src] to tgt
  //   if convFn, haveIt := convertFromTo[srcMaybe.String()+tgt.String()]; haveIt {
  //     return func(s, t any) error {
  //       maybeVal := reflect.GetMaybeValue(goreflect.ValueOf(s))
  //
  //       if !maybeVal.IsValid() {
  //         // If tgt is a pointer, set it to nil
  //         if tgtIsPtr {
  //           goreflect.ValueOf(t).Elem().SetZero()
  //           return nil
  //
  //         // If tgt is a Maybe, set it to empty
  //         } else if tgtMaybe != nil {
  //           reflect.SetMaybeValueEmpty(goreflect.ValueOf(t))
  //           return nil
  //
  //         // If tgt is not a pointer or maybe, it is an error to convert an empty Maybe to target
  //         } else {
  //           return fmt.Errorf(errEmptyMaybeMsg, src, tgt)
  //         }
  //       }
  //
  //       // Convert present Maybe[S] -> T
  //       return convFn(maybeVal.Interface(), t)
  //     }, nil
  //   }
  // }
  //
  // if tgtMaybe {
  //   // Check for conversion from src to Maybe[tgt]
  //   if convFn, haveIt := convertFromTo[src.String()+tgtMaybe.String()]; haveIt {
  //     return func(s, t any) error {
  //       // If src is a pointer
  //     }, nil
  //   }
  // }

	// var (
	// 	// // Max one ptrs
  //   srcStr = src.String()
  //   tgtStr = tgt.String()
  //
	// 	// Are they ptrs?
	// 	srcIsPtr = src.Kind() == goreflect.Pointer
  //   srcDeref = funcs.TernaryResult(srcIsPtr, src.Elem, nil)
  //   srcDerefStr = funcs.TernaryResult(srcIsPtr, srcDeref.String, nil)
  //
	// 	tgtIsPtr = tgt.Kind() == goreflect.Pointer
  //   tgtDeref = funcs.TernaryResult(tgtIsPtr, src.Elem, nil)
  //   tgtDerefStr = funcs.TernaryResult(tgtIsPtr, tgtDeref.String, nil)
  //
  //   // Are they Maybes?
	// 	srcMaybe = reflect.GetMaybeType(src)
  //   srcMaybeStr = funcs.TernaryResult(srcMaybe != nil, srcMaybe.String, nil)
  //
	// 	tgtMaybe = reflect.GetMaybeType(tgt)
  //   tgtMaybeStr = funcs.TernaryResult(tgtMaybe != nil, tgtMaybe.String, nil)
  //
  //   // Are they subtypes?
	// 	srcBase = reflect.TypeToBaseType(src)
  //   srcBaseStr = srcBase.String()
  //
	// 	tgtBase = reflect.TypeToBaseType(tgt)
  //   tgtBaseStr = tgtBase.String()
  //
  //   // Error message
  //   errMsg = fmt.Errorf(errLookupMsg, src, tgt)
  //
  //   // Scenarios to handle, described by 3 bits for src, and 3 bits for tgt.
  //   // Bits are PMS, where 1 values mean P = pointer, M = Maybe, S = subtype.
  //   // 6 bits have 2^6 = 64 combinations.
  //   // The 4 patterns * (is/not a subtype) = 8 combinations that cover all 3 bit patterns.
  //   // - T         = 000/001
  //   // - *T        = 100/101
  //   // - Maybe[T]  = 010/011
  //   // - Maybe[*T] = 110/111
  //   // 8 combinations of src * 8 combinations of tgt = 64 combinations that cover all 6 bit patterns.
  //   srcScenario = funcs.Ternary(srcIsPtr, 0b100, 0) | funcs.Ternary(srcMaybe != nil, 0b010, 0) | funcs.Ternary(src != srcBase, 0b001, 0)
  //   tgtScenario = funcs.Ternary(tgtIsPtr, 0b100, 0) | funcs.Ternary(tgtMaybe != nil, 0b010, 0) | funcs.Ternary(tgt != tgtBase, 0b001, 0)
	// )
  //
  // // Notation used for each combination listed below:
  // // (V|*|M|M*)(B|U), where:
  // //
  // // V  = value
  // // *  = pointer
  // // M  = Maybe
  // // M* = Maybe[pointer]
  // //
  // // B  = base type
  // // U  = sub  type
  // switch (srcScenario << 3) | tgtScenario {
  // case 0b000_000: { // VB -> VB
  //
  //     // Try S = T
  //     if src == tgt {
  //       // Copy S -> *T
  //       return func(s, t any) error {
  //         goreflect.ValueOf(t).Elem().Set(goreflect.ValueOf(s))
  //         return nil
  //       }, nil
  //     }
  //   }
  //
  // case 0b001_000: { // VU -> VB
  //     // Try S -> T
  //     if conv, haveIt := convertFromTo[srcStr+tgtStr]; haveIt {
  //       return conv, nil
  //     }
  //
  //     // Try SB -> T
  //     if conv, haveIt := convertFromTo[srcBaseStr+tgtStr]; haveIt {
  //       return func(s, t any) error {
  //         return conv(goreflect.ValueOf(s).Convert(srcBase).Interface(), t)
  //       }, nil
  //     }
  //
  //     // Try SB = TB
  //     if srcBase == tgt {
  //       // Copy SB -> *T
  //       return func(s, t any) error {
  //         goreflect.ValueOf(t).Elem().Set(goreflect.ValueOf(s).Convert(tgt))
  //         return nil
  //       }, nil
  //     }
  //   }
  //
  // case 0b010_000: { // MB -> VB
  //     // Try S -> T
  //     if conv, haveIt := convertFromTo[srcStr+tgtStr]; haveIt {
  //       return conv, nil
  //     }
  //
  //     // Try SM -> T
  //     if conv, haveIt := convertFromTo[srcMay+tgtStr]; haveIt {
  //       return func(s, t any) error {
  //         return conv(goreflect.ValueOf(s).Elem().Interface(), t)
  //       }, nil
  //     }
  //
  //     if srcDeref == tgt {
  //       // Copy *S -> *T
  //       return func(s, t any) error {
  //         goreflect.ValueOf(t).Elem().Set(goreflect.ValueOf(s).Elem())
  //         return nil
  //       }, nil
  //     }
  //   }
  //
  // case 0b100_000: { // PB -> VB
  //     if conv, haveIt := convertFromTo[srcStr+tgtStr]; haveIt {
  //       return conv, nil
  //     }
  //
  //     if conv, haveIt := convertFromTo[srcDerefStr+tgtStr]; haveIt {
  //       return func(s, t any) error {
  //         return conv(goreflect.ValueOf(s).Elem().Interface(), t)
  //       }, nil
  //     }
  //
  //     if srcDeref == tgt {
  //       // Copy *S -> *T
  //       return func(s, t any) error {
  //         goreflect.ValueOf(t).Elem().Set(goreflect.ValueOf(s).Elem())
  //         return nil
  //       }, nil
  //     }
  //   }
  // }
  //
  // return nil, errMsg
}

// func LookupConversion(src, tgt goreflect.Type) (func(any, any) error, error) {
// 	// Verify the types are not more than one pointer
// 	if (reflect.NumPointers(src) > 1) || (reflect.NumPointers(tgt) > 1) {
// 		// A conversion CANNOT be registered for multiple pointers
// 		return nil, fmt.Errorf(errLookupMsg, src, tgt)
// 	}
//
// 	// Verify the types are not uintptr, chan, func, or unsafe pointer
// 	for _, check := range []goreflect.Type{src, tgt} {
// 		if badConversionKinds[check.Kind()] {
// 			// A conversion CANNOT be registered for these kinds
// 			return nil, fmt.Errorf(errLookupMsg, src, tgt)
// 		}
// 	}
//
// 	// Get max one ptr derefd types, are they pointers, are they maybes, base types
// 	var (
// 		// Max one ptrs
// 		maxOnePtrSrc = reflect.DerefTypeMaxOnePtr(src)
// 		maxOnePtrTgt = reflect.DerefTypeMaxOnePtr(tgt)
//
// 		// Are the max ones ptrs?
// 		srcIsPtr = maxOnePtrSrc.Kind() == goreflect.Pointer
// 		tgtIsPtr = maxOnePtrTgt.Kind() == goreflect.Pointer
// 	)
//
// 	// Vsrs that change in each loop iteration below
// 	var (
// 		// Check original src and tgt types to see if they are a Maybe[T] - a *Maybe[T] makes no sense
// 		// If not, then they will be nil
// 		srcMaybeTyp, tgtMaybeTyp goreflect.Type
//
// 		// Actual types are generic type of Maybe or the original type
// 		// srcActTyp, tgtActTyp goreflect.Type
//
// 		// Base types, which may be different
// 		srcBaseType, tgtBaseType goreflect.Type
//
// 		// Types that we actually work with after possibly converting to a base type
// 		srcTyp, tgtTyp goreflect.Type
//
// 		// Func to return to caller, and do we have a func
// 		fn     func(any, any) error
// 		haveIt bool
// 	)
//
// 	for i := 'A'; i <= 'D'; i++ {
//
// 		switch i {
// 		case 'A': // A: Use types as given
//   		// Are src or tgt Maybe?
//   		srcMaybeTyp = reflect.GetMaybeType(src)
//   		tgtMaybeTyp = reflect.GetMaybeType(tgt)
//
//   		// Actual types are generic type of Maybe or the original type
//   		srcTyp = funcs.Ternary(srcMaybeTyp == nil, src, srcMaybeTyp)
//   		tgtTyp = funcs.Ternary(tgtMaybeTyp == nil, tgt, tgtMaybeTyp)
//
// 		case 'B': // B: If src is a subtype, use base type
//   		// Are src or tgt Maybe?
//   		srcMaybeTyp = reflect.GetMaybeType(src)
//   		tgtMaybeTyp = reflect.GetMaybeType(tgt)
//
//   		// Actual types are generic type of Maybe or the original type
//   		srcTyp = funcs.Ternary(srcMaybeTyp == nil, src, srcMaybeTyp)
//   		tgtTyp = funcs.Ternary(tgtMaybeTyp == nil, tgt, tgtMaybeTyp)
//
//   		// Base types
//   		srcBaseType = reflect.TypeToBaseType(srcTyp)
//   		tgtBaseType = reflect.TypeToBaseType(tgtTyp)
//
// 			if (srcTyp == srcBaseType) || (tgt != tgtBaseType) {
// 				// Skip if src is not a subtype
// 				continue
// 			}
//
// 		case 'C': // C: If tgt is a subtype, use base type
// 			if (src != srcBaseType) || (tgt == tgtBaseType) {
// 				// Skip if tgt is not a subtype
// 				continue
// 			}
// 			srcTyp = src
// 			tgtTyp = tgtBaseType
//
// 		default: // D: If both are subtypes, use base types
// 			srcTyp = srcBaseType
// 			tgtTyp = tgtBaseType
// 		}
//
// 		// 1. source -> target
// 		if fn, haveIt = convertFromTo[srcTyp.String()+tgtTyp.String()]; (!haveIt) && (!srcIsPtr) && (!tgtIsPtr) && (srcTyp == tgtTyp) {
// 			// Copy source -> target
// 			fn = func(in, out any) error {
// 				inVal, outVal := goreflect.ValueOf(in), goreflect.ValueOf(out)
//
// 				if srcMaybeTyp != nil {
// 					if mVal := reflect.GetMaybeValue(inVal); !mVal.IsValid() {
// 						// Can't set target to a value that doesn't exist
// 						return fmt.Errorf(errEmptyMaybeMsg, srcTyp.String(), tgtTyp.String())
// 					} else {
// 						// Replace inVal with present Maybe value
// 						inVal = mVal
// 					}
// 				}
//
// 				// // Convert inVal
// 				// inVal = inVal.Convert(tgtTyp)
//
// 				if tgtMaybeTyp != nil {
// 					// conv.To Target *Maybe[T] cannot be derefd, as Maybe.Set requires a pointer receiver
// 					reflect.SetMaybeValue(outVal, inVal)
// 				} else {
// 					// Target is not a *Maybe[T], so deref to set it
// 					outVal.Elem().Set(inVal)
// 				}
//
// 				return nil
// 			}
//
// 			haveIt = true
// 		}
//
// 		// 2. derefd source -> target
// 		if (!haveIt) && srcIsPtr && (!tgtIsPtr) {
// 			if lufn, luHaveIt := convertFromTo[srcTyp.Elem().String()+tgtTyp.String()]; luHaveIt {
// 				// Have to use wrapper func that derefs source
// 				fn = func(in, out any) error {
// 					inVal, outVal := goreflect.ValueOf(in), goreflect.ValueOf(out)
//
// 					if inVal.IsNil() {
// 						return fmt.Errorf(errConvertNilSourceMsg, srcTyp, tgtTyp)
// 					}
//
// 					return lufn(inVal.Elem().Interface(), outVal)
// 				}
//
// 				haveIt = true
// 			} else if srcTyp.Elem() == tgtTyp {
// 				// Copy *source -> target
// 				fn = func(in, out any) error {
// 					rin := goreflect.ValueOf(in)
//
// 					if rin.IsNil() {
// 						return fmt.Errorf(errCopyNilSourceMsg, srcTyp, tgtTyp)
// 					}
//
// 					goreflect.ValueOf(out).Elem().Set(goreflect.ValueOf(in).Elem().Convert(tgtTyp))
//
// 					return nil
// 				}
//
// 				haveIt = true
// 			}
// 		}
//
// 		// 3. source -> derefd target
// 		if (!haveIt) && (!srcIsPtr) && tgtIsPtr {
// 			if lufn, luHaveIt := convertFromTo[srcTyp.String()+tgtTyp.Elem().String()]; luHaveIt {
// 				// Have to use wrapper func that derefs target
// 				fn = func(in, out any) error {
// 					return lufn(in, goreflect.ValueOf(out).Elem().Interface())
// 				}
//
// 				haveIt = true
// 			} else if srcTyp == tgtTyp.Elem() {
// 				// Copy source -> *target
// 				fn = func(in, out any) error {
// 					rout := goreflect.ValueOf(out)
//
// 					if rout.IsNil() || rout.Elem().IsNil() {
// 						return fmt.Errorf(errCopyNilTargetMsg, srcTyp, tgtTyp)
// 					}
//
// 					rout.Elem().Elem().Set(goreflect.ValueOf(in).Convert(tgtTyp.Elem()))
//
// 					return nil
// 				}
// 				haveIt = true
// 			}
// 		}
//
// 		// 4. derefd source -> derefd target
// 		if (!haveIt) && srcIsPtr && tgtIsPtr {
// 			if lufn, luHaveIt := convertFromTo[srcTyp.Elem().String()+tgtTyp.Elem().String()]; luHaveIt {
// 				// Have to use wrapper func that derefs source and target
// 				fn = func(in, out any) error {
// 					return lufn(goreflect.ValueOf(in).Elem().Interface(), goreflect.ValueOf(out).Elem().Interface())
// 				}
//
// 				haveIt = true
// 			} else if srcTyp.Elem() == tgtTyp.Elem() {
// 				// Copy *source -> *target
// 				fn = func(in, out any) error {
// 					rin, rout := goreflect.ValueOf(in), goreflect.ValueOf(out)
//
// 					if rout.IsNil() || ((!rin.IsNil()) && rout.Elem().IsNil()) {
// 						return fmt.Errorf(errCopyNilTargetMsg, srcTyp, tgtTyp)
// 					}
//
// 					if rin.IsNil() {
// 						rout.Elem().Set(rin)
// 					} else {
// 						rout.Elem().Elem().Set(rin.Elem().Convert(tgtTyp.Elem()))
// 					}
//
// 					return nil
// 				}
//
// 				haveIt = true
// 			}
// 		}
//
// 		if haveIt {
// 			// Recheck i, and further wrap if subtypes are used
// 			switch i {
// 			case 'A': // Use types as given
//
// 			case 'B':
// 				{ // src is subtype
// 					// Have to use a wrapper func that converts src subtype to src basetype
// 					lufn := fn
// 					sbfn := func(in, out any) error {
// 						// in = subtype, tgt = *type
// 						// lufn = func(basetype, tgt)
// 						return lufn(reflect.ValueToBaseType(goreflect.ValueOf(in)).Interface(), out)
// 					}
// 					fn = sbfn
// 				}
//
// 			case 'C':
// 				{ // tgt is subtype
// 					// Have to use a wrapper func that converts tgt subtype to tgt basetype
// 					lufn := fn
// 					tbfn := func(in, out any) error {
// 						// in = src, tgt = *subtype
// 						// lufn = func(src, *basetype)
// 						return lufn(in, reflect.ValueToBaseType(goreflect.ValueOf(out)).Interface())
// 					}
// 					fn = tbfn
// 				}
//
// 			default:
// 				{ // src and tgt are subtypes
// 					// Have to use a wrapper func that converts src subtype to src basetype and tgt subtype to tgt basetype
// 					lufn := fn
// 					tbfn := func(in, out any) error {
// 						// in = subtype, tgt = *subtype
// 						// lufn = func(basetype, *basetype)
// 						return lufn(
// 							reflect.ValueToBaseType(goreflect.ValueOf(in)).Interface(),
// 							reflect.ValueToBaseType(goreflect.ValueOf(out)).Interface(),
// 						)
// 					}
// 					fn = tbfn
// 				}
// 			}
//
// 			break
// 		}
// 	}
//
// 	// Ensure we have found a conversion (possibly wrapped for derefing)
// 	if !haveIt {
// 		// Conversion is not registered
// 		return nil, nil
// 	}
//
// 	// Conversion is registered (or it is a copy)
// 	return fn, nil
// }

// RegisterConversion registers a conversion from a value of type S to a value of type T.
// Note the conversion function must accept a pointer type for the target.
// If the target is already a pointer type, an additional level of pointer is required.
//
// Examples:
// RegisterConversion(0, Foo{}, func(int, *Foo) error {...})
// RegisterConversion((*int)(nil), (*Foo)(nil), func(*int, **Foo) error {...})
//
// See LookupConversion for details
func RegisterConversion[S, T any](convFn func(S, *T) error) error {
	var (
		fn         = goreflect.TypeOf(convFn)
		sTyp, tTyp = fn.In(0), fn.In(1).Elem()
		convKey    = sTyp.String() + tTyp.String()
	)

	// See if a conversion exists for the exact types given.
	// If not, register it without bothering to use LookupConversion.
	// This allows registering functions when the types or the same, or a similar conversion exists for convertible type(s).

	if _, haveIt := convertFromTo[convKey]; haveIt {
		// Return error that the conversion already exists
		return fmt.Errorf(errRegisterExistsMsg, sTyp, tTyp)
	}

	// Store the conversion in the map - we have to store a func(any, any), so generate one
	convertFromTo[convKey] = func(src, tgt any) error {
		return convFn(src.(S), tgt.(*T))
	}

	return nil
}

// MustRegisterConversion is a must version of RegisterConversion
func MustRegisterConversion[S, T any](convFn func(S, *T) error) {
	funcs.Must(RegisterConversion(convFn))
}

// To converts any supported combination of source and target types.
//
// The actual conversion is performed by:
// - functions declared in this source file
// - functions registered by other packages in this library
// - functions registered by other packages outside this library
//
// The source is typed any for two reasons:
// - to allow for cases where the caller accepts type any
// - arbitrary new conversions can be registered (ideally via an init function)
//
// The target is a *T because Go can infer generic parameters, but not generic return types.
// So instead of writing this:
//
//	var str = conv.To[int, string](0)
//
// We write this:
//
//	var str string
//	conv.To(0, &str)
//
// It is a design choice to not make the user constantly repeat generic types for every conversion.
//
// See LookupConversion for the algorithm to find a registered conversion function.
// There are 4 cases:
// 1. LookupConversion does not find a conversion, the conversion is allowable, the types are the same
//   - The source value is copied to the target
//   - If the source is a pointer, the target gets a copy of the pointer, so source and target point to same address
//
// 2. LookupConversion does not find a conversion, the conversion is allowable, the types are different
//   - Returns error that source cannot be converted to target
//
// 3. LookupConversion returns an error
//   - The error is returned as is
//
// 4. LookupConversion finds a conversion:
//   - The conversion is applied
func To[T any](src any, tgt *T) error {
	var (
		valsrc = goreflect.ValueOf(src)
		valtgt = goreflect.ValueOf(tgt)
		srcTyp = valsrc.Type()
		tgtTyp = valtgt.Type().Elem()
	)
	// fmt.Printf("To: %s -> %s\n", srcTyp, tgtTyp)

	// Use LookupConversion to find the conversion function, if it exists
	fn, err := LookupConversion(srcTyp, tgtTyp)
	// fmt.Printf("Lookup: %p, %s\n", fn, err)
	switch {
	case (fn == nil) && (err == nil):
		// No conversion exists, but could be registered
		return fmt.Errorf(errLookupMsg, srcTyp, tgtTyp)

	case err != nil:
		// Conversion will never be possible
		return err

	default:
		// Conversion exists
		return fn(src, tgt)
	}
}

// MustTo is a Must version of To
func MustTo[T any](src any, tgt *T) {
	funcs.Must(To(src, tgt))
}

// ToBigOps is the BigOps version of To
func ToBigOps[S constraint.Numeric | ~string, T constraint.BigOps[T]](src S, tgt *T) error {
	var (
		valsrc = goreflect.ValueOf(src)
		valtgt = goreflect.ValueOf(tgt)
	)

	// Convert source to base type
	valsrc = reflect.ValueToBaseType(valsrc)

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
