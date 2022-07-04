// SPDX-License-Identifier: Apache-2.0

package iter

import (
	//	"bytes"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/bantling/micro/go/funcs"
	"github.com/stretchr/testify/assert"
)

func TestSliceIterGen(t *testing.T) {
	// nil
	var slc []int
	assert.Nil(t, slc)
	iter := SliceIterGen(slc)

	i, haveIt := iter()
	assert.Zero(t, i)
	assert.False(t, haveIt)

	i, haveIt = iter()
	assert.Zero(t, i)
	assert.False(t, haveIt)

	// empty
	slc = []int{}
	iter = SliceIterGen(slc)

	i, haveIt = iter()
	assert.Zero(t, i)
	assert.False(t, haveIt)

	i, haveIt = iter()
	assert.Zero(t, i)
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
	assert.Zero(t, i)
	assert.False(t, haveIt)

	i, haveIt = iter()
	assert.Zero(t, i)
	assert.False(t, haveIt)
}

func TestMapIterGen(t *testing.T) {
	// nil
	var src map[string]int
	assert.Nil(t, src)
	iter := MapIterGen(src)

	kv, haveIt := iter()
	assert.Zero(t, kv)
	assert.False(t, haveIt)

	kv, haveIt = iter()
	assert.Zero(t, kv)
	assert.False(t, haveIt)

	// empty
	src = map[string]int{}
	iter = MapIterGen(src)

	kv, haveIt = iter()
	assert.Zero(t, kv)
	assert.False(t, haveIt)

	kv, haveIt = iter()
	assert.Zero(t, kv)
	assert.False(t, haveIt)

	// one pair
	src = map[string]int{"a": 1}
	iter = MapIterGen(src)

	kv, haveIt = iter()
	assert.Equal(t, KeyValue[string, int]{"a", 1}, kv)
	assert.True(t, haveIt)

	kv, haveIt = iter()
	assert.Zero(t, kv)
	assert.False(t, haveIt)

	kv, haveIt = iter()
	assert.Zero(t, kv)
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

	kv, haveIt = iter()
	assert.Zero(t, kv)
	assert.False(t, haveIt)

	kv, haveIt = iter()
	assert.Zero(t, kv)
	assert.False(t, haveIt)
}

func TestReaderIterGen(t *testing.T) {
	// nil
	var src io.Reader
	assert.Zero(t, src)
	iter := ReaderIterGen(src)

	val, haveIt := iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)

	val, haveIt = iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)

	// empty
	src = strings.NewReader("")
	iter = ReaderIterGen(src)

	val, haveIt = iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)

	val, haveIt = iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)

	// one byte
	src = strings.NewReader("a")
	iter = ReaderIterGen(src)

	val, haveIt = iter()
	assert.Equal(t, byte('a'), val)
	assert.True(t, haveIt)

	val, haveIt = iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)

	val, haveIt = iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)
}

func TestReaderAsRunesIterGen(t *testing.T) {
	// nil
	var src io.Reader
	assert.Zero(t, src)
	iter := ReaderAsRunesIterGen(src)

	val, haveIt := iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)

	val, haveIt = iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)

	// empty
	src = strings.NewReader("")
	iter = ReaderAsRunesIterGen(src)

	val, haveIt = iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)

	val, haveIt = iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)
	inputs := []string{
		"",
		// 1 byte UTF8
		"a",
		"ab",
		"abc",
		"abcd",
		"abcde",
		"abcdef",
		"abcdefg",
		"abcdefgh",
		"abcdefghi",
		// 2 byte UTF8
		"\u00e0",
		"\u00e0\u00e0",
		"\u00e0\u00e0\u00e0",
		"\u00e0\u00e0\u00e0\u00e0",
		// 3 byte UTF8
		"\u1e01",
		"\u1e01\u1e01",
		"\u1e01\u1e01\u1e01",
		"\u1e01\u1e01\u1e01\u1e01",
		// 4 bytes UTF8
		"\u10348",
		"\u10348\u10348",
		"\u10348\u10348\u10348",
		"\u10348\u10348\u10348\u10348",
	}

	for _, input := range inputs {
		var (
			iterFunc = ReaderAsRunesIterGen(strings.NewReader(input))
			val      rune
			haveIt   bool
		)

		for _, char := range []rune(input) {
			val, haveIt = iterFunc()
			assert.Equal(t, char, val)
			assert.True(t, haveIt)
		}

		val, haveIt = iterFunc()
		assert.Equal(t, rune(0), val)
		assert.False(t, haveIt)

		val, haveIt = iterFunc()
		assert.Equal(t, rune(0), val)
		assert.False(t, haveIt)
	}
}

func TestReaderAsLinesIterGen(t *testing.T) {
	var (
		inputs = []string{
			"",
			"oneline",
			"two\rline cr",
			"two\nline lf",
			"two\r\nline crlf",
		}
		linesRegex, _ = regexp.Compile("\r\n|\r|\n")
	)

	for _, input := range inputs {
		var (
			iterFunc = ReaderAsLinesIterGen(strings.NewReader(input))
			lines    = linesRegex.Split(input, -1)
			val      string
			haveIt   bool
		)

		for _, line := range lines {
			val, haveIt = iterFunc()
			assert.Equal(t, line, val)
			assert.Equal(t, input != "", haveIt)
		}

		val, haveIt = iterFunc()
		assert.Equal(t, "", val)
		assert.False(t, haveIt)

		val, haveIt = iterFunc()
		assert.Equal(t, "", val)
		assert.False(t, haveIt)
	}
}

func TestFlattenSlice(t *testing.T) {
	assert.Equal(t, []int{}, FlattenSlice[int](nil))

	// Check that one dimensional slice is returned as (same address)
	oneDim := []int{}
	assert.Equal(t, fmt.Sprintf("%p", oneDim), fmt.Sprintf("%p", FlattenSlice[int](oneDim)))

	assert.Equal(t, []int{1, 2, 3, 4}, FlattenSlice[int]([][]int{{1, 2}, {3, 4}}))

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, FlattenSlice[int]([][][]int{{{1, 2}, {3, 4}}, {{5}}, {{6}}}))

	// Die if a value that is not a slice is passed
	funcs.TryTo(
		func() {
			FlattenSlice[int](0)
			assert.Fail(t, "Must die")
		},
		func(err any) {
			assert.Equal(t, fmt.Errorf(flattenSliceArgNotSliceMsg, 0), err)
		},
	)

	// Die if expecting a []int but passed a []string
	funcs.TryTo(
		func() {
			FlattenSlice[int]([]string{})
			assert.Fail(t, "Must die")
		},
		func(err any) {
			assert.Equal(t, fmt.Errorf(flattenSliceArgNotTMsg, reflect.TypeOf(0), reflect.TypeOf("")), err)
		},
	)
}
