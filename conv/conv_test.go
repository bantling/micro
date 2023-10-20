package conv

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math"
	"math/big"
	goreflect "reflect"
	"testing"
	"unsafe"

	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/union"
	"github.com/stretchr/testify/assert"
)

// ==== LookupConversion errors

func TestLookupConversionErrNumPtrs_(t *testing.T) {
  // Too many src ptrs
  fn, err := LookupConversion(goreflect.TypeOf((**int)(nil)), goreflect.TypeOf(0))
  assert.Nil(t, fn)
  assert.Equal(t, fmt.Errorf("**int cannot be converted to int"), err)

  // Too many tgt ptrs
  fn, err = LookupConversion(goreflect.TypeOf(0), goreflect.TypeOf((**int)(nil)))
  assert.Nil(t, fn)
  assert.Equal(t, fmt.Errorf("int cannot be converted to **int"), err)

  // Too many src and tgt ptrs
  fn, err = LookupConversion(goreflect.TypeOf((**int)(nil)), goreflect.TypeOf((**int)(nil)))
  assert.Nil(t, fn)
  assert.Equal(t, fmt.Errorf("**int cannot be converted to **int"), err)
}

func TestLookupConversionErrBadTypes_(t *testing.T) {
  badTypes := []goreflect.Type{
    goreflect.TypeOf((uintptr)(0)),
    goreflect.TypeOf((chan int)(nil)),
    goreflect.TypeOf((func())(nil)),
    goreflect.TypeOf(unsafe.Pointer((*int)(nil))),
  }

  // error: src cannot be nil
  {
    fn, err := LookupConversion(nil, goreflect.TypeOf(0))
    assert.Nil(t, fn)
    assert.Equal(t, fmt.Errorf("<nil> cannot be converted to int"), err)
  }

  // error: tgt cannot be nil
  {
    fn, err := LookupConversion(goreflect.TypeOf(0), nil)
    assert.Nil(t, fn)
    assert.Equal(t, fmt.Errorf("int cannot be converted to <nil>"), err)
  }

  // error: src and tgt cannot be nil
  {
    fn, err := LookupConversion(nil, nil)
    assert.Nil(t, fn)
    assert.Equal(t, fmt.Errorf("<nil> cannot be converted to <nil>"), err)
  }

  // error: src cannot be uintptr, chan, func, or unsafe pointer
  {
    for _, styp := range badTypes {
      fn, err := LookupConversion(styp, goreflect.TypeOf(0))
      assert.Nil(t, fn)
      assert.Equal(t, fmt.Errorf("%s cannot be converted to int", styp), err)
    }
  }

  // error: tgt cannot be uintptr, chan, func, or unsafe pointer
  {
    for _, ttyp := range badTypes {
      fn, err := LookupConversion(goreflect.TypeOf(0), ttyp)
      assert.Nil(t, fn)
      assert.Equal(t, fmt.Errorf("int cannot be converted to %s", ttyp), err)
    }
  }

  // error: src and tgt cannot be uintptr, chan, func, or unsafe pointer
  {
    for _, typ := range badTypes {
      fn, err := LookupConversion(typ, typ)
      assert.Nil(t, fn)
      assert.Equal(t, fmt.Errorf("%s cannot be converted to %s", typ, typ), err)
    }
  }
}

// ==== LookupConversion exists

func TestLookupConversionExists_(t *testing.T) {
  var tgt string
  fn, err := LookupConversion(goreflect.TypeOf(0), goreflect.TypeOf(""))
  assert.NotNil(t, fn)
  assert.Nil(t, err)
  fn(1, &tgt)
  assert.Equal(t, "1", tgt)
}

// ==== LookupConversion copy

func TestLookupConversionCopy_(t *testing.T) {
  var tgt string
  fn, err := LookupConversion(goreflect.TypeOf(""), goreflect.TypeOf(""))
  assert.NotNil(t, fn)
  assert.Nil(t, err)
  fn("1", &tgt)
  assert.Equal(t, "1", tgt)
}

// ==== LookupConversion from Val

func TestLookupConversionVal2Val_(t *testing.T) {
  var src = 1
  var tgt string

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionVal2Base_(t *testing.T) {
  var src = 1
  type substring string
  var tgt substring

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionVal2Ptr_(t *testing.T) {
  var src = 1
  var tgt string
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionVal2PtrBase_(t *testing.T) {
  var src = 1
  type substring string
  var tgt substring
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, tgtp)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionVal2Maybe_(t *testing.T) {
  var src = 1
  var tgt union.Maybe[string]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", tgt.Get())
}

func TestLookupConversionVal2MaybeBase_(t *testing.T) {
  var src = 1
  type substring string
  var tgt union.Maybe[substring]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), tgt.Get())
}

func TestLookupConversionVal2MaybePtr_(t *testing.T) {
  var src = 1
  var tgt union.Maybe[*string]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", *tgt.Get())
}

func TestLookupConversionVal2MaybePtrBase_(t *testing.T) {
  var src = 1
  type substring string
  var tgt union.Maybe[*substring]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), *tgt.Get())
}

// ==== LookupConversion from Base

func TestLookupConversionBase2Val_(t *testing.T) {
  type subint int
  var src = subint(1)
  var tgt string

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionBase2Base_(t *testing.T) {
  type subint int
  var src = subint(1)
  type substring string
  var tgt substring

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionBase2Ptr_(t *testing.T) {
  type subint int
  var src = subint(1)
  var tgt string
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, tgtp)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionBase2PtrBase_(t *testing.T) {
  type subint int
  var src = subint(1)
  type substring string
  var tgt substring
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, tgtp)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionBase2Maybe_(t *testing.T) {
  type subint int
  var src = subint(1)
  var tgt union.Maybe[string]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", tgt.Get())
}

func TestLookupConversionBase2MaybeBase_(t *testing.T) {
  type subint int
  var src = subint(1)
  type substring string
  var tgt union.Maybe[substring]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), tgt.Get())
}

func TestLookupConversionBase2MaybePtr_(t *testing.T) {
  type subint int
  var src = subint(1)
  var tgt union.Maybe[*string]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", *tgt.Get())
}

func TestLookupConversionBase2MaybePtrBase_(t *testing.T) {
  type subint int
  var src = subint(1)
  type substring string
  var tgt union.Maybe[*substring]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), *tgt.Get())
}

// ==== LookupConversion from Ptr

func TestLookupConversionPtr2Val_(t *testing.T) {
  var src int = 1
  var srcp = &src
  var tgt string

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, &tgt)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionPtr2Base_(t *testing.T) {
  var src int = 1
  var srcp = &src
  type substring string
  var tgt substring

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, &tgt)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionPtr2Ptr_(t *testing.T) {
  var src int = 1
  var srcp = &src
  var tgt string
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, tgtp)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionPtr2PtrBase_(t *testing.T) {
  var src int = 1
  var srcp = &src
  type substring string
  var tgt substring
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, tgtp)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionPtr2Maybe_(t *testing.T) {
  var src int = 1
  var srcp = &src
  var tgt union.Maybe[string]

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", tgt.Get())
}

func TestLookupConversionPtr2MaybeBase_(t *testing.T) {
  var src int = 1
  var srcp = &src
  type substring string
  var tgt union.Maybe[substring]

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), tgt.Get())
}

func TestLookupConversionPtr2MaybePtr_(t *testing.T) {
  var src int = 1
  var srcp = &src
  var tgt union.Maybe[*string]

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", *tgt.Get())
}

func TestLookupConversionPtr2MaybePtrBase_(t *testing.T) {
  var src int = 1
  var srcp = &src
  type substring string
  var tgt union.Maybe[*substring]

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), *tgt.Get())
}

// ==== LookupConversion from Ptr Base

func TestLookupConversionPtrBase2Val_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcp = &src
  var tgt string

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, &tgt)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionPtrBase2Base_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcp = &src
  type substring string
  var tgt substring

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, &tgt)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionPtrBase2Ptr_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcp = &src
  var tgt string
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, tgtp)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionPtrBase2PtrBase_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcp = &src
  type substring string
  var tgt substring
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, tgtp)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionPtrBase2Maybe_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcp = &src
  var tgt union.Maybe[string]

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", tgt.Get())
}

func TestLookupConversionPtrBase2MaybeBase_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcp = &src
  type substring string
  var tgt union.Maybe[substring]

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), tgt.Get())
}

func TestLookupConversionPtrBase2MaybePtr_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcp = &src
  var tgt union.Maybe[*string]

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", *tgt.Get())
}

func TestLookupConversionPtrBase2MaybePtrBase_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcp = &src
  type substring string
  var tgt union.Maybe[*substring]

  fn, err := LookupConversion(goreflect.TypeOf(srcp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), *tgt.Get())
}

// ==== LookupConversion from Maybe

func TestLookupConversionMaybe2Val_(t *testing.T) {
  var src = union.Of(1)
  var tgt string

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionMaybe2Base_(t *testing.T) {
  var src = union.Of(1)
  type substring string
  var tgt substring

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionMaybe2Ptr_(t *testing.T) {
  var src = union.Of(1)
  var tgt string
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, tgtp)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionMaybe2PtrBase_(t *testing.T) {
  var src = union.Of(1)
  type substring string
  var tgt substring
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, tgtp)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionMaybe2Maybe_(t *testing.T) {
  var src = union.Of(1)
  var tgt union.Maybe[string]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", tgt.Get())
}

func TestLookupConversionMaybe2MaybeBase_(t *testing.T) {
  var src = union.Of(1)
  type substring string
  var tgt union.Maybe[substring]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), tgt.Get())
}

func TestLookupConversionMaybe2MaybePtr_(t *testing.T) {
  var src = union.Of(1)
  var tgt union.Maybe[*string]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", *tgt.Get())
}

func TestLookupConversionMaybe2MaybePtrBase_(t *testing.T) {
  var src = union.Of(1)
  type substring string
  var tgt union.Maybe[*substring]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), *tgt.Get())
}

// ==== LookupConversion from Maybe Base

func TestLookupConversionMaybeBase2Val_(t *testing.T) {
  type subint int
  var src = union.Of(subint(1))
  var tgt string

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionMaybeBase2Base_(t *testing.T) {
  type subint int
  var src = union.Of(subint(1))
  type substring string
  var tgt substring

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionMaybeBase2Ptr_(t *testing.T) {
  type subint int
  var src = union.Of(subint(1))
  var tgt string
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, tgtp)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionMaybeBase2PtrBase_(t *testing.T) {
  type subint int
  var src = union.Of(subint(1))
  type substring string
  var tgt substring
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, tgtp)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionMaybeBase2Maybe_(t *testing.T) {
  type subint int
  var src = union.Of(subint(1))
  var tgt union.Maybe[string]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", tgt.Get())
}

func TestLookupConversionMaybeBase2MaybeBase_(t *testing.T) {
  type subint int
  var src = union.Of(subint(1))
  type substring string
  var tgt union.Maybe[substring]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), tgt.Get())
}

func TestLookupConversionMaybeBase2MaybePtr_(t *testing.T) {
  type subint int
  var src = union.Of(subint(1))
  var tgt union.Maybe[*string]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", *tgt.Get())
}

func TestLookupConversionMaybeBase2MaybePtrBase_(t *testing.T) {
  type subint int
  var src = union.Of(subint(1))
  type substring string
  var tgt union.Maybe[*substring]

  fn, err := LookupConversion(goreflect.TypeOf(src), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(src, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), *tgt.Get())
}

// ==== LookupConversion from Maybe Ptr

func TestLookupConversionMaybePtr2Val_(t *testing.T) {
  var src = 1
  var srcmp = union.Of(&src)
  var tgt string

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, &tgt)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionMaybePtr2Base_(t *testing.T) {
  var src = 1
  var srcmp = union.Of(&src)
  type substring string
  var tgt substring

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, &tgt)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionMaybePtr2Ptr_(t *testing.T) {
  var src = 1
  var srcmp = union.Of(&src)
  var tgt string
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, tgtp)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionMaybePtr2PtrBase_(t *testing.T) {
  var src = 1
  var srcmp = union.Of(&src)
  type substring string
  var tgt substring
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, tgtp)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionMaybePtr2Maybe_(t *testing.T) {
  var src = 1
  var srcmp = union.Of(&src)
  var tgt union.Maybe[string]

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", tgt.Get())
}

func TestLookupConversionMaybePtr2MaybeBase_(t *testing.T) {
  var src = 1
  var srcmp = union.Of(&src)
  type substring string
  var tgt union.Maybe[substring]

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), tgt.Get())
}

func TestLookupConversionMaybePtr2MaybePtr_(t *testing.T) {
  var src = 1
  var srcmp = union.Of(&src)
  var tgt union.Maybe[*string]

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", *tgt.Get())
}

func TestLookupConversionMaybePtr2MaybePtrBase_(t *testing.T) {
  var src = 1
  var srcmp = union.Of(&src)
  type substring string
  var tgt union.Maybe[*substring]

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), *tgt.Get())
}

// ==== LookupConversion from Maybe Ptr Base

func TestLookupConversionMaybePtrBase2Val_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcmp = union.Of(&src)
  var tgt string

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, &tgt)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionMaybePtrBase2Base_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcmp = union.Of(&src)
  type substring string
  var tgt substring

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, &tgt)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionMaybePtrBase2Ptr_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcmp = union.Of(&src)
  var tgt string
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, tgtp)
  assert.Equal(t, "1", tgt)
}

func TestLookupConversionMaybePtrBase2PtrBase_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcmp = union.Of(&src)
  type substring string
  var tgt substring
  var tgtp = &tgt

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgtp))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, tgtp)
  assert.Equal(t, substring("1"), tgt)
}

func TestLookupConversionMaybePtrBase2Maybe_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcmp = union.Of(&src)
  var tgt union.Maybe[string]

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", tgt.Get())
}

func TestLookupConversionMaybePtrBase2MaybeBase_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcmp = union.Of(&src)
  type substring string
  var tgt union.Maybe[substring]

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), tgt.Get())
}

func TestLookupConversionMaybePtrBase2MaybePtr_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcmp = union.Of(&src)
  var tgt union.Maybe[*string]

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, "1", *tgt.Get())
}

func TestLookupConversionMaybePtrBase2MaybePtrBase_(t *testing.T) {
  type subint int
  var src = subint(1)
  var srcmp = union.Of(&src)
  type substring string
  var tgt union.Maybe[*substring]

  fn, err := LookupConversion(goreflect.TypeOf(srcmp), goreflect.TypeOf(tgt))

  assert.NotNil(t, fn)
  assert.Nil(t, err)

  fn(srcmp, &tgt)
  assert.True(t, tgt.Present())
  assert.Equal(t, substring("1"), *tgt.Get())
}

// ==== Other functions

func TestRegisterConversion_(t *testing.T) {
	type Conv_Reg_Foo struct{ fld int }

	{
		// Working conversion
		fn := func(src int, tgt *Conv_Reg_Foo) error { (*tgt).fld = src; return nil }
		assert.Nil(t, RegisterConversion(fn))
		var f Conv_Reg_Foo
		assert.Nil(t, To(5, &f))
		assert.Equal(t, Conv_Reg_Foo{5}, f)
	}

	{
		// Working conversion
		fn := func(src uint, tgt *Conv_Reg_Foo) error { (*tgt).fld = int(src); return nil }
		MustRegisterConversion(fn)
		var f Conv_Reg_Foo
		assert.Nil(t, To(uint(6), &f))
		assert.Equal(t, Conv_Reg_Foo{6}, f)

		// Can't register same conversion twice
		assert.Equal(t, fmt.Errorf("The conversion from uint to conv.Conv_Reg_Foo has already been registered"), RegisterConversion(fn))
	}

	{
		// Conversion for same type
		fn := func(src Conv_Reg_Foo, tgt *Conv_Reg_Foo) error { (*tgt).fld = src.fld + 1; return nil }
		MustRegisterConversion(fn)
		var f Conv_Reg_Foo
		assert.Nil(t, To(Conv_Reg_Foo{7}, &f))
		assert.Equal(t, Conv_Reg_Foo{8}, f)
	}
}

func TestTo_(t *testing.T) {
	// == int
	{
		var i int

		// ints
		assert.Nil(t, To(-1, &i))
		assert.Equal(t, -1, i)

		assert.Nil(t, To(int8(-2), &i))
		assert.Equal(t, -2, i)

		assert.Nil(t, To(int16(-3), &i))
		assert.Equal(t, -3, i)

		assert.Nil(t, To(int32(-4), &i))
		assert.Equal(t, -4, i)

		assert.Nil(t, To(int64(-5), &i))
		assert.Equal(t, -5, i)

		// uints
		assert.Nil(t, To(uint(1), &i))
		assert.Equal(t, 1, i)

		assert.Nil(t, To(uint8(2), &i))
		assert.Equal(t, 2, i)

		assert.Nil(t, To(uint16(3), &i))
		assert.Equal(t, 3, i)

		assert.Nil(t, To(uint32(4), &i))
		assert.Equal(t, 4, i)

		assert.Nil(t, To(uint64(5), &i))
		assert.Equal(t, 5, i)

		// floats
		assert.Nil(t, To(float32(1), &i))
		assert.Equal(t, 1, i)

		assert.Nil(t, To(2.0, &i))
		assert.Equal(t, 2, i)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &i))
		assert.Equal(t, 1, i)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(2))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 18446744073709551614 cannot be converted to int64"), To(bi, &i))

		funcs.TryTo(
			func() {
				MustTo(bi, &i)
				assert.Fail(t, "Never execute")
			},
			func(e any) {
				assert.Equal(t, "The *big.Int value of 18446744073709551614 cannot be converted to int64", e.(error).Error())
			},
		)

		assert.Nil(t, To(big.NewFloat(2), &i))
		assert.Equal(t, 2, i)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to int64"), To(big.NewFloat(1.25), &i))

		assert.Nil(t, To(big.NewRat(3, 1), &i))
		assert.Equal(t, 3, i)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to int64"), To(big.NewRat(5, 4), &i))

		// string
		assert.Nil(t, To("1", &i))
		assert.Equal(t, 1, i)

		assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "int64"), To("a", &i))
		assert.Equal(t, 1, i)
	}

	// == int8
	{
		var i8 int8

		// ints
		assert.Nil(t, To(-1, &i8))
		assert.Equal(t, int8(-1), i8)

		assert.Nil(t, To(int8(-2), &i8))
		assert.Equal(t, int8(-2), i8)

		assert.Nil(t, To(int16(-3), &i8))
		assert.Equal(t, int8(-3), i8)

		assert.Nil(t, To(int32(-4), &i8))
		assert.Equal(t, int8(-4), i8)

		assert.Nil(t, To(int64(-5), &i8))
		assert.Equal(t, int8(-5), i8)

		// uints
		assert.Nil(t, To(uint(1), &i8))
		assert.Equal(t, int8(1), i8)

		assert.Nil(t, To(uint8(2), &i8))
		assert.Equal(t, int8(2), i8)

		assert.Nil(t, To(uint16(3), &i8))
		assert.Equal(t, int8(3), i8)

		assert.Nil(t, To(uint32(4), &i8))
		assert.Equal(t, int8(4), i8)

		assert.Nil(t, To(uint64(5), &i8))
		assert.Equal(t, int8(5), i8)

		// floats
		assert.Nil(t, To(float32(1), &i8))
		assert.Equal(t, int8(1), i8)

		assert.Nil(t, To(2.0, &i8))
		assert.Equal(t, int8(2), i8)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &i8))
		assert.Equal(t, int8(1), i8)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(2))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 18446744073709551614 cannot be converted to int64"), To(bi, &i8))

		assert.Nil(t, To(big.NewFloat(2), &i8))
		assert.Equal(t, int8(2), i8)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to int64"), To(big.NewFloat(1.25), &i8))

		assert.Nil(t, To(big.NewRat(3, 1), &i8))
		assert.Equal(t, int8(3), i8)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to int64"), To(big.NewRat(5, 4), &i8))

		// string
		assert.Nil(t, To("1", &i8))
		assert.Equal(t, int8(1), i8)

		assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "int64"), To("a", &i8))
		assert.Equal(t, int8(1), i8)
	}

	// == int16
	{
		var i16 int16

		// ints
		assert.Nil(t, To(-1, &i16))
		assert.Equal(t, int16(-1), i16)

		assert.Nil(t, To(int8(-2), &i16))
		assert.Equal(t, int16(-2), i16)

		assert.Nil(t, To(int16(-3), &i16))
		assert.Equal(t, int16(-3), i16)

		assert.Nil(t, To(int32(-4), &i16))
		assert.Equal(t, int16(-4), i16)

		assert.Nil(t, To(int64(-5), &i16))
		assert.Equal(t, int16(-5), i16)

		// uints
		assert.Nil(t, To(uint(1), &i16))
		assert.Equal(t, int16(1), i16)

		assert.Nil(t, To(uint8(2), &i16))
		assert.Equal(t, int16(2), i16)

		assert.Nil(t, To(uint16(3), &i16))
		assert.Equal(t, int16(3), i16)

		assert.Nil(t, To(uint32(4), &i16))
		assert.Equal(t, int16(4), i16)

		assert.Nil(t, To(uint64(5), &i16))
		assert.Equal(t, int16(5), i16)

		// floats
		assert.Nil(t, To(float32(1), &i16))
		assert.Equal(t, int16(1), i16)

		assert.Nil(t, To(2.0, &i16))
		assert.Equal(t, int16(2), i16)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &i16))
		assert.Equal(t, int16(1), i16)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(2))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 18446744073709551614 cannot be converted to int64"), To(bi, &i16))

		assert.Nil(t, To(big.NewFloat(2), &i16))
		assert.Equal(t, int16(2), i16)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to int64"), To(big.NewFloat(1.25), &i16))

		assert.Nil(t, To(big.NewRat(3, 1), &i16))
		assert.Equal(t, int16(3), i16)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to int64"), To(big.NewRat(5, 4), &i16))

		// string
		assert.Nil(t, To("1", &i16))
		assert.Equal(t, int16(1), i16)

		assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "int64"), To("a", &i16))
		assert.Equal(t, int16(1), i16)
	}

	// == int32
	{
		var i32 int32

		// ints
		assert.Nil(t, To(-1, &i32))
		assert.Equal(t, int32(-1), i32)

		assert.Nil(t, To(int8(-2), &i32))
		assert.Equal(t, int32(-2), i32)

		assert.Nil(t, To(int16(-3), &i32))
		assert.Equal(t, int32(-3), i32)

		assert.Nil(t, To(int32(-4), &i32))
		assert.Equal(t, int32(-4), i32)

		assert.Nil(t, To(int64(-5), &i32))
		assert.Equal(t, int32(-5), i32)

		// uints
		assert.Nil(t, To(uint(1), &i32))
		assert.Equal(t, int32(1), i32)

		assert.Nil(t, To(uint8(2), &i32))
		assert.Equal(t, int32(2), i32)

		assert.Nil(t, To(uint16(3), &i32))
		assert.Equal(t, int32(3), i32)

		assert.Nil(t, To(uint32(4), &i32))
		assert.Equal(t, int32(4), i32)

		assert.Nil(t, To(uint64(5), &i32))
		assert.Equal(t, int32(5), i32)

		// floats
		assert.Nil(t, To(float32(1), &i32))
		assert.Equal(t, int32(1), i32)

		assert.Nil(t, To(2.0, &i32))
		assert.Equal(t, int32(2), i32)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &i32))
		assert.Equal(t, int32(1), i32)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(2))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 18446744073709551614 cannot be converted to int64"), To(bi, &i32))

		assert.Nil(t, To(big.NewFloat(2), &i32))
		assert.Equal(t, int32(2), i32)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to int64"), To(big.NewFloat(1.25), &i32))

		assert.Nil(t, To(big.NewRat(3, 1), &i32))
		assert.Equal(t, int32(3), i32)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to int64"), To(big.NewRat(5, 4), &i32))

		// string
		assert.Nil(t, To("1", &i32))
		assert.Equal(t, int32(1), i32)

		assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "int64"), To("a", &i32))
		assert.Equal(t, int32(1), i32)
	}

	// == int64
	{
		var i64 int64

		// ints
		assert.Nil(t, To(-1, &i64))
		assert.Equal(t, int64(-1), i64)

		assert.Nil(t, To(int8(-2), &i64))
		assert.Equal(t, int64(-2), i64)

		assert.Nil(t, To(int16(-3), &i64))
		assert.Equal(t, int64(-3), i64)

		assert.Nil(t, To(int32(-4), &i64))
		assert.Equal(t, int64(-4), i64)

		assert.Nil(t, To(int64(-5), &i64))
		assert.Equal(t, int64(-5), i64)

		// uints
		assert.Nil(t, To(uint(1), &i64))
		assert.Equal(t, int64(1), i64)

		assert.Nil(t, To(uint8(2), &i64))
		assert.Equal(t, int64(2), i64)

		assert.Nil(t, To(uint16(3), &i64))
		assert.Equal(t, int64(3), i64)

		assert.Nil(t, To(uint32(4), &i64))
		assert.Equal(t, int64(4), i64)

		assert.Nil(t, To(uint64(5), &i64))
		assert.Equal(t, int64(5), i64)

		// floats
		assert.Nil(t, To(float32(1), &i64))
		assert.Equal(t, int64(1), i64)

		assert.Nil(t, To(2.0, &i64))
		assert.Equal(t, int64(2), i64)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &i64))
		assert.Equal(t, int64(1), i64)

		assert.Nil(t, To(big.NewFloat(2), &i64))
		assert.Equal(t, int64(2), i64)

		assert.Nil(t, To(big.NewRat(3, 1), &i64))
		assert.Equal(t, int64(3), i64)

		// string
		assert.Nil(t, To("1", &i64))
		assert.Equal(t, int64(1), i64)

		assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "int64"), To("a", &i64))
		assert.Equal(t, int64(0), i64)
	}

	// == uint
	{
		var ui uint

		// ints
		assert.Nil(t, To(1, &ui))
		assert.Equal(t, uint(1), ui)

		assert.Nil(t, To(int8(2), &ui))
		assert.Equal(t, uint(2), ui)

		assert.Nil(t, To(int16(3), &ui))
		assert.Equal(t, uint(3), ui)

		assert.Nil(t, To(int32(4), &ui))
		assert.Equal(t, uint(4), ui)

		assert.Nil(t, To(int64(5), &ui))
		assert.Equal(t, uint(5), ui)

		// uints
		assert.Nil(t, To(uint(1), &ui))
		assert.Equal(t, uint(1), ui)

		assert.Nil(t, To(uint8(2), &ui))
		assert.Equal(t, uint(2), ui)

		assert.Nil(t, To(uint16(3), &ui))
		assert.Equal(t, uint(3), ui)

		assert.Nil(t, To(uint32(4), &ui))
		assert.Equal(t, uint(4), ui)

		assert.Nil(t, To(uint64(5), &ui))
		assert.Equal(t, uint(5), ui)

		// floats
		assert.Nil(t, To(float32(1), &ui))
		assert.Equal(t, uint(1), ui)

		assert.Nil(t, To(2.0, &ui))
		assert.Equal(t, uint(2), ui)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &ui))
		assert.Equal(t, uint(1), ui)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(4))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 36893488147419103228 cannot be converted to uint64"), To(bi, &ui))

		assert.Nil(t, To(big.NewFloat(2), &ui))
		assert.Equal(t, uint(2), ui)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to uint64"), To(big.NewFloat(1.25), &ui))

		assert.Nil(t, To(big.NewRat(3, 1), &ui))
		assert.Equal(t, uint(3), ui)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to uint64"), To(big.NewRat(5, 4), &ui))

		// string
		assert.Nil(t, To("1", &ui))
		assert.Equal(t, uint(1), ui)

		assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "uint64"), To("a", &ui))
		assert.Equal(t, uint(1), ui)
	}

	// == uint8
	{
		var ui8 uint8

		// ints
		assert.Nil(t, To(1, &ui8))
		assert.Equal(t, uint8(1), ui8)

		assert.Nil(t, To(int8(2), &ui8))
		assert.Equal(t, uint8(2), ui8)

		assert.Nil(t, To(int16(3), &ui8))
		assert.Equal(t, uint8(3), ui8)

		assert.Nil(t, To(int32(4), &ui8))
		assert.Equal(t, uint8(4), ui8)

		assert.Nil(t, To(int64(5), &ui8))
		assert.Equal(t, uint8(5), ui8)

		// uints
		assert.Nil(t, To(uint(1), &ui8))
		assert.Equal(t, uint8(1), ui8)

		assert.Nil(t, To(uint8(2), &ui8))
		assert.Equal(t, uint8(2), ui8)

		assert.Nil(t, To(uint16(3), &ui8))
		assert.Equal(t, uint8(3), ui8)

		assert.Nil(t, To(uint32(4), &ui8))
		assert.Equal(t, uint8(4), ui8)

		assert.Nil(t, To(uint64(5), &ui8))
		assert.Equal(t, uint8(5), ui8)

		// floats
		assert.Nil(t, To(float32(1), &ui8))
		assert.Equal(t, uint8(1), ui8)

		assert.Nil(t, To(2.0, &ui8))
		assert.Equal(t, uint8(2), ui8)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &ui8))
		assert.Equal(t, uint8(1), ui8)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(4))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 36893488147419103228 cannot be converted to uint64"), To(bi, &ui8))

		assert.Nil(t, To(big.NewFloat(2), &ui8))
		assert.Equal(t, uint8(2), ui8)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to uint64"), To(big.NewFloat(1.25), &ui8))

		assert.Nil(t, To(big.NewRat(3, 1), &ui8))
		assert.Equal(t, uint8(3), ui8)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to uint64"), To(big.NewRat(5, 4), &ui8))

		// string
		assert.Nil(t, To("1", &ui8))
		assert.Equal(t, uint8(1), ui8)

		assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "uint64"), To("a", &ui8))
		assert.Equal(t, uint8(1), ui8)
	}

	// == uint16
	{
		var ui16 uint16

		// ints
		assert.Nil(t, To(1, &ui16))
		assert.Equal(t, uint16(1), ui16)

		assert.Nil(t, To(int8(2), &ui16))
		assert.Equal(t, uint16(2), ui16)

		assert.Nil(t, To(int16(3), &ui16))
		assert.Equal(t, uint16(3), ui16)

		assert.Nil(t, To(int32(4), &ui16))
		assert.Equal(t, uint16(4), ui16)

		assert.Nil(t, To(int64(5), &ui16))
		assert.Equal(t, uint16(5), ui16)

		// uints
		assert.Nil(t, To(uint(1), &ui16))
		assert.Equal(t, uint16(1), ui16)

		assert.Nil(t, To(uint8(2), &ui16))
		assert.Equal(t, uint16(2), ui16)

		assert.Nil(t, To(uint16(3), &ui16))
		assert.Equal(t, uint16(3), ui16)

		assert.Nil(t, To(uint32(4), &ui16))
		assert.Equal(t, uint16(4), ui16)

		assert.Nil(t, To(uint64(5), &ui16))
		assert.Equal(t, uint16(5), ui16)

		// floats
		assert.Nil(t, To(float32(1), &ui16))
		assert.Equal(t, uint16(1), ui16)

		assert.Nil(t, To(2.0, &ui16))
		assert.Equal(t, uint16(2), ui16)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &ui16))
		assert.Equal(t, uint16(1), ui16)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(4))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 36893488147419103228 cannot be converted to uint64"), To(bi, &ui16))

		assert.Nil(t, To(big.NewFloat(2), &ui16))
		assert.Equal(t, uint16(2), ui16)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to uint64"), To(big.NewFloat(1.25), &ui16))

		assert.Nil(t, To(big.NewRat(3, 1), &ui16))
		assert.Equal(t, uint16(3), ui16)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to uint64"), To(big.NewRat(5, 4), &ui16))

		// string
		assert.Nil(t, To("1", &ui16))
		assert.Equal(t, uint16(1), ui16)

		assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "uint64"), To("a", &ui16))
		assert.Equal(t, uint16(1), ui16)
	}

	// == uint32
	{
		var ui32 uint32

		// ints
		assert.Nil(t, To(1, &ui32))
		assert.Equal(t, uint32(1), ui32)

		assert.Nil(t, To(int8(2), &ui32))
		assert.Equal(t, uint32(2), ui32)

		assert.Nil(t, To(int16(3), &ui32))
		assert.Equal(t, uint32(3), ui32)

		assert.Nil(t, To(int32(4), &ui32))
		assert.Equal(t, uint32(4), ui32)

		assert.Nil(t, To(int64(5), &ui32))
		assert.Equal(t, uint32(5), ui32)

		// uints
		assert.Nil(t, To(uint(1), &ui32))
		assert.Equal(t, uint32(1), ui32)

		assert.Nil(t, To(uint8(2), &ui32))
		assert.Equal(t, uint32(2), ui32)

		assert.Nil(t, To(uint16(3), &ui32))
		assert.Equal(t, uint32(3), ui32)

		assert.Nil(t, To(uint32(4), &ui32))
		assert.Equal(t, uint32(4), ui32)

		assert.Nil(t, To(uint64(5), &ui32))
		assert.Equal(t, uint32(5), ui32)

		// floats
		assert.Nil(t, To(float32(1), &ui32))
		assert.Equal(t, uint32(1), ui32)

		assert.Nil(t, To(2.0, &ui32))
		assert.Equal(t, uint32(2), ui32)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &ui32))
		assert.Equal(t, uint32(1), ui32)

		bi := big.NewInt(math.MaxInt64)
		bi = bi.Mul(bi, big.NewInt(4))
		assert.Equal(t, fmt.Errorf("The *big.Int value of 36893488147419103228 cannot be converted to uint64"), To(bi, &ui32))

		assert.Nil(t, To(big.NewFloat(2), &ui32))
		assert.Equal(t, uint32(2), ui32)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 1.25 cannot be converted to uint64"), To(big.NewFloat(1.25), &ui32))

		assert.Nil(t, To(big.NewRat(3, 1), &ui32))
		assert.Equal(t, uint32(3), ui32)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 5/4 cannot be converted to uint64"), To(big.NewRat(5, 4), &ui32))

		// string
		assert.Nil(t, To("1", &ui32))
		assert.Equal(t, uint32(1), ui32)

		assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "uint64"), To("a", &ui32))
		assert.Equal(t, uint32(1), ui32)
	}

	// == uint64
	{
		var ui64 uint64

		// ints
		assert.Nil(t, To(1, &ui64))
		assert.Equal(t, uint64(1), ui64)

		assert.Nil(t, To(int8(2), &ui64))
		assert.Equal(t, uint64(2), ui64)

		assert.Nil(t, To(int16(3), &ui64))
		assert.Equal(t, uint64(3), ui64)

		assert.Nil(t, To(int32(4), &ui64))
		assert.Equal(t, uint64(4), ui64)

		assert.Nil(t, To(int64(5), &ui64))
		assert.Equal(t, uint64(5), ui64)

		// uints
		assert.Nil(t, To(uint(1), &ui64))
		assert.Equal(t, uint64(1), ui64)

		assert.Nil(t, To(uint8(2), &ui64))
		assert.Equal(t, uint64(2), ui64)

		assert.Nil(t, To(uint16(3), &ui64))
		assert.Equal(t, uint64(3), ui64)

		assert.Nil(t, To(uint32(4), &ui64))
		assert.Equal(t, uint64(4), ui64)

		assert.Nil(t, To(uint64(5), &ui64))
		assert.Equal(t, uint64(5), ui64)

		// floats
		assert.Nil(t, To(float32(1), &ui64))
		assert.Equal(t, uint64(1), ui64)

		assert.Nil(t, To(2.0, &ui64))
		assert.Equal(t, uint64(2), ui64)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &ui64))
		assert.Equal(t, uint64(1), ui64)

		assert.Nil(t, To(big.NewFloat(2), &ui64))
		assert.Equal(t, uint64(2), ui64)

		assert.Nil(t, To(big.NewRat(3, 1), &ui64))
		assert.Equal(t, uint64(3), ui64)

		// string
		assert.Nil(t, To("1", &ui64))
		assert.Equal(t, uint64(1), ui64)

		assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "uint64"), To("a", &ui64))
		assert.Equal(t, uint64(0), ui64)
	}

	// == float32
	{
		var f32 float32

		// ints
		assert.Nil(t, To(1, &f32))
		assert.Equal(t, float32(1), f32)

		assert.Nil(t, To(int8(2), &f32))
		assert.Equal(t, float32(2), f32)

		assert.Nil(t, To(int16(3), &f32))
		assert.Equal(t, float32(3), f32)

		assert.Nil(t, To(int32(4), &f32))
		assert.Equal(t, float32(4), f32)

		assert.Nil(t, To(int64(5), &f32))
		assert.Equal(t, float32(5), f32)

		// uints
		assert.Nil(t, To(uint(1), &f32))
		assert.Equal(t, float32(1), f32)

		assert.Nil(t, To(uint8(2), &f32))
		assert.Equal(t, float32(2), f32)

		assert.Nil(t, To(uint16(3), &f32))
		assert.Equal(t, float32(3), f32)

		assert.Nil(t, To(uint32(4), &f32))
		assert.Equal(t, float32(4), f32)

		assert.Nil(t, To(uint64(5), &f32))
		assert.Equal(t, float32(5), f32)

		// floats
		assert.Nil(t, To(float32(1.25), &f32))
		assert.Equal(t, float32(1.25), f32)

		assert.Nil(t, To(2.5, &f32))
		assert.Equal(t, float32(2.5), f32)

		// *bigs
		assert.Nil(t, To(big.NewInt(1), &f32))
		assert.Equal(t, float32(1), f32)
		assert.Equal(t, fmt.Errorf("The *big.Int value of 9223372036854775807 cannot be converted to float64"), To(big.NewInt(math.MaxInt64), &f32))

		assert.Nil(t, To(big.NewFloat(1.25), &f32))
		assert.Equal(t, float32(1.25), f32)

		bf := big.NewFloat(0)
		IntToBigFloat(math.MaxInt64, &bf)
		assert.Equal(t, fmt.Errorf("The *big.Float value of 9223372036854775807 cannot be converted to float64"), To(bf, &f32))

		assert.Nil(t, To(big.NewRat(250, 100), &f32))
		assert.Equal(t, float32(2.5), f32)
		assert.Equal(t, fmt.Errorf("The *big.Rat value of 9223372036854775807/1 cannot be converted to float64"), To(big.NewRat(math.MaxInt64, 1), &f32))

		// string
		assert.Nil(t, To("1.25", &f32))
		assert.Equal(t, float32(1.25), f32)

		assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "float32"), To("a", &f32))
		assert.Equal(t, float32(1.25), f32)
	}

	// == float64
	{
		var f64 float64

		// ints
		assert.Nil(t, To(1, &f64))
		assert.Equal(t, 1.0, f64)

		assert.Nil(t, To(int8(2), &f64))
		assert.Equal(t, 2.0, f64)

		assert.Nil(t, To(int16(3), &f64))
		assert.Equal(t, 3.0, f64)

		assert.Nil(t, To(int32(4), &f64))
		assert.Equal(t, 4.0, f64)

		assert.Nil(t, To(int64(5), &f64))
		assert.Equal(t, 5.0, f64)

		// uints
		assert.Nil(t, To(uint(1), &f64))
		assert.Equal(t, 1.0, f64)

		assert.Nil(t, To(uint8(2), &f64))
		assert.Equal(t, 2.0, f64)

		assert.Nil(t, To(uint16(3), &f64))
		assert.Equal(t, 3.0, f64)

		assert.Nil(t, To(uint32(4), &f64))
		assert.Equal(t, 4.0, f64)

		assert.Nil(t, To(uint64(5), &f64))
		assert.Equal(t, 5.0, f64)

		// floats
		assert.Nil(t, To(float32(1.25), &f64))
		assert.Equal(t, 1.25, f64)

		assert.Nil(t, To(2.5, &f64))
		assert.Equal(t, 2.5, f64)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &f64))
		assert.Equal(t, 1.0, f64)

		assert.Nil(t, To(big.NewFloat(1.25), &f64))
		assert.Equal(t, 1.25, f64)

		assert.Nil(t, To(big.NewRat(250, 100), &f64))
		assert.Equal(t, 2.5, f64)

		// string
		assert.Nil(t, To("1.25", &f64))
		assert.Equal(t, 1.25, f64)

		assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "float64"), To("a", &f64))
		assert.Equal(t, 1.25, f64)
	}

	// == *big.Int
	{
		var bi *big.Int

		// ints
		assert.Nil(t, To(1, &bi))
		assert.Equal(t, big.NewInt(1), bi)

		assert.Nil(t, To(int8(2), &bi))
		assert.Equal(t, big.NewInt(2), bi)

		assert.Nil(t, To(int16(3), &bi))
		assert.Equal(t, big.NewInt(3), bi)

		assert.Nil(t, To(int32(4), &bi))
		assert.Equal(t, big.NewInt(4), bi)

		assert.Nil(t, To(int64(5), &bi))
		assert.Equal(t, big.NewInt(5), bi)

		// uints
		assert.Nil(t, To(uint(1), &bi))
		assert.Equal(t, big.NewInt(1), bi)

		assert.Nil(t, To(uint8(2), &bi))
		assert.Equal(t, big.NewInt(2), bi)

		assert.Nil(t, To(uint16(3), &bi))
		assert.Equal(t, big.NewInt(3), bi)

		assert.Nil(t, To(uint32(4), &bi))
		assert.Equal(t, big.NewInt(4), bi)

		assert.Nil(t, To(uint64(5), &bi))
		assert.Equal(t, big.NewInt(5), bi)

		// floats
		assert.Nil(t, To(float32(125), &bi))
		assert.Equal(t, big.NewInt(125), bi)

		assert.Nil(t, To(25.0, &bi))
		assert.Equal(t, big.NewInt(25), bi)

		// bigs
		bisrc := big.NewInt(1)
		assert.Nil(t, To(bisrc, &bi))
		assert.False(t, bisrc == bi)
		assert.Equal(t, big.NewInt(1), bi)

		assert.Nil(t, To(big.NewFloat(125), &bi))
		assert.Equal(t, big.NewInt(125), bi)

		assert.Nil(t, To(big.NewRat(250, 1), &bi))
		assert.Equal(t, big.NewInt(250), bi)

		// string
		assert.Nil(t, To("1", &bi))
		assert.Equal(t, big.NewInt(1), bi)

		assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "*big.Int"), To("a", &bi))
		assert.Equal(t, big.NewInt(0), bi)
	}

	// == *big.Float
	{
		var bf *big.Float

		// ints
		assert.Nil(t, To(1, &bf))
		cmp := big.NewFloat(1)
		cmp.SetPrec(4)
		assert.Equal(t, cmp, bf)

		assert.Nil(t, To(int8(2), &bf))
		cmp = big.NewFloat(2)
		cmp.SetPrec(4)
		assert.Equal(t, cmp, bf)

		assert.Nil(t, To(int16(3), &bf))
		cmp = big.NewFloat(3)
		cmp.SetPrec(4)
		assert.Equal(t, cmp, bf)

		assert.Nil(t, To(int32(4), &bf))
		cmp = big.NewFloat(4)
		cmp.SetPrec(4)
		assert.Equal(t, cmp, bf)

		assert.Nil(t, To(int64(5), &bf))
		cmp = big.NewFloat(5)
		cmp.SetPrec(4)
		assert.Equal(t, cmp, bf)

		// uints
		assert.Nil(t, To(uint(1), &bf))
		assert.Equal(t, big.NewFloat(1), bf)

		assert.Nil(t, To(uint8(2), &bf))
		assert.Equal(t, big.NewFloat(2), bf)

		assert.Nil(t, To(uint16(3), &bf))
		assert.Equal(t, big.NewFloat(3), bf)

		assert.Nil(t, To(uint32(4), &bf))
		assert.Equal(t, big.NewFloat(4), bf)

		assert.Nil(t, To(uint64(5), &bf))
		assert.Equal(t, big.NewFloat(5), bf)

		// floats
		assert.Nil(t, To(float32(1.25), &bf))
		assert.Equal(t, big.NewFloat(1.25), bf)

		assert.Nil(t, To(2.5, &bf))
		assert.Equal(t, big.NewFloat(2.5), bf)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &bf))
		assert.Equal(t, big.NewFloat(1), bf)

		bfsrc := big.NewFloat(1.25)
		assert.Nil(t, To(bfsrc, &bf))
		assert.False(t, bfsrc == bf)
		assert.Equal(t, big.NewFloat(1.25), bf)

		assert.Nil(t, To(big.NewRat(250, 100), &bf))
		assert.Equal(t, big.NewFloat(2.5), bf)

		// string
		assert.Nil(t, To("1.25", &bf))
		assert.Equal(t, big.NewFloat(1.25), bf)

		assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "*big.Float"), To("a", &bf))
		assert.Equal(t, (*big.Float)(nil), bf)
	}

	// == *big.Rat
	{
		var br *big.Rat

		// ints
		assert.Nil(t, To(1, &br))
		assert.Equal(t, big.NewRat(1, 1), br)

		assert.Nil(t, To(int8(2), &br))
		assert.Equal(t, big.NewRat(2, 1), br)

		assert.Nil(t, To(int16(3), &br))
		assert.Equal(t, big.NewRat(3, 1), br)

		assert.Nil(t, To(int32(4), &br))
		assert.Equal(t, big.NewRat(4, 1), br)

		assert.Nil(t, To(int64(5), &br))
		assert.Equal(t, big.NewRat(5, 1), br)

		// uints
		assert.Nil(t, To(uint(1), &br))
		assert.Equal(t, big.NewRat(1, 1), br)

		assert.Nil(t, To(uint8(2), &br))
		assert.Equal(t, big.NewRat(2, 1), br)

		assert.Nil(t, To(uint16(3), &br))
		assert.Equal(t, big.NewRat(3, 1), br)

		assert.Nil(t, To(uint32(4), &br))
		assert.Equal(t, big.NewRat(4, 1), br)

		assert.Nil(t, To(uint64(5), &br))
		assert.Equal(t, big.NewRat(5, 1), br)

		// floats
		assert.Nil(t, To(float32(1.25), &br))
		assert.Equal(t, big.NewRat(125, 100), br)

		assert.Nil(t, To(2.5, &br))
		assert.Equal(t, big.NewRat(25, 10), br)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &br))
		assert.Equal(t, big.NewRat(1, 1), br)

		assert.Nil(t, To(big.NewFloat(1.25), &br))
		assert.Equal(t, big.NewRat(125, 100), br)

		brsrc := big.NewRat(25, 10)
		assert.Nil(t, To(brsrc, &br))
		assert.False(t, brsrc == br)
		assert.Equal(t, big.NewRat(25, 10), br)

		// string
		assert.Nil(t, To("5/4", &br))
		assert.Equal(t, big.NewRat(5, 4), br)

		assert.Equal(t, fmt.Errorf(errMsg, "a", "a", "*big.Rat"), To("a", &br))
		assert.Equal(t, (*big.Rat)(nil), br)
	}

	// == string
	{
		var s string

		// ints
		assert.Nil(t, To(1, &s))
		assert.Equal(t, "1", s)

		assert.Nil(t, To(int8(2), &s))
		assert.Equal(t, "2", s)

		assert.Nil(t, To(int16(3), &s))
		assert.Equal(t, "3", s)

		assert.Nil(t, To(int32(4), &s))
		assert.Equal(t, "4", s)

		assert.Nil(t, To(int64(5), &s))
		assert.Equal(t, "5", s)

		// uints
		assert.Nil(t, To(uint(1), &s))
		assert.Equal(t, "1", s)

		assert.Nil(t, To(uint8(2), &s))
		assert.Equal(t, "2", s)

		assert.Nil(t, To(uint16(3), &s))
		assert.Equal(t, "3", s)

		assert.Nil(t, To(uint32(4), &s))
		assert.Equal(t, "4", s)

		assert.Nil(t, To(uint64(5), &s))
		assert.Equal(t, "5", s)

		// floats
		assert.Nil(t, To(float32(1.25), &s))
		assert.Equal(t, "1.25", s)

		assert.Nil(t, To(2.5, &s))
		assert.Equal(t, "2.5", s)

		// bigs
		assert.Nil(t, To(big.NewInt(1), &s))
		assert.Equal(t, "1", s)

		assert.Nil(t, To(big.NewFloat(1.25), &s))
		assert.Equal(t, "1.25", s)

		assert.Nil(t, To(big.NewRat(25, 10), &s))
		assert.Equal(t, "5/2", s)

		// string
		assert.Nil(t, To("foo", &s))
		assert.Equal(t, "foo", s)
	}

	// source type = target type (int -> int)
	{
		var o int
		assert.Nil(t, To(1, &o))
		assert.Equal(t, 1, o)
	}

	// Derfd source type to other target type with a conversion (*int -> string)
	{
		var i int = 1
		var o string

		// source exists
		assert.Nil(t, To(&i, &o))
		assert.Equal(t, "1", o)

		// source cannot be nil
		assert.Equal(t, fmt.Errorf("A nil *int cannot be converted to a(n) string"), To((*int)(nil), &o))
		assert.Equal(t, "1", o)
	}

	// Derefd source type = target type (*int -> int)
	{
		var i int = 1
		var o int
		assert.Nil(t, To(&i, &o))
		assert.Equal(t, 1, o)

		// source cannot be nil
		assert.Equal(t, fmt.Errorf("A nil *int cannot be copied to a(n) int"), To((*int)(nil), &o))
		assert.Equal(t, 1, o)
	}

	// source type = derefd target type (int -> *int)
	{
		var o int
		var po = &o
		assert.Nil(t, To(1, &po))
		assert.Equal(t, 1, o)

		// target cannot be nil
		po = nil
		assert.Equal(t, fmt.Errorf("A(n) int cannot be copied to a nil *int"), To(2, &po))
		assert.Equal(t, 1, o)

		assert.Equal(t, fmt.Errorf("A(n) int cannot be copied to a nil *int"), To(2, (**int)(nil)))
		assert.Equal(t, 1, o)
	}

	// derefd source type = derefd target type (*int -> *int)
	{
		var i int
		var o int
		var po *int

		// source is nil, target is not nil
		i = 1
		po = &o
		assert.Nil(t, To((*int)(nil), &po))
		assert.Nil(t, po)
		assert.Equal(t, 0, o)

		// source is not nil, target is not nil
		i = 2
		po = &o
		assert.Nil(t, To(&i, &po))
		assert.Equal(t, &o, po)
		assert.Equal(t, 2, i)

		// source is nil, target is **nil
		i = 3
		po = &o
		assert.Equal(t, fmt.Errorf("A(n) *int cannot be copied to a nil *int"), To((*int)(nil), (**int)(nil)))
		assert.Equal(t, &o, po)
		assert.Equal(t, 3, i)

		// source is nil, target is *nil
		i = 4
		po = nil
		assert.Nil(t, To((*int)(nil), &po))
		assert.Nil(t, po)
		assert.Equal(t, 4, i)
	}

	{
		// byte to rune, which is really uint8 to int32
		// it is not a subtype, reflection sees uint8 and int32
		var r rune
		assert.Nil(t, To(byte('A'), &r))
		assert.Equal(t, 'A', r)
	}

	{
		var c chan bool
		assert.Equal(t, fmt.Errorf("string cannot be converted to chan bool"), To("str", &c))
	}

	{
		type Conv_To_Foo struct{ Bar int }
		var f Conv_To_Foo
		assert.Equal(t, fmt.Errorf("int cannot be converted to conv.Conv_To_Foo"), To(1, &f))
	}

	// Subtypes where no conversion exists, base types are the same
	{
		type foo int
		type bar int
		var b bar
		assert.Nil(t, To(foo(1), &b))
		assert.Equal(t, bar(1), b)
	}

	// Subtypes where no conversion exists, base types are different
	{
		type foo uint
		type bar int
		var b bar
		assert.Nil(t, To(foo(1), &b))
		assert.Equal(t, bar(1), b)
	}
}

func TestToBigOps_(t *testing.T) {
	{
		var bi *big.Int
		assert.Nil(t, ToBigOps(1, &bi))
		assert.Equal(t, big.NewInt(1), bi)

		assert.Nil(t, ToBigOps(bi, &bi))
		assert.Equal(t, big.NewInt(1), bi)

		// byte to *big.Int, which is relly uint8 to *big.Int
		// verify subtypes are handled correctly
		assert.Nil(t, To(byte('A'), &bi))
		assert.Equal(t, big.NewInt('A'), bi)
	}

	{
		var bf *big.Float
		assert.Nil(t, ToBigOps(2, &bf))
		cmp := big.NewFloat(2)
		cmp.SetPrec(uint(math.Ceil(1 * log2Of10)))
		assert.Equal(t, cmp, bf)

		assert.Nil(t, ToBigOps(bf, &bf))
		assert.Equal(t, cmp, bf)
	}

	{
		var br *big.Rat
		assert.Nil(t, ToBigOps(3, &br))
		assert.Equal(t, big.NewRat(3, 1), br)

		assert.Nil(t, ToBigOps(br, &br))
		assert.Equal(t, big.NewRat(3, 1), br)
	}

	{
		var br *big.Rat
		assert.Nil(t, ToBigOps(3, &br))
		assert.Equal(t, big.NewRat(3, 1), br)

		assert.Nil(t, ToBigOps(br, &br))
		assert.Equal(t, big.NewRat(3, 1), br)
	}

	funcs.TryTo(
		func() {
			var br *big.Int
			MustToBigOps("a", &br)
			assert.Fail(t, "Never execute")
		},
		func(e any) {
			assert.Equal(t, fmt.Sprintf("The string value of a cannot be converted to *big.Int"), e.(error).Error())
		},
	)
}
