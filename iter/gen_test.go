package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/bantling/micro/tuple"
	"github.com/bantling/micro/util"
	"github.com/stretchr/testify/assert"
)

func TestSliceIterGen_(t *testing.T) {
	// nil
	var slc []int
	assert.Nil(t, slc)
	iter := SliceIterGen(slc)

	i, err := iter()
	assert.Zero(t, i)
	assert.Equal(t, EOI, err)

	i, err = iter()
	assert.Zero(t, i)
	assert.Equal(t, EOI, err)

	// empty
	slc = []int{}
	iter = SliceIterGen(slc)

	i, err = iter()
	assert.Zero(t, i)
	assert.Equal(t, EOI, err)

	i, err = iter()
	assert.Zero(t, i)
	assert.Equal(t, EOI, err)

	// one element
	slc = []int{1}
	iter = SliceIterGen(slc)

	i, err = iter()
	assert.Equal(t, 1, i)
	assert.Nil(t, err)

	i, err = iter()
	assert.Zero(t, i)
	assert.Equal(t, EOI, err)

	i, err = iter()
	assert.Zero(t, i)
	assert.Equal(t, EOI, err)

	// two elements
	slc = []int{1, 2}
	iter = SliceIterGen(slc)

	i, err = iter()
	assert.Equal(t, 1, i)
	assert.Nil(t, err)

	i, err = iter()
	assert.Equal(t, 2, i)
	assert.Nil(t, err)

	i, err = iter()
	assert.Zero(t, i)
	assert.Equal(t, EOI, err)

	i, err = iter()
	assert.Zero(t, i)
	assert.Equal(t, EOI, err)
}

func TestMapIterGen_(t *testing.T) {
	// nil
	var src map[string]int
	assert.Nil(t, src)
	iter := MapIterGen(src)

	kv, err := iter()
	assert.Zero(t, kv)
	assert.Equal(t, EOI, err)

	kv, err = iter()
	assert.Zero(t, kv)
	assert.Equal(t, EOI, err)

	// empty
	src = map[string]int{}
	iter = MapIterGen(src)

	kv, err = iter()
	assert.Zero(t, kv)
	assert.Equal(t, EOI, err)

	kv, err = iter()
	assert.Zero(t, kv)
	assert.Equal(t, EOI, err)

	// one pair
	src = map[string]int{"a": 1}
	iter = MapIterGen(src)

	kv, err = iter()
	assert.Equal(t, tuple.Of2("a", 1), kv)
	assert.Nil(t, err)

	kv, err = iter()
	assert.Zero(t, kv)
	assert.Equal(t, EOI, err)

	kv, err = iter()
	assert.Zero(t, kv)
	assert.Equal(t, EOI, err)

	// two pairs
	src = map[string]int{"a": 1, "b": 2}
	dst := map[string]int{}
	iter = MapIterGen(src)

	kv, err = iter()
	assert.Nil(t, err)
	dst[kv.T] = kv.U

	kv, err = iter()
	assert.Nil(t, err)
	dst[kv.T] = kv.U

	assert.Equal(t, src, dst)

	kv, err = iter()
	assert.Zero(t, kv)
	assert.Equal(t, EOI, err)

	kv, err = iter()
	assert.Zero(t, kv)
	assert.Equal(t, EOI, err)
}

func TestNoValueIterGen_(t *testing.T) {
	iter := NoValueIterGen[int]()

	val, err := iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)
}

func TestSingleValueIterGen_(t *testing.T) {
	iter := SingleValueIterGen(1)

	val, err := iter()
	assert.Equal(t, 1, val)
	assert.Nil(t, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)
}

func TestInfiniteIterGen_(t *testing.T) {
	// Func to return (seed + 1, seed + 2, seed + 3, ...
	fn := func(prev int) int {
		return prev + 1
	}

	// Generate {1,2,3, ...}, which does not require a seed value
	iter := InfiniteIterGen(fn)

	for _, e := range []int{1, 2, 3, 4, 5, 6, 7} {
		v, err := iter()
		assert.Equal(t, e, v)
		assert.Nil(t, err)
	}

	// Generate {2,3,4, ...}, which requires a seed value of 1
	iter = InfiniteIterGen(fn, 1)

	for _, e := range []int{2, 3, 4, 5, 6, 7, 8} {
		v, err := iter()
		assert.Equal(t, e, v)
		assert.Nil(t, err)
	}

	// Generate {100,1,2,3, ...}, which requires a literal value of 100 and a seed value of 0
	iter = InfiniteIterGen(fn, 100, 0)

	for _, e := range []int{100, 1, 2, 3, 4, 5, 6} {
		v, err := iter()
		assert.Equal(t, e, v)
		assert.Nil(t, err)
	}

	// Generate {100,101,1,2,3, ...}, which requires literal values of 100,101 and a seed value of 0
	iter = InfiniteIterGen(fn, 100, 101, 0)

	for _, e := range []int{100, 101, 1, 2, 3, 4, 5} {
		v, err := iter()
		assert.Equal(t, e, v)
		assert.Nil(t, err)
	}

	// Generate {100,101,2,3,4, ...}, which requires literal values of 100,101 and a seed value of 1
	iter = InfiniteIterGen(fn, 100, 101, 1)

	for _, e := range []int{100, 101, 2, 3, 4, 5, 6} {
		v, err := iter()
		assert.Equal(t, e, v)
		assert.Nil(t, err)
	}
}

func TestFibonnaciIterGen_(t *testing.T) {
	// Fibonnaci
	iter := FibonnaciIterGen()

	for _, e := range []int{1, 1, 2, 3, 5, 8, 13} {
		v, err := iter()
		assert.Equal(t, e, v)
		assert.Nil(t, err)
	}
}

func TestReaderIterGen_(t *testing.T) {
	// nil
	var src io.Reader
	assert.Zero(t, src)
	iter := ReaderIterGen(src)

	val, err := iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	// empty
	src = strings.NewReader("")
	iter = ReaderIterGen(src)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	// one byte
	src = strings.NewReader("a")
	iter = ReaderIterGen(src)

	val, err = iter()
	assert.Equal(t, byte('a'), val)
	assert.Nil(t, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	// non-eof-eoi error occurs after one byte
	anErr := fmt.Errorf("An error")
	src = util.NewErrorReader([]byte("a"), anErr)
	iter = ReaderIterGen(src)

	val, err = iter()
	assert.Equal(t, byte('a'), val)
	assert.Nil(t, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, anErr, err)
}

func TestReaderAsRunesIterGen_(t *testing.T) {
	// nil
	var src io.Reader
	assert.Zero(t, src)
	iter := ReaderAsRunesIterGen(src)

	val, err := iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	// empty
	src = strings.NewReader("")
	iter = ReaderAsRunesIterGen(src)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	inputs := []string{
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
		"\U00010348",
		"\U00010348\U00010348",
		"\U00010348\U00010348\U00010348",
		"\U00010348\U00010348\U00010348\U00010348",
	}

	for _, input := range inputs {
		var (
			iter = ReaderAsRunesIterGen(strings.NewReader(input))
			val  rune
			err  error
		)

		for _, char := range []rune(input) {
			val, err = iter()
			assert.Equal(t, fmt.Sprintf("%b", char), fmt.Sprintf("%b", val))
			assert.Nil(t, err)
		}

		val, err = iter()
		assert.Zero(t, val)
		assert.Equal(t, EOI, err)

		val, err = iter()
		assert.Zero(t, val)
		assert.Equal(t, EOI, err)
	}

	// non-eof error occurs after one byte
	anErr := fmt.Errorf("An error")
	src = util.NewErrorReader([]byte("a"), anErr)
	iter = ReaderAsRunesIterGen(src)

	val, err = iter()
	assert.Equal(t, 'a', val)
	assert.Nil(t, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, anErr, err)

	// decoding error occurs on first byte
	src = strings.NewReader("a\x80")
	iter = ReaderAsRunesIterGen(src)

	val, err = iter()
	assert.Equal(t, 'a', val)
	assert.Nil(t, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, InvalidUTF8EncodingError, err)

	// decoding error when two bytes required, only one provided
	// 110 00000 = C0
	src = strings.NewReader("\xc0")
	iter = ReaderAsRunesIterGen(src)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, InvalidUTF8EncodingError, err)

	// decoding error if an extra byte does not begin with 10
	// 110 00000, 11 000000 = C0 C0
	src = strings.NewReader("\xc0\xc0")
	iter = ReaderAsRunesIterGen(src)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, InvalidUTF8EncodingError, err)

	// decoding error if value > max allowed 10FFFF
	// 11FFFF = 11110 100, 10 011111, 10 111111, 10 111111
	src = strings.NewReader("\xf4\x9f\xbf\xbf")
	iter = ReaderAsRunesIterGen(src)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, InvalidUTF8EncodingError, err)
}

func TestStringAsRunesIterGen_(t *testing.T) {
	// nil
	var src string
	assert.Zero(t, src)
	iter := StringAsRunesIterGen(src)

	val, err := iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	// empty
	src = ""
	iter = StringAsRunesIterGen(src)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

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
		"\U00010348",
		"\U00010348\U00010348",
		"\U00010348\U00010348\U00010348",
		"\U00010348\U00010348\U00010348\U00010348",
	}

	for _, input := range inputs {
		var (
			iter = StringAsRunesIterGen(input)
			val  rune
			err  error
		)

		for _, char := range []rune(input) {
			val, err = iter()
			assert.Equal(t, char, val)
			assert.Nil(t, err)
		}

		val, err = iter()
		assert.Zero(t, val)
		assert.Equal(t, EOI, err)

		val, err = iter()
		assert.Zero(t, val)
		assert.Equal(t, EOI, err)
	}

	// utf8 decoding error occurs after one byte
	src = "a\x80"
	iter = StringAsRunesIterGen(src)

	val, err = iter()
	assert.Equal(t, 'a', val)
	assert.Nil(t, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, InvalidUTF8EncodingError, err)
}

func TestReaderAsLinesIterGen_(t *testing.T) {
	var (
		inputs = []string{
			"oneline",
			"two\rline cr",
			"two\nline lf",
			"two\r\nline crlf",
		}
		linesRegex, _ = regexp.Compile("\r\n|\r|\n")
	)

	for _, input := range inputs {
		var (
			iter  = ReaderAsLinesIterGen(strings.NewReader(input))
			lines = linesRegex.Split(input, -1)
			val   string
			err   error
		)

		for _, line := range lines {
			val, err = iter()
			assert.Equal(t, line, val)
			assert.Equal(t, nil, err)
		}

		val, err = iter()
		assert.Zero(t, val)
		assert.Equal(t, EOI, err)

		val, err = iter()
		assert.Zero(t, val)
		assert.Equal(t, EOI, err)
	}

	iter := ReaderAsLinesIterGen(strings.NewReader(""))
	val, err := iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	iter = ReaderAsLinesIterGen(strings.NewReader("a\x80"))
	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, InvalidUTF8EncodingError, err)
}

func TestStringAsLinesIterGen_(t *testing.T) {
	var (
		inputs = []string{
			"oneline",
			"two\rline cr",
			"two\nline lf",
			"two\r\nline crlf",
		}
		linesRegex, _ = regexp.Compile("\r\n|\r|\n")
	)

	for _, input := range inputs {
		var (
			iter  = StringAsLinesIterGen(input)
			lines = linesRegex.Split(input, -1)
			val   string
			err   error
		)

		for _, line := range lines {
			val, err = iter()
			assert.Equal(t, line, val)
			assert.Nil(t, err)
		}

		val, err = iter()
		assert.Zero(t, val)
		assert.Equal(t, EOI, err)

		val, err = iter()
		assert.Zero(t, val)
		assert.Equal(t, EOI, err)
	}

	iter := StringAsLinesIterGen("")
	val, err := iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	iter = StringAsLinesIterGen("a\x80")
	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, InvalidUTF8EncodingError, err)
}

func TestConcatIterGen_(t *testing.T) {
	iter := ConcatIterGen(
		[]Iter[string]{
			NewIter(NoValueIterGen[string]()),
			NewIter(SliceIterGen([]string{"foo", "bar"})),
			NewIter(NoValueIterGen[string]()),
			NewIter(SingleValueIterGen("baz")),
			NewIter(NoValueIterGen[string]()),
		},
	)

	val, err := iter()
	assert.Equal(t, "foo", val)
	assert.Nil(t, err)

	val, err = iter()
	assert.Equal(t, "bar", val)
	assert.Nil(t, err)

	val, err = iter()
	assert.Equal(t, "baz", val)
	assert.Nil(t, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, EOI, err)

	anErr := fmt.Errorf("An error")
	iter = ConcatIterGen(
		[]Iter[string]{
			Of("1"),
			SetError(OfEmpty[string](), anErr),
			Of("2"),
		},
	)

	val, err = iter()
	assert.Equal(t, "1", val)
	assert.Nil(t, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, anErr, err)

	val, err = iter()
	assert.Zero(t, val)
	assert.Equal(t, anErr, err)
}
