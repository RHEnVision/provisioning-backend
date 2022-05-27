BEGIN;

-- Unicode-safe empty check
CREATE OR REPLACE FUNCTION empty(t TEXT)
  RETURNS BOOLEAN AS
$empty$
BEGIN
  RETURN t ~ '^[[:space:]]*$';
END;
$empty$ LANGUAGE 'plpgsql';

-- Provider constant check
CREATE OR REPLACE FUNCTION valid_provider(i INTEGER)
  RETURNS BOOLEAN AS
$valid_provider$
BEGIN
  RETURN i BETWEEN 1 AND 3;
END;
$valid_provider$ LANGUAGE 'plpgsql';

-- Reset all sequences to the maximum value, works on empty tables too
CREATE OR REPLACE FUNCTION reset_sequences()
  RETURNS void AS
$reset_sequences$
DECLARE
  tn text;
BEGIN
  FOR tn IN SELECT table_name
            FROM information_schema.tables
            WHERE table_schema = 'public'
              AND table_type = 'BASE TABLE'
              AND table_name != 'schema_migrations'
    LOOP
      EXECUTE format(
        'SELECT setval(pg_get_serial_sequence(''"%s"'', ''id''), (SELECT COALESCE(MAX("id"), 1) from "%s"))', tn, tn);
    END LOOP;
END ;
$reset_sequences$ LANGUAGE 'plpgsql';

-- Resource tags must never be changed, this function allows triggers to enforce it
CREATE OR REPLACE FUNCTION prevent_tag_update()
  RETURNS trigger AS
$prevent_tag_update$
BEGIN
  IF OLD.tag != '' AND OLD.tag != NEW.TAG THEN
    RAISE EXCEPTION 'tag is read-only';
  END IF;
  RETURN NEW;
END;
$prevent_tag_update$ LANGUAGE 'plpgsql';

COMMIT;
