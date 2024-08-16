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
	errMsg        = "The %T value of %s cannot be converted to %s"
	errONonNilMsg = "The target value of type %T cannot be nil"
	errAnyToInvalidIMsg = "AnyTo cannot convert the input type %s"
	errReflectToInvalidSrc = fmt.Errorf("ReflectTo source cannot be Invalid")
  errReflectToInvalidTgt = fmt.Errorf("ReflectTo target cannot be Invalid")
  errReflectToTgtMustBePtr = fmt.Errorf("ReflectTo target must be be a pointer")
  errReflectToTgtBigTypeMsg = "The target value of type %T is invalid: big types have to be a **"
  errReflectToLookupMsg = "There is no conversion function from %s to %s"

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
)

// To converts any Numeric type or string to any Numeric type or string
// If the types are the same, a copy by value is performed, unless they are big types.
// For big type copies, a new pointer is constructed with a copy of the input value.
// This allows the big copy to be modified without affecting the original big value.
//
// Note that subtypes are handled automatically by the generic constraints.
func To[I, O constraint.Numeric | string](i I, o *O) error {
	// Target cannot be nil
	if o == nil {
		return fmt.Errorf(errONonNilMsg, o)
	}

	// Get reflection info on i and o
  // If i and/or o is a subtype, convert it to the base type, so we can find a conversion function
	var (
		ival = reflect.ValueToBaseType(goreflect.ValueOf(i))
		ityp = ival.Type()
		oval = reflect.ValueToBaseType(goreflect.ValueOf(o)) // o cannot be nil, but o.Elem() can be nil
		otyp = oval.Type().Elem()   // o.Type().Elem() cannot be nil
	)

	// If the types are the same, then a simple copy will suffice
	if ityp == otyp {
		// Are they big types?
		if reflect.IsBigPtr(ityp) {
			// If the input is nil, make the output nil
			if ival.IsNil() {
				oval.Elem().Set(ival)
			} else {
				// If the oval is the same pointer as ival or nil, allocate it
				if (ival.Interface() == oval.Elem().Interface()) || oval.Elem().IsNil() {
					oval.Elem().Set(goreflect.New(otyp.Elem()))
				}

				// ival is non-nil, oval is non-nil and not the same pointer as ival
				// Copy the value
				oval.Elem().Elem().Set(ival.Elem())
			}
		} else {
			// All non-big types are value types, just copy the value
			oval.Elem().Set(ival)
		}

		return nil
	}

	// Construct a string of the input and output types (eg "int8int" means int8 -> int)
	// Use the string as an index into the convertFromTo map
	return convertFromTo[ityp.String()+otyp.String()](any(ival.Interface()), any(oval.Interface()))
}

// MustTo is a Must version of To
func MustTo[I, O constraint.Numeric | string](i I, o *O) {
	funcs.Must(To(i, o))
}

// AnyTo is a version of To that accepts input values of type any
// The input value must still satisfy constraint.Numeric | string
func AnyTo[O constraint.Numeric | string](i any, o *O) error {
  var (
    iv   =  goreflect.ValueOf(i)
    ival = reflect.ValueToBaseType(iv).Interface()
  )
  
	switch iv.Kind() {
	case goreflect.Int:
		return To(ival.(int), o)
	case goreflect.Int8:
		return To(ival.(int8), o)
	case goreflect.Int16:
		return To(ival.(int16), o)
	case goreflect.Int32:
		return To(ival.(int32), o)
	case goreflect.Int64:
		return To(ival.(int64), o)
	case goreflect.Uint:
		return To(ival.(uint), o)
	case goreflect.Uint8:
		return To(ival.(uint8), o)
	case goreflect.Uint16:
		return To(ival.(uint16), o)
	case goreflect.Uint32:
		return To(ival.(uint32), o)
	case goreflect.Uint64:
		return To(ival.(uint64), o)
	case goreflect.Float32:
		return To(ival.(float32), o)
	case goreflect.Float64:
		return To(ival.(float64), o)
  case goreflect.String:
    return To(ival.(string), o)
	case goreflect.Ptr:
		if bi, isa := i.(*big.Int); isa {
			return To(bi, o)
		} else if bf, isa := i.(*big.Float); isa {
			return To(bf, o)
		} else if br, isa := i.(*big.Rat); isa {
			return To(br, o)
		}
	}

	return fmt.Errorf(errAnyToInvalidIMsg, iv.Type())
}

// ToBigOps is the BigOps version of To
func ToBigOps[I constraint.Numeric | string, O constraint.BigOps[O]](i I, o *O) error {
	// Target cannot be nil
	if o == nil {
		return fmt.Errorf(errONonNilMsg, o)
	}

	var (
		ival = goreflect.ValueOf(i)
		ityp = ival.Type()
		oval = goreflect.ValueOf(o) // o cannot be nil, but o.Elem() can be nil
		otyp = oval.Type().Elem()   // o.Type().Elem() cannot be nil
	)

	// If the types are the same, then a simple copy will suffice
	if ityp == otyp {
		// If the input is nil, make the output nil
		if ival.IsNil() {
			oval.Elem().Set(ival)
		} else {
			// If the oval is the same pointer as ival or nil, allocate it
			if (ival.Interface() == oval.Elem().Interface()) || oval.Elem().IsNil() {
				oval.Elem().Set(goreflect.New(otyp.Elem()))
			}

			// ival is non-nil, oval is non-nil and not the same pointer as ival
			// Copy the value
			oval.Elem().Elem().Set(ival.Elem())
		}

		return nil
	}

	// Types differ, lookup conversion using types and execute it, returning result
	return convertFromTo[ityp.String()+otyp.String()](ival.Interface(), oval.Interface())
}

// MustToBigOps is a Must version of ToBigOps
func MustToBigOps[I constraint.Numeric | string, O constraint.BigOps[O]](i I, o *O) {
	funcs.Must(ToBigOps(i, o))
}

// ReflectTo uses reflection objects to convert from source to target.
// This function is useful for reflection algorithms that need to do conversions.
// The tgt must wrap a pointer.
func ReflectTo(i, o goreflect.Value) error {
  // Die if i is invalid
  if !i.IsValid() {
    return errReflectToInvalidSrc
  }

  // Die if o is invalid
  if !o.IsValid() {
    return errReflectToInvalidTgt
  }
  
  // Die if o is not a pointer
  if o.Kind() != goreflect.Pointer {
    return errReflectToTgtMustBePtr
  }
  
  var (
   ityp = i.Type()
   otyp = o.Type()
  )
  
  // Die if o is nil
  if o.IsNil() {
    return fmt.Errorf(errONonNilMsg, otyp)
  }
  
  // Die if o is a big type that is a value or only one pointer
  if (
    ((otyp.Kind() == goreflect.Struct) && reflect.IsBigPtr(goreflect.PointerTo(otyp))) ||
    reflect.IsBigPtr(otyp)) {
    return fmt.Errorf(errReflectToTgtBigTypeMsg, otyp)
  }

  // Convert output to a base type
  ob := reflect.ValueToBaseType(o)

  // Locate a conversion function in convertFromTo map  
  convFn := convertFromTo[(ityp.String()+ob.Type().Elem().String())]
  if convFn == nil {
    return fmt.Errorf(errReflectToLookupMsg, ityp, otyp.Elem())
  }
  
  return convFn(i.Interface(), ob.Interface())
}

// MustReflectTo is a Must version of ReflectTo
func MustReflectTo(i, o goreflect.Value) {
  funcs.Must(ReflectTo(i, o))
}
