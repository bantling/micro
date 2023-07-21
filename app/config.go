package app

// SPDX-License-Identifier: Apache-2.0

import (
	"io"
  "regexp"

	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/funcs"
	"github.com/mitchellh/mapstructure"
	"github.com/pelletier/go-toml/v2"
)

var (
  errDescriptorMustHaveTermsAndDescriptionMsg = "%s: descriptor_ must contain terms array of at least one string and non-empty description string"
  errUniqueMustHaveAtLeastOneColumnMsg = "%s: unique_ must contain at least key of at least one column"
  errColumnTypeNotRecognizedMsg = "%s: the column %s is not a valid column name, or the type ias not a recognized type"
)

// Vendor is a database vendor
type Vendor uint

// Vendor constants
const (
  Postgres Vendor = iota
)

// CustomType is any type that is not directly supported, that can have different names for different vendors
type CustomType struct {
  Name string
  VendorTypes map[Vendor]string
}

// Database contains the database portion of configuration
type Database struct {
	Name            string
	Description     string
	Locale          string
	Encoding        string
	AccentSensitive bool `mapstructure:"accent_sensitive"`
	CaseSensitive   bool `mapstructure:"case_sensitive"`
  Schemas         []string
  VendorTypes     []CustomType
}

// FieldType is an enum of known field types
type FieldType uint

// FieldType constants
const (
  Bool FieldType = iota
  Date
  Decimal
  Float32
  Float64
  Int32
  Int64
  JSON
  RefOne
  RefManyToOne
  RefManyToMany
  Row
  String
  TableRow
  Time
  Timestamp
  UUID
  Custom
)

// String to FieldType mapping for all FieldType except Ref* and Custom
var (
  fieldStrings = map[string]FieldType {
    "bool": Bool,
    "date": Date,
    "float32": Float32,
    "float64": Float64,
    "int32": Int32,
    "int64": Int64,
    "json": JSON,
    "row": Row,
    "string": String,
    "table_row": TableRow,
    "time": Time,
    "timestamp": Timestamp,
    "uuid": UUID,
  },

  decimalPrecisionScaleRegex = regexp.MustCompile(`decimal[(]([0-9]+)(?:, *([0-9]+))?[)]`)
  refRegex = regexp.MustCompile(`ref:(one|manyToMany|many)`) // have to put prefix many after manyToMany
  stringLengthRegex = regexp.MustCompile(`string[(]([0-9]+)[)]`)
)

// stringToTypeDef converts a field string to a field type def, as follows:
//
// decimal(precision) -> Decimal, precision, 0
// decimal(precision, scale) -> Decimal, precision, scale
// ref:one -> RefOne
// ref:many -> RefManyToOne
// ref:manyToMany -> RefManyToMany
// string -> String, 0, 0
// string(limit) -> String, limit, 0
//
// any type can be followed by ? to make it nullable, else it is non-nullable
func stringToTypeDef(str string) typ FieldType, length, scale int, nullable bool {
  // Handle common cases, including string - but not string(limit)
  // Strip ? suffix if it exists
  nullable = str[len(str)-1] == '?'
  if nullable {
    str = str[:len(str)-1]
  }

  switch {
  case ft, isa := fieldStrings[str]; isa
    typ = ft

  // Try decimal(precision, scale)
  case match := decimalPrecisionScaleRegex.FindStringSubmatch(str)); match != nil
    typ = Decimal
    conv.To(match[1], &length)
    conv.To(match[2], &scale) // Ignore error if string is empty, leaving scale at 0

  // Try ref:one, ref:many, ref:manyToMany
  case match := refRegex.FindStringSubmatch(str)); match != nil
    switch match[1] {
    case "one":
      typ = RefOne
    case "many":
      typ = RefManyToOne
    default:
      typ = RefManyToMany
    }

  // Try string(limit)
  case match := stringLengthRegex.FindStringSubmatch(str)); match != nil
    typ = String
    conv.To(match[1], &length)
  }

  return
}

// Field describes a single database field
type Field struct {
  Name string
  Type FieldType
  TypeName string // The other type name for ref*, custom type name, else empty
  Precision int // precision of decimal
  Scale int // scale of decimal
  Length // length of string
  Nullable bool // true if nullable
}

// Descriptor describes the optional search terms and description for each object
type Descriptor struct {
  Terms []string // Individual terms that can be used to search for each object
  Description string // The description of each object
}

// UserDefinedType contains the details of a single user defined type
type UserDefinedType struct {
  Name string // the type name
  Fields []Field // the fields of the type
  Descriptor *Descriptor // the descriptor, if any
  UniqueKeys [][]string // the set of unique keys, if any
}

// Configuration contains a combination of knowable stuff like the general database config, and unknowable stuff like
// whatever objects are stored in the database.
type Configuration struct {
	Database    Database
	UserDefinedTypes []UserDefinedType
}

var (
	// defaultConfiguration is the default Configuration, where default values are not necessarily zero values.
	defaultConfiguration = Configuration{
		Database: Database{
			Name:            "myapp",
			Locale:          "en_US",
			AccentSensitive: true,
			CaseSensitive:   true,
		},
		UserDefinedTypes: []UserDefinedType{},
	}
)

// Load a TOML file into a Configuration.
// The approach used is to simply decode into a map[string]any, and look for knowable stuff like the database config
// (which is necessarily a sub map[string]any), which gets converted into Configuration.Database via mapstruct.
// All unrecognized top level keys are assumed to be user defined types.
func Load(src io.Reader) Configuration {
	var (
		config      = defaultConfiguration
		configMap   = map[string]any{}
		tomlDecoder = toml.NewDecoder(src)
	)

	funcs.Must(tomlDecoder.Decode(&configMap))

	// Iterate all top level keys of configMap:
	// - Recognize keys are decoded into appropriate field
	// - Remaining keys are mapped in the UserDefinedTypes field
	for k, v := range configMap {
		switch k {
		case "database_":
      // Decode into config.Database automatically
			{
				var (
					msdc      = mapstructure.DecoderConfig{ErrorUnused: true, Result: &config.Database}
					msDecoder = funcs.MustValue(mapstructure.NewDecoder(&msdc))
				)
				funcs.Must(msDecoder.Decode(v))
			}

		default:
      // Decode into config.UserDefinedTypes manually
			{
        var (
          udf UserDefinedType
          data = v.(map[string]any)
        )

        // Top level table name is type name
        udf.Name = k

        // Iterate keys of table
        for fk, fv := range data {
          switch fk {
            // descriptor_ -> *Descriptor
            case "descriptor_": {
              // Grab terms and description
              var (
                fdata = fv.(map[string]any)
                terms, haveTerms = fdata["terms"]
                desc, haveDesc = fdata["description"]
                err = fmt.Errorf(errDescriptorMustHaveTermsAndDescriptionMsg, udf.Name)
              )

              // Must have terms and description
              if !(haveTerms && haveDesc) {
                panic(err)
              }

              // Terms must be a []string and Description must be a string
              var (
                slcTerms, isSlc = terms.([]string)
                strDesc, isStr = desc.(string)
              )

              // Terms and Description must not be empty
              if !(isSlc && isStr && (len(slcTerms) > 0) && (len(strDesc) > 0)) {
                panic(err)
              }

              udf.Descriptor = &Descriptor{
                Terms = slcTerms,
                Description = strDesc,
              }
            }

            case "unique_": {
              // Grab set of unique keys
              var (
                uks, isa = fv.([][]string)
                err = fmt.Errorf(errUniqueMustHaveAtLeastOneColumnMsg, udf.Name)
              )

              // Must have at least one unique key
              if !(isa && (len(uks) > 0)) {
                panic(err)
              }

              // Each unique key must contain at least one column name, all column names must be non-empty
              for _, uk := range uks {
                if len(uk) == 0 {
                  panic(err)
                }

                for _, column := range uk {
                  if len(column) == 0 {
                    panic(err)
                  }
                }
              }

              udf.UniqueKeys = uks
            }

            default: {
              // Must be a column, the value must be a string of a recognized type
              var (
                colName, isa := fv.(string)
                err = fmt.Errorf(errColumnTypeNotRecognizedMsg, udf.Name, colName)
              )

              if !(isa && (len(colName) > 0) && (colName[len(colName)-1] != '_')) {
                panic(err)
              }


            }
          }
        }

  			config.UserDefinedTypes = append(config.UserDefinedTypes, udf)
			}
		}
	}

	return config
}
