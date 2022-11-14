package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bantling/micro/go/funcs"
	"github.com/bantling/micro/go/util"
	"github.com/stretchr/testify/assert"
)

func TestNewIter(t *testing.T) {
	it := NewIter(SliceIterGen[int]([]int{1, 2}))
	assert.Equal(t, util.Of2Error(1, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(2, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	it = NewIter(FibonnaciIterGen())
	assert.Equal(t, util.Of2Error(1, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(1, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(2, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(3, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(5, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(8, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(13, nil), Maybe(it))

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
}

func TestOf(t *testing.T) {
	it := Of(3, 4)
	assert.Equal(t, util.Of2Error(3, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(4, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
}

func TestOfEmpty(t *testing.T) {
	it := OfEmpty[int]()
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
}

func TestOfOne(t *testing.T) {
	it := OfOne(5)
	assert.Equal(t, util.Of2Error(5, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
}

func TestOfMap(t *testing.T) {
	var (
		src = map[string]int{"a": 1, "b": 2}
		it  = OfMap(src)
		dst = map[string]int{}
	)

	kv, err := it.Next()
	assert.Nil(t, err)
	dst[kv.T] = kv.U

	kv, err = it.Next()
	assert.Nil(t, err)
	dst[kv.T] = kv.U

	assert.Equal(t, util.Of2Error(util.Of2("", 0), EOI), Maybe(it))
	assert.Equal(t, src, dst)
}

func TestOfReader(t *testing.T) {
	it := OfReader(strings.NewReader("ab"))
	assert.Equal(t, util.Of2Error(byte('a'), nil), Maybe(it))
	assert.Equal(t, util.Of2Error(byte('b'), nil), Maybe(it))
	assert.Equal(t, util.Of2Error(byte(0), EOI), Maybe(it))
}

func TestOfReaderAsRunes(t *testing.T) {
	it := OfReaderAsRunes(strings.NewReader("ab"))
	assert.Equal(t, util.Of2Error('a', nil), Maybe(it))
	assert.Equal(t, util.Of2Error('b', nil), Maybe(it))
	assert.Equal(t, util.Of2Error('\x00', EOI), Maybe(it))
}

func TestOfStringAsRunes(t *testing.T) {
	it := OfStringAsRunes("ab")
	assert.Equal(t, util.Of2Error('a', nil), Maybe(it))
	assert.Equal(t, util.Of2Error('b', nil), Maybe(it))
	assert.Equal(t, util.Of2Error('\x00', EOI), Maybe(it))
}

func TestOfReaderAsLines(t *testing.T) {
	it := OfReaderAsLines(strings.NewReader("ab\ncd\ref\r\ngh"))
	assert.Equal(t, util.Of2Error("ab", nil), Maybe(it))
	assert.Equal(t, util.Of2Error("cd", nil), Maybe(it))
	assert.Equal(t, util.Of2Error("ef", nil), Maybe(it))
	assert.Equal(t, util.Of2Error("gh", nil), Maybe(it))
	assert.Equal(t, util.Of2Error("", EOI), Maybe(it))
}

func TestOfStringAsLines(t *testing.T) {
	it := OfStringAsLines("ab\ncd\ref\r\ngh")
	assert.Equal(t, util.Of2Error("ab", nil), Maybe(it))
	assert.Equal(t, util.Of2Error("cd", nil), Maybe(it))
	assert.Equal(t, util.Of2Error("ef", nil), Maybe(it))
	assert.Equal(t, util.Of2Error("gh", nil), Maybe(it))
	assert.Equal(t, util.Of2Error("", EOI), Maybe(it))
}

func TestConcat(t *testing.T) {
	it := Concat(Of(1), Of(2, 3))
	assert.Equal(t, util.Of2Error(1, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(2, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(3, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
}

func TestUnread(t *testing.T) {
	// Unread before next
	it := OfEmpty[int]()
	it.Unread(1)
	assert.Equal(t, util.Of2Error(1, nil), Maybe(it))

	// Unread after next
	it.Unread(2)
	assert.Equal(t, util.Of2Error(2, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	// Unread two values to test order, after next returns EOI
	it.Unread(3)
	it.Unread(4)
	assert.Equal(t, util.Of2Error(4, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(3, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
}

func TestMaybe(t *testing.T) {
	it := OfEmpty[int]()
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	it = OfOne(1)
	assert.Equal(t, util.Of2Error(1, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
}

func TestSetError(t *testing.T) {
	anErr := fmt.Errorf("An err")
	it := SetError(OfEmpty[int](), anErr)
	assert.Equal(t, util.Of2Error(0, anErr), Maybe(it))
	assert.Equal(t, util.Of2Error(0, anErr), Maybe(it))

	it = SetError(OfOne(1), anErr)
	assert.Equal(t, util.Of2Error(1, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, anErr), Maybe(it))
	assert.Equal(t, util.Of2Error(0, anErr), Maybe(it))
}
