package app

// SPDX-License-Identifier: Apache-2.0

import (
	"strings"
	"testing"

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
vendor_types = {"whatever" = {"postgres" = "psqlType"}}

[address]
id = "uuid"
country = "ref:many"
region = "ref:many?"
line = "string"
city = "string"
mail_code = "string?"
descriptor_ = {terms = ["line", "city", "region", "country", "mail_code"], description = "$line $city(, $region), $country(, $mail_code)"}
unique_ = [["id"], ["country", "region", "line", "city", "mail_code"]]
`)

		config = Load(data)
	)

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
          UniqueKeys: [][]string{{"id"}, {"country", "region", "line", "city", "mail_code"}},
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
