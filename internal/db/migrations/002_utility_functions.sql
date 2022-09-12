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
  RETURN i BETWEEN 1 AND 4;
END;
$valid_provider$ LANGUAGE 'plpgsql';

-- A global constant functions for each provider
CREATE OR REPLACE FUNCTION provider_type_noop()
  RETURNS INTEGER AS
$provider_type_noop$
BEGIN
  RETURN(SELECT 1);
END;
$provider_type_noop$ LANGUAGE 'plpgsql' IMMUTABLE PARALLEL SAFE;

CREATE OR REPLACE FUNCTION provider_type_aws()
  RETURNS INTEGER AS
$provider_type_aws$
BEGIN
  RETURN(SELECT 2);
END;
$provider_type_aws$ LANGUAGE 'plpgsql' IMMUTABLE PARALLEL SAFE;

CREATE OR REPLACE FUNCTION provider_type_azure()
  RETURNS INTEGER AS
$provider_type_azure$
BEGIN
  RETURN(SELECT 3);
END;
$provider_type_azure$ LANGUAGE 'plpgsql' IMMUTABLE PARALLEL SAFE;

CREATE OR REPLACE FUNCTION provider_type_gcp()
  RETURNS INTEGER AS
$provider_type_gcp$
BEGIN
  RETURN(SELECT 4);
END;
$provider_type_gcp$ LANGUAGE 'plpgsql' IMMUTABLE PARALLEL SAFE;

-- Reset all sequences to the maximum value, works on empty tables too
CREATE OR REPLACE FUNCTION reset_sequences(schema TEXT)
  RETURNS void AS
$reset_sequences$
DECLARE
  tn text;
BEGIN
  FOR tn IN SELECT table_name
            FROM information_schema.tables
            WHERE table_schema = schema
              AND table_type = 'BASE TABLE'
              AND table_name !~* '^(schema_version|jobs|job_dependencies|heartbeats)$'
              AND table_name !~* '_?reservation_(details|instances)$'
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
