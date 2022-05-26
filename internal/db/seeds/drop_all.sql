--
-- Drops ALL data and tables! Only use for development and testing purposes!
-- An attempt to run this seed script in non-development mode will return an error.
--
BEGIN;

DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO public;
COMMENT ON SCHEMA public IS 'standard public schema';

COMMIT;
