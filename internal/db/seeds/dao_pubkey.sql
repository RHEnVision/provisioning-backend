--
-- Testing seed data, to execute this file during migration set DB_SEED variable to "dao_pubkey".
--
BEGIN;

INSERT INTO pubkeys(id, account_id, name, body)
VALUES (1, 1, 'firstkey', 'sha-rsa body')
  ON CONFLICT DO NOTHING;

COMMIT;
