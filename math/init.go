package math

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math/big"

	"github.com/bantling/micro/conv"
)

// init registers conversions between math types (decimal and range) and go standard types
func init() {
	// ==== To Decimal

	// Signed integers
	conv.MustRegisterConversion(func(src int, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })
	conv.MustRegisterConversion(func(src int8, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })
	conv.MustRegisterConversion(func(src int16, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })
	conv.MustRegisterConversion(func(src int32, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })
	conv.MustRegisterConversion(func(src int64, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })

	// Unsigned integers
	conv.MustRegisterConversion(func(src uint, tgt *Decimal) (err error) {
		var isrc int64
		if err = conv.To(src, &isrc); err != nil {
			return
		}

		*tgt, err = OfDecimal(isrc, 0)
		return
	})
	conv.MustRegisterConversion(func(src uint8, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })
	conv.MustRegisterConversion(func(src uint16, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })
	conv.MustRegisterConversion(func(src uint32, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })
	conv.MustRegisterConversion(func(src uint64, tgt *Decimal) (err error) {
		var isrc int64
		if err = conv.To(src, &isrc); err != nil {
			return
		}

		*tgt, err = OfDecimal(isrc, 0)
		return
	})

	// *big.Int and *big.Rat
	conv.MustRegisterConversion(func(src *big.Int, tgt *Decimal) (err error) {
		var isrc int64
		if err = conv.To(src, &isrc); err != nil {
			return
		}

		*tgt, err = OfDecimal(isrc, 0)
		return
	})

	conv.MustRegisterConversion(func(src *big.Rat, tgt *Decimal) (err error) {
		var isrc int64
		if err = conv.To(src, &isrc); err != nil {
			return
		}

		*tgt, err = OfDecimal(isrc, 0)
		return
	})

	// String
	conv.MustRegisterConversion(func(src string, tgt *Decimal) (err error) {
		*tgt, err = StringToDecimal(src)
		return
	})

	// ==== From Decimal

	// *big.Int and *big.Rat
	conv.MustRegisterConversion(func(src Decimal, tgt **big.Int) (err error) {
		if src.scale > 0 {
			err = fmt.Errorf(errToBigIntMsg, src)
			return
		}

		err = conv.To(src.value, tgt)

		return
	})

	conv.MustRegisterConversion(func(src Decimal, tgt **big.Rat) (err error) {
		err = conv.To(src.String(), tgt)
		return
	})

	// String
	conv.MustRegisterConversion(func(src Decimal, tgt *string) (err error) {
		*tgt = src.String()
		return
	})
}
