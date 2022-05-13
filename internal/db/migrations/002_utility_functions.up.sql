BEGIN;

CREATE OR REPLACE FUNCTION empty(t TEXT)
    RETURNS BOOLEAN AS
$empty$
BEGIN
    RETURN t ~ '^[[:space:]]*$';
END;
$empty$ LANGUAGE 'plpgsql';

COMMIT;
