package code

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperator_(t *testing.T) {
	assert.Equal(t, uint(afterUnary), uint(Add))
	assert.Equal(t, uint(afterUnary+1), uint(Sub))

	assert.Equal(t, uint(afterBinary), uint(Not))
	assert.Equal(t, uint(afterBinary+1), uint(And))
}
