package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"strings"
	"testing"

	"github.com/bantling/micro/go/funcs"
	"github.com/stretchr/testify/assert"
)

func TestNewIter(t *testing.T) {
	it := NewIter(SliceIterGen[int]([]int{1, 2}))
	assert.True(t, it.Next())
	assert.Equal(t, 1, it.Value())
	assert.True(t, it.Next())
	assert.Equal(t, 2, it.Value())
	assert.False(t, it.Next())
	assert.False(t, it.Next())

	it = Of(3)
	assert.True(t, it.Next())
	assert.Equal(t, 3, it.Value())
	assert.False(t, it.Next())
	assert.False(t, it.Next())

	it = OfEmpty[int]()
	assert.False(t, it.Next())
	assert.False(t, it.Next())

	it = OfOne(4)
	assert.True(t, it.Next())
	assert.Equal(t, 4, it.Value())
	assert.False(t, it.Next())
	assert.False(t, it.Next())

	it = NewIter(FibonnaciIterGen())
	assert.Equal(t, 1, it.Must())
	assert.Equal(t, 1, it.Must())
	assert.Equal(t, 2, it.Must())
	assert.Equal(t, 3, it.Must())
	assert.Equal(t, 5, it.Must())
	assert.Equal(t, 8, it.Must())
	assert.Equal(t, 13, it.Must())
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

func TestOfStringAsRunes(t *testing.T) {
	it := OfStringAsRunes("ab")
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

func TestOfStringAsLines(t *testing.T) {
	it := OfStringAsLines("ab\ncd\ref\r\ngh")
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
	it := Concat(Of(1), Of(2, 3))
	assert.True(t, it.Next())
	assert.Equal(t, 1, it.Value())
	assert.True(t, it.Next())
	assert.Equal(t, 2, it.Value())
	assert.True(t, it.Next())
	assert.Equal(t, 3, it.Value())
	assert.False(t, it.Next())
}

func TestUnread(t *testing.T) {
	// Unread without next returning false
	it := OfEmpty[int]()
	it.Unread(1)
	assert.True(t, it.Next())
	assert.Equal(t, 1, it.Value())

	// Unread with next returning false
	it.Unread(2)
	assert.True(t, it.Next())
	assert.Equal(t, 2, it.Value())
	assert.False(t, it.Next())

	// Unread two values to test order, after next returns false
	it.Unread(3)
	it.Unread(4)
	assert.True(t, it.Next())
	assert.Equal(t, 3, it.Value())
	assert.True(t, it.Next())
	assert.Equal(t, 4, it.Value())
	assert.False(t, it.Next())
}

func TestNextValue(t *testing.T) {
	it := OfEmpty[int]()
	val, haveIt := it.NextValue()
	assert.False(t, haveIt)
	assert.Equal(t, 0, val)

	val, haveIt = it.NextValue()
	assert.False(t, haveIt)
	assert.Equal(t, 0, val)

	it = Of(1)
	val, haveIt = it.NextValue()
	assert.True(t, haveIt)
	assert.Equal(t, 1, val)

	val, haveIt = it.NextValue()
	assert.False(t, haveIt)
	assert.Equal(t, 0, val)
}

func TestMust(t *testing.T) {
	funcs.TryTo(
		func() {
			OfEmpty[int]().Must()
			assert.Fail(t, "Must die")
		},
		func(err any) {
			assert.Equal(t, errNoMoreValues, err)
		},
	)

	it := Of(1)
	assert.Equal(t, 1, it.Must())
	assert.False(t, it.Next())
}

func TestFailure(t *testing.T) {
	// Nil iter func
	funcs.TryTo(
		func() {
			NewIter[int](nil)
			assert.Fail(t, "Must die")
		},
		func(err any) {
			assert.Equal(t, errNewIterNeedsIterator, err)
		},
	)

	// Call Next twice without calling Value when there is a value to read
	funcs.TryTo(
		func() {
			it := Of(1)
			it.Next()
			it.Next()
			assert.Fail(t, "Must die")
		},
		func(err any) {
			assert.Equal(t, errValueExpected, err)
		},
	)

	// Call Value before Next
	funcs.TryTo(
		func() {
			it := Of(1)
			it.Value()
			assert.Fail(t, "Must die")
		},
		func(err any) {
			assert.Equal(t, errNextExpected, err)
		},
	)

	// Call Value twice
	funcs.TryTo(
		func() {
			it := Of(1)
			it.Next()
			it.Value()
			it.Value()
			assert.Fail(t, "Must die")
		},
		func(err any) {
			assert.Equal(t, errNextExpected, err)
		},
	)
}
