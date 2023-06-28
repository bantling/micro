package golang

// SPDX-License-Identifier: Apache-2.0

import (
  "testing"

  "github.com/bantling/micro/generate"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
  dt := DataType{
    Comment: "Foo is very fooey",
    Name: "Foo",
    Fields: []Field{
      Field{
        Comment: "Bar is very bary",
        Name: "Bar",
        Type: generate.String_t,
      },
    },
  }

  es := `// Foo is very fooey
type Foo struct {
  Bar string // Bar is very bary
}
`

  assert.Equal(t, es, dt.String())
}
