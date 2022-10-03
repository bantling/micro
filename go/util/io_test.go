package util

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorReader(t *testing.T) {
	// Read first char, then remaining, then error
	var (
		err = fmt.Errorf("an error")
		rdr = NewErrorReader([]byte("hello"), err)
		p   = make([]byte, 5)
	)

	n, e := rdr.Read(p[0:1])
	assert.Equal(t, 1, n)
	assert.Equal(t, byte('h'), p[0])
	assert.Nil(t, e)

	n, e = rdr.Read(p[0:4])
	assert.Equal(t, 4, n)
	assert.Equal(t, []byte("ello"), p[0:4])
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

	// Read all 4 chars at once with a 5 byte buffer, then error

	rdr = NewErrorReader([]byte("hell"), err)
	n, e = rdr.Read(p)
	assert.Equal(t, 4, n)
	assert.Equal(t, []byte("hell"), p[0:4])

	n, e = rdr.Read(p)
	assert.Equal(t, 0, n)
	assert.Equal(t, err, e)
}
