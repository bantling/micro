package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/bantling/micro/go/funcs"
	"github.com/bantling/micro/go/util"
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

func TestNoValueIterGen(t *testing.T) {
	iter := NoValueIterGen[int]()

	val, haveIt := iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)

	val, haveIt = iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)
}

func TestSingleValueIterGen(t *testing.T) {
	iter := SingleValueIterGen(1)

	val, haveIt := iter()
	assert.Equal(t, 1, val)
	assert.True(t, haveIt)

	val, haveIt = iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)

	val, haveIt = iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)
}

func TestInfiniteIterGen(t *testing.T) {
	// Func to return (seed + 1, seed + 2, seed + 3, ...
	fn := func(prev int) int {
		return prev + 1
	}

	// Generate {1,2,3, ...}, which does not require a seed value
	iter := InfiniteIterGen(fn)

	for _, e := range []int{1, 2, 3, 4, 5, 6, 7} {
		v, h := iter()
		assert.Equal(t, e, v)
		assert.True(t, h)
	}

	// Generate {2,3,4, ...}, which requires a seed value of 1
	iter = InfiniteIterGen(fn, 1)

	for _, e := range []int{2, 3, 4, 5, 6, 7, 8} {
		v, h := iter()
		assert.Equal(t, e, v)
		assert.True(t, h)
	}

	// Generate {100,1,2,3, ...}, which requires a literal value of 100 and a seed value of 0
	iter = InfiniteIterGen(fn, 100, 0)

	for _, e := range []int{100, 1, 2, 3, 4, 5, 6} {
		v, h := iter()
		assert.Equal(t, e, v)
		assert.True(t, h)
	}

	// Generate {100,101,1,2,3, ...}, which requires a literal values of 100,101 and a seed value of 0
	iter = InfiniteIterGen(fn, 100, 101, 0)

	for _, e := range []int{100, 101, 1, 2, 3, 4, 5} {
		v, h := iter()
		assert.Equal(t, e, v)
		assert.True(t, h)
	}

	// Generate {100,101,2,3,4, ...}, which requires a literal values of 100,101 and a seed value of 1
	iter = InfiniteIterGen(fn, 100, 101, 1)

	for _, e := range []int{100, 101, 2, 3, 4, 5, 6} {
		v, h := iter()
		assert.Equal(t, e, v)
		assert.True(t, h)
	}
}

func TestFibonnaciIterGen(t *testing.T) {
	// Fibonnaci
	iter := FibonnaciIterGen()

	for _, e := range []int{1, 1, 2, 3, 5, 8, 13} {
		v, h := iter()
		assert.Equal(t, e, v)
		assert.True(t, h)
	}
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

	// non-eof error occurs after one byte

	err := fmt.Errorf("An error")
	src = util.NewErrorReader([]byte("a"), err)
	iter = ReaderIterGen(src)

	val, haveIt = iter()
	assert.Equal(t, byte('a'), val)
	assert.True(t, haveIt)

	funcs.TryTo(
		func() { iter() },
		func(e any) { assert.Equal(t, err, e) },
	)
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
			iter   = ReaderAsRunesIterGen(strings.NewReader(input))
			val    rune
			haveIt bool
		)

		for _, char := range []rune(input) {
			val, haveIt = iter()
			assert.Equal(t, char, val)
			assert.True(t, haveIt)
		}

		val, haveIt = iter()
		assert.Equal(t, rune(0), val)
		assert.False(t, haveIt)

		val, haveIt = iter()
		assert.Equal(t, rune(0), val)
		assert.False(t, haveIt)
	}

	// non-eof error occurs after one byte

	err := fmt.Errorf("An error")
	src = util.NewErrorReader([]byte("a"), err)
	iter = ReaderAsRunesIterGen(src)

	val, haveIt = iter()
	assert.Equal(t, 'a', val)
	assert.True(t, haveIt)

	funcs.TryTo(
		func() {
			iter()
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, err, e) },
	)

	// utf8 decoding error occurs after one byte

	src = strings.NewReader("a\x80")
	iter = ReaderAsRunesIterGen(src)

	val, haveIt = iter()
	assert.Equal(t, 'a', val)
	assert.True(t, haveIt)

	funcs.TryTo(
		func() {
			iter()
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, InvalidUTF8EncodingError, e) },
	)
}

func TestStringAsRunesIterGen(t *testing.T) {
	// nil
	var src string
	assert.Zero(t, src)
	iter := StringAsRunesIterGen(src)

	val, haveIt := iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)

	val, haveIt = iter()
	assert.Zero(t, val)
	assert.False(t, haveIt)

	// empty
	src = ""
	iter = StringAsRunesIterGen(src)

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
			iter   = StringAsRunesIterGen(input)
			val    rune
			haveIt bool
		)

		for _, char := range []rune(input) {
			val, haveIt = iter()
			assert.Equal(t, char, val)
			assert.True(t, haveIt)
		}

		val, haveIt = iter()
		assert.Equal(t, rune(0), val)
		assert.False(t, haveIt)

		val, haveIt = iter()
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

func TestStringAsLinesIterGen(t *testing.T) {
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
			iterFunc = StringAsLinesIterGen(input)
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
