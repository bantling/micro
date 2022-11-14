package iter

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"github.com/bantling/micro/go/funcs"
	"github.com/bantling/micro/go/util"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strconv"
	"testing"
)

// ==== Foundation funcs

func TestMap(t *testing.T) {
	it := Map(strconv.Itoa)(Of(1, 2))
	assert.Equal(t, util.Of2Error("1", nil), Maybe(it))
	assert.Equal(t, util.Of2Error("2", nil), Maybe(it))
	assert.Equal(t, util.Of2Error("", EOI), Maybe(it))
}

func TestFilter(t *testing.T) {
	it := Filter(func(val int) bool { return val > 1 })(Of(1, 2))
	assert.Equal(t, util.Of2Error(2, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
}

func TestReduce(t *testing.T) {
	// Reeducer func
	fn := func(i, j int) int { return i + j }

	// No identity, () => sum() = empty
	it := Reduce(fn)(OfEmpty[int]())
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	// No identity, (1) => sum(1) = 1
	it = Reduce(fn)(OfOne(1))
	assert.Equal(t, util.Of2Error(1, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	// No identity, (1,  2, 3) => sum(1, 2, 3) = 6
	it = Reduce(fn)(Of(1, 2, 3))
	assert.Equal(t, util.Of2Error(6, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	// Identity 4, () => sum(4) = 4
	it = Reduce(fn, 4)(OfEmpty[int]())
	assert.Equal(t, util.Of2Error(4, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	// Identity 4, (1) => sum(4, 1) = 5
	it = Reduce(fn, 4)(OfOne(1))
	assert.Equal(t, util.Of2Error(5, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	// Identity 4, (1, 2, 3) => sum(4, 1, 2, 3) = 10
	it = Reduce(fn, 4)(Of(1, 2, 3))
	assert.Equal(t, util.Of2Error(10, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	// test non-EOI errors
	{
		var (
			anErr = fmt.Errorf("An err")
			fn    = func(i, j int) int { return i + j }
		)
		// no identity, error after first element
		it := Reduce(fn)(SetError(Of(1), anErr))
		assert.Equal(t, util.Of2Error(0, anErr), Maybe(it))

		// identity, error before first element
		it = Reduce(fn, 0)(SetError(OfEmpty[int](), anErr))
		assert.Equal(t, util.Of2Error(0, anErr), Maybe(it))

		// identity, error after first element
		it = Reduce(fn, 0)(SetError(Of(1), anErr))
		assert.Equal(t, util.Of2Error(0, anErr), Maybe(it))
	}
}

func TestReduceTo(t *testing.T) {
	// No identity, () => concat() = ()
	fn := func(i string, j int) string { return i + strconv.Itoa(j) }
	it := ReduceTo(fn)(OfEmpty[int]())
	assert.Equal(t, util.Of2Error("", EOI), Maybe(it))

	// No identity, (1) => concat(1) = ("1")
	it = ReduceTo(fn)(OfOne(1))
	assert.Equal(t, util.Of2Error("1", nil), Maybe(it))
	assert.Equal(t, util.Of2Error("", EOI), Maybe(it))

	// No identity, (1, 2, 3) => concat(1, 2, 3) = "123"
	it = ReduceTo(fn)(Of(1, 2, 3))
	assert.Equal(t, util.Of2Error("123", nil), Maybe(it))
	assert.Equal(t, util.Of2Error("", EOI), Maybe(it))

	// Identity "4", () => concat(4) = "4"
	it = ReduceTo(fn, "4")(OfEmpty[int]())
	assert.Equal(t, util.Of2Error("4", nil), Maybe(it))
	assert.Equal(t, util.Of2Error("", EOI), Maybe(it))

	// Identity "4", (1) => concat(4, 1) = "41"
	it = ReduceTo(fn, "4")(OfOne(1))
	assert.Equal(t, util.Of2Error("41", nil), Maybe(it))
	assert.Equal(t, util.Of2Error("", EOI), Maybe(it))

	// Identity "4", (1, 2, 3) => concat(4, 1, 2, 3) = "4123"
	it = ReduceTo(fn, "4")(Of(1, 2, 3))
	assert.Equal(t, util.Of2Error("4123", nil), Maybe(it))
	assert.Equal(t, util.Of2Error("", EOI), Maybe(it))

	// test non-EOI errors
	{
		var (
			anErr = fmt.Errorf("An err")
			fn    = func(i, j int) int { return i + j }
		)
		// no identity, error after first element
		it := ReduceTo(fn)(SetError(Of(1), anErr))
		assert.Equal(t, util.Of2Error(0, anErr), Maybe(it))

		// identity, error before first element
		it = ReduceTo(fn, 0)(SetError(OfEmpty[int](), anErr))
		assert.Equal(t, util.Of2Error(0, anErr), Maybe(it))

		// identity, error after first element
		it = ReduceTo(fn, 0)(SetError(Of(1), anErr))
		assert.Equal(t, util.Of2Error(0, anErr), Maybe(it))
	}
}

func TestReduceToBool(t *testing.T) {
	// And logic (all match): identity = true, stop on false
	it := ReduceToBool(func(i int) bool { return i < 3 }, true, false)(OfEmpty[int]())
	assert.Equal(t, util.Of2Error(true, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(false, EOI), Maybe(it))

	it = ReduceToBool(func(i int) bool { return i < 3 }, true, false)(Of(1, 2, 3))
	assert.Equal(t, util.Of2Error(false, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(false, EOI), Maybe(it))

	// Or logic (at least one match): identity = false, stop on true
	it = ReduceToBool(func(i int) bool { return i < 3 }, false, true)(Of(4))
	assert.Equal(t, util.Of2Error(false, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(false, EOI), Maybe(it))

	it = ReduceToBool(func(i int) bool { return i < 3 }, false, true)(Of(1))
	assert.Equal(t, util.Of2Error(true, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(false, EOI), Maybe(it))

	// test non-EOI errors
	{
		var (
			anErr = fmt.Errorf("An err")
			fn    = func(i int) bool { return i > 0 }
		)
		// error before first element
		it := ReduceToBool(fn, true, false)(SetError(OfEmpty[int](), anErr))
		assert.Equal(t, util.Of2Error(false, anErr), Maybe(it))

		// error after first element
		it = ReduceToBool(fn, true, false)(SetError(Of(1), anErr))
		assert.Equal(t, util.Of2Error(false, anErr), Maybe(it))
	}
}

func TestReduceToSlice(t *testing.T) {
	// Reduce into a new generated slice, no error
	it := ReduceToSlice(Of(1, 2))
	assert.Equal(t, util.Of2Error([]int{1, 2}, nil), Maybe(it))
	assert.Equal(t, util.Of2Error([]int(nil), EOI), Maybe(it))

	// Error before first element
	{
		anErr := fmt.Errorf("An err")
		it := ReduceToSlice(SetError(OfEmpty[int](), anErr))
		assert.Equal(t, util.Of2Error([]int(nil), anErr), Maybe(it))
	}
}

func TestReduceIntoSlice(t *testing.T) {
	// Reduce into a new generated slice, no error
	{
		slc := make([]int, 2)
		it := ReduceIntoSlice(slc)(Of(1, 2))
		assert.Equal(t, util.Of2Error([]int{1, 2}, nil), Maybe(it))
		assert.Equal(t, util.Of2Error([]int(nil), EOI), Maybe(it))
	}

	// Error before first element
	{
		slc := []int{2}
		anErr := fmt.Errorf("An err")
		it := ReduceIntoSlice(slc)(SetError(OfEmpty[int](), anErr))
		assert.Equal(t, util.Of2Error([]int(nil), anErr), Maybe(it))
		// ReduceIntoSlice sets every element of target slice to zero val if an error occurs
		assert.Equal(t, []int{0}, slc)
	}

	// Panics if target slice is not large enough
	{
		slc := []int{1}
		it := ReduceIntoSlice(slc)(Of(1, 2))

		funcs.TryTo(
			func() {
				it.Next()
				assert.Fail(t, "Must die")
			},
			func(e any) {
				assert.Equal(t, "runtime.boundsError{x:1, y:1, signed:true, code:0x0}", fmt.Sprintf("%#v", e))
			},
		)
	}
}

func TestExpandSlices(t *testing.T) {
	it := ReduceToSlice(ExpandSlices(Of([]int{1, 2, 3}, nil, []int{}, []int{4, 5})))
	assert.Equal(t, util.Of2Error([]int{1, 2, 3, 4, 5}, nil), Maybe(it))
	assert.Equal(t, util.Of2Error([]int(nil), EOI), Maybe(it))

	// Error before first element
	{
		anErr := fmt.Errorf("An err")
		it := ExpandSlices(SetError(OfEmpty[[]int](), anErr))
		assert.Equal(t, util.Of2Error(0, anErr), Maybe(it))
	}
}

func TestReduceToMap(t *testing.T) {
	it := ReduceToMap(Of(util.Of2(1, "1"), util.Of2(2, "2"), util.Of2(3, "3")))
	assert.Equal(t, util.Of2Error(map[int]string{1: "1", 2: "2", 3: "3"}, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(map[int]string(nil), EOI), Maybe(it))

	{
		anErr := fmt.Errorf("An err")
		it := ReduceToMap(SetError(OfEmpty[util.Tuple2[int, string]](), anErr))
		assert.Equal(t, util.Of2Error(map[int]string(nil), anErr), Maybe(it))
	}
}

func TestExpandMaps(t *testing.T) {
	it := ReduceToMap(ExpandMaps(Of(map[int]string{1: "1", 2: "2"}, nil, map[int]string{}, map[int]string{3: "3"})))
	assert.Equal(t, util.Of2Error(map[int]string{1: "1", 2: "2", 3: "3"}, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(map[int]string(nil), EOI), Maybe(it))

	{
		anErr := fmt.Errorf("An err")
		it := ExpandMaps(SetError(OfEmpty[map[int]int](), anErr))
		assert.Equal(t, util.Of2Error(util.Of2(0, 0), anErr), Maybe(it))
	}
}

func TestSkip(t *testing.T) {
	fn := Skip[int](3)
	it := fn(OfEmpty[int]())
	assert.Equal(t, util.Of2Error([]int{}, nil), Maybe(ReduceToSlice(it)))

	it = fn(OfOne(1))
	assert.Equal(t, util.Of2Error([]int{}, nil), Maybe(ReduceToSlice(it)))

	it = fn(Of(1, 2))
	assert.Equal(t, util.Of2Error([]int{}, nil), Maybe(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3))
	assert.Equal(t, util.Of2Error([]int{}, nil), Maybe(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3, 4))
	assert.Equal(t, util.Of2Error([]int{4}, nil), Maybe(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3, 4, 5))
	assert.Equal(t, util.Of2Error([]int{4, 5}, nil), Maybe(ReduceToSlice(it)))
}

func TestLimit(t *testing.T) {
	fn := Limit[int](3)
	it := fn(OfEmpty[int]())
	assert.Equal(t, util.Of2Error([]int{}, nil), Maybe(ReduceToSlice(it)))

	it = fn(OfOne(1))
	assert.Equal(t, util.Of2Error([]int{1}, nil), Maybe(ReduceToSlice(it)))

	it = fn(Of(1, 2))
	assert.Equal(t, util.Of2Error([]int{1, 2}, nil), Maybe(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3))
	assert.Equal(t, util.Of2Error([]int{1, 2, 3}, nil), Maybe(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3, 4))
	assert.Equal(t, util.Of2Error([]int{1, 2, 3}, nil), Maybe(ReduceToSlice(it)))

	it = fn(Of(1, 2, 3, 4, 5))
	assert.Equal(t, util.Of2Error([]int{1, 2, 3}, nil), Maybe(ReduceToSlice(it)))
}

func TestPeek(t *testing.T) {
	slc := []int{}
	fn := Peek(func(val int) { slc = append(slc, val) })
	it := fn(OfEmpty[int]())
	// assertNext(t, 0, EOI)(it.Next())
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
	assert.Equal(t, []int{}, slc)

	it = fn(OfOne(1))
	assert.Equal(t, util.Of2Error(1, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
	assert.Equal(t, []int{1}, slc)

	it = fn(Of(2, 3, 4))
	assert.Equal(t, util.Of2Error(2, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(3, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(4, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
	assert.Equal(t, []int{1, 2, 3, 4}, slc)
}

func TestGenerator(t *testing.T) {
	called := 0
	fn := Generator(func() func(Iter[int]) Iter[int] {
		return func(it Iter[int]) Iter[int] {
			called++
			return it
		}
	})

	it := fn(Of(1))
	assert.Equal(t, util.Of2Error(1, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
	assert.Equal(t, 1, called)

	it = fn(Of(3))
	assert.Equal(t, util.Of2Error(3, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
	assert.Equal(t, 2, called)
}

func TestTransform(t *testing.T) {
	it := Transform(func(it Iter[int]) (int, error) {
		// Sum pairs of ints
		val, err := it.Next()
		if err != nil {
			return 0, err
		}

		res := val
		val, err = it.Next()
		if err != nil {
			if err == EOI {
				return res, nil
			}
			return 0, err
		}
		res += val

		return res, nil
	})(Of(1, 2, 3, 4, 5))

	assert.Equal(t, util.Of2Error(3, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(7, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(5, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
}

// // ==== Funcs based on foundational funcs

func TestAllMatch(t *testing.T) {
	fn := AllMatch(func(i int) bool { return i < 3 })
	it := fn(OfEmpty[int]())
	assert.Equal(t, util.Of2Error(true, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(false, EOI), Maybe(it))

	it = fn(Of(1))
	assert.Equal(t, util.Of2Error(true, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(false, EOI), Maybe(it))

	it = fn(Of(1, 2, 3))
	assert.Equal(t, util.Of2Error(false, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(false, EOI), Maybe(it))
}

func TestAnyMatch(t *testing.T) {
	fn := AnyMatch(func(i int) bool { return i < 3 })
	it := fn(OfEmpty[int]())
	assert.Equal(t, util.Of2Error(false, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(false, EOI), Maybe(it))

	it = fn(Of(1))
	assert.Equal(t, util.Of2Error(true, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(false, EOI), Maybe(it))

	it = fn(Of(4, 5, 6))
	assert.Equal(t, util.Of2Error(false, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(false, EOI), Maybe(it))
}

func TestNoneMatch(t *testing.T) {
	fn := NoneMatch(func(i int) bool { return i < 3 })
	it := fn(OfEmpty[int]())
	assert.Equal(t, util.Of2Error(true, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(false, EOI), Maybe(it))

	it = fn(Of(1))
	assert.Equal(t, util.Of2Error(false, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(false, EOI), Maybe(it))

	it = fn(Of(4, 5, 6))
	assert.Equal(t, util.Of2Error(true, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(false, EOI), Maybe(it))
}

func TestCount(t *testing.T) {
	fn := Count[int]()
	it := fn(OfEmpty[int]())
	assert.Equal(t, util.Of2Error(0, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	it = fn(Of(1))
	assert.Equal(t, util.Of2Error(1, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	it = fn(Of(1, 2, 3))
	assert.Equal(t, util.Of2Error(3, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
}

func TestDistinct(t *testing.T) {
	// Distinct
	fn := Distinct[int]()
	it := fn(OfEmpty[int]())
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	it = fn(OfOne(1))
	assert.Equal(t, util.Of2Error(1, nil), Maybe(it))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	it = fn(Of(1, 3, 2, 3, 2, 1))
	assert.Equal(t, util.Of2Error([]int{1, 3, 2}, nil), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	// Distinct across multiple iters via Concat
	it = fn(Concat(OfEmpty[int](), OfOne(1), Of(1, 2, 3, 3, 2, 1), Of(1, 4)))
	assert.Equal(t, util.Of2Error([]int{1, 2, 3, 4}, nil), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
}

func TestDuplicate(t *testing.T) {
	// Duplicate
	fn := Duplicate[int]()
	it := fn(OfEmpty[int]())
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	it = fn(OfOne(1))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	it = fn(Of(1, 2, 3, 3, 2))
	assert.Equal(t, util.Of2Error([]int{3, 2}, nil), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	// Duplicate across multiple iters via Concat
	it = fn(Concat(OfEmpty[int](), OfOne(1), Of(1, 2, 3, 3, 2), Of(1, 4)))
	assert.Equal(t, util.Of2Error([]int{1, 3, 2}, nil), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
}

func TestReverse(t *testing.T) {
	fn := Reverse[int]()
	it := fn(OfEmpty[int]())
	assert.Equal(t, util.Of2Error([]int{}, nil), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	it = fn(OfOne(1))
	assert.Equal(t, util.Of2Error([]int{1}, nil), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	it = fn(Of(1, 2))
	assert.Equal(t, util.Of2Error([]int{2, 1}, nil), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	it = fn(Of(1, 2, 3))
	assert.Equal(t, util.Of2Error([]int{3, 2, 1}, nil), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	it = fn(Of(1, 2, 3, 4))
	assert.Equal(t, util.Of2Error([]int{4, 3, 2, 1}, nil), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	anErr := fmt.Errorf("An err")
	it = fn(SetError(OfEmpty[int](), anErr))
	assert.Equal(t, util.Of2Error([]int(nil), anErr), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error(0, anErr), Maybe(it))
}

func TestSortOrdered(t *testing.T) {
	fn := SortOrdered[int]()
	it := fn(Of(1, 3, 2))
	assert.Equal(t, util.Of2Error([]int{1, 2, 3}, nil), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	anErr := fmt.Errorf("An err")
	it = fn(SetError(OfEmpty[int](), anErr))
	assert.Equal(t, util.Of2Error([]int(nil), anErr), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error(0, anErr), Maybe(it))
}

func TestSortComplex(t *testing.T) {
	fn := SortComplex[complex128]()
	it := fn(Of(1+0i, 3+1i, 2+0i))
	assert.Equal(t, util.Of2Error([]complex128{1 + 0i, 2 + 0i, 3 + 1i}, nil), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error(0+0i, EOI), Maybe(it))

	anErr := fmt.Errorf("An err")
	it = fn(SetError(OfEmpty[complex128](), anErr))
	assert.Equal(t, util.Of2Error([]complex128(nil), anErr), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error(0+0i, anErr), Maybe(it))
}

func TestSortCmp(t *testing.T) {
	fn := SortCmp[*big.Int]()
	it := fn(Of(big.NewInt(2), big.NewInt(3), big.NewInt(1)))
	assert.Equal(t, util.Of2Error([]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}, nil), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error((*big.Int)(nil), EOI), Maybe(it))

	anErr := fmt.Errorf("An err")
	it = fn(SetError(OfEmpty[*big.Int](), anErr))
	assert.Equal(t, util.Of2Error([]*big.Int(nil), anErr), Maybe(ReduceToSlice(it)))
}

func TestSortBy(t *testing.T) {
	fn := SortBy(func(i, j int) bool { return j < i })
	it := fn(Of(1, 3, 2))
	assert.Equal(t, util.Of2Error([]int{3, 2, 1}, nil), Maybe(ReduceToSlice(it)))
	assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

	anErr := fmt.Errorf("An err")
	it = fn(SetError(OfEmpty[int](), anErr))
	assert.Equal(t, util.Of2Error([]int(nil), anErr), Maybe(ReduceToSlice(it)))
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

		uintFn       = Map(func(i int) uint { return uint(i * 3) })
		pUintSqrt    = Parallel(uintFn)
		pUintThreads = Parallel(uintFn, infoThreads)
		pUintItems   = Parallel(uintFn, infoItems)
	)
	for _, i := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100, 103} {
		inInt, outInt, outUint := make([]int, i), make([]int, i), make([]uint, i)
		for j := 1; j <= i; j++ {
			inInt[j-1] = j
			outInt[j-1] = j * 2
			outUint[j-1] = uint(j * 3)
		}

		// Same type, modify slice in place
		assert.Equal(t, util.Of2Error(outInt, nil), Maybe(ReduceToSlice(pIntSqrt(Of(inInt...)))))
		assert.Equal(t, util.Of2Error(outInt, nil), Maybe(ReduceToSlice(pIntThreads(Of(inInt...)))))
		assert.Equal(t, util.Of2Error(outInt, nil), Maybe(ReduceToSlice(pIntItems(Of(inInt...)))))

		// Different type, generate a new slice
		assert.Equal(t, util.Of2Error(outUint, nil), Maybe(ReduceToSlice(pUintSqrt(Of(inInt...)))))
		assert.Equal(t, util.Of2Error(outUint, nil), Maybe(ReduceToSlice(pUintThreads(Of(inInt...)))))
		assert.Equal(t, util.Of2Error(outUint, nil), Maybe(ReduceToSlice(pUintItems(Of(inInt...)))))
	}

	// Error on source iter
	anErr := fmt.Errorf("An err")
	it := ReduceToSlice(pIntSqrt(SetError(OfEmpty[int](), anErr)))
	assert.Equal(t, util.Of2Error([]int(nil), anErr), Maybe(it))

	// Error during transform that executes in a thread
	fn := func(Iter[int]) Iter[int] {
		return SetError(OfEmpty[int](), anErr)
	}
	it = ReduceToSlice(Parallel(fn)(Of(1, 2)))
	assert.Equal(t, util.Of2Error([]int(nil), anErr), Maybe(it))
}

// // ==== Composition

func TestStreamCompose(t *testing.T) {
	{
		fn := funcs.Compose2(Skip[int](1), Limit[int](3))
		it := fn(OfEmpty[int]())
		assert.Equal(t, util.Of2Error([]int{}, nil), Maybe(ReduceToSlice(it)))
		assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

		it = fn(OfOne(1))
		assert.Equal(t, util.Of2Error([]int{}, nil), Maybe(ReduceToSlice(it)))
		assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

		it = fn(Of(1, 2))
		assert.Equal(t, util.Of2Error([]int{2}, nil), Maybe(ReduceToSlice(it)))
		assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

		it = fn(Of(1, 2, 3))
		assert.Equal(t, util.Of2Error([]int{2, 3}, nil), Maybe(ReduceToSlice(it)))
		assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

		it = fn(Of(1, 2, 3, 4))
		assert.Equal(t, util.Of2Error([]int{2, 3, 4}, nil), Maybe(ReduceToSlice(it)))
		assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

		it = fn(Of(1, 2, 3, 4, 5))
		assert.Equal(t, util.Of2Error([]int{2, 3, 4}, nil), Maybe(ReduceToSlice(it)))
		assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
	}

	{
		fn := funcs.Compose2(Limit[int](3), Skip[int](1))
		it := fn(OfEmpty[int]())
		assert.Equal(t, util.Of2Error([]int{}, nil), Maybe(ReduceToSlice(it)))
		assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

		it = fn(OfOne(1))
		assert.Equal(t, util.Of2Error([]int{}, nil), Maybe(ReduceToSlice(it)))
		assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

		it = fn(Of(1, 2))
		assert.Equal(t, util.Of2Error([]int{2}, nil), Maybe(ReduceToSlice(it)))
		assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

		it = fn(Of(1, 2, 3))
		assert.Equal(t, util.Of2Error([]int{2, 3}, nil), Maybe(ReduceToSlice(it)))
		assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

		it = fn(Of(1, 2, 3, 4))
		assert.Equal(t, util.Of2Error([]int{2, 3}, nil), Maybe(ReduceToSlice(it)))
		assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))

		it = fn(Of(1, 2, 3, 4, 5))
		assert.Equal(t, util.Of2Error([]int{2, 3}, nil), Maybe(ReduceToSlice(it)))
		assert.Equal(t, util.Of2Error(0, EOI), Maybe(it))
	}

	{
		fn := funcs.Compose5(
			Map(strconv.Itoa),
			Map(func(s string) int { i, _ := strconv.Atoi(s); return i }),
			Filter(func(val int) bool { return val&1 == 1 }),
			ReduceToSlice[int],
			Maybe[[]int],
		)

		assert.Equal(t, util.Of2Error([]int{1, 3}, nil), fn(Of(1, 2, 3)))
	}
}
