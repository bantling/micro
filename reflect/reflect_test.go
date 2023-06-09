package reflect

// SPDX-License-Identifier: Apache-2.0

import (
	"math/big"
	goreflect "reflect"
	"testing"

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

func TestIsBigPtr_(t *testing.T) {
	assert.True(t, IsBigPtr(goreflect.TypeOf((*big.Int)(nil))))
	assert.True(t, IsBigPtr(goreflect.TypeOf((*big.Float)(nil))))
	assert.True(t, IsBigPtr(goreflect.TypeOf((*big.Rat)(nil))))
	assert.False(t, IsBigPtr(goreflect.TypeOf(0)))
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

func TestIsPrimitive_(t *testing.T) {
	assert.True(t, IsPrimitive(goreflect.TypeOf(0)))
	assert.False(t, IsPrimitive(goreflect.PtrTo(goreflect.TypeOf(0))))
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

func TestToBaseType_(t *testing.T) {
	{
		// int
		val := goreflect.ValueOf(0)
		ToBaseType(&val)
		assert.Equal(t, goreflect.TypeOf(0), val.Type())
		assert.Equal(t, 0, val.Interface())

		// *int
		var i int
		val = goreflect.ValueOf(&i)
		ToBaseType(&val)
		assert.Equal(t, goreflect.PtrTo(goreflect.TypeOf(0)), val.Type())
		assert.Equal(t, &i, val.Interface())
	}

	{
		// subString
		val := goreflect.ValueOf(subString("foo"))
		ToBaseType(&val)
		assert.Equal(t, goreflect.TypeOf("foo"), val.Type())
		assert.Equal(t, "foo", val.Interface())

		// *subString
		i := subString("foo")
		val = goreflect.ValueOf(&i)
		ToBaseType(&val)
		assert.Equal(t, goreflect.PtrTo(goreflect.TypeOf("foo")), val.Type())
		assert.Equal(t, "foo", *(val.Interface().(*string)))
	}

	{
		// rune
		val := goreflect.ValueOf('0')
		ToBaseType(&val)
		assert.Equal(t, goreflect.TypeOf(int32('0')), val.Type())
		assert.Equal(t, int32('0'), val.Interface())

		// *rune
		i := '0'
		val = goreflect.ValueOf(&i)
		ToBaseType(&val)
		assert.Equal(t, goreflect.PtrTo(goreflect.TypeOf(int32('0'))), val.Type())
		assert.Equal(t, int32('0'), *(val.Interface().(*int32)))
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
