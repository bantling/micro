package app

// SPDX-License-Identifier: Apache-2.0

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAllValues(t *testing.T) {
	// Load that specifies all values
	var (
		data = strings.NewReader(`
[database_]
name = "mydb"
description = "my great database"
locale = "en_CA"
accent_sensitive = false
case_sensitive = false

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
			Database: Database{
        Name:"mydb",
				Description:     "my great database",
				Locale:          "en_CA",
				AccentSensitive: false,
				CaseSensitive:   false,
			},
			UserDefinedTypes: []UserDefinedType{
        UserDefinedType {
          Name: "address",
          Fields: []Field {
            Field {
              Name: "id",
              Type: UUID,
            },
            Field {
              Name: "country",
              Type: RefMany,
            },
            Field {
              Name: "region",
              Type: RefMany,
              Nullable: true,
            },
            Field {
              Name: "line",
              Type: String,
            },
            Field {
              Name: "city",
              Type: String,
            },
            Field {
              Name: "mail_code",
              Type: String,
              Nullable: true,
            },
          },
          Descriptor: Descriptor {
            Terms: []string{"line", "city", "region", "country", "mail_code"},
            Description: "$line $city(, $region), $country(, $mail_code)",
          },
        },
			},
		},
		config,
	)
}

// func TestLoadDefaults(t *testing.T) {
// 	// Load that specifies only values with no default
// 	var (
// 		data = strings.NewReader(`
// [address]
// id = "uuid"
// country = "ref:many"
// region = "ref:many?"
// line = "string"
// city = "string"
// mail_code = "string?"
// ext_descriptor = {terms = ["line", "city", "region", "country", "mail_code"], description = "$line $city(, $region), $country(, $mail_code)"}
// `)
//
// 		config = Load(data)
// 	)
//
// 	assert.Equal(
// 		t,
// 		Configuration{
// 			Database: defaultConfiguration.Database,
// 			UserDefinedTypes: map[string]any{
// 				"address": map[string]any{
// 					"id":        "uuid",
// 					"country":   "ref:many",
// 					"region":    "ref:many?",
// 					"line":      "string",
// 					"city":      "string",
// 					"mail_code": "string?",
// 					"ext_descriptor": map[string]any{
// 						"terms":       []any{"line", "city", "region", "country", "mail_code"},
// 						"description": "$line $city(, $region), $country(, $mail_code)",
// 					},
// 				},
// 			},
// 		},
// 		config,
// 	)
// }
