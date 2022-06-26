// SPDX-License-Identifier: Apache-2.0

package funcs

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceIndex(t *testing.T) {
	slc := []int{}
	assert.Equal(t, 0, SliceIndex(slc, 0))
	assert.Equal(t, 0, SliceIndex(slc, 1))
	assert.Equal(t, 1, SliceIndex(slc, 0, 1))

	slc = []int{1, 2}
	assert.Equal(t, 1, SliceIndex(slc, 0))
	assert.Equal(t, 2, SliceIndex(slc, 1))
	assert.Equal(t, 0, SliceIndex(slc, 2))
	assert.Equal(t, 3, SliceIndex(slc, 2, 3))
}

func TestMapValue(t *testing.T) {
	mp := map[string]int{}
	assert.Equal(t, 0, MapValue(mp, ""))
	assert.Equal(t, 0, MapValue(mp, "a"))
	assert.Equal(t, 3, MapValue(mp, "b", 3))

	mp = map[string]int{"": 1, "a": 2}
	assert.Equal(t, 1, MapValue(mp, ""))
	assert.Equal(t, 2, MapValue(mp, "a"))
	assert.Equal(t, 3, MapValue(mp, "b", 3))
}

func lessThan5(i int) bool {
	return i < 5
}

func lessThan10(i int) bool {
	return i < 10
}

func greaterThan5(i int) bool {
	return i > 5
}

func greaterThan10(i int) bool {
	return i > 10
}

func TestAnd(t *testing.T) {
	lt5_10 := And(lessThan5, lessThan10)
	assert.True(t, lt5_10(3))
	assert.False(t, lt5_10(5))
	assert.False(t, lt5_10(7))
	assert.False(t, lt5_10(10))
	assert.False(t, lt5_10(12))
}

func TestOr(t *testing.T) {
	lt5_gt10 := Or(lessThan5, greaterThan10)
	assert.True(t, lt5_gt10(3))
	assert.False(t, lt5_gt10(5))
	assert.False(t, lt5_gt10(7))
	assert.False(t, lt5_gt10(10))
	assert.True(t, lt5_gt10(12))
}

func TestNot(t *testing.T) {
	nlt5 := Not(lessThan5)
	assert.False(t, nlt5(3))
	assert.True(t, nlt5(5))
	assert.True(t, nlt5(7))
	assert.True(t, nlt5(10))
	assert.True(t, nlt5(12))
}

func TestTernary(t *testing.T) {
	assert.Equal(t, 1, Ternary(1 < 2, 1, 2))
	assert.Equal(t, 1, TernaryResult(1 < 2, func() int { return 1 }, func() int { return 2 }))

	assert.Equal(t, 2, Ternary(1 > 2, 1, 2))
	assert.Equal(t, 2, TernaryResult(1 > 2, func() int { return 1 }, func() int { return 2 }))
}

func TestLessThan(t *testing.T) {
	lt5 := LessThan(5)
	assert.True(t, lt5(3))
	assert.False(t, lt5(5))
	assert.False(t, lt5(7))
	assert.False(t, lt5(10))
	assert.False(t, lt5(12))
}

func TestLessThanEqual(t *testing.T) {
	lte5 := LessThanEqual(5)
	assert.True(t, lte5(3))
	assert.True(t, lte5(5))
	assert.False(t, lte5(7))
	assert.False(t, lte5(10))
	assert.False(t, lte5(12))
}

func TestEqual(t *testing.T) {
	eq5 := Equal(5)
	assert.False(t, eq5(3))
	assert.True(t, eq5(5))
	assert.False(t, eq5(7))
	assert.False(t, eq5(10))
	assert.False(t, eq5(12))
}

func TestGreaterThan(t *testing.T) {
	gt5 := GreaterThan(5)
	assert.False(t, gt5(3))
	assert.False(t, gt5(5))
	assert.True(t, gt5(7))
	assert.True(t, gt5(10))
	assert.True(t, gt5(12))
}

func TestGreaterThanEqual(t *testing.T) {
	gte5 := GreaterThanEqual(5)
	assert.False(t, gte5(3))
	assert.True(t, gte5(5))
	assert.True(t, gte5(7))
	assert.True(t, gte5(10))
	assert.True(t, gte5(12))
}

func TestIsNegative(t *testing.T) {
	neg := IsNegative[int]()
	assert.True(t, neg(-3))
	assert.False(t, neg(0))
	assert.False(t, neg(3))
}

func TestIsNonNegative(t *testing.T) {
	nneg := IsNonNegative[int]()
	assert.False(t, nneg(-3))
	assert.True(t, nneg(0))
	assert.True(t, nneg(3))
}

func TestIsPositive(t *testing.T) {
	pos := IsPositive[int]()
	assert.False(t, pos(-3))
	assert.False(t, pos(0))
	assert.True(t, pos(3))
}

func TestSort(t *testing.T) {
	// Ordered
	{
		slc := []int{2, 3, 1}
		SortOrdered(slc)
		assert.Equal(t, []int{1, 2, 3}, slc)
	}

	// Complex
	{
		slc := []complex64{2, 3, 1}
		SortComplex(slc)
		assert.Equal(t, []complex64{1, 2, 3}, slc)
	}

	// Cmp
	{
		slc := []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(1)}
		SortCmp(slc)
		assert.Equal(t, []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}, slc)
	}

	// By
	{
		slc := []int{2, 3, 1}
		SortBy(slc, func(i, j int) bool { return j < i })
		assert.Equal(t, []int{3, 2, 1}, slc)
	}
}

func TestNillable(t *testing.T) {
	var (
		cn chan int
		c  = make(chan int)
		fn func()
		f  func() = func() {}
		mn map[int]int
		m  map[int]int = map[int]int{}
		i  int         = 0
		pn *int
		p  *int = &i
		sn []int
		s  []int = []int{}
		a  any   = s
	)
	assert.True(t, Nillable(reflect.TypeOf(cn)))
	assert.True(t, IsNil[chan int]()(cn))
	assert.True(t, Nillable(reflect.TypeOf(c)))
	assert.True(t, IsNonNil[chan int]()(c))

	assert.True(t, Nillable(reflect.TypeOf(fn)))
	assert.True(t, IsNil[func()]()(fn))
	assert.True(t, Nillable(reflect.TypeOf(f)))
	assert.True(t, IsNonNil[func()]()(f))

	assert.True(t, Nillable(reflect.TypeOf(mn)))
	assert.True(t, IsNil[map[int]int]()(mn))
	assert.True(t, Nillable(reflect.TypeOf(m)))
	assert.True(t, IsNonNil[map[int]int]()(m))

	assert.True(t, Nillable(reflect.TypeOf(pn)))
	assert.True(t, IsNil[*int]()(pn))
	assert.True(t, Nillable(reflect.TypeOf(p)))
	assert.True(t, IsNonNil[*int]()(p))

	assert.True(t, Nillable(reflect.TypeOf(sn)))
	assert.True(t, IsNil[[]int]()(sn))
	assert.True(t, Nillable(reflect.TypeOf(s)))
	assert.True(t, IsNonNil[[]int]()(s))

	assert.True(t, Nillable(reflect.TypeOf(a)))
	assert.True(t, IsNonNil[[]int]()(a.([]int)))

	assert.False(t, Nillable(reflect.TypeOf(0)))
	func() {
		defer func() {
			assert.Equal(t, fmt.Errorf(notNilableMsg, "int"), recover())
		}()

		IsNil[int]()
		assert.Fail(t, "int cannot be Nillable")
	}()
	func() {
		defer func() {
			assert.Equal(t, fmt.Errorf(notNilableMsg, "int"), recover())
		}()

		IsNonNil[int]()
		assert.Fail(t, "int cannot be Nillable")
	}()
}

func TestMust(t *testing.T) {
	var e error
	Must(e)

	e = fmt.Errorf("bob")
	func() {
		defer func() {
			assert.Equal(t, e, recover())
		}()
		Must(e)
		assert.Fail(t, "Must die")
	}()

	e = nil
	var i int
	assert.Equal(t, i, MustValue(i, e))

	e = fmt.Errorf("bob")
	func() {
		defer func() {
			assert.Equal(t, e, recover())
		}()
		MustValue(i, e)
		assert.Fail(t, "Must die")
	}()
}

func TestSupplier(t *testing.T) {
	supplier := SupplierOf(5)
	assert.Equal(t, 5, supplier())
	assert.Equal(t, 5, supplier())

	var called bool
	supplier = CachingSupplier(func() int { called = true; return 7 })

	assert.False(t, called)
	assert.Equal(t, 7, supplier())
	assert.True(t, called)

	called = false
	assert.False(t, called)
	assert.Equal(t, 7, supplier())
	assert.False(t, called)
}

func TestIgnoreResult(t *testing.T) {
	called := false
	IgnoreResult(func() int { called = true; return 0 })()
	assert.True(t, called)
}

func TestTryTo(t *testing.T) {
	var (
		tryCalled     bool
		panicValue    any
		closersCalled = []int{0}
		theError      = fmt.Errorf("The error")
	)

	TryTo(
		func() { tryCalled = true },
		func(err any) { panicValue = err },
		func() { closersCalled[0] = 1 },
	)
	assert.True(t, tryCalled)
	assert.Nil(t, panicValue)
	assert.Equal(t, 1, closersCalled[0])

	tryCalled, panicValue, closersCalled = false, nil, []int{0}
	TryTo(
		func() { tryCalled = true; panic(theError) },
		func(err any) { panicValue = err },
	)
	assert.True(t, tryCalled)
	assert.Equal(t, theError, panicValue)
	assert.Equal(t, 0, closersCalled[0])

	tryCalled, panicValue, closersCalled = false, nil, []int{}
	TryTo(
		func() { tryCalled = true },
		func(err any) { panicValue = err },
		func() { closersCalled = append(closersCalled, 1) },
		func() { closersCalled = append(closersCalled, 2) },
	)
	assert.True(t, tryCalled)
	assert.Nil(t, panicValue)
	assert.Equal(t, []int{2, 1}, closersCalled)

	tryCalled, panicValue, closersCalled = false, nil, []int{}
	TryTo(
		func() { tryCalled = true; panic(theError) },
		func(err any) { panicValue = err },
		func() { closersCalled = append(closersCalled, 1) },
		func() { closersCalled = append(closersCalled, 2) },
	)
	assert.True(t, tryCalled)
	assert.Equal(t, theError, panicValue)
	assert.Equal(t, []int{2, 1}, closersCalled)
}
