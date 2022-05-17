--
-- Drops ALL data and tables! Only use for development and testing purposes!
-- An attempt to run this seed script in non-development mode will return an error.
--
BEGIN;

DROP TABLE IF EXISTS
  accounts,
  pubkeys,
  schema_migrations_history,
  schema_migrations CASCADE;

COMMIT;
