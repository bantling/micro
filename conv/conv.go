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
	unionreflect "github.com/bantling/micro/union/reflect"
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
  errRegisterNilFuncMsg          = "The conversion from %s to %s requires a non-nil conversion function"
  errRegisterWrapperInfoExistsMsg = "The wrapper type %s has already been registered"

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

	// map strings of types that may be a nil wrapper to a func(any) bool that tests if the instance wraps nil.
  // this map is only populated by other packages, as Go has no such standard types.
	wrapperTypes = map[string]WrapperInfo{}

	badConversionKinds = map[goreflect.Kind]bool{
		goreflect.Uintptr:       true,
		goreflect.Chan:          true,
		goreflect.Func:          true,
		goreflect.UnsafePointer: true,
	}
)

// ==== To and functions that support it

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
    zvs S
    zvt T
		sTyp, tTyp = goreflect.TypeOf(zvs), goreflect.TypeOf(zvt)
		convKey    = sTyp.String() + tTyp.String()
	)

  // Error if the conversion func is nil
  if convFn == nil {
    return fmt.Errorf(errRegisterNilFuncMsg, sTyp, tTyp)
  }

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

// WrapperInfo describes a wrapper type that can store one or more specific types, and may be able to represent an empty value
type WrapperInfo interface {
  // PackagePath returns the package path of the wrapper type
  PackagePath() string

  // TypeNamePrefix returns the type name prefix. EG, if wrapper is Foo[T] then the prefix is Foo.
  TypeNamePrefix() string

  // AcceptsType returns true if the given wrapper instance is capable of storing the given type.
  // A wrapper type may be capable of storing any type, or only a set of specific - possibly unrelated - types.
  AcceptsType(instance goreflect.Value, typ goreflect.Type) bool

  // CanBeEmpty returns true if the wrapper can hold an empty value.
  CanBeEmpty(instance goreflect.Value) bool

  // ConvertibleTo returns true if the given wrapper instance is capable of returning the given type.
  // Actually converting to the specified type can still fail.
  // EG, a type may be capable of converting to int in general, but the current value may lie outside the range of an int.
  // Passing any type where AcceptsType returns true will always return true.
  ConvertibleTo(instance goreflect.Value, typ goreflect.Type) bool

  // Get a value of the given type, which may require a conversion.
  // Even if the wrapper can be converted to the specified type, it may still fail to convert.
  // See ConvertibleTo.
  //
  // Errors if:
  // - ConvertibleTo(Type) returns false
  // - ConvertibleTo(Type) returns true, but the specific value stored cannot be converted to Type.
  Get(instance goreflect.Value, typ goreflect.Type) (goreflect.Value, bool, error)

  // Set to the given value if the bool is true, else store an empty value (ignoring the value passed) if the bool is false.
  // When the bool is true, if AcceptsType returns true for the type of value given, Set cannot fail.
  // If CanBeEmpty returns true, then providing a bool flag of false cannot fail.
  // See AcceptsType.
  //
  // Errors if:
  // - AcceptsType returns false for the type of value given and the bool is true
  // - CanBeEmpty returns false and the bool is false (the value is irrelevant)
  Set(instance, val goreflect.Value, present bool) error
}

// RegisterWrapper allows other packages to register types that hold an empty value.
// The only error condition is if the same type is registered twice.
func RegisterWrapper(wi WrapperInfo) error {
  // Map key is package name "." type name prefix
  key := fmt.Sprintf("%s.%s", wi.PackagePath(), wi.TypeNamePrefix())

  // Error the wrapper type has already been registered
  if _, haveIt := wrapperTypes[key]; haveIt {
    return fmt.Errorf(errRegisterWrapperInfoExistsMsg, key)
  }

  // Register funcs
  wrapperTypes[key] = wi

  return nil
}

// MustRegisterEmptyWrapper is a must verison of RegisterEmptyWrapper
func MustRegisterWrapper(wi WrapperInfo) {
  funcs.Must(RegisterWrapper(wi))
}

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
// Example lookups, listing possible conversions in search order (shown without the extra * the target type has to have):
// - int to string -> int to string
// - subint to *string -> subint to string, int to string
// - Maybe[int] to string -> int to string
// - int to int -> copy
// - subint to Maybe[*int] -> subint to *int, subint to int, copy
//
// Additionally, other packages can register custom functions to test if an instance of a type is effectively nil, similar
// to an empty Maybe. The returned conversion will do a precheck to see if the src value is effectively nil, and if so,
// proceed as if it is nil.
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
		srcVal, srcBase, srcPtr, srcPtrBase, srcMaybe, srcMaybeBase, srcMaybePtr, srcMaybePtrBase goreflect.Type
		tgtVal, tgtBase, tgtPtr, tgtPtrBase, tgtMaybe, tgtMaybeBase, tgtMaybePtr, tgtMaybePtrBase goreflect.Type
		convFn                                                                                    func(any, any) error
		haveIt                                                                                    bool
	)

	srcVal = src
	srcBase = reflect.TypeToBaseType(src)

	if srcPtr = funcs.TernaryResult(src.Kind() == goreflect.Pointer, src.Elem, nil); srcPtr != nil {
		srcVal = nil
		srcPtrBase = reflect.TypeToBaseType(srcPtr)

		if unionreflect.GetMaybeType(srcPtr) != nil {
			// Cannot have a *Maybe, that makes no sense
			return nil, fmt.Errorf(errLookupMsg, src, tgt)
		}
	}

	if srcMaybe = unionreflect.GetMaybeType(src); srcMaybe != nil {
		srcVal = nil
		srcMaybeBase = reflect.TypeToBaseType(srcMaybe)

		if srcMaybePtr = funcs.TernaryResult(srcMaybe.Kind() == goreflect.Pointer, srcMaybe.Elem, nil); srcMaybePtr != nil {
			srcMaybePtrBase = reflect.TypeToBaseType(srcMaybePtr)
		}
	}

	tgtVal = tgt
	tgtBase = reflect.TypeToBaseType(tgt)

	if tgtPtr = funcs.TernaryResult(tgt.Kind() == goreflect.Pointer, tgt.Elem, nil); tgtPtr != nil {
		tgtVal = nil
		tgtPtrBase = reflect.TypeToBaseType(tgtPtr)

		if unionreflect.GetMaybeType(tgtPtr) != nil {
			// Cannot have a *Maybe, that makes no sense
			return nil, fmt.Errorf(errLookupMsg, src, tgt)
		}
	}

	if tgtMaybe = unionreflect.GetMaybeType(tgt); tgtMaybe != nil {
		tgtVal = nil
		tgtMaybeBase = reflect.TypeToBaseType(tgtMaybe)

		if tgtMaybePtr = funcs.TernaryResult(tgtMaybe.Kind() == goreflect.Pointer, tgtMaybe.Elem, nil); tgtMaybePtr != nil {
			tgtMaybePtrBase = reflect.TypeToBaseType(tgtMaybePtr)
		}
	}

	for _, srcTyp := range []goreflect.Type{
		srcVal, srcBase, srcPtr, srcPtrBase, srcMaybe, srcMaybeBase, srcMaybePtr, srcMaybePtrBase,
	} {
		for _, tgtTyp := range []goreflect.Type{
			tgtVal, tgtBase, tgtPtr, tgtPtrBase, tgtMaybe, tgtMaybeBase, tgtMaybePtr, tgtMaybePtrBase,
		} {
			// Cannot lookup conversions for types that don't exist
			if (srcTyp != nil) && (tgtTyp != nil) {
				convFn, haveIt = nil, srcTyp == tgtTyp
				if !haveIt {
					convFn, haveIt = convertFromTo[srcTyp.String()+tgtTyp.String()]
				}

				if haveIt {
					// Generate a function to unwrap the src type and read it
					var srcFn func(goreflect.Value) goreflect.Value
					switch srcTyp {
					case src:
						srcFn = func(s goreflect.Value) goreflect.Value { return s }
					case srcBase:
						srcFn = func(s goreflect.Value) goreflect.Value { return s.Convert(srcBase) }
					case srcPtr:
						srcFn = func(s goreflect.Value) goreflect.Value {
							if s.IsValid() {
								return s.Elem()
							}
							return s
						}
					case srcPtrBase:
						srcFn = func(s goreflect.Value) goreflect.Value {
							if s.IsValid() && (!s.IsNil()) {
								return s.Elem().Convert(srcPtrBase)
							}
							return s
						}
					case srcMaybe:
						srcFn = func(s goreflect.Value) goreflect.Value { return unionreflect.GetMaybeValue(s) }
					case srcMaybeBase:
						srcFn = func(s goreflect.Value) goreflect.Value {
							if temp := unionreflect.GetMaybeValue(s); temp.IsValid() {
								return temp.Convert(srcMaybeBase)
							} else {
								return temp
							}
						}
					case srcMaybePtr:
						srcFn = func(s goreflect.Value) goreflect.Value {
							if temp := unionreflect.GetMaybeValue(s); temp.IsValid() {
								return temp.Elem()
							} else {
								return temp
							}
						}
					case srcMaybePtrBase:
						srcFn = func(s goreflect.Value) goreflect.Value {
							if temp := unionreflect.GetMaybeValue(s); temp.IsValid() {
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
						tgtFn = func(temp, t goreflect.Value) {
							if t.Elem().IsNil() {
								t.Elem().Set(temp)
							} else {
								t.Elem().Elem().Set(temp.Elem())
							}
						}
					case tgtPtrBase:
						tgtFn = func(temp, t goreflect.Value) {
							if t.Elem().IsNil() {
								t.Elem().Set(temp.Convert(t.Elem().Type()))
							} else {
								t.Elem().Elem().Set(temp.Elem().Convert(tgt.Elem()))
							}
						}
					case tgtMaybe:
						tgtFn = func(temp, t goreflect.Value) { unionreflect.SetMaybeValue(t, temp.Elem()) }
					case tgtMaybeBase:
						tgtFn = func(temp, t goreflect.Value) { unionreflect.SetMaybeValue(t, temp.Elem().Convert(tgtMaybe)) }
					case tgtMaybePtr:
						tgtFn = func(temp, t goreflect.Value) { unionreflect.SetMaybeValue(t, temp) }
					case tgtMaybePtrBase:
						tgtFn = func(temp, t goreflect.Value) { unionreflect.SetMaybeValue(t, temp.Convert(goreflect.PtrTo(tgtMaybePtr))) }
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
						// - Registered conversions panic with errors that indicate if the source or target type is the problem
						srcVal, tgtVal := goreflect.ValueOf(s), goreflect.ValueOf(t)

						// The source may be an untyped nil
						if srcVal.IsValid() {
							// If not, then assert types match
							reflect.MustTypeAssert(srcVal, src, "source")
						}

						// The target cannot be an untyped nil:
						// - conv.To accepts *T, which camn only be a typed nil
						// - reflect.Value.Call requires all arguments to be a valid Value object

						// The target may be a typed nil
						if reflect.IsNil(tgtVal) {
							return fmt.Errorf(errCopyNilTargetMsg, src, tgt)
						}
						// If the target is not nil, then assert types match
						reflect.MustTypeAssert(tgtVal, goreflect.PtrTo(tgt), "target")

						// Unwrap src value, which will be invalid for nil ptr or empty maybe
						srcVal = srcFn(srcVal)

						if reflect.IsNil(srcVal) {
							// Tgt must be nillable or maybe
							if reflect.IsNillable(tgt) {
								// Tgt is nillable)
								tgtVal.Elem().SetZero()
							} else if tgtMaybe != nil {
								// Tgt is a Maybe
								unionreflect.SetMaybeValueEmpty(tgtVal)
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

						// Convert source -> unwrapped target
						err := convFn(funcs.TernaryResult(srcVal.IsValid(), srcVal.Interface, nil), temp.Interface())

						// Wrap target value only if no error occurred - the target is unmodified if the conversion fails
						if err == nil {
							tgtFn(temp, tgtVal)
						}

						// Return any error
						return err
					}, nil
				}
			}
		}
	}

	return nil, nil
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
