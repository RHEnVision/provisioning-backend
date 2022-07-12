--
-- Testing seed data, to execute this file during migration set DB_SEED variable to "dao_pubkey".
--
BEGIN;

INSERT INTO accounts(id, account_number, org_id)
VALUES (1, '13', '000013')
  ON CONFLICT DO NOTHING;

COMMIT;
