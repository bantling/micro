package tuple

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

func TestOf2_(t *testing.T) {
	assert.Equal(t, Two[string, int]{"a", 1}, Of2("a", 1))
}

func TestOf2Same_(t *testing.T) {
	assert.Equal(t, Two[string, string]{"a", "b"}, Of2Same("a", "b"))
}

func TestOf3_(t *testing.T) {
	assert.Equal(t, Three[string, int, uint]{"a", 1, 2}, Of3("a", 1, uint(2)))
}

func TestOf3Same_(t *testing.T) {
	assert.Equal(t, Three[string, string, string]{"a", "b", "c"}, Of3Same("a", "b", "c"))
}

func TestOf4_(t *testing.T) {
	assert.Equal(t, Four[string, int, uint, string]{"a", 1, 2, "b"}, Of4("a", 1, uint(2), "b"))
}

func TestOf4Same_(t *testing.T) {
	assert.Equal(t, Four[string, string, string, string]{"a", "b", "c", "d"}, Of4Same("a", "b", "c", "d"))
}

// ==== methods

func TestValues_(t *testing.T) {
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
