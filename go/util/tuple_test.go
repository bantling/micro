package util

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	anErr = fmt.Errorf("An error")
)

// ==== Constructors

func TestOf2(t *testing.T) {
	assert.Equal(t, Tuple2[string, int]{"a", 1}, Of2("a", 1))
}

func TestOf2Same(t *testing.T) {
	assert.Equal(t, Tuple2[string, string]{"a", "b"}, Of2Same("a", "b"))
}

func TestOf2Error(t *testing.T) {
	assert.Equal(t, Tuple2[string, error]{"a", anErr}, Of2Error("a", anErr))
}

func TestOf3(t *testing.T) {
	assert.Equal(t, Tuple3[string, int, uint]{"a", 1, 2}, Of3("a", 1, uint(2)))
}

func TestOf3Same(t *testing.T) {
	assert.Equal(t, Tuple3[string, string, string]{"a", "b", "c"}, Of3Same("a", "b", "c"))
}

func TestOf3Error(t *testing.T) {
	assert.Equal(t, Tuple3[string, int, error]{"a", 1, anErr}, Of3Error("a", 1, anErr))
}

func TestOf4(t *testing.T) {
	assert.Equal(t, Tuple4[string, int, uint, string]{"a", 1, 2, "b"}, Of4("a", 1, uint(2), "b"))
}

func TestOf4Same(t *testing.T) {
	assert.Equal(t, Tuple4[string, string, string, string]{"a", "b", "c", "d"}, Of4Same("a", "b", "c", "d"))
}

func TestOf4Error(t *testing.T) {
	assert.Equal(t, Tuple4[string, int, uint, error]{"a", 1, 2, anErr}, Of4Error("a", 1, uint(2), anErr))
}

// ==== methods

func TestValues(t *testing.T) {
	at, au := Of2("a", 1).Values()
	assert.Equal(t, "a", at)
	assert.Equal(t, 1, au)

	at, au, av := Of3("a", 1, 2).Values()
	assert.Equal(t, "a", at)
	assert.Equal(t, 1, au)
	assert.Equal(t, 2, av)

	at, au, av, aw := Of4("a", 1, 2, "b").Values()
	assert.Equal(t, "a", at)
	assert.Equal(t, 1, au)
	assert.Equal(t, 2, av)
	assert.Equal(t, "b", aw)
}
