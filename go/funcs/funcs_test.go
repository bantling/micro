package funcs

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestEqualTo(t *testing.T) {
	eq5 := EqualTo(5)
	assert.False(t, eq5(3))
	assert.True(t, eq5(5))
	assert.False(t, eq5(7))
	assert.False(t, eq5(10))
	assert.False(t, eq5(12))
}

func TestLessThan(t *testing.T) {
	lt5 := LessThan(5)
	assert.True(t, lt5(3))
	assert.False(t, lt5(5))
	assert.False(t, lt5(7))
	assert.False(t, lt5(10))
	assert.False(t, lt5(12))
}

func TestLessEqual(t *testing.T) {
	lte5 := LessThanEqual(5)
	assert.True(t, lte5(3))
	assert.True(t, lte5(5))
	assert.False(t, lte5(7))
	assert.False(t, lte5(10))
	assert.False(t, lte5(12))
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

func TestSupplier(t *testing.T) {
	supplier := Supplier(5)
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

func TestTryTo(t *testing.T) {
	var (
		tryCalled    bool
		panicValue   error
		closerCalled bool
		theError     = fmt.Errorf("The error")
	)

	TryTo(
		func() { tryCalled = true },
		func(err error) { panicValue = err },
		func() error { closerCalled = true; return nil },
	)
	assert.True(t, tryCalled)
	assert.Nil(t, panicValue)
	assert.True(t, closerCalled)

	tryCalled, panicValue, closerCalled = false, nil, false
	TryTo(
		func() { tryCalled = true; panic(theError) },
		func(err error) { panicValue = err },
	)
	assert.True(t, tryCalled)
	assert.Equal(t, theError, panicValue)
	assert.False(t, closerCalled)

	tryCalled, panicValue, closerCalled = false, nil, false
	TryTo(
		func() { tryCalled = true },
		func(err error) { panicValue = err },
		func() error { return theError },
		func() error { closerCalled = true; return nil },
	)
	assert.True(t, tryCalled)
	assert.Equal(t, theError, panicValue)
	assert.False(t, closerCalled)

	tryCalled, panicValue, closerCalled = false, nil, false
	TryTo(
		func() { tryCalled = true },
		func(err error) { panicValue = err },
		func() error { return nil },
		func() error { closerCalled = true; return nil },
	)
	assert.True(t, tryCalled)
	assert.Nil(t, panicValue)
	assert.True(t, closerCalled)
}
