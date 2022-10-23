package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/bantling/micro/go/funcs"
	"github.com/bantling/micro/go/util"
	"github.com/stretchr/testify/assert"
)

// ==== Foundation funcs

func TestMap(t *testing.T) {
	it := Map(strconv.Itoa)(Of(1, 2))
	assert.Equal(t, "1", it.Must())
	assert.Equal(t, "2", it.Must())
	assert.False(t, it.Next())
}

func TestFilter(t *testing.T) {
	it := Filter(func(val int) bool { return val > 1 })(Of(1, 2))
	assert.Equal(t, 2, it.Must())
	assert.False(t, it.Next())
}

func TestReduce(t *testing.T) {
	// Reeducer func
	fn := func(i, j int) int { return i + j }

	// No identity => sum(1, 2, 3) = 6
	it := Reduce(fn)(Of(1, 2, 3))
	assert.Equal(t, 6, it.Must())
	assert.False(t, it.Next())

	// Identity = 4 => sum(4, 1, 2, 3) = 10
	it = Reduce(fn, 4)(Of(1, 2, 3))
	assert.Equal(t, 10, it.Must())
	assert.False(t, it.Next())

	// Empty set, no identity
	it = Reduce(fn)(OfEmpty[int]())
	assert.False(t, it.Next())

	// Empty set, identity = 4
	it = Reduce(fn, 4)(OfEmpty[int]())
	assert.Equal(t, 4, it.Must())
	assert.False(t, it.Next())
}

func TestReduceTo(t *testing.T) {
	// No identity => concat(1, 2, 3) = "123"
	fn := func(i string, j int) string { return i + strconv.Itoa(j) }
	it := ReduceTo(fn)(Of(1, 2, 3))
	assert.Equal(t, "123", it.Must())
	assert.False(t, it.Next())

	// Identity = "4" => concat(4, 1, 2, 3) = "4123"
	it = ReduceTo(fn, "4")(Of(1, 2, 3))
	assert.Equal(t, "4123", it.Must())
	assert.False(t, it.Next())

	// Empty set, no identity
	it = ReduceTo(fn)(OfEmpty[int]())
	assert.Equal(t, "", it.Must())
	assert.False(t, it.Next())

	// Empty set, identity = "4"
	it = ReduceTo(fn, "4")(OfEmpty[int]())
	assert.Equal(t, "4", it.Must())
	assert.False(t, it.Next())
}

func TestReduceToBool(t *testing.T) {
	// And logic (all match): identity = true, stop on false
	it := ReduceToBool(func(i int) bool { return i < 3 }, true, false)(Of(1, 2, 3))
	assert.False(t, it.Must())
	assert.False(t, it.Next())

	// Or logic (at least one match): identity = falsee, stop on true
	it = ReduceToBool(func(i int) bool { return i < 3 }, false, true)(Of(1, 2, 3))
	assert.True(t, it.Must())
	assert.False(t, it.Next())
}

func TestReduceExpandSlice(t *testing.T) {
	// Reduce into a new generated slice
	it := ReduceToSlice(Of(1, 2))
	assert.Equal(t, []int{1, 2}, it.Must())
	assert.False(t, it.Next())

	// Reduce into an existing slice
	slc := make([]int, 2)
	it = ReduceIntoSlice(slc)(Of(1, 2))
	cmp := it.Must()
	assert.Equal(t, []int{1, 2}, cmp)
	assert.Equal(t, fmt.Sprintf("%p", slc), fmt.Sprintf("%p", cmp))
	assert.False(t, it.Next())

	it = ReduceToSlice(ExpandSlices(Of([]int{1, 2, 3}, nil, []int{}, []int{4, 5})))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, it.Must())
	assert.False(t, it.Next())
}

func TestReduceExpandMap(t *testing.T) {
	it := ReduceToMap(Of(util.KVOf(1, "1"), util.KVOf(2, "2"), util.KVOf(3, "3")))
	assert.Equal(t, map[int]string{1: "1", 2: "2", 3: "3"}, it.Must())
	assert.False(t, it.Next())

	it = ReduceToMap(ExpandMaps(Of(map[int]string{1: "1", 2: "2"}, nil, map[int]string{}, map[int]string{3: "3"})))
	assert.Equal(t, map[int]string{1: "1", 2: "2", 3: "3"}, it.Must())
	assert.False(t, it.Next())
}

func TestTransform(t *testing.T) {
	it := Transform(func(it Iter[int]) (int, bool) {
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

	assert.Equal(t, 3, it.Must())
	assert.Equal(t, 7, it.Must())
	assert.Equal(t, 5, it.Must())
	assert.False(t, it.Next())
}

// ==== Funcs based on foundational funcs

func TestAllMatch(t *testing.T) {
	fn := AllMatch(func(i int) bool { return i < 3 })
	it := fn(OfEmpty[int]())
	assert.True(t, it.Must())
	assert.False(t, it.Next())

	it = fn(Of(1))
	assert.True(t, it.Must())
	assert.False(t, it.Next())

	it = fn(Of(1, 2, 3))
	assert.False(t, it.Must())
	assert.False(t, it.Next())
}

func TestAnyMatch(t *testing.T) {
	fn := AnyMatch(func(i int) bool { return i < 3 })
	it := fn(OfEmpty[int]())
	assert.False(t, it.Must())
	assert.False(t, it.Next())

	it = fn(Of(1))
	assert.True(t, it.Must())
	assert.False(t, it.Next())

	it = fn(Of(4, 5, 6))
	assert.False(t, it.Must())
	assert.False(t, it.Next())
}

func TestNoneMatch(t *testing.T) {
	fn := NoneMatch(func(i int) bool { return i < 3 })
	it := fn(OfEmpty[int]())
	assert.True(t, it.Must())
	assert.False(t, it.Next())

	it = fn(Of(1))
	assert.False(t, it.Must())
	assert.False(t, it.Next())

	it = fn(Of(4, 5, 6))
	assert.True(t, it.Must())
	assert.False(t, it.Next())
}

func TestCount(t *testing.T) {
	fn := Count[int]()
	it := fn(OfEmpty[int]())
	assert.Equal(t, 0, it.Must())
	assert.False(t, it.Next())

	it = fn(Of(1))
	assert.Equal(t, 1, it.Must())
	assert.False(t, it.Next())

	it = fn(Of(1, 2, 3))
	assert.Equal(t, 3, it.Must())
	assert.False(t, it.Next())
}

func TestGeneratorDistinctDuplicate(t *testing.T) {
	// Distinct
	fn := Distinct[int]()
	it := fn(OfEmpty[int]())
	assert.False(t, it.Next())

	it = fn(OfOne(1))
	assert.Equal(t, 1, it.Must())
	assert.False(t, it.Next())

	it = fn(Of(1, 3, 2, 3, 2, 1))
	assert.Equal(t, []int{1, 3, 2}, First(ReduceToSlice(it)))
	assert.False(t, it.Next())

	// Distinct across multiple iters via Concat
	it = fn(Concat(OfEmpty[int](), OfOne(1), Of(1, 2, 3, 3, 2, 1), Of(1, 4)))
	assert.Equal(t, []int{1, 2, 3, 4}, First(ReduceToSlice(it)))

	// Duplicate
	fn = Duplicate[int]()
	it = fn(OfEmpty[int]())
	assert.False(t, it.Next())

	it = fn(OfOne(1))
	assert.False(t, it.Next())

	it = fn(Of(1, 2, 3, 3, 2))
	assert.Equal(t, []int{3, 2}, First(ReduceToSlice(it)))
	assert.False(t, it.Next())

	// Duplicate across multiple iters via Concat
	it = fn(Concat(OfEmpty[int](), OfOne(1), Of(1, 2, 3, 3, 2), Of(1, 4)))
	assert.Equal(t, []int{1, 3, 2}, First(ReduceToSlice(it)))
}

func TestPeek(t *testing.T) {
	slc := []int{}
	fn := Peek(func(val int) { slc = append(slc, val) })
	it := fn(OfEmpty[int]())
	First(ReduceToSlice(it))
	assert.Equal(t, []int{}, slc)

	it = fn(OfOne(1))
	First(ReduceToSlice(it))
	assert.Equal(t, []int{1}, slc)

	it = fn(Of(2, 3, 4))
	First(ReduceToSlice(it))
	assert.Equal(t, []int{1, 2, 3, 4}, slc)
}

func TestSkip(t *testing.T) {
	fn := Skip[int](3)
	it := fn(OfEmpty[int]())
	assert.Equal(t, []int{}, First(ReduceToSlice(it)))

	it = fn(OfOne(1))
	assert.Equal(t, []int{}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2))
	assert.Equal(t, []int{}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3))
	assert.Equal(t, []int{}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3, 4))
	assert.Equal(t, []int{4}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3, 4, 5))
	assert.Equal(t, []int{4, 5}, First(ReduceToSlice(it)))

	funcs.TryTo(
		func() {
			Skip[int](-1)
			assert.Fail(t, "Must die")
		},
		func(err any) {
			assert.Equal(t, errSkipLimitValueCannotBeNegative, err)
		},
	)

	fn = Limit[int](3)
	it = fn(OfEmpty[int]())
	assert.Equal(t, []int{}, First(ReduceToSlice(it)))

	it = fn(OfOne(1))
	assert.Equal(t, []int{1}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2))
	assert.Equal(t, []int{1, 2}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3))
	assert.Equal(t, []int{1, 2, 3}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3, 4))
	assert.Equal(t, []int{1, 2, 3}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3, 4, 5))
	assert.Equal(t, []int{1, 2, 3}, First(ReduceToSlice(it)))

	funcs.TryTo(
		func() {
			Limit[int](-1)
			assert.Fail(t, "Must die")
		},
		func(err any) {
			assert.Equal(t, errSkipLimitValueCannotBeNegative, err)
		},
	)

	fn = funcs.Compose2(Skip[int](1), Limit[int](3))
	it = fn(OfEmpty[int]())
	assert.Equal(t, []int{}, First(ReduceToSlice(it)))

	it = fn(OfOne(1))
	assert.Equal(t, []int{}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2))
	assert.Equal(t, []int{2}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3))
	assert.Equal(t, []int{2, 3}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3, 4))
	assert.Equal(t, []int{2, 3, 4}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3, 4, 5))
	assert.Equal(t, []int{2, 3, 4}, First(ReduceToSlice(it)))

	fn = funcs.Compose2(Limit[int](3), Skip[int](1))
	it = fn(OfEmpty[int]())
	assert.Equal(t, []int{}, First(ReduceToSlice(it)))

	it = fn(OfOne(1))
	assert.Equal(t, []int{}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2))
	assert.Equal(t, []int{2}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3))
	assert.Equal(t, []int{2, 3}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3, 4))
	assert.Equal(t, []int{2, 3}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3, 4, 5))
	assert.Equal(t, []int{2, 3}, First(ReduceToSlice(it)))
}

func TestReverse(t *testing.T) {
	fn := Reverse[int]()
	it := fn(OfEmpty[int]())
	assert.Equal(t, []int{}, First(ReduceToSlice(it)))

	it = fn(OfOne(1))
	assert.Equal(t, []int{1}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2))
	assert.Equal(t, []int{2, 1}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3))
	assert.Equal(t, []int{3, 2, 1}, First(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3, 4))
	assert.Equal(t, []int{4, 3, 2, 1}, First(ReduceToSlice(it)))
}

func TestSort(t *testing.T) {
	{
		fn := SortOrdered[int]()
		it := fn(Of(1, 3, 2))
		assert.Equal(t, []int{1, 2, 3}, First(ReduceToSlice(it)))
	}

	{
		fn := SortComplex[complex128]()
		it := fn(Of(1+0i, 3+1i, 2+0i))
		assert.Equal(t, []complex128{1 + 0i, 2 + 0i, 3 + 1i}, First(ReduceToSlice(it)))
	}

	{
		fn := SortCmp[*big.Int]()
		it := fn(Of(big.NewInt(2), big.NewInt(3), big.NewInt(1)))
		assert.Equal(t, []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}, First(ReduceToSlice(it)))
	}

	{
		fn := SortBy(func(i, j int) bool { return j < i })
		it := fn(Of(1, 3, 2))
		assert.Equal(t, []int{3, 2, 1}, First(ReduceToSlice(it)))
	}
}

func TestGenerateRanges(t *testing.T) {
	// ==== square root method
	assert.Equal(t, [][]uint{{0, 1}, {1, 2}}, generateRanges(2, []PInfo{}))
	assert.Equal(t, [][]uint{{0, 2}, {2, 3}}, generateRanges(3, []PInfo{}))
	assert.Equal(t, [][]uint{{0, 2}, {2, 4}}, generateRanges(4, []PInfo{}))
	assert.Equal(t, [][]uint{{0, 3}, {3, 6}, {6, 9}, {9, 10}}, generateRanges(10, []PInfo{}))
	assert.Equal(t, [][]uint{{0, 4}, {4, 8}, {8, 12}, {12, 15}}, generateRanges(15, []PInfo{}))

	// ==== threads method - first remainder threads get one extra item each
	assert.Equal(t, [][]uint{{0, 1}, {1, 2}}, generateRanges(2, []PInfo{{2, Threads}}))
	assert.Equal(t, [][]uint{{0, 2}, {2, 3}}, generateRanges(3, []PInfo{{2, Threads}}))
	assert.Equal(t, [][]uint{{0, 2}, {2, 4}}, generateRanges(4, []PInfo{{2, Threads}}))
	assert.Equal(t, [][]uint{{0, 5}, {5, 10}}, generateRanges(10, []PInfo{{2, Threads}}))
	assert.Equal(t, [][]uint{{0, 8}, {8, 15}}, generateRanges(15, []PInfo{{2, Threads}}))

	// 2 items with 3 threads = 2 threads, can't have more threads than items
	assert.Equal(t, [][]uint{{0, 1}, {1, 2}}, generateRanges(2, []PInfo{{3, Threads}}))

	assert.Equal(t, [][]uint{{0, 1}, {1, 2}, {2, 3}}, generateRanges(3, []PInfo{{3, Threads}}))
	assert.Equal(t, [][]uint{{0, 2}, {2, 3}, {3, 4}}, generateRanges(4, []PInfo{{3, Threads}}))
	assert.Equal(t, [][]uint{{0, 4}, {4, 7}, {7, 10}}, generateRanges(10, []PInfo{{3, Threads}}))
	assert.Equal(t, [][]uint{{0, 5}, {5, 10}, {10, 15}}, generateRanges(15, []PInfo{{3, Threads}}))

	// ==== items per thread method - any remainder is an additional thread
	assert.Equal(t, [][]uint{{0, 2}}, generateRanges(2, []PInfo{{2, Items}}))
	assert.Equal(t, [][]uint{{0, 2}, {2, 3}}, generateRanges(3, []PInfo{{2, Items}}))
	assert.Equal(t, [][]uint{{0, 2}, {2, 4}}, generateRanges(4, []PInfo{{2, Items}}))
	assert.Equal(t, [][]uint{{0, 2}, {2, 4}, {4, 6}, {6, 8}, {8, 10}}, generateRanges(10, []PInfo{{2, Items}}))
	assert.Equal(t, [][]uint{{0, 2}, {2, 4}, {4, 6}, {6, 8}, {8, 10}, {10, 12}, {12, 14}, {14, 15}}, generateRanges(15, []PInfo{{2, Items}}))

	// 2 items with 3 items per thread = 2 items in 1 thread, bucket size can't exceed number of items
	assert.Equal(t, [][]uint{{0, 2}}, generateRanges(2, []PInfo{{3, Items}}))

	assert.Equal(t, [][]uint{{0, 3}}, generateRanges(3, []PInfo{{3, Items}}))
	assert.Equal(t, [][]uint{{0, 3}, {3, 4}}, generateRanges(4, []PInfo{{3, Items}}))
	assert.Equal(t, [][]uint{{0, 3}, {3, 6}, {6, 9}, {9, 10}}, generateRanges(10, []PInfo{{3, Items}}))
	assert.Equal(t, [][]uint{{0, 3}, {3, 6}, {6, 9}, {9, 12}, {12, 15}}, generateRanges(15, []PInfo{{3, Items}}))
}

func TestParallel(t *testing.T) {
	var (
		infoThreads = PInfo{5, Threads}
		infoItems   = PInfo{5, Items}

		intFn       = Map(func(i int) int { return i * 2 })
		pIntSqrt    = Parallel(intFn)
		pIntThreads = Parallel(intFn, infoThreads)
		pIntItems   = Parallel(intFn, infoItems)

		uintFn       = Map(func(i int) uint { return uint(i * 2) })
		pUintSqrt    = Parallel(uintFn)
		pUintThreads = Parallel(uintFn, infoThreads)
		pUintItems   = Parallel(uintFn, infoItems)
	)
	for _, i := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100, 103} {
		inInt, outInt, outUint := make([]int, i), make([]int, i), make([]uint, i)
		for j := 1; j <= i; j++ {
			inInt[j-1] = j
			outInt[j-1] = j * 2
			outUint[j-1] = uint(j * 2)
		}

		// Same type, modify slice in place
		assert.Equal(t, outInt, First(ReduceToSlice(pIntSqrt(Of(inInt...)))))
		assert.Equal(t, outInt, First(ReduceToSlice(pIntThreads(Of(inInt...)))))
		assert.Equal(t, outInt, First(ReduceToSlice(pIntItems(Of(inInt...)))))

		// Different type, generate a new slice
		assert.Equal(t, outUint, First(ReduceToSlice(pUintSqrt(Of(inInt...)))))
		assert.Equal(t, outUint, First(ReduceToSlice(pUintThreads(Of(inInt...)))))
		assert.Equal(t, outUint, First(ReduceToSlice(pUintItems(Of(inInt...)))))
	}
}

// ==== Composition

func TestCompose(t *testing.T) {
	fn := funcs.Compose5(
		Map(strconv.Itoa),
		Map(func(s string) int { i, _ := strconv.Atoi(s); return i }),
		Filter(func(val int) bool { return val&1 == 1 }),
		ReduceToSlice[int],
		First[[]int],
	)

	assert.Equal(t, []int{1, 3}, fn(Of(1, 2, 3)))
	assert.Equal(t, []int{1, 3}, fn(Of(1, 2, 3)))
}
