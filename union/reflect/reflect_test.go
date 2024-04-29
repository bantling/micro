package reflect

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math/big"
	goreflect "reflect"
	"testing"

	"github.com/bantling/micro/union"
	"github.com/stretchr/testify/assert"
)

func TestGetMaybeType_(t *testing.T) {
	assert.Equal(t, goreflect.TypeOf(0), GetMaybeType(goreflect.TypeOf(union.Maybe[int]{})))
	assert.Equal(t, goreflect.TypeOf(0), GetMaybeType(goreflect.TypeOf((*union.Maybe[int])(nil))))

	assert.Equal(t, goreflect.TypeOf((*big.Int)(nil)), GetMaybeType(goreflect.TypeOf(union.Maybe[*big.Int]{})))
	assert.Equal(t, goreflect.TypeOf((*big.Int)(nil)), GetMaybeType(goreflect.TypeOf((*union.Maybe[*big.Int])(nil))))

	assert.Nil(t, GetMaybeType(goreflect.TypeOf(0)))
}

func TestSafeMaybeAccess_(t *testing.T) {
	var (
		val  = union.Empty[int]()
		ptr1 = &val
		ptr2 = &ptr1
	)

	assert.True(t, safeMaybeAccess(goreflect.ValueOf(val)))
	assert.True(t, safeMaybeAccess(goreflect.ValueOf(ptr1)))
	assert.True(t, safeMaybeAccess(goreflect.ValueOf(ptr2)))

	assert.False(t, safeMaybeAccess(goreflect.Value{}))
	assert.False(t, safeMaybeAccess(goreflect.ValueOf(0)))
	assert.False(t, safeMaybeAccess(goreflect.ValueOf((*union.Maybe[int])(nil))))
	assert.False(t, safeMaybeAccess(goreflect.ValueOf((**union.Maybe[int])(nil))))
}

func TestMaybeValueIsPresent_(t *testing.T) {
	var (
		eval  = union.Empty[int]()
		eptr1 = &eval
		eptr2 = &eptr1

		pval  = union.Of(0)
		pptr1 = &pval
		pptr2 = &pptr1
	)

	assert.False(t, MaybeValueIsPresent(goreflect.Value{}))
	assert.False(t, MaybeValueIsPresent(goreflect.ValueOf(0)))
	assert.False(t, MaybeValueIsPresent(goreflect.ValueOf((*union.Maybe[int])(nil))))
	assert.False(t, MaybeValueIsPresent(goreflect.ValueOf((**union.Maybe[int])(nil))))
	assert.False(t, MaybeValueIsPresent(goreflect.ValueOf(eval)))
	assert.False(t, MaybeValueIsPresent(goreflect.ValueOf(eptr1)))
	assert.False(t, MaybeValueIsPresent(goreflect.ValueOf(eptr2)))

	assert.True(t, MaybeValueIsPresent(goreflect.ValueOf(pval)))
	assert.True(t, MaybeValueIsPresent(goreflect.ValueOf(pptr1)))
	assert.True(t, MaybeValueIsPresent(goreflect.ValueOf(pptr2)))
}

func TestGetMaybeValue_(t *testing.T) {
	assert.Equal(t, 1, GetMaybeValue(goreflect.ValueOf(union.Of(1))).Interface())
	assert.False(t, GetMaybeValue(goreflect.ValueOf(union.Empty[int]())).IsValid())
}

func TestSetMaybeValue_(t *testing.T) {
	// Set int value
	{
		m := union.Maybe[int]{}
		assert.False(t, MaybeValueIsPresent(goreflect.ValueOf(m)))
		assert.Nil(t, SetMaybeValue(goreflect.ValueOf(&m), goreflect.ValueOf(1)))
		assert.True(t, m.Present())
		assert.Equal(t, 1, m.Get())
	}

	{
		type Foo struct {
			Bar union.Maybe[int]
		}
		f := Foo{}

		assert.Nil(t, SetMaybeValue(goreflect.ValueOf(&f).Elem().FieldByName("Bar").Addr(), goreflect.ValueOf(2)))
		assert.True(t, f.Bar.Present())
		assert.Equal(t, 2, f.Bar.Get())
	}

	{
		type Foo struct {
			Bar union.Maybe[*int]
		}
		f := Foo{}
		i := 3
		assert.Nil(t, SetMaybeValue(goreflect.ValueOf(&f).Elem().FieldByName("Bar").Addr(), goreflect.ValueOf(&i)))
		assert.True(t, f.Bar.Present())
		assert.Equal(t, union.Of(&i), f.Bar)
	}

	{
		type Bar struct {
			Int int
			Str string
		}

		type Foo struct {
			Fld union.Maybe[Bar]
		}

		var (
			f       = Foo{Fld: union.Of(Bar{Int: 1})}
			barAddr = &f.Fld
		)

		// Copy Fld to get preinitialized object with Int = 1 and Str = ""
		barAddrVal := goreflect.ValueOf(&f).Elem().FieldByName("Fld").Addr()
		assert.True(t, MaybeValueIsPresent(barAddrVal))

		assert.Nil(t, SetMaybeValue(barAddrVal, goreflect.ValueOf(Bar{Int: 1})))
		assert.True(t, MaybeValueIsPresent(barAddrVal))
		assert.Equal(t, Bar{Int: 1}, GetMaybeValue(barAddrVal).Interface())
		assert.Equal(t, Bar{Int: 1}, f.Fld.Get())
		assert.True(t, barAddr == &f.Fld)
	}

	{
		type Bar struct {
			Int int
		}

		type Foo struct {
			Fld union.Maybe[*Bar]
		}

		var (
			f       = Foo{}
			barAddr = &f.Fld
		)

		barAddrVal := goreflect.ValueOf(&f).Elem().FieldByName("Fld").Addr()
		assert.False(t, MaybeValueIsPresent(barAddrVal))
		assert.Nil(t, SetMaybeValue(barAddrVal, goreflect.ValueOf(&Bar{Int: 1})))
		assert.True(t, MaybeValueIsPresent(barAddrVal))
		assert.Equal(t, Bar{Int: 1}, GetMaybeValue(barAddrVal).Elem().Interface())
		assert.Equal(t, Bar{Int: 1}, *f.Fld.Get())
		assert.True(t, barAddr == &f.Fld)
	}

	// Failures
	assert.Equal(
		t,
		fmt.Errorf("Cannot set the Maybe value of type <invalid Value>"),
		SetMaybeValue(goreflect.Value{}, goreflect.Value{}),
	)

	assert.Equal(
		t,
		fmt.Errorf("Cannot set the Maybe value of type union.Maybe[int] as it is not a pointer and not addressable"),
		SetMaybeValue(goreflect.ValueOf(union.Of(1)), goreflect.Value{}),
	)
}

func TestSetMaybeValueEmpty_(t *testing.T) {
	m := union.Of(1)
	SetMaybeValueEmpty(goreflect.ValueOf(&m))
	assert.False(t, m.Present())

	// Failures
	assert.Equal(
		t,
		fmt.Errorf("Cannot set the Maybe value of type <invalid Value>"),
		SetMaybeValueEmpty(goreflect.Value{}),
	)

	assert.Equal(
		t,
		fmt.Errorf("Cannot set the Maybe value of type union.Maybe[int] as it is not a pointer and not addressable"),
		SetMaybeValueEmpty(goreflect.ValueOf(union.Of(1))),
	)
}
