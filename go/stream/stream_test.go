package stream

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	gomath "math"

	"github.com/bantling/micro/go/funcs"
	"github.com/bantling/micro/go/iter"
	"github.com/bantling/micro/go/math"
	"github.com/bantling/micro/go/tuple"
	"github.com/bantling/micro/go/union"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strconv"
	"testing"
)

// ==== Foundation funcs

func TestMap_(t *testing.T) {
	it := Map(strconv.Itoa)(iter.Of(1, 2))
	assert.Equal(t, union.OfResult("1"), iter.Maybe(it))
	assert.Equal(t, union.OfResult("2"), iter.Maybe(it))
	assert.Equal(t, union.OfError[string](iter.EOI), iter.Maybe(it))
}

func TestMapError_(t *testing.T) {
	it := MapError(strconv.Atoi)(iter.Of("1", "2"))
	assert.Equal(t, union.OfResult(1), iter.Maybe(it))
	assert.Equal(t, union.OfResult(2), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	it = MapError(strconv.Atoi)(iter.Of("1", "3.25"))
	assert.Equal(t, union.OfResult(1), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](&strconv.NumError{Func: "Atoi", Num: "3.25", Err: strconv.ErrSyntax}), iter.Maybe(it))
}

func TestFilter_(t *testing.T) {
	it := Filter(func(val int) bool { return val > 1 })(iter.Of(1, 2))
	assert.Equal(t, union.OfResult(2), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))
}

func TestReduce_(t *testing.T) {
	// Reeducer func
	fn := func(i, j int) int { return i + j }

	// No identity, () => sum() = empty
	it := Reduce(fn)(iter.OfEmpty[int]())
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	// No identity, (1) => sum(1) = 1
	it = Reduce(fn)(iter.OfOne(1))
	assert.Equal(t, union.OfResult(1), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	// No identity, (1,  2, 3) => sum(1, 2, 3) = 6
	it = Reduce(fn)(iter.Of(1, 2, 3))
	assert.Equal(t, union.OfResult(6), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	// Identity 4, () => sum(4) = 4
	it = Reduce(fn, 4)(iter.OfEmpty[int]())
	assert.Equal(t, union.OfResult(4), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	// Identity 4, (1) => sum(4, 1) = 5
	it = Reduce(fn, 4)(iter.OfOne(1))
	assert.Equal(t, union.OfResult(5), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	// Identity 4, (1, 2, 3) => sum(4, 1, 2, 3) = 10
	it = Reduce(fn, 4)(iter.Of(1, 2, 3))
	assert.Equal(t, union.OfResult(10), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	// test non-iter.EOI errors
	{
		var (
			anErr = fmt.Errorf("An err")
			fn    = func(i, j int) int { return i + j }
		)
		// no identity, error after first element
		it := Reduce(fn)(iter.SetError(iter.Of(1), anErr))
		assert.Equal(t, union.OfError[int](anErr), iter.Maybe(it))

		// identity, error before first element
		it = Reduce(fn, 0)(iter.SetError(iter.OfEmpty[int](), anErr))
		assert.Equal(t, union.OfError[int](anErr), iter.Maybe(it))

		// identity, error after first element
		it = Reduce(fn, 0)(iter.SetError(iter.Of(1), anErr))
		assert.Equal(t, union.OfError[int](anErr), iter.Maybe(it))
	}
}

func TestReduceTo_(t *testing.T) {
	// No identity, () => concat() = ()
	fn := func(i string, j int) string { return i + strconv.Itoa(j) }
	it := ReduceTo(fn)(iter.OfEmpty[int]())
	assert.Equal(t, union.OfError[string](iter.EOI), iter.Maybe(it))

	// No identity, (1) => concat(1) = ("1")
	it = ReduceTo(fn)(iter.OfOne(1))
	assert.Equal(t, union.OfResult("1"), iter.Maybe(it))
	assert.Equal(t, union.OfError[string](iter.EOI), iter.Maybe(it))

	// No identity, (1, 2, 3) => concat(1, 2, 3) = "123"
	it = ReduceTo(fn)(iter.Of(1, 2, 3))
	assert.Equal(t, union.OfResult("123"), iter.Maybe(it))
	assert.Equal(t, union.OfError[string](iter.EOI), iter.Maybe(it))

	// Identity "4", () => concat(4) = "4"
	it = ReduceTo(fn, "4")(iter.OfEmpty[int]())
	assert.Equal(t, union.OfResult("4"), iter.Maybe(it))
	assert.Equal(t, union.OfError[string](iter.EOI), iter.Maybe(it))

	// Identity "4", (1) => concat(4, 1) = "41"
	it = ReduceTo(fn, "4")(iter.OfOne(1))
	assert.Equal(t, union.OfResult("41"), iter.Maybe(it))
	assert.Equal(t, union.OfError[string](iter.EOI), iter.Maybe(it))

	// Identity "4", (1, 2, 3) => concat(4, 1, 2, 3) = "4123"
	it = ReduceTo(fn, "4")(iter.Of(1, 2, 3))
	assert.Equal(t, union.OfResult("4123"), iter.Maybe(it))
	assert.Equal(t, union.OfError[string](iter.EOI), iter.Maybe(it))

	// test non-iter.EOI errors
	{
		var (
			anErr = fmt.Errorf("An err")
			fn    = func(i, j int) int { return i + j }
		)
		// no identity, error after first element
		it := ReduceTo(fn)(iter.SetError(iter.Of(1), anErr))
		assert.Equal(t, union.OfError[int](anErr), iter.Maybe(it))

		// identity, error before first element
		it = ReduceTo(fn, 0)(iter.SetError(iter.OfEmpty[int](), anErr))
		assert.Equal(t, union.OfError[int](anErr), iter.Maybe(it))

		// identity, error after first element
		it = ReduceTo(fn, 0)(iter.SetError(iter.Of(1), anErr))
		assert.Equal(t, union.OfError[int](anErr), iter.Maybe(it))
	}
}

func TestReduceToBool_(t *testing.T) {
	// And logic (all match): identity = true, stop on false
	it := ReduceToBool(func(i int) bool { return i < 3 }, true, false)(iter.OfEmpty[int]())
	assert.Equal(t, union.OfResult(true), iter.Maybe(it))
	assert.Equal(t, union.OfError[bool](iter.EOI), iter.Maybe(it))

	it = ReduceToBool(func(i int) bool { return i < 3 }, true, false)(iter.Of(1, 2, 3))
	assert.Equal(t, union.OfResult(false), iter.Maybe(it))
	assert.Equal(t, union.OfError[bool](iter.EOI), iter.Maybe(it))

	// Or logic (at least one match): identity = false, stop on true
	it = ReduceToBool(func(i int) bool { return i < 3 }, false, true)(iter.Of(4))
	assert.Equal(t, union.OfResult(false), iter.Maybe(it))
	assert.Equal(t, union.OfError[bool](iter.EOI), iter.Maybe(it))

	it = ReduceToBool(func(i int) bool { return i < 3 }, false, true)(iter.Of(1))
	assert.Equal(t, union.OfResult(true), iter.Maybe(it))
	assert.Equal(t, union.OfError[bool](iter.EOI), iter.Maybe(it))

	// test non-iter.EOI errors
	{
		var (
			anErr = fmt.Errorf("An err")
			fn    = func(i int) bool { return i > 0 }
		)
		// error before first element
		it := ReduceToBool(fn, true, false)(iter.SetError(iter.OfEmpty[int](), anErr))
		assert.Equal(t, union.OfError[bool](anErr), iter.Maybe(it))

		// error after first element
		it = ReduceToBool(fn, true, false)(iter.SetError(iter.Of(1), anErr))
		assert.Equal(t, union.OfError[bool](anErr), iter.Maybe(it))
	}
}

func TestReduceToSlice_(t *testing.T) {
	// Reduce into a new generated slice, no error
	it := ReduceToSlice(iter.Of(1, 2))
	assert.Equal(t, union.OfResult([]int{1, 2}), iter.Maybe(it))
	assert.Equal(t, union.OfError[[]int](iter.EOI), iter.Maybe(it))

	// Error before first element
	{
		anErr := fmt.Errorf("An err")
		it := ReduceToSlice(iter.SetError(iter.OfEmpty[int](), anErr))
		assert.Equal(t, union.OfError[[]int](anErr), iter.Maybe(it))
	}
}

func TestReduceIntoSlice_(t *testing.T) {
	// Reduce into a new generated slice, no error
	{
		slc := make([]int, 2)
		it := ReduceIntoSlice(slc)(iter.Of(1, 2))
		assert.Equal(t, union.OfResult([]int{1, 2}), iter.Maybe(it))
		assert.Equal(t, union.OfError[[]int](iter.EOI), iter.Maybe(it))
	}

	// Error before first element
	{
		slc := []int{2}
		anErr := fmt.Errorf("An err")
		it := ReduceIntoSlice(slc)(iter.SetError(iter.OfEmpty[int](), anErr))
		assert.Equal(t, union.OfError[[]int](anErr), iter.Maybe(it))
		// ReduceIntoSlice sets every element of target slice to zero val if an error occurs
		assert.Equal(t, []int{0}, slc)
	}

	// Panics if target slice is not large enough
	{
		slc := []int{1}
		it := ReduceIntoSlice(slc)(iter.Of(1, 2))

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

func TestExpandSlices_(t *testing.T) {
	it := ReduceToSlice(ExpandSlices(iter.Of([]int{1, 2, 3}, nil, []int{}, []int{4, 5})))
	assert.Equal(t, union.OfResult([]int{1, 2, 3, 4, 5}), iter.Maybe(it))
	assert.Equal(t, union.OfError[[]int](iter.EOI), iter.Maybe(it))

	// Error before first element
	{
		anErr := fmt.Errorf("An err")
		it := ExpandSlices(iter.SetError(iter.OfEmpty[[]int](), anErr))
		assert.Equal(t, union.OfError[int](anErr), iter.Maybe(it))
	}
}

func TestReduceToMap_(t *testing.T) {
	it := ReduceToMap(iter.Of(tuple.Of2(1, "1"), tuple.Of2(2, "2"), tuple.Of2(3, "3")))
	assert.Equal(t, union.OfResult(map[int]string{1: "1", 2: "2", 3: "3"}), iter.Maybe(it))
	assert.Equal(t, union.OfError[map[int]string](iter.EOI), iter.Maybe(it))

	{
		anErr := fmt.Errorf("An err")
		it := ReduceToMap(iter.SetError(iter.OfEmpty[tuple.Two[int, string]](), anErr))
		assert.Equal(t, union.OfError[map[int]string](anErr), iter.Maybe(it))
	}
}

func TestExpandMaps_(t *testing.T) {
	it := ReduceToMap(ExpandMaps(iter.Of(map[int]string{1: "1", 2: "2"}, nil, map[int]string{}, map[int]string{3: "3"})))
	assert.Equal(t, union.OfResult(map[int]string{1: "1", 2: "2", 3: "3"}), iter.Maybe(it))
	assert.Equal(t, union.OfError[map[int]string](iter.EOI), iter.Maybe(it))

	{
		anErr := fmt.Errorf("An err")
		it := ExpandMaps(iter.SetError(iter.OfEmpty[map[int]int](), anErr))
		assert.Equal(t, union.OfError[tuple.Two[int, int]](anErr), iter.Maybe(it))
	}
}

func TestSkip_(t *testing.T) {
	fn := Skip[int](3)
	it := fn(iter.OfEmpty[int]())
	assert.Equal(t, union.OfResult([]int{}), iter.Maybe(ReduceToSlice(it)))

	it = fn(iter.OfOne(1))
	assert.Equal(t, union.OfResult([]int{}), iter.Maybe(ReduceToSlice(it)))

	it = fn(iter.Of(1, 2))
	assert.Equal(t, union.OfResult([]int{}), iter.Maybe(ReduceToSlice(it)))

	it = fn(iter.Of(1, 2, 3))
	assert.Equal(t, union.OfResult([]int{}), iter.Maybe(ReduceToSlice(it)))

	it = fn(iter.Of(1, 2, 3, 4))
	assert.Equal(t, union.OfResult([]int{4}), iter.Maybe(ReduceToSlice(it)))

	it = fn(iter.Of(1, 2, 3, 4, 5))
	assert.Equal(t, union.OfResult([]int{4, 5}), iter.Maybe(ReduceToSlice(it)))
}

func TestLimit_(t *testing.T) {
	fn := Limit[int](3)
	it := fn(iter.OfEmpty[int]())
	assert.Equal(t, union.OfResult([]int{}), iter.Maybe(ReduceToSlice(it)))

	it = fn(iter.OfOne(1))
	assert.Equal(t, union.OfResult([]int{1}), iter.Maybe(ReduceToSlice(it)))

	it = fn(iter.Of(1, 2))
	assert.Equal(t, union.OfResult([]int{1, 2}), iter.Maybe(ReduceToSlice(it)))

	it = fn(iter.Of(1, 2, 3))
	assert.Equal(t, union.OfResult([]int{1, 2, 3}), iter.Maybe(ReduceToSlice(it)))

	it = fn(iter.Of(1, 2, 3, 4))
	assert.Equal(t, union.OfResult([]int{1, 2, 3}), iter.Maybe(ReduceToSlice(it)))

	it = fn(iter.Of(1, 2, 3, 4, 5))
	assert.Equal(t, union.OfResult([]int{1, 2, 3}), iter.Maybe(ReduceToSlice(it)))
}

func TestPeek_(t *testing.T) {
	slc := []int{}
	fn := Peek(func(val int) { slc = append(slc, val) })
	it := fn(iter.OfEmpty[int]())
	// assertNext(t, 0, iter.EOI)(it.Next())
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))
	assert.Equal(t, []int{}, slc)

	it = fn(iter.OfOne(1))
	assert.Equal(t, union.OfResult(1), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))
	assert.Equal(t, []int{1}, slc)

	it = fn(iter.Of(2, 3, 4))
	assert.Equal(t, union.OfResult(2), iter.Maybe(it))
	assert.Equal(t, union.OfResult(3), iter.Maybe(it))
	assert.Equal(t, union.OfResult(4), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))
	assert.Equal(t, []int{1, 2, 3, 4}, slc)
}

func TestGenerator_(t *testing.T) {
	called := 0
	fn := Generator(func() func(iter.Iter[int]) iter.Iter[int] {
		return func(it iter.Iter[int]) iter.Iter[int] {
			called++
			return it
		}
	})

	it := fn(iter.Of(1))
	assert.Equal(t, union.OfResult(1), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))
	assert.Equal(t, 1, called)

	it = fn(iter.Of(3))
	assert.Equal(t, union.OfResult(3), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))
	assert.Equal(t, 2, called)
}

// // ==== Funcs based on foundational funcs

func TestAllMatch_(t *testing.T) {
	fn := AllMatch(func(i int) bool { return i < 3 })
	it := fn(iter.OfEmpty[int]())
	assert.Equal(t, union.OfResult(true), iter.Maybe(it))
	assert.Equal(t, union.OfError[bool](iter.EOI), iter.Maybe(it))

	it = fn(iter.Of(1))
	assert.Equal(t, union.OfResult(true), iter.Maybe(it))
	assert.Equal(t, union.OfError[bool](iter.EOI), iter.Maybe(it))

	it = fn(iter.Of(1, 2, 3))
	assert.Equal(t, union.OfResult(false), iter.Maybe(it))
	assert.Equal(t, union.OfError[bool](iter.EOI), iter.Maybe(it))
}

func TestAnyMatch_(t *testing.T) {
	fn := AnyMatch(func(i int) bool { return i < 3 })
	it := fn(iter.OfEmpty[int]())
	assert.Equal(t, union.OfResult(false), iter.Maybe(it))
	assert.Equal(t, union.OfError[bool](iter.EOI), iter.Maybe(it))

	it = fn(iter.Of(1))
	assert.Equal(t, union.OfResult(true), iter.Maybe(it))
	assert.Equal(t, union.OfError[bool](iter.EOI), iter.Maybe(it))

	it = fn(iter.Of(4, 5, 6))
	assert.Equal(t, union.OfResult(false), iter.Maybe(it))
	assert.Equal(t, union.OfError[bool](iter.EOI), iter.Maybe(it))
}

func TestNoneMatch_(t *testing.T) {
	fn := NoneMatch(func(i int) bool { return i < 3 })
	it := fn(iter.OfEmpty[int]())
	assert.Equal(t, union.OfResult(true), iter.Maybe(it))
	assert.Equal(t, union.OfError[bool](iter.EOI), iter.Maybe(it))

	it = fn(iter.Of(1))
	assert.Equal(t, union.OfResult(false), iter.Maybe(it))
	assert.Equal(t, union.OfError[bool](iter.EOI), iter.Maybe(it))

	it = fn(iter.Of(4, 5, 6))
	assert.Equal(t, union.OfResult(true), iter.Maybe(it))
	assert.Equal(t, union.OfError[bool](iter.EOI), iter.Maybe(it))
}

func TestCount_(t *testing.T) {
	it := Count(iter.OfEmpty[int]())
	assert.Equal(t, union.OfResult(0), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	it = Count(iter.Of(1))
	assert.Equal(t, union.OfResult(1), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	it = Count(iter.Of(1, 2, 3))
	assert.Equal(t, union.OfResult(3), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))
}

func TestDistinct_(t *testing.T) {
	// Distinct
	it := Distinct(iter.OfEmpty[int]())
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	it = Distinct(iter.OfOne(1))
	assert.Equal(t, union.OfResult(1), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	it = Distinct(iter.Of(1, 3, 2, 3, 2, 1))
	assert.Equal(t, union.OfResult([]int{1, 3, 2}), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	// Distinct across multiple iters via iter.Concat
	it = Distinct(iter.Concat(iter.OfEmpty[int](), iter.OfOne(1), iter.Of(1, 2, 3, 3, 2, 1), iter.Of(1, 4)))
	assert.Equal(t, union.OfResult([]int{1, 2, 3, 4}), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))
}

func TestDuplicate_(t *testing.T) {
	// Duplicate
	it := Duplicate(iter.OfEmpty[int]())
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	it = Duplicate(iter.OfOne(1))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	it = Duplicate(iter.Of(1, 2, 3, 3, 2))
	assert.Equal(t, union.OfResult([]int{3, 2}), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	// Duplicate across multiple iters via iter.Concat
	it = Duplicate(iter.Concat(iter.OfEmpty[int](), iter.OfOne(1), iter.Of(1, 2, 3, 3, 2), iter.Of(1, 4)))
	assert.Equal(t, union.OfResult([]int{1, 3, 2}), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))
}

func TestReverse_(t *testing.T) {
	it := Reverse(iter.OfEmpty[int]())
	assert.Equal(t, union.OfResult([]int{}), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	it = Reverse(iter.OfOne(1))
	assert.Equal(t, union.OfResult([]int{1}), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	it = Reverse(iter.Of(1, 2))
	assert.Equal(t, union.OfResult([]int{2, 1}), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	it = Reverse(iter.Of(1, 2, 3))
	assert.Equal(t, union.OfResult([]int{3, 2, 1}), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	it = Reverse(iter.Of(1, 2, 3, 4))
	assert.Equal(t, union.OfResult([]int{4, 3, 2, 1}), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	anErr := fmt.Errorf("An err")
	it = Reverse(iter.SetError(iter.OfEmpty[int](), anErr))
	assert.Equal(t, union.OfError[[]int](anErr), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[int](anErr), iter.Maybe(it))
}

func TestSortOrdered_(t *testing.T) {
	it := SortOrdered(iter.Of(1, 3, 2))
	assert.Equal(t, union.OfResult([]int{1, 2, 3}), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	anErr := fmt.Errorf("An err")
	it = SortOrdered(iter.SetError(iter.OfEmpty[int](), anErr))
	assert.Equal(t, union.OfError[[]int](anErr), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[int](anErr), iter.Maybe(it))
}

func TestSortComplex_(t *testing.T) {
	it := SortComplex(iter.Of(1+0i, 3+1i, 2+0i))
	assert.Equal(t, union.OfResult([]complex128{1 + 0i, 2 + 0i, 3 + 1i}), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[complex128](iter.EOI), iter.Maybe(it))

	anErr := fmt.Errorf("An err")
	it = SortComplex(iter.SetError(iter.OfEmpty[complex128](), anErr))
	assert.Equal(t, union.OfError[[]complex128](anErr), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[complex128](anErr), iter.Maybe(it))
}

func TestSortCmp_(t *testing.T) {
	it := SortCmp(iter.Of(big.NewInt(2), big.NewInt(3), big.NewInt(1)))
	assert.Equal(t, union.OfResult([]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[*big.Int](iter.EOI), iter.Maybe(it))

	anErr := fmt.Errorf("An err")
	it = SortCmp(iter.SetError(iter.OfEmpty[*big.Int](), anErr))
	assert.Equal(t, union.OfError[[]*big.Int](anErr), iter.Maybe(ReduceToSlice(it)))
}

func TestSortBy_(t *testing.T) {
	fn := SortBy(func(i, j int) bool { return j < i })
	it := fn(iter.Of(1, 3, 2))
	assert.Equal(t, union.OfResult([]int{3, 2, 1}), iter.Maybe(ReduceToSlice(it)))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	anErr := fmt.Errorf("An err")
	it = fn(iter.SetError(iter.OfEmpty[int](), anErr))
	assert.Equal(t, union.OfError[[]int](anErr), iter.Maybe(ReduceToSlice(it)))
}

// ==== Math

func TestAbs_(t *testing.T) {
	it := Abs(iter.Of(-1, 5, gomath.MinInt))
	assert.Equal(t, union.OfResult(1), iter.Maybe(it))
	assert.Equal(t, union.OfResult(5), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](fmt.Errorf("Absolute value error for -9223372036854775808: there is no corresponding positive value in type int")), iter.Maybe(it))
}

func TestAbsBigOps_(t *testing.T) {
	it := AbsBigOps(iter.Of(big.NewInt(-1), big.NewInt(5)))
	assert.Equal(t, union.OfResult(big.NewInt(1)), iter.Maybe(it))
	assert.Equal(t, union.OfResult(big.NewInt(5)), iter.Maybe(it))
}

func TestAvgInt_(t *testing.T) {
	it := AvgInt(iter.Of(-1, 5))
	assert.Equal(t, union.OfResult(2), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	it = AvgInt(iter.Of(1, gomath.MaxInt))
	assert.Equal(t, union.OfError[int](math.OverflowErr), iter.Maybe(it))

	it = AvgInt(iter.OfEmpty[int]())
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	err := fmt.Errorf("An err")
	it = AvgInt(iter.SetError(iter.Of(1), err))
	assert.Equal(t, union.OfError[int](err), iter.Maybe(it))
}

func TestAvgUint_(t *testing.T) {
	it := AvgUint(iter.Of(uint(1), uint(5)))
	assert.Equal(t, union.OfResult(uint(3)), iter.Maybe(it))
	assert.Equal(t, union.OfError[uint](iter.EOI), iter.Maybe(it))

	it = AvgUint(iter.Of(uint(1), gomath.MaxUint))
	assert.Equal(t, union.OfError[uint](math.OverflowErr), iter.Maybe(it))

	it = AvgUint(iter.OfEmpty[uint]())
	assert.Equal(t, union.OfError[uint](iter.EOI), iter.Maybe(it))

	err := fmt.Errorf("An err")
	it = AvgUint(iter.SetError(iter.Of(uint(1)), err))
	assert.Equal(t, union.OfError[uint](err), iter.Maybe(it))
}

func TestAvgBigOps_(t *testing.T) {
	it := AvgBigOps(iter.Of(big.NewInt(-1), big.NewInt(5)))
	assert.Equal(t, union.OfResult(big.NewInt(2)), iter.Maybe(it))
	assert.Equal(t, union.OfError[*big.Int](iter.EOI), iter.Maybe(it))

	it = AvgBigOps(iter.OfEmpty[*big.Int]())
	assert.Equal(t, union.OfError[*big.Int](iter.EOI), iter.Maybe(it))

	err := fmt.Errorf("An err")
	it = AvgBigOps(iter.SetError(iter.Of(big.NewInt(1)), err))
	assert.Equal(t, union.OfError[*big.Int](err), iter.Maybe(it))
}

func TestMax_(t *testing.T) {
	it := Max(iter.Of(2, 3, 1))
	assert.Equal(t, union.OfResult(3), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	it = Max(iter.OfEmpty[int]())
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	err := fmt.Errorf("An err")
	it = Max(iter.SetError(iter.Of(1), err))
	assert.Equal(t, union.OfError[int](err), iter.Maybe(it))
}

func TestMaxCmp_(t *testing.T) {
	it := MaxCmp(iter.Of(big.NewInt(2), big.NewInt(3), big.NewInt(1)))
	assert.Equal(t, union.OfResult(big.NewInt(3)), iter.Maybe(it))
	assert.Equal(t, union.OfError[*big.Int](iter.EOI), iter.Maybe(it))

	it = MaxCmp(iter.OfEmpty[*big.Int]())
	assert.Equal(t, union.OfError[*big.Int](iter.EOI), iter.Maybe(it))

	err := fmt.Errorf("An err")
	it = MaxCmp(iter.SetError(iter.Of(big.NewInt(1)), err))
	assert.Equal(t, union.OfError[*big.Int](err), iter.Maybe(it))
}

func TestMin_(t *testing.T) {
	it := Min(iter.Of(2, 3, 1))
	assert.Equal(t, union.OfResult(1), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	it = Min(iter.OfEmpty[int]())
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

	err := fmt.Errorf("An err")
	it = Min(iter.SetError(iter.Of(1), err))
	assert.Equal(t, union.OfError[int](err), iter.Maybe(it))
}

func TestMinCmp_(t *testing.T) {
	it := MinCmp(iter.Of(big.NewInt(2), big.NewInt(3), big.NewInt(1)))
	assert.Equal(t, union.OfResult(big.NewInt(1)), iter.Maybe(it))
	assert.Equal(t, union.OfError[*big.Int](iter.EOI), iter.Maybe(it))

	it = MinCmp(iter.OfEmpty[*big.Int]())
	assert.Equal(t, union.OfError[*big.Int](iter.EOI), iter.Maybe(it))

	err := fmt.Errorf("An err")
	it = MinCmp(iter.SetError(iter.Of(big.NewInt(1)), err))
	assert.Equal(t, union.OfError[*big.Int](err), iter.Maybe(it))
}

func TestSum_(t *testing.T) {
	it := Sum(iter.Of(-1, 5))
	assert.Equal(t, union.OfResult(4), iter.Maybe(it))
	assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))
}

func TestSumBigOps_(t *testing.T) {
	it := SumBigOps(iter.Of(big.NewInt(-1), big.NewInt(5)))
	assert.Equal(t, union.OfResult(big.NewInt(4)), iter.Maybe(it))
	assert.Equal(t, union.OfError[*big.Int](iter.EOI), iter.Maybe(it))
}

func TestGenerateRanges_(t *testing.T) {
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

func TestParallel_(t *testing.T) {
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
		assert.Equal(t, union.OfResult(outInt), iter.Maybe(ReduceToSlice(pIntSqrt(iter.Of(inInt...)))))
		assert.Equal(t, union.OfResult(outInt), iter.Maybe(ReduceToSlice(pIntThreads(iter.Of(inInt...)))))
		assert.Equal(t, union.OfResult(outInt), iter.Maybe(ReduceToSlice(pIntItems(iter.Of(inInt...)))))

		// Different type, generate a new slice
		assert.Equal(t, union.OfResult(outUint), iter.Maybe(ReduceToSlice(pUintSqrt(iter.Of(inInt...)))))
		assert.Equal(t, union.OfResult(outUint), iter.Maybe(ReduceToSlice(pUintThreads(iter.Of(inInt...)))))
		assert.Equal(t, union.OfResult(outUint), iter.Maybe(ReduceToSlice(pUintItems(iter.Of(inInt...)))))
	}

	// Error on source iter
	anErr := fmt.Errorf("An err")
	it := ReduceToSlice(pIntSqrt(iter.SetError(iter.OfEmpty[int](), anErr)))
	assert.Equal(t, union.OfError[[]int](anErr), iter.Maybe(it))

	// Error during transform that executes in a thread
	fn := func(iter.Iter[int]) iter.Iter[int] {
		return iter.SetError(iter.OfEmpty[int](), anErr)
	}
	it = ReduceToSlice(Parallel(fn)(iter.Of(1, 2)))
	assert.Equal(t, union.OfError[[]int](anErr), iter.Maybe(it))
}

// // ==== Composition

func TestStreamCompose_(t *testing.T) {
	{
		fn := funcs.Compose2(Skip[int](1), Limit[int](3))
		it := fn(iter.OfEmpty[int]())
		assert.Equal(t, union.OfResult([]int{}), iter.Maybe(ReduceToSlice(it)))
		assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

		it = fn(iter.OfOne(1))
		assert.Equal(t, union.OfResult([]int{}), iter.Maybe(ReduceToSlice(it)))
		assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

		it = fn(iter.Of(1, 2))
		assert.Equal(t, union.OfResult([]int{2}), iter.Maybe(ReduceToSlice(it)))
		assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

		it = fn(iter.Of(1, 2, 3))
		assert.Equal(t, union.OfResult([]int{2, 3}), iter.Maybe(ReduceToSlice(it)))
		assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

		it = fn(iter.Of(1, 2, 3, 4))
		assert.Equal(t, union.OfResult([]int{2, 3, 4}), iter.Maybe(ReduceToSlice(it)))
		assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

		it = fn(iter.Of(1, 2, 3, 4, 5))
		assert.Equal(t, union.OfResult([]int{2, 3, 4}), iter.Maybe(ReduceToSlice(it)))
		assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))
	}

	{
		fn := funcs.Compose2(Limit[int](3), Skip[int](1))
		it := fn(iter.OfEmpty[int]())
		assert.Equal(t, union.OfResult([]int{}), iter.Maybe(ReduceToSlice(it)))
		assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

		it = fn(iter.OfOne(1))
		assert.Equal(t, union.OfResult([]int{}), iter.Maybe(ReduceToSlice(it)))
		assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

		it = fn(iter.Of(1, 2))
		assert.Equal(t, union.OfResult([]int{2}), iter.Maybe(ReduceToSlice(it)))
		assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

		it = fn(iter.Of(1, 2, 3))
		assert.Equal(t, union.OfResult([]int{2, 3}), iter.Maybe(ReduceToSlice(it)))
		assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

		it = fn(iter.Of(1, 2, 3, 4))
		assert.Equal(t, union.OfResult([]int{2, 3}), iter.Maybe(ReduceToSlice(it)))
		assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))

		it = fn(iter.Of(1, 2, 3, 4, 5))
		assert.Equal(t, union.OfResult([]int{2, 3}), iter.Maybe(ReduceToSlice(it)))
		assert.Equal(t, union.OfError[int](iter.EOI), iter.Maybe(it))
	}

	{
		fn := funcs.Compose5(
			Map(strconv.Itoa),
			Map(func(s string) int { i, _ := strconv.Atoi(s); return i }),
			Filter(func(val int) bool { return val&1 == 1 }),
			ReduceToSlice[int],
			iter.Maybe[[]int],
		)

		assert.Equal(t, union.OfResult([]int{1, 3}), fn(iter.Of(1, 2, 3)))
	}
}
