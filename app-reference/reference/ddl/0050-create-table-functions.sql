\c mydb

-- trigger functions

-- Statement level trigger for delete events to delete from shared_info table.
-- Any table this trigger is applied to has to specify that the name of the table of data being deleted is OLD_TABLE.
CREATE OR REPLACE FUNCTION myapp_delete_shared_info() RETURNS TRIGGER AS
$$
BEGIN
  DELETE FROM myapp.shared_info msr
  -- TG_RELID is pre-defined var containing object id of table the trigger is applied to
   WHERE msr.tbl_oid = TG_RELID
     AND msr.rel_id IN (SELECT rel_id FROM OLD_TABLE);
  RETURN NULL;
END
$$ LANGUAGE plpgsql;

-- Statement level trigger for customer insert and update events to generate description and search terms in shared_info table
CREATE OR REPLACE FUNCTION myapp_update_customer() RETURNS TRIGGER AS
$$
BEGIN
    INSERT
      INTO myapp.shared_info(
             tbl_oid
            ,rel_id
            ,description
            ,search_terms
           )
    SELECT
           TG_RELID
          ,NEW_TABLE.rel_id
          ,CONCAT(NEW_TABLE.first_name, COALESCE(CONCAT(' ', NEW_TABLE.middle_name), ''), ' ', NEW_TABLE.last_name)
          ,TO_TSVECTOR(CONCAT(NEW_TABLE.first_name, ' ', NEW_TABLE.last_name))
      FROM NEW_TABLE
        ON CONFLICT ON CONSTRAINT shared_info_pk DO
    UPDATE SET
           description  = (SELECT CONCAT(NEW_TABLE.first_name, COALESCE(CONCAT(' ', NEW_TABLE.middle_name), ''), ' ', NEW_TABLE.last_name)
                             FROM NEW_TABLE
                            WHERE NEW_TABLE.rel_id = EXCLUDED.rel_id)
          ,search_terms = (SELECT TO_TSVECTOR(CONCAT(NEW_TABLE.first_name, ' ', NEW_TABLE.last_name))
                             FROM NEW_TABLE
                            WHERE NEW_TABLE.rel_id = EXCLUDED.rel_id);

  RETURN NULL;
END
$$ LANGUAGE plpgsql;


-- Statement level trigger for book insert and update events to generate description and search terms in shared_info table
CREATE OR REPLACE FUNCTION myapp_update_book() RETURNS TRIGGER AS
$$
BEGIN
    INSERT
      INTO myapp.shared_info(
             tbl_oid
            ,rel_id
            ,description
            ,search_terms
           )
    SELECT
           TG_RELID
          ,NEW_TABLE.rel_id
          ,CONCAT(NEW_TABLE.name, ' by ', NEW_TABLE.author, ' in ', NEW_TABLE.theyear, ' pp ', NEW_TABLE.pages, ' isbn ', NEW_TABLE.isbn)
          ,TO_TSVECTOR(CONCAT(NEW_TABLE.name, ' ', NEW_TABLE.author, ' ', NEW_TABLE.theyear, ' ', NEW_TABLE. isbn))
      FROM NEW_TABLE
        ON CONFLICT ON CONSTRAINT shared_info_pk DO
    UPDATE SET
           description  = (SELECT CONCAT(NEW_TABLE.name, ' by ', NEW_TABLE.author, ' in ', NEW_TABLE.theyear, ' pp ', NEW_TABLE.pages, ' isbn ', NEW_TABLE.isbn)
                             FROM NEW_TABLE
                            WHERE NEW_TABLE.rel_id = EXCLUDED.rel_id)
          ,search_terms = (SELECT TO_TSVECTOR(CONCAT(NEW_TABLE.name, ' ', NEW_TABLE.author, ' ', NEW_TABLE.theyear, ' ', NEW_TABLE.isbn))
                             FROM NEW_TABLE
                            WHERE NEW_TABLE.rel_id = EXCLUDED.rel_id);

  RETURN NULL;
END
$$ LANGUAGE plpgsql;

-- Statement level trigger for movie insert and update events to generate description and search terms in shared_info table
CREATE OR REPLACE FUNCTION myapp_update_movie() RETURNS TRIGGER AS
$$
BEGIN
    INSERT
      INTO myapp.shared_info(
             tbl_oid
            ,rel_id
            ,description
            ,search_terms
           )
    SELECT
           TG_RELID
          ,NEW_TABLE.rel_id
          ,CONCAT(NEW_TABLE.name, ' directed by ', NEW_TABLE.director, ' in ', NEW_TABLE.theyear, ' ', NEW_TABLE.duration, ' imdb ', NEW_TABLE.imdb)
          ,TO_TSVECTOR(CONCAT(NEW_TABLE.name, ' ', NEW_TABLE.director, ' ', NEW_TABLE.theyear, ' ', NEW_TABLE.imdb))
      FROM NEW_TABLE
        ON CONFLICT ON CONSTRAINT shared_info_pk DO
    UPDATE SET
           description  = (SELECT CONCAT(NEW_TABLE.name, ' directed by ', NEW_TABLE.director, ' in ', NEW_TABLE.theyear, ' ', duration, ' imdb ', NEW_TABLE.imdb)
                             FROM NEW_TABLE
                            WHERE NEW_TABLE.rel_id = EXCLUDED.rel_id)
          ,search_terms = (SELECT TO_TSVECTOR(CONCAT(NEW_TABLE.name, ' ', NEW_TABLE.director, ' ', NEW_TABLE.theyear, ' ', NEW_TABLE.imdb))
                             FROM NEW_TABLE
                            WHERE NEW_TABLE.rel_id = EXCLUDED.rel_id);

  RETURN NULL;
END
$$ LANGUAGE plpgsql;

-- Statement level trigger for invoice insert and update events to generate description and search terms in shared_info table
CREATE OR REPLACE FUNCTION myapp_update_invoice() RETURNS TRIGGER AS
$$
BEGIN
    INSERT
      INTO myapp.shared_info(
             tbl_oid
            ,rel_id
            ,description
            ,search_terms
           )
    SELECT
           TG_RELID
          ,NEW_TABLE.rel_id
          ,CONCAT(NEW_TABLE.invoice_number, ' purchased on ', NEW_TABLE.purchased_on, ' by ', CONCAT(c.first_name, ' ', c.last_name))
          ,TO_TSVECTOR(CONCAT(NEW_TABLE.invoice_number, ' ', c.first_name, ' ', c.last_name))
      FROM NEW_TABLE
      JOIN myapp.customer c
        ON c.rel_id = NEW_TABLE.customer_rel_id
        ON CONFLICT ON CONSTRAINT shared_info_pk DO
    UPDATE SET
           description  = (SELECT CONCAT(NEW_TABLE.invoice_number, ' purchased on ', NEW_TABLE.purchased_on, ' by ', CONCAT(c.first_name, ' ', c.last_name))
                             FROM NEW_TABLE
                             JOIN myapp.customer c
                               ON c.rel_id = NEW_TABLE.customer_rel_id
                            WHERE NEW_TABLE.rel_id = EXCLUDED.rel_id)
          ,search_terms = (SELECT TO_TSVECTOR(CONCAT(NEW_TABLE.invoice_number, ' ', c.first_name, ' ', c.last_name))
                             FROM NEW_TABLE
                             JOIN myapp.customer c
                               ON c.rel_id = NEW_TABLE.customer_rel_id
                            WHERE NEW_TABLE.rel_id = EXCLUDED.rel_id);

  RETURN NULL;
END
$$ LANGUAGE plpgsql;

-- Statement level trigger for line insert and update events to generate description and search terms in shared_info table
CREATE OR REPLACE FUNCTION myapp_update_invoice_line() RETURNS TRIGGER AS
$$
BEGIN
    INSERT
      INTO myapp.shared_info(
             tbl_oid
            ,rel_id
            ,description
            ,search_terms
           )
    SELECT
           TG_RELID
          ,NEW_TABLE.rel_id
          ,CONCAT(i.invoice_number, ' line ', NEW_TABLE.line, ' of ', NEW_TABLE.quantity, ' ', ss.description) description
          ,TO_TSVECTOR(CONCAT(i.invoice_number, ' ', NEW_TABLE.line, ' ', NEW_TABLE.quantity)) search_terms
      FROM NEW_TABLE
      JOIN myapp.invoice i
        ON i.rel_id = NEW_TABLE.invoice_rel_id
      JOIN myapp.shared_info ss
        ON ss.tbl_oid = NEW_TABLE.product_oid
       AND ss.rel_id = NEW_TABLE.product_rel_id
        ON CONFLICT ON CONSTRAINT shared_info_pk DO
    UPDATE SET
           description  = (SELECT CONCAT(i.invoice_number, ' line ', NEW_TABLE.line, ' of ', NEW_TABLE.quantity, ' ', ss.description)
                             FROM NEW_TABLE
                             JOIN myapp.invoice i
                               ON i.rel_id = NEW_TABLE.invoice_rel_id
                             JOIN myapp.shared_info ss
                               ON ss.tbl_oid = NEW_TABLE.product_oid
                              AND ss.rel_id = NEW_TABLE.product_rel_id
                            WHERE NEW_TABLE.rel_id = EXCLUDED.rel_id)
          ,search_terms = (SELECT TO_TSVECTOR(CONCAT(i.invoice_number, ' ', NEW_TABLE.line, ' ', NEW_TABLE.quantity))
                             FROM NEW_TABLE
                             JOIN myapp.invoice i
                               ON i.rel_id = NEW_TABLE.invoice_rel_id
                             JOIN myapp.shared_info ss
                               ON ss.tbl_oid = NEW_TABLE.product_oid
                              AND ss.rel_id = NEW_TABLE.product_rel_id
                            WHERE NEW_TABLE.rel_id = EXCLUDED.rel_id);

  RETURN NULL;
END
$$ LANGUAGE plpgsql;
