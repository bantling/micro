package math

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"

	"github.com/bantling/micro/conv"
)

// init registers conversions between decimal and other types
func init() {
  // Signed integers
	conv.RegisterConversion[int, Decimal](0, nil, func(src int, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })
	conv.RegisterConversion[int8, Decimal](0, nil, func(src int8, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })
	conv.RegisterConversion[int16, Decimal](0, nil, func(src int16, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })
	conv.RegisterConversion[int32, Decimal](0, nil, func(src int32, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })
	conv.RegisterConversion[int64, Decimal](0, nil, func(src int64, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })

  // Unsigned integers
	conv.RegisterConversion[uint, Decimal](0, nil, func(src uint, tgt *Decimal) (err error) {
    var isrc int64
    if err = conv.To(src, &isrc); err != nil {
      return
    }

    *tgt, err = OfDecimal(isrc, 0)
    return
  })
	conv.RegisterConversion[uint8, Decimal](0, nil, func(src uint8, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })
	conv.RegisterConversion[uint16, Decimal](0, nil, func(src uint16, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })
	conv.RegisterConversion[uint32, Decimal](0, nil, func(src uint32, tgt *Decimal) (err error) { *tgt, err = OfDecimal(int64(src), 0); return })
	conv.RegisterConversion[uint64, Decimal](0, nil, func(src uint64, tgt *Decimal) (err error) {
    var isrc int64
    if err = conv.To(src, &isrc); err != nil {
      return
    }

    *tgt, err = OfDecimal(isrc, 0)
    return
  })

  // *big.Int and *big.Rat
	conv.RegisterConversion[*big.Int, Decimal](nil, nil, func(src *big.Int, tgt *Decimal) (err error) {
    var isrc int64
    if err = conv.To(src, &isrc); err != nil {
      return
    }

    *tgt, err = OfDecimal(isrc, 0)
    return
  })

	conv.RegisterConversion[*big.Rat, Decimal](nil, nil, func(src *big.Rat, tgt *Decimal) (err error) {
    var isrc int64
    if err = conv.To(src, &isrc); err != nil {
      return
    }

    *tgt, err = OfDecimal(isrc, 0)
    return
  })
}
