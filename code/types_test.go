package code

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/bantling/micro/union"
	"github.com/stretchr/testify/assert"
)

func TestType_(t *testing.T) {
	assert.Equal(t, uint(afterScalar), uint(JSON))
	assert.Equal(t, uint(afterScalar+1), uint(Array))

	assert.True(t, IsScalar(Type(Bool)))
	assert.True(t, IsScalar(Type(Uint)))
	assert.False(t, IsScalar(Type(JSON)))
	assert.False(t, IsScalar(Type(Array)))

	assert.False(t, IsAggregate(Type(Bool)))
	assert.False(t, IsAggregate(Type(Uint)))
	assert.True(t, IsAggregate(Type(JSON)))
	assert.True(t, IsAggregate(Type(Array)))

	assert.True(t, IsScalar(OfScalarType(Bool).Typ))
	assert.True(t, IsAggregate(OfJSONType().Typ))
	assert.True(t, IsAggregate(OfEnumType("e", OfScalarType(Uint), []string{"foo"}).Typ))
	assert.True(t, IsAggregate(OfListType(OfScalarType(Uint)).Typ))
	assert.True(t, IsAggregate(OfMapType(OfScalarType(Uint), OfScalarType(String)).Typ))
	assert.True(t, IsAggregate(OfMaybeType(OfScalarType(Uint)).Typ))
	assert.True(t, IsAggregate(OfObjectType(OfScalarType(Uint), "foo", []string{"bar"}).Typ))
	assert.True(t, IsAggregate(OfSetType(OfScalarType(Uint)).Typ))
}

func TestVal_(t *testing.T) {
	typ := OfScalarType(Int)
	assert.Equal(
		t,
		Val{Access: union.Empty[AccessLevel](), Kind: LitVal, Typ: typ, Value: "1"},
		OfLitVal(typ, "1"),
	)

	assert.Equal(
		t,
		Val{Access: union.Of(Public), Kind: VarVal, Typ: typ, Value: "2"},
		OfVarVal(false, typ, "2", Public),
	)

	assert.Equal(
		t,
		Val{Access: union.Of(Private), Kind: VarVal, Typ: typ, Value: "3"},
		OfVarVal(false, typ, "3", Private),
	)

	assert.Equal(
		t,
		Val{Access: union.Empty[AccessLevel](), Kind: VarConst, Typ: typ, Value: "4"},
		OfVarVal(true, typ, "4"),
	)
}
