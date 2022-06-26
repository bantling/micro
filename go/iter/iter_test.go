// SPDX-License-Identifier: Apache-2.0

package iter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceIterGen(t *testing.T) {
	// nil
	var slc []int
	assert.Nil(t, slc)
	iter := SliceIterGen(slc)
	i, haveIt := iter()
	assert.Equal(t, 0, i)
	assert.False(t, haveIt)
	assert.Equal(t, 0, i)
	assert.False(t, haveIt)

	// two elements
	slc = []int{1, 2}
	iter = SliceIterGen(slc)
	i, haveIt = iter()
	assert.Equal(t, 1, i)
	assert.True(t, haveIt)
	i, haveIt = iter()
	assert.Equal(t, 2, i)
	assert.True(t, haveIt)
	i, haveIt = iter()
	assert.Equal(t, 0, i)
	assert.False(t, haveIt)
	assert.Equal(t, 0, i)
	assert.False(t, haveIt)
}

func TestMapIterGen(t *testing.T) {
	// nil
	var src map[string]int
	assert.Nil(t, src)
	iter := MapIterGen(src)
	kv, haveIt := iter()
	assert.Equal(t, KeyValue[string, int]{}, kv)
	assert.False(t, haveIt)

	// two pairs
	src = map[string]int{"a": 1, "b": 2}
	dst := map[string]int{}
	iter = MapIterGen(src)
	kv, haveIt = iter()
	assert.True(t, haveIt)
	dst[kv.Key] = kv.Value
	kv, haveIt = iter()
	assert.True(t, haveIt)
	dst[kv.Key] = kv.Value
	assert.Equal(t, src, dst)
}

func TestReaderIterGen(t *testing.T) {
	//
}
