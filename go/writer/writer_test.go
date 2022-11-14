package writer

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bantling/micro/go/funcs"
	"github.com/bantling/micro/go/util"
	"github.com/stretchr/testify/assert"
)

func TestNewWriter(t *testing.T) {
	var (
		slc []int
		w   = NewWriter(SliceWriterGen(&slc))
	)
	assert.NotNil(t, w)

	w.Write(2)
	assert.Equal(t, []int{2}, slc)

	funcs.TryTo(
		func() {
			NewWriter[any](nil)
			assert.Fail(t, "Must die")
		},
		func(e any) { assert.Equal(t, errNewWriterNeedsFn, e) },
	)
}

func TestOfSliceWriter(t *testing.T) {
	var (
		slc []int
		w   = OfSliceWriter(&slc)
	)
	assert.NotNil(t, w)

	w.Write(2)
	assert.Equal(t, []int{2}, slc)
}

func TestOfMapWriter(t *testing.T) {
	var (
		m = map[int]string{}
		w = OfMapWriter(m)
	)
	assert.NotNil(t, w)

	w.Write(util.Of2(1, "2"))
	assert.Equal(t, map[int]string{1: "2"}, m)
}

func TestOfIOWriterAsBytes(t *testing.T) {
	var (
		str strings.Builder
		w   = OfIOWriterAsBytes(&str)
	)
	assert.NotNil(t, w)

	w.Write(0x41)
	assert.Equal(t, "A", str.String())
}

func TestOfIOWriterAsRunes(t *testing.T) {
	var (
		str strings.Builder
		w   = OfIOWriterAsRunes(&str)
	)
	assert.NotNil(t, w)

	w.Write('A')
	assert.Equal(t, "A", str.String())
}

func TestOfIOWriterAsStrings(t *testing.T) {
	var (
		str strings.Builder
		w   = OfIOWriterAsStrings(&str)
	)
	assert.NotNil(t, w)

	w.Write("A")
	assert.Equal(t, "A", str.String())
}

func TestOfIOWriterAsLines(t *testing.T) {
	var (
		str strings.Builder
		w   = OfIOWriterAsLines(&str)
	)
	assert.NotNil(t, w)

	w.Write("A")
	assert.Equal(t, "A"+osEOLSequence, str.String())
}

func TestWrite(t *testing.T) {
	var (
		err = fmt.Errorf("died")
		w   = NewWriter(func(any) error { return err })
	)
	assert.NotNil(t, w)
	assert.Equal(t, err, w.Write(0))
}
