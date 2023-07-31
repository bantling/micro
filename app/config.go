package app

// SPDX-License-Identifier: Apache-2.0

import (
  "fmt"
	"io"
  "regexp"

	"github.com/bantling/micro/conv"
	"github.com/bantling/micro/funcs"
	"github.com/pelletier/go-toml/v2"
)

var (
  errNoSuchVendorMsg = "%q is not a recognized database vendor name"
  errDescriptorMustHaveTermsAndDescriptionMsg = "%s: descriptor_ must contain terms array of at least one string and non-empty description string"
  errUniqueMustHaveAtLeastOneColumnMsg = "%s: unique_ must contain at least one key, and each key must have at least one column"
  errColumnTypeNotRecognizedMsg = "%s: the column %s is not a valid column name, or the type is not a recognized type"
  errUnrecognizedDatabaseKeyMsg = "%s is not a valid database_ configuration key"
  errUDTCannotEndWithUnderscoreMsg = "%s: user defined type names cannot end with an underscore"
)

// Vendor is a database vendor
type Vendor uint

// Vendor constants
const (
  Postgres Vendor = iota
)

// Vendor string names
var (
  vendorStrings = map[string]Vendor{
    "postgres": Postgres,
  }
)

// CustomType is a custom type, a string name associated with a string column definition, one per vendior
type CustomType struct {
  Name string
  VendorColDefs map[Vendor]string
}

// Database contains the database portion of configuration
type Database struct {
	Name            string
	Description     string
	Locale          string
	AccentSensitive bool
	CaseSensitive   bool
  Schemas         []string
  Vendors         []Vendor
  CustomTypes     []CustomType
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
  Interval
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
    "interval": Interval,
    "json": JSON,
    "ref:one": RefOne,
    "ref:many": RefManyToOne,
    "ref:manyToMany": RefManyToMany,
    "row": Row,
    "string": String,
    "table_row": TableRow,
    "time": Time,
    "timestamp": Timestamp,
    "uuid": UUID,
  }

  decimalPrecisionScaleRegex = regexp.MustCompile(`decimal[(]([0-9]+)(?:, *([0-9]+))?[)]`)
  //refRegex = regexp.MustCompile(`ref:(one|manyToMany|many)`) // have to put prefix many after manyToMany
  stringLengthRegex = regexp.MustCompile(`string[(]([0-9]+)[)]`)
)

// stringToTypeDef converts a field string to a field type def, as follows:
//
// decimal(precision) -> Decimal, precision, 0
// decimal(precision, scale) -> Decimal, precision, scale
// string -> String, 0, 0
// string(limit) -> String, limit, 0
//
// any type can be followed by ? to make it nullable, else it is non-nullable
func stringToTypeDef(str string) (typ FieldType, length, scale int, nullable bool) {
  // Handle common cases, including string - but not string(limit)
  // Strip ? suffix if it exists, setting nullable to true/false
  nullable = str[len(str)-1] == '?'
  if nullable {
    str = str[:len(str)-1]
  }

  // Try fieldStrings map
  if ft, isa := fieldStrings[str]; isa {
    typ = ft

    // Try decimal(precision, scale)
  } else if match := decimalPrecisionScaleRegex.FindStringSubmatch(str); match != nil {
    typ = Decimal
    conv.To(match[1], &length)
    conv.To(match[2], &scale) // Ignore error if string is empty, leaving scale at 0

  // // Try ref:one, ref:many, ref:manyToMany
  // } else if match := refRegex.FindStringSubmatch(str); match != nil {
  //   switch match[1] {
  //   case "one":
  //     typ = RefOne
  //   case "many":
  //     typ = RefManyToOne
  //   default:
  //     typ = RefManyToMany
  //   }

    // Try string(limit)
  } else if match := stringLengthRegex.FindStringSubmatch(str); match != nil {
    typ = String
    conv.To(match[1], &length)
  } else {
    // Must be a custom or invalid type, treat it as custom
    typ = Custom
  }

  return
}

// Field describes a single database field
type Field struct {
  Name string
  Type FieldType
  TypeName string // custom type name
  Precision int // precision of decimal
  Scale int // scale of decimal
  Length int // length of string
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
      Schemas:         []string{},
      Vendors: []Vendor{Postgres},
      CustomTypes: []CustomType{},
		},
		UserDefinedTypes: []UserDefinedType{},
	}
)

// Load a TOML file into a Configuration.
// The approach used is to simply decode into a map[string]any, and look for knowable stuff like the database config
// (which is necessarily a sub map[string]any), and manually convert it into Configuration.Database.
// All unrecognized top level keys are manually converted into Configuration.Database.UserDefinedTypes.
func Load(src io.Reader) Configuration {
  config := defaultConfiguration

	var (
		configMap   = map[string]any{}
		tomlDecoder = toml.NewDecoder(src)
	)

	funcs.Must(tomlDecoder.Decode(&configMap))

	// Iterate all top level keys of configMap:
	// - Recognized keys are decoded into appropriate field
	// - Remaining keys are mapped in the UserDefinedTypes field
	for k, v := range configMap {
		switch k {
		case "database_":
      // Decode into config.Database
			{
        database := funcs.MustAssertType[map[string]any](k, v)

        for fk, fv := range database {
          databasePath := k + "." + fk

          switch fk {
          case "name":
            config.Database.Name = funcs.MustAssertType[string](databasePath, fv)

          case "description":
            config.Database.Description = funcs.MustAssertType[string](databasePath, fv)

          case "locale":
            config.Database.Locale = funcs.MustAssertType[string](databasePath, fv)

          case "accent_sensitive":
            config.Database.AccentSensitive = funcs.MustAssertType[bool](databasePath, fv)

          case "case_sensitive":
            config.Database.CaseSensitive = funcs.MustAssertType[bool](databasePath, fv)

          case "schemas":
            config.Database.Schemas = funcs.MustAssertType[[]string](databasePath, fv)

          case "vendors": {
            var (
              vendors = funcs.MustAssertType[[]string](databasePath, fv)
              uniqueVendors = map[Vendor]int{}
            )

            // Collect unique vendor names, don't care if same vendor specified multiple times
            for _, vendorName := range vendors {
              if vendor, hasIt := vendorStrings[vendorName]; hasIt {
                uniqueVendors[vendor] = 0
              } else {
                panic(fmt.Errorf(errNoSuchVendorMsg, vendorName))
              }
            }

            config.Database.Vendors = funcs.MapKeysToSlice(uniqueVendors)
          }

        case "custom_types": {
          var (
            customTypeDefs = funcs.MustAssertType[map[string]any](databasePath, fv)
            customTypes = []CustomType{}
          )

          for customTypeName, vendorColDefsVal := range customTypeDefs {
            var (
              ctPath = databasePath + "." + customTypeName
              vendorColDefsVal = funcs.MustAssertType[map[string]string](ctPath, vendorColDefsVal)
              vendorColDefs = map[Vendor]string{}
            )

            for vendorName, colDef := range vendorColDefsVal {
              if vendor, hasIt := vendorStrings[vendorName]; hasIt {
                vendorColDefs[vendor] = colDef
              } else {
                panic(fmt.Errorf(errNoSuchVendorMsg, vendorName))
              }
            }

            customTypes = append(
              customTypes,
              CustomType{
                Name: customTypeName,
                VendorColDefs: vendorColDefs,
              },
            )
          }
        }

          default:
            panic(fmt.Errorf(errUnrecognizedDatabaseKeyMsg, fk))
          }
        }
			}

		default:
      // Decode into config.UserDefinedTypes manually
			{
        if k[len(k)-1] == '_' {
          panic(fmt.Errorf(errUDTCannotEndWithUnderscoreMsg, k))
        }

        var (
          udf UserDefinedType
          data = funcs.MustAssertType[map[string]any](k, v)
        )

        // Top level table name is type name
        udf.Name = k

        // Iterate keys of table
        for fk, fv := range data {
          udtPath := k + "." + fk

          switch fk {
            // descriptor_ -> *Descriptor
            case "descriptor_": {
              // Grab terms and description
              var (
                fdata = funcs.MustAssertType[map[string]any](udtPath, fv)
                terms = funcs.MustAssertType[[]string](udtPath + ".terms", fdata["terms"])
                desc = funcs.MustAssertType[string](udtPath + ".description", fdata["description"])
              )

              // Must have terms and description
              if (len(terms) == 0) || (len(desc) == 0) {
                panic(fmt.Errorf(errDescriptorMustHaveTermsAndDescriptionMsg, udtPath))
              }

              udf.Descriptor = &Descriptor{
                Terms: terms,
                Description: desc,
              }
            }

            case "unique_": {
              // Grab set of unique keys
              var (
                uks = funcs.MustAssertType[[][]string](udtPath, fv)
                err = fmt.Errorf(errUniqueMustHaveAtLeastOneColumnMsg, udf.Name)
              )

              // Must have at least one unique key
              if len(uks) == 0 {
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

              if str, isa := fv.(string); isa {
                var (
                  typ, length, scale, nullable = stringToTypeDef(str)
                )

                // If the type name is empty or ends in an underscore, it is invalid
                if ((fk == "") || (fk[len(fk)-1] != '_')) {
                  panic(fmt.Errorf(errColumnTypeNotRecognizedMsg, udf.Name, fk))
                } else {
                    // Assume it is a valid definition
                    fld := Field{
                      Name: fk,
                      Type: typ,
                      TypeName: "",
                      Nullable: nullable,
                    }

                    // Copy length and scale values, if relevant
                    switch typ {
                    case String:
                      fld.Length = length
                    case Decimal:
                      fld.Precision = length
                      fld.Scale = scale
                    }

                    udf.Fields = append(udf.Fields, fld)
                }
              } else {
                // Not a string, reject it
                panic(fmt.Errorf(errColumnTypeNotRecognizedMsg, udf.Name, fk))
              }
            }
          }
        }

		    config.UserDefinedTypes = append(config.UserDefinedTypes, udf)
			}
		}
	}

  // Validate that the udfs make sense

  // 1. If a field is a Ref* type, the field name is a type name
  // 2. If a field is a Custom type, then the field name is defined in Database.VendorTypes

	return config
}
