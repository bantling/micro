package funcs

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"testing"

	"github.com/bantling/micro/tuple"
	"github.com/stretchr/testify/assert"
)

// ==== Slices

func TestSliceAdd_(t *testing.T) {
	var slc []int
	SliceAdd(&slc, 3)
	assert.NotNil(t, slc)
	assert.Equal(t, []int{3}, slc)

	SliceAdd(&slc, 4)
	assert.Equal(t, []int{3, 4}, slc)
}

func TestSliceCopy_(t *testing.T) {
	assert.Equal(t, []int{}, SliceCopy([]int(nil)))

	var (
		slc  = []int{1, 2}
		slc2 = SliceCopy(slc)
	)
	assert.NotEqual(t, fmt.Sprintf("%p", slc), fmt.Sprintf("%p", slc2))
	assert.Equal(t, slc, slc2)
}

func TestSliceFlatten_(t *testing.T) {
	assert.Equal(t, []int{}, SliceFlatten[int](nil))

	// Check that one dimensional slice is returned as is (same address)
	oneDim := []int{}
	assert.Equal(t, fmt.Sprintf("%p", oneDim), fmt.Sprintf("%p", SliceFlatten[int](oneDim)))

	assert.Equal(t, []int{1, 2, 3, 4}, SliceFlatten[int]([][]int{{1, 2}, {3, 4}}))

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, SliceFlatten[int]([][][]int{{{1, 2}, {3, 4}}, {{5}}, {{6}}}))

	// Die if a value that is not a slice is passed
	TryTo(
		func() {
			SliceFlatten[int](0)
			assert.Fail(t, "Must die")
		},
		func(err any) {
			assert.Equal(t, fmt.Errorf("SliceFlatten argument must be a slice, not type int"), err)
		},
	)

	// Die if expecting a []int but passed a []string
	TryTo(
		func() {
			SliceFlatten[int]([]string{})
			assert.Fail(t, "Must die")
		},
		func(err any) {
			assert.Equal(t, fmt.Errorf("SliceFlatten argument must be slice of int, not a slice of string"), err)
		},
	)
}

func TestSliceIndex_(t *testing.T) {
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

func TestSliceOf_(t *testing.T) {
	assert.Equal(t, []string{"a"}, SliceOf("a"))
	assert.Equal(t, []int{1, 2}, SliceOf(1, 2))
}

func TestSliceRemove_(t *testing.T) {
	assert.Equal(t, []int{}, SliceRemove([]int{}, 0))
	assert.Equal(t, []int{1, 2, 4}, SliceRemove([]int{1, 2, 3, 4}, 3))
	assert.Equal(t, []int{1, 2, 4, 3, 3}, SliceRemove([]int{1, 2, 3, 4, 3, 3}, 3))
	assert.Equal(t, []int{1, 2, 4}, SliceRemove([]int{1, 2, 3, 4, 3, 3}, 3, true))
}

type Uncomparable[T any] interface {
	Op(T) T
}

type UncomparableFunc[T any] func(T) T

func (u UncomparableFunc[T]) Op(t T) T {
	return u(t)
}

func TestSliceRemoveUncomparable_(t *testing.T) {
	var (
		f1  Uncomparable[int] = UncomparableFunc[int](func(i int) int { return i + 1 })
		f2  Uncomparable[int] = UncomparableFunc[int](func(i int) int { return i + i })
		f3  Uncomparable[int] = UncomparableFunc[int](func(i int) int { return i + i })
		slc                   = []Uncomparable[int]{f1, f2, f1, f2}
	)

	assert.Equal(t, 3, len(SliceRemoveUncomparable(slc, f1)))       // One removed
	assert.Equal(t, 4, len(SliceRemoveUncomparable(slc, f3)))       // None removed
	assert.Equal(t, 2, len(SliceRemoveUncomparable(slc, f2, true))) // Two removed
}

func TestSliceReverse_(t *testing.T) {
	slc := []int{}
	assert.Equal(t, []int{}, SliceReverse(slc))
	assert.Equal(t, []int{}, slc)

	slc = []int{1}
	assert.Equal(t, []int{1}, SliceReverse(slc))
	assert.Equal(t, []int{1}, slc)

	slc = []int{1, 2}
	assert.Equal(t, []int{2, 1}, SliceReverse(slc))
	assert.Equal(t, []int{2, 1}, slc)

	slc = []int{1, 2, 3}
	assert.Equal(t, []int{3, 2, 1}, SliceReverse(slc))
	assert.Equal(t, []int{3, 2, 1}, slc)

	slc = []int{1, 2, 3, 4}
	assert.Equal(t, []int{4, 3, 2, 1}, SliceReverse(slc))
	assert.Equal(t, []int{4, 3, 2, 1}, slc)
}

func TestSliceSort_(t *testing.T) {
	// Ordered
	{
		slc := []int{2, 3, 1}
		assert.Equal(t, []int{1, 2, 3}, SliceSortOrdered(slc))
		assert.Equal(t, []int{1, 2, 3}, slc)
	}

	// Complex
	{
		slc := []complex64{2, 3, 1}
		assert.Equal(t, []complex64{1, 2, 3}, SliceSortComplex(slc))
		assert.Equal(t, []complex64{1, 2, 3}, slc)
	}

	// Cmp
	{
		slc := []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(1)}
		assert.Equal(t, []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}, SliceSortCmp(slc))
		assert.Equal(t, []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}, slc)
	}

	// By
	{
		slc := []int{2, 3, 1}
		assert.Equal(t, []int{3, 2, 1}, SliceSortBy(slc, func(i, j int) bool { return j < i }))
		assert.Equal(t, []int{3, 2, 1}, slc)
	}
}

func TestSliceToMap_(t *testing.T) {
	// ToMap
	assert.Equal(t, map[int]bool{}, SliceToMap([]int(nil)))
	assert.Equal(t, map[int]bool{}, SliceToMap([]int{}))
	assert.Equal(t, map[int]bool{1: true}, SliceToMap([]int{1}))
	assert.Equal(t, map[int]bool{1: true, 2: true, 3: true}, SliceToMap([]int{3, 1, 2}))

	// ToMapBy
	fn := func(i int) string { return fmt.Sprintf("%d", i) }
	assert.Equal(t, map[string]bool{}, SliceToMapBy([]int(nil), fn))
	assert.Equal(t, map[string]bool{}, SliceToMapBy([]int{}, fn))
	assert.Equal(t, map[string]bool{"1": true}, SliceToMapBy([]int{1}, fn))
	assert.Equal(t, map[string]bool{"1": true, "2": true, "3": true}, SliceToMapBy([]int{3, 1, 2}, fn))
}

func TestSliceUniqueValues_(t *testing.T) {
	slc := []int{1, 2, 3, 2, 1}
	assert.Equal(t, []int{1, 2, 3}, SliceSortOrdered(SliceUniqueValues(slc)))
	assert.Equal(t, []int{1, 2, 3, 2, 1}, slc)
}

func TestMapIndex_(t *testing.T) {
	mp := map[string]int{}
	assert.Equal(t, 0, MapIndex(mp, ""))
	assert.Equal(t, 0, MapIndex(mp, "a"))
	assert.Equal(t, 3, MapIndex(mp, "b", 3))

	mp = map[string]int{"": 1, "a": 2}
	assert.Equal(t, 1, MapIndex(mp, ""))
	assert.Equal(t, 2, MapIndex(mp, "a"))
	assert.Equal(t, 3, MapIndex(mp, "b", 3))
}

func TestMapSet_(t *testing.T) {
	{
		var mp map[string]int
		MapSet(&mp, "foo", 3)
		assert.NotNil(t, mp)
		assert.Equal(t, map[string]int{"foo": 3}, mp)

		MapSet(&mp, "bar", 4)
		assert.Equal(t, map[string]int{"foo": 3, "bar": 4}, mp)
	}

	{
		var mp map[string]map[int]bool
		Map2Set(&mp, "foo", 3, true)
		assert.NotNil(t, mp)
		assert.NotNil(t, mp["foo"])
		assert.Equal(t, map[string]map[int]bool{"foo": {3: true}}, mp)

		Map2Set(&mp, "foo", 4, false)
		assert.Equal(t, map[string]map[int]bool{"foo": {3: true, 4: false}}, mp)

		Map2Set(&mp, "bar", 5, true)
		assert.Equal(t, map[string]map[int]bool{"foo": {3: true, 4: false}, "bar": {5: true}}, mp)
	}
}

func TestMapSliceAdd_(t *testing.T) {
	{
		var mp map[string][]int
		MapSliceAdd(&mp, "foo", 3)
		assert.NotNil(t, mp)
		assert.NotNil(t, mp["foo"])
		assert.Equal(t, map[string][]int{"foo": {3}}, mp)

		MapSliceAdd(&mp, "foo", 4)
		assert.Equal(t, map[string][]int{"foo": {3, 4}}, mp)

		MapSliceAdd(&mp, "bar", 5)
		assert.Equal(t, map[string][]int{"foo": {3, 4}, "bar": {5}}, mp)
	}

	{
		var mp map[string]map[int][]bool
		Map2SliceAdd(&mp, "foo", 3, true)
		assert.NotNil(t, mp)
		assert.NotNil(t, mp["foo"])
		assert.NotNil(t, mp["foo"][3])
		assert.Equal(t, map[string]map[int][]bool{"foo": {3: []bool{true}}}, mp)

		Map2SliceAdd(&mp, "foo", 3, false)
		assert.Equal(t, map[string]map[int][]bool{"foo": {3: []bool{true, false}}}, mp)

		Map2SliceAdd(&mp, "bar", 4, true)
		assert.Equal(t, map[string]map[int][]bool{"foo": {3: []bool{true, false}}, "bar": {4: []bool{true}}}, mp)
	}
}

func TestMapSort_(t *testing.T) {
	// Ordered
	{
		mp := map[int]int{2: 2, 3: 3, 1: 1}
		slc := MapSortOrdered(mp)
		assert.Equal(t, []tuple.Two[int, int]{tuple.Of2(1, 1), tuple.Of2(2, 2), tuple.Of2(3, 3)}, slc)
	}

	// Complex
	{
		mp := map[complex64]int{2: 2, 3: 3, 1: 1}
		slc := MapSortComplex(mp)
		assert.Equal(t, []tuple.Two[complex64, int]{tuple.Of2(complex64(1+0i), 1), tuple.Of2(complex64(2+0i), 2), tuple.Of2(complex64(3+0i), 3)}, slc)
	}

	// Cmp
	{
		mp := map[*big.Int]int{big.NewInt(2): 2, big.NewInt(3): 3, big.NewInt(1): 1}
		slc := MapSortCmp(mp)
		assert.Equal(t, []tuple.Two[*big.Int, int]{tuple.Of2(big.NewInt(1), 1), tuple.Of2(big.NewInt(2), 2), tuple.Of2(big.NewInt(3), 3)}, slc)
	}

	// By
	{
		mp := map[int]int{2: 2, 3: 3, 1: 1}
		slc := MapSortBy(mp, func(i, j int) bool { return i < j })
		assert.Equal(t, []tuple.Two[int, int]{tuple.Of2(1, 1), tuple.Of2(2, 2), tuple.Of2(3, 3)}, slc)
	}
}

func TestMapKeysToSlice_(t *testing.T) {
	assert.Equal(t, []int{}, MapKeysToSlice((map[int]int)(nil)))
	assert.Equal(t, []int{}, MapKeysToSlice(map[int]int{}))
	assert.Equal(t, []int{1}, MapKeysToSlice(map[int]int{1: 0}))
	assert.Equal(t, []int{1, 2, 3}, SliceSortOrdered(MapKeysToSlice(map[int]int{1: 0, 2: 0, 3: 0})))
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

func TestAnd_(t *testing.T) {
	lt5_10 := And(lessThan5, lessThan10)
	assert.True(t, lt5_10(3))
	assert.False(t, lt5_10(5))
	assert.False(t, lt5_10(7))
	assert.False(t, lt5_10(10))
	assert.False(t, lt5_10(12))
}

func TestNot_(t *testing.T) {
	nlt5 := Not(lessThan5)
	assert.False(t, nlt5(3))
	assert.True(t, nlt5(5))
	assert.True(t, nlt5(7))
	assert.True(t, nlt5(10))
	assert.True(t, nlt5(12))
}

func TestOr_(t *testing.T) {
	lt5_gt10 := Or(lessThan5, greaterThan10)
	assert.True(t, lt5_gt10(3))
	assert.False(t, lt5_gt10(5))
	assert.False(t, lt5_gt10(7))
	assert.False(t, lt5_gt10(10))
	assert.True(t, lt5_gt10(12))
}

func TestLessThan_(t *testing.T) {
	lt5 := LessThan(5)
	assert.True(t, lt5(3))
	assert.False(t, lt5(5))
	assert.False(t, lt5(7))
	assert.False(t, lt5(10))
	assert.False(t, lt5(12))
}

func TestLessThanEqual_(t *testing.T) {
	lte5 := LessThanEqual(5)
	assert.True(t, lte5(3))
	assert.True(t, lte5(5))
	assert.False(t, lte5(7))
	assert.False(t, lte5(10))
	assert.False(t, lte5(12))
}

func TestEqual_(t *testing.T) {
	eq5 := Equal(5)
	assert.False(t, eq5(3))
	assert.True(t, eq5(5))
	assert.False(t, eq5(7))
	assert.False(t, eq5(10))
	assert.False(t, eq5(12))
}

func TestIn_(t *testing.T) {
	in := In("foo", "bar")
	assert.False(t, in("baz"))
	assert.True(t, in("foo"))
	assert.True(t, in("bar"))
}

func TestGreaterThan_(t *testing.T) {
	gt5 := GreaterThan(5)
	assert.False(t, gt5(3))
	assert.False(t, gt5(5))
	assert.True(t, gt5(7))
	assert.True(t, gt5(10))
	assert.True(t, gt5(12))
}

func TestGreaterThanEqual_(t *testing.T) {
	gte5 := GreaterThanEqual(5)
	assert.False(t, gte5(3))
	assert.True(t, gte5(5))
	assert.True(t, gte5(7))
	assert.True(t, gte5(10))
	assert.True(t, gte5(12))
}

func TestIsNegative_(t *testing.T) {
	neg := IsNegative[int]()
	assert.True(t, neg(-3))
	assert.False(t, neg(0))
	assert.False(t, neg(3))
}

func TestIsNonNegative_(t *testing.T) {
	nneg := IsNonNegative[int]()
	assert.False(t, nneg(-3))
	assert.True(t, nneg(0))
	assert.True(t, nneg(3))
}

func TestIsPositive_(t *testing.T) {
	pos := IsPositive[int]()
	assert.False(t, pos(-3))
	assert.False(t, pos(0))
	assert.True(t, pos(3))
}

func TestCompose_(t *testing.T) {
	var (
		fn1 = func(i int) int { return i + 2 }
		fn2 = func(i int) int { return i * 3 }
		fn3 = func(i int) int { return i - 4 }
	)

	fn := Compose(fn1)
	assert.Equal(t, 3, fn(1))

	fn = Compose(fn1, fn2)
	assert.Equal(t, 9, fn(1))

	fn = Compose(fn1, fn2, fn3)
	assert.Equal(t, 5, fn(1))

	fn = Compose(fn, fn)
	assert.Equal(t, 17, fn(1))
}

func stringToInt(t string) int {
	return MustValue(strconv.Atoi(t))
}

func TestCompose2_(t *testing.T) {
	fn := Compose2(
		strconv.Itoa,
		stringToInt,
	)

	assert.Equal(t, 1, fn(1))
}

func TestCompose3_(t *testing.T) {
	fn := Compose3(
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
	)

	assert.Equal(t, "1", fn(1))
}

func TestCompose4_(t *testing.T) {
	fn := Compose4(
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
	)

	assert.Equal(t, 1, fn(1))
}

func TestCompose5_(t *testing.T) {
	fn := Compose5(
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
	)

	assert.Equal(t, "1", fn(1))
}

func TestCompose6_(t *testing.T) {
	fn := Compose6(
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
	)

	assert.Equal(t, 1, fn(1))
}

func TestCompose7_(t *testing.T) {
	fn := Compose7(
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
	)

	assert.Equal(t, "1", fn(1))
}

func TestCompose8_(t *testing.T) {
	fn := Compose8(
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
	)

	assert.Equal(t, 1, fn(1))
}

func TestCompose9_(t *testing.T) {
	fn := Compose9(
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
	)

	assert.Equal(t, "1", fn(1))
}

func TestCompose10_(t *testing.T) {
	fn := Compose10(
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
		strconv.Itoa,
		stringToInt,
	)

	assert.Equal(t, 1, fn(1))
}

func TestTernary_(t *testing.T) {
	assert.Equal(t, 1, Ternary(1 < 2, 1, 2))
	assert.Equal(t, 1, TernaryResult(1 < 2, func() int { return 1 }, func() int { return 2 }))

	assert.Equal(t, 2, Ternary(1 > 2, 1, 2))
	assert.Equal(t, 2, TernaryResult(1 > 2, func() int { return 1 }, func() int { return 2 }))
}

func TestNillable_(t *testing.T) {
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

	TryTo(
		func() {
			IsNil[int]()
			assert.Fail(t, "Must die")
		},
		func(err any) { assert.Equal(t, fmt.Errorf("Type int is not a nillable type"), err) },
	)

  assert.False(t, IsNilValue(0))
  assert.False(t, IsNilValue(c))
  assert.True(t, IsNilValue(cn))

	TryTo(
		func() {
			IsNonNil[int]()
			assert.Fail(t, "Must die")
		},
		func(err any) { assert.Equal(t, fmt.Errorf("Type int is not a nillable type"), err) },
	)

	TryTo(
		func() {
			MustBeNillable(reflect.TypeOf(0))
			assert.Fail(t, "Must die")
		},
		func(err any) { assert.Equal(t, fmt.Errorf("Type int is not a nillable type"), err) },
	)

	assert.False(t, Nillable(reflect.TypeOf(0)))

	assert.True(t, IsNil[chan int]()(cn))
	assert.True(t, IsNonNil[chan int]()(c))
	MustBeNillable(reflect.TypeOf(cn))
	MustBeNillable(reflect.TypeOf(c))
	assert.True(t, Nillable(reflect.TypeOf(cn)))
	assert.True(t, Nillable(reflect.TypeOf(c)))

	assert.True(t, IsNil[func()]()(fn))
	assert.True(t, IsNonNil[func()]()(f))
	MustBeNillable(reflect.TypeOf(fn))
	MustBeNillable(reflect.TypeOf(f))
	assert.True(t, Nillable(reflect.TypeOf(fn)))
	assert.True(t, Nillable(reflect.TypeOf(f)))

	assert.True(t, IsNil[map[int]int]()(mn))
	assert.True(t, IsNonNil[map[int]int]()(m))
	MustBeNillable(reflect.TypeOf(mn))
	MustBeNillable(reflect.TypeOf(m))
	assert.True(t, Nillable(reflect.TypeOf(mn)))
	assert.True(t, Nillable(reflect.TypeOf(m)))

	assert.True(t, IsNil[*int]()(pn))
	assert.True(t, IsNonNil[*int]()(p))
	MustBeNillable(reflect.TypeOf(pn))
	MustBeNillable(reflect.TypeOf(p))
	assert.True(t, Nillable(reflect.TypeOf(pn)))
	assert.True(t, Nillable(reflect.TypeOf(p)))

	assert.True(t, IsNil[[]int]()(sn))
	assert.True(t, IsNonNil[[]int]()(s))
	MustBeNillable(reflect.TypeOf(sn))
	MustBeNillable(reflect.TypeOf(s))
	assert.True(t, Nillable(reflect.TypeOf(sn)))
	assert.True(t, Nillable(reflect.TypeOf(s)))

	assert.True(t, IsNonNil[[]int]()(a.([]int)))
	MustBeNillable(reflect.TypeOf(a))
	assert.True(t, Nillable(reflect.TypeOf(a)))
}

func TestMust_(t *testing.T) {
	var e error
	Must(e)

	e = fmt.Errorf("bob")
	TryTo(
		func() {
			Must(e)
			assert.Fail(t, "Must die")
		},
		func(err any) { assert.Equal(t, e, err) },
	)
}

func TestMustValue_(t *testing.T) {
	var (
		e error
		i int
	)
	assert.Equal(t, i, MustValue(i, e))

	e = fmt.Errorf("bob")
	TryTo(
		func() {
			MustValue(i, e)
			assert.Fail(t, "Must die")
		},
		func(err any) { assert.Equal(t, e, err) },
	)
}

func TestMustValue2_(t *testing.T) {
	var (
		e      error
		p1, p2 = 1, 2
		r1, r2 int
	)
	r1, r2 = MustValue2(p1, p2, e)
	assert.Equal(t, p1, r1)
	assert.Equal(t, p2, r2)

	e = fmt.Errorf("bob")
	TryTo(
		func() {
			MustValue2(p1, p2, e)
			assert.Fail(t, "Must die")
		},
		func(err any) { assert.Equal(t, e, err) },
	)
}

func TestMustValue3_(t *testing.T) {
	var (
		e          error
		p1, p2, p3 = 1, 2, 3
		r1, r2, r3 int
	)
	r1, r2, r3 = MustValue3(p1, p2, p3, e)
	assert.Equal(t, p1, r1)
	assert.Equal(t, p2, r2)
	assert.Equal(t, p3, r3)

	e = fmt.Errorf("bob")
	TryTo(
		func() {
			MustValue3(p1, p2, p3, e)
			assert.Fail(t, "Must die")
		},
		func(err any) { assert.Equal(t, e, err) },
	)
}

func TestAssertType_(t *testing.T) {
	assert.Equal(t, tuple.Of2[bool, error](true, nil), tuple.Of2(AssertType[bool]("", true)))
	assert.Equal(t, tuple.Of2[bool, error](false, fmt.Errorf("expected foo to be bool, not string")), tuple.Of2(AssertType[bool]("foo", "")))

	assert.Equal(t, tuple.Of2[int8, error](1, nil), tuple.Of2(AssertType[int8]("", int8(1))))
	assert.Equal(t, tuple.Of2[int8, error](0, fmt.Errorf("expected foo to be int8, not bool")), tuple.Of2(AssertType[int8]("foo", false)))

	assert.Equal(t, true, MustAssertType[bool]("", true))
	TryTo(
		func() {
			MustAssertType[bool]("foo", "")
			assert.Fail(t, "Must Die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf("expected foo to be bool, not string"), e)
		},
	)

	assert.Equal(t, int8(1), MustAssertType[int8]("", int8(1)))
	TryTo(
		func() {
			MustAssertType[int8]("foo", false)
			assert.Fail(t, "Must Die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf("expected foo to be int8, not bool"), e)
		},
	)
}

func TestConvertToSlice_(t *testing.T) {
	// One dimension
	assert.Equal(t, tuple.Of2[[]bool, error]([]bool{true}, nil), tuple.Of2(ConvertToSlice[bool]("", any([]any{true}))))
	assert.Equal(
		t,
		tuple.Of2[[]bool, error]([]bool(nil),
			fmt.Errorf("expected foo to be []interface {}, not string")), tuple.Of2(ConvertToSlice[bool]("foo", any(""))),
	)
	assert.Equal(
		t,
		tuple.Of2[[]bool, error]([]bool(nil),
			fmt.Errorf("expected foo[0] to be bool, not string")), tuple.Of2(ConvertToSlice[bool]("foo", any([]any{""}))),
	)

	assert.Equal(t, []int8{1}, MustConvertToSlice[int8]("", any([]any{int8(1)})))
	TryTo(
		func() {
			MustConvertToSlice[int8]("foo", any([]any{false}))
			assert.Fail(t, "Must Die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf("expected foo[0] to be int8, not bool"), e)
		},
	)

	// Two dimensions
	assert.Equal(t, tuple.Of2[[][]bool, error]([][]bool{{true}}, nil), tuple.Of2(ConvertToSlice2[bool]("", any([]any{[]any{true}}))))
	assert.Equal(
		t,
		tuple.Of2[[][]bool, error]([][]bool(nil),
			fmt.Errorf("expected foo to be []interface {}, not string")), tuple.Of2(ConvertToSlice2[bool]("foo", any(""))),
	)
	assert.Equal(
		t,
		tuple.Of2[[][]bool, error]([][]bool(nil),
			fmt.Errorf("expected foo[0] to be []interface {}, not string")), tuple.Of2(ConvertToSlice2[bool]("foo", any([]any{""}))),
	)
	assert.Equal(
		t,
		tuple.Of2[[][]bool, error]([][]bool(nil),
			fmt.Errorf("expected foo[0][0] to be bool, not string")), tuple.Of2(ConvertToSlice2[bool]("foo", any([]any{[]any{""}}))),
	)

	assert.Equal(t, [][]int8{{1}}, MustConvertToSlice2[int8]("", any([]any{[]any{int8(1)}})))
	TryTo(
		func() {
			MustConvertToSlice2[int8]("foo", any([]any{[]any{false}}))
			assert.Fail(t, "Must Die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf("expected foo[0][0] to be int8, not bool"), e)
		},
	)
}

func TestConvertToMap_(t *testing.T) {
	assert.Equal(
		t,
		tuple.Of2[map[string]bool, error](map[string]bool{"bar": true}, nil),
		tuple.Of2(ConvertToMap[string, bool]("", any(map[string]any{"bar": true}))),
	)
	assert.Equal(
		t,
		tuple.Of2[map[string]bool, error](map[string]bool(nil),
			fmt.Errorf("expected foo to be map[string]interface {}, not string")), tuple.Of2(ConvertToMap[string, bool]("foo", any(""))),
	)
	assert.Equal(
		t,
		tuple.Of2[map[string]bool, error](map[string]bool(nil),
			fmt.Errorf("expected foo[bar] to be bool, not string")), tuple.Of2(ConvertToMap[string, bool]("foo", any(map[string]any{"bar": ""}))),
	)

	assert.Equal(
		t, map[string]int8{"bar": 1},
		MustConvertToMap[string, int8]("", any(map[string]any{"bar": int8(1)})),
	)
	TryTo(
		func() {
			MustConvertToMap[string, int8]("foo", any(map[string]any{"bar": false}))
			assert.Fail(t, "Must Die")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf("expected foo[bar] to be int8, not bool"), e)
		},
	)
}

func TestSupplier_(t *testing.T) {
	supplier := SupplierOf(5)
	assert.Equal(t, 5, supplier())
	assert.Equal(t, 5, supplier())

	var called bool
	supplier = SupplierCached(func() int { called = true; return 7 })

	assert.False(t, called)
	assert.Equal(t, 7, supplier())
	assert.True(t, called)

	called = false
	assert.False(t, called)
	assert.Equal(t, 7, supplier())
	assert.False(t, called)
}

func TestFirstValue2_(t *testing.T) {
	assert.Equal(t, 1, FirstValue2(1, 2))
}

func TestFirstValue3_(t *testing.T) {
	assert.Equal(t, 1, FirstValue3(1, 2, 3))
}

func TestSecondValue2_(t *testing.T) {
	assert.Equal(t, 2, SecondValue2(1, 2))
}

func TestSecondValue3_(t *testing.T) {
	assert.Equal(t, 2, SecondValue3(1, 2, 3))
}

func TestThirdValue3_(t *testing.T) {
	assert.Equal(t, 3, ThirdValue3(1, 2, 3))
}

func TestIgnoreResult_(t *testing.T) {
	called := false
	IgnoreResult(func() int { called = true; return 0 })()
	assert.True(t, called)
}

func TestTryTo_(t *testing.T) {
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
