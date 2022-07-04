// SPDX-License-Identifier: Apache-2.0

package iter

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIter(t *testing.T) {
	it := NewIter(SliceIterGen[int]([]int{1, 2}))
	assert.True(t, it.Next())
	assert.Equal(t, 1, it.Value())
	assert.True(t, it.Next())
	assert.Equal(t, 2, it.Value())
	assert.False(t, it.Next())

	it = Of(3)
	assert.True(t, it.Next())
	assert.Equal(t, 3, it.Value())
	assert.False(t, it.Next())

	it = OfEmpty[int]()
	assert.False(t, it.Next())

	it = OfOne(4)
	assert.True(t, it.Next())
	assert.Equal(t, 4, it.Value())
	assert.False(t, it.Next())
}

func TestOfMap(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2}
	it := OfMap(src)
	dst := map[string]int{}
	assert.True(t, it.Next())
	kv := it.Value()
	dst[kv.Key] = kv.Value
	assert.True(t, it.Next())
	kv = it.Value()
	dst[kv.Key] = kv.Value
	assert.False(t, it.Next())
	assert.Equal(t, src, dst)
}

func TestOfReader(t *testing.T) {
	it := OfReader(strings.NewReader("ab"))
	assert.True(t, it.Next())
	assert.Equal(t, byte('a'), it.Value())
	assert.True(t, it.Next())
	assert.Equal(t, byte('b'), it.Value())
	assert.False(t, it.Next())
}

func TestOfReaderAsRunes(t *testing.T) {
	it := OfReaderAsRunes(strings.NewReader("ab"))
	assert.True(t, it.Next())
	assert.Equal(t, 'a', it.Value())
	assert.True(t, it.Next())
	assert.Equal(t, 'b', it.Value())
	assert.False(t, it.Next())
}

func TestOfReaderAsLines(t *testing.T) {
	it := OfReaderAsLines(strings.NewReader("ab\ncd\ref\r\ngh"))
	assert.True(t, it.Next())
	assert.Equal(t, "ab", it.Value())
	assert.True(t, it.Next())
	assert.Equal(t, "cd", it.Value())
	assert.True(t, it.Next())
	assert.Equal(t, "ef", it.Value())
	assert.True(t, it.Next())
	assert.Equal(t, "gh", it.Value())
	assert.False(t, it.Next())
}

func TestConcat(t *testing.T) {
	
}
