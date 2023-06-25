package io

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorReader_(t *testing.T) {
	// Read first char, then remaining, then error
	var (
		err = fmt.Errorf("an error")
		rdr = NewErrorReader([]byte("hello"), err)
		p   = make([]byte, 5)
	)

	n, e := rdr.Read(p[:1])
	assert.Equal(t, 1, n)
	assert.Equal(t, byte('h'), p[0])
	assert.Nil(t, e)

	n, e = rdr.Read(p[:4])
	assert.Equal(t, 4, n)
	assert.Equal(t, []byte("ello"), p[:4])
	assert.Nil(t, e)

	n, e = rdr.Read(p)
	assert.Equal(t, 0, n)
	assert.Equal(t, err, e)

	// Read all 5 chars at once, then error

	rdr = NewErrorReader([]byte("hello"), err)
	n, e = rdr.Read(p)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("hello"), p)

	n, e = rdr.Read(p)
	assert.Equal(t, 0, n)
	assert.Equal(t, err, e)

	// Read all 4 chars at once into a 5 byte buffer, then error

	rdr = NewErrorReader([]byte("hell"), err)
	n, e = rdr.Read(p)
	assert.Equal(t, 4, n)
	assert.Equal(t, []byte("hell"), p[:4])

	n, e = rdr.Read(p)
	assert.Equal(t, 0, n)
	assert.Equal(t, err, e)

	// We can return eof if we want to
	rdr = NewErrorReader([]byte("hello"), io.EOF)
	n, e = rdr.Read(p)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("hello"), p)

	n, e = rdr.Read(p)
	assert.Equal(t, 0, n)
	assert.Equal(t, io.EOF, e)
}

func TestErrorWriter_(t *testing.T) {
	// Write first char, then remaining, then error
	var (
		err = fmt.Errorf("an error")
		p   = []byte("hello")
		w   = NewErrorWriter(len(p), err)
	)

	n, e := w.Write(p[:1])
	assert.Equal(t, 1, n)
	assert.Equal(t, []byte("h"), w.Output())
	assert.Nil(t, e)

	n, e = w.Write(p[1:])
	assert.Equal(t, 4, n)
	assert.Equal(t, []byte("hello"), w.Output())
	assert.Nil(t, e)

	n, e = w.Write(p)
	assert.Equal(t, 0, n)
	assert.Equal(t, err, e)

	// Write all 5 chars at once, then error

	w = NewErrorWriter(len(p), err)
	n, e = w.Write(p)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("hello"), w.Output())
	assert.Nil(t, e)

	n, e = w.Write(p)
	assert.Equal(t, 0, n)
	assert.Equal(t, err, e)

	// Write all 5 chars at once into a 4 byte ouput, then error

	w = NewErrorWriter(len(p)-1, err)
	n, e = w.Write(p)
	assert.Equal(t, 4, n)
	assert.Equal(t, []byte("hell"), w.Output())

	n, e = w.Write(p)
	assert.Equal(t, 0, n)
	assert.Equal(t, err, e)

	// We can return eof if we want to

	w = NewErrorWriter(len(p), io.EOF)
	n, e = w.Write(p)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("hello"), w.Output())
	assert.Nil(t, e)

	n, e = w.Write(p)
	assert.Equal(t, 0, n)
	assert.Equal(t, io.EOF, e)

	// We can return 0, nil which shd never happen in real life

	w = NewErrorWriter(-1, io.EOF)
	n, e = w.Write(p)
	assert.Equal(t, 0, n)
	assert.Equal(t, []byte{}, w.Output())
	assert.Nil(t, e)

	n, e = w.Write(p)
	assert.Equal(t, []byte{}, w.Output())
	assert.Equal(t, 0, n)
	assert.Nil(t, e)
}
