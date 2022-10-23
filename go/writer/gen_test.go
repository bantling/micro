package writer

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bantling/micro/go/conv"
	"github.com/bantling/micro/go/util"
	"github.com/stretchr/testify/assert"
)

func TestSliceWriterGen(t *testing.T) {
	var (
		slc = make([]int, 0, 1)
		w   = SliceWriterGen(&slc)
		p   = fmt.Sprintf("%p", slc) // address of empty slice
	)

	for i := 1; i <= 100; i++ {
		assert.Nil(t, w(i))
		assert.Equal(t, i, len(slc))
		assert.Equal(t, i, slc[i-1])
	}

	assert.NotEqual(t, p, fmt.Sprintf("%p", slc)) // address has changed
}

func TestMapWriterGen(t *testing.T) {
	var (
		m = map[int]string{}
		w = MapWriterGen(m)
	)

	for i := 1; i <= 100; i++ {
		s := conv.IntToString(i)
		assert.Nil(t, w(util.KVOf(i, s)))
		assert.Equal(t, i, len(m))
		assert.Equal(t, s, m[i])
	}
}

func TestIOWriterGen(t *testing.T) {
	// Write alphabet to a normal writer
	var (
		sb strings.Builder
		w  = IOWriterGen(&sb)
	)

	// 'A' through 'Z'
	for i := 0; i <= 25; i++ {
		assert.Nil(t, w(byte(i+0x41)))
		assert.Equal(t, i+1, sb.Len())
	}
	assert.Equal(t, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", sb.String())

	// Write partial alphabet to an error writer
	var (
		err = fmt.Errorf("An error")
		ew  = util.NewErrorWriter(5, err)
	)
	w = IOWriterGen(ew)

	for i := 0; i <= 4; i++ {
		assert.Nil(t, w(byte(i+0x41)))
		assert.Equal(t, i+1, len(ew.Output()))
	}
	assert.Equal(t, "ABCDE", string(ew.Output()))

	assert.Equal(t, err, w(0x46))
	assert.Equal(t, "ABCDE", string(ew.Output()))

	// Write 0 bytes with no error
	ew = util.NewErrorWriter(-1, err)
	w = IOWriterGen(ew)

	assert.Equal(t, fmt.Errorf(errIOByteWriterMsg, byte(0x47)), w(0x47))
}

func TestIOWriterAsRunesGen(t *testing.T) {
	var (
		sb    strings.Builder
		w     = IOWriterAsRunesGen(&sb)
		runes = []rune{
			'a',          // 1 byte
			'\u00e0',     // 2 bytes
			'\u1e01',     // 3 bytes
			'\U00010348', // 4 bytes
		}
	)

	lens := 0
	for _, r := range runes {
		assert.Nil(t, w(r))
		lens += len(string(r))
		assert.Equal(t, lens, sb.Len())
	}
	assert.Equal(t, "a\u00e0\u1e01\U00010348", sb.String())

	// Write to an error writer
	var (
		err = fmt.Errorf("An error")
		ew  = util.NewErrorWriter(len(string(runes)), err)
	)
	w = IOWriterAsRunesGen(ew)

	for _, r := range runes {
		assert.Nil(t, w(r))
	}
	assert.Equal(t, "a\u00e0\u1e01\U00010348", string(ew.Output()))

	assert.Equal(t, err, w(0x98))
	assert.Equal(t, "a\u00e0\u1e01\U00010348", string(ew.Output()))

	// Write 0 runes with no error
	ew = util.NewErrorWriter(-1, err)
	w = IOWriterAsRunesGen(ew)

	assert.Equal(t, fmt.Errorf(errIORuneWriterMsg, rune(0x99)), w(0x99))
}

func TestIOWriterAsStringsGen(t *testing.T) {
	var (
		sb      strings.Builder
		w       = IOWriterAsStringsGen(&sb)
		strings = []string{
			"a",          // 1 byte
			"\u00e0",     // 2 bytes
			"\u1e01",     // 3 bytes
			"\U00010348", // 4 bytes
		}
	)

	lens := 0
	for _, s := range strings {
		assert.Nil(t, w(s))
		lens += len(string(s))
		assert.Equal(t, lens, sb.Len())
	}
	assert.Equal(t, "a\u00e0\u1e01\U00010348", sb.String())

	// Write to an error writer
	var (
		err = fmt.Errorf("An error")
		ew  = util.NewErrorWriter(10, err)
	)
	w = IOWriterAsStringsGen(ew)

	for _, s := range strings {
		assert.Nil(t, w(s))
	}
	assert.Equal(t, "a\u00e0\u1e01\U00010348", string(ew.Output()))

	assert.Equal(t, err, w("b"))
	assert.Equal(t, "a\u00e0\u1e01\U00010348", string(ew.Output()))

	// Write 0 runes with no error
	ew = util.NewErrorWriter(-1, err)
	w = IOWriterAsStringsGen(ew)

	assert.Equal(t, fmt.Errorf(errIOStringWriterMsg, 0, 1, 1), w("c"))
}

func TestIOWriterAsLinesGen(t *testing.T) {
	var (
		sb      strings.Builder
		w       = IOWriterAsLinesGen(&sb)
		strings = []string{
			"a",          // 1 byte
			"\u00e0",     // 2 bytes
			"\u1e01",     // 3 bytes
			"\U00010348", // 4 bytes
		}
	)

	lens := 0
	for _, r := range strings {
		assert.Nil(t, w(r))
		lens += len(string(r)) + len(osEOLSequence)
		assert.Equal(t, lens, sb.Len())
	}
	assert.Equal(
		t,
		"a"+osEOLSequence+
			"\u00e0"+osEOLSequence+
			"\u1e01"+osEOLSequence+
			"\U00010348"+osEOLSequence,
		sb.String(),
	)

	// Write to an error writer
	var (
		err = fmt.Errorf("An error")
		ew  = util.NewErrorWriter(10+len(strings)*len(osEOLSequence), err)
	)
	w = IOWriterAsLinesGen(ew)

	for _, s := range strings {
		assert.Nil(t, w(s))
	}
	assert.Equal(
		t,
		"a"+osEOLSequence+
			"\u00e0"+osEOLSequence+
			"\u1e01"+osEOLSequence+
			"\U00010348"+osEOLSequence,
		string(ew.Output()),
	)

	assert.Equal(t, err, w("b"))
	assert.Equal(
		t,
		"a"+osEOLSequence+
			"\u00e0"+osEOLSequence+
			"\u1e01"+osEOLSequence+
			"\U00010348"+osEOLSequence,
		string(ew.Output()),
	)

	// Write 0 runes with no error
	ew = util.NewErrorWriter(-1, err)
	w = IOWriterAsLinesGen(ew)

	assert.Equal(t, fmt.Errorf(errIOStringWriterMsg, 0, 1+len(osEOLSequence), 1+len(osEOLSequence)), w("c"))
}
