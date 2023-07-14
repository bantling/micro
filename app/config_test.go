package app

// SPDX-License-Identifier: Apache-2.0

import (
  "strings"
  "testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
  var (
  data = strings.NewReader(`
[database]
name = "mydb"
description = "my great database"
locale = "en_US"
encoding = "UTF8"
accent_sensitive = false
case_sensitive = false

[address]
id = "uuid"
country = "ref:many"
region = "ref:many?"
line = "string"
city = "string"
mail_code = "string?"
ext_descriptor = {terms = ["line", "city", "region", "country", "mail_code"], description = "$line $city(, $region), $country(, $mail_code)"}
`)

  config = Load(data)
)

  assert.Equal(
    t,
    Configuration{
      Database: Database {
        Name: "mydb",
        Description: "my great database",
        Locale: "en_US",
        Encoding: "UTF8",
        AccentSensitive: false,
        CaseSensitive: false,
      },
      UserDefined: map[string]any{
        "address": map[string]any {
          "id": "uuid",
          "country": "ref:many",
          "region": "ref:many?",
          "line": "string",
          "city": "string",
          "mail_code": "string?",
          "ext_descriptor": map[string]any {
            "terms": []any{"line", "city", "region", "country", "mail_code"},
            "description": "$line $city(, $region), $country(, $mail_code)",
          },
        },
      },
    },
    config,
  )
}
