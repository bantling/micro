package util

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVOf(t *testing.T) {
	kv := KVOf(5, "5")
	assert.Equal(t, 5, kv.Key)
	assert.Equal(t, "5", kv.Value)
}
