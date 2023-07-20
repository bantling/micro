package golang

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/bantling/micro/generate"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	var (
		dt generate.DataType
		es string
	)

	// One string field, no comments
	dt = generate.DataType{
		Name: "Foo",
		Fields: []generate.Field{
			{
				Name: "Bar",
				Type: generate.String_t,
			},
		},
	}

	es = `type Foo struct {
  Bar string
}
`

	assert.Equal(t, es, DataTypeString(dt))

	// One []*string field, comments on type and field
	dt = generate.DataType{
		Comment: "Foo is very fooey",
		Name:    "Foo",
		Fields: []generate.Field{
			{
				Comment: "Bar is very bary",
				Name:    "Bar",
				Type:    generate.String_t,
				Array:   true,
				Ref:     true,
			},
		},
	}

	es = `// Foo is very fooey
type Foo struct {
  Bar []*string // Bar is very bary
}
`

	assert.Equal(t, es, DataTypeString(dt))
}
