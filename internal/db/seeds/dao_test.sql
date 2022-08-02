--
-- Testing seed data, to execute this file during migration set DB_SEED variable to "dao_test".
--
BEGIN;

INSERT INTO accounts(id, account_number, org_id)
VALUES (1, '1', '1')
  ON CONFLICT DO NOTHING;

COMMIT;
