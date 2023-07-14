package app

// SPDX-License-Identifier: Apache-2.0

import (
  "io"

  "github.com/bantling/micro/funcs"
  "github.com/mitchellh/mapstructure"
  "github.com/pelletier/go-toml/v2"
)

// Database contains the database portion of configuration
type Database struct {
  Name string
  Description string
  Locale string
  Encoding string
  AccentSensitive bool `mapstructure:"accent_sensitive"`
  CaseSensitive bool `mapstructure:"case_sensitive"`
}

// Configuration contains a combination of knowable stuff like the general database config, and unknowable stuff like
// whatever objects are stored in the database.
type Configuration struct {
  Database Database
  UserDefined map[string]any
}

var (
  // defaultConfiguration is the default Configuration, where default values are not necessarily zero values.
  defaultConfiguration = Configuration {
    Database: Database {
      Name: "myapp",
      Locale: "en_US",
      Encoding: "UTF8",
      AccentSensitive: true,
      CaseSensitive: true,
    },
    UserDefined: map[string]any{},
  }
)

// Load a TOML file into a Configuration.
// The approach used is to simply decode into a map[string]any, and look for knowable stuff like the database config (which
// are necessarily a sub map[string]any), which gets converted into Configuration.Database via mapstruct.
// All unrecognized top level keys are assumed to be user defined types.
func Load(src io.Reader) Configuration {
  var (
    config = defaultConfiguration
    configMap = map[string]any{}
    tomlDecoder = toml.NewDecoder(src)
  )

  funcs.Must(tomlDecoder.Decode(&configMap))

  // Iterate all top level keys of configMap:
  // - Recognize keys are decoded into appropriate field
  // - Remaining keys are mapped in the UserDefined field
  for k, v := range configMap {
    switch k {
    case "database": {
      var (
        msdc = mapstructure.DecoderConfig{ErrorUnused: true, Result: &config.Database}
        msDecoder = funcs.MustValue(mapstructure.NewDecoder(&msdc))
      )
      funcs.Must(msDecoder.Decode(v))
    }

    default:
      config.UserDefined[k] = v
    }
  }

  return config
}
