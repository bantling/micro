// Package ddl compares the database schema (using information_schema tables) against a toml file and makes changes
//
// SPDX-License-Identifier: Apache-2.0
//
// Example toml file:
//
// [database]
// name = "mydb"
// description = "my great database"
// locale = "en_US" // defaul;t is en_US
// encoding = "UTF8" // default is UTF8
// accent_sensitive = boolean // default is true
// case_sensitive = boolean   // default is true
//
// when accent_sensitive and/or case_sensitive is false:
// (see https://stackoverflow.com/questions/11005036/does-postgresql-support-accent-insensitive-collations)
//
//   postgres create database <name> locale=en_US.UTF8;
//   createdb dbname [description] -E UTF8 -l en_US
//
//   CREATE EXTENSION unaccent; -- accent insensitive
//   CREATE EXTENSION pg_trgm; -- case insensitive
//
//   -- accent insensitive
//   CREATE OR REPLACE FUNCTION public.immutable_unaccent(REGDICTIONARY, TEXT)
//   RETURNS TEXT
//   LANGUAGE C IMMUTABLE PARALLEL SAFE STRICT AS
//   '$libdir/unaccent', 'unaccent_dict';
//
//   -- accent insensitive
//   CREATE OR REPLACE FUNCTION public.f_unaccent(TEXT)
//   RETURNS TEXT
//   LANGUAGE SQL IMMUTABLE PARALLEL SAFE STRICT
//   BEGIN ATOMIC
//     SELECT public.immutable_unaccent(REGDICTIONARY 'public.unaccent', $1); // shd this be 'public.immutable_unaccent'?
//   END;
//
//   -- accent insensitive index
//   CREATE INDEX <table>_<column>_ai_idx ON <table>(public.f_unaccent(<column>));
//
//   -- case insensitive index
//   CREATE INDEX <table>_<column>_ci_idx ON <table> USING GIN(<column> GIN_TRGM_OPS);
//
//   -- accent and case insensitive index
//   CREATE INDEX <table>_<name>_ai_ci__idx ON <table> USING GIN (public.f_unaccent(<column>) GIN_TRGM_OPS);
//
//   -- accent insensitive query
//   SELECT * FROM <table> WHERE f_unaccent(<column>) = f_unaccent(<value>);
//
//   -- case insensitive query - can also do ci regex using ~*
//   SELECT * FROM <table> WHERE <column> ILIKE ('%' || <value> || '%');
//
//   -- accent and case insensitive query
//   SELECT * FROM <table> WHERE f_unaccent(name) ILIKE ('%' || f_unaccent(<value>) || '%');
//
// schemas = ["myapp"] // default is no schema name. If multiple schemas provided, each type has to state schema.
// vendor_types = [{"custom_type_name" = {"vendor_name" = "vendor_type"}+} // define custom type names for other types, per vendor
//
// example address type
// if one schema is specified in [database], then address would actually be s.address
// For two schemas s1 and s2, we'd have to specify [s1.address] or [s2.address], else an error occurs
// table names cannot begin with ext_, that prefix is reserved for baked in functionality
// column names cannot begin with ext_, except for specific cases related to baked in functionality (currently limited to ext_descriptor)
// column names cannot be ext_rel_id, that is reserved for the surrogate primary key of each table, which is used in foreign keys
// columns are not nullable by default, names ending in ? are nullable (the ? is not part of the name)
//
// example:

// [address]
// id = "uuid"
// country = "ref:many"
// region = "ref:many?"
// line = "string"
// city = "string"
// mail_code = "string?"
// ext_descriptor = {terms = ["line", "city", "region", "country", "mail_code"], description = "$line $city(, $region), $country(, $mail_code)}"
//
// Note: since region description is $code, the $region in the address description will be the region code.
//
// [country]
// id = "uuid"
// name = "string"
// code = "string"
// has_regions = bool
// unique = [["name"],["code"]]
// ext_descriptor = {terms = ["name", "code"], description = "$name"}
//
// [region]
// id = "uuid"
// country = "ref:many"
// name = "string"
// code = "string"
// unique = [["country","name"],["country", "code"]]
// ext_descriptor = {terms = ["name", "code"], description = "$code"}
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
// - ref
// - string, string(limit)
// - table_row
// - time
// - timestamp
// - uuid
// - custom type names defined in [database].vendor_types
//
// Each table has an implicit primary key column named ext_rel_id that is an identity column. It's sole purpose is for
// relationships between tables. It is an error for a table to define a column name starting with ext_.
//
// Notes about types:
//   - The type names provided are logical type names, they get translated to whatever the database calls them.
//   - The decimal scale defaults to 0, and must be in range [0, precision] for portability.
//   - The string type with no limit is translated to some column type that allows for some large upper limit, with
//     reasonable efficiency. Limiting the length of the string often results in pain and suffering down the line, and is generally not needed.
//
// The ext_descriptor column defines the columns used for search terms, and the format of a text description.
// If at least one table has an ext_descriptor, full text searching is enabled:
//   - <schema_name.>ext_shared_search table has tbl_oid oid, tbl_rel_id integer, descriptor varchar, terms <full text search type>
//   - all tables use statement level triggers to insert/update/delete entries in ext_shared_search table.
//   - format is recursive - eg, in example above, the format of a region is code, while the format of country is name.
//     So when address format refers to $region it means the code, and $country means the name.
//   - format has three special values:
//   - $column is replaced by the column value
//   - (...) is an optional string that is only displayed if at least one $ expression inside the () is non-empty.
//     any characters inside () except $column is literal text.
//
// The ref type supports 3 types of references (implicitly using the ext_rel_id identity column):
// ref:many   - many rows in this table can refer to the same row of the target table (this child table contains ref to parent table)
// ref:one    - one row in this table can refer to a given row in the target table (parent target table contains ref to this child table)
// ref:bridge - many rows in this table can refer to the same row of the target table, and vice-versa (bridge table)
// foreign keys are defined for ref constraints
// foreign key columns begin with ext_fk_
//
// The table_row type is a reference to a row of another table, such that both the table name and primary key are stored,
// referred to as <name>_table and <name>_ext_rel_id. Useful for scenarios like graphing.
//
// The name unique is reserved for defining one or more unique constraints in the form of a two dimensional array.
// Each row of the array contains the names of one or more columns that must be a unique combination for each row.
// It is an error if the value of unique is not a two dimensional array of strings.
// If multiple unique keys contain a common subset, then only one unique key is generated where subset comes first when possible.
// - Examples:
// - unique(foo, bar) and unique (bar, baz, foo) can be handled by one unique key (foo, bar, baz)
package ddl
