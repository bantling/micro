package reflect

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math/big"
	goreflect "reflect"
	"testing"

	"github.com/bantling/micro/funcs"
	"github.com/stretchr/testify/assert"
)

type subString string

func TestDerefType_(t *testing.T) {
	{
		typ := goreflect.TypeOf(0)
		assert.True(t, typ == DerefType(typ))
	}

	{
		i := 1
		typ := goreflect.TypeOf(&i)
		assert.Equal(t, goreflect.TypeOf(0), DerefType(typ))
	}

	{
		i := 2
		p := &i
		typ := goreflect.TypeOf(&p)
		assert.Equal(t, goreflect.TypeOf(0), DerefType(typ))
	}
}

func TestDerefTypeMaxOnePtr_(t *testing.T) {
	{
		typ := goreflect.TypeOf(0)
		assert.True(t, typ == DerefTypeMaxOnePtr(typ))
	}

	{
		i := 1
		typ := goreflect.TypeOf(&i)
		assert.True(t, typ == DerefTypeMaxOnePtr(typ))
	}

	{
		i := 2
		p := &i
		typ := goreflect.TypeOf(&p)
		assert.Equal(t, goreflect.TypeOf((*int)(nil)), DerefTypeMaxOnePtr(typ))
	}
}

func TestDerefValue_(t *testing.T) {
	{
		v := goreflect.ValueOf(0)
		assert.True(t, v == DerefValue(v))
	}

	{
		i := 1
		v := goreflect.ValueOf(&i)
		assert.Equal(t, 1, DerefValue(v).Interface())
	}

	assert.Equal(t, goreflect.Value{}, DerefValue(goreflect.ValueOf((*int)(nil))))

	{
		i := 2
		p := &i
		v := goreflect.ValueOf(&p)
		assert.Equal(t, 2, DerefValue(v).Interface())
	}

	assert.Equal(t, goreflect.Value{}, DerefValue(goreflect.ValueOf((**int)(nil))))

	{
		var p *int
		p2 := &p
		assert.Equal(t, goreflect.Value{}, DerefValue(goreflect.ValueOf(p2)))
	}
}

func TestDerefValueMaxOnePtr_(t *testing.T) {
	{
		v := goreflect.ValueOf(0)
		assert.True(t, v == DerefValueMaxOnePtr(v))
	}

	{
		var p *int
		v := goreflect.ValueOf(p)
		assert.Equal(t, v, DerefValueMaxOnePtr(v))
	}

	{
		var p *int
		p2 := &p
		v := DerefValueMaxOnePtr(goreflect.ValueOf(p2))
		assert.True(t, v.IsValid())
		assert.True(t, v.IsNil())
	}

	{
		i := 1
		p := &i
		v := goreflect.ValueOf(&p)
		assert.Equal(t, 1, DerefValueMaxOnePtr(v).Elem().Interface())
	}

	{
		var p **int
		assert.False(t, DerefValueMaxOnePtr(goreflect.ValueOf(p)).IsValid())
	}
}

func TestFieldsByName_(t *testing.T) {
	{
		str := ""

		typ := goreflect.TypeOf(str)
		assert.Equal(t, map[string]goreflect.StructField(nil), FieldsByName(typ))
	}

	{
		str := struct {
		}{}

		typ := goreflect.TypeOf(str)
		assert.Equal(t, map[string]goreflect.StructField(nil), FieldsByName(typ))
	}

	{
		str := struct {
			Foo string
			Bar int
		}{}

		typ := goreflect.TypeOf(str)
		fooFld, _ := typ.FieldByName("Foo")
		barFld, _ := typ.FieldByName("Bar")

		assert.Equal(t, map[string]goreflect.StructField{"Foo": fooFld, "Bar": barFld}, FieldsByName(typ))
	}
}

func TestSetPointerValue_(t *testing.T) {
	m := 0
	SetPointerValue(goreflect.ValueOf(&m), goreflect.ValueOf(1))
	assert.Equal(t, 1, m)

	p := &m
	SetPointerValue(goreflect.ValueOf(&p), goreflect.ValueOf(2))
	assert.Equal(t, 2, m)
}

func TestIsBigPtr_(t *testing.T) {
	assert.True(t, IsBigPtr(goreflect.TypeOf((*big.Int)(nil))))
	assert.True(t, IsBigPtr(goreflect.TypeOf((*big.Float)(nil))))
	assert.True(t, IsBigPtr(goreflect.TypeOf((*big.Rat)(nil))))
	assert.False(t, IsBigPtr(goreflect.TypeOf(0)))
}

func TestIsSignedInteger_(t *testing.T) {
	assert.True(t, IsSignedInteger(goreflect.TypeOf(int(0))))
	assert.True(t, IsSignedInteger(goreflect.TypeOf(int8(0))))
	assert.True(t, IsSignedInteger(goreflect.TypeOf(int16(0))))
	assert.True(t, IsSignedInteger(goreflect.TypeOf(int32(0))))
	assert.True(t, IsSignedInteger(goreflect.TypeOf(int64(0))))

	assert.False(t, IsSignedInteger(goreflect.TypeOf(uint(0))))
	assert.False(t, IsSignedInteger(goreflect.TypeOf("")))
}

func TestIsUnsignedInteger_(t *testing.T) {
	assert.True(t, IsUnsignedInteger(goreflect.TypeOf(uint(0))))
	assert.True(t, IsUnsignedInteger(goreflect.TypeOf(uint8(0))))
	assert.True(t, IsUnsignedInteger(goreflect.TypeOf(uint16(0))))
	assert.True(t, IsUnsignedInteger(goreflect.TypeOf(uint32(0))))
	assert.True(t, IsUnsignedInteger(goreflect.TypeOf(uint64(0))))

	assert.False(t, IsUnsignedInteger(goreflect.TypeOf(int(0))))
	assert.False(t, IsUnsignedInteger(goreflect.TypeOf("")))
}

func TestIsNumeric_(t *testing.T) {
	assert.True(t, IsNumeric(goreflect.TypeOf(int(0))))
	assert.True(t, IsNumeric(goreflect.TypeOf(int8(0))))
	assert.True(t, IsNumeric(goreflect.TypeOf(int16(0))))
	assert.True(t, IsNumeric(goreflect.TypeOf(int32(0))))
	assert.True(t, IsNumeric(goreflect.TypeOf(int64(0))))

	assert.True(t, IsNumeric(goreflect.TypeOf(uint(0))))
	assert.True(t, IsNumeric(goreflect.TypeOf(uint8(0))))
	assert.True(t, IsNumeric(goreflect.TypeOf(uint16(0))))
	assert.True(t, IsNumeric(goreflect.TypeOf(uint32(0))))
	assert.True(t, IsNumeric(goreflect.TypeOf(uint64(0))))

	assert.True(t, IsNumeric(goreflect.TypeOf(float32(0))))
	assert.True(t, IsNumeric(goreflect.TypeOf(float64(0))))

	assert.True(t, IsNumeric(goreflect.TypeOf((*big.Int)(nil))))
	assert.True(t, IsNumeric(goreflect.TypeOf((*big.Float)(nil))))
	assert.True(t, IsNumeric(goreflect.TypeOf((*big.Rat)(nil))))

	assert.True(t, IsNumeric(goreflect.TypeOf(byte(0))))
	assert.True(t, IsNumeric(goreflect.TypeOf(rune(0))))

	assert.False(t, IsNumeric(goreflect.TypeOf("")))
}

func TestIsNillable_(t *testing.T) {
	assert.False(t, IsNillable(goreflect.TypeOf(0)))
	assert.True(t, IsNillable(goreflect.TypeOf((chan int)(nil))))
	assert.True(t, IsNillable(goreflect.TypeOf((func())(nil))))
	assert.True(t, IsNillable(goreflect.TypeOf((goreflect.Type)(nil))))
	assert.True(t, IsNillable(goreflect.TypeOf((map[int]any)(nil))))
	assert.True(t, IsNillable(goreflect.TypeOf((*int)(nil))))
	assert.True(t, IsNillable(goreflect.TypeOf(([]int)(nil))))

	assert.False(t, IsNillable(goreflect.ValueOf(0)))
	assert.True(t, IsNillable(goreflect.ValueOf((chan int)(nil))))
	assert.True(t, IsNillable(goreflect.ValueOf((func())(nil))))
	assert.True(t, IsNillable(goreflect.ValueOf((goreflect.Type)(nil))))
	assert.True(t, IsNillable(goreflect.ValueOf((map[int]any)(nil))))
	assert.True(t, IsNillable(goreflect.ValueOf((*int)(nil))))
	assert.True(t, IsNillable(goreflect.ValueOf(([]int)(nil))))
}

func TestIsNil_(t *testing.T) {
	assert.False(t, IsNil(goreflect.ValueOf(0)))
	assert.True(t, IsNil(goreflect.ValueOf((chan int)(nil))))
	assert.True(t, IsNil(goreflect.ValueOf((func())(nil))))
	assert.True(t, IsNil(goreflect.ValueOf((goreflect.Type)(nil))))
	assert.True(t, IsNil(goreflect.ValueOf((map[int]any)(nil))))
	assert.True(t, IsNil(goreflect.ValueOf((*int)(nil))))
	assert.True(t, IsNil(goreflect.ValueOf(([]int)(nil))))
}

func TestIsPrimitive_(t *testing.T) {
	assert.True(t, IsPrimitive(goreflect.TypeOf(0)))
	assert.False(t, IsPrimitive(goreflect.PtrTo(goreflect.TypeOf(0))))
}

func TestNumPointers_(t *testing.T) {
	assert.Equal(t, 0, NumPointers(goreflect.TypeOf(0)))
	assert.Equal(t, 1, NumPointers(goreflect.TypeOf((*int)(nil))))
	assert.Equal(t, 2, NumPointers(goreflect.TypeOf((**int)(nil))))
	assert.Equal(t, 3, NumPointers(goreflect.TypeOf((***int)(nil))))
}

func TestResolveValueType_(t *testing.T) {
	// Test special case
	slc := []any{"foo", 1}
	rslc := goreflect.ValueOf(slc)
	assert.Equal(t, goreflect.String, ResolveValueType(rslc.Index(0)).Kind())
	assert.Equal(t, goreflect.Int, ResolveValueType(rslc.Index(1)).Kind())

	// Test normal case
	assert.Equal(t, goreflect.String, ResolveValueType(goreflect.ValueOf("foo")).Kind())
	assert.Equal(t, goreflect.Int, ResolveValueType(goreflect.ValueOf(1)).Kind())
}

func TestTypeAssert_(t *testing.T) {
	{
		// any containing an int is not a string
		var i = []any{0}
		assert.Equal(t, fmt.Errorf("interface {} is int, not string"), TypeAssert(goreflect.ValueOf(i).Index(0), goreflect.TypeOf("")))
		assert.Equal(t, fmt.Errorf("foo: interface {} is int, not string"), TypeAssert(goreflect.ValueOf(i).Index(0), goreflect.TypeOf(""), "foo"))

		var failed bool
		funcs.TryTo(
			func() {
				MustTypeAssert(goreflect.ValueOf(i).Index(0), goreflect.TypeOf(""))
				assert.Fail(t, "Must die")
			},
			func(e any) {
				assert.Equal(t, fmt.Errorf("interface {} is int, not string"), e)
				failed = true
			},
		)
		assert.True(t, failed)
	}

	{
		// any containing an int is an int
		var i = []any{0}
		assert.Nil(t, TypeAssert(goreflect.ValueOf(i).Index(0), goreflect.TypeOf(0)))

		var failed bool
		funcs.TryTo(
			func() {
				MustTypeAssert(goreflect.ValueOf(i).Index(0), goreflect.TypeOf(""), "bar")
				assert.Fail(t, "Must die")
			},
			func(e any) {
				assert.Equal(t, fmt.Errorf("bar: interface {} is int, not string"), e)
				failed = true
			},
		)
		assert.True(t, failed)
	}

	{
		// int is not a string
		var i = 0
		assert.Equal(t, fmt.Errorf("int is int, not string"), TypeAssert(goreflect.ValueOf(i), goreflect.TypeOf("")))
		assert.Equal(t, fmt.Errorf("bar: int is int, not string"), TypeAssert(goreflect.ValueOf(i), goreflect.TypeOf(""), "bar"))
	}

	{
		// int is an int
		var i = 0
		assert.Nil(t, TypeAssert(goreflect.ValueOf(i), goreflect.TypeOf(0)))
	}
}

func TestTypeOf_(t *testing.T) {
	assert.Equal(t, "<invalid Value>", TypeOf(goreflect.Value{}))
	assert.Equal(t, "string", TypeOf(goreflect.ValueOf("")))
}

func TestTypeToBaseType_(t *testing.T) {
	{
		// int
		typ := goreflect.TypeOf(0)
		assert.Nil(t, TypeToBaseType(typ))
	}

	{
		// subString
		typ := goreflect.TypeOf(subString("foo"))
		assert.Equal(t, goreflect.TypeOf("foo"), TypeToBaseType(typ))
	}

	{
		// rune is an alias for int32, which cannot be distinguished from int32 by reflection
		typ := goreflect.TypeOf('0')
		assert.Nil(t, TypeToBaseType(typ))
	}
}

func TestValueToBaseType_(t *testing.T) {
	{
		// int
		val := goreflect.ValueOf(0)
		rval := ValueToBaseType(val)
		assert.Equal(t, goreflect.TypeOf(0), rval.Type())
		assert.Equal(t, 0, rval.Interface())

		// *int
		var i int
		val = goreflect.ValueOf(&i)
		rval = ValueToBaseType(val)
		assert.Equal(t, goreflect.PtrTo(goreflect.TypeOf(0)), rval.Type())
		assert.Equal(t, &i, rval.Interface())
	}

	{
		// subString
		val := goreflect.ValueOf(subString("foo"))
		rval := ValueToBaseType(val)
		assert.Equal(t, goreflect.TypeOf("foo"), rval.Type())
		assert.Equal(t, "foo", rval.Interface())

		// *subString
		i := subString("foo")
		val = goreflect.ValueOf(&i)
		rval = ValueToBaseType(val)
		assert.Equal(t, goreflect.PtrTo(goreflect.TypeOf("foo")), rval.Type())
		assert.Equal(t, "foo", *(rval.Interface().(*string)))
	}

	{
		// rune
		val := goreflect.ValueOf('0')
		rval := ValueToBaseType(val)
		assert.Equal(t, goreflect.TypeOf(int32('0')), rval.Type())
		assert.Equal(t, int32('0'), rval.Interface())

		// *rune
		i := '0'
		val = goreflect.ValueOf(&i)
		rval = ValueToBaseType(val)
		assert.Equal(t, goreflect.PtrTo(goreflect.TypeOf(int32('0'))), rval.Type())
		assert.Equal(t, int32('0'), *(rval.Interface().(*int32)))
	}
}

func TestValueMaxOnePtrType_(t *testing.T) {
	{
		var v goreflect.Value
		assert.Equal(t, goreflect.Type(nil), ValueMaxOnePtrType(v))
	}

	{
		v := goreflect.ValueOf(0)
		assert.Equal(t, goreflect.TypeOf(0), ValueMaxOnePtrType(v))
	}

	{
		var p *int
		v := goreflect.ValueOf(p)
		assert.Equal(t, goreflect.TypeOf(0), ValueMaxOnePtrType(v))
	}

	{
		var p *int
		p2 := &p
		v := goreflect.ValueOf(p2)
		assert.Equal(t, goreflect.Type(nil), ValueMaxOnePtrType(v))
	}

	{
		var i = 1
		p := &i
		v := goreflect.ValueOf(&p)
		assert.Equal(t, goreflect.Type(nil), ValueMaxOnePtrType(v))
	}

	{
		var p **int
		assert.Equal(t, goreflect.Type(nil), ValueMaxOnePtrType(goreflect.ValueOf(p)))
	}
}

func TestGetFieldByName_(t *testing.T) {
	type Foo struct {
		Bar int
	}

	typ := goreflect.TypeOf(Foo{})
	sf, _ := typ.FieldByName("Bar")
	assert.Equal(t, sf, GetFieldByName(typ, "Bar"))

	var failed bool
	funcs.TryTo(
		func() {
			GetFieldByName(typ, "Foo")
		},
		func(e any) {
			assert.Equal(t, fmt.Errorf("type %s does not have a field named Foo", typ.String()), e)
			failed = true
		},
	)
	assert.True(t, failed)
}
