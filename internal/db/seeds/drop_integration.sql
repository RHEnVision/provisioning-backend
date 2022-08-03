--
-- Drops ALL data and tables in integration schema. Used in tests.
--
BEGIN;

DROP SCHEMA IF EXISTS integration CASCADE;
CREATE SCHEMA integration;
GRANT ALL ON SCHEMA integration TO postgres;
GRANT ALL ON SCHEMA integration TO public;
COMMENT ON SCHEMA integration IS 'integration tests schema';

COMMIT;
