--
-- Testing seed data, to execute this file during migration set DB_SEED variable to "dao_account".
--
BEGIN;

INSERT INTO accounts(id, account_number, org_id)
VALUES (2, '10', '000010')
  ON CONFLICT DO NOTHING;

COMMIT;
