\c mydb

--
-- Shared information
--
-- tbl_oid      object id for the source table
-- rel_id       PK of the source table row
-- descriptor   full text description of the source table row
-- search_terms terms to search by
--
-- PK: tbl_oid, rel_id
--

-- Table
CREATE TABLE IF NOT EXISTS myapp.shared_info(
  tbl_oid      OID      NOT NULL,
  rel_id       INTEGER  NOT NULL,
  description  VARCHAR  NOT NULL,
  search_terms TSVECTOR NOT NULL
);

-- PK
SELECT 'ALTER TABLE myapp.shared_info ADD CONSTRAINT shared_info_pk PRIMARY KEY (tbl_oid, rel_id)'
 WHERE NOT EXISTS (
   SELECT
     FROM information_schema.table_constraints
    WHERE table_schema    = 'myapp'
      AND table_name      = 'shared_info'
      AND constraint_name = 'shared_info_pk'
 )
\gexec

-- Index on search_terms field for full text searches
CREATE INDEX IF NOT EXISTS shared_info_ix_search_terms ON myapp.shared_info USING GIN(search_terms);

--
-- Customer
--
-- rel_id      PK
-- id          UUID
-- first_name  first name
-- middle_name middle name
-- last_name   last name
--
-- PK: rel_id
-- UK: id
--

-- Table
CREATE TABLE IF NOT EXISTS myapp.customer(
  rel_id INTEGER GENERATED ALWAYS AS IDENTITY
);

-- Add columns to / remove columns from customer table
SELECT CONCAT(
         'ALTER TABLE "myapp"."customer" ',
         CASE
           -- schema does not have desired column, add it
           WHEN u.column_name IS NULL
           THEN CONCAT('ADD COLUMN "', t.column_name, '" ', t.column_type)
           -- schema has additional column, remove it
           WHEN t.column_name IS NULL
           THEN CONCAT('DROP COLUMN "', t.column_name, '"')
           -- schema has different type, same nullability
           WHEN (t.column_type != u.column_type) AND (t.is_nullable = u.is_nullable)
           THEN CONCAT('ALTER COLUMN "', t.column_name, '" SET TYPE ', t.column_type)
           -- schema has same type, different nullability
           WHEN (t.column_type = u.column_type)  AND (t.is_nullable != u.is_nullable)
           THEN CONCAT('ALTER COLUMN "', t.column_name, '" ', CASE WHEN t.is_nullable THEN 'DROP NOT NULL' ELSE 'SET NOT NULL' END)
           -- schema has different type and different nullability
           ELSE CONCAT('ALTER COLUMN "', t.column_name, '" SET TYPE ' t.column2_type, '; ALTER COLUMN "', t.column_name, '" ', CASE WHEN t.is_nullable THEN 'DROP NOT NULL' ELSE 'SET NOT NULL' END)
         END
       ) stmt
  FROM (
    SELECT column1 AS column_name
          ,column2 AS column_type
          ,column3 AS is_nullable
      FROM (
        VALUES
               ('id',          'UUID',    FALSE)
              ,('first_name',  'VARCHAR', FALSE)
              ,('middle_name', 'VARCHAR', FALSE)
              ,('last_name',   'VARCHAR', FALSE)
      ) c
  ) t
  FULL JOIN (
    SELECT column_name
          ,CONCAT(
            UPPER(CASE WHEN data_type = 'character varying' THEN 'varchar' ELSE data_type END)
           ,CASE WHEN data_type = 'numeric' AND numeric_precision IS NOT NULL THEN
              CONCAT('(', numeric_precision, CASE WHEN numeric_scale > 0 THEN CONCAT(',', numeric_scale) END, ')')
            END
          ) column_type
          ,is_nullable
      FROM information_schema.columns c
     WHERE c.table_schema = 'myapp'
       AND c.table_name   = 'customer'
  ) u
  ON u.column_name = t.column_name;

-- PK
SELECT 'ALTER TABLE myapp.customer ADD CONSTRAINT customer_pk PRIMARY KEY (rel_id)'
 WHERE NOT EXISTS (
   SELECT
     FROM information_schema.table_constraints
    WHERE table_schema    = 'myapp'
      AND table_name      = 'customer'
      AND constraint_name = 'customer_pk'
 )
\gexec

-- UK
SELECT 'ALTER TABLE myapp.customer ADD CONSTRAINT customer_uk_id UNIQUE (id)'
 WHERE NOT EXISTS (
   SELECT
     FROM information_schema.table_constraints
    WHERE table_schema    = 'myapp'
      AND table_name      = 'customer'
      AND constraint_name = 'customer_uk_id'
 )
\gexec

-- Statement level trigger to insert descriptor in shared_info table
CREATE OR REPLACE TRIGGER myapp_customer_search_insert_tg
  AFTER INSERT ON myapp.customer
  REFERENCING NEW TABLE AS NEW_TABLE
  FOR EACH STATEMENT EXECUTE FUNCTION myapp_update_customer();

-- Statement level trigger to update descriptor in shared_info table
CREATE OR REPLACE TRIGGER myapp_customer_search_update_tg
  AFTER UPDATE ON myapp.customer
  REFERENCING NEW TABLE AS NEW_TABLE
  FOR EACH STATEMENT EXECUTE FUNCTION myapp_update_customer();

-- Statement level trigger to delete descriptor from shared_info table
CREATE OR REPLACE TRIGGER myapp_customer_search_delete_tg
  AFTER DELETE ON myapp.customer
  REFERENCING OLD TABLE AS OLD_TABLE
  FOR EACH STATEMENT EXECUTE FUNCTION myapp_delete_shared_info();

--
-- Book
--
-- rel_id     PK
-- descriptor full text description
-- id         UUID
-- name       name of the book
-- author     author of the book
-- theyear    year of the book
-- pages      number of pages in the book
-- isbn       ISBN number of the book
--
-- PK: rel_id
-- UK: id
--

-- Table
CREATE TABLE IF NOT EXISTS myapp.book(
  rel_id  INTEGER GENERATED ALWAYS AS IDENTITY,
  id      UUID    NOT NULL,
  name    VARCHAR NOT NULL,
  author  VARCHAR NOT NULL,
  theyear INTEGER NOT NULL,
  pages   INTEGER NOT NULL,
  isbn    VARCHAR NOT NULL
);

-- PK
SELECT 'ALTER TABLE myapp.book ADD CONSTRAINT book_pk PRIMARY KEY (rel_id)'
 WHERE NOT EXISTS (
   SELECT
     FROM information_schema.table_constraints
    WHERE table_schema    = 'myapp'
      AND table_name      = 'book'
      AND constraint_name = 'book_pk'
 )
\gexec

-- UK
SELECT 'ALTER TABLE myapp.book ADD CONSTRAINT book_uk_id UNIQUE (id)'
 WHERE NOT EXISTS (
   SELECT
     FROM information_schema.table_constraints
    WHERE table_schema    = 'myapp'
      AND table_name      = 'book'
      AND constraint_name = 'book_uk_id'
 )
\gexec

-- Statement level trigger to insert descriptor in shared_info table
CREATE OR REPLACE TRIGGER myapp_book_search_insert_tg
  AFTER INSERT ON myapp.book
  REFERENCING NEW TABLE AS NEW_TABLE
  FOR EACH STATEMENT EXECUTE FUNCTION myapp_update_book();

-- Statement level trigger to update descriptor in shared_info table
CREATE OR REPLACE TRIGGER myapp_book_search_update_tg
  AFTER UPDATE ON myapp.book
  REFERENCING NEW TABLE AS NEW_TABLE
  FOR EACH STATEMENT EXECUTE FUNCTION myapp_update_book();

-- Statement level trigger to delete descriptor from shared_info table
CREATE OR REPLACE TRIGGER myapp_book_search_delete_tg
  AFTER DELETE ON myapp.book
  REFERENCING OLD TABLE AS OLD_TABLE
  FOR EACH STATEMENT EXECUTE FUNCTION myapp_delete_shared_info();

--
-- Movie
--
-- rel_id     PK
-- descriptor full text description
-- id         UUID
-- name       name of the movie
-- director   director of the movie
-- theyear    year of the movie
-- duration   length of the movie
-- imdb       IMDB number of the movie
--
-- PK: rel_id
-- UK: id
--

-- Table
CREATE TABLE IF NOT EXISTS myapp.movie(
  rel_id   INTEGER  GENERATED ALWAYS AS IDENTITY,
  id       UUID     NOT NULL,
  name     VARCHAR  NOT NULL,
  director VARCHAR  NOT NULL,
  theyear  INTEGER  NOT NULL,
  duration INTERVAL NOT NULL,
  imdb     VARCHAR  NOT NULL
);

-- PK
SELECT 'ALTER TABLE myapp.movie ADD CONSTRAINT movie_pk PRIMARY KEY (rel_id)'
 WHERE NOT EXISTS (
   SELECT
     FROM information_schema.table_constraints
    WHERE table_schema    = 'myapp'
      AND table_name      = 'movie'
      AND constraint_name = 'movie_pk'
 )
\gexec

-- UK
SELECT 'ALTER TABLE myapp.movie ADD CONSTRAINT movie_uk_id UNIQUE (id)'
 WHERE NOT EXISTS (
   SELECT
     FROM information_schema.table_constraints
    WHERE table_schema    = 'myapp'
      AND table_name      = 'movie'
      AND constraint_name = 'movie_uk_id'
 )
\gexec

-- Statement level trigger to insert descriptor in shared_info table
CREATE OR REPLACE TRIGGER myapp_movie_search_insert_tg
  AFTER INSERT ON myapp.movie
  REFERENCING NEW TABLE AS NEW_TABLE
  FOR EACH STATEMENT EXECUTE FUNCTION myapp_update_movie();

-- Statement level trigger to update descriptor in shared_info table
CREATE OR REPLACE TRIGGER myapp_movie_search_update_tg
  AFTER UPDATE ON myapp.movie
  REFERENCING NEW TABLE AS NEW_TABLE
  FOR EACH STATEMENT EXECUTE FUNCTION myapp_update_movie();

-- Statement level trigger to delete descriptor from shared_info table
CREATE OR REPLACE TRIGGER myapp_movie_search_delete_tg
  AFTER DELETE ON myapp.movie
  REFERENCING OLD TABLE AS OLD_TABLE
  FOR EACH STATEMENT EXECUTE FUNCTION myapp_delete_shared_info();

--
-- Invoice
--
-- rel_id          PK
-- descriptor      full text description
-- id              UUID
-- customer_rel_id primary key of the customer who paid for the invoice
-- purchased_on    date of ther invoice
-- invoice_number  number of the invoice
--
-- PK: rel_id
-- UK: id
-- FK: customer_rel_id
--

CREATE TABLE IF NOT EXISTS myapp.invoice(
  rel_id          INTEGER                     GENERATED ALWAYS AS IDENTITY,
  id              UUID                        NOT NULL,
  customer_rel_id INTEGER                     NOT NULL,
  purchased_on    TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  invoice_number  VARCHAR                     NOT NULL
);

-- PK
SELECT 'ALTER TABLE myapp.invoice ADD CONSTRAINT invoice_pk PRIMARY KEY (rel_id)'
 WHERE NOT EXISTS (
   SELECT
     FROM information_schema.table_constraints
    WHERE table_schema    = 'myapp'
      AND table_name      = 'invoice'
      AND constraint_name = 'invoice_pk'
 )
\gexec

-- UK
SELECT 'ALTER TABLE myapp.invoice ADD CONSTRAINT invoice_uk_id UNIQUE (id)'
 WHERE NOT EXISTS (
   SELECT
     FROM information_schema.table_constraints
    WHERE table_schema    = 'myapp'
      AND table_name      = 'invoice'
      AND constraint_name = 'invoice_uk_id'
 )
\gexec

-- FK
SELECT 'ALTER TABLE myapp.invoice ADD CONSTRAINT invoice_fk_customer_rel_id FOREIGN KEY(customer_rel_id) REFERENCES myapp.customer(rel_id)'
 WHERE NOT EXISTS (
   SELECT
     FROM information_schema.table_constraints
    WHERE table_schema    = 'myapp'
      AND table_name      = 'invoice'
      AND constraint_name = 'invoice_fk_customer_rel_id'
 )
\gexec

-- Statement level trigger to insert descriptor in shared_info table
CREATE OR REPLACE TRIGGER myapp_invoice_search_insert_tg
  AFTER INSERT ON myapp.invoice
  REFERENCING NEW TABLE AS NEW_TABLE
  FOR EACH STATEMENT EXECUTE FUNCTION myapp_update_invoice();

-- Statement level trigger to update descriptor in shared_info table
CREATE OR REPLACE TRIGGER myapp_invoice_search_update_tg
  AFTER UPDATE ON myapp.invoice
  REFERENCING NEW TABLE AS NEW_TABLE
  FOR EACH STATEMENT EXECUTE FUNCTION myapp_update_invoice();

-- Statement level trigger to delete descriptor from shared_info table
CREATE OR REPLACE TRIGGER myapp_invoice_search_delete_tg
  AFTER DELETE ON myapp.invoice
  REFERENCING OLD TABLE AS OLD_TABLE
  FOR EACH STATEMENT EXECUTE FUNCTION myapp_delete_shared_info();

--
-- Invoice Line
--
-- rel_id         PK
-- descriptor     full text description
-- id             UUID
-- product_oid    object id of the product table
-- product_rel_id PK of the product
-- invoice_rel_id primary key of the invoice containing the line
-- quantity       quantity of items of this type being purchased
-- price          cost of a single item of this type
-- extended       rounded quantity * price value
--
-- PK: rel_id
-- UK: id
-- FK: invoice_rel_id
--

CREATE TABLE IF NOT EXISTS myapp.invoice_line(
  rel_id         INTEGER      GENERATED ALWAYS AS IDENTITY,
  id             UUID         NOT NULL,
  invoice_rel_id INTEGER      NOT NULL,
  product_oid    OID          NOT NULL,
  product_rel_id INTEGER      NOT NULL,
  line           INTEGER      NOT NULL,
  quantity       INTEGER      NOT NULL,
  price          DECIMAL(8,2) NOT NULL,
  extended       DECIMAL(8,2) GENERATED ALWAYS AS (quantity * price) STORED
);

-- PK
SELECT 'ALTER TABLE myapp.invoice_line ADD CONSTRAINT invoice_line_pk PRIMARY KEY (rel_id)'
 WHERE NOT EXISTS (
   SELECT
     FROM information_schema.table_constraints
    WHERE table_schema    = 'myapp'
      AND table_name      = 'invoice_line'
      AND constraint_name = 'invoice_line_pk'
 )
\gexec

-- UK
SELECT 'ALTER TABLE myapp.invoice_line ADD CONSTRAINT invoice_line_uk_id UNIQUE (id)'
 WHERE NOT EXISTS (
   SELECT
     FROM information_schema.table_constraints
    WHERE table_schema    = 'myapp'
      AND table_name      = 'invoice_line'
      AND constraint_name = 'invoice_line_uk_id'
 )
\gexec

-- FK
SELECT 'ALTER TABLE myapp.invoice_line ADD CONSTRAINT invoice_line_fk_invoice_rel_id FOREIGN KEY(invoice_rel_id) REFERENCES myapp.invoice(rel_id)'
 WHERE NOT EXISTS (
   SELECT
     FROM information_schema.table_constraints
    WHERE table_schema    = 'myapp'
      AND table_name      = 'invoice_line'
      AND constraint_name = 'invoice_line_fk_invoice_rel_id'
 )
\gexec

-- CK
SELECT 'ALTER TABLE myapp.invoice_line ADD CONSTRAINT invoice_line_ck_product_oid CHECK (product_oid IN (''myapp.book''::regclass::oid, ''myapp.movie''::regclass::oid))'
 WHERE NOT EXISTS (
   SELECT
     FROM information_schema.table_constraints
    WHERE table_schema    = 'myapp'
      AND table_name      = 'invoice_line'
      AND constraint_name = 'invoice_line_ck_product_oid'
 )
\gexec

-- Statement level trigger to insert descriptor in shared_info table
CREATE OR REPLACE TRIGGER myapp_invoice_line_search_insert_tg
  AFTER INSERT ON myapp.invoice_line
  REFERENCING NEW TABLE AS NEW_TABLE
  FOR EACH STATEMENT EXECUTE FUNCTION myapp_update_invoice_line();

-- Statement level trigger to update descriptor in shared_info table
CREATE OR REPLACE TRIGGER myapp_invoice_line_search_update_tg
  AFTER UPDATE ON myapp.invoice_line
  REFERENCING NEW TABLE AS NEW_TABLE
  FOR EACH STATEMENT EXECUTE FUNCTION myapp_update_invoice_line();

-- Statement level trigger to delete descriptor from shared_info table
CREATE OR REPLACE TRIGGER myapp_invoice_line_search_delete_tg
  AFTER DELETE ON myapp.invoice_line
  REFERENCING OLD TABLE AS OLD_TABLE
  FOR EACH STATEMENT EXECUTE FUNCTION myapp_delete_shared_info();
