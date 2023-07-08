// Package ddl compares the database schema (using information_schema tables) against a toml file and makes changes
//
// SPDX-License-Identifier: Apache-2.0
//
// Example toml file:
//
// [database]
// name = "mydb"
// encoding = "en_US.UTF8" // default is en_US.UTF8
// case_sensitive = true   // default is true
// // when false: postgres create database locale=en_US.UTF8 lc-collate=C
//                use ilike and ~* for case insensitive matching (~* for ci regex)
//                create extension if not exists pg_trgm
//                create index table_column_idx on <table> using gin (<column> gin_trgm_ops) to speed up ilike and ~*
// schemas = ["myapp"] // default is no schema name. If multiple schemas provided, each type has to state schema.
// vendor_types = [{"custom_type_name" = {"vendor_name" = "vendor_type"}+} // define custom type names for vendor types
//
// example address type
// if one schema s specified in [database], then address would actually be s.address
// For multiple schemas s1 and s2, we'd have to specify [s1.address] or [s2.address], else an error occurs.
// table names cannot begin with ext_, that prefix is reserved for baked in functionality
//
// [address]
// id = "uuid"
// country = "ref:many"
// region = "ref:many"
// line = "string"
// city = "string"
// mail_code = "string"
// descriptor = {terms = ["line", "city", "region", "country", "mail_code"], format = "$line $city(, $region), $country(, $mail_code)}"
//
// [country]
// id = "uuid"
// name = "string"
// code = "string"
// unique = [["name"],["code"]]
// has_regions = bool
//
// [region]
// id = "uuid"
// country = "ref:many"
// name = "string"
// code = "string"
// unique = [["country","name"],["code"]]
//
// The following types are supported for each column:
// - bool
// - date
// - decimal(5), decimal(5,2)
// - float32
// - float64
// - int32
// - int64
// - json
// - ref (see below)
// - string
// - string(limit)
// - table (see below)
// - time
// - timestamp
// - uuid
// - custom type names defined in [database].vendor_types
//
// Each table has an implicit primary key column named rel_id that is an identity column. It's sole purpose is for
// relationships between tables. It is an error for a table to define a column named rel_id.
//
// Notes about types:
// - The type names provided are logical type names, they get translated to whatever the database calls them.
// - The string type with no limit is translated to some column type that allows for some large upper limit, with
//   reasonable efficiency.
//   Limiting the length of the string often results in pain and suffering down the line, and is generally not needed.
//
// The descriptor column defines the columns used for search terms, and the format of a text description.
// If at least one table has a descriptor, full text searching is enabled:
// - schema_name.ext_shared_search table has tbl_oid oid, tbl_rel_id integer, descriptor varchar, terms tsvector
// - all tables use statement level triggers to insert/update/delete entries in ext_shared_search table.
//
// The ref type supports 3 types of references (implicitly using the rel_id identity column):
// ref:many   - many rows in this table can refer to the same row of the target table (this table contains ref)
// ref:one    - one row in this table can refer to a given row in the target table (target table contains ref)
// ref:bridge - many rows in this table can refer to the same row of the target table, and vice-versa (bridge table)
// foreign keys are defined for ref:many and ref:one constraints
//
// The table type is a reference to another table, not a row of it. Useful for scenarios like graphing.
//
// The name unique is reserved for defining one or more unique constraints in the form of a two dimensional array.
// Each row of the array contains the names of one or more columns that must be a unique combination for each row.
// It is an error if the value of unique is not a two dimensional array of strings.
package ddl
