package app

// SPDX-License-Identifier: Apache-2.0

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/bantling/micro/funcs"
	"github.com/bantling/micro/tuple"
	"github.com/stretchr/testify/assert"
)

func TestStringToTypeDef_(t *testing.T) {
	// typ FieldType, length, scale int, nullable bool
	assert.Equal(t, tuple.Of4(stringToTypeDef("bool")), tuple.Of4(Bool, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("bool?")), tuple.Of4(Bool, 0, 0, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("date")), tuple.Of4(Date, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("date?")), tuple.Of4(Date, 0, 0, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("decimal(5)")), tuple.Of4(Decimal, 5, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("decimal(5,2)")), tuple.Of4(Decimal, 5, 2, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("decimal(5, 2)?")), tuple.Of4(Decimal, 5, 2, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("float32")), tuple.Of4(Float32, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("float32?")), tuple.Of4(Float32, 0, 0, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("float64")), tuple.Of4(Float64, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("float64?")), tuple.Of4(Float64, 0, 0, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("int32")), tuple.Of4(Int32, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("int32?")), tuple.Of4(Int32, 0, 0, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("int64")), tuple.Of4(Int64, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("int64?")), tuple.Of4(Int64, 0, 0, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("interval")), tuple.Of4(Interval, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("interval?")), tuple.Of4(Interval, 0, 0, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("json")), tuple.Of4(JSON, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("json?")), tuple.Of4(JSON, 0, 0, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("ref:one")), tuple.Of4(RefOne, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("ref:one?")), tuple.Of4(RefOne, 0, 0, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("ref:many")), tuple.Of4(RefManyToOne, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("ref:many?")), tuple.Of4(RefManyToOne, 0, 0, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("ref:manyToMany")), tuple.Of4(RefManyToMany, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("ref:manyToMany?")), tuple.Of4(RefManyToMany, 0, 0, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("string")), tuple.Of4(String, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("string(5)?")), tuple.Of4(String, 5, 0, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("table_row")), tuple.Of4(TableRow, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("table_row?")), tuple.Of4(TableRow, 0, 0, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("time")), tuple.Of4(Time, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("time?")), tuple.Of4(Time, 0, 0, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("timestamp")), tuple.Of4(Timestamp, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("timestamp?")), tuple.Of4(Timestamp, 0, 0, true))

	assert.Equal(t, tuple.Of4(stringToTypeDef("uuid")), tuple.Of4(UUID, 0, 0, false))
	assert.Equal(t, tuple.Of4(stringToTypeDef("uuid?")), tuple.Of4(UUID, 0, 0, true))
}

func TestLoadAllValues_(t *testing.T) {
	// Load that specifies all values
	var (
		data = strings.NewReader(`
[database_]
name = "mydb"
description = "my great database"
locale = "en_CA"
accent_sensitive = false
case_sensitive = false
schemas = ["app1", "app2"]
vendors = ["postgres"]
vendor_types = {"whatever" = {"postgres" = "psqlType"}, sowhat = {postgres = "psqlType"}}

[country]
id = "uuid"
name = "string"
code = "string"
has_regions = "bool"
descriptor_ = {terms = ["name", "code"], description = "$name"}
unique_ = [["name"], ["code"]]

[region]
id = "uuid"
country = "ref:many"
name = "string"
code = "string"
descriptor_ = {terms = ["name", "code"], description = "$code"}
unique_ = [["country", "name"], ["country", "code"]]

[address]
id = "uuid"
country = "ref:many"
region = "ref:many?"
line = "string"
city = "string"
mail_code = "string?"
extra = "whatever"
descriptor_ = {terms = ["line", "city", "region", "country", "mail_code"], description = "$line $city(, $region), $country(, $mail_code)"}
unique_ = [["id"]]

[product]
id = "uuid"
price = "decimal(6,2)"
`)

		config = Load(data)
	)
	// Validate before asserting equality, to ensure it does not alter the configuration
	Validate(config)

	assert.Equal(
		t,
		Configuration{
			Database: Database{
				Name:            "mydb",
				Description:     "my great database",
				Locale:          "en_CA",
				AccentSensitive: false,
				CaseSensitive:   false,
				Schemas:         []string{"app1", "app2"},
				Vendors:         []Vendor{Postgres},
				VendorTypes: []VendorType{
					{
						Name: "sowhat",
						VendorColDefs: map[Vendor]string{
							Postgres: "psqlType",
						},
					},
					{
						Name: "whatever",
						VendorColDefs: map[Vendor]string{
							Postgres: "psqlType",
						},
					},
				},
			},
			UserDefinedTypes: []UserDefinedType{
				{
					Name: "address",
					Fields: []Field{
						{
							Name: "city",
							Type: String,
						},
						{
							Name: "country",
							Type: RefManyToOne,
						},
						{
							Name:     "extra",
							Type:     VendorTypeRef,
							TypeName: "whatever",
						},
						{
							Name: "id",
							Type: UUID,
						},
						{
							Name: "line",
							Type: String,
						},
						{
							Name:     "mail_code",
							Type:     String,
							Nullable: true,
						},
						{
							Name:     "region",
							Type:     RefManyToOne,
							Nullable: true,
						},
					},
					Descriptor: &Descriptor{
						Terms:       []string{"line", "city", "region", "country", "mail_code"},
						Description: "$line $city(, $region), $country(, $mail_code)",
					},
					UniqueKeys: [][]string{{"id"}},
				},
				{
					Name: "country",
					Fields: []Field{
						{
							Name: "code",
							Type: String,
						},
						{
							Name: "has_regions",
							Type: Bool,
						},
						{
							Name: "id",
							Type: UUID,
						},
						{
							Name: "name",
							Type: String,
						},
					},
					Descriptor: &Descriptor{
						Terms:       []string{"name", "code"},
						Description: "$name",
					},
					UniqueKeys: [][]string{{"name"}, {"code"}},
				},
				{
					Name: "product",
					Fields: []Field{
						{
							Name: "id",
							Type: UUID,
						},
						{
							Name:      "price",
							Type:      Decimal,
							Precision: 6,
							Scale:     2,
						},
					},
				},
				{
					Name: "region",
					Fields: []Field{
						{
							Name: "code",
							Type: String,
						},
						{
							Name: "country",
							Type: RefManyToOne,
						},
						{
							Name: "id",
							Type: UUID,
						},
						{
							Name: "name",
							Type: String,
						},
					},
					Descriptor: &Descriptor{
						Terms:       []string{"name", "code"},
						Description: "$code",
					},
					UniqueKeys: [][]string{{"country", "name"}, {"country", "code"}},
				},
			},
		},
		config,
	)
}

func TestLoadDefaults_(t *testing.T) {
	// Load that specifies only values with no default
	var (
		data = strings.NewReader(`
[address]
id = "uuid"
country = "ref:many"
region = "ref:many?"
line = "string"
city = "string"
mail_code = "string?"
descriptor_ = {terms = ["line", "city", "region", "country", "mail_code"], description = "$line $city(, $region), $country(, $mail_code)"}
`)

		config = Load(data)
	)

	assert.Equal(
		t,
		Configuration{
			Database: defaultConfiguration.Database,
			UserDefinedTypes: []UserDefinedType{
				{
					Name: "address",
					Fields: []Field{
						{
							Name: "city",
							Type: String,
						},
						{
							Name: "country",
							Type: RefManyToOne,
						},
						{
							Name: "id",
							Type: UUID,
						},
						{
							Name: "line",
							Type: String,
						},
						{
							Name:     "mail_code",
							Type:     String,
							Nullable: true,
						},
						{
							Name:     "region",
							Type:     RefManyToOne,
							Nullable: true,
						},
					},
					Descriptor: &Descriptor{
						Terms:       []string{"line", "city", "region", "country", "mail_code"},
						Description: "$line $city(, $region), $country(, $mail_code)",
					},
				},
			},
		},
		config,
	)
}

func TestLoadErrors_(t *testing.T) {
	var (
		failed   bool
		data     *strings.Reader
		loadData = func() { Validate(Load(data)); assert.Fail(t, "Must die") }
	)

	// ==== errColumnNameInvalidMsg

	// Empty column name
	data = strings.NewReader(`
[address]
"" = "uuid"
`)

	funcs.TryTo(
		loadData,
		func(e any) {
			failed = true
			assert.Equal(t, fmt.Errorf("address: a column name cannot be an empty string or end in an underscore"), e)
		},
	)
	assert.True(t, failed)

	// Column name is underscore
	data = strings.NewReader(`
[address]
_ = "uuid"
`)

	funcs.TryTo(
		loadData,
		func(e any) {
			failed = true
			assert.Equal(t, fmt.Errorf("address: a column name cannot be an empty string or end in an underscore"), e)
		},
	)
	assert.True(t, failed)

	// Column ends in underscore
	data = strings.NewReader(`
[address]
foo_ = "uuid"
`)

	funcs.TryTo(
		loadData,
		func(e any) {
			failed = true
			assert.Equal(t, fmt.Errorf("address: a column name cannot be an empty string or end in an underscore"), e)
		},
	)
	assert.True(t, failed)

	// ==== errDescriptorMustHaveTermsAndDescriptionMsg

	data = strings.NewReader(`
  [address]
  descriptor_ = {terms = [], description = ""}
  `)

	funcs.TryTo(
		loadData,
		func(e any) {
			failed = true
			assert.Equal(t, fmt.Errorf("address.descriptor_: terms array must have at least one string, and description must be a non-empty string"), e)
		},
	)
	assert.True(t, failed)

	// ==== errDuplicateUniqueKeyMsg

	data = strings.NewReader(`
  [address]
  unique_ = [["a", "b"], ["c"], ["b", "a"]]
  `)

	funcs.TryTo(
		loadData,
		func(e any) {
			failed = true
			assert.Equal(t, errors.Join(fmt.Errorf("address.unique_ has a duplicate key [a b] at indexes [0 2] (the order of columns is not significant)")), e)
		},
	)
	assert.True(t, failed)

	// ==== errEmptyUniqueMsg

	data = strings.NewReader(`
  [address]
  unique_ = []
  `)

	funcs.TryTo(
		loadData,
		func(e any) {
			failed = true
			assert.Equal(t, errors.Join(fmt.Errorf("address.unique_ must have at least one key")), e)
		},
	)
	assert.True(t, failed)

	// ==== errEmptyUniqueKeyMsg

	data = strings.NewReader(`
  [address]
  unique_ = [["a", "b"], [], ["c", "a"]]
  `)

	funcs.TryTo(
		loadData,
		func(e any) {
			failed = true
			assert.Equal(t, errors.Join(fmt.Errorf("address.unique_[1] has an empty key")), e)
		},
	)
	assert.True(t, failed)

	// ==== errFieldOfUndefinedVendorTypeMsg

	data = strings.NewReader(`
  [address]
  foo = "bar"
  `)

	funcs.TryTo(
		loadData,
		func(e any) {
			failed = true
			assert.Equal(t, errors.Join(fmt.Errorf("address.foo refers to undefined vendor type bar")), e)
		},
	)
	assert.True(t, failed)

	// ==== errNoSuchVendorMsg

	data = strings.NewReader(`
  [database_]
  vendors = ["noSuchVendor"]
  `)

	funcs.TryTo(
		loadData,
		func(e any) {
			failed = true
			assert.Equal(t, fmt.Errorf(`"noSuchVendor" is not a recognized database vendor name`), e)
		},
	)
	assert.True(t, failed)

	data = strings.NewReader(`
  [database_]
  vendor_types = {whatever = {noSuchVendor = "noSuchType"}}
  `)

	funcs.TryTo(
		loadData,
		func(e any) {
			failed = true
			assert.Equal(t, fmt.Errorf(`"noSuchVendor" is not a recognized database vendor name`), e)
		},
	)
	assert.True(t, failed)

	// ==== errRefToUndefinedTypeMsg

	data = strings.NewReader(`
  [address]
  region = "ref:one"
  `)

	funcs.TryTo(
		loadData,
		func(e any) {
			failed = true
			assert.Equal(t, errors.Join(fmt.Errorf("address.region is a reference field, but there is no User Defined Type by that name")), e)
		},
	)
	assert.True(t, failed)

	// ==== errUDTNameInvalidMsg

	data = strings.NewReader(`
  [""]
  `)

	funcs.TryTo(
		loadData,
		func(e any) {
			failed = true
			assert.Equal(t, fmt.Errorf(`"": user defined type names cannot be empty or end with an underscore`), e)
		},
	)
	assert.True(t, failed)

	data = strings.NewReader(`
  [foo_]
  `)

	funcs.TryTo(
		loadData,
		func(e any) {
			failed = true
			assert.Equal(t, fmt.Errorf(`"foo_": user defined type names cannot be empty or end with an underscore`), e)
		},
	)
	assert.True(t, failed)

	// ==== errUnrecognizedDatabaseKeyMsg

	data = strings.NewReader(`
  [database_]
  foo = "bar"
  `)

	funcs.TryTo(
		loadData,
		func(e any) {
			failed = true
			assert.Equal(t, fmt.Errorf("foo is not a valid database_ configuration key"), e)
		},
	)
	assert.True(t, failed)
}
