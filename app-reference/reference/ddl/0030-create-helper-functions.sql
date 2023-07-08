\c mydb

-- Determine the difference between the desired schema provided and the actual schema of the provided table name.
-- Update the table schema to conform to the desired schema, as follows:
--
-- - Add any columns in desired schema that are not in actual schema
-- - Drop any columns in actual schema that are not in desired schema
-- - Alter type of any columns in both schemas where the types differ
-- - Alter nullability of any columns in both schemas where the nullability differs
--
-- Format of JSON is [{"column_name": {"type": "type_name", "nullable": boolean}}, ...]
-- Supported type names are:
-- - VARCHAR, VARCHAR(LENGTH)
-- - NUMERIC, NUMERIC(PRECISION), NUMERIC(PRECISION, SCALE)
-- - TIMESTAMP WITHOUT TIME ZONE
-- - DATE
-- - TIME WITHOUT TIME ZONE
-- - INTERVAL
-- - BOOLEAN
-- - INTEGER
-- - BIGINT
-- - REAL
-- - DOUBLE_PRECISION
-- - JSONB
-- - OID
--
-- Notes:
-- - Use type names exactly as specified (case does not matter), as that is how they are reported by information_schema.columns
-- - EG, do not use TEXT, DECIMAL, BOOL, or INT
-- - For a NUMERIC with SCALE = 0, do not use NUMERIC(PRECISION, 0), instead use NUMERIC(PRECISION)
-- - It is assumed the table will have a PRIMARY KEY of type INTEGER named REL_ID, do not provide it
CREATE OR REPLACE PROCEDURE myapp.update_table_schema(
  table_name VARCHAR
 ,desired_schema JSON
) AS $$
DECLARE
  rec RECORD;
BEGIN
  FOR rec IN SELECT
END
$$ LANGUAGE PLPGSQL;
