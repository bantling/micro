package build

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/bantling/micro/go/json"
)

type ObjectBuilder struct {
	val map[string]json.JSONValue
}

type ObjectKeyBuilder struct {
	p *ObjectBuilder
}

type ObjectValueBuilder struct {
	p *ObjectBuilder
}

func Object() *ObjectBuilder {
	return &ObjectBuilder{val: map[string]json.JSONValue{}}
}

func (ob *ObjectBuilder) Map(key string) ObjectKeyBuilder {
	return ObjectKeyBuilder{p: ob}
}

func (okb ObjectKeyBuilder) To() ObjectValueBuilder {
	return ObjectValueBuilder{p: okb.p}
}

func (ovb ObjectValueBuilder) Object() *ObjectBuilder {
	return Object()
}
