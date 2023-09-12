package json

// SPDX-License-Identifier: Apache-2.0

import (
  "github.com/bantling/micro/conv"
)

// init registers conversions between json.Value and other types
func init() {
      conv.RegisterConversion[int, Value](0, nil, func(src int, tgt *Value) error {*tgt = MustNumberToValue(src); return nil})
}
