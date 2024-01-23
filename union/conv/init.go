package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/conv"
)

func init() {
	// Register the MaybeWrapperInfo so conv.To can handle Maybe wrappers
	var mwi MaybeWrapperInfo
	conv.MustRegisterWrapper(mwi)
}
