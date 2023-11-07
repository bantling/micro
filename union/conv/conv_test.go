package conv

// SPDX-License-Identifier: Apache-2.0

import (
  "fmt"
  goreflect "reflect"
  "testing"

  "github.com/bantling/micro/union"
	"github.com/stretchr/testify/assert"
)

func TestMaybeWrapperInfo_(t *testing.T) {
  // Maybe[int]
  var (
    maybeVal = union.Of(1)
    instance = goreflect.ValueOf(&maybeVal)
    info MaybeWrapperInfo
    intTyp, strTyp = goreflect.TypeOf(0), goreflect.TypeOf("")
  )

  // Package path and prefix are correct
  assert.Equal(t, goreflect.ValueOf(maybeVal).Type().PkgPath(), info.PackagePath())
  assert.Equal(t, "Maybe", info.TypeNamePrefix())

  // Accepts only int type
  assert.True(t, info.AcceptsType(instance, intTyp))
  assert.False(t, info.AcceptsType(instance, strTyp))

  // Can always be empty
  assert.True(t, info.CanBeEmpty(instance))

  // Converts only to int type
  assert.True(t, info.ConvertibleTo(instance, intTyp))
  assert.False(t, info.ConvertibleTo(instance, strTyp))

  // Get value as int
  val, pres, err := info.Get(instance, intTyp)
  assert.Equal(t, 1, val.Interface())
  assert.True(t, pres)
  assert.Nil(t, err)

  // Get value as string fails
  val, pres, err = info.Get(instance, strTyp)
  assert.False(t, val.IsValid())
  assert.False(t, pres)
  assert.Equal(t, fmt.Errorf("A *union.Maybe[int] cannot be converted to string"), err)

  // Set value to present int
  assert.Nil(t, info.Set(instance, goreflect.ValueOf(2), true))
  val, pres, err = info.Get(instance, intTyp)
  assert.Equal(t, 2, val.Interface())
  assert.True(t, pres)
  assert.Nil(t, err)

  // Set value to present string
  assert.Equal(t, fmt.Errorf("A *union.Maybe[int] cannot be set to a(n) string"), info.Set(instance, goreflect.ValueOf("3"), true))

  // Set value to empty
  for _, v := range []goreflect.Value{goreflect.Value{}, goreflect.ValueOf(4), goreflect.ValueOf("4")} {
    assert.Nil(t, info.Set(instance, v, false))
    val, pres, err = info.Get(instance, intTyp)
    assert.False(t, val.IsValid())
    assert.False(t, pres)
    assert.Nil(t, err)
  }

  // Errors
}
