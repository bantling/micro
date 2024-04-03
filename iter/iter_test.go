package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/tuple"
	"github.com/bantling/micro/union"
	"github.com/stretchr/testify/assert"
)

func TestOfIter_(t *testing.T) {
	it := OfIter(SliceIterGen[int]([]int{1, 2}))
	assert.Equal(t, union.OfResult(1), Maybe(it))
	assert.Equal(t, union.OfResult(2), Maybe(it))
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))

	it = OfIter(FibonnaciIterGen())
	assert.Equal(t, union.OfResult(1), Maybe(it))
	assert.Equal(t, union.OfResult(1), Maybe(it))
	assert.Equal(t, union.OfResult(2), Maybe(it))
	assert.Equal(t, union.OfResult(3), Maybe(it))
	assert.Equal(t, union.OfResult(5), Maybe(it))
	assert.Equal(t, union.OfResult(8), Maybe(it))
	assert.Equal(t, union.OfResult(13), Maybe(it))

	// Nil iter func
	funcs.TryTo(
		func() {
			OfIter[int](nil)
			assert.Fail(t, "Must die")
		},
		func(err any) {
			assert.Equal(t, errOfIterNeedsIterator, err)
		},
	)
}

func TestOf_(t *testing.T) {
	it := Of(3, 4)
	assert.Equal(t, union.OfResult(3), Maybe(it))
	assert.Equal(t, union.OfResult(4), Maybe(it))
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))
}

func TestOfEmpty_(t *testing.T) {
	it := OfEmpty[int]()
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))
}

func TestOfOne_(t *testing.T) {
	it := OfOne(5)
	assert.Equal(t, union.OfResult(5), Maybe(it))
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))
}

func TestOfSlice_(t *testing.T) {
	it := OfSlice([]int{3, 4})
	assert.Equal(t, union.OfResult(3), Maybe(it))
	assert.Equal(t, union.OfResult(4), Maybe(it))
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))
}

func TestOfMap_(t *testing.T) {
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

	assert.Equal(t, union.OfError[tuple.Two[string, int]](EOI), Maybe(it))
	assert.Equal(t, src, dst)
}

func TestOfReader_(t *testing.T) {
	it := OfReader(strings.NewReader("ab"))
	assert.Equal(t, union.OfResult(byte('a')), Maybe(it))
	assert.Equal(t, union.OfResult(byte('b')), Maybe(it))
	assert.Equal(t, union.OfError[byte](EOI), Maybe(it))
}

func TestOfReaderAsRunes_(t *testing.T) {
	it := OfReaderAsRunes(strings.NewReader("ab"))
	assert.Equal(t, union.OfResult('a'), Maybe(it))
	assert.Equal(t, union.OfResult('b'), Maybe(it))
	assert.Equal(t, union.OfError[rune](EOI), Maybe(it))
}

func TestOfStringAsRunes_(t *testing.T) {
	it := OfStringAsRunes("ab")
	assert.Equal(t, union.OfResult('a'), Maybe(it))
	assert.Equal(t, union.OfResult('b'), Maybe(it))
	assert.Equal(t, union.OfError[rune](EOI), Maybe(it))
}

func TestOfReaderAsLines_(t *testing.T) {
	it := OfReaderAsLines(strings.NewReader("ab\ncd\ref\r\ngh"))
	assert.Equal(t, union.OfResult("ab"), Maybe(it))
	assert.Equal(t, union.OfResult("cd"), Maybe(it))
	assert.Equal(t, union.OfResult("ef"), Maybe(it))
	assert.Equal(t, union.OfResult("gh"), Maybe(it))
	assert.Equal(t, union.OfError[string](EOI), Maybe(it))
}

func TestOfStringAsLines_(t *testing.T) {
	it := OfStringAsLines("ab\ncd\ref\r\ngh")
	assert.Equal(t, union.OfResult("ab"), Maybe(it))
	assert.Equal(t, union.OfResult("cd"), Maybe(it))
	assert.Equal(t, union.OfResult("ef"), Maybe(it))
	assert.Equal(t, union.OfResult("gh"), Maybe(it))
	assert.Equal(t, union.OfError[string](EOI), Maybe(it))
}

func TestConcat_(t *testing.T) {
	it := Concat(Of(1), Of(2, 3))
	assert.Equal(t, union.OfResult(1), Maybe(it))
	assert.Equal(t, union.OfResult(2), Maybe(it))
	assert.Equal(t, union.OfResult(3), Maybe(it))
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))
}

func TestUnread_(t *testing.T) {
	// Unread before next
	it := OfEmpty[int]()
	it.Unread(1)
	assert.Equal(t, union.OfResult(1), Maybe(it))

	// Unread after next
	it.Unread(2)
	assert.Equal(t, union.OfResult(2), Maybe(it))
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))

	// Unread two values to test order, after next returns EOI
	it.Unread(3)
	it.Unread(4)
	assert.Equal(t, union.OfResult(4), Maybe(it))
	assert.Equal(t, union.OfResult(3), Maybe(it))
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))
}

func TestMaybe_(t *testing.T) {
	it := OfEmpty[int]()
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))

	it = OfOne(1)
	assert.Equal(t, union.OfResult(1), Maybe(it))
	assert.Equal(t, union.OfError[int](EOI), Maybe(it))
}

func TestSetError_(t *testing.T) {
	anErr := fmt.Errorf("An err")
	it := SetError(OfEmpty[int](), anErr)
	assert.Equal(t, union.OfError[int](anErr), Maybe(it))
	assert.Equal(t, union.OfError[int](anErr), Maybe(it))

	it = SetError(OfOne(1), anErr)
	assert.Equal(t, union.OfResult(1), Maybe(it))
	assert.Equal(t, union.OfError[int](anErr), Maybe(it))
	assert.Equal(t, union.OfError[int](anErr), Maybe(it))
}
