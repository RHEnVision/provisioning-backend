BEGIN;

-- Unicode-safe empty check
CREATE OR REPLACE FUNCTION empty(t TEXT)
  RETURNS BOOLEAN AS
$empty$
BEGIN
  RETURN t ~ '^[[:space:]]*$';
END;
$empty$ LANGUAGE 'plpgsql';

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

COMMIT;
