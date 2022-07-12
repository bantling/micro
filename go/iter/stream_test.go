package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"strconv"
	"testing"

	"github.com/bantling/micro/go/funcs"
	"github.com/stretchr/testify/assert"
)

// ==== Foundation funcs

func TestMap(t *testing.T) {
	it := Map(strconv.Itoa)(Of(1, 2))
	assert.True(t, it.Next())
	assert.Equal(t, "1", it.Value())
	assert.True(t, it.Next())
	assert.Equal(t, "2", it.Value())
	assert.False(t, it.Next())
}

func TestFilter(t *testing.T) {
	it := Filter(func(val int) bool { return val > 1 })(Of(1, 2))
	assert.True(t, it.Next())
	assert.Equal(t, 2, it.Value())
	assert.False(t, it.Next())
}

func TestReduce(t *testing.T) {
	// No identity => sum(1, 2, 3) = 6
	it := Reduce(func(i, j int) int { return i + j })(Of(1, 2, 3))
	assert.True(t, it.Next())
	assert.Equal(t, 6, it.Value())
	assert.False(t, it.Next())

	// Identity = 4 => sum(4, 1, 2, 3) = 10
	it = Reduce(func(i, j int) int { return i + j }, 4)(Of(1, 2, 3))
	assert.True(t, it.Next())
	assert.Equal(t, 10, it.Value())
	assert.False(t, it.Next())
}

func TestReduceTo(t *testing.T) {
	// No identity => concat(1, 2, 3) = "123"
	it := ReduceTo(func(i string, j int) string { return i + strconv.Itoa(j) })(Of(1, 2, 3))
	assert.True(t, it.Next())
	assert.Equal(t, "123", it.Value())
	assert.False(t, it.Next())

	// Identity = "4" => concat(4, 1, 2, 3) = "4123"
	it = ReduceTo(func(i string, j int) string { return i + strconv.Itoa(j) }, "4")(Of(1, 2, 3))
	assert.True(t, it.Next())
	assert.Equal(t, "4123", it.Value())
	assert.False(t, it.Next())
}

func TestReduceExpandSlice(t *testing.T) {
	it := ReduceToSlice(Of(1, 2))
	assert.True(t, it.Next())
	assert.Equal(t, []int{1, 2}, it.Value())
	assert.False(t, it.Next())

	it = ReduceToSlice(ExpandSlices(Of([]int{1, 2, 3}, nil, []int{}, []int{4, 5})))
	assert.True(t, it.Next())
	assert.Equal(t, []int{1, 2, 3, 4, 5}, it.Value())
	assert.False(t, it.Next())
}

func TestReduceExpandMap(t *testing.T) {
	it := ReduceToMap(Of(KVOf(1, "1"), KVOf(2, "2"), KVOf(3, "3")))
	assert.True(t, it.Next())
	assert.Equal(t, map[int]string{1: "1", 2: "2", 3: "3"}, it.Value())
	assert.False(t, it.Next())

	it = ReduceToMap(ExpandMaps(Of(map[int]string{1: "1", 2: "2"}, nil, map[int]string{}, map[int]string{3: "3"})))
	assert.True(t, it.Next())
	assert.Equal(t, map[int]string{1: "1", 2: "2", 3: "3"}, it.Value())
	assert.False(t, it.Next())
}

func TestTransform(t *testing.T) {
	it := Transform(func(it *Iter[int]) (int, bool) {
		// Sum pairs of ints
		if !it.Next() {
			return 0, false
		}

		res := it.Value()
		if it.Next() {
			res += it.Value()
		}

		return res, true
	})(Of(1, 2, 3, 4, 5))

	assert.True(t, it.Next())
	assert.Equal(t, 3, it.Value())
	assert.True(t, it.Next())
	assert.Equal(t, 7, it.Value())
	assert.True(t, it.Next())
	assert.Equal(t, 5, it.Value())
	assert.False(t, it.Next())
}

// ==== Composition

func TestCompose(t *testing.T) {
	slc := funcs.Compose5(
		Map(strconv.Itoa),
		Map(func(s string) int { i, _ := strconv.Atoi(s); return i }),
		Filter(func(val int) bool { return val&1 == 1 }),
		ReduceToSlice[int],
		First[[]int],
	)(Of(1, 2, 3))

	assert.Equal(t, []int{1, 3}, slc)
}
